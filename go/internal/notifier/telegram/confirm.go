package telegram

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// This file implements the generic confirm → undo-window → run flow used by
// every dangerous command (/kapat, /yeniden, and later /kaldir, /aktar, /imha).
//
// Flow:
//  1. The command calls Bot.requestDanger, which sends an "Onayla / Vazgeç"
//     message and registers the op in phaseConfirm.
//  2. "Onayla" promotes the op to phaseUndo and starts the undo window; the
//     message turns into "uygulanacak — ↩️ Geri Al / /iptal".
//  3. When the window elapses the action's run callback fires exactly once.
//     "Geri Al", /iptal, or "Vazgeç" before that aborts it so run never fires.
//
// The dangerManager below is the client-free state core (owns the op map and
// timers); the Bot layers Telegram messaging on top via the onExpire/onFire
// hooks. Keeping state separate makes the whole flow testable without a network.

const (
	// confirmTTL is how long the initial Onayla/Vazgeç prompt stays valid before
	// it auto-expires (the user walked away without pressing anything).
	confirmTTL = 60 * time.Second

	// undoWindowSystem is the geri-alma window for /kapat and /yeniden.
	undoWindowSystem = 15 * time.Second
)

// dangerPhase tracks where an op is in the confirm → undo → run flow.
type dangerPhase int

const (
	phaseConfirm dangerPhase = iota // waiting for Onayla/Vazgeç
	phaseUndo                       // confirmed; runs when the undo window elapses
)

// dangerOp is one pending dangerous action. All fields are read/written only
// under dangerManager.mu (except run, which is immutable after creation).
type dangerOp struct {
	id         string
	kind       string // "kapat","yeniden","kaldir","aktar","imha"
	title      string // human title shown in messages, e.g. "Sistemi kapat"
	chatID     int64
	messageID  int64 // the confirm/undo message we keep editing in place
	phase      dangerPhase
	undoWindow time.Duration
	run        func(context.Context) // executed once, when the undo window elapses
	timer      *time.Timer           // current phase's timer (confirm TTL or undo window)
}

func (op *dangerOp) stopTimer() {
	if op.timer != nil {
		op.timer.Stop()
		op.timer = nil
	}
}

// dangerManager is the client-free state core of the confirm flow. It owns the
// op map and the per-phase timers. The hooks are invoked when a timer fires so
// the Bot can update the Telegram message; they are nil in unit tests.
type dangerManager struct {
	mu  sync.Mutex
	ops map[string]*dangerOp
	seq atomic.Uint64

	onExpire func(op *dangerOp) // confirm window elapsed, op never confirmed
	onFire   func(op *dangerOp) // undo window elapsed, about to run the action
}

func newDangerManager() *dangerManager {
	return &dangerManager{ops: make(map[string]*dangerOp)}
}

// add registers a new op in phaseConfirm and arms its confirm-expiry timer.
// Any prior pending op from the same chat is dropped so prompts never stack.
func (m *dangerManager) add(chatID int64, kind, title string, undoWindow time.Duration, run func(context.Context)) *dangerOp {
	id := "op" + strconv.FormatUint(m.seq.Add(1), 10)
	op := &dangerOp{
		id: id, kind: kind, title: title, chatID: chatID,
		phase: phaseConfirm, undoWindow: undoWindow, run: run,
	}

	m.mu.Lock()
	for oid, o := range m.ops {
		if o.chatID == chatID {
			o.stopTimer()
			delete(m.ops, oid)
		}
	}
	m.ops[id] = op
	op.timer = time.AfterFunc(confirmTTL, func() { m.expire(id) })
	m.mu.Unlock()
	return op
}

// setMessageID records the Telegram message ID once the prompt has been sent.
func (m *dangerManager) setMessageID(id string, messageID int64) {
	m.mu.Lock()
	if op, ok := m.ops[id]; ok {
		op.messageID = messageID
	}
	m.mu.Unlock()
}

// confirm promotes an op from phaseConfirm to phaseUndo and (re)arms the timer
// for the undo window. Returns the op and true on success; false if the op is
// gone or already confirmed (double-tap, expired).
func (m *dangerManager) confirm(id string) (*dangerOp, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	op, ok := m.ops[id]
	if !ok || op.phase != phaseConfirm {
		return nil, false
	}
	op.phase = phaseUndo
	op.stopTimer()
	op.timer = time.AfterFunc(op.undoWindow, func() { m.fire(id) })
	return op, true
}

// fire is the undo-window timer callback: it removes the op and runs its action
// exactly once. Guarded so a concurrent abort/cancel wins the race cleanly.
func (m *dangerManager) fire(id string) {
	m.mu.Lock()
	op, ok := m.ops[id]
	if !ok || op.phase != phaseUndo {
		m.mu.Unlock()
		return
	}
	delete(m.ops, id)
	m.mu.Unlock()

	if m.onFire != nil {
		m.onFire(op)
	}
	if op.run != nil {
		op.run(context.Background())
	}
}

// expire is the confirm-TTL timer callback: it drops an op the user never
// confirmed. No-op if the op was already confirmed or removed.
func (m *dangerManager) expire(id string) {
	m.mu.Lock()
	op, ok := m.ops[id]
	if !ok || op.phase != phaseConfirm {
		m.mu.Unlock()
		return
	}
	delete(m.ops, id)
	m.mu.Unlock()

	if m.onExpire != nil {
		m.onExpire(op)
	}
}

// abort removes an op by ID (Vazgeç / Geri Al) so its action never runs.
// Returns the removed op and true, or false if it was already gone.
func (m *dangerManager) abort(id string) (*dangerOp, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	op, ok := m.ops[id]
	if !ok {
		return nil, false
	}
	op.stopTimer()
	delete(m.ops, id)
	return op, true
}

// cancelChat aborts the newest pending op belonging to chatID (used by /iptal).
// Returns the canceled op and true, or false if the chat has none pending.
func (m *dangerManager) cancelChat(chatID int64) (*dangerOp, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var found *dangerOp
	for _, op := range m.ops {
		if op.chatID == chatID && (found == nil || opSeqNum(op.id) > opSeqNum(found.id)) {
			found = op
		}
	}
	if found == nil {
		return nil, false
	}
	found.stopTimer()
	delete(m.ops, found.id)
	return found, true
}

// ---- Bot integration (Telegram messaging on top of dangerManager) ----

// requestDanger starts the confirm flow for a dangerous action: it sends the
// Onayla/Vazgeç prompt and registers the op. The run callback fires only if the
// user confirms and lets the undo window elapse without canceling.
func (b *Bot) requestDanger(ctx context.Context, chatID int64, kind, title string, undoWindow time.Duration, run func(context.Context)) {
	op := b.danger.add(chatID, kind, title, undoWindow, run)

	text := "⚠️ <b>" + safeHTML(title) + "</b>\n\n" +
		"Onaylarsan <b>" + undoLabel(undoWindow) + "</b> geri alma süren olur " +
		"(↩️ Geri Al ya da /iptal). Onaylıyor musun?"
	msg, err := b.client.SendMessageWithKeyboard(ctx, chatID, text, confirmKeyboard(op.id))
	if err != nil || msg == nil {
		b.danger.abort(op.id) // prompt failed to show — make sure nothing runs
		b.log.Warn("onay istemi gönderilemedi", "kind", kind, "err", err)
		return
	}
	b.danger.setMessageID(op.id, msg.MessageID)
}

// handleDangerCallback routes the op:ok / op:no / op:undo inline-button presses.
func (b *Bot) handleDangerCallback(ctx context.Context, chatID, messageID int64, callbackID, data string) {
	parts := strings.Split(data, ":") // op:<action>:<id>
	if len(parts) != 3 {
		b.client.AnswerCallbackQuery(ctx, callbackID, "")
		return
	}
	id := parts[2]

	switch parts[1] {
	case "ok":
		op, ok := b.danger.confirm(id)
		if !ok {
			b.client.AnswerCallbackQuery(ctx, callbackID, "Süresi doldu")
			return
		}
		b.client.AnswerCallbackQuery(ctx, callbackID, "Onaylandı")
		text := "⏳ <b>" + safeHTML(op.title) + "</b>\n\n" +
			"<b>" + undoLabel(op.undoWindow) + "</b> içinde uygulanacak. " +
			"Vazgeçmek için ↩️ Geri Al ya da /iptal."
		kb := undoKeyboard(id)
		b.client.EditMessageText(ctx, chatID, op.messageID, text, &kb)

	case "no":
		op, ok := b.danger.abort(id)
		b.client.AnswerCallbackQuery(ctx, callbackID, "Vazgeçildi")
		if ok {
			b.client.EditMessageText(ctx, chatID, op.messageID,
				"✅ <b>Vazgeçildi.</b> Hiçbir şey yapılmadı.", nil)
		}

	case "undo":
		op, ok := b.danger.abort(id)
		b.client.AnswerCallbackQuery(ctx, callbackID, "Geri alındı")
		if ok {
			b.client.EditMessageText(ctx, chatID, op.messageID,
				"↩️ <b>Geri alındı.</b> İşlem uygulanmadı.", nil)
		}

	default:
		b.client.AnswerCallbackQuery(ctx, callbackID, "")
	}
}

// onDangerExpire updates the prompt when the confirm window lapses unconfirmed.
func (b *Bot) onDangerExpire(op *dangerOp) {
	b.client.EditMessageText(context.Background(), op.chatID, op.messageID,
		"⌛ <b>"+safeHTML(op.title)+"</b> — onay süresi doldu, iptal edildi.", nil)
}

// onDangerFire updates the message just before the action runs.
func (b *Bot) onDangerFire(op *dangerOp) {
	b.client.EditMessageText(context.Background(), op.chatID, op.messageID,
		"▶️ <b>"+safeHTML(op.title)+"</b> — uygulanıyor…", nil)
}

// ---- helpers ----

func confirmKeyboard(id string) InlineKeyboardMarkup {
	return InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{
		{Text: "✅ Onayla", CallbackData: "op:ok:" + id},
		{Text: "❌ Vazgeç", CallbackData: "op:no:" + id},
	}}}
}

func undoKeyboard(id string) InlineKeyboardMarkup {
	return InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{
		{Text: "↩️ Geri Al", CallbackData: "op:undo:" + id},
	}}}
}

// undoLabel renders a window duration as a short Turkish label, e.g. "15 saniye".
func undoLabel(d time.Duration) string {
	return strconv.Itoa(int(d.Seconds())) + " saniye"
}

// opSeqNum extracts the numeric sequence from an op ID ("op42" → 42) so /iptal
// can pick the newest op deterministically.
func opSeqNum(id string) uint64 {
	n, _ := strconv.ParseUint(strings.TrimPrefix(id, "op"), 10, 64)
	return n
}
