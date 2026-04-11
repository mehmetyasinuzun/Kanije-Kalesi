//go:build !windows && !linux

// Package stub provides no-op listener implementations for platforms where
// native monitoring is not yet supported (macOS, FreeBSD, etc.).
// The application starts without error; only Telegram bot commands work.
package stub

import (
	"context"
	"log/slog"

	"github.com/kanije-kalesi/sentinel/internal/event"
)

// NopListener is a Listener that immediately returns nil.
type NopListener struct {
	name string
	log  *slog.Logger
}

func New(name string, log *slog.Logger) *NopListener {
	return &NopListener{name: name, log: log}
}

func (n *NopListener) Name() string { return n.name }

func (n *NopListener) Start(ctx context.Context, bus *event.Bus) error {
	n.log.Warn("listener bu platformda desteklenmiyor — atlanıyor",
		"listener", n.name,
		"platform", "darwin/other")
	<-ctx.Done()
	return nil
}
