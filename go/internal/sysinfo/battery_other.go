//go:build !windows

package sysinfo

// CollectBattery is implemented via the Win32 power API on Windows. Elsewhere it
// reports an unknown/absent battery.
func CollectBattery() Battery {
	return Battery{Percent: -1, Remaining: -1}
}
