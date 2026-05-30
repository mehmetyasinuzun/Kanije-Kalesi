package telegram

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/kanije-kalesi/kanije/internal/config"
	"github.com/kanije-kalesi/kanije/internal/event"
	"github.com/kanije-kalesi/kanije/internal/storage"
)

// pendingActionTTL is how long a dangerous action (restart/shutdown) waits for
// confirmation before auto-expiring.
const pendingActionTTL = 15 * time.Second

// ActionState tracks a pending dangerous action (shutdown/restart).
type ActionState struct {
	kind    string // "kapat" | "yeniden"
	expires time.Time
	cancel  context.CancelFunc
}

// Bot handles all Telegram interactions: event notifications, commands,
// callback queries, and the setup wizard.
type Bot struct {
	cfg    *config.Config
	client *Client
	wizard *SetupWizard
	store  storage.Storage
	log    *slog.Logger

	// System services provided by the app layer
	lockScreen    func() error
	capturePhoto  func(ctx context.Context) ([]byte, error)
	captureScreen func(ctx context.Context) ([]byte, error)
	getStatus     func() StatusInfo

	// Pending action state. Guarded by mu because Poll dispatches each update
	// in its own goroutine, so the menu, confirm, cancel, and expiry paths all
	// touch pendingAction concurrently.
	mu            sync.Mutex
	pendingAction *ActionState

	// Per-chat command rate limiter (anti-abuse).
	cmdLimiter *cmdRateLimiter
}

// BotConfig holds dependencies for the Bot.
type BotConfig struct {
	Config        *config.Config
	Client        *Client
	Wizard        *SetupWizard
	Store         storage.Storage
	Log           *slog.Logger
	LockScreen    func() error
	CapturePhoto  func(ctx context.Context) ([]byte, error)
	CaptureScreen func(ctx context.Context) ([]byte, error)
	GetStatus     func() StatusInfo
}

// NewBot creates a fully wired Bot.
func NewBot(cfg BotConfig) *Bot {
	return &Bot{
		cfg:           cfg.Config,
		client:        cfg.Client,
		wizard:        cfg.Wizard,
		store:         cfg.Store,
		log:           cfg.Log,
		lockScreen:    cfg.LockScreen,
		capturePhoto:  cfg.CapturePhoto,
		captureScreen: cfg.CaptureScreen,
		getStatus:     cfg.GetStatus,
		cmdLimiter:    newCmdRateLimiter(cfg.Config.MaxCommandsPerMinute()),
	}
}

// SendEvent sends a formatted event notification to the configured chat.
// If there are media attachments (photo/screenshot), they are sent after the text.
func (b *Bot) SendEvent(ctx context.Context, ev event.Event) error {
	chatID := b.cfg.ChatID()
	text := FormatEvent(ev)

	if _, err := b.client.SendMessage(ctx, chatID, text, "HTML"); err != nil {
		return err
	}

	for _, att := range ev.Attachments {
		switch att.Type {
		case event.AttachmentPhoto:
			if err := b.client.SendPhoto(ctx, chatID, att.Data, att.Caption); err != nil {
				b.log.Warn("fotoğraf gönderilemedi", "err", err)
			}
		case event.AttachmentScreenshot:
			if err := b.client.SendPhoto(ctx, chatID, att.Data, att.Caption); err != nil {
				b.log.Warn("ekran görüntüsü gönderilemedi", "err", err)
			}
		}
	}

	return nil
}

// SendMessage sends a plain text message to the configured chat.
func (b *Bot) SendMessage(ctx context.Context, text string) error {
	_, err := b.client.SendMessage(ctx, b.cfg.ChatID(), text, "HTML")
	return err
}

// TestConnection verifies the bot token is valid.
func (b *Bot) TestConnection(ctx context.Context) error {
	_, err := b.client.GetMe(ctx)
	return err
}

// Poll runs the long-polling update loop until ctx is canceled.
// Each incoming update is dispatched in a short-lived goroutine.
func (b *Bot) Poll(ctx context.Context) error {
	var offset int64
	backoff := 2 * time.Second
	maxBackoff := 60 * time.Second

	b.log.Info("Telegram bot polling başlatıldı")

	for {
		if ctx.Err() != nil {
			return nil
		}

		updates, err := b.client.GetUpdates(ctx, offset, 30)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			b.log.Debug("getUpdates hatası (geçici)", "err", err, "bekleme", backoff)
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(backoff):
			}
			backoff = min(backoff*2, maxBackoff)
			continue
		}
		backoff = 2 * time.Second // Reset on success

		for _, u := range updates {
			offset = u.UpdateID + 1
			go b.handleUpdate(ctx, u)
		}
	}
}

// handleUpdate routes a single update to the appropriate handler.
func (b *Bot) handleUpdate(ctx context.Context, u Update) {
	switch {
	case u.Message != nil:
		b.handleMessage(ctx, u.Message)
	case u.CallbackQuery != nil:
		b.handleCallback(ctx, u.CallbackQuery)
	}
}

// handleMessage processes incoming text messages.
func (b *Bot) handleMessage(ctx context.Context, m *Message) {
	if m.From == nil || m.Chat == nil {
		return
	}

	// Security: only accept messages from authorized chat IDs
	if !b.isAuthorized(m.Chat.ID) {
		b.log.Warn("yetkisiz mesaj alındı",
			"chat_id", m.Chat.ID,
			"from", m.From.Username)
		return
	}

	// Per-chat command rate limit (anti-abuse — bir saldırgan /foto spam'i ile
	// kamerayı kilitleyemesin). Aşımda mesaj sessizce düşürülür.
	if !b.cmdLimiter.allow(m.Chat.ID) {
		b.log.Debug("komut hız sınırı aşıldı, mesaj atlanıyor", "chat_id", m.Chat.ID)
		return
	}

	text := m.Text
	chatID := m.Chat.ID

	// Let the setup wizard intercept non-command messages when waiting for input
	if !isCommand(text) && b.wizard.HandleText(ctx, chatID, text) {
		return
	}

	// Route commands
	cmd := extractCommand(text)
	switch cmd {
	case "/start", "/yardim", "/help":
		b.cmdHelp(ctx, chatID)
	case "/status", "/durum":
		b.cmdStatus(ctx, chatID)
	case "/foto", "/photo":
		b.cmdFoto(ctx, chatID)
	case "/ekran", "/screenshot":
		b.cmdEkran(ctx, chatID)
	case "/olaylar", "/events":
		b.cmdOlaylar(ctx, chatID)
	case "/ayarlar", "/config":
		b.cmdAyarlar(ctx, chatID)
	case "/kurulum", "/setup":
		b.wizard.SendMainMenu(ctx, chatID)
	case "/kilitle", "/lock":
		b.cmdKilitle(ctx, chatID)
	case "/ping":
		b.cmdPing(ctx, chatID)
	case "/yeniden", "/restart":
		b.cmdYeniden(ctx, chatID)
	case "/kapat", "/shutdown":
		b.cmdKapat(ctx, chatID)
	case "/iptal", "/cancel":
		b.cmdIptal(ctx, chatID)
	default:
		if text != "" && isCommand(text) {
			b.reply(ctx, chatID, "❓ Bilinmeyen komut. /yardim yazın.")
		}
	}
}

// handleCallback processes inline keyboard button presses.
func (b *Bot) handleCallback(ctx context.Context, cq *CallbackQuery) {
	if cq.From == nil || cq.Message == nil {
		return
	}

	chatID := cq.Message.Chat.ID
	if !b.isAuthorized(chatID) {
		return
	}

	// Delegate to the setup wizard for wizard: callbacks
	if len(cq.Data) >= 7 && cq.Data[:7] == "wizard:" {
		b.wizard.HandleCallback(ctx, chatID, cq.Message.MessageID, cq.ID, cq.Data)
		return
	}

	// Handle confirmation callbacks for dangerous operations
	switch cq.Data {
	case "confirm:yeniden":
		b.executeRestart(ctx, chatID)
		b.client.AnswerCallbackQuery(ctx, cq.ID, "Yeniden başlatılıyor...")
	case "confirm:kapat":
		b.executeShutdown(ctx, chatID)
		b.client.AnswerCallbackQuery(ctx, cq.ID, "Kapatılıyor...")
	case "confirm:iptal":
		b.cancelPendingAction(ctx, chatID)
		b.client.AnswerCallbackQuery(ctx, cq.ID, "İptal edildi")
	default:
		b.client.AnswerCallbackQuery(ctx, cq.ID, "")
	}
}

// ---- Command handlers ----

func (b *Bot) cmdHelp(ctx context.Context, chatID int64) {
	b.reply(ctx, chatID, FormatHelp())
}

func (b *Bot) cmdStatus(ctx context.Context, chatID int64) {
	info := b.getStatus()
	b.reply(ctx, chatID, FormatStatus(info))
}

func (b *Bot) cmdFoto(ctx context.Context, chatID int64) {
	if b.capturePhoto == nil {
		b.reply(ctx, chatID, "❌ Kamera desteği bu derlemede yok.")
		return
	}

	b.reply(ctx, chatID, "📷 Fotoğraf çekiliyor…")

	data, err := b.capturePhoto(ctx)
	if err != nil {
		b.reply(ctx, chatID, "❌ Kamera hatası: "+safeHTML(err.Error()))
		return
	}

	if err := b.client.SendPhoto(ctx, chatID, data, "📷 Anlık kamera görüntüsü"); err != nil {
		b.reply(ctx, chatID, "❌ Fotoğraf gönderilemedi: "+safeHTML(err.Error()))
	}
}

func (b *Bot) cmdEkran(ctx context.Context, chatID int64) {
	if b.captureScreen == nil {
		b.reply(ctx, chatID, "❌ Ekran görüntüsü desteği bu derlemede yok.")
		return
	}

	b.reply(ctx, chatID, "🖥️ Ekran görüntüsü alınıyor…")

	data, err := b.captureScreen(ctx)
	if err != nil {
		b.reply(ctx, chatID, "❌ Ekran görüntüsü hatası: "+safeHTML(err.Error()))
		return
	}

	if err := b.client.SendPhoto(ctx, chatID, data, "🖥️ Anlık ekran görüntüsü"); err != nil {
		b.reply(ctx, chatID, "❌ Ekran görüntüsü gönderilemedi: "+safeHTML(err.Error()))
	}
}

func (b *Bot) cmdOlaylar(ctx context.Context, chatID int64) {
	events, err := b.store.RecentEvents(ctx, b.cfg.MaxRecentEvents())
	if err != nil {
		b.reply(ctx, chatID, "❌ Olaylar alınamadı: "+safeHTML(err.Error()))
		return
	}
	b.reply(ctx, chatID, FormatRecentEvents(events))
}

func (b *Bot) cmdAyarlar(ctx context.Context, chatID int64) {
	json := b.cfg.GetSafeJSON()
	b.reply(ctx, chatID, "⚙️ <b>Mevcut Yapılandırma</b>\n\n<pre>"+safeHTML(json)+"</pre>")
}

func (b *Bot) cmdKilitle(ctx context.Context, chatID int64) {
	if b.lockScreen == nil {
		b.reply(ctx, chatID, "❌ Ekran kilitleme bu platformda desteklenmiyor.")
		return
	}
	if err := b.lockScreen(); err != nil {
		b.reply(ctx, chatID, "❌ Ekran kilitlenemedi: "+safeHTML(err.Error()))
		return
	}
	b.reply(ctx, chatID, "🔒 Ekran kilitlendi.")
}

func (b *Bot) cmdPing(ctx context.Context, chatID int64) {
	info := b.getStatus()
	b.reply(ctx, chatID, "🏓 Pong! Çalışma süresi: <b>"+formatDuration(info.Uptime)+"</b>")
}

// cmdYeniden sends a confirmation keyboard before actually rebooting.
func (b *Bot) cmdYeniden(ctx context.Context, chatID int64) {
	kb := InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{Text: "✅ Evet, yeniden başlat", CallbackData: "confirm:yeniden"},
				{Text: "❌ İptal", CallbackData: "confirm:iptal"},
			},
		},
	}
	b.client.SendMessageWithKeyboard(ctx, chatID,
		"⚠️ <b>Sistemi yeniden başlatmak istediğinizden emin misiniz?</b>\n\n"+
			"Bu işlem geri alınamaz. 15 saniye içinde onaylamazsanız iptal olur.",
		kb)

	b.armPendingAction("yeniden")
}

// cmdKapat sends a confirmation keyboard before shutting down.
func (b *Bot) cmdKapat(ctx context.Context, chatID int64) {
	kb := InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{Text: "✅ Evet, kapat", CallbackData: "confirm:kapat"},
				{Text: "❌ İptal", CallbackData: "confirm:iptal"},
			},
		},
	}
	b.client.SendMessageWithKeyboard(ctx, chatID,
		"⚠️ <b>Sistemi kapatmak istediğinizden emin misiniz?</b>\n\n"+
			"15 saniye içinde onaylamazsanız iptal olur.",
		kb)

	b.armPendingAction("kapat")
}

func (b *Bot) cmdIptal(ctx context.Context, chatID int64) {
	if b.cancelPendingAction(ctx, chatID) {
		return
	}
	// Also cancel wizard input state
	if b.wizard.CancelInput(chatID) {
		b.reply(ctx, chatID, "✅ Bekleyen giriş iptal edildi.")
		return
	}
	b.reply(ctx, chatID, "ℹ️ İptal edilecek bekleyen işlem yok.")
}

func (b *Bot) cancelPendingAction(ctx context.Context, chatID int64) bool {
	state, ok := b.clearPendingAction()
	if !ok {
		return false
	}
	state.cancel()
	b.reply(ctx, chatID, "✅ İşlem iptal edildi.")
	return true
}

func (b *Bot) executeRestart(ctx context.Context, chatID int64) {
	state, ok := b.takePendingAction("yeniden")
	if !ok {
		b.reply(ctx, chatID, "⚠️ Onaylanacak işlem bulunamadı. /yeniden yazın.")
		return
	}
	state.cancel()
	b.reply(ctx, chatID, "🔄 Sistem yeniden başlatılıyor… İyi günler!")
	systemRestart()
}

func (b *Bot) executeShutdown(ctx context.Context, chatID int64) {
	state, ok := b.takePendingAction("kapat")
	if !ok {
		b.reply(ctx, chatID, "⚠️ Onaylanacak işlem bulunamadı. /kapat yazın.")
		return
	}
	state.cancel()
	b.reply(ctx, chatID, "🔴 Sistem kapatılıyor… Güle güle!")
	systemShutdown()
}

// ---- Pending-action state (mutex-guarded) ----

// armPendingAction registers a dangerous action and schedules its auto-expiry.
// Any previously pending action is canceled.
func (b *Bot) armPendingAction(kind string) {
	actionCtx, cancelAction := context.WithTimeout(context.Background(), pendingActionTTL)
	state := &ActionState{
		kind:    kind,
		expires: time.Now().Add(pendingActionTTL),
		cancel:  cancelAction,
	}
	b.setPendingAction(state)

	// Expire by identity: only clear the slot if it still holds *this* action,
	// so a newer action started within the TTL is not wrongly cleared.
	go func() {
		<-actionCtx.Done()
		b.mu.Lock()
		if b.pendingAction == state {
			b.pendingAction = nil
		}
		b.mu.Unlock()
	}()
}

func (b *Bot) setPendingAction(s *ActionState) {
	b.mu.Lock()
	old := b.pendingAction
	b.pendingAction = s
	b.mu.Unlock()
	if old != nil {
		old.cancel() // release the previous action's expiry goroutine
	}
}

// takePendingAction atomically consumes the pending action if it matches kind.
func (b *Bot) takePendingAction(kind string) (*ActionState, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.pendingAction == nil || b.pendingAction.kind != kind {
		return nil, false
	}
	s := b.pendingAction
	b.pendingAction = nil
	return s, true
}

// clearPendingAction atomically consumes any pending action.
func (b *Bot) clearPendingAction() (*ActionState, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.pendingAction == nil {
		return nil, false
	}
	s := b.pendingAction
	b.pendingAction = nil
	return s, true
}

// ---- Authorization ----

func (b *Bot) isAuthorized(chatID int64) bool {
	if chatID == b.cfg.ChatID() {
		return true
	}
	for _, id := range b.cfg.AllowedChatIDs() {
		if chatID == id {
			return true
		}
	}
	return false
}

// ---- Helpers ----

func (b *Bot) reply(ctx context.Context, chatID int64, text string) {
	if _, err := b.client.SendMessage(ctx, chatID, text, "HTML"); err != nil {
		b.log.Warn("mesaj gönderilemedi", "err", err, "chat_id", chatID)
	}
}

func isCommand(text string) bool {
	return len(text) > 1 && text[0] == '/'
}

func extractCommand(text string) string {
	if !isCommand(text) {
		return ""
	}
	for i, r := range text {
		if r == ' ' || r == '\n' {
			return text[:i]
		}
	}
	return text
}
