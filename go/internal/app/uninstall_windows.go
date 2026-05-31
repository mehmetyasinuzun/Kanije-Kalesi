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

// systemUninstall launches a detached PowerShell helper that waits for this
// process to exit, deletes the scheduled task, every file/dir in the plan, the
// updater's temp scripts, the executable, and finally its own install directory
// (if empty) and itself — leaving no trace. The caller must shut the app down
// shortly after this returns so the helper can remove the now-unlocked exe.
func systemUninstall(plan removalPlan) error {
	scriptPath := filepath.Join(os.TempDir(), "kanije-uninstall.ps1")
	if err := os.WriteFile(scriptPath, []byte(buildUninstallScript(os.Getpid(), plan)), 0o600); err != nil {
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

// buildUninstallScript renders the self-deleting cleanup helper. It runs after
// the agent (pid) exits so file locks are released.
func buildUninstallScript(pid int, plan removalPlan) string {
	var b strings.Builder
	b.WriteString("$ErrorActionPreference='SilentlyContinue'\n")
	b.WriteString("$p=" + strconv.Itoa(pid) + "\n")
	b.WriteString("while (Get-Process -Id $p -ErrorAction SilentlyContinue) { Start-Sleep -Milliseconds 300 }\n")
	b.WriteString("Start-Sleep -Milliseconds 800\n")

	// Kill the auto-start so a leftover task can never relaunch us.
	if plan.Task != "" {
		b.WriteString("schtasks /delete /tn " + psQuote(plan.Task) + " /f | Out-Null\n")
		b.WriteString("Unregister-ScheduledTask -TaskName " + psQuote(plan.Task) + " -Confirm:$false\n")
	}

	for _, f := range plan.Files {
		b.WriteString("Remove-Item -Force -LiteralPath " + psQuote(f) + "\n")
	}
	for _, d := range plan.Dirs {
		b.WriteString("Remove-Item -Recurse -Force -LiteralPath " + psQuote(d) + "\n")
	}

	// The updater's staged binaries and swap scripts.
	b.WriteString("Remove-Item -Force \"$env:TEMP\\kanije-swap.ps1\"\n")
	b.WriteString("Remove-Item -Force \"$env:TEMP\\kanije-*.ps1\"\n")

	if plan.Exe != "" {
		b.WriteString("Remove-Item -Force -LiteralPath " + psQuote(plan.Exe) + "\n")
		// Remove the install directory too, but only if our exe was the last thing
		// in it — never wipe a shared/user directory that still holds other files.
		dir := filepath.Dir(plan.Exe)
		b.WriteString("$d=" + psQuote(dir) + "\n")
		b.WriteString("if ((Get-ChildItem -Force -LiteralPath $d -ErrorAction SilentlyContinue | Measure-Object).Count -eq 0) { Remove-Item -Force -Recurse -LiteralPath $d }\n")
	}

	// The helper deletes itself last.
	b.WriteString("Remove-Item -Force -LiteralPath $MyInvocation.MyCommand.Path\n")
	return b.String()
}

// psQuote wraps s in a PowerShell single-quoted literal, escaping embedded quotes.
func psQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
