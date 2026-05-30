//go:build windows

package telegram

import (
	"os/exec"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

func systemRestart() {
	cmd := exec.Command("shutdown", "/r", "/t", "5")
	sysproc.Hide(cmd)
	cmd.Start()
}

func systemShutdown() {
	cmd := exec.Command("shutdown", "/s", "/t", "5")
	sysproc.Hide(cmd)
	cmd.Start()
}
