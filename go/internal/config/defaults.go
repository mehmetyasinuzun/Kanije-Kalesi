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
		Motion: MotionConfig{
			Enabled:     false, // opt-in via /hareket ac
			IntervalSec: 3,
			Threshold:   12,
		},
		Heartbeat: HeartbeatConfig{
			Enabled:       true,
			IntervalHours: 6,
			IncludeUptime: true,
			IncludeDisk:   true,
		},
		Storage: StorageConfig{
			DBPath:             "./kanije.db",
			MaxRecentEvents:    10,
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
			MaxCommandsPerMinute:    20,
			DedupWindowSec:          3,
			SingleInstance:          true,
			TamperWatch:             true,
		},
		Protection: ProtectionConfig{
			Enabled:              false, // opt-in via /koruma
			DeadManEnabled:       false,
			DeadManHours:         72, // 3 gün
			DeadManAction:        "lock_alert",
			USBEnabled:           false,
			USBAction:            "lock_alert",
			FailedLoginEnabled:   false,
			FailedLoginThreshold: 5,
			FailedLoginAction:    "lock_alert",
			RAMOnly:              false,
		},
		QuietHours: QuietHoursConfig{
			Enabled:   false,
			StartHour: 23,
			EndHour:   7,
		},
		Update: UpdateConfig{
			Enabled:            true,
			CheckIntervalHours: 24,
			AutoInstall:        false, // güvenli varsayılan: sadece bildir, /guncelle ile kur
		},
		GeoIP: GeoIPConfig{
			Enabled: true,
		},
		Metrics: MetricsConfig{
			Enabled: false, // opt-in: açıldığında bir port dinler
			Addr:    "127.0.0.1:9099",
		},
		Tray: TrayConfig{
			// Stealth by default: no system-tray icon. Combined with the
			// -H=windowsgui build (no console window), the agent runs fully
			// hidden — an intruder sees nothing on the desktop.
			Enabled: false,
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
		"network_up": {
			Enabled: true,
		},
		"network_down": {
			Enabled: true,
		},
		"network_changed": {
			Enabled: true,
		},
		"tamper_alert": {
			Enabled:       true,
			CaptureCamera: true, // catch whoever is tampering
		},
		"panic_triggered": {
			Enabled: true,
		},
		"motion_detected": {
			Enabled: true,
		},
		"protection_triggered": {
			Enabled:       true,
			CaptureCamera: true, // capture whoever triggered the protection
		},
	}
}
