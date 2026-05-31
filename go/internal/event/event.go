// Package event defines the core event types and structures used throughout
// the application. All security events flow through this package.
package event

import (
	"fmt"
	"time"
)

// Type identifies what kind of security event occurred.
type Type string

const (
	// Authentication
	TypeLoginSuccess Type = "login_success"
	TypeLoginFailed  Type = "login_failed"
	TypeLogoff       Type = "logoff"

	// Screen
	TypeScreenLock   Type = "screen_lock"
	TypeScreenUnlock Type = "screen_unlock"

	// System
	TypeSystemBoot     Type = "system_boot"
	TypeSystemShutdown Type = "system_shutdown"
	TypeSystemSleep    Type = "system_sleep"
	TypeSystemWake     Type = "system_wake"

	// Hardware
	TypeUSBInserted Type = "usb_inserted"
	TypeUSBRemoved  Type = "usb_removed"

	// Any device interface (mouse/keyboard/phone/audio/HID/...), not just storage.
	TypeDeviceConnected    Type = "device_connected"
	TypeDeviceDisconnected Type = "device_disconnected"

	// Network
	TypeNetworkUp      Type = "network_up"
	TypeNetworkDown    Type = "network_down"
	TypeNetworkChanged Type = "network_changed"

	// Security / anti-theft
	TypeTamperAlert     Type = "tamper_alert"         // exe/config/task tampered with or unexpected kill
	TypePanic           Type = "panic_triggered"      // owner triggered the panic sweep
	TypeMotionDetected  Type = "motion_detected"      // camera motion detector fired
	TypeProtectionFired Type = "protection_triggered" // a protection policy trigger fired (deadman/usb/login)

	// Internal
	TypeHeartbeat Type = "heartbeat"
	TypeError     Type = "internal_error"
)

// Severity represents the urgency level.
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityAlert
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "bilgi"
	case SeverityWarning:
		return "uyarı"
	case SeverityAlert:
		return "alarm"
	case SeverityCritical:
		return "kritik"
	default:
		return "bilinmiyor"
	}
}

func (s Severity) Emoji() string {
	switch s {
	case SeverityInfo:
		return "ℹ️"
	case SeverityWarning:
		return "⚠️"
	case SeverityAlert:
		return "🚨"
	case SeverityCritical:
		return "🔴"
	default:
		return "❓"
	}
}

// LogonType maps Windows logon type IDs to human-readable names.
type LogonType int

const (
	LogonInteractive       LogonType = 2
	LogonNetwork           LogonType = 3
	LogonBatch             LogonType = 4
	LogonService           LogonType = 5
	LogonUnlock            LogonType = 7
	LogonNetworkCleartext  LogonType = 8
	LogonNewCredentials    LogonType = 9
	LogonRemoteInteractive LogonType = 10
	LogonCachedInteractive LogonType = 11
)

func (l LogonType) String() string {
	switch l {
	case LogonInteractive:
		return "Etkileşimli"
	case LogonNetwork:
		return "Ağ"
	case LogonBatch:
		return "Batch"
	case LogonService:
		return "Servis"
	case LogonUnlock:
		return "Kilit Açma"
	case LogonNetworkCleartext:
		return "Ağ (Açık Metin)"
	case LogonNewCredentials:
		return "Yeni Kimlik"
	case LogonRemoteInteractive:
		return "Uzak Masaüstü"
	case LogonCachedInteractive:
		return "Önbellekli Etkileşimli"
	default:
		return fmt.Sprintf("Tip-%d", int(l))
	}
}

// AttachmentType identifies the kind of media attachment.
type AttachmentType string

const (
	AttachmentPhoto      AttachmentType = "photo"
	AttachmentScreenshot AttachmentType = "screenshot"
)

// Attachment represents a media file to be sent alongside an event notification.
type Attachment struct {
	Type    AttachmentType
	Data    []byte // Raw image bytes — no temp file paths
	Caption string
}

// Event is the central data structure representing a single security event.
// All fields use UTF-8 strings. Zero values are safe (omitted in output).
type Event struct {
	// Identity
	ID     int64  `json:"id"`
	Source string `json:"source"` // Listener name that produced this

	// Classification
	Type     Type     `json:"type"`
	Severity Severity `json:"severity"`

	// Timing
	Timestamp time.Time `json:"timestamp"`

	// System context
	Hostname string `json:"hostname,omitempty"`
	Username string `json:"username,omitempty"`
	OS       string `json:"os,omitempty"` // friendly OS name, e.g. "Windows 11"

	// Network context
	SourceIP    string `json:"source_ip,omitempty"`
	NetworkSSID string `json:"network_ssid,omitempty"`
	NetworkType string `json:"network_type,omitempty"` // "WiFi" | "Ethernet"
	LocalIP     string `json:"local_ip,omitempty"`
	PublicIP    string `json:"public_ip,omitempty"` // external/public IP (network events)

	// Authentication context
	LogonType LogonType `json:"logon_type,omitempty"`
	Domain    string    `json:"domain,omitempty"`

	// USB context
	DeviceName  string `json:"device_name,omitempty"`
	DeviceLabel string `json:"device_label,omitempty"`
	DeviceSize  int64  `json:"device_size,omitempty"` // total bytes
	DeviceFree  int64  `json:"device_free,omitempty"` // free bytes (fullness analysis)
	DeviceFS    string `json:"device_fs,omitempty"`   // NTFS, FAT32, exFAT…
	DevicePath  string `json:"device_path,omitempty"` // D:\

	// Power context
	WakeType string `json:"wake_type,omitempty"` // "manuel" | "otomatik"

	// Generic extra fields (for extensibility)
	Extra map[string]string `json:"extra,omitempty"`

	// Media attachments (captured synchronously before queuing)
	Attachments []Attachment `json:"-"`
}

// DefaultSeverity returns the expected severity for a given event type.
func DefaultSeverity(t Type) Severity {
	switch t {
	case TypeLoginFailed:
		return SeverityAlert
	case TypeUSBInserted, TypeDeviceConnected:
		return SeverityWarning
	case TypeSystemBoot, TypeSystemWake, TypeScreenUnlock:
		return SeverityInfo
	case TypeSystemShutdown, TypeSystemSleep:
		return SeverityWarning
	case TypeNetworkDown:
		return SeverityWarning
	case TypeTamperAlert, TypePanic, TypeProtectionFired:
		return SeverityCritical
	case TypeMotionDetected:
		return SeverityAlert
	case TypeError:
		return SeverityAlert
	default:
		return SeverityInfo
	}
}

// Label returns the human-readable Turkish name of an event type.
func (t Type) Label() string {
	switch t {
	case TypeLoginSuccess:
		return "Başarılı Giriş"
	case TypeLoginFailed:
		return "Başarısız Giriş Denemesi"
	case TypeLogoff:
		return "Oturum Kapatıldı"
	case TypeScreenLock:
		return "Ekran Kilitlendi"
	case TypeScreenUnlock:
		return "Ekran Kilidi Açıldı"
	case TypeSystemBoot:
		return "Sistem Başlatıldı"
	case TypeSystemShutdown:
		return "Sistem Kapatılıyor"
	case TypeSystemSleep:
		return "Uyku Moduna Girdi"
	case TypeSystemWake:
		return "Uykudan Uyandı"
	case TypeUSBInserted:
		return "USB Cihazı Takıldı"
	case TypeUSBRemoved:
		return "USB Cihazı Çıkarıldı"
	case TypeDeviceConnected:
		return "Aygıt Bağlandı"
	case TypeDeviceDisconnected:
		return "Aygıt Çıkarıldı"
	case TypeNetworkUp:
		return "İnternet Bağlantısı Kuruldu"
	case TypeNetworkDown:
		return "İnternet Bağlantısı Kesildi"
	case TypeNetworkChanged:
		return "Ağ Değişti"
	case TypeTamperAlert:
		return "Kurcalama Alarmı"
	case TypePanic:
		return "Panik Modu"
	case TypeMotionDetected:
		return "Hareket Algılandı"
	case TypeProtectionFired:
		return "Koruma Tetiklendi"
	case TypeHeartbeat:
		return "Sistem Nabzı"
	default:
		return string(t)
	}
}

// Emoji returns an emoji representing the event type.
func (t Type) Emoji() string {
	switch t {
	case TypeLoginSuccess:
		return "✅"
	case TypeLoginFailed:
		return "🚨"
	case TypeLogoff:
		return "👋"
	case TypeScreenLock:
		return "🔒"
	case TypeScreenUnlock:
		return "🔓"
	case TypeSystemBoot:
		return "🖥️"
	case TypeSystemShutdown:
		return "🔴"
	case TypeSystemSleep:
		return "😴"
	case TypeSystemWake:
		return "☀️"
	case TypeUSBInserted:
		return "🔌"
	case TypeUSBRemoved:
		return "⏏️"
	case TypeDeviceConnected:
		return "🧩"
	case TypeDeviceDisconnected:
		return "🧩"
	case TypeNetworkUp:
		return "🌐"
	case TypeNetworkDown:
		return "📡"
	case TypeNetworkChanged:
		return "🔄"
	case TypeTamperAlert:
		return "🛑"
	case TypePanic:
		return "🆘"
	case TypeMotionDetected:
		return "🎥"
	case TypeProtectionFired:
		return "🛡️"
	case TypeHeartbeat:
		return "💓"
	default:
		return "📌"
	}
}

// New creates a new Event with sensible defaults.
func New(t Type, source string) Event {
	return Event{
		Type:      t,
		Severity:  DefaultSeverity(t),
		Source:    source,
		Timestamp: time.Now(),
	}
}

func (e *Event) String() string {
	return fmt.Sprintf("[%s] %s @ %s (kaynak=%s, kullanıcı=%s)",
		e.Severity, e.Type, e.Timestamp.Format("2006-01-02 15:04:05"), e.Source, e.Username)
}
