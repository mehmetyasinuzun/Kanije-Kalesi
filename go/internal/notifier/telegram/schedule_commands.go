package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// schedulableCmds is the whitelist of commands /zamanla may run automatically —
// non-interactive, owner-safe actions only.
var schedulableCmds = map[string]bool{
	"/foto": true, "/ekran": true, "/seskayit": true,
	"/status": true, "/pil": true, "/panik": true,
}

// cmdZamanla manages scheduled tasks:
//
//	/zamanla <aralık> <komut>   → ekle (örn: /zamanla 30dk /foto)
//	/zamanla liste              → listele
//	/zamanla sil <id>           → kaldır
func (b *Bot) cmdZamanla(ctx context.Context, chatID int64, text string) {
	if b.sched == nil {
		b.reply(ctx, chatID, "❌ Zamanlama bu derlemede yok.")
		return
	}

	fields := strings.Fields(text)
	if len(fields) < 2 {
		b.reply(ctx, chatID, scheduleUsage())
		return
	}

	switch strings.ToLower(fields[1]) {
	case "liste", "list":
		b.listSchedules(ctx, chatID)
		return
	case "sil", "del", "delete":
		if len(fields) < 3 {
			b.reply(ctx, chatID, "Kullanım: <code>/zamanla sil &lt;id&gt;</code>")
			return
		}
		id, err := strconv.Atoi(fields[2])
		if err != nil {
			b.reply(ctx, chatID, "❌ Geçersiz id.")
			return
		}
		ok, err := b.sched.Remove(id)
		switch {
		case err != nil:
			b.reply(ctx, chatID, "❌ Silinemedi: "+safeHTML(err.Error()))
		case !ok:
			b.reply(ctx, chatID, "ℹ️ Bu id'de görev yok.")
		default:
			b.reply(ctx, chatID, fmt.Sprintf("🗑️ Görev #%d silindi.", id))
		}
		return
	}

	// /zamanla <aralık> <komut>
	if len(fields) < 3 {
		b.reply(ctx, chatID, scheduleUsage())
		return
	}
	intervalSec, err := parseInterval(fields[1])
	if err != nil {
		b.reply(ctx, chatID, "❌ "+safeHTML(err.Error())+"\nÖrnek: <code>/zamanla 30dk /foto</code> (sn/dk/sa/g)")
		return
	}
	cmd := strings.ToLower(fields[2])
	if !strings.HasPrefix(cmd, "/") {
		cmd = "/" + cmd
	}
	if !schedulableCmds[cmd] {
		b.reply(ctx, chatID, "❌ <code>"+safeHTML(cmd)+"</code> zamanlanamaz.\nDesteklenen: /foto /ekran /seskayit /status /pil /panik")
		return
	}

	t, err := b.sched.Add(cmd, intervalSec, time.Now())
	if err != nil {
		b.reply(ctx, chatID, "❌ Eklenemedi: "+safeHTML(err.Error()))
		return
	}
	b.reply(ctx, chatID, fmt.Sprintf(
		"⏰ <b>Görev #%d eklendi</b>\nKomut: <code>%s</code>\nHer: %s\nİlk çalışma: %s",
		t.ID, safeHTML(cmd), formatInterval(intervalSec), t.NextRun.Local().Format("02.01 15:04")))
}

func (b *Bot) listSchedules(ctx context.Context, chatID int64) {
	tasks := b.sched.List()
	if len(tasks) == 0 {
		b.reply(ctx, chatID, "📭 Zamanlanmış görev yok.\nEkle: <code>/zamanla 30dk /foto</code>")
		return
	}
	var sb strings.Builder
	sb.WriteString("⏰ <b>Zamanlanmış Görevler</b>\n\n")
	for _, t := range tasks {
		sb.WriteString(fmt.Sprintf("<b>#%d</b> <code>%s</code> — her %s\n   <i>sonraki: %s</i>\n",
			t.ID, safeHTML(t.Command), formatInterval(t.IntervalSec), t.NextRun.Local().Format("02.01 15:04")))
	}
	sb.WriteString("\n<i>Sil: /zamanla sil &lt;id&gt;</i>")
	b.reply(ctx, chatID, sb.String())
}

// ExecScheduled runs a scheduled command as the owner. Called by the app's
// scheduler tick; only whitelisted, non-interactive commands are honored.
func (b *Bot) ExecScheduled(ctx context.Context, command string) {
	chatID := b.cfg.ChatID()
	if b.acl != nil && b.acl.RootID() != 0 {
		chatID = b.acl.RootID()
	}
	if chatID == 0 {
		return
	}

	b.reply(ctx, chatID, "⏰ Zamanlanmış: <code>"+safeHTML(command)+"</code>")
	switch command {
	case "/foto":
		b.cmdFoto(ctx, chatID)
	case "/ekran":
		b.cmdEkran(ctx, chatID)
	case "/seskayit":
		b.cmdSesKayit(ctx, chatID, "")
	case "/status":
		b.cmdStatus(ctx, chatID)
	case "/pil":
		b.cmdPil(ctx, chatID)
	case "/panik":
		b.cmdPanik(ctx, chatID, "")
	}
}

func scheduleUsage() string {
	return "⏰ <b>Zamanlama</b>\n" +
		"<code>/zamanla &lt;aralık&gt; &lt;komut&gt;</code> — ekle (örn: /zamanla 30dk /foto)\n" +
		"<code>/zamanla liste</code> — görevleri gör\n" +
		"<code>/zamanla sil &lt;id&gt;</code> — kaldır\n\n" +
		"<i>Aralık: 45sn · 30dk · 2sa · 1g · Komut: /foto /ekran /seskayit /status /pil /panik</i>"
}

// parseInterval converts "30dk"/"2sa"/"1g"/"45sn" (Turkish) or s/m/h/d to
// seconds, with a 60-second floor to avoid hammering. Two-letter Turkish units
// are matched before single-letter ones so "sa" never collides with "s".
func parseInterval(s string) (int, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	var unit int
	var num string
	switch {
	case strings.HasSuffix(s, "sn"):
		unit, num = 1, strings.TrimSuffix(s, "sn")
	case strings.HasSuffix(s, "dk"):
		unit, num = 60, strings.TrimSuffix(s, "dk")
	case strings.HasSuffix(s, "sa"):
		unit, num = 3600, strings.TrimSuffix(s, "sa")
	case strings.HasSuffix(s, "g"):
		unit, num = 86400, strings.TrimSuffix(s, "g")
	case strings.HasSuffix(s, "s"):
		unit, num = 1, strings.TrimSuffix(s, "s")
	case strings.HasSuffix(s, "m"):
		unit, num = 60, strings.TrimSuffix(s, "m")
	case strings.HasSuffix(s, "h"):
		unit, num = 3600, strings.TrimSuffix(s, "h")
	case strings.HasSuffix(s, "d"):
		unit, num = 86400, strings.TrimSuffix(s, "d")
	default:
		return 0, fmt.Errorf("birim belirt (sn/dk/sa/g)")
	}
	n, err := strconv.Atoi(num)
	if err != nil || n < 1 {
		return 0, fmt.Errorf("geçersiz sayı")
	}
	sec := n * unit
	if sec < 60 {
		return 0, fmt.Errorf("en az 60 saniye")
	}
	return sec, nil
}

// formatInterval renders a period in seconds as a short Turkish label.
func formatInterval(sec int) string {
	switch {
	case sec%86400 == 0:
		return fmt.Sprintf("%d gün", sec/86400)
	case sec%3600 == 0:
		return fmt.Sprintf("%d saat", sec/3600)
	case sec%60 == 0:
		return fmt.Sprintf("%d dk", sec/60)
	default:
		return fmt.Sprintf("%d sn", sec)
	}
}
