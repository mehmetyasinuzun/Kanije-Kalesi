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

// cmdImha is the nuclear option: secure-wipe the agent's data and factory-reset
// the device. THREE gates protect it — the explicit "ONAYLA" keyword in the
// command, then the Onayla button, then the 15s undo window — on top of being
// owner-only. This is the Find My / MDM "remote wipe" equivalent.
func (b *Bot) cmdImha(ctx context.Context, chatID int64, text string) {
	if !b.requireOwner(ctx, chatID) {
		return
	}
	if b.destroy == nil {
		b.reply(ctx, chatID, "❌ İmha bu derlemede desteklenmiyor (yalnızca Windows).")
		return
	}

	// First gate: the command must carry the explicit confirmation keyword.
	if !strings.EqualFold(strings.TrimSpace(commandArg(text)), "ONAYLA") {
		b.reply(ctx, chatID, "💥 <b>İMHA — GERİ DÖNÜŞ YOK</b>\n\n"+
			"Bu komut Kanije'nin tüm verisini (bot token, olay geçmişi) <b>kurtarılamaz</b> şekilde siler ve "+
			"cihazı <b>fabrika ayarlarına</b> sıfırlar (Windows: her şeyi kaldır).\n\n"+
			"Gerçekten istiyorsan şöyle yaz: <code>/imha ONAYLA</code>")
		return
	}

	b.requestDanger(ctx, chatID, "imha",
		"💥 CİHAZI İMHA ET — veriyi güvenli sil + fabrika sıfırlama (GERİ DÖNÜŞ YOK)",
		undoWindowSystem,
		func(c context.Context) {
			b.reply(c, chatID, "💥 İmha başlatılıyor… Hassas veriler güvenli siliniyor ve cihaz fabrika ayarlarına sıfırlanıyor.")
			if err := b.destroy(c); err != nil {
				b.reply(c, chatID, "❌ İmha başlatılamadı: "+safeHTML(err.Error()))
			}
		})
}
