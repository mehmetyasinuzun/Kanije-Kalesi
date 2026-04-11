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

	"github.com/BurntSushi/toml"
)

// Config is the root configuration structure.
// All nested structs use value semantics for thread-safe copying.
type Config struct {
	mu sync.RWMutex `toml:"-" json:"-"`

	Telegram   TelegramConfig            `toml:"telegram"   json:"telegram"`
	Triggers   map[string]TriggerConfig  `toml:"triggers"   json:"triggers"`
	Camera     CameraConfig              `toml:"camera"     json:"camera"`
	Screenshot ScreenshotConfig          `toml:"screenshot" json:"screenshot"`
	Heartbeat  HeartbeatConfig           `toml:"heartbeat"  json:"heartbeat"`
	Storage    StorageConfig             `toml:"storage"    json:"storage"`
	Logging    LoggingConfig             `toml:"logging"    json:"logging"`
	Security   SecurityConfig            `toml:"security"   json:"security"`
	Tray       TrayConfig                `toml:"tray"       json:"tray"`
	Network    NetworkConfig             `toml:"network"    json:"network"`

	// Runtime path — where to save changes
	filePath string `toml:"-" json:"-"`
}

type TelegramConfig struct {
	BotToken       string  `toml:"bot_token"        json:"bot_token"`
	ChatID         int64   `toml:"chat_id"          json:"chat_id"`
	AllowedChatIDs []int64 `toml:"allowed_chat_ids" json:"allowed_chat_ids"`
	SendTimeoutSec int     `toml:"send_timeout_sec" json:"send_timeout_sec"`
	RetryCount     int     `toml:"retry_count"      json:"retry_count"`
	RetryDelaySec  int     `toml:"retry_delay_sec"  json:"retry_delay_sec"`
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

type ScreenshotConfig struct {
	JPEGQuality int    `toml:"jpeg_quality" json:"jpeg_quality"`
	SaveLocal   bool   `toml:"save_local"   json:"save_local"`
	LocalPath   string `toml:"local_path"   json:"local_path"`
}

type HeartbeatConfig struct {
	Enabled       bool `toml:"enabled"        json:"enabled"`
	IntervalHours int  `toml:"interval_hours" json:"interval_hours"`
	IncludeUptime bool `toml:"include_uptime" json:"include_uptime"`
	IncludeDisk   bool `toml:"include_disk"   json:"include_disk"`
}

type StorageConfig struct {
	DBPath              string `toml:"db_path"               json:"db_path"`
	MaxRecentEvents     int    `toml:"max_recent_events"     json:"max_recent_events"`
	EventRetentionDays  int    `toml:"event_retention_days"  json:"event_retention_days"`
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
	DedupWindowSec          int  `toml:"dedup_window_sec"           json:"dedup_window_sec"`
	SingleInstance          bool `toml:"single_instance"            json:"single_instance"`
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
// Missing values are filled with defaults. Environment variables override
// TOML values with KANIJE_ prefix (e.g. KANIJE_TELEGRAM_BOT_TOKEN).
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

	// Ensure Triggers map is initialized
	if cfg.Triggers == nil {
		cfg.Triggers = defaultTriggers()
	}
	for k, v := range defaultTriggers() {
		if _, ok := cfg.Triggers[k]; !ok {
			cfg.Triggers[k] = v
		}
	}

	return cfg, nil
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
// The write is performed to a temp file first, then renamed — no partial writes.
func (c *Config) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.filePath == "" {
		return fmt.Errorf("config dosya yolu belirtilmemiş")
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
	f.Close()

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
	c.Triggers[name] = trig
	c.mu.Unlock()
	return c.Save()
}

// SetField updates a single config field by dot-notation key and string value.
// This is used by the Telegram setup wizard to apply changes from user input.
// Returns a user-friendly error if the key is unknown or the value is invalid.
func (c *Config) SetField(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch strings.ToLower(key) {
	case "telegram.bot_token":
		c.Telegram.BotToken = value
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
	default:
		return fmt.Errorf("bilinmeyen ayar: %q", key)
	}

	return c.save()
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

// save is the internal (unlocked-by-caller) version of Save.
func (c *Config) save() error {
	if c.filePath == "" {
		return nil // In-memory only — silently skip
	}
	tmp := c.filePath + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if err := toml.NewEncoder(f).Encode(c); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, c.filePath)
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

func contains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
