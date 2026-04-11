//go:build windows

package app

import "golang.org/x/sys/windows"

var user32dll = windows.NewLazySystemDLL("user32.dll")
var procLockWorkStation = user32dll.NewProc("LockWorkStation")

func lockScreen() error {
	ret, _, err := procLockWorkStation.Call()
	if ret == 0 {
		return err
	}
	return nil
}
