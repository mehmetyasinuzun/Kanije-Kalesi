//go:build !windows

package sysproc

import "os/exec"

// Hide is a no-op on non-Windows platforms, which have no console-window concept.
func Hide(cmd *exec.Cmd) {}
