package config

// Defaults returns a Config with all safe production defaults applied.
// These values are used when the config file is missing or a key is absent.
func Defaults() *Config {
	return &Config{
		Telegram: TelegramConfig{
			SendTimeoutSec: 15,
			RetryCount:     3,
			RetryDelaySec:  5,
		},
		Triggers: defaultTriggers(),
		Camera: CameraConfig{
			FFmpegPath:   "ffmpeg",
			DeviceIndex:  0,
			Width:        640,
			Height:       480,
			WarmupFrames: 5,
			JPEGQuality:  85,
			SaveLocal:    false,
			LocalPath:    "./captures/",
		},
		Screenshot: ScreenshotConfig{
			JPEGQuality: 75,
			SaveLocal:   false,
			LocalPath:   "./captures/",
		},
		Heartbeat: HeartbeatConfig{
			Enabled:       true,
			IntervalHours: 6,
			IncludeUptime: true,
			IncludeDisk:   true,
		},
		Storage: StorageConfig{
			DBPath:             "./kanije.db",
			MaxRecentEvents:    100,
			EventRetentionDays: 30,
		},
		Logging: LoggingConfig{
			Level:         "info",
			File:          "./kanije.log",
			MaxSizeMB:     10,
			BackupCount:   3,
			ConsoleOutput: true,
			JSONFormat:    false,
		},
		Security: SecurityConfig{
			DeleteCapturesAfterSend: true,
			MaxEventsPerMinute:      10,
			DedupWindowSec:          3,
			SingleInstance:          true,
		},
		Tray: TrayConfig{
			Enabled: true,
		},
		Network: NetworkConfig{
			CheckIntervalSec: 5,
			CheckHost:        "api.telegram.org",
			CheckPort:        443,
		},
	}
}

func defaultTriggers() map[string]TriggerConfig {
	return map[string]TriggerConfig{
		"login_success": {
			Enabled:           true,
			CaptureCamera:     false,
			CaptureScreenshot: false,
		},
		"login_failed": {
			Enabled:         true,
			CaptureCamera:   true,
			MaxPhotosPerMin: 3,
		},
		"screen_lock": {
			Enabled: true,
		},
		"screen_unlock": {
			Enabled:       true,
			CaptureCamera: false,
		},
		"system_boot": {
			Enabled: true,
		},
		"system_shutdown": {
			Enabled: true,
		},
		"system_sleep": {
			Enabled: true,
		},
		"system_wake": {
			Enabled: true,
		},
		"usb_inserted": {
			Enabled: true,
		},
		"usb_removed": {
			Enabled: true,
		},
		"network_changed": {
			Enabled: true,
		},
	}
}
