//go:build windows

package sysinfo

import (
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32              = windows.NewLazySystemDLL("kernel32.dll")
	procGetSystemPowerStatus = modkernel32.NewProc("GetSystemPowerStatus")
)

// SYSTEM_POWER_STATUS, as returned by GetSystemPowerStatus.
type systemPowerStatus struct {
	ACLineStatus        byte   // 0=offline, 1=online (AC), 255=unknown
	BatteryFlag         byte   // bit 3 (8)=charging, 128=no battery, 255=unknown
	BatteryLifePercent  byte   // 0-100, 255=unknown
	SystemStatusFlag    byte   // power-saver flag (reserved here)
	BatteryLifeTime     uint32 // seconds remaining on battery, 0xFFFFFFFF=unknown
	BatteryFullLifeTime uint32 // seconds at full charge, 0xFFFFFFFF=unknown
}

const (
	batteryFlagCharging  = 0x08
	batteryFlagNoBattery = 0x80
	powerValueUnknown    = 0xFF
	lifeTimeUnknown      = 0xFFFFFFFF
)

// CollectBattery reads the system power status via the Win32 API (no CGo).
func CollectBattery() Battery {
	b := Battery{Percent: -1, Remaining: -1}

	var sps systemPowerStatus
	r, _, _ := procGetSystemPowerStatus.Call(uintptr(unsafe.Pointer(&sps)))
	if r == 0 {
		return b // call failed — leave as unknown
	}

	b.Present = sps.BatteryFlag != batteryFlagNoBattery && sps.BatteryFlag != powerValueUnknown
	b.Charging = sps.BatteryFlag&batteryFlagCharging != 0
	b.OnAC = sps.ACLineStatus == 1
	if sps.BatteryLifePercent != powerValueUnknown {
		b.Percent = int(sps.BatteryLifePercent)
	}
	if sps.BatteryLifeTime != lifeTimeUnknown {
		b.Remaining = time.Duration(sps.BatteryLifeTime) * time.Second
	}
	return b
}
