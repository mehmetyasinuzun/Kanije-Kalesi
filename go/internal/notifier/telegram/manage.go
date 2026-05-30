package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/access"
)

// This file implements the multi-user delegation UX: /yonetim (manage tree),
// /ekle (invite), /loglar (audit). All mutations re-verify capability and
// subtree visibility server-side inside the access.Manager — the callback data
// is never trusted to carry authority, only a target chat ID.

// ---- /yonetim : list the caller's subtree ----

func (b *Bot) cmdYonetim(ctx context.Context, chatID int64) {
	if b.acl == nil {
		b.reply(ctx, chatID, "❌ Erişim yönetimi kullanılamıyor.")
		return
	}
	text, kb := b.buildManageView(chatID)
	b.client.SendMessageWithKeyboard(ctx, chatID, text, kb)
}

// buildManageView renders the caller's own row plus its subtree, with a button
// per managed user and an "add" hint.
func (b *Bot) buildManageView(chatID int64) (string, InlineKeyboardMarkup) {
	var sb strings.Builder
	sb.WriteString("👥 <b>Kişi Yönetimi</b>\n\n")

	if me, ok := b.acl.Get(chatID); ok {
		sb.WriteString("👤 <b>Sen</b> — ")
		sb.WriteString(safeHTML(me.Name))
		sb.WriteString("\n<i>")
		sb.WriteString(capsSummary(me.Caps))
		sb.WriteString("</i>\n\n")
	}

	sub := b.acl.Subtree(chatID)
	var rows [][]InlineKeyboardButton
	if len(sub) == 0 {
		sb.WriteString("<i>Henüz kimseyi eklemedin.</i>\n")
	} else {
		sb.WriteString("<b>Eklediklerin (ve onların eklediği):</b>\n")
		for _, e := range sub {
			indent := strings.Repeat("　", e.Depth) // full-width space for tree depth
			branch := "├ "
			sb.WriteString(fmt.Sprintf("%s%s%s %s — <i>%s</i>\n",
				indent, branch, userIcon(e.User), safeHTML(e.User.Name), capsSummary(e.User.Caps)))
			rows = append(rows, []InlineKeyboardButton{{
				Text:         "⚙️ " + truncateLabel(e.User.Name, 24),
				CallbackData: "mng:user:" + strconv.FormatInt(e.User.ChatID, 10),
			}})
		}
	}

	sb.WriteString("\n<i>Eklemek için: /ekle &lt;chat_id&gt; [isim]</i>")
	rows = append(rows, []InlineKeyboardButton{{Text: "🔄 Yenile", CallbackData: "mng:list"}})
	return sb.String(), InlineKeyboardMarkup{InlineKeyboard: rows}
}

// ---- /ekle : invite a new user (least-privilege default, then fine-tune) ----

func (b *Bot) cmdEkle(ctx context.Context, chatID int64, text string) {
	if b.acl == nil {
		b.reply(ctx, chatID, "❌ Erişim yönetimi kullanılamıyor.")
		return
	}
	fields := strings.Fields(text)
	if len(fields) < 2 {
		b.reply(ctx, chatID, "Kullanım: <code>/ekle &lt;chat_id&gt; [isim]</code>\n\nKişinin chat ID'sini öğrenmek için ona @userinfobot'u kullandırabilirsin.")
		return
	}
	newID, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil || newID == 0 {
		b.reply(ctx, chatID, "❌ Geçersiz chat ID. Sayısal olmalı (örn. 123456789).")
		return
	}
	name := ""
	if len(fields) > 2 {
		name = strings.Join(fields[2:], " ")
	}

	// Invite with the least-privilege "izleyici" role (clamped to inviter), then
	// open the capability editor so the inviter can grant more if desired.
	u, err := b.acl.Invite(ctx, chatID, newID, name, access.RoleCaps("izleyici"))
	if err != nil {
		b.reply(ctx, chatID, "❌ "+safeHTML(err.Error()))
		return
	}
	text2, kb := b.buildCapsEditor(chatID, u.ChatID)
	b.client.SendMessageWithKeyboard(ctx,
		chatID,
		fmt.Sprintf("✅ <b>%s</b> eklendi (İzleyici).\nŞimdi yetkilerini ayarlayabilirsin:\n\n%s", safeHTML(u.Name), text2),
		kb)
}

// ---- /loglar : recent privileged actions ----

func (b *Bot) cmdLoglar(ctx context.Context, chatID int64) {
	if b.acl == nil {
		b.reply(ctx, chatID, "❌ Erişim yönetimi kullanılamıyor.")
		return
	}
	entries, err := b.acl.RecentAudit(ctx, 20)
	if err != nil {
		b.reply(ctx, chatID, "❌ Günlük alınamadı: "+safeHTML(err.Error()))
		return
	}
	if len(entries) == 0 {
		b.reply(ctx, chatID, "📭 Henüz kayıtlı işlem yok.")
		return
	}
	var sb strings.Builder
	sb.WriteString("🧾 <b>Son İşlemler</b>\n\n")
	for _, e := range entries {
		sb.WriteString(e.At.Local().Format("02.01 15:04"))
		sb.WriteString(" · ")
		sb.WriteString(safeHTML(actorOrID(e.Actor, e.ChatID)))
		sb.WriteString(" → <b>")
		sb.WriteString(safeHTML(e.Action))
		sb.WriteString("</b>")
		if e.Target != "" {
			sb.WriteString(" (")
			sb.WriteString(safeHTML(e.Target))
			sb.WriteString(")")
		}
		sb.WriteString("\n")
	}
	b.reply(ctx, chatID, sb.String())
}

// ---- callback router : mng:* ----

func (b *Bot) handleManageCallback(ctx context.Context, chatID, messageID int64, callbackID, data string) {
	b.client.AnswerCallbackQuery(ctx, callbackID, "")
	if b.acl == nil {
		return
	}
	// Top-level gate: you must hold at least one delegation power to use this UI.
	if !b.acl.Can(chatID, access.CapInvite) && !b.acl.Can(chatID, access.CapManage) {
		return
	}
	parts := strings.Split(data, ":") // mng:<action>[:target[:extra]]
	if len(parts) < 2 {
		return
	}

	switch parts[1] {
	case "list":
		text, kb := b.buildManageView(chatID)
		b.client.EditMessageText(ctx, chatID, messageID, text, &kb)

	case "user":
		target := parseID(parts, 2)
		if !b.canSeeTarget(chatID, target) {
			b.client.EditMessageText(ctx, chatID, messageID, "⛔ Bu kişiyi yönetme yetkin yok.", nil)
			return
		}
		text, kb := b.buildUserDetail(chatID, target)
		b.client.EditMessageText(ctx, chatID, messageID, text, &kb)

	case "caps":
		target := parseID(parts, 2)
		if !b.canEditTarget(chatID, target) {
			return
		}
		text, kb := b.buildCapsEditor(chatID, target)
		b.client.EditMessageText(ctx, chatID, messageID, text, &kb)

	case "ct": // toggle a capability: mng:ct:<target>:<cap>
		target := parseID(parts, 2)
		if len(parts) < 4 || !b.canEditTarget(chatID, target) {
			return
		}
		b.toggleCap(ctx, chatID, target, access.Capability(parts[3]))
		text, kb := b.buildCapsEditor(chatID, target)
		b.client.EditMessageText(ctx, chatID, messageID, text, &kb)

	case "role": // apply a preset role: mng:role:<target>:<roleKey>
		target := parseID(parts, 2)
		if len(parts) < 4 || !b.canEditTarget(chatID, target) {
			return
		}
		if err := b.acl.SetCaps(ctx, chatID, target, access.RoleCaps(parts[3])); err != nil {
			b.log.Warn("rol uygulanamadı", "hata", err)
		}
		text, kb := b.buildCapsEditor(chatID, target)
		b.client.EditMessageText(ctx, chatID, messageID, text, &kb)

	case "rm": // open removal menu (asks reparent if needed)
		target := parseID(parts, 2)
		if !b.canRemoveTarget(chatID, target) {
			return
		}
		text, kb := b.buildRemoveMenu(chatID, target)
		b.client.EditMessageText(ctx, chatID, messageID, text, &kb)

	case "rmgo": // execute removal: mng:rmgo:<target>:<mode>
		target := parseID(parts, 2)
		if len(parts) < 4 || !b.canRemoveTarget(chatID, target) {
			return
		}
		mode := reparentMode(parts[3])
		if err := b.acl.Remove(ctx, chatID, target, mode); err != nil {
			b.client.EditMessageText(ctx, chatID, messageID, "❌ "+safeHTML(err.Error()), nil)
			return
		}
		text, kb := b.buildManageView(chatID)
		b.client.EditMessageText(ctx, chatID, messageID, "✅ Çıkarıldı.\n\n"+text, &kb)
	}
}

// ---- views ----

func (b *Bot) buildUserDetail(viewer, target int64) (string, InlineKeyboardMarkup) {
	u, ok := b.acl.Get(target)
	if !ok {
		return "Kişi bulunamadı.", InlineKeyboardMarkup{}
	}
	var sb strings.Builder
	sb.WriteString("👤 <b>")
	sb.WriteString(safeHTML(u.Name))
	sb.WriteString("</b>\n")
	sb.WriteString(fmt.Sprintf("🆔 <code>%d</code>\n", u.ChatID))
	if inv, ok := b.acl.Get(u.InvitedBy); ok {
		sb.WriteString("➕ Ekleyen: ")
		sb.WriteString(safeHTML(inv.Name))
		sb.WriteString("\n")
	}
	sb.WriteString("\n<b>Yetkiler:</b>\n")
	if len(u.Caps) == 0 {
		sb.WriteString("<i>(hiç yetki yok)</i>\n")
	} else {
		for _, c := range u.Caps {
			sb.WriteString(c.Emoji() + " " + safeHTML(c.Label()) + "\n")
		}
	}

	idStr := strconv.FormatInt(target, 10)
	var rows [][]InlineKeyboardButton
	if b.acl.Can(viewer, access.CapInvite) {
		rows = append(rows, []InlineKeyboardButton{{Text: "✏️ Yetkileri düzenle", CallbackData: "mng:caps:" + idStr}})
	}
	if b.acl.Can(viewer, access.CapManage) {
		rows = append(rows, []InlineKeyboardButton{{Text: "🗑️ Çıkar", CallbackData: "mng:rm:" + idStr}})
	}
	rows = append(rows, []InlineKeyboardButton{{Text: "◀️ Geri", CallbackData: "mng:list"}})
	return sb.String(), InlineKeyboardMarkup{InlineKeyboard: rows}
}

// buildCapsEditor renders a toggle grid containing only the capabilities the
// editor itself holds (you cannot grant what you lack). ✅ = target has it.
func (b *Bot) buildCapsEditor(editor, target int64) (string, InlineKeyboardMarkup) {
	ed, _ := b.acl.Get(editor)
	u, ok := b.acl.Get(target)
	if !ok {
		return "Kişi bulunamadı.", InlineKeyboardMarkup{}
	}

	text := fmt.Sprintf("✏️ <b>%s — Yetkiler</b>\n<i>Hazır rol seç ya da tek tek aç/kapat. Yalnızca sende olan yetkileri verebilirsin; bir yetkiyi kısarsan bu kişinin eklediklerinde de otomatik kısılır.</i>", safeHTML(u.Name))

	idStr := strconv.FormatInt(target, 10)
	rows := [][]InlineKeyboardButton{{
		{Text: "👁️ İzleyici", CallbackData: "mng:role:" + idStr + ":izleyici"},
		{Text: "🛠️ Operatör", CallbackData: "mng:role:" + idStr + ":operator"},
		{Text: "👑 Yönetici", CallbackData: "mng:role:" + idStr + ":yonetici"},
	}}
	var row []InlineKeyboardButton
	for _, meta := range access.AllCaps {
		if !hasCapability(ed.Caps, meta.Cap) {
			continue // editor lacks this → cannot delegate it
		}
		mark := "⬜"
		if hasCapability(u.Caps, meta.Cap) {
			mark = "✅"
		}
		row = append(row, InlineKeyboardButton{
			Text:         mark + " " + meta.Label,
			CallbackData: "mng:ct:" + strconv.FormatInt(target, 10) + ":" + string(meta.Cap),
		})
		if len(row) == 2 {
			rows = append(rows, row)
			row = nil
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}
	rows = append(rows, []InlineKeyboardButton{{Text: "◀️ Geri", CallbackData: "mng:user:" + strconv.FormatInt(target, 10)}})
	return text, InlineKeyboardMarkup{InlineKeyboard: rows}
}

func (b *Bot) buildRemoveMenu(viewer, target int64) (string, InlineKeyboardMarkup) {
	u, ok := b.acl.Get(target)
	if !ok {
		return "Kişi bulunamadı.", InlineKeyboardMarkup{}
	}
	idStr := strconv.FormatInt(target, 10)

	if b.acl.HasChildren(target) {
		text := fmt.Sprintf("🗑️ <b>%s</b> çıkarılıyor.\n\nBu kişinin eklediği kişiler <b>silinmez</b> — kime bağlansın?", safeHTML(u.Name))
		kb := InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{
			{{Text: "⬆️ En tepe (varsayılan)", CallbackData: "mng:rmgo:" + idStr + ":root"}},
			{{Text: "↖️ Bir üst", CallbackData: "mng:rmgo:" + idStr + ":gp"}},
			{{Text: "👤 Bana bağla", CallbackData: "mng:rmgo:" + idStr + ":me"}},
			{{Text: "◀️ İptal", CallbackData: "mng:user:" + idStr}},
		}}
		return text, kb
	}

	text := fmt.Sprintf("🗑️ <b>%s</b> çıkarılsın mı?", safeHTML(u.Name))
	kb := InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{
		{{Text: "✅ Evet, çıkar", CallbackData: "mng:rmgo:" + idStr + ":root"}},
		{{Text: "◀️ İptal", CallbackData: "mng:user:" + idStr}},
	}}
	return text, kb
}

// ---- helpers ----

// Server-side gates. Two distinct delegation powers (per the user's model that
// "removal is its own permission"): CapInvite = add + edit your subtree's caps;
// CapManage = remove from your subtree. Both also require the target to be in
// the caller's subtree (visibility).
func (b *Bot) canSeeTarget(caller, target int64) bool {
	return target != 0 && b.acl.Visible(caller, target) &&
		(b.acl.Can(caller, access.CapInvite) || b.acl.Can(caller, access.CapManage))
}
func (b *Bot) canEditTarget(caller, target int64) bool {
	return target != 0 && b.acl.Can(caller, access.CapInvite) && b.acl.Visible(caller, target)
}
func (b *Bot) canRemoveTarget(caller, target int64) bool {
	return target != 0 && b.acl.Can(caller, access.CapManage) && b.acl.Visible(caller, target)
}

func (b *Bot) toggleCap(ctx context.Context, editor, target int64, cap access.Capability) {
	u, ok := b.acl.Get(target)
	if !ok {
		return
	}
	var next []access.Capability
	if hasCapability(u.Caps, cap) {
		for _, c := range u.Caps {
			if c != cap {
				next = append(next, c)
			}
		}
	} else {
		next = append(append([]access.Capability(nil), u.Caps...), cap)
	}
	if err := b.acl.SetCaps(ctx, editor, target, next); err != nil {
		b.log.Warn("yetki değiştirilemedi", "hata", err)
	}
}

func reparentMode(s string) access.ReparentMode {
	switch s {
	case "gp":
		return access.ReparentGrandparent
	case "me":
		return access.ReparentRemover
	default:
		return access.ReparentRoot
	}
}

func parseID(parts []string, idx int) int64 {
	if idx >= len(parts) {
		return 0
	}
	id, _ := strconv.ParseInt(parts[idx], 10, 64)
	return id
}

func hasCapability(caps []access.Capability, c access.Capability) bool {
	for _, x := range caps {
		if x == c {
			return true
		}
	}
	return false
}

func capsSummary(caps []access.Capability) string {
	if len(caps) == 0 {
		return "yetki yok"
	}
	if len(caps) >= len(access.AllCaps) {
		return "tüm yetkiler 👑"
	}
	parts := make([]string, 0, len(caps))
	for _, c := range caps {
		parts = append(parts, c.Emoji())
	}
	return strings.Join(parts, "")
}

func userIcon(u access.User) string {
	if hasCapability(u.Caps, access.CapInvite) {
		return "👥"
	}
	return "👤"
}

func actorOrID(name string, id int64) string {
	if name != "" {
		return name
	}
	return strconv.FormatInt(id, 10)
}

func truncateLabel(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "…"
}
