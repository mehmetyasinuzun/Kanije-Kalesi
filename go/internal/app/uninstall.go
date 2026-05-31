package app

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

// removalPlan is the full footprint to erase on /kaldir — every file, directory,
// auto-start entry and leftover the agent ever creates. Built by removalPaths and
// executed by the platform-specific systemUninstall, which runs detached so it can
// delete the (now-exited) executable and its own install directory.
type removalPlan struct {
	Exe   string   // the running executable
	Files []string // config, DB (+wal/shm/journal), log, version + run markers, old-version leftovers
	Dirs  []string // capture output directories
	Task  string   // Windows scheduled-task name to delete
}

// removalPaths gathers everything the agent has written so /kaldir leaves no
// trace. It deliberately enumerates known artifacts (never a blind directory
// wipe) so it can never delete unrelated user files.
func (a *App) removalPaths() removalPlan {
	exe, _ := os.Executable()

	var files []string
	add := func(p string) {
		if p != "" {
			files = append(files, p)
		}
	}

	cfgPath := a.cfg.FilePath()
	add(cfgPath)

	db := a.cfg.Storage.DBPath
	add(db)
	if db != "" {
		// SQLite side files (WAL mode + rollback journal).
		add(db + "-wal")
		add(db + "-shm")
		add(db + "-journal")
	}

	add(a.cfg.Logging.File)

	// Markers and any older-version leftovers next to the config / executable.
	if dir := filepath.Dir(cfgPath); dir != "" && dir != "." {
		add(filepath.Join(dir, "version.txt"))
		add(filepath.Join(dir, runMarkerName))
	}
	if exe != "" {
		// The updater stages the next binary as exe+".new" and may leave .old/.bak.
		add(exe + ".new")
		add(exe + ".old")
		add(exe + ".bak")
	}

	return removalPlan{Exe: exe, Files: files, Dirs: a.captureDirs(), Task: scheduledTaskName}
}

// captureDirs returns the configured local capture output directories (camera and
// screenshot), de-duplicated. Shared by /kaldir and /imha.
func (a *App) captureDirs() []string {
	var dirs []string
	if d := a.cfg.Camera.LocalPath; d != "" {
		dirs = append(dirs, d)
	}
	if d := a.cfg.Screenshot.LocalPath; d != "" && d != a.cfg.Camera.LocalPath {
		dirs = append(dirs, d)
	}
	return dirs
}

// uninstall erases the agent's entire footprint and then shuts down so the
// detached helper can finish (delete the now-unlocked exe and its directory).
// Returns an error only if the helper could not be launched.
func (a *App) uninstall(_ context.Context) error {
	if err := systemUninstall(a.removalPaths()); err != nil {
		return err
	}
	a.updating.Store(true) // suppress the "shutting down" notification on the way out
	go func() {
		time.Sleep(3 * time.Second) // let the confirmation reply send first
		if a.cancel != nil {
			a.cancel()
		}
	}()
	return nil
}
