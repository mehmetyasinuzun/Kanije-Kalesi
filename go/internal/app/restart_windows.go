//go:build windows

package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

// systemRestartApp relaunches the agent with its freshly-saved config. A detached
// helper waits for THIS process to fully exit (releasing the single-instance
// lock) before starting it again — preferring the scheduled task, falling back to
// launching the exe directly. Used by /aktar after the owner/token changes.
func systemRestartApp() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	scriptPath := filepath.Join(os.TempDir(), "kanije-restart.ps1")
	if err := os.WriteFile(scriptPath, []byte(buildRestartScript(os.Getpid(), exe)), 0o600); err != nil {
		return err
	}

	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	sysproc.Hide(cmd)
	if err := cmd.Start(); err != nil {
		os.Remove(scriptPath)
		return err
	}
	return nil
}

func buildRestartScript(pid int, exe string) string {
	var b strings.Builder
	b.WriteString("$ErrorActionPreference='SilentlyContinue'\n")
	b.WriteString("$p=" + strconv.Itoa(pid) + "\n")
	b.WriteString("while (Get-Process -Id $p -ErrorAction SilentlyContinue) { Start-Sleep -Milliseconds 300 }\n")
	b.WriteString("Start-Sleep -Milliseconds 800\n")
	b.WriteString("$t=Get-ScheduledTask -TaskName " + psQuote(scheduledTaskName) + " -ErrorAction SilentlyContinue\n")
	b.WriteString("if ($t) { Start-ScheduledTask -TaskName " + psQuote(scheduledTaskName) + " } else { Start-Process -FilePath " + psQuote(exe) + " }\n")
	b.WriteString("Remove-Item -Force -LiteralPath $MyInvocation.MyCommand.Path\n")
	return b.String()
}
