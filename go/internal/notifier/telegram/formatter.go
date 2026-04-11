package telegram

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"time"

	"github.com/kanije-kalesi/kanije/internal/event"
)

// version is embedded at build time via ldflags.
var version = "dev"

// FormatEvent converts a security event into a Telegram HTML message.
// All string values go through SafeText() before interpolation — no XSS.
func FormatEvent(ev event.Event) string {
	var b strings.Builder
	b.Grow(512)

	// Header: emoji + event type label
	b.WriteString(eventEmoji(ev.Type))
	b.WriteString(" <b>")
	b.WriteString(eventLabel(ev.Type))
	b.WriteString("</b>\n")

	// Timestamp
	b.WriteString("🕐 <i>")
	b.WriteString(ev.Timestamp.Format("02.01.2006 15:04:05"))
	b.WriteString("</i>\n")

	// System context
	if ev.Hostname != "" {
		b.WriteString("💻 <code>")
		b.WriteString(SafeText(ev.Hostname))
		b.WriteString("</code>\n")
	}
	if ev.Username != "" {
		u := SafeText(ev.Username)
		if ev.Domain != "" {
			u = SafeText(ev.Domain) + `\` + u
		}
		b.WriteString("👤 <code>")
		b.WriteString(u)
		b.WriteString("</code>\n")
	}

	// Event-specific fields
	switch ev.Type {
	case event.TypeLoginSuccess, event.TypeLoginFailed:
		if ev.SourceIP != "" {
			b.WriteString("🌐 IP: <code>")
			b.WriteString(SafeText(ev.SourceIP))
			b.WriteString("</code>\n")
		}
		if ev.LogonType != 0 {
			b.WriteString("🔑 Giriş tipi: ")
			b.WriteString(SafeText(ev.LogonType.String()))
			b.WriteString("\n")
		}

	case event.TypeUSBInserted:
		if ev.DeviceName != "" || ev.DeviceLabel != "" {
			label := ev.DeviceLabel
			if label == "" {
				label = ev.DeviceName
			}
			b.WriteString("💾 Cihaz: <b>")
			b.WriteString(SafeText(label))
			b.WriteString("</b>\n")
		}
		if ev.DevicePath != "" {
			b.WriteString("📂 Yol: <code>")
			b.WriteString(SafeText(ev.DevicePath))
			b.WriteString("</code>\n")
		}
		if ev.DeviceFS != "" {
			b.WriteString("🗂️ Dosya sistemi: ")
			b.WriteString(SafeText(ev.DeviceFS))
			b.WriteString("\n")
		}
		if ev.DeviceSize > 0 {
			b.WriteString("📦 Boyut: ")
			b.WriteString(formatBytes(ev.DeviceSize))
			b.WriteString("\n")
		}

	case event.TypeUSBRemoved:
		if ev.DeviceName != "" {
			b.WriteString("💾 Cihaz: <b>")
			b.WriteString(SafeText(ev.DeviceName))
			b.WriteString("</b>\n")
		}
		if ev.DevicePath != "" {
			b.WriteString("📂 Yol: <code>")
			b.WriteString(SafeText(ev.DevicePath))
			b.WriteString("</code>\n")
		}

	case event.TypeSystemWake:
		if ev.WakeType != "" {
			b.WriteString("⚡ Uyanış: ")
			b.WriteString(SafeText(ev.WakeType))
			b.WriteString("\n")
		}

	case event.TypeNetworkUp, event.TypeNetworkChanged:
		if ev.NetworkSSID != "" {
			b.WriteString("📶 Ağ: <b>")
			b.WriteString(SafeText(ev.NetworkSSID))
			b.WriteString("</b>")
			if ev.NetworkType != "" {
				b.WriteString(" (")
				b.WriteString(SafeText(ev.NetworkType))
				b.WriteString(")")
			}
			b.WriteString("\n")
		}
		if ev.LocalIP != "" {
			b.WriteString("🔌 IP: <code>")
			b.WriteString(SafeText(ev.LocalIP))
			b.WriteString("</code>\n")
		}
	}

	// Extra fields
	for k, v := range ev.Extra {
		b.WriteString("  • ")
		b.WriteString(SafeText(k))
		b.WriteString(": ")
		b.WriteString(SafeText(v))
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n<i>🏰 Kanije Kalesi v")
	b.WriteString(version)
	b.WriteString("</i>")

	return b.String()
}

// FormatHeartbeat creates the periodic status message.
func FormatHeartbeat(uptime time.Duration, diskFree, diskTotal uint64, evCount int64, platform string) string {
	var b strings.Builder

	b.WriteString("💓 <b>Sistem Nabzı</b>\n\n")
	b.WriteString("⏱️ Çalışma süresi: <b>")
	b.WriteString(formatDuration(uptime))
	b.WriteString("</b>\n")

	b.WriteString("🖥️ Platform: <code>")
	b.WriteString(SafeText(platform))
	b.WriteString("</code>\n")

	if diskTotal > 0 {
		usedPct := float64(diskTotal-diskFree) / float64(diskTotal) * 100
		b.WriteString("💿 Disk: ")
		b.WriteString(formatBytes(int64(diskFree)))
		b.WriteString(" boş / ")
		b.WriteString(formatBytes(int64(diskTotal)))
		b.WriteString(fmt.Sprintf(" (%%%.0f dolu)\n", usedPct))
	}

	b.WriteString("📊 Toplam olay: <b>")
	b.WriteString(fmt.Sprintf("%d", evCount))
	b.WriteString("</b>\n")

	b.WriteString("\n<i>🏰 Kanije Kalesi v")
	b.WriteString(version)
	b.WriteString("</i>")

	return b.String()
}

// FormatStatus creates the /status command response.
func FormatStatus(info StatusInfo) string {
	var b strings.Builder

	b.WriteString("📊 <b>Sistem Durumu</b>\n\n")

	// System
	b.WriteString("🖥️ <b>Sistem</b>\n")
	b.WriteString(fmt.Sprintf("  CPU: <b>%.1f%%</b>\n", info.CPUPercent))
	b.WriteString(fmt.Sprintf("  RAM: <b>%.1f%%</b> (%s / %s)\n",
		info.MemPercent,
		formatBytes(int64(info.MemUsed)),
		formatBytes(int64(info.MemTotal))))

	if len(info.Disks) > 0 {
		b.WriteString("  Disk:\n")
		for _, d := range info.Disks {
			b.WriteString(fmt.Sprintf("    <code>%s</code> %s boş / %s\n",
				SafeText(d.Path),
				formatBytes(int64(d.Free)),
				formatBytes(int64(d.Total))))
		}
	}

	b.WriteString("  Çalışma süresi: <b>")
	b.WriteString(formatDuration(info.Uptime))
	b.WriteString("</b>\n\n")

	// Bus stats
	b.WriteString("📨 <b>Olaylar</b>\n")
	b.WriteString(fmt.Sprintf("  Toplam: <b>%d</b>\n", info.BusReceived))
	b.WriteString(fmt.Sprintf("  Düşürülen: %d\n", info.BusDropped))

	if info.LastEvent != nil {
		b.WriteString("  Son olay: <i>")
		b.WriteString(eventLabel(info.LastEvent.Type))
		b.WriteString(" (")
		b.WriteString(info.LastEvent.Timestamp.Format("15:04:05"))
		b.WriteString(")</i>\n")
	}

	b.WriteString("\n<i>🏰 Kanije Kalesi v")
	b.WriteString(version)
	b.WriteString("</i>")

	return b.String()
}

// FormatRecentEvents formats the last N events for the /olaylar command.
func FormatRecentEvents(events []event.Event) string {
	if len(events) == 0 {
		return "📭 Henüz kaydedilmiş olay yok."
	}

	var b strings.Builder
	b.WriteString("📋 <b>Son Olaylar</b>\n\n")

	for i, ev := range events {
		b.WriteString(fmt.Sprintf("%d. %s <b>%s</b>",
			i+1,
			eventEmoji(ev.Type),
			eventLabel(ev.Type)))

		if ev.Username != "" {
			b.WriteString(" — ")
			b.WriteString(SafeText(ev.Username))
		}

		b.WriteString("\n   <i>")
		b.WriteString(ev.Timestamp.Format("02.01 15:04:05"))
		b.WriteString("</i>\n")
	}

	return b.String()
}

// FormatHelp generates the /yardim command response.
func FormatHelp() string {
	return `🏰 <b>Kanije Kalesi — Komut Listesi</b>

<b>📊 İzleme</b>
/status — Sistem durumu (CPU, RAM, disk)
/olaylar — Son 10 güvenlik olayı
/ping — Bağlantı kontrolü

<b>📷 Medya</b>
/foto — Kameradan anlık fotoğraf
/ekran — Ekran görüntüsü

<b>⚙️ Yönetim</b>
/kilitle — Ekranı kilitle
/yeniden — Sistemi yeniden başlat
/kapat — Sistemi kapat
/iptal — Bekleyen işlemi iptal et

<b>🔧 Ayarlar</b>
/ayarlar — Mevcut yapılandırma
/kurulum — Etkileşimli kurulum menüsü

<i>🏰 Kanije Kalesi — Siber kale muhafızı</i>`
}

// StatusInfo is the data bag for FormatStatus.
type StatusInfo struct {
	CPUPercent  float64
	MemPercent  float64
	MemUsed     uint64
	MemTotal    uint64
	Disks       []DiskInfo
	Uptime      time.Duration
	BusReceived int64
	BusDropped  int64
	LastEvent   *event.Event
}

// DiskInfo holds disk usage for one mount point.
type DiskInfo struct {
	Path  string
	Free  uint64
	Total uint64
}

// ---- Helpers ----

func eventEmoji(t event.Type) string {
	switch t {
	case event.TypeLoginSuccess:
		return "✅"
	case event.TypeLoginFailed:
		return "🚨"
	case event.TypeLogoff:
		return "👋"
	case event.TypeScreenLock:
		return "🔒"
	case event.TypeScreenUnlock:
		return "🔓"
	case event.TypeSystemBoot:
		return "🖥️"
	case event.TypeSystemShutdown:
		return "🔴"
	case event.TypeSystemSleep:
		return "😴"
	case event.TypeSystemWake:
		return "☀️"
	case event.TypeUSBInserted:
		return "🔌"
	case event.TypeUSBRemoved:
		return "⏏️"
	case event.TypeNetworkUp:
		return "🌐"
	case event.TypeNetworkDown:
		return "📡"
	case event.TypeNetworkChanged:
		return "🔄"
	case event.TypeHeartbeat:
		return "💓"
	default:
		return "📌"
	}
}

func eventLabel(t event.Type) string {
	switch t {
	case event.TypeLoginSuccess:
		return "Başarılı Giriş"
	case event.TypeLoginFailed:
		return "Başarısız Giriş Denemesi"
	case event.TypeLogoff:
		return "Oturum Kapatıldı"
	case event.TypeScreenLock:
		return "Ekran Kilitlendi"
	case event.TypeScreenUnlock:
		return "Ekran Kilidi Açıldı"
	case event.TypeSystemBoot:
		return "Sistem Başlatıldı"
	case event.TypeSystemShutdown:
		return "Sistem Kapatılıyor"
	case event.TypeSystemSleep:
		return "Uyku Moduna Girdi"
	case event.TypeSystemWake:
		return "Uykudan Uyandı"
	case event.TypeUSBInserted:
		return "USB Cihazı Takıldı"
	case event.TypeUSBRemoved:
		return "USB Cihazı Çıkarıldı"
	case event.TypeNetworkUp:
		return "İnternet Bağlantısı Kuruldu"
	case event.TypeNetworkDown:
		return "İnternet Bağlantısı Kesildi"
	case event.TypeNetworkChanged:
		return "Ağ Değişti"
	case event.TypeHeartbeat:
		return "Sistem Nabzı"
	default:
		return string(t)
	}
}

func formatBytes(b int64) string {
	if b < 0 {
		return "?"
	}
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	days := hours / 24
	hours = hours % 24

	if days > 0 {
		return fmt.Sprintf("%dg %ds %dm", days, hours, minutes)
	}
	return fmt.Sprintf("%ds %dm", hours, minutes)
}

// Ensure math is used (log2 for formatBytes)
var _ = math.Log2
var _ = runtime.GOOS
