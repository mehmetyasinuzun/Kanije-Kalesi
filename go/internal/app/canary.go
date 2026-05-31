package app

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/kanije-kalesi/kanije/internal/fileaudit"
)

// ignoredCanaryProcs are background processes that routinely touch files (AV
// scans, search indexing, thumbnails, OS services). They must NOT trip the
// honeypot, or it fires constantly on its own. Lower-case exe base names.
var ignoredCanaryProcs = map[string]bool{
	"msmpeng.exe": true, "mpdefendercoreservice.exe": true, // Defender
	"searchindexer.exe": true, "searchprotocolhost.exe": true, "searchfilterhost.exe": true,
	"svchost.exe": true, "dllhost.exe": true, "taskhostw.exe": true, "runtimebroker.exe": true,
	"explorer.exe": true, // thumbnail/preview — too noisy to treat as intrusion
	"system":       true, "smss.exe": true, "csrss.exe": true, "wininit.exe": true,
	"services.exe": true, "lsass.exe": true, "fontdrvhost.exe": true, "dwm.exe": true,
	"sihost.exe": true, "ctfmon.exe": true,
}

func isIgnoredCanaryProc(procPath string) bool {
	base := strings.ToLower(filepath.Base(strings.TrimSpace(procPath)))
	return base == "" || ignoredCanaryProcs[base]
}

// canaryWatch polls the honeypot folder's SACL access log; when a NEW access by a
// NON-system process appears (a human or unknown tool opened/copied a decoy), it
// fires the configured action. Background/AV/indexer/thumbnail access is ignored
// so the trap doesn't trigger on its own. Windows-only; a no-op elsewhere.
func (a *App) canaryWatch(ctx context.Context) error {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	var lastSeen string
	baselined := false

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			p := a.cfg.ProtectionPolicy()
			if !p.Enabled || !p.CanaryEnabled || p.CanaryPath == "" {
				baselined = false
				continue
			}

			evs, err := fileaudit.RecentAccess([]string{p.CanaryPath}, 10)
			if err != nil || len(evs) == 0 {
				continue
			}

			// Find the most recent access by a non-ignored process.
			var hit *fileaudit.AccessEvent
			for i := range evs {
				if !isIgnoredCanaryProc(evs[i].Process) {
					hit = &evs[i]
					break
				}
			}
			if hit == nil {
				continue // only background/AV touched it — not an intrusion
			}

			sig := hit.Time + "|" + hit.Object + "|" + hit.Process
			if !baselined {
				lastSeen = sig // pre-existing accesses must not trigger
				baselined = true
				continue
			}
			if sig != lastSeen {
				lastSeen = sig
				proc := filepath.Base(hit.Process)
				reason := "Tuzağa erişildi: " + hit.Object +
					" · erişen: " + proc
				if hit.User != "" {
					reason += " (kullanıcı: " + hit.User + ")"
				}
				a.executeProtection(p.CanaryAction, reason)
			}
		}
	}
}
