package app

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kanije-kalesi/kanije/internal/event"
)

// protectionWipeDelay is the cancelable window before a protection-triggered
// secure wipe actually runs — a safety net against a false trigger.
const protectionWipeDelay = 60 * time.Second

// executeProtection runs a protection action chain in response to a fired trigger.
// It ALWAYS publishes a critical event (which notifies the owner and — per the
// default trigger — snaps a camera photo). Then, per the action string:
//
//	"lock" → lock the screen immediately (buys time)
//	"wipe" → schedule a secure wipe behind a cancelable window
//
// ("alert" needs no extra step — the event above is the alert.)
func (a *App) executeProtection(action, reason string) {
	a.log.Warn("KORUMA TETİKLENDİ", "aksiyon", action, "sebep", reason)

	ev := event.New(event.TypeProtectionFired, "Protection")
	ev.Hostname, _ = os.Hostname()
	ev.Extra = map[string]string{"🛡️ Tetik": reason, "🎬 Aksiyon": actionLabel(action)}

	// Grab evidence NOW, before any lock (a locked screen would blank the
	// screenshot). If a capture fails (no camera / busy / no ffmpeg), say so in
	// the message instead of silently sending nothing.
	a.attachEvidence(&ev)

	a.bus.Publish(ev)

	// "lockdown" → persistent lock mode (lockdownWatch keeps re-locking). It also
	// locks right now. ("lockdown" contains "lock", so the immediate lock below
	// fires too — harmless.)
	if strings.Contains(action, "lockdown") {
		if err := a.cfg.SetLockdown(true); err != nil {
			a.log.Warn("koruma: lockdown ayarlanamadı", "err", err)
		} else {
			a.log.Warn("koruma: LOCKDOWN açıldı", "sebep", reason)
		}
	}

	if strings.Contains(action, "lock") {
		if err := lockScreen(); err != nil {
			a.log.Warn("koruma: ekran kilitlenemedi", "err", err)
		}
	}

	if strings.Contains(action, "wipe") {
		a.scheduleProtectionWipe(reason)
	}
}

// attachEvidence captures a camera photo + screenshot and attaches them to ev.
// On failure it records the reason in ev.Extra (so "no camera" is reported, not
// silently dropped). Runs with its own timeout — independent of any caller ctx.
func (a *App) attachEvidence(ev *event.Event) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if a.camera != nil {
		if data, err := a.camera.Capture(ctx); err == nil && len(data) > 0 {
			ev.Attachments = append(ev.Attachments, event.Attachment{
				Type: event.AttachmentPhoto, Data: data, Caption: "📷 Koruma anı",
			})
		} else if err != nil {
			ev.Extra["📷 Kamera"] = "alınamadı (" + shortErr(err) + ")"
		}
	} else {
		ev.Extra["📷 Kamera"] = "bu derlemede yok"
	}

	if a.screen != nil {
		if data, err := a.screen.Capture(ctx); err == nil && len(data) > 0 {
			ev.Attachments = append(ev.Attachments, event.Attachment{
				Type: event.AttachmentScreenshot, Data: data, Caption: "🖥️ Koruma anı",
			})
		} else if err != nil {
			ev.Extra["🖥️ Ekran"] = "alınamadı (" + shortErr(err) + ")"
		}
	}
}

// shortErr trims an error message to a short, UTF-8-safe snippet for messages.
func shortErr(err error) string {
	s := strings.TrimSpace(err.Error())
	if len(s) > 90 {
		s = strings.ToValidUTF8(s[:90], "") + "…"
	}
	return s
}

// lockdownWatch enforces persistent lock mode: while Lockdown is on it re-locks
// the screen every 2s, so an unlock (even with the correct password) is reversed
// within seconds. Disabled only by /kilit tam kapat.
func (a *App) lockdownWatch(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if a.cfg.LockdownActive() {
				_ = lockScreen() // already locked → no-op; unlocked → re-lock
			}
		}
	}
}

// scheduleProtectionWipe starts a secure wipe after protectionWipeDelay unless
// canceled (/koruma iptal). Only one wipe can be pending at a time.
func (a *App) scheduleProtectionWipe(reason string) {
	if !a.wipePending.CompareAndSwap(false, true) {
		return // a wipe is already pending
	}

	if a.cfg.IsConfigured() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		a.bot.SendMessage(ctx, fmt.Sprintf(
			"⏳ <b>GÜVENLİ SİLME %d sn içinde başlayacak</b>\nSebep: %s\nİptal için: <code>/koruma iptal</code>",
			int(protectionWipeDelay.Seconds()), reason))
		cancel()
	}

	go func() {
		time.Sleep(protectionWipeDelay)
		// CompareAndSwap so a concurrent /koruma iptal wins cleanly.
		if a.wipePending.CompareAndSwap(true, false) {
			a.log.Warn("koruma: güvenli silme yürütülüyor", "sebep", reason)
			_ = a.destroy(context.Background())
		}
	}()
}

// CancelProtectionWipe aborts a pending protection wipe. Returns true if one was
// actually pending. Called by /koruma iptal.
func (a *App) CancelProtectionWipe() bool {
	return a.wipePending.CompareAndSwap(true, false)
}

// checkProtectionTriggers inspects an incoming event and fires the matching
// protection trigger (USB removal / failed-login threshold). Called from
// handleEvent. The dead-man trigger lives in deadManWatch instead.
func (a *App) checkProtectionTriggers(ev event.Event) {
	p := a.cfg.ProtectionPolicy()
	if !p.Enabled {
		return
	}

	switch ev.Type {
	case event.TypeUSBRemoved:
		if p.USBEnabled && usbMatches(p.USBDevice, ev) {
			label := ev.DeviceLabel
			if label == "" {
				label = ev.DeviceName
			}
			reason := "USB aygıtı çıkarıldı"
			if label != "" {
				reason = "USB çıkarıldı: " + label
			}
			a.executeProtection(p.USBAction, reason)
		}

	case event.TypeLoginFailed:
		if p.FailedLoginEnabled {
			n := a.failedLogins.Add(1)
			if int(n) >= p.FailedLoginThreshold {
				a.failedLogins.Store(0)
				a.executeProtection(p.FailedLoginAction,
					fmt.Sprintf("%d ardışık başarısız giriş denemesi", n))
			}
		}

	case event.TypeLoginSuccess, event.TypeScreenUnlock:
		a.failedLogins.Store(0) // a successful login clears the streak
	}
}

// usbMatches reports whether a USB-removal event matches the configured filter
// (empty filter = any USB removal triggers).
func usbMatches(filter string, ev event.Event) bool {
	if strings.TrimSpace(filter) == "" {
		return true
	}
	return strings.EqualFold(filter, ev.DeviceName) || strings.EqualFold(filter, ev.DeviceLabel)
}

// actionLabel renders an action string as a readable Turkish chain. "lockdown"
// supersedes the plain "lock" label (it contains the substring "lock").
func actionLabel(action string) string {
	var parts []string
	switch {
	case strings.Contains(action, "lockdown"):
		parts = append(parts, "Tam Kilit")
	case strings.Contains(action, "lock"):
		parts = append(parts, "Kilitle")
	}
	if strings.Contains(action, "alert") {
		parts = append(parts, "Alarm+Foto")
	}
	if strings.Contains(action, "wipe") {
		parts = append(parts, "Güvenli Sil")
	}
	if len(parts) == 0 {
		return action
	}
	return strings.Join(parts, " + ")
}
