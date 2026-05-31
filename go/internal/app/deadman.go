package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const checkInFileName = "checkin.dat"

// checkInPath returns where the last-check-in timestamp is stored (next to the
// config, else the executable).
func (a *App) checkInPath() string {
	dir := filepath.Dir(a.cfg.FilePath())
	if dir == "" || dir == "." {
		if exe, err := os.Executable(); err == nil {
			dir = filepath.Dir(exe)
		}
	}
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, checkInFileName)
}

// loadCheckIn seeds the last check-in time from disk, or "now" on first run.
func (a *App) loadCheckIn() {
	if p := a.checkInPath(); p != "" {
		if b, err := os.ReadFile(p); err == nil {
			if t, err := strconv.ParseInt(strings.TrimSpace(string(b)), 10, 64); err == nil && t > 0 {
				a.lastCheckIn.Store(t)
				return
			}
		}
	}
	a.lastCheckIn.Store(time.Now().Unix())
}

// recordCheckIn marks the owner as present (called on any authorized command) and
// resets the dead-man latch. Persisted so a restart doesn't reset the timer.
func (a *App) recordCheckIn() {
	now := time.Now().Unix()
	a.lastCheckIn.Store(now)
	a.deadManFired.Store(false)
	if p := a.checkInPath(); p != "" {
		_ = os.WriteFile(p, []byte(strconv.FormatInt(now, 10)), 0o600)
	}
}

// LastCheckIn returns the time of the owner's last check-in (for /koruma status).
func (a *App) LastCheckIn() time.Time {
	return time.Unix(a.lastCheckIn.Load(), 0)
}

// deadManWatch fires the dead-man action when no check-in has arrived within the
// configured window. It checks every minute and re-reads config live, so /koruma
// changes take effect without a restart. Fires once per lapse (reset on check-in).
func (a *App) deadManWatch(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			p := a.cfg.ProtectionPolicy()
			if !p.Enabled || !p.DeadManEnabled {
				continue
			}
			elapsed := time.Since(a.LastCheckIn())
			if elapsed >= time.Duration(p.DeadManHours)*time.Hour {
				if a.deadManFired.CompareAndSwap(false, true) {
					a.executeProtection(p.DeadManAction,
						fmt.Sprintf("Dead-man switch — %d saattir check-in yok", p.DeadManHours))
				}
			}
		}
	}
}
