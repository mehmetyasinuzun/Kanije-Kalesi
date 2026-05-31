package app

import (
	"context"
	"time"

	"github.com/kanije-kalesi/kanije/internal/fileaudit"
)

// canaryWatch polls the honeypot folder's SACL access log; when a NEW access
// appears (someone opened/copied a decoy file), it fires the configured action —
// catching an intruder reaching for "juicy" files WITHOUT touching their data.
// Windows-only effect (fileaudit); a harmless no-op elsewhere.
func (a *App) canaryWatch(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
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
				baselined = false // re-baseline on next enable
				continue
			}

			evs, err := fileaudit.RecentAccess([]string{p.CanaryPath}, 5)
			if err != nil || len(evs) == 0 {
				continue
			}

			// Signature of the most recent access; a change means a new access.
			sig := evs[0].Time + "|" + evs[0].Object + "|" + evs[0].Process
			if !baselined {
				lastSeen = sig // pre-existing accesses must not trigger
				baselined = true
				continue
			}
			if sig != lastSeen {
				lastSeen = sig
				a.executeProtection(p.CanaryAction, "Tuzak (honeypot) dosyaya erişildi: "+evs[0].Object)
			}
		}
	}
}
