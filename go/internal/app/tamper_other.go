//go:build !windows

package app

// scheduledTaskState is Windows-specific (Task Scheduler). On other platforms the
// auto-start mechanism differs (systemd, launchd), so the watchdog skips this
// probe by reporting an indeterminate state.
func scheduledTaskState() taskState { return taskUnknown }
