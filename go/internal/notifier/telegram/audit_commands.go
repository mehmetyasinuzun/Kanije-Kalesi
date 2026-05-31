package telegram

import (
	"context"
	"strconv"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/fileaudit"
)

// cmdErisim manages SACL file-access auditing — "who opened/copied my files".
//
//	/erisim kur <yol>     → klasörü denetime al (yönetici/SYSTEM gerekir)
//	/erisim durdur <yol>  → denetimi kaldır
//	/erisim liste         → denetlenen klasörler
//	/erisim [sayı]        → son erişimler (kim/ne zaman/hangi program/ne yaptı)
//
// Setup/teardown change the system audit policy, so they are owner-only; querying
// is open to CapFiles holders.
func (b *Bot) cmdErisim(ctx context.Context, chatID int64, text string) {
	fields := strings.Fields(text)
	sub := ""
	if len(fields) >= 2 {
		sub = strings.ToLower(fields[1])
	}

	switch sub {
	case "kur", "ekle", "add", "setup":
		if !b.requireOwner(ctx, chatID) {
			return
		}
		if len(fields) < 3 {
			b.reply(ctx, chatID, "Kullanım: <code>/erisim kur &lt;klasör_yolu&gt;</code>\nÖrn: <code>/erisim kur C:\\Users\\Yasin\\Music</code>")
			return
		}
		path := strings.Join(fields[2:], " ")
		b.reply(ctx, chatID, "🔧 Denetim kuruluyor… (yönetici/SYSTEM yetkisi gerekir)")
		if err := fileaudit.EnableAndAudit(path); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		if err := b.cfg.AddAuditPath(path); err != nil {
			b.log.Warn("denetim yolu kaydedilemedi", "err", err)
		}
		b.reply(ctx, chatID, "✅ Denetim açıldı: <code>"+safeHTML(path)+"</code>\nBu klasöre kim erişirse <code>/erisim</code> ile görürsün.")

	case "durdur", "kaldir", "sil", "stop", "remove":
		if !b.requireOwner(ctx, chatID) {
			return
		}
		if len(fields) < 3 {
			b.reply(ctx, chatID, "Kullanım: <code>/erisim durdur &lt;klasör_yolu&gt;</code>")
			return
		}
		path := strings.Join(fields[2:], " ")
		if err := fileaudit.DisableAudit(path); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		if _, err := b.cfg.RemoveAuditPath(path); err != nil {
			b.log.Warn("denetim yolu kaldırılamadı", "err", err)
		}
		b.reply(ctx, chatID, "🛑 Denetim kaldırıldı: <code>"+safeHTML(path)+"</code>")

	case "liste", "list":
		paths := b.cfg.AuditPaths()
		if len(paths) == 0 {
			b.reply(ctx, chatID, "📭 Denetlenen klasör yok.\nEkle: <code>/erisim kur &lt;yol&gt;</code>")
			return
		}
		var sb strings.Builder
		sb.WriteString("👁️ <b>Denetlenen Klasörler</b>\n\n")
		for _, p := range paths {
			sb.WriteString("• <code>" + safeHTML(p) + "</code>\n")
		}
		sb.WriteString("\n<i>Son erişimler: /erisim · Kaldır: /erisim durdur &lt;yol&gt;</i>")
		b.reply(ctx, chatID, sb.String())

	default: // /erisim [sayı] → son erişimler
		paths := b.cfg.AuditPaths()
		if len(paths) == 0 {
			b.reply(ctx, chatID, "ℹ️ Henüz denetlenen klasör yok.\nKur: <code>/erisim kur &lt;klasör_yolu&gt;</code> (yönetici gerekir)")
			return
		}
		limit := 15
		if sub != "" {
			if n, err := strconv.Atoi(sub); err == nil && n > 0 {
				limit = n
				if limit > 50 {
					limit = 50
				}
			}
		}
		b.reply(ctx, chatID, "👁️ Son erişimler okunuyor…")
		evs, err := fileaudit.RecentAccess(paths, limit)
		if err != nil {
			b.reply(ctx, chatID, "❌ Okunamadı: "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, FormatAccessEvents(evs))
	}
}
