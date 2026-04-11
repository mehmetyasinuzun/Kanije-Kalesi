//go:build linux

package main

import (
	"log/slog"

	"github.com/kanije-kalesi/sentinel/internal/listener"
	linuxlistener "github.com/kanije-kalesi/sentinel/internal/listener/linux"
)

func buildListeners(log *slog.Logger) []listener.Listener {
	return []listener.Listener{
		linuxlistener.NewJournaldListener(log.With("listener", "Journald")),
		linuxlistener.NewUdevListener(log.With("listener", "Udev")),
		linuxlistener.NewLogindListener(log.With("listener", "Logind")),
	}
}
