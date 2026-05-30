//go:build !windows

// Package sysproc provides subprocess helpers. On non-Windows platforms hiding
// the console window is a no-op (there is no console-window concept).
package sysproc

import "os/exec"

// Hide is a no-op on non-Windows platforms, which have no console-window concept.
func Hide(cmd *exec.Cmd) {}
