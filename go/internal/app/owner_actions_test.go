package app

import (
	"testing"

	"github.com/kanije-kalesi/kanije/internal/config"
)

func contains(s []string, want string) bool {
	for _, v := range s {
		if v == want {
			return true
		}
	}
	return false
}

func testCfg() *config.Config {
	cfg := config.Defaults()
	cfg.SetFilePath("base/config.toml")
	cfg.Storage.DBPath = "base/kanije.db"
	cfg.Logging.File = "base/kanije.log"
	cfg.Camera.LocalPath = "base/cam"
	cfg.Screenshot.LocalPath = "base/shot"
	return cfg
}

func TestRemovalPaths(t *testing.T) {
	a := &App{cfg: testCfg()}
	plan := a.removalPaths()

	if plan.Task != scheduledTaskName {
		t.Errorf("Task = %q, beklenen %q", plan.Task, scheduledTaskName)
	}
	// Config, DB + SQLite side files, and log must all be scheduled for deletion.
	for _, want := range []string{
		"base/config.toml", "base/kanije.db",
		"base/kanije.db-wal", "base/kanije.db-shm", "base/kanije.db-journal",
		"base/kanije.log",
	} {
		if !contains(plan.Files, want) {
			t.Errorf("removal files %q içermiyor: %v", want, plan.Files)
		}
	}
	// Both distinct capture dirs are included.
	if !contains(plan.Dirs, "base/cam") || !contains(plan.Dirs, "base/shot") {
		t.Errorf("capture dizinleri eksik: %v", plan.Dirs)
	}
}

func TestSecureFiles(t *testing.T) {
	a := &App{cfg: testCfg()}
	got := a.secureFiles()

	// Sensitive files (token-bearing config, history DB + side files, log).
	for _, want := range []string{
		"base/config.toml", "base/kanije.db",
		"base/kanije.db-wal", "base/kanije.db-shm", "base/kanije.db-journal",
		"base/kanije.log",
	} {
		if !contains(got, want) {
			t.Errorf("secureFiles %q içermiyor: %v", want, got)
		}
	}
}

func TestCaptureDirsDedup(t *testing.T) {
	cfg := config.Defaults()
	cfg.Camera.LocalPath = "same"
	cfg.Screenshot.LocalPath = "same" // identical → de-duplicated to one
	a := &App{cfg: cfg}
	if dirs := a.captureDirs(); len(dirs) != 1 {
		t.Errorf("aynı dizin tekilleştirilmeli, alınan: %v", dirs)
	}
}
