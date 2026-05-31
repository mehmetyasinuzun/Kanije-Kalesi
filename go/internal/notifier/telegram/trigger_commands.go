package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// cmdTetikKamera toggles and configures motion-triggered camera (burst) mode.
// Subcommands: ac | kapat | esik <n> | seri <n>. Bare /tetikkamera shows status.
func (b *Bot) cmdTetikKamera(ctx context.Context, chatID int64, text string) {
	fields := strings.Fields(text)
	sub := ""
	if len(fields) > 1 {
		sub = strings.ToLower(fields[1])
	}

	switch sub {
	case "ac", "aç", "on", "baslat", "başlat":
		if err := b.cfg.SetMotionEnabled(true); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, b.tetikKameraStatus("📷 Hareket-tetikli kamera <b>AÇILDI</b> — hareket görülünce foto serisi çekilir."))

	case "kapat", "kapa", "off", "dur":
		if err := b.cfg.SetMotionEnabled(false); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, "📷 Hareket-tetikli kamera <b>KAPATILDI</b>.")

	case "esik", "eşik":
		if len(fields) < 3 {
			b.reply(ctx, chatID, "Kullanım: <code>/tetikkamera esik 12</code>\nDüşük değer = daha hassas (1–255).")
			return
		}
		v, err := strconv.ParseFloat(fields[2], 64)
		if err != nil || v <= 0 {
			b.reply(ctx, chatID, "❌ Geçersiz eşik. Örn: <code>/tetikkamera esik 12</code>")
			return
		}
		if err := b.cfg.SetMotionThreshold(v); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, fmt.Sprintf("✅ Hareket hassasiyet eşiği: <b>%.0f</b>", v))

	case "seri", "burst", "foto":
		if len(fields) < 3 {
			b.reply(ctx, chatID, "Kullanım: <code>/tetikkamera seri 3</code>\nHareket başına kaç foto çekilsin (1–10).")
			return
		}
		n, err := strconv.Atoi(fields[2])
		if err != nil || n < 1 {
			b.reply(ctx, chatID, "❌ Geçersiz sayı. Örn: <code>/tetikkamera seri 3</code>")
			return
		}
		if err := b.cfg.SetMotionBurst(n); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, fmt.Sprintf("✅ Hareket başına <b>%d</b> foto (seri).", b.cfg.MotionBurst()))

	default:
		b.reply(ctx, chatID, b.tetikKameraStatus(""))
	}
}

func (b *Bot) tetikKameraStatus(prefix string) string {
	iv, th := b.cfg.MotionSettings()
	var sb strings.Builder
	if prefix != "" {
		sb.WriteString(prefix)
		sb.WriteString("\n\n")
	}
	sb.WriteString("📷 <b>Hareket-Tetikli Kamera</b>\n")
	sb.WriteString(onOff("Durum", b.cfg.MotionEnabled()))
	sb.WriteString(fmt.Sprintf("⏱️ Tarama aralığı: <b>%.0f sn</b>\n", iv.Seconds()))
	sb.WriteString(fmt.Sprintf("🎚️ Hassasiyet eşiği: <b>%.0f</b> <i>(düşük = daha hassas)</i>\n", th))
	sb.WriteString(fmt.Sprintf("📸 Hareket başına: <b>%d</b> foto (seri)\n", b.cfg.MotionBurst()))
	sb.WriteString("\n<i>Komutlar:</i> /tetikkamera ac · kapat · esik &lt;n&gt; · seri &lt;n&gt;")
	return sb.String()
}

// cmdTetikSes toggles and configures voice-activated recording (VOX).
// Subcommands: ac | kapat | esik <dB>. Bare /tetikses shows status.
func (b *Bot) cmdTetikSes(ctx context.Context, chatID int64, text string) {
	fields := strings.Fields(text)
	sub := ""
	if len(fields) > 1 {
		sub = strings.ToLower(fields[1])
	}

	switch sub {
	case "ac", "aç", "on", "baslat", "başlat":
		if err := b.cfg.SetVOXEnabled(true); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, b.tetikSesStatus("🔊 Ses-tetikli kayıt <b>AÇILDI</b> — sessizken kayıt yok; ses gelince parça parça kaydeder."))

	case "kapat", "kapa", "off", "dur":
		if err := b.cfg.SetVOXEnabled(false); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, "🔊 Ses-tetikli kayıt <b>KAPATILDI</b>.")

	case "esik", "eşik":
		if len(fields) < 3 {
			b.reply(ctx, chatID, "Kullanım: <code>/tetikses esik -35</code>\nYüksek (sıfıra yakın) = daha hassas. Sessiz oda ~ -55 dB, konuşma ~ -25 dB.")
			return
		}
		v, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			b.reply(ctx, chatID, "❌ Geçersiz dB. Örn: <code>/tetikses esik -35</code>")
			return
		}
		if err := b.cfg.SetVOXThreshold(v); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, fmt.Sprintf("✅ Ses tetik eşiği: <b>%.0f dB</b>", v))

	default:
		b.reply(ctx, chatID, b.tetikSesStatus(""))
	}
}

func (b *Bot) tetikSesStatus(prefix string) string {
	thr, partSec, _, maxParts := b.cfg.VOXSettings()
	var sb strings.Builder
	if prefix != "" {
		sb.WriteString(prefix)
		sb.WriteString("\n\n")
	}
	sb.WriteString("🔊 <b>Ses-Tetikli Kayıt (VOX)</b>\n")
	sb.WriteString(onOff("Durum", b.cfg.VOXEnabled()))
	sb.WriteString(fmt.Sprintf("🎚️ Tetik eşiği: <b>%.0f dB</b> <i>(üstü = kayıt başlar)</i>\n", thr))
	sb.WriteString(fmt.Sprintf("🧩 Parça uzunluğu: <b>%d sn</b> · en çok <b>%d</b> parça\n", partSec, maxParts))
	sb.WriteString("\n<i>Ses devam ettikçe parça parça kaydeder, sessizlikte durur. Sessizken hiçbir şey kaydedilmez.</i>\n")
	sb.WriteString("<i>Komutlar:</i> /tetikses ac · kapat · esik &lt;dB&gt;")
	return sb.String()
}
