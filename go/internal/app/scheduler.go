package app

import (
	"context"
	"time"
)

// runScheduler fires due scheduled tasks every 30s, asking the bot to execute
// each command for the owner. It is a no-op when no store is configured.
func (a *App) runScheduler(ctx context.Context) error {
	if a.schedule == nil {
		<-ctx.Done()
		return nil
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for _, t := range a.schedule.Due(time.Now()) {
				a.bot.ExecScheduled(ctx, t.Command)
			}
		}
	}
}
