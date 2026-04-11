//go:build windows

package telegram

import "os/exec"

func systemRestart() {
	exec.Command("shutdown", "/r", "/t", "5").Start()
}

func systemShutdown() {
	exec.Command("shutdown", "/s", "/t", "5").Start()
}
