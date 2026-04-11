//go:build windows

package main

import (
	"log/slog"

	"github.com/kanije-kalesi/kanije/internal/listener"
	winlistener "github.com/kanije-kalesi/kanije/internal/listener/windows"
)

func buildListeners(log *slog.Logger) []listener.Listener {
	return []listener.Listener{
		winlistener.NewEventLogListener(log.With("listener", "EventLog")),
		winlistener.NewUSBMonitor(log.With("listener", "USB")),
		winlistener.NewPowerMonitor(log.With("listener", "Power")),
	}
}
