//go:build windows

package app

import (
	"os/exec"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

// scheduledTaskState queries the Windows Task Scheduler for the auto-start task.
// It uses the ScheduledTasks PowerShell module whose State enum (Ready, Running,
// Disabled) is NOT localized, so parsing is robust across system languages —
// unlike schtasks.exe text output. The console window is suppressed via Hide.
func scheduledTaskState() taskState {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
		"(Get-ScheduledTask -TaskName "+scheduledTaskName+" -ErrorAction SilentlyContinue).State")
	sysproc.Hide(cmd)

	out, err := cmd.Output()
	if err != nil {
		return taskUnknown
	}
	switch strings.TrimSpace(string(out)) {
	case "":
		return taskMissing
	case "Disabled":
		return taskDisabled
	default: // Ready, Running
		return taskEnabled
	}
}
