package telegram

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/fileaudit"
)

// canaryDecoys are tempting fake files dropped into the honeypot folder. The NAME
// is the lure; content is generic. Opening/copying any of them trips the SACL
// watch → the configured canary action fires. We never read the intruder's data —
// we only detect that THEY touched OUR decoy.
var canaryDecoys = map[string]string{
	"Sifreler.txt":        "hesap parolalari\n",
	"Banka Bilgileri.txt": "iban / kart\n",
	"Ozel - GIZLI.txt":    "kisisel\n",
}

// cmdTuzak manages the canary honeypot (owner-only).
//
//	/tuzak kur [klasör]                 → tuzak dosyaları oluştur + denetime al
//	/tuzak aksiyon <alarm|kilitmodu|sil> → tetiklenince ne olsun
//	/tuzak kapat                        → denetimi kaldır
//	/tuzak                              → durum
func (b *Bot) cmdTuzak(ctx context.Context, chatID int64, text string) {
	if !b.requireOwner(ctx, chatID) {
		return
	}

	fields := strings.Fields(text)
	sub := ""
	if len(fields) >= 2 {
		sub = strings.ToLower(fields[1])
	}

	p := b.cfg.ProtectionPolicy()
	switch sub {
	case "kur", "ac", "aç", "setup":
		dir := defaultCanaryDir()
		if len(fields) >= 3 {
			dir = strings.Join(fields[2:], " ")
		}
		if dir == "" {
			b.reply(ctx, chatID, "❌ Klasör belirlenemedi. Örn: <code>/tuzak kur C:\\Users\\Yasin\\Desktop</code>")
			return
		}

		b.reply(ctx, chatID, "🍯 Tuzak kuruluyor… (denetim için yönetici/SYSTEM gerekir)")
		if err := createCanaryDecoys(dir); err != nil {
			b.reply(ctx, chatID, "❌ Tuzak dosyaları oluşturulamadı: "+safeHTML(err.Error()))
			return
		}
		if err := fileaudit.EnableAndAudit(dir); err != nil {
			b.reply(ctx, chatID, "❌ Denetim kurulamadı: "+safeHTML(err.Error()))
			return
		}
		p.CanaryEnabled = true
		p.CanaryPath = dir
		p.Enabled = true // ana koruma şalterini de aç
		if err := b.cfg.SetProtection(p); err != nil {
			b.reply(ctx, chatID, "❌ Kaydedilemedi: "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, "🍯 <b>Tuzak (honeypot) AKTİF</b>\n\n"+
			"📂 Klasör: <code>"+safeHTML(dir)+"</code>\n"+
			"📄 Tuzak dosyalar: <code>"+safeHTML(canaryFileList())+"</code>\n"+
			"🎬 Biri dokunursa → <b>"+actionTR(p.CanaryAction)+"</b>\n\n"+
			"<i>⚠️ Sen bu klasöre dokunma (tuzak senin). AV/sistem süreçleri yok sayılır — yalnızca kişi/şüpheli erişim tetikler.</i>\n"+
			"<i>Aksiyonu değiştir:</i> <code>/tuzak aksiyon kilitmodu</code>")

	case "kapat", "kaldir", "off":
		if p.CanaryPath != "" {
			_ = fileaudit.DisableAudit(p.CanaryPath)
		}
		p.CanaryEnabled = false
		if err := b.cfg.SetProtection(p); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, "🍯 Tuzak KAPALI. (Sahte dosyalar yerinde — istersen elle sil.)")

	case "aksiyon", "action":
		if len(fields) < 3 {
			b.reply(ctx, chatID, "Kullanım: <code>/tuzak aksiyon alarm|kilitmodu|sil</code>")
			return
		}
		a, ok := parseAction(fields[2])
		if !ok {
			b.reply(ctx, chatID, "❌ Geçersiz aksiyon. <b>alarm</b> · <b>kilitmodu</b> · <b>sil</b>")
			return
		}
		p.CanaryAction = a
		if err := b.cfg.SetProtection(p); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, "🍯 Tuzak aksiyonu: <b>"+actionTR(a)+"</b>")

	default:
		durum := "🔴 kapalı"
		if p.CanaryEnabled {
			durum = "🟢 açık"
		}
		msg := "🍯 <b>Tuzak (Honeypot)</b>\nDurum: <b>" + durum + "</b>\n"
		if p.CanaryEnabled {
			msg += "\n📂 Klasör: <code>" + safeHTML(p.CanaryPath) + "</code>\n" +
				"📄 Tuzak dosyalar: <code>" + safeHTML(canaryFileList()) + "</code>\n" +
				"🎬 Dokununca: <b>" + actionTR(p.CanaryAction) + "</b>\n"
		}
		msg += "\n<i>Dokun-kopyala:</i> <code>/tuzak kur</code> · <code>/tuzak aksiyon kilitmodu</code> · <code>/tuzak kapat</code>\n" +
			"<i>Gerçek honeypot: saldırganın verisini KOPYALAMAZ, yalnızca tuzağa dokunan kişi/programı yakalar (AV/sistem yok sayılır).</i>"
		b.reply(ctx, chatID, msg)
	}
}

// defaultCanaryDir returns a DEDICATED honeypot subfolder on the Desktop — never
// the whole Desktop, so normal Desktop activity (Telegram cache, media metadata,
// AV scans) can't trip the trap. The folder name is tempting enough to lure a
// snoop, isolated enough to stay quiet.
func defaultCanaryDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "Desktop", "Özel Belgeler")
}

// canaryFileList returns the decoy file names (for status/setup messages).
func canaryFileList() string {
	names := make([]string, 0, len(canaryDecoys))
	for n := range canaryDecoys {
		names = append(names, n)
	}
	return strings.Join(names, ", ")
}

// createCanaryDecoys drops the decoy files into dir (skips any that already exist).
func createCanaryDecoys(dir string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	for name, content := range canaryDecoys {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			continue // don't overwrite a real file with the same name
		}
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			return err
		}
	}
	return nil
}
