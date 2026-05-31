package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kanije-kalesi/kanije/internal/config"
)

// cmdKoruma manages the physical-threat protection policy (owner-only — these
// settings are system-level self-defense). Subcommands:
//
//	/koruma                              → durum paneli
//	/koruma ac | kapat                   → ana şalter
//	/koruma deadman ac [saat] [aksiyon]  → ölü-adam anahtarı
//	/koruma usb ac [aksiyon] | kapat     → USB dead-man
//	/koruma giris ac [eşik] [aksiyon]    → yanlış-giriş sayacı
//	/koruma ramonly ac | kapat           → DB'yi RAM'de tut
//	/koruma checkin                      → manuel check-in (dead-man sıfırla)
//	/koruma iptal                        → bekleyen güvenli silmeyi durdur
//
// Action shorthands: kilit | alarm | sil.
func (b *Bot) cmdKoruma(ctx context.Context, chatID int64, text string) {
	if !b.requireOwner(ctx, chatID) {
		return
	}

	fields := strings.Fields(text)
	if len(fields) < 2 {
		b.reply(ctx, chatID, b.formatProtectionStatus())
		return
	}

	p := b.cfg.ProtectionPolicy()
	switch strings.ToLower(fields[1]) {
	case "ac", "aç", "on", "başlat", "baslat":
		p.Enabled = true
		b.applyProtection(ctx, chatID, p, "🛡️ <b>Koruma AÇIK.</b> Etkin tetikleyiciler çalışıyor.")
	case "kapat", "kapa", "off", "durdur":
		p.Enabled = false
		b.applyProtection(ctx, chatID, p, "🛡️ <b>Koruma KAPALI.</b>")
	case "iptal", "cancel":
		if b.cancelWipe != nil && b.cancelWipe() {
			b.reply(ctx, chatID, "✅ Bekleyen güvenli silme <b>iptal edildi</b>.")
		} else {
			b.reply(ctx, chatID, "ℹ️ Bekleyen güvenli silme yok.")
		}
	case "checkin", "buradayim", "buradayım":
		if b.checkIn != nil {
			b.checkIn()
		}
		b.reply(ctx, chatID, "✅ Check-in alındı — dead-man sayacı sıfırlandı.")
	case "deadman", "oluadam":
		b.korumaDeadman(ctx, chatID, p, fields)
	case "usb":
		b.korumaUSB(ctx, chatID, p, fields)
	case "giris", "giriş", "login":
		b.korumaLogin(ctx, chatID, p, fields)
	case "ramonly", "ram":
		b.korumaRAM(ctx, chatID, p, fields)
	default:
		b.reply(ctx, chatID, korumaUsage())
	}
}

func (b *Bot) korumaDeadman(ctx context.Context, chatID int64, p config.ProtectionConfig, fields []string) {
	if len(fields) < 3 {
		b.reply(ctx, chatID, "Kullanım: <code>/koruma deadman ac [saat] [kilit|alarm|sil]</code> · <code>kapat</code>")
		return
	}
	switch strings.ToLower(fields[2]) {
	case "ac", "aç", "on":
		p.DeadManEnabled = true
		if len(fields) >= 4 {
			if h, err := strconv.Atoi(fields[3]); err == nil && h > 0 {
				p.DeadManHours = h
			}
		}
		if len(fields) >= 5 {
			if a, ok := parseAction(fields[4]); ok {
				p.DeadManAction = a
			}
		}
		b.applyProtection(ctx, chatID, p, fmt.Sprintf(
			"⏳ <b>Dead-man AÇIK</b>: %d saat check-in gelmezse → %s\nCanlı tutmak için ara sıra <code>/koruma checkin</code>.",
			p.DeadManHours, actionTR(p.DeadManAction)))
	case "kapat", "off":
		p.DeadManEnabled = false
		b.applyProtection(ctx, chatID, p, "⏳ Dead-man KAPALI.")
	default:
		b.reply(ctx, chatID, "Kullanım: <code>/koruma deadman ac [saat] [kilit|alarm|sil]</code> · <code>kapat</code>")
	}
}

func (b *Bot) korumaUSB(ctx context.Context, chatID int64, p config.ProtectionConfig, fields []string) {
	if len(fields) < 3 {
		b.reply(ctx, chatID, "Kullanım: <code>/koruma usb ac [kilit|alarm|sil]</code> · <code>kapat</code>")
		return
	}
	switch strings.ToLower(fields[2]) {
	case "ac", "aç", "on":
		p.USBEnabled = true
		if len(fields) >= 4 {
			if a, ok := parseAction(fields[3]); ok {
				p.USBAction = a
			}
		}
		b.applyProtection(ctx, chatID, p, fmt.Sprintf(
			"🔌 <b>USB dead-man AÇIK</b>: bir USB çıkınca → %s", actionTR(p.USBAction)))
	case "kapat", "off":
		p.USBEnabled = false
		b.applyProtection(ctx, chatID, p, "🔌 USB dead-man KAPALI.")
	default:
		b.reply(ctx, chatID, "Kullanım: <code>/koruma usb ac [kilit|alarm|sil]</code> · <code>kapat</code>")
	}
}

func (b *Bot) korumaLogin(ctx context.Context, chatID int64, p config.ProtectionConfig, fields []string) {
	if len(fields) < 3 {
		b.reply(ctx, chatID, "Kullanım: <code>/koruma giris ac [eşik] [kilit|alarm|sil]</code> · <code>kapat</code>")
		return
	}
	switch strings.ToLower(fields[2]) {
	case "ac", "aç", "on":
		p.FailedLoginEnabled = true
		if len(fields) >= 4 {
			if n, err := strconv.Atoi(fields[3]); err == nil && n > 0 {
				p.FailedLoginThreshold = n
			}
		}
		if len(fields) >= 5 {
			if a, ok := parseAction(fields[4]); ok {
				p.FailedLoginAction = a
			}
		}
		b.applyProtection(ctx, chatID, p, fmt.Sprintf(
			"🔢 <b>Yanlış-giriş AÇIK</b>: %d başarısız denemede → %s", p.FailedLoginThreshold, actionTR(p.FailedLoginAction)))
	case "kapat", "off":
		p.FailedLoginEnabled = false
		b.applyProtection(ctx, chatID, p, "🔢 Yanlış-giriş KAPALI.")
	default:
		b.reply(ctx, chatID, "Kullanım: <code>/koruma giris ac [eşik] [kilit|alarm|sil]</code> · <code>kapat</code>")
	}
}

func (b *Bot) korumaRAM(ctx context.Context, chatID int64, p config.ProtectionConfig, fields []string) {
	if len(fields) < 3 {
		b.reply(ctx, chatID, "Kullanım: <code>/koruma ramonly ac</code> · <code>kapat</code>")
		return
	}
	switch strings.ToLower(fields[2]) {
	case "ac", "aç", "on":
		p.RAMOnly = true
		b.applyProtection(ctx, chatID, p, "💾 <b>RAM-only AÇIK.</b> Sonraki başlatmada olay DB'si diske yazılmaz (kapanınca silinir). Etkinleştirmek için <code>/yeniden</code>.")
	case "kapat", "off":
		p.RAMOnly = false
		b.applyProtection(ctx, chatID, p, "💾 RAM-only KAPALI. Sonraki başlatmada DB normal (diske) yazılır.")
	default:
		b.reply(ctx, chatID, "Kullanım: <code>/koruma ramonly ac</code> · <code>kapat</code>")
	}
}

// applyProtection persists a modified policy and confirms.
func (b *Bot) applyProtection(ctx context.Context, chatID int64, p config.ProtectionConfig, msg string) {
	if err := b.cfg.SetProtection(p); err != nil {
		b.reply(ctx, chatID, "❌ Kaydedilemedi: "+safeHTML(err.Error()))
		return
	}
	b.reply(ctx, chatID, msg)
}

func (b *Bot) formatProtectionStatus() string {
	p := b.cfg.ProtectionPolicy()
	var sb strings.Builder
	sb.WriteString("🛡️ <b>Koruma Politikası</b>\n")
	sb.WriteString("Ana durum: " + onOffTR(p.Enabled) + "\n\n")

	sb.WriteString("⏳ <b>Dead-man switch</b>: " + onOffTR(p.DeadManEnabled) + "\n")
	if p.DeadManEnabled {
		sb.WriteString(fmt.Sprintf("   %d saat check-in yoksa → %s\n", p.DeadManHours, actionTR(p.DeadManAction)))
		if b.lastCheckIn != nil {
			elapsed := time.Since(b.lastCheckIn())
			remaining := time.Duration(p.DeadManHours)*time.Hour - elapsed
			if remaining < 0 {
				remaining = 0
			}
			sb.WriteString("   Son check-in: " + formatDuration(elapsed) + " önce · Kalan: <b>" + formatDuration(remaining) + "</b>\n")
		}
	}

	sb.WriteString("🔌 <b>USB dead-man</b>: " + onOffTR(p.USBEnabled) + "\n")
	if p.USBEnabled {
		dev := p.USBDevice
		if dev == "" {
			dev = "herhangi USB"
		}
		sb.WriteString("   " + safeHTML(dev) + " çıkınca → " + actionTR(p.USBAction) + "\n")
	}

	sb.WriteString("🔢 <b>Yanlış-giriş</b>: " + onOffTR(p.FailedLoginEnabled) + "\n")
	if p.FailedLoginEnabled {
		sb.WriteString(fmt.Sprintf("   %d başarısız giriş → %s\n", p.FailedLoginThreshold, actionTR(p.FailedLoginAction)))
	}

	sb.WriteString("💾 <b>RAM-only mod</b>: " + onOffTR(p.RAMOnly) + "\n")
	sb.WriteString("\n<i>" + korumaUsage() + "</i>")
	return sb.String()
}

func korumaUsage() string {
	return "Komutlar:\n" +
		"/koruma ac · kapat\n" +
		"/koruma deadman ac [saat] [kilit|alarm|sil] · kapat\n" +
		"/koruma usb ac [kilit|alarm|sil] · kapat\n" +
		"/koruma giris ac [eşik] [kilit|alarm|sil] · kapat\n" +
		"/koruma ramonly ac · kapat\n" +
		"/koruma checkin · /koruma iptal"
}

// parseAction maps a user shorthand to a canonical action chain.
//
//	kilit → lock          (sadece kilitle)
//	alarm → lock_alert    (kilitle + alarm + foto)  [varsayılan, veri kaybı yok]
//	sil   → lock_alert_wipe (kilitle + alarm + güvenli sil)  [agresif]
func parseAction(s string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "kilit", "lock":
		return "lock", true
	case "alarm", "alert":
		return "lock_alert", true
	case "sil", "imha", "wipe":
		return "lock_alert_wipe", true
	default:
		return "", false
	}
}

// actionTR renders an action chain in Turkish.
func actionTR(action string) string {
	var parts []string
	if strings.Contains(action, "lock") {
		parts = append(parts, "Kilitle")
	}
	if strings.Contains(action, "alert") {
		parts = append(parts, "Alarm+Foto")
	}
	if strings.Contains(action, "wipe") {
		parts = append(parts, "Güvenli Sil")
	}
	if len(parts) == 0 {
		return action
	}
	return strings.Join(parts, " + ")
}

func onOffTR(on bool) string {
	if on {
		return "🟢 açık"
	}
	return "🔴 kapalı"
}
