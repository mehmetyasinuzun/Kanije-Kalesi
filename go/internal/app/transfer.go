package app

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// transfer hands the device to a new owner: it replaces the access tree with the
// new owner as sole root (all capabilities), updates the saved config's chat ID
// (and bot token, if one was supplied), then restarts the agent so it reconnects
// under the new identity. The old owner and any sub-users lose all access.
func (a *App) transfer(_ context.Context, newOwnerID int64, newToken string) error {
	if newOwnerID == 0 {
		return fmt.Errorf("geçersiz yeni sahip ID")
	}

	// Use a background context so the hand-off finishes even though the caller's
	// request context is short-lived and we're about to restart.
	ctx := context.Background()
	if err := a.acl.TransferRoot(ctx, newOwnerID); err != nil {
		return err
	}
	if err := a.cfg.SetField("telegram.chat_id", strconv.FormatInt(newOwnerID, 10)); err != nil {
		return fmt.Errorf("yeni sahip ID kaydedilemedi: %w", err)
	}
	if newToken != "" {
		if err := a.cfg.SetField("telegram.bot_token", newToken); err != nil {
			return fmt.Errorf("yeni bot token kaydedilemedi: %w", err)
		}
	}

	a.updating.Store(true) // restarting on purpose — don't send the shutdown notice
	go func() {
		time.Sleep(3 * time.Second) // let the confirmation reply send first
		if err := systemRestartApp(); err != nil {
			a.log.Warn("devir sonrası yeniden başlatma başlatılamadı", "err", err)
		}
		if a.cancel != nil {
			a.cancel()
		}
	}()
	return nil
}
