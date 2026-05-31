//go:build windows

// Package clipboard reads the system clipboard. Windows-only; other platforms
// return an unsupported error.
package clipboard

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

// ReadText returns the clipboard's text via PowerShell's Get-Clipboard. Console
// output is forced to UTF-8 so Turkish characters survive, and the console window
// is hidden. This avoids hand-rolled Win32 clipboard/unsafe.Pointer handling —
// CGo-free and vet-clean.
func ReadText() (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
		"[Console]::OutputEncoding=[System.Text.Encoding]::UTF8; Get-Clipboard -Raw")
	sysproc.Hide(cmd)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("pano okunamadı: %w", err)
	}
	return strings.TrimRight(string(out), "\r\n"), nil
}
