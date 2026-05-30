//go:build windows

// Package sysproc provides cross-platform helpers for launching subprocesses
// cleanly. On Windows its main job is to suppress the console window that the OS
// would otherwise flash every time a GUI binary (built with -H=windowsgui)
// spawns a console program such as netsh, ffmpeg or shutdown.
package sysproc

import (
	"os/exec"
	"syscall"
)

// createNoWindow is the CREATE_NO_WINDOW process creation flag: the child runs
// without allocating a console, so no window appears.
const createNoWindow = 0x08000000

// Hide configures cmd so that starting it does not flash a console window.
// Call it after constructing the *exec.Cmd and before Run/Start/Output.
func Hide(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true
	cmd.SysProcAttr.CreationFlags |= createNoWindow
}
