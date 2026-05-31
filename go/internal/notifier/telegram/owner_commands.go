package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// This file groups the owner-only "last resort" commands — /kaldir (clean
// removal), /aktar (hand-off) and /imha (wipe). Beyond the normal capability
// gate each ALSO requires the caller to be the delegation-tree root (the device
// owner): these powers are never delegated, even if the capability bit somehow
// gets set on a sub-user. Every one runs through the confirm → undo-window flow
// so a mis-tap is always recoverable.

// requireOwner replies and returns false unless chatID is the device owner (root
// of the access tree). Used to gate the irreversible commands above.
func (b *Bot) requireOwner(ctx context.Context, chatID int64) bool {
	if b.acl != nil && chatID == b.acl.RootID() {
		return true
	}
	b.reply(ctx, chatID, "⛔ Bu komut yalnızca <b>cihaz sahibine</b> özeldir.")
	return false
}

// cmdKaldir uninstalls the agent completely after owner confirmation and the
// standard 15s undo window: the scheduled task, all files, any older-version
// leftovers and every trace are erased, then the agent exits.
func (b *Bot) cmdKaldir(ctx context.Context, chatID int64) {
	if !b.requireOwner(ctx, chatID) {
		return
	}
	if b.uninstall == nil {
		b.reply(ctx, chatID, "❌ Kaldırma bu derlemede desteklenmiyor (yalnızca Windows).")
		return
	}
	b.requestDanger(ctx, chatID, "kaldir",
		"Kanije'yi tamamen kaldır — görev, dosyalar, eski sürümler ve tüm izler silinecek",
		undoWindowSystem,
		func(c context.Context) {
			b.reply(c, chatID, "🧹 Kaldırılıyor… İz bırakmadan siliniyorum. Hoşça kal! 🏰")
			if err := b.uninstall(c); err != nil {
				b.reply(c, chatID, "❌ Kaldırma başlatılamadı: "+safeHTML(err.Error()))
			}
		})
}

// cmdAktar hands the device to a new owner after owner confirmation and the 15s
// undo window. Usage: /aktar <yeni_chat_id> [yeni_bot_token]. The new owner gets
// full control; the current owner (and all sub-users) lose access.
func (b *Bot) cmdAktar(ctx context.Context, chatID int64, text string) {
	if !b.requireOwner(ctx, chatID) {
		return
	}
	if b.transfer == nil {
		b.reply(ctx, chatID, "❌ Devretme bu derlemede desteklenmiyor.")
		return
	}

	fields := strings.Fields(text)
	if len(fields) < 2 {
		b.reply(ctx, chatID, "🔁 <b>Devretme</b>\nKullanım: <code>/aktar &lt;yeni_chat_id&gt; [yeni_bot_token]</code>\n\n"+
			"Yeni sahip cihazın tam kontrolünü alır; senin (ve eklediklerinin) erişimi <b>kalkar</b>. "+
			"Bot token verirsen güvenlik için o mesajı sonra sil.")
		return
	}

	newID, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil || newID == 0 {
		b.reply(ctx, chatID, "❌ Geçersiz chat ID. Sayısal olmalı (örn. 123456789).")
		return
	}
	if newID == chatID {
		b.reply(ctx, chatID, "ℹ️ Bu zaten senin ID'in — devretmeye gerek yok.")
		return
	}
	newToken := ""
	if len(fields) > 2 {
		newToken = fields[2]
	}

	tokenNote := ""
	if newToken != "" {
		tokenNote = " + yeni bot token"
	}
	title := fmt.Sprintf("Botu yeni sahibe devret (%d)%s — kendi erişimini kaybedeceksin", newID, tokenNote)

	b.requestDanger(ctx, chatID, "aktar", title, undoWindowSystem,
		func(c context.Context) {
			b.reply(c, chatID, fmt.Sprintf("🔁 Devrediliyor… Yeni sahip: <code>%d</code>. Yeni kimlikle yeniden başlatılıyorum.", newID))
			if err := b.transfer(c, newID, newToken); err != nil {
				b.reply(c, chatID, "❌ Devir başarısız: "+safeHTML(err.Error()))
			}
		})
}

// cmdImha securely wipes the agent's own data plus the user's Music folder (and
// any extra configured targets). No factory reset — the OS stays intact. THREE
// gates protect it — the explicit "ONAYLA" keyword, the Onayla button, then the
// 15s undo window — on top of being owner-only.
func (b *Bot) cmdImha(ctx context.Context, chatID int64, text string) {
	if !b.requireOwner(ctx, chatID) {
		return
	}
	if b.destroy == nil {
		b.reply(ctx, chatID, "❌ İmha bu derlemede desteklenmiyor (yalnızca Windows).")
		return
	}

	// /imha hedef ... → manage the wipe target list (no confirmation needed).
	fields := strings.Fields(text)
	if len(fields) >= 2 && (strings.EqualFold(fields[1], "hedef") || strings.EqualFold(fields[1], "target")) {
		b.imhaTarget(ctx, chatID, fields)
		return
	}

	// First gate: the command must carry the explicit confirmation keyword.
	if !strings.EqualFold(strings.TrimSpace(commandArg(text)), "ONAYLA") {
		b.reply(ctx, chatID, "💥 <b>İMHA — GERİ DÖNÜŞ YOK</b>\n\n"+
			"Şunlar <b>kurtarılamaz</b> şekilde (üzerine yazılarak) silinir:\n"+
			"• Kanije verisi — bot token + olay geçmişi\n"+
			"• <b>Müzik klasörünün içi</b> (kritik dosyaların orada)\n"+
			"• Eklediğin ek hedefler (bkz. <code>/imha hedef</code>)\n\n"+
			"✅ İşletim sistemi <b>silinmez</b> — fabrika sıfırlama yok.\n\n"+
			"Onaylıyorsan: <code>/imha ONAYLA</code>")
		return
	}

	b.requestDanger(ctx, chatID, "imha",
		"💥 İMHA ET — Kanije verisi + Müzik klasörü + hedef klasörler güvenli silinecek (GERİ DÖNÜŞ YOK)",
		undoWindowSystem,
		func(c context.Context) {
			b.reply(c, chatID, "💥 İmha başlatılıyor… Hassas veriler ve hedef klasörler güvenli siliniyor. İşletim sistemi korunuyor.")
			if err := b.destroy(c); err != nil {
				b.reply(c, chatID, "❌ İmha başlatılamadı: "+safeHTML(err.Error()))
			}
		})
}

// imhaTarget manages the extra /imha wipe-target directories.
//
//	/imha hedef                → listele
//	/imha hedef ekle <yol>     → ekle
//	/imha hedef sil <yol>      → kaldır
func (b *Bot) imhaTarget(ctx context.Context, chatID int64, fields []string) {
	sub := ""
	if len(fields) >= 3 {
		sub = strings.ToLower(fields[2])
	}

	switch sub {
	case "ekle", "add":
		if len(fields) < 4 {
			b.reply(ctx, chatID, "Kullanım: <code>/imha hedef ekle &lt;yol&gt;</code>")
			return
		}
		path := strings.Join(fields[3:], " ")
		if err := b.cfg.AddWipeTarget(path); err != nil {
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
			return
		}
		b.reply(ctx, chatID, "✅ İmha hedefi eklendi: <code>"+safeHTML(path)+"</code>\nBu klasörün içi de <code>/imha ONAYLA</code> ile güvenli silinir.")

	case "sil", "kaldir", "remove", "del":
		if len(fields) < 4 {
			b.reply(ctx, chatID, "Kullanım: <code>/imha hedef sil &lt;yol&gt;</code>")
			return
		}
		path := strings.Join(fields[3:], " ")
		ok, err := b.cfg.RemoveWipeTarget(path)
		switch {
		case err != nil:
			b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
		case !ok:
			b.reply(ctx, chatID, "ℹ️ Bu hedef listede yok.")
		default:
			b.reply(ctx, chatID, "🗑️ İmha hedefi kaldırıldı: <code>"+safeHTML(path)+"</code>")
		}

	default: // listele
		var sb strings.Builder
		sb.WriteString("💥 <b>İmha Hedefleri</b>\n\n")
		sb.WriteString("• <b>Müzik klasörü</b> — içi (her zaman, otomatik)\n")
		sb.WriteString("• Kanije verisi — token + geçmiş (her zaman)\n")
		targets := b.cfg.WipeTargets()
		if len(targets) == 0 {
			sb.WriteString("\n<i>Ek hedef yok.</i>")
		} else {
			sb.WriteString("\n<b>Ek hedefler:</b>\n")
			for _, t := range targets {
				sb.WriteString("• <code>" + safeHTML(t) + "</code>\n")
			}
		}
		sb.WriteString("\n<i>Ekle: /imha hedef ekle &lt;yol&gt; · Sil: /imha hedef sil &lt;yol&gt;</i>")
		b.reply(ctx, chatID, sb.String())
	}
}
