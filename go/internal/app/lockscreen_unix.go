//go:build !windows

package app

import (
	"fmt"
	"os/exec"
)

func lockScreen() error {
	// Try common Linux screen lock commands in order
	commands := [][]string{
		{"loginctl", "lock-session"},
		{"xdg-screensaver", "lock"},
		{"gnome-screensaver-command", "--lock"},
		{"xscreensaver-command", "-lock"},
		{"qdbus", "org.kde.screensaver", "/ScreenSaver", "Lock"},
	}

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("ekran kilitleme komutu bulunamadı")
}
