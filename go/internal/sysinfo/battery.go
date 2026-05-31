package sysinfo

import "time"

// Battery is a snapshot of the system battery, for the /pil command.
type Battery struct {
	Present   bool          // false if the machine has no battery (desktop)
	Percent   int           // 0-100, or -1 if unknown
	Charging  bool          // currently charging
	OnAC      bool          // plugged into AC power
	Remaining time.Duration // estimated time left on battery, or -1 if unknown
}
