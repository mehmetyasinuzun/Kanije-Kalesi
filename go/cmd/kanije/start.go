package main

import (
	"log/slog"

	"github.com/kanije-kalesi/sentinel/internal/app"
	"github.com/kanije-kalesi/sentinel/internal/config"
)

// startApp creates and runs the full application.
// Separated to keep main.go clean.
func startApp(cfg *config.Config, log *slog.Logger) error {
	application, err := app.New(cfg, log)
	if err != nil {
		return err
	}
	application.SetListeners(buildListeners(log)...)
	return application.Run()
}
