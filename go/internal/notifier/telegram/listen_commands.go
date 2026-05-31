package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	listenSegmentSec = 20 // length of each streamed audio chunk
	listenDefaultMin = 5  // default session length
	listenMaxMinutes = 60 // hard cap
)

// cmdDinle starts or stops live audio "streaming": the mic is recorded in
// listenSegmentSec chunks and each chunk is sent to the owner until the time runs
// out or /dinle kapat is issued.
//
//	/dinle [dakika]   → başlat (varsayılan 5 dk, en çok 60)
//	/dinle kapat      → durdur
func (b *Bot) cmdDinle(ctx context.Context, chatID int64, text string) {
	if b.captureAudio == nil {
		b.reply(ctx, chatID, "❌ Ses kaydı bu derlemede yok.")
		return
	}

	arg := strings.ToLower(commandArg(text))
	if arg == "kapat" || arg == "kapa" || arg == "dur" || arg == "stop" || arg == "off" {
		b.stopListen(ctx, chatID)
		return
	}

	minutes := listenDefaultMin
	if arg != "" {
		n, err := strconv.Atoi(arg)
		if err != nil || n < 1 {
			b.reply(ctx, chatID, "❌ Geçersiz süre. Örnek: <code>/dinle 10</code> (dakika) · durdur: <code>/dinle kapat</code>")
			return
		}
		minutes = n
	}
	if minutes > listenMaxMinutes {
		minutes = listenMaxMinutes
	}

	b.startListen(chatID, minutes)
}

// startListen launches a listening session, replacing any session already running.
func (b *Bot) startListen(chatID int64, minutes int) {
	lctx, cancel := context.WithTimeout(context.Background(), time.Duration(minutes)*time.Minute)

	b.listenMu.Lock()
	if b.listenCancel != nil {
		b.listenCancel() // stop the previous session first
	}
	b.listenCancel = cancel
	b.listenMu.Unlock()

	b.reply(context.Background(), chatID, fmt.Sprintf(
		"🎧 <b>Canlı dinleme başladı</b>\nSüre: %d dk · ~%d sn'lik parçalar halinde gelecek.\nDurdurmak için: <code>/dinle kapat</code>",
		minutes, listenSegmentSec))

	go b.listenLoop(lctx, cancel, chatID)
}

// listenLoop records and sends audio chunks until ctx ends (timeout or stop).
func (b *Bot) listenLoop(ctx context.Context, cancel context.CancelFunc, chatID int64) {
	defer cancel() // release the context's resources on exit

	count := 0
	for ctx.Err() == nil {
		data, err := b.captureAudio(ctx, listenSegmentSec)
		if err != nil {
			if ctx.Err() != nil {
				break // stopped/expired mid-segment — expected
			}
			b.reply(ctx, chatID, "⚠️ Dinleme parçası alınamadı: "+safeHTML(err.Error()))
			break
		}
		count++
		caption := fmt.Sprintf("🎧 Canlı dinleme #%d", count)
		if err := b.client.SendAudio(ctx, chatID, data, fmt.Sprintf("dinle_%d.mp3", count), caption); err != nil {
			if ctx.Err() != nil {
				break
			}
			b.log.Warn("canlı dinleme parçası gönderilemedi", "err", err)
		}
	}

	// Clear the session pointer if it's still ours.
	b.listenMu.Lock()
	if b.listenCancel != nil {
		b.listenCancel = nil
	}
	b.listenMu.Unlock()

	b.reply(context.Background(), chatID, fmt.Sprintf("🎧 Canlı dinleme bitti (%d parça gönderildi).", count))
}

// stopListen cancels the active listening session, if any.
func (b *Bot) stopListen(ctx context.Context, chatID int64) {
	b.listenMu.Lock()
	cancel := b.listenCancel
	b.listenMu.Unlock()

	if cancel == nil {
		b.reply(ctx, chatID, "ℹ️ Aktif canlı dinleme yok.")
		return
	}
	cancel()
	b.reply(ctx, chatID, "🎧 Canlı dinleme durduruluyor…")
}
