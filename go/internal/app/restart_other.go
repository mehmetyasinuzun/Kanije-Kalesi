//go:build !windows

package app

import (
	"os"
	"os/exec"
)

// systemRestartApp re-execs the agent on non-Windows platforms. A managed service
// (systemd/launchd) normally restarts it anyway; this is a best-effort fallback.
func systemRestartApp() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	return exec.Command(exe).Start()
}
