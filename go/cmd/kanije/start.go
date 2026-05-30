package main

import (
	"log/slog"

	"github.com/kanije-kalesi/kanije/internal/app"
	"github.com/kanije-kalesi/kanije/internal/config"
)

// startApp creates and runs the full application.
// Separated to keep main.go clean.
func startApp(cfg *config.Config, log *slog.Logger, version string) error {
	application, err := app.New(cfg, log, version)
	if err != nil {
		return err
	}
	application.SetListeners(buildListeners(log)...)
	return application.Run()
}
