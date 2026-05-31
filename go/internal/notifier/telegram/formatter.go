package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/kanije-kalesi/kanije/internal/defender"
	"github.com/kanije-kalesi/kanije/internal/event"
	"github.com/kanije-kalesi/kanije/internal/fileaudit"
	"github.com/kanije-kalesi/kanije/internal/storage"
	"github.com/kanije-kalesi/kanije/internal/sysinfo"
)

// version is the application version shown in message footers. It is set once
// at startup via SetVersion, before the bot begins sending messages.
var version = "dev"

// SetVersion sets the version string shown in Telegram message footers.
// Call once during startup (before any goroutine sends a message).
func SetVersion(v string) {
	if v != "" {
		version = v
	}
}

// FormatEvent converts a security event into a Telegram HTML message.
// Every dynamic value goes through safeHTML() (UTF-8 repair + HTML escaping)
// before interpolation, so attacker-controlled fields such as a USB volume
// label or Wi-Fi SSID containing &, < or > cannot break the message or inject
// markup. The static <b>/<code> scaffolding is written literally.
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
		b.WriteString(safeHTML(ev.Hostname))
		b.WriteString("</code>")
		if ev.OS != "" {
			b.WriteString(" · ")
			b.WriteString(safeHTML(ev.OS))
		}
		b.WriteString("\n")
	}
	if ev.Username != "" {
		u := safeHTML(ev.Username)
		if ev.Domain != "" {
			u = safeHTML(ev.Domain) + `\` + u
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
			b.WriteString(safeHTML(ev.SourceIP))
			b.WriteString("</code>\n")
		}
		if ev.LogonType != 0 {
			b.WriteString("🔑 Giriş tipi: ")
			b.WriteString(safeHTML(ev.LogonType.String()))
			b.WriteString("\n")
		}

	case event.TypeUSBInserted:
		if ev.DeviceName != "" || ev.DeviceLabel != "" {
			label := ev.DeviceLabel
			if label == "" {
				label = ev.DeviceName
			}
			b.WriteString("💾 Cihaz: <b>")
			b.WriteString(safeHTML(label))
			b.WriteString("</b>\n")
		}
		if ev.DevicePath != "" {
			b.WriteString("📂 Yol: <code>")
			b.WriteString(safeHTML(ev.DevicePath))
			b.WriteString("</code>\n")
		}
		if ev.DeviceFS != "" {
			b.WriteString("🗂️ Dosya sistemi: ")
			b.WriteString(safeHTML(ev.DeviceFS))
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
			b.WriteString(safeHTML(ev.DeviceName))
			b.WriteString("</b>\n")
		}
		if ev.DevicePath != "" {
			b.WriteString("📂 Yol: <code>")
			b.WriteString(safeHTML(ev.DevicePath))
			b.WriteString("</code>\n")
		}

	case event.TypeSystemWake:
		if ev.WakeType != "" {
			b.WriteString("⚡ Uyanış: ")
			b.WriteString(safeHTML(ev.WakeType))
			b.WriteString("\n")
		}

	case event.TypeNetworkUp, event.TypeNetworkChanged:
		// Show the medium for EVERY connection type (USB tethering / Bluetooth /
		// Ethernet have no SSID), with the WiFi SSID appended when present.
		if ev.NetworkType != "" {
			b.WriteString("📡 Bağlantı: <b>")
			b.WriteString(safeHTML(ev.NetworkType))
			b.WriteString("</b>")
			if ev.NetworkSSID != "" {
				b.WriteString(" · 📶 ")
				b.WriteString(safeHTML(ev.NetworkSSID))
			}
			b.WriteString("\n")
		} else if ev.NetworkSSID != "" {
			b.WriteString("📶 Ağ: <b>")
			b.WriteString(safeHTML(ev.NetworkSSID))
			b.WriteString("</b>\n")
		}
		if ev.LocalIP != "" {
			b.WriteString("🔌 İç IP: <code>")
			b.WriteString(safeHTML(ev.LocalIP))
			b.WriteString("</code>\n")
		}
		if ev.PublicIP != "" {
			b.WriteString("🌍 Dış IP: <code>")
			b.WriteString(safeHTML(ev.PublicIP))
			b.WriteString("</code>\n")
		}
	}

	// Extra fields
	for k, v := range ev.Extra {
		b.WriteString("  • ")
		b.WriteString(safeHTML(k))
		b.WriteString(": ")
		b.WriteString(safeHTML(v))
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n<i>🏰 Kanije Kalesi ")
	b.WriteString(version)
	b.WriteString("</i>")

	return b.String()
}

// FormatHeartbeat creates the periodic status message.
func FormatHeartbeat(uptime time.Duration, diskFree, diskTotal uint64, evCount int64, hostname, osName, platform string) string {
	var b strings.Builder

	b.WriteString("💓 <b>Sistem Nabzı</b>\n\n")

	if hostname != "" {
		b.WriteString("💻 <code>")
		b.WriteString(safeHTML(hostname))
		b.WriteString("</code>")
		if osName != "" {
			b.WriteString(" · ")
			b.WriteString(safeHTML(osName))
		}
		b.WriteString("\n")
	}

	b.WriteString("⏱️ Çalışma süresi: <b>")
	b.WriteString(formatDuration(uptime))
	b.WriteString("</b>\n")

	b.WriteString("🖥️ Platform: <code>")
	b.WriteString(safeHTML(platform))
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

	b.WriteString("\n<i>🏰 Kanije Kalesi ")
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
				safeHTML(d.Path),
				formatBytes(int64(d.Free)),
				formatBytes(int64(d.Total))))
		}
	}

	b.WriteString("  Çalışma süresi: <b>")
	b.WriteString(formatDuration(info.Uptime))
	b.WriteString("</b>\n")
	if info.PublicIP != "" {
		b.WriteString("  🌍 Dış IP: <code>")
		b.WriteString(safeHTML(info.PublicIP))
		b.WriteString("</code>\n")
	}
	b.WriteString("\n")

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

	b.WriteString("\n<i>🏰 Kanije Kalesi ")
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
			b.WriteString(safeHTML(ev.Username))
		}

		b.WriteString("\n   <i>")
		b.WriteString(ev.Timestamp.Format("02.01 15:04:05"))
		b.WriteString("</i>")
		if ev.SourceIP != "" {
			b.WriteString(" · 🌐 <code>")
			b.WriteString(safeHTML(ev.SourceIP))
			b.WriteString("</code>")
		}
		if ev.DeviceLabel != "" {
			b.WriteString(" · 💾 ")
			b.WriteString(safeHTML(ev.DeviceLabel))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n<i>İpucu: aşağıdaki butonla bir olayın tüm detayını aç. /olaylar &lt;tip&gt; ile filtrele, /ozet ile özet.</i>")
	return b.String()
}

// FormatEventDetail renders every populated field of a single event for the
// /olaylar drill-down. This is a CTI tool — detail is the point.
func FormatEventDetail(ev event.Event) string {
	var b strings.Builder
	b.WriteString(ev.Type.Emoji())
	b.WriteString(" <b>")
	b.WriteString(safeHTML(ev.Type.Label()))
	b.WriteString("</b>\n🕐 ")
	b.WriteString(ev.Timestamp.Local().Format("02.01.2006 15:04:05"))
	b.WriteString("\n\n")

	row := func(label, val string) {
		if val == "" {
			return
		}
		b.WriteString(label)
		b.WriteString(": ")
		b.WriteString(safeHTML(val))
		b.WriteString("\n")
	}
	codeRow := func(label, val string) {
		if val == "" {
			return
		}
		b.WriteString(label)
		b.WriteString(": <code>")
		b.WriteString(safeHTML(val))
		b.WriteString("</code>\n")
	}

	row("🖥️ Cihaz", ev.Hostname)
	row("💿 OS", ev.OS)
	row("👤 Kullanıcı", ev.Username)
	row("🏢 Domain", ev.Domain)
	if ev.LogonType != 0 {
		row("🔑 Oturum tipi", ev.LogonType.String())
	}
	codeRow("🌐 Kaynak IP", ev.SourceIP)

	if ev.NetworkType != "" || ev.NetworkSSID != "" {
		row("📡 Bağlantı", networkLabelText(ev.NetworkType, ev.NetworkSSID))
	}
	codeRow("🔌 İç IP", ev.LocalIP)
	codeRow("🌍 Dış IP", ev.PublicIP)

	row("💾 Aygıt", ev.DeviceName)
	row("🏷️ Etiket", ev.DeviceLabel)
	if ev.DeviceSize > 0 {
		row("📦 Boyut", formatBytes(ev.DeviceSize))
	}
	row("🗂️ Dosya sistemi", ev.DeviceFS)
	codeRow("📂 Yol", ev.DevicePath)
	row("⚡ Uyanış", ev.WakeType)

	// Extras: geo location, interface, transition, etc.
	for k, v := range ev.Extra {
		row("• "+k, v)
	}

	b.WriteString(fmt.Sprintf("\n<i>Olay #%d</i>", ev.ID))
	return b.String()
}

// networkLabelText renders "medium · SSID" for the detail view.
func networkLabelText(medium, ssid string) string {
	switch {
	case medium != "" && ssid != "":
		return medium + " · " + ssid
	case ssid != "":
		return ssid
	default:
		return medium
	}
}

// FormatSummary renders a per-type event summary (the /ozet command).
func FormatSummary(counts []storage.TypeCount, total int64, days int) string {
	if total == 0 {
		return fmt.Sprintf("📭 Son %d günde kayıtlı olay yok.", days)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("📊 <b>Son %d Gün — Olay Özeti</b>\n\n", days))
	for _, c := range counts {
		b.WriteString(c.Type.Emoji())
		b.WriteString(" ")
		b.WriteString(safeHTML(c.Type.Label()))
		b.WriteString(": <b>")
		b.WriteString(fmt.Sprintf("%d", c.Count))
		b.WriteString("</b>\n")
	}
	b.WriteString(fmt.Sprintf("\n<b>Toplam: %d olay</b>", total))
	return b.String()
}

// FormatHelp generates the /yardim command response.
func FormatHelp() string {
	return `🏰 <b>Kanije Kalesi — Komut Listesi</b>

<b>📊 İzleme</b>
/status — Sistem durumu (CPU, RAM, disk)
/pil — Pil durumu (yüzde, şarj, kalan süre)
/olaylar — Son olaylar (/olaylar &lt;tip&gt; veya &lt;sayı&gt; ile filtrele)
/ozet — Son 7 günün olay özeti (tip bazlı)
/dogrula — Olay günlüğü bütünlüğünü doğrula
/defender — 🛡️ Microsoft Defender durumu + son taramalar + tespitler
/ping — Bağlantı kontrolü

<b>📷 Medya</b>
/foto — Kameradan anlık fotoğraf
/ekran — Ekran görüntüsü
/seskayit — Mikrofon kaydı (örn. /seskayit 30 → 30 sn · varsayılan 30 · en çok 600)
/pano — Panodaki metni getir
/panik — Tek komutla kanıt topla: foto + ekran + ses + dış IP (/panik kilit → ekranı da kilitler)
/hareket — Kamera hareket algılama (/hareket ac · kapat · esik &lt;n&gt;)
/dinle — Canlı dinleme (/dinle 10 → 10 dk parça parça ses · /dinle kapat)

<b>⚙️ Yönetim</b>
/dosya — Dosya gez / indir (/dosya &lt;yol&gt; · /dosya al &lt;yol&gt;)
/erisim — Dosyana kim erişti/kopyaladı (/erisim kur &lt;yol&gt; · /erisim · durdur) — yönetici gerekir
/zamanla — Komut zamanla (/zamanla 30dk /foto · /zamanla liste · sil &lt;id&gt;)
/terminal — Uzak komut çalıştır (/terminal whoami · cd kalıcı)
/terminalix — Yönetici terminali
/kilitle — Ekranı kilitle
/yeniden — Sistemi yeniden başlat
/kapat — Sistemi kapat
/iptal — Bekleyen işlemi iptal et
/guncelle — Yeni sürümü kontrol et ve kur

<b>🛡️ Sahip — Geri Dönüşsüz</b>
/kaldir — Kanije'yi izsiz kaldır (görev + dosya + eski sürümler)
/aktar — Botu yeni sahibe devret (/aktar &lt;chat_id&gt; [token])
/imha — Kanije verisi + Müzik klasörü içini güvenli sil (/imha ONAYLA · hedef ekle/sil) — OS silinmez
/koruma — Fiziksel tehdit koruması: dead-man switch · USB dead-man · yanlış-giriş · RAM-only (/koruma)
<i>Bu üçü yalnızca cihaz sahibine özeldir; çift onay + geri-alma ile korunur.</i>

<b>🔧 Ayarlar</b>
/ayarlar — Mevcut yapılandırma
/kurulum — Etkileşimli kurulum menüsü

<b>🛰️ Fleet (Çok Cihaz)</b>
/cihazlar — Tüm cihazları listele
<i>Grupta komut cihaz adıyla: /foto dizustu</i>

<b>👥 Çok Kullanıcı</b>
/ekle — Kişi ekle (/ekle &lt;chat_id&gt; [isim])
/yonetim — Eklediklerini gör/düzenle/çıkar
/loglar — Son işlemler (kim ne yaptı)

<i>Not: Her komut yetkiye bağlıdır. Eklediğin kişiye yalnızca sendeki yetkileri verebilirsin; kişi kendi alt ağacını görür.</i>

<i>🏰 Kanije Kalesi — Siber kale muhafızı</i>`
}

// FormatBattery renders the /pil battery snapshot.
func FormatBattery(b sysinfo.Battery) string {
	var sb strings.Builder
	sb.WriteString("🔋 <b>Pil Durumu</b>\n\n")

	if !b.Present {
		sb.WriteString("Bu cihazda pil algılanmadı (masaüstü olabilir).\n")
		if b.OnAC {
			sb.WriteString("🔌 AC güç bağlı.\n")
		}
		return sb.String()
	}

	if b.Percent >= 0 {
		sb.WriteString(fmt.Sprintf("%s <b>%%%d</b>\n", batteryEmoji(b.Percent), b.Percent))
	}
	switch {
	case b.Charging:
		sb.WriteString("⚡ Şarj oluyor\n")
	case b.OnAC:
		sb.WriteString("🔌 Fişe takılı\n")
	default:
		sb.WriteString("🔋 Pilden çalışıyor\n")
	}
	if b.Remaining > 0 {
		sb.WriteString("⏳ Tahmini kalan: <b>")
		sb.WriteString(formatDuration(b.Remaining))
		sb.WriteString("</b>\n")
	}
	return sb.String()
}

// FormatDefender renders the /defender Microsoft Defender status panel.
func FormatDefender(s defender.Status) string {
	if !s.Available {
		return "🛡️ <b>Microsoft Defender</b>\n\nDurum okunamadı — Defender yok, devre dışı, ya da bu platform desteklemiyor."
	}

	var b strings.Builder
	b.WriteString("🛡️ <b>Microsoft Defender</b>\n\n")
	b.WriteString(onOff("Gerçek zamanlı koruma", s.RealTimeProtection))
	b.WriteString(onOff("Antivirüs motoru", s.AntivirusEnabled))
	b.WriteString(onOff("Kurcalama koruması", s.TamperProtection))

	if s.LastQuickScan != "" || s.LastFullScan != "" {
		b.WriteString("\n")
		if s.LastQuickScan != "" {
			b.WriteString("⚡ Son hızlı tarama: <b>" + safeHTML(s.LastQuickScan) + "</b>\n")
		}
		if s.LastFullScan != "" {
			b.WriteString("🔍 Son tam tarama: <b>" + safeHTML(s.LastFullScan) + "</b>\n")
		}
	}
	if s.SignatureVersion != "" {
		b.WriteString("🧬 İmza: <code>" + safeHTML(s.SignatureVersion) + "</code>")
		if s.SignatureUpdated != "" {
			b.WriteString(" · " + safeHTML(s.SignatureUpdated))
		}
		b.WriteString("\n")
	}

	if len(s.RecentThreats) > 0 {
		b.WriteString("\n⚠️ <b>Son tespitler:</b>\n")
		for _, t := range s.RecentThreats {
			b.WriteString("• " + safeHTML(t.Name))
			if t.Time != "" {
				b.WriteString(" <i>(" + safeHTML(t.Time) + ")</i>")
			}
			b.WriteString("\n")
		}
	} else {
		b.WriteString("\n✅ Defender geçmişinde kayıtlı tehdit yok.\n")
	}
	return b.String()
}

// FormatAccessEvents renders the /erisim file-access list (who/when/what/where).
func FormatAccessEvents(evs []fileaudit.AccessEvent) string {
	if len(evs) == 0 {
		return "👁️ <b>Dosya Erişimleri</b>\n\nKayıt yok. (Denetimi yeni açtıysan, biri erişene kadar boş kalır.)"
	}
	var b strings.Builder
	b.WriteString("👁️ <b>Dosya Erişimleri</b> — kim · ne zaman · hangi program\n\n")
	for _, e := range evs {
		b.WriteString("🔹 <b>")
		b.WriteString(safeHTML(e.Access))
		b.WriteString("</b>")
		if e.User != "" {
			b.WriteString(" · 👤 ")
			b.WriteString(safeHTML(e.User))
		}
		b.WriteString("\n")
		if e.Process != "" {
			b.WriteString("   💻 <code>")
			b.WriteString(safeHTML(shortPath(e.Process)))
			b.WriteString("</code>\n")
		}
		if e.Object != "" {
			b.WriteString("   📄 <code>")
			b.WriteString(safeHTML(shortPath(e.Object)))
			b.WriteString("</code>\n")
		}
		if e.Time != "" {
			b.WriteString("   🕐 <i>")
			b.WriteString(safeHTML(e.Time))
			b.WriteString("</i>\n")
		}
	}
	return b.String()
}

// shortPath trims a long path for display, keeping the tail (most informative).
func shortPath(p string) string {
	const max = 52
	if len([]rune(p)) <= max {
		return p
	}
	r := []rune(p)
	return "…" + string(r[len(r)-(max-1):])
}

// onOff renders an enabled/disabled status line.
func onOff(label string, on bool) string {
	state := "🔴 kapalı"
	if on {
		state = "🟢 açık"
	}
	return label + ": <b>" + state + "</b>\n"
}

// batteryEmoji picks a color by charge level.
func batteryEmoji(pct int) string {
	switch {
	case pct >= 80:
		return "🟢"
	case pct >= 40:
		return "🟡"
	case pct >= 15:
		return "🟠"
	default:
		return "🔴"
	}
}

// FormatDeviceCard renders this device's identity line for /cihazlar. In a fleet
// group every device's bot answers, so the replies together list the whole fleet.
func FormatDeviceCard(label, osName string, info StatusInfo) string {
	if label == "" {
		label = "cihaz"
	}
	var b strings.Builder
	b.WriteString("🖥️ <b>")
	b.WriteString(safeHTML(label))
	b.WriteString("</b> — çevrimiçi ✅\n")
	if osName != "" {
		b.WriteString("💿 ")
		b.WriteString(safeHTML(osName))
		b.WriteString("\n")
	}
	b.WriteString("⏱️ Çalışma süresi: <b>")
	b.WriteString(formatDuration(info.Uptime))
	b.WriteString("</b>\n")
	if info.PublicIP != "" {
		b.WriteString("🌍 Dış IP: <code>")
		b.WriteString(safeHTML(info.PublicIP))
		b.WriteString("</code>\n")
	}
	b.WriteString("<i>Bu cihaza komut: /foto ")
	b.WriteString(safeHTML(label))
	b.WriteString("</i>")
	return b.String()
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
	PublicIP    string // external/public IP (best-effort)
}

// DiskInfo holds disk usage for one mount point.
type DiskInfo struct {
	Path  string
	Free  uint64
	Total uint64
}

// ---- Helpers ----

// eventEmoji and eventLabel delegate to the canonical definitions on event.Type
// so Telegram and other notifiers (webhooks) share one source of truth.
func eventEmoji(t event.Type) string { return t.Emoji() }

func eventLabel(t event.Type) string { return t.Label() }

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
