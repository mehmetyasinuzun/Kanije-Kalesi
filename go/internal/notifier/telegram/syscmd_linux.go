//go:build linux || darwin

package telegram

import "os/exec"

func systemRestart() {
	exec.Command("systemctl", "reboot").Start()
}

func systemShutdown() {
	exec.Command("systemctl", "poweroff").Start()
}
