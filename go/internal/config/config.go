// Package config handles loading, validating, and persisting the application
// configuration. The primary format is TOML; sensitive values can be supplied
// via environment variables or configured at runtime through the Telegram bot.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/kanije-kalesi/kanije/internal/secret"
)

// Config is the root configuration structure.
// All nested structs use value semantics for thread-safe copying.
type Config struct {
	mu sync.RWMutex `toml:"-" json:"-"`

	Telegram   TelegramConfig           `toml:"telegram"   json:"telegram"`
	Triggers   map[string]TriggerConfig `toml:"triggers"   json:"triggers"`
	Device     DeviceConfig             `toml:"device"     json:"device"`
	Camera     CameraConfig             `toml:"camera"     json:"camera"`
	Audio      AudioConfig              `toml:"audio"      json:"audio"`
	Screenshot ScreenshotConfig         `toml:"screenshot" json:"screenshot"`
	Motion     MotionConfig             `toml:"motion"     json:"motion"`
	Heartbeat  HeartbeatConfig          `toml:"heartbeat"  json:"heartbeat"`
	Storage    StorageConfig            `toml:"storage"    json:"storage"`
	Logging    LoggingConfig            `toml:"logging"     json:"logging"`
	Security   SecurityConfig           `toml:"security"    json:"security"`
	Protection ProtectionConfig         `toml:"protection"  json:"protection"`
	Tray       TrayConfig               `toml:"tray"        json:"tray"`
	Network    NetworkConfig            `toml:"network"     json:"network"`
	QuietHours QuietHoursConfig         `toml:"quiet_hours" json:"quiet_hours"`
	Update     UpdateConfig             `toml:"update"      json:"update"`
	GeoIP      GeoIPConfig              `toml:"geoip"       json:"geoip"`
	Metrics    MetricsConfig            `toml:"metrics"     json:"metrics"`
	Webhooks   []WebhookConfig          `toml:"webhooks"    json:"webhooks"`

	// Runtime path — where to save changes
	filePath string `toml:"-" json:"-"`
}

type TelegramConfig struct {
	BotToken       string  `toml:"bot_token"        json:"bot_token"`
	ChatID         int64   `toml:"chat_id"          json:"chat_id"`  // owner's private chat/user ID
	GroupID        int64   `toml:"group_id"         json:"group_id"` // optional fleet group; notifications go here when set
	AllowedChatIDs []int64 `toml:"allowed_chat_ids" json:"allowed_chat_ids"`
	SendTimeoutSec int     `toml:"send_timeout_sec" json:"send_timeout_sec"`
	RetryCount     int     `toml:"retry_count"      json:"retry_count"`
	RetryDelaySec  int     `toml:"retry_delay_sec"  json:"retry_delay_sec"`
}

// DeviceConfig identifies this machine within a multi-device fleet.
type DeviceConfig struct {
	Label string `toml:"label" json:"label"` // friendly name; empty = hostname
}

// TriggerConfig controls what happens when a specific event fires.
type TriggerConfig struct {
	Enabled           bool `toml:"enabled"              json:"enabled"`
	CaptureCamera     bool `toml:"capture_camera"       json:"capture_camera"`
	CaptureScreenshot bool `toml:"capture_screenshot"   json:"capture_screenshot"`
	MaxPhotosPerMin   int  `toml:"max_photos_per_minute" json:"max_photos_per_minute"`
}

type CameraConfig struct {
	FFmpegPath   string `toml:"ffmpeg_path"   json:"ffmpeg_path"`
	DeviceIndex  int    `toml:"device_index"  json:"device_index"`
	DeviceName   string `toml:"device_name"   json:"device_name"` // Windows dshow device name
	Width        int    `toml:"width"         json:"width"`
	Height       int    `toml:"height"        json:"height"`
	WarmupFrames int    `toml:"warmup_frames" json:"warmup_frames"`
	JPEGQuality  int    `toml:"jpeg_quality"  json:"jpeg_quality"`
	SaveLocal    bool   `toml:"save_local"    json:"save_local"`
	LocalPath    string `toml:"local_path"    json:"local_path"`
}

// AudioConfig holds microphone (voice recording) settings for /seskayit.
type AudioConfig struct {
	FFmpegPath string `toml:"ffmpeg_path" json:"ffmpeg_path"` // empty = reuse camera's / "ffmpeg"
	DeviceName string `toml:"device_name" json:"device_name"` // Windows dshow mic name; empty = auto-detect
	Bitrate    string `toml:"bitrate"     json:"bitrate"`     // MP3 bitrate, e.g. "96k"
}

type ScreenshotConfig struct {
	JPEGQuality int    `toml:"jpeg_quality" json:"jpeg_quality"`
	SaveLocal   bool   `toml:"save_local"   json:"save_local"`
	LocalPath   string `toml:"local_path"   json:"local_path"`
}

// MotionConfig controls the camera motion detector (/hareket). When enabled, the
// agent grabs a frame every IntervalSec and raises a motion_detected event (with
// the photo) whenever the mean pixel change exceeds Threshold.
type MotionConfig struct {
	Enabled     bool    `toml:"enabled"      json:"enabled"`
	IntervalSec int     `toml:"interval_sec" json:"interval_sec"` // seconds between frames
	Threshold   float64 `toml:"threshold"    json:"threshold"`    // mean luma diff (0-255) to trigger
}

type HeartbeatConfig struct {
	Enabled       bool `toml:"enabled"        json:"enabled"`
	IntervalHours int  `toml:"interval_hours" json:"interval_hours"`
	IncludeUptime bool `toml:"include_uptime" json:"include_uptime"`
	IncludeDisk   bool `toml:"include_disk"   json:"include_disk"`
}

type StorageConfig struct {
	DBPath             string `toml:"db_path"               json:"db_path"`
	MaxRecentEvents    int    `toml:"max_recent_events"     json:"max_recent_events"`
	EventRetentionDays int    `toml:"event_retention_days"  json:"event_retention_days"`
}

type LoggingConfig struct {
	Level         string `toml:"level"          json:"level"`
	File          string `toml:"file"           json:"file"`
	MaxSizeMB     int    `toml:"max_size_mb"    json:"max_size_mb"`
	BackupCount   int    `toml:"backup_count"   json:"backup_count"`
	ConsoleOutput bool   `toml:"console_output" json:"console_output"`
	JSONFormat    bool   `toml:"json_format"    json:"json_format"`
}

type SecurityConfig struct {
	DeleteCapturesAfterSend bool `toml:"delete_captures_after_send" json:"delete_captures_after_send"`
	MaxEventsPerMinute      int  `toml:"max_events_per_minute"      json:"max_events_per_minute"`
	MaxCommandsPerMinute    int  `toml:"max_commands_per_minute"    json:"max_commands_per_minute"`
	DedupWindowSec          int  `toml:"dedup_window_sec"           json:"dedup_window_sec"`
	SingleInstance          bool `toml:"single_instance"            json:"single_instance"`
	// TamperWatch, when true, periodically verifies the agent's own binary,
	// config and DB are present and (Windows) its scheduled task is still enabled,
	// and detects an unclean previous shutdown — alerting the owner on tampering.
	TamperWatch bool `toml:"tamper_watch" json:"tamper_watch"`
	// WipeTargets are extra directories whose CONTENTS /imha securely erases, on
	// top of the user's Music folder and the agent's own data. There is NO factory
	// reset — this keeps the OS intact and lowers the AV detection surface.
	WipeTargets []string `toml:"wipe_targets" json:"wipe_targets"`
	// AuditPaths are folders under Windows SACL file-access auditing (/erisim), so
	// the owner can see who opened/copied files there.
	AuditPaths []string `toml:"audit_paths" json:"audit_paths"`
	// TOTPSecret, when set (base32), requires a valid 2FA code for dangerous
	// commands (/kapat, /yeniden). Empty disables 2FA.
	TOTPSecret string `toml:"totp_secret" json:"-"`
}

// ProtectionConfig is the physical-threat protection policy. Each trigger
// (dead-man switch, USB removal, failed logins) maps to an action chain encoded
// as a string containing any of "lock", "alert", "wipe" (e.g. "lock_alert" =
// lock + alert; "lock_alert_wipe" = all three). Wipe-bearing actions are opt-in.
// Managed via /koruma. See app/protect.go for execution.
type ProtectionConfig struct {
	Enabled bool `toml:"enabled" json:"enabled"` // master switch

	DeadManEnabled bool   `toml:"deadman_enabled" json:"deadman_enabled"`
	DeadManHours   int    `toml:"deadman_hours"   json:"deadman_hours"` // no check-in for this long → fire
	DeadManAction  string `toml:"deadman_action"  json:"deadman_action"`

	USBEnabled bool   `toml:"usb_enabled" json:"usb_enabled"`
	USBDevice  string `toml:"usb_device"  json:"usb_device"` // empty = any USB removal triggers
	USBAction  string `toml:"usb_action"  json:"usb_action"`

	FailedLoginEnabled   bool   `toml:"failed_login_enabled"   json:"failed_login_enabled"`
	FailedLoginThreshold int    `toml:"failed_login_threshold" json:"failed_login_threshold"`
	FailedLoginAction    string `toml:"failed_login_action"    json:"failed_login_action"`

	RAMOnly bool `toml:"ram_only" json:"ram_only"` // keep DB in memory, never persist to disk

	// Lockdown is a persistent "lost mode": while on, the screen is re-locked the
	// moment it's unlocked, so the device is effectively unusable until the owner
	// turns it off from Telegram (/kilit tam kapat).
	Lockdown bool `toml:"lockdown" json:"lockdown"`

	// Canary (honeypot): decoy files under CanaryPath are watched via SACL; on
	// access, CanaryAction fires (catch an intruder reaching for "juicy" files).
	CanaryEnabled bool   `toml:"canary_enabled" json:"canary_enabled"`
	CanaryPath    string `toml:"canary_path"    json:"canary_path"`
	CanaryAction  string `toml:"canary_action"  json:"canary_action"`
}

// QuietHoursConfig suppresses non-critical notifications during a daily window
// (e.g. overnight). Alert/critical events always get through.
type QuietHoursConfig struct {
	Enabled   bool `toml:"enabled"    json:"enabled"`
	StartHour int  `toml:"start_hour" json:"start_hour"` // 0-23, dahil
	EndHour   int  `toml:"end_hour"   json:"end_hour"`   // 0-23, hariç
}

// UpdateConfig controls automatic update checks against GitHub Releases.
type UpdateConfig struct {
	Enabled            bool `toml:"enabled"              json:"enabled"`
	CheckIntervalHours int  `toml:"check_interval_hours" json:"check_interval_hours"`
	// AutoInstall, when true, installs a newer version unattended (Windows only).
	// When false (default) the agent only notifies; the user runs /guncelle.
	AutoInstall bool `toml:"auto_install" json:"auto_install"`
}

// GeoIPConfig controls IP geolocation enrichment of login events.
type GeoIPConfig struct {
	// Enabled, when true, looks up the source IP of login events and adds the
	// country/city/ISP to the notification. Uses the free ipwho.is HTTPS API.
	Enabled bool `toml:"enabled" json:"enabled"`
}

// MetricsConfig controls the optional Prometheus /metrics HTTP endpoint.
type MetricsConfig struct {
	Enabled bool   `toml:"enabled" json:"enabled"`
	Addr    string `toml:"addr"    json:"addr"` // listen address, e.g. "127.0.0.1:9099"
}

// WebhookConfig is one outgoing webhook destination (Discord or generic JSON).
type WebhookConfig struct {
	Name   string `toml:"name"   json:"name"`
	URL    string `toml:"url"    json:"url"`
	Format string `toml:"format" json:"format"` // "discord" | "json"
}

type TrayConfig struct {
	Enabled bool `toml:"enabled" json:"enabled"`
}

type NetworkConfig struct {
	CheckIntervalSec int    `toml:"check_interval_sec" json:"check_interval_sec"`
	CheckHost        string `toml:"check_host"         json:"check_host"`
	CheckPort        int    `toml:"check_port"         json:"check_port"`
}

// ----- Loading -----

// Load reads configuration from the given TOML file path.
// Missing values are filled with defaults. Environment variables (KANIJE_BOT_TOKEN,
// KANIJE_CHAT_ID, KANIJE_LOG_LEVEL, KANIJE_DB_PATH) override TOML values. All
// values are clamped to safe ranges so a hand-edited config can never crash the app.
func Load(path string) (*Config, error) {
	cfg := Defaults()
	cfg.filePath = path

	if _, err := os.Stat(path); err == nil {
		if _, err := toml.DecodeFile(path, cfg); err != nil {
			return nil, fmt.Errorf("config okuma hatası (%s): %w", path, err)
		}
	}
	// Always try to fill from environment
	applyEnvOverrides(cfg)

	// Ensure Triggers map is initialized and back-fill any newly added triggers
	// (e.g. network_up/network_down) into older config files.
	if cfg.Triggers == nil {
		cfg.Triggers = defaultTriggers()
	}
	for k, v := range defaultTriggers() {
		if _, ok := cfg.Triggers[k]; !ok {
			cfg.Triggers[k] = v
		}
	}

	cfg.validate()
	return cfg, nil
}

// validate clamps configuration values to safe operating ranges. It runs once
// during Load (single-goroutine) so a partial or hand-edited config — a zero
// rate limit, a negative interval, an out-of-range JPEG quality — can never
// panic or misbehave at runtime.
func (c *Config) validate() {
	// clampMin returns def when v is below min, otherwise v.
	clampMin := func(v, min, def int) int {
		if v < min {
			return def
		}
		return v
	}
	clampQuality := func(q, def int) int {
		if q < 1 || q > 100 {
			return def
		}
		return q
	}

	c.Telegram.SendTimeoutSec = clampMin(c.Telegram.SendTimeoutSec, 1, 15)
	if c.Telegram.RetryCount < 0 {
		c.Telegram.RetryCount = 0
	}
	if c.Telegram.RetryDelaySec < 0 {
		c.Telegram.RetryDelaySec = 0
	}

	if c.Camera.FFmpegPath == "" {
		c.Camera.FFmpegPath = "ffmpeg"
	}
	c.Camera.Width = clampMin(c.Camera.Width, 1, 640)
	c.Camera.Height = clampMin(c.Camera.Height, 1, 480)
	if c.Camera.WarmupFrames < 0 {
		c.Camera.WarmupFrames = 0
	}
	c.Camera.JPEGQuality = clampQuality(c.Camera.JPEGQuality, 85)

	// Audio: fall back to the camera's ffmpeg binary, then PATH.
	if c.Audio.FFmpegPath == "" {
		c.Audio.FFmpegPath = c.Camera.FFmpegPath
	}
	if c.Audio.Bitrate == "" {
		c.Audio.Bitrate = "96k"
	}

	c.Screenshot.JPEGQuality = clampQuality(c.Screenshot.JPEGQuality, 75)

	c.Motion.IntervalSec = clampMin(c.Motion.IntervalSec, 1, 3)
	if c.Motion.Threshold <= 0 {
		c.Motion.Threshold = 12
	}

	c.Protection.DeadManHours = clampMin(c.Protection.DeadManHours, 1, 72)
	c.Protection.FailedLoginThreshold = clampMin(c.Protection.FailedLoginThreshold, 1, 5)
	c.Protection.DeadManAction = validProtectionAction(c.Protection.DeadManAction)
	c.Protection.USBAction = validProtectionAction(c.Protection.USBAction)
	c.Protection.FailedLoginAction = validProtectionAction(c.Protection.FailedLoginAction)
	c.Protection.CanaryAction = validProtectionAction(c.Protection.CanaryAction)

	c.Heartbeat.IntervalHours = clampMin(c.Heartbeat.IntervalHours, 1, 6)

	if c.Storage.DBPath == "" {
		c.Storage.DBPath = "./kanije.db"
	}
	c.Storage.MaxRecentEvents = clampMin(c.Storage.MaxRecentEvents, 1, 10)
	if c.Storage.MaxRecentEvents > 50 {
		c.Storage.MaxRecentEvents = 50 // keep the /olaylar message within Telegram limits
	}
	c.Storage.EventRetentionDays = clampMin(c.Storage.EventRetentionDays, 1, 30)

	c.Security.MaxEventsPerMinute = clampMin(c.Security.MaxEventsPerMinute, 1, 10)
	c.Security.MaxCommandsPerMinute = clampMin(c.Security.MaxCommandsPerMinute, 1, 20)
	if c.Security.DedupWindowSec < 0 {
		c.Security.DedupWindowSec = 0 // 0 disables dedup
	}

	if c.QuietHours.StartHour < 0 || c.QuietHours.StartHour > 23 {
		c.QuietHours.StartHour = 23
	}
	if c.QuietHours.EndHour < 0 || c.QuietHours.EndHour > 23 {
		c.QuietHours.EndHour = 7
	}

	c.Network.CheckIntervalSec = clampMin(c.Network.CheckIntervalSec, 1, 5)
	if c.Network.CheckHost == "" {
		c.Network.CheckHost = "api.telegram.org"
	}
	if c.Network.CheckPort <= 0 || c.Network.CheckPort > 65535 {
		c.Network.CheckPort = 443
	}

	c.Logging.Level = strings.ToLower(c.Logging.Level)
	if !contains([]string{"debug", "info", "warn", "error"}, c.Logging.Level) {
		c.Logging.Level = "info"
	}

	c.Update.CheckIntervalHours = clampMin(c.Update.CheckIntervalHours, 1, 24)
}

// applyEnvOverrides reads KANIJE_* environment variables and overrides config.
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("KANIJE_BOT_TOKEN"); v != "" {
		cfg.Telegram.BotToken = v
	}
	if v := os.Getenv("KANIJE_CHAT_ID"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			cfg.Telegram.ChatID = id
		}
	}
	if v := os.Getenv("KANIJE_LOG_LEVEL"); v != "" {
		cfg.Logging.Level = v
	}
	if v := os.Getenv("KANIJE_DB_PATH"); v != "" {
		cfg.Storage.DBPath = v
	}
}

// ----- Persistence -----

// Save atomically writes the current configuration back to its source file.
func (c *Config) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.filePath == "" {
		return fmt.Errorf("config dosya yolu belirtilmemiş")
	}
	return c.persist()
}

// persist atomically writes the config to its file (temp file + rename, so no
// partial writes are ever observed). The caller MUST hold at least the read
// lock. With no file path configured it is a no-op (env-only configuration).
func (c *Config) persist() error {
	if c.filePath == "" {
		return nil
	}

	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("dizin oluşturma hatası: %w", err)
	}

	tmp := c.filePath + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("geçici dosya oluşturma hatası: %w", err)
	}

	if err := toml.NewEncoder(f).Encode(c); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("config yazma hatası: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("config kapatma hatası: %w", err)
	}

	if err := os.Rename(tmp, c.filePath); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("config kaydetme hatası: %w", err)
	}
	return nil
}

// ----- Runtime mutation (for Telegram config wizard) -----

// SetTrigger enables or updates a trigger at runtime and persists the change.
func (c *Config) SetTrigger(name string, trig TriggerConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Triggers[name] = trig
	return c.persist()
}

// SetField updates a single config field by dot-notation key and string value.
// This is used by the Telegram setup wizard to apply changes from user input.
// Returns a user-friendly error if the key is unknown or the value is invalid.
func (c *Config) SetField(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch strings.ToLower(key) {
	case "telegram.bot_token":
		c.Telegram.BotToken = secret.Protect(value)
	case "telegram.chat_id":
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("geçersiz chat_id: sayı olmalı")
		}
		c.Telegram.ChatID = id
	case "heartbeat.interval_hours":
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return fmt.Errorf("geçersiz değer: 1 veya daha büyük bir sayı girin")
		}
		c.Heartbeat.IntervalHours = n
	case "heartbeat.enabled":
		c.Heartbeat.Enabled = parseBool(value)
	case "camera.device_index":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("geçersiz kamera indeksi")
		}
		c.Camera.DeviceIndex = n
	case "camera.device_name":
		c.Camera.DeviceName = value
	case "audio.device_name":
		c.Audio.DeviceName = value
	case "device.label":
		c.Device.Label = strings.TrimSpace(value)
	case "telegram.group_id":
		n, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
		if err != nil {
			return fmt.Errorf("geçersiz grup ID (sayısal olmalı, örn: -1001234567890)")
		}
		c.Telegram.GroupID = n
	case "logging.level":
		value = strings.ToLower(value)
		if !contains([]string{"debug", "info", "warn", "error"}, value) {
			return fmt.Errorf("geçersiz log seviyesi: debug/info/warn/error olmalı")
		}
		c.Logging.Level = value
	case "security.max_events_per_minute":
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return fmt.Errorf("geçersiz değer: 1+ sayı girin")
		}
		c.Security.MaxEventsPerMinute = n
	case "security.delete_captures_after_send":
		c.Security.DeleteCapturesAfterSend = parseBool(value)
	case "camera.save_local":
		c.Camera.SaveLocal = parseBool(value)
	default:
		return fmt.Errorf("bilinmeyen ayar: %q", key)
	}

	return c.persist()
}

// GetBool returns the current value of a boolean setting by dot-notation key.
// Used by the Telegram wizard to read-then-flip toggle switches.
func (c *Config) GetBool(key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	switch strings.ToLower(key) {
	case "heartbeat.enabled":
		return c.Heartbeat.Enabled, nil
	case "security.delete_captures_after_send":
		return c.Security.DeleteCapturesAfterSend, nil
	case "camera.save_local":
		return c.Camera.SaveLocal, nil
	default:
		return false, fmt.Errorf("bilinmeyen bool ayar: %q", key)
	}
}

// GetSafeJSON returns the config as JSON with the bot token masked.
// Used by the /ayarlar Telegram command.
func (c *Config) GetSafeJSON() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	type safe struct {
		Telegram struct {
			BotToken string `json:"bot_token"`
			ChatID   int64  `json:"chat_id"`
		} `json:"telegram"`
		Heartbeat HeartbeatConfig `json:"heartbeat"`
		Camera    CameraConfig    `json:"camera"`
		Audio     AudioConfig     `json:"audio"`
		Security  SecurityConfig  `json:"security"`
		Logging   LoggingConfig   `json:"logging"`
	}

	s := safe{}
	if len(c.Telegram.BotToken) > 8 {
		s.Telegram.BotToken = c.Telegram.BotToken[:4] + "****" + c.Telegram.BotToken[len(c.Telegram.BotToken)-4:]
	} else {
		s.Telegram.BotToken = "****"
	}
	s.Telegram.ChatID = c.Telegram.ChatID
	s.Heartbeat = c.Heartbeat
	s.Camera = c.Camera
	s.Audio = c.Audio
	s.Security = c.Security
	s.Logging = c.Logging

	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}

// GetTrigger returns a copy of the named trigger config.
func (c *Config) GetTrigger(name string) (TriggerConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	t, ok := c.Triggers[name]
	return t, ok
}

// ----- Thread-safe accessors -----
//
// These getters take the read lock so callers can safely read individual
// settings while the Telegram setup wizard mutates the config concurrently.
// Always call them at the point of use — never cache the result, or runtime
// changes made via /kurulum will not take effect.

// ChatID returns the owner's private chat/user ID (the RBAC root).
func (c *Config) ChatID() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Telegram.ChatID
}

// GroupID returns the optional fleet group chat ID (0 if unset).
func (c *Config) GroupID() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Telegram.GroupID
}

// NotifyChatID returns where event notifications should be sent: the fleet
// group when configured, otherwise the owner's chat.
func (c *Config) NotifyChatID() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.Telegram.GroupID != 0 {
		return c.Telegram.GroupID
	}
	return c.Telegram.ChatID
}

// DeviceLabel returns the configured friendly device label, or "" if unset
// (callers fall back to the hostname).
func (c *Config) DeviceLabel() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Device.Label
}

// BotToken returns the decrypted Telegram bot token. The stored value may be
// DPAPI-encrypted (Windows); decryption happens here at point of use.
func (c *Config) BotToken() string {
	c.mu.RLock()
	stored := c.Telegram.BotToken
	c.mu.RUnlock()

	tok, err := secret.Unprotect(stored)
	if err != nil {
		return ""
	}
	return tok
}

// AllowedChatIDs returns a copy of the additional authorized chat IDs.
func (c *Config) AllowedChatIDs() []int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.Telegram.AllowedChatIDs) == 0 {
		return nil
	}
	out := make([]int64, len(c.Telegram.AllowedChatIDs))
	copy(out, c.Telegram.AllowedChatIDs)
	return out
}

// HeartbeatEnabled reports whether periodic heartbeat messages are enabled.
func (c *Config) HeartbeatEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Heartbeat.Enabled
}

// HeartbeatInterval returns the heartbeat interval as a duration.
func (c *Config) HeartbeatInterval() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Duration(c.Heartbeat.IntervalHours) * time.Hour
}

// MaxEventsPerMinute returns the per-type event rate limit.
func (c *Config) MaxEventsPerMinute() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Security.MaxEventsPerMinute
}

// DedupWindow returns the deduplication window as a duration.
func (c *Config) DedupWindow() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Duration(c.Security.DedupWindowSec) * time.Second
}

// MaxRecentEvents returns how many recent events the /olaylar command shows.
func (c *Config) MaxRecentEvents() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Storage.MaxRecentEvents
}

// MaxCommandsPerMinute returns the per-chat bot-command rate limit.
func (c *Config) MaxCommandsPerMinute() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Security.MaxCommandsPerMinute
}

// TOTPSecret returns the configured 2FA secret (base32), or "" if disabled.
func (c *Config) TOTPSecret() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Security.TOTPSecret
}

// TamperWatchEnabled reports whether the anti-tamper watchdog is active.
func (c *Config) TamperWatchEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Security.TamperWatch
}

// WipeTargets returns a copy of the extra directories /imha securely erases.
func (c *Config) WipeTargets() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.Security.WipeTargets) == 0 {
		return nil
	}
	out := make([]string, len(c.Security.WipeTargets))
	copy(out, c.Security.WipeTargets)
	return out
}

// AddWipeTarget appends a directory to the /imha wipe list (de-duplicated) and
// persists the change.
func (c *Config) AddWipeTarget(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range c.Security.WipeTargets {
		if strings.EqualFold(p, path) {
			return nil // already present
		}
	}
	c.Security.WipeTargets = append(c.Security.WipeTargets, path)
	return c.persist()
}

// RemoveWipeTarget drops a directory from the /imha wipe list (case-insensitive)
// and persists. Returns true if one was removed.
func (c *Config) RemoveWipeTarget(path string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, p := range c.Security.WipeTargets {
		if strings.EqualFold(p, path) {
			c.Security.WipeTargets = append(c.Security.WipeTargets[:i], c.Security.WipeTargets[i+1:]...)
			return true, c.persist()
		}
	}
	return false, nil
}

// AuditPaths returns a copy of the folders under file-access auditing (/erisim).
func (c *Config) AuditPaths() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.Security.AuditPaths) == 0 {
		return nil
	}
	out := make([]string, len(c.Security.AuditPaths))
	copy(out, c.Security.AuditPaths)
	return out
}

// AddAuditPath records a folder as being audited (de-duplicated) and persists.
func (c *Config) AddAuditPath(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, p := range c.Security.AuditPaths {
		if strings.EqualFold(p, path) {
			return nil
		}
	}
	c.Security.AuditPaths = append(c.Security.AuditPaths, path)
	return c.persist()
}

// RemoveAuditPath drops a folder from the audited list and persists. Returns true
// if one was removed.
func (c *Config) RemoveAuditPath(path string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, p := range c.Security.AuditPaths {
		if strings.EqualFold(p, path) {
			c.Security.AuditPaths = append(c.Security.AuditPaths[:i], c.Security.AuditPaths[i+1:]...)
			return true, c.persist()
		}
	}
	return false, nil
}

// ProtectionPolicy returns a copy of the physical-threat protection policy.
func (c *Config) ProtectionPolicy() ProtectionConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Protection
}

// SetProtection replaces the protection policy and persists it.
func (c *Config) SetProtection(p ProtectionConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Protection = p
	return c.persist()
}

// RAMOnly reports whether the DB should be kept in memory only (never persisted).
func (c *Config) RAMOnly() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Protection.RAMOnly
}

// LockdownActive reports whether persistent lock mode is on.
func (c *Config) LockdownActive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Protection.Lockdown
}

// SetLockdown toggles persistent lock mode and persists it. Called by /kilit tam
// and by the "lockdown" protection action.
func (c *Config) SetLockdown(on bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Protection.Lockdown = on
	return c.persist()
}

// DeleteCapturesAfterSend reports whether locally-saved captures are removed
// right after a successful send (anti-forensics; only meaningful with SaveLocal).
func (c *Config) DeleteCapturesAfterSend() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Security.DeleteCapturesAfterSend
}

// CameraSaveLocal reports whether camera photos are also written to disk, and to
// which directory.
func (c *Config) CameraSaveLocal() (bool, string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Camera.SaveLocal, c.Camera.LocalPath
}

// ScreenshotSaveLocal reports whether screenshots are also written to disk, and
// to which directory.
func (c *Config) ScreenshotSaveLocal() (bool, string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Screenshot.SaveLocal, c.Screenshot.LocalPath
}

// MotionEnabled reports whether the camera motion detector is running.
func (c *Config) MotionEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Motion.Enabled
}

// MotionSettings returns the frame interval and trigger threshold (with safe
// fallbacks), for the motion-detection loop.
func (c *Config) MotionSettings() (time.Duration, float64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	iv := c.Motion.IntervalSec
	if iv < 1 {
		iv = 3
	}
	th := c.Motion.Threshold
	if th <= 0 {
		th = 12
	}
	return time.Duration(iv) * time.Second, th
}

// SetMotionEnabled turns the motion detector on/off and persists the change.
func (c *Config) SetMotionEnabled(on bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Motion.Enabled = on
	return c.persist()
}

// SetMotionThreshold updates the trigger threshold and persists the change.
func (c *Config) SetMotionThreshold(th float64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Motion.Threshold = th
	return c.persist()
}

// TOTPEnabled reports whether dangerous commands require a 2FA code.
func (c *Config) TOTPEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Security.TOTPSecret != ""
}

// UpdateEnabled reports whether periodic update checks are enabled.
func (c *Config) UpdateEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Update.Enabled
}

// UpdateInterval returns how often to check for updates.
func (c *Config) UpdateInterval() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Duration(c.Update.CheckIntervalHours) * time.Hour
}

// UpdateAutoInstall reports whether newer versions install unattended.
func (c *Config) UpdateAutoInstall() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Update.AutoInstall
}

// GeoIPEnabled reports whether login events are enriched with IP geolocation.
func (c *Config) GeoIPEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.GeoIP.Enabled
}

// MetricsEnabled reports whether the Prometheus endpoint is enabled.
func (c *Config) MetricsEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Metrics.Enabled
}

// MetricsAddr returns the metrics listen address.
func (c *Config) MetricsAddr() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.Metrics.Addr == "" {
		return "127.0.0.1:9099"
	}
	return c.Metrics.Addr
}

// WebhookTargets returns a copy of the configured outgoing webhook destinations.
func (c *Config) WebhookTargets() []WebhookConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.Webhooks) == 0 {
		return nil
	}
	out := make([]WebhookConfig, len(c.Webhooks))
	copy(out, c.Webhooks)
	return out
}

// InQuietHours reports whether the given time falls inside the configured quiet
// window. Handles overnight wraparound (e.g. 23:00–07:00). During quiet hours
// only alert/critical events notify; everything else is stored silently.
func (c *Config) InQuietHours(now time.Time) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.QuietHours.Enabled {
		return false
	}
	s, e, h := c.QuietHours.StartHour, c.QuietHours.EndHour, now.Hour()
	if s == e {
		return false
	}
	if s < e {
		return h >= s && h < e // aynı gün içinde
	}
	return h >= s || h < e // gece yarısını aşan pencere
}

// FilePath returns the config file path in use.
func (c *Config) FilePath() string { return c.filePath }

// SetFilePath sets the config file path (used by cmdSetup before first Save).
func (c *Config) SetFilePath(path string) { c.filePath = path }

// IsConfigured returns true if the minimum required settings are present.
func (c *Config) IsConfigured() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Telegram.BotToken != "" && c.Telegram.ChatID != 0
}

// helpers

func parseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "evet" || s == "açık" || s == "on"
}

// validProtectionAction returns a recognized protection action chain, defaulting
// to the safe "lock_alert" (no data loss) for unknown/empty values.
func validProtectionAction(a string) string {
	switch a {
	case "lock", "alert", "lock_alert", "lock_wipe", "alert_wipe", "lock_alert_wipe", "wipe",
		"lockdown", "lockdown_alert", "alert_lockdown", "lock_alert_lockdown", "lockdown_alert_wipe":
		return a
	default:
		return "lock_alert"
	}
}

func contains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
