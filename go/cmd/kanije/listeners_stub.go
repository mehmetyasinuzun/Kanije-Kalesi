//go:build !windows && !linux

package main

import (
	"log/slog"

	"github.com/kanije-kalesi/sentinel/internal/listener"
	"github.com/kanije-kalesi/sentinel/internal/listener/stub"
)

func buildListeners(log *slog.Logger) []listener.Listener {
	return []listener.Listener{
		stub.New("EventLog", log),
		stub.New("USBMonitor", log),
		stub.New("PowerMonitor", log),
	}
}
