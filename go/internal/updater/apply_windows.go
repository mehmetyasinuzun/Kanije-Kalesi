//go:build windows

package updater

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

// CanSelfInstall reports whether the agent can replace its own binary in place.
func CanSelfInstall() bool { return true }

// SelfInstall downloads (and verifies) the new binary, then spawns a detached
// helper that waits for this process to exit, swaps the executable, and restarts
// the scheduled task. The caller must shut the app down shortly after this
// returns so the helper can replace the now-unlocked exe.
func (u *Updater) SelfInstall(ctx context.Context, rel *Release) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	newPath := exe + ".new"
	if err := u.downloadVerified(ctx, rel, newPath); err != nil {
		return err
	}

	scriptPath := filepath.Join(os.TempDir(), "kanije-swap.ps1")
	if err := os.WriteFile(scriptPath, []byte(buildSwapScript(os.Getpid(), newPath, exe)), 0o600); err != nil {
		os.Remove(newPath)
		return err
	}

	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	sysproc.Hide(cmd)
	if err := cmd.Start(); err != nil {
		os.Remove(newPath)
		os.Remove(scriptPath)
		return err
	}
	u.log.Info("güncelleme yardımcısı başlatıldı", "sürüm", rel.Version)
	return nil
}

// buildSwapScript waits for the agent (pid) to exit, replaces exePath with
// newPath, and restarts the scheduled task.
func buildSwapScript(pid int, newPath, exePath string) string {
	var b strings.Builder
	b.WriteString("$ErrorActionPreference='SilentlyContinue'\n")
	b.WriteString("$p=" + strconv.Itoa(pid) + "\n")
	b.WriteString("while (Get-Process -Id $p -ErrorAction SilentlyContinue) { Start-Sleep -Milliseconds 300 }\n")
	b.WriteString("Start-Sleep -Milliseconds 800\n")
	b.WriteString("Move-Item -Force " + psQuote(newPath) + " " + psQuote(exePath) + "\n")
	b.WriteString("Start-ScheduledTask -TaskName KanijeKalesi\n")
	return b.String()
}

// psQuote wraps s in a PowerShell single-quoted literal, escaping embedded quotes.
func psQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
