package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/kanije-kalesi/kanije/internal/event"
)

// scheduledTaskName is the Windows Task Scheduler entry that auto-starts the
// agent. The anti-tamper watchdog checks it is still present and enabled, and
// /kaldir removes it. Keep in sync with the installer and updater swap script.
const scheduledTaskName = "KanijeKalesi"

// runMarkerName is a sentinel written while the agent runs and removed on a clean
// shutdown. If it already exists at boot, the previous run ended uncleanly — a
// kill, crash or power loss — which for an anti-theft guard is worth reporting.
const runMarkerName = "kanije.run"

// taskState is the health of the auto-start scheduled task.
type taskState int

const (
	taskEnabled  taskState = iota // present and runnable (Ready/Running)
	taskDisabled                  // present but switched off — someone disabled it
	taskMissing                   // not installed (e.g. running manually) — not an alert
	taskUnknown                   // could not determine (non-Windows, query failed)
)

// tamperWatch runs the anti-tamper watchdog until ctx is canceled. It detects an
// unclean previous shutdown at boot, then periodically verifies the agent's
// executable, config and DB are present and (Windows) its scheduled task is still
// enabled. Each distinct problem alerts the owner once — on the transition into
// the bad state — so a persistent issue never spams.
func (a *App) tamperWatch(ctx context.Context) error {
	if !a.cfg.TamperWatchEnabled() {
		<-ctx.Done()
		return nil
	}

	a.checkRunMarker() // unclean-shutdown detection + (re)arm the marker

	exePath, _ := os.Executable()
	// Baseline the running binary's hash at boot. A later mismatch means the
	// on-disk executable was swapped/patched while we run — code injection or a
	// trojanized replacement. (Self-update never trips this: it exits first, then
	// the helper swaps the file, so the new binary re-baselines on next boot.)
	baselineHash := exeHash(exePath)
	alerted := make(map[string]bool)

	a.tamperScan(exePath, baselineHash, alerted) // first pass right away

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			a.clearRunMarker()
			return nil
		case <-ticker.C:
			a.tamperScan(exePath, baselineHash, alerted)
		}
	}
}

// tamperScan performs one round of integrity checks. Each probe fires once, on
// the transition into the bad state (tracked in alerted), and clears silently on
// recovery so the owner is told both when something breaks and isn't re-spammed.
func (a *App) tamperScan(exePath, baselineHash string, alerted map[string]bool) {
	type probe struct {
		key  string
		bad  bool
		desc string
	}
	var probes []probe

	if exePath != "" {
		probes = append(probes, probe{"exe", !fileExists(exePath),
			"Ajan dosyası bulunamadı (silinmiş/taşınmış): " + exePath})

		// Integrity: a changed hash means the binary was replaced/patched at runtime.
		if baselineHash != "" {
			if cur := exeHash(exePath); cur != "" && cur != baselineHash {
				probes = append(probes, probe{"exehash", true,
					"Ajan dosyası DEĞİŞTİRİLDİ — hash uyuşmuyor (olası kod enjeksiyonu / truva'lı kopya): " + exePath})
			} else {
				probes = append(probes, probe{"exehash", false, ""})
			}
		}
	}
	if cf := a.cfg.FilePath(); cf != "" {
		probes = append(probes, probe{"config", !fileExists(cf),
			"Yapılandırma dosyası silinmiş: " + cf})
	}
	if db := a.cfg.Storage.DBPath; db != "" {
		probes = append(probes, probe{"db", !fileExists(db),
			"Veritabanı dosyası silinmiş: " + db})
	}
	// Only "Disabled" is an alert: missing/unknown means it's likely not installed
	// (manual run) rather than someone switching off the auto-start.
	probes = append(probes, probe{"task", scheduledTaskState() == taskDisabled,
		"Otomatik başlatma görevi devre dışı bırakıldı (" + scheduledTaskName + ")"})

	for _, p := range probes {
		switch {
		case p.bad && !alerted[p.key]:
			alerted[p.key] = true
			a.raiseTamper(p.desc)
		case !p.bad:
			alerted[p.key] = false
		}
	}
}

// raiseTamper publishes a critical tamper event (which, per the default trigger,
// also snaps a camera photo of whoever is at the machine).
func (a *App) raiseTamper(detail string) {
	a.log.Warn("KURCALAMA tespit edildi", "detay", detail)
	ev := event.New(event.TypeTamperAlert, "TamperWatch")
	ev.Hostname, _ = os.Hostname()
	ev.Extra = map[string]string{"🛑 Kurcalama": detail}
	a.bus.Publish(ev)
}

// checkRunMarker reports an unclean previous shutdown, then (re)writes the marker
// so the next unclean exit is caught too.
func (a *App) checkRunMarker() {
	path := a.runMarkerPath()
	if path == "" {
		return
	}
	if fileExists(path) {
		a.raiseTamper("Önceki oturum düzgün kapanmadı — süreç sonlandırılmış, çökmüş veya güç kesilmiş olabilir.")
	}
	_ = os.WriteFile(path, []byte(time.Now().UTC().Format(time.RFC3339)), 0o600)
}

// clearRunMarker removes the run marker on a clean shutdown.
func (a *App) clearRunMarker() {
	if path := a.runMarkerPath(); path != "" {
		os.Remove(path)
	}
}

// runMarkerPath returns the run-marker location: next to the config when known,
// else next to the executable.
func (a *App) runMarkerPath() string {
	dir := filepath.Dir(a.cfg.FilePath())
	if dir == "" || dir == "." {
		if exe, err := os.Executable(); err == nil {
			dir = filepath.Dir(exe)
		}
	}
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, runMarkerName)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// exeHash returns the SHA-256 of the file at path (hex), or "" if it can't be
// read. Used to baseline and re-verify the agent's own binary.
func exeHash(path string) string {
	if path == "" {
		return ""
	}
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}
