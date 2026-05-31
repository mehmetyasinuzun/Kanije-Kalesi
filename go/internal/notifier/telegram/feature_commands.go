package telegram

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/clipboard"
	"github.com/kanije-kalesi/kanije/internal/defender"
	"github.com/kanije-kalesi/kanije/internal/event"
	"github.com/kanije-kalesi/kanije/internal/sysinfo"
)

// maxDocBytes caps /dosya downloads below Telegram's ~50 MB bot-upload limit.
const maxDocBytes = 45 * 1024 * 1024

// This file holds the M4 device-feature commands: /pil (battery), /pano
// (clipboard), /panik (evidence sweep), and (added below) /dosya.

// cmdPil reports the battery status.
func (b *Bot) cmdPil(ctx context.Context, chatID int64) {
	b.reply(ctx, chatID, FormatBattery(sysinfo.CollectBattery()))
}

// cmdDefender reports Microsoft Defender's posture: whether real-time protection
// is on, when it last scanned, signature freshness, and recent detections.
func (b *Bot) cmdDefender(ctx context.Context, chatID int64) {
	b.reply(ctx, chatID, "🛡️ Defender durumu okunuyor…")
	b.reply(ctx, chatID, FormatDefender(defender.GetStatus()))
}

// cmdPano reads and returns the current clipboard text.
func (b *Bot) cmdPano(ctx context.Context, chatID int64) {
	text, err := clipboard.ReadText()
	if err != nil {
		b.reply(ctx, chatID, "❌ Pano okunamadı: "+safeHTML(err.Error()))
		return
	}
	text = strings.TrimRight(text, "\x00")
	if strings.TrimSpace(text) == "" {
		b.reply(ctx, chatID, "📋 Pano boş veya metin içermiyor.")
		return
	}

	const maxRunes = 3500 // keep within Telegram's message limit
	note := ""
	if r := []rune(text); len(r) > maxRunes {
		text = string(r[:maxRunes])
		note = "\n<i>(uzun içerik kısaltıldı)</i>"
	}
	b.reply(ctx, chatID, "📋 <b>Pano içeriği</b>\n<pre>"+safeHTML(text)+"</pre>"+note)
}

// cmdPanik is the one-shot panic sweep: it gathers a camera photo, a screenshot,
// the public IP, and a short audio clip and sends them to the owner, then records
// a panic event. With the "kilit"/"lock" argument it also locks the screen.
// Every step is best-effort — one failure never stops the rest.
func (b *Bot) cmdPanik(ctx context.Context, chatID int64, text string) {
	b.reply(ctx, chatID, "🆘 <b>PANİK MODU</b>\nKanıt toplanıyor — fotoğraf, ekran, dış IP ve ses birazdan gelecek…")

	// 1) Camera photo (fast).
	if b.capturePhoto != nil {
		if data, err := b.capturePhoto(ctx); err == nil {
			b.client.SendPhoto(ctx, chatID, data, "🆘 Panik — kamera")
		} else {
			b.reply(ctx, chatID, "⚠️ Kamera alınamadı: "+safeHTML(err.Error()))
		}
	}

	// 2) Screenshot (fast).
	if b.captureScreen != nil {
		if data, err := b.captureScreen(ctx); err == nil {
			b.client.SendPhoto(ctx, chatID, data, "🆘 Panik — ekran")
		} else {
			b.reply(ctx, chatID, "⚠️ Ekran alınamadı: "+safeHTML(err.Error()))
		}
	}

	// 3) Public IP (approximate location).
	publicIP := ""
	if b.getStatus != nil {
		if info := b.getStatus(); info.PublicIP != "" {
			publicIP = info.PublicIP
			b.reply(ctx, chatID, "🌍 Dış IP: <code>"+safeHTML(publicIP)+"</code>")
		}
	}

	// 4) Short audio (slower — kept last so the visuals arrive first).
	if b.captureAudio != nil {
		b.reply(ctx, chatID, "🎤 10 sn ses kaydı alınıyor…")
		if data, err := b.captureAudio(ctx, 10); err == nil {
			b.client.SendAudio(ctx, chatID, data, "panik.mp3", "🆘 Panik — ses")
		} else {
			b.reply(ctx, chatID, "⚠️ Ses alınamadı: "+safeHTML(err.Error()))
		}
	}

	// 5) Optional screen lock — only when the FIRST argument is exactly
	// "kilit"/"lock", not anywhere in the text (so "/panik şu kilitli oda" can't
	// accidentally lock the screen).
	if fields := strings.Fields(text); len(fields) > 1 && b.lockScreen != nil {
		if a := strings.ToLower(fields[1]); a == "kilit" || a == "lock" {
			if err := b.lockScreen(); err == nil {
				b.reply(ctx, chatID, "🔒 Ekran kilitlendi.")
			}
		}
	}

	// Persist a panic event so it shows up in /olaylar.
	ev := event.New(event.TypePanic, "panik")
	ev.Hostname, _ = os.Hostname()
	ev.PublicIP = publicIP
	if err := b.store.SaveEvent(ctx, ev); err != nil {
		b.log.Warn("panik olayı kaydedilemedi", "err", err)
	}

	b.reply(ctx, chatID, "✅ Panik kanıtı toplandı.")
}

// cmdDosya browses the filesystem and downloads files.
//
//	/dosya <yol>      → list a directory or show file info
//	/dosya al <yol>   → send the file to Telegram
func (b *Bot) cmdDosya(ctx context.Context, chatID int64, text string) {
	rest := commandRest(text)
	if rest == "" {
		b.reply(ctx, chatID, "📁 <b>Dosya</b>\nKullanım:\n"+
			"<code>/dosya &lt;yol&gt;</code> — dizini listele / dosya bilgisi\n"+
			"<code>/dosya al &lt;yol&gt;</code> — dosyayı indir")
		return
	}

	// "al <yol>" / "get <yol>" → download.
	if low := strings.ToLower(rest); strings.HasPrefix(low, "al ") || strings.HasPrefix(low, "get ") {
		path := strings.TrimSpace(rest[strings.IndexByte(rest, ' ')+1:])
		b.sendFile(ctx, chatID, path)
		return
	}
	b.listPath(ctx, chatID, rest)
}

// listPath shows a directory listing, or a single file's info.
func (b *Bot) listPath(ctx context.Context, chatID int64, path string) {
	info, err := os.Stat(path)
	if err != nil {
		b.reply(ctx, chatID, "❌ Bulunamadı: "+safeHTML(err.Error()))
		return
	}
	if !info.IsDir() {
		b.reply(ctx, chatID, fmt.Sprintf(
			"📄 <code>%s</code>\nBoyut: <b>%s</b>\nDeğişim: %s\n\nİndir: <code>/dosya al %s</code>",
			safeHTML(path), formatBytes(info.Size()),
			info.ModTime().Format("02.01.2006 15:04"), safeHTML(path)))
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		b.reply(ctx, chatID, "❌ Dizin okunamadı: "+safeHTML(err.Error()))
		return
	}

	var sb strings.Builder
	sb.WriteString("📂 <b>" + safeHTML(path) + "</b>\n\n")
	const maxEntries = 100
	for i, e := range entries {
		if i >= maxEntries {
			sb.WriteString(fmt.Sprintf("\n<i>…ve %d öğe daha</i>", len(entries)-maxEntries))
			break
		}
		if e.IsDir() {
			sb.WriteString("📁 " + safeHTML(e.Name()) + "/\n")
			continue
		}
		sb.WriteString("📄 " + safeHTML(e.Name()))
		if fi, err := e.Info(); err == nil {
			sb.WriteString(" <i>(" + formatBytes(fi.Size()) + ")</i>")
		}
		sb.WriteString("\n")
	}
	if len(entries) == 0 {
		sb.WriteString("<i>(boş dizin)</i>")
	}
	b.reply(ctx, chatID, sb.String())
}

// sendFile uploads a single file to Telegram, guarding the size limit.
func (b *Bot) sendFile(ctx context.Context, chatID int64, path string) {
	info, err := os.Stat(path)
	if err != nil {
		b.reply(ctx, chatID, "❌ Bulunamadı: "+safeHTML(err.Error()))
		return
	}
	if info.IsDir() {
		b.reply(ctx, chatID, "❌ Bu bir dizin. Listelemek için: <code>/dosya "+safeHTML(path)+"</code>")
		return
	}
	if info.Size() > maxDocBytes {
		b.reply(ctx, chatID, fmt.Sprintf("❌ Dosya çok büyük (%s). Telegram bot limiti ~50 MB.", formatBytes(info.Size())))
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		b.reply(ctx, chatID, "❌ Okunamadı: "+safeHTML(err.Error()))
		return
	}

	b.reply(ctx, chatID, "📤 Gönderiliyor: <code>"+safeHTML(filepath.Base(path))+"</code> ("+formatBytes(info.Size())+")")
	if err := b.client.SendDocument(ctx, chatID, data, filepath.Base(path), "📄 "+path); err != nil {
		b.reply(ctx, chatID, "❌ Gönderilemedi: "+safeHTML(err.Error()))
	}
}
