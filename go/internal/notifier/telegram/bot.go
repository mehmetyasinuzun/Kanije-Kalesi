package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kanije-kalesi/kanije/internal/access"
	"github.com/kanije-kalesi/kanije/internal/config"
	"github.com/kanije-kalesi/kanije/internal/event"
	"github.com/kanije-kalesi/kanije/internal/storage"
	"github.com/kanije-kalesi/kanije/internal/totp"
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
	acl    *access.Manager
	log    *slog.Logger

	// System services provided by the app layer
	lockScreen    func() error
	capturePhoto  func(ctx context.Context) ([]byte, error)
	captureScreen func(ctx context.Context) ([]byte, error)
	captureAudio  func(ctx context.Context, seconds int) ([]byte, error)
	getStatus     func() StatusInfo
	checkUpdate   func(ctx context.Context) string

	// Pending action state. Guarded by mu because Poll dispatches each update
	// in its own goroutine, so the menu, confirm, cancel, and expiry paths all
	// touch pendingAction concurrently.
	mu            sync.Mutex
	pendingAction *ActionState

	// Per-chat command rate limiter (anti-abuse).
	cmdLimiter *cmdRateLimiter

	// Fleet identity — used to route group commands to the right device.
	deviceLabel string
	osName      string
}

// BotConfig holds dependencies for the Bot.
type BotConfig struct {
	Config        *config.Config
	Client        *Client
	Wizard        *SetupWizard
	Store         storage.Storage
	Access        *access.Manager
	Log           *slog.Logger
	DeviceLabel   string
	OSName        string
	LockScreen    func() error
	CapturePhoto  func(ctx context.Context) ([]byte, error)
	CaptureScreen func(ctx context.Context) ([]byte, error)
	CaptureAudio  func(ctx context.Context, seconds int) ([]byte, error)
	GetStatus     func() StatusInfo
	CheckUpdate   func(ctx context.Context) string
}

// NewBot creates a fully wired Bot.
func NewBot(cfg BotConfig) *Bot {
	return &Bot{
		cfg:           cfg.Config,
		client:        cfg.Client,
		wizard:        cfg.Wizard,
		store:         cfg.Store,
		acl:           cfg.Access,
		log:           cfg.Log,
		deviceLabel:   cfg.DeviceLabel,
		osName:        cfg.OSName,
		lockScreen:    cfg.LockScreen,
		capturePhoto:  cfg.CapturePhoto,
		captureScreen: cfg.CaptureScreen,
		captureAudio:  cfg.CaptureAudio,
		getStatus:     cfg.GetStatus,
		checkUpdate:   cfg.CheckUpdate,
		cmdLimiter:    newCmdRateLimiter(cfg.Config.MaxCommandsPerMinute()),
	}
}

// SendEvent sends a formatted event notification to the configured chat.
// If there are media attachments (photo/screenshot), they are sent after the text.
func (b *Bot) SendEvent(ctx context.Context, ev event.Event) error {
	chatID := b.cfg.NotifyChatID() // fleet group when configured, else owner
	text := FormatEvent(ev)

	// In a shared fleet group, tag which device this event came from.
	if b.cfg.GroupID() != 0 && b.deviceLabel != "" {
		text = "🏷️ <b>" + safeHTML(b.deviceLabel) + "</b>\n" + text
	}

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
	_, err := b.client.SendMessage(ctx, b.cfg.NotifyChatID(), text, "HTML")
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
	b.registerCommands(ctx) // populate the "/" autocomplete menu

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
// commandCaps maps a command to the capability required to run it. Commands
// absent from this map (e.g. /yardim, /ping, /iptal) are open to any authorized
// tree member. Both Turkish and English aliases are listed.
var commandCaps = map[string]access.Capability{
	"/status": access.CapStatus, "/durum": access.CapStatus,
	"/olaylar": access.CapEvents, "/events": access.CapEvents,
	"/ozet": access.CapEvents, "/summary": access.CapEvents,
	"/foto": access.CapPhoto, "/photo": access.CapPhoto,
	"/ekran": access.CapScreen, "/screenshot": access.CapScreen,
	"/seskayit": access.CapAudio, "/record": access.CapAudio,
	"/kilitle": access.CapLock, "/lock": access.CapLock,
	"/yeniden": access.CapRestart, "/restart": access.CapRestart,
	"/kapat": access.CapShutdown, "/shutdown": access.CapShutdown,
	"/ayarlar": access.CapConfig, "/config": access.CapConfig,
	"/kurulum": access.CapConfig, "/setup": access.CapConfig,
	"/guncelle": access.CapUpdate, "/update": access.CapUpdate,
	"/dogrula": access.CapVerify, "/verify": access.CapVerify,
	"/ekle": access.CapInvite, "/davet": access.CapInvite,
	"/yonetim": access.CapInvite, "/kisiler": access.CapInvite,
	"/loglar": access.CapManage, "/audit": access.CapManage,
}

// groupBroadcastCmds run on every device in a shared group (no device target).
var groupBroadcastCmds = map[string]bool{
	"/cihazlar": true, "/devices": true,
	"/yardim": true, "/help": true,
}

// privateOnlyCmds manage configuration or the user tree and only make sense in a
// one-on-one chat, never in a shared fleet group.
var privateOnlyCmds = map[string]bool{
	"/kurulum": true, "/setup": true,
	"/ayarlar": true, "/config": true,
	"/yonetim": true, "/kisiler": true,
	"/ekle": true, "/davet": true,
	"/loglar": true, "/audit": true,
}

// routeGroupCommand decides whether this device should act on a group command.
// Broadcast commands always run. Otherwise the first argument must equal this
// device's label; if so it is consumed and the remaining text is returned.
func (b *Bot) routeGroupCommand(cmd, text string) (handled bool, stripped string) {
	if groupBroadcastCmds[cmd] {
		return true, text
	}
	fields := strings.Fields(text)
	if len(fields) < 2 {
		return false, text // no device named → ignore (avoid every device replying)
	}
	if !strings.EqualFold(fields[1], b.deviceLabel) {
		return false, text // addressed to another device
	}
	return true, strings.Join(append([]string{fields[0]}, fields[2:]...), " ")
}

func (b *Bot) handleMessage(ctx context.Context, m *Message) {
	if m.From == nil || m.Chat == nil {
		return
	}

	chatID := m.Chat.ID  // reply / notification target (group or private)
	actorID := m.From.ID // the human; in a group this differs from chatID
	isGroup := m.Chat.Type == "group" || m.Chat.Type == "supergroup"

	// Security: authorize the actual user — works in both private and group chats.
	if !b.isAuthorized(actorID) {
		b.log.Warn("yetkisiz mesaj alındı", "actor_id", actorID, "chat_id", chatID, "from", m.From.Username)
		return
	}

	// Per-chat command rate limit (anti-abuse). Keyed on the chat.
	if !b.cmdLimiter.allow(chatID) {
		b.log.Debug("komut hız sınırı aşıldı, mesaj atlanıyor", "chat_id", chatID)
		return
	}

	text := m.Text

	// In private chats the wizard may be waiting for typed input.
	if !isGroup && !isCommand(text) && b.wizard.HandleText(ctx, chatID, text) {
		return
	}

	cmd := extractCommand(text)

	// Fleet routing: in a shared group every device's bot sees the command.
	// Targeted commands must name this device (e.g. "/foto dizustu"); broadcast
	// commands (/cihazlar) run everywhere; an untargeted action command is
	// ignored so the devices don't all answer at once.
	if isGroup {
		handled, stripped := b.routeGroupCommand(cmd, text)
		if !handled {
			return
		}
		text = stripped
		cmd = extractCommand(text)
		if privateOnlyCmds[cmd] {
			b.reply(ctx, chatID, "🔒 Bu komutu bota <b>özelden</b> yaz (grupta değil).")
			return
		}
	}

	// Capability gate (server-side, authorized by the human's user ID). Ungated
	// commands (help/ping/cancel/cihazlar) are open to any tree member.
	if reqCap, gated := commandCaps[cmd]; gated && b.acl != nil {
		if !b.acl.Can(actorID, reqCap) {
			b.reply(ctx, chatID, "⛔ <b>Yetki yok.</b> Bu komut için gerekli izne sahip değilsin.")
			return
		}
		b.acl.Log(ctx, actorID, strings.TrimPrefix(cmd, "/"), "")
	}

	switch cmd {
	case "/start", "/yardim", "/help", "/komutlar", "/komut", "/menu":
		b.cmdHelp(ctx, chatID)
	case "/status", "/durum":
		b.cmdStatus(ctx, chatID)
	case "/foto", "/photo":
		b.cmdFoto(ctx, chatID)
	case "/ekran", "/screenshot":
		b.cmdEkran(ctx, chatID)
	case "/seskayit", "/record":
		b.cmdSesKayit(ctx, chatID, text)
	case "/olaylar", "/events":
		b.cmdOlaylar(ctx, chatID, text)
	case "/ozet", "/summary":
		b.cmdOzet(ctx, chatID)
	case "/ayarlar", "/config":
		b.cmdAyarlar(ctx, chatID)
	case "/kurulum", "/setup":
		b.wizard.SendMainMenu(ctx, chatID)
	case "/kilitle", "/lock":
		b.cmdKilitle(ctx, chatID)
	case "/ping":
		b.cmdPing(ctx, chatID)
	case "/yeniden", "/restart":
		b.cmdYeniden(ctx, chatID, text)
	case "/kapat", "/shutdown":
		b.cmdKapat(ctx, chatID, text)
	case "/iptal", "/cancel":
		b.cmdIptal(ctx, chatID)
	case "/guncelle", "/update":
		b.cmdGuncelle(ctx, chatID)
	case "/dogrula", "/verify":
		b.cmdDogrula(ctx, chatID)
	case "/yonetim", "/kisiler":
		b.cmdYonetim(ctx, chatID)
	case "/ekle", "/davet":
		b.cmdEkle(ctx, chatID, text)
	case "/loglar", "/audit":
		b.cmdLoglar(ctx, chatID)
	case "/cihazlar", "/devices":
		b.cmdCihazlar(ctx, chatID)
	default:
		if text != "" && isCommand(text) {
			b.reply(ctx, chatID, "❓ <b>Bilinmeyen komut:</b> <code>"+safeHTML(cmd)+"</code>\n\n"+
				"Tüm komutlar için /yardim yaz. İpucu: giriş kutusuna <b>/</b> yazınca menü açılır.")
		}
	}
}

// handleCallback processes inline keyboard button presses.
func (b *Bot) handleCallback(ctx context.Context, cq *CallbackQuery) {
	if cq.From == nil || cq.Message == nil {
		return
	}

	chatID := cq.Message.Chat.ID // reply target (group or private)
	actorID := cq.From.ID        // the human; authorize by this in groups too
	if !b.isAuthorized(actorID) {
		return
	}

	// Event drill-down: olay:<id> — show the full stored detail of one event.
	if strings.HasPrefix(cq.Data, "olay:") {
		if b.acl != nil && !b.acl.Can(actorID, access.CapEvents) {
			b.client.AnswerCallbackQuery(ctx, cq.ID, "Yetki yok")
			return
		}
		b.client.AnswerCallbackQuery(ctx, cq.ID, "")
		b.showEventDetail(ctx, chatID, strings.TrimPrefix(cq.Data, "olay:"))
		return
	}

	// Delegate to the multi-user management UI for mng: callbacks
	if strings.HasPrefix(cq.Data, "mng:") {
		b.handleManageCallback(ctx, chatID, cq.Message.MessageID, cq.ID, cq.Data)
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

// registerCommands publishes the command list to Telegram so users get a "/"
// autocomplete menu and a Menu button. Best-effort; failure is non-fatal.
func (b *Bot) registerCommands(ctx context.Context) {
	cmds := []BotCommand{
		{"yardim", "📜 Komut listesi"},
		{"status", "📊 Sistem durumu (CPU/RAM/disk)"},
		{"olaylar", "📋 Son olaylar (tip/sayı ile filtre)"},
		{"ozet", "📈 Son 7 günün olay özeti"},
		{"foto", "📷 Kameradan anlık fotoğraf"},
		{"ekran", "🖥️ Ekran görüntüsü"},
		{"seskayit", "🎤 Mikrofon kaydı (saniye)"},
		{"cihazlar", "🛰️ Tüm cihazları listele"},
		{"terminal", "💻 Uzak komut çalıştır"},
		{"kilitle", "🔒 Ekranı kilitle"},
		{"dogrula", "🛡️ Olay günlüğü bütünlüğü"},
		{"guncelle", "⬆️ Güncelleme kontrol/kur"},
		{"yeniden", "🔄 Yeniden başlat (onaylı)"},
		{"kapat", "⏻ Kapat (onaylı)"},
		{"ekle", "➕ Kişi ekle"},
		{"yonetim", "👥 Kişi yönetimi"},
		{"loglar", "🧾 İşlem günlüğü"},
		{"kurulum", "⚙️ Ayar menüsü"},
	}
	if err := b.client.SetMyCommands(ctx, cmds); err != nil {
		b.log.Debug("setMyCommands başarısız (önemsiz)", "err", err)
	}
}

func (b *Bot) cmdStatus(ctx context.Context, chatID int64) {
	info := b.getStatus()
	b.reply(ctx, chatID, FormatStatus(info))
}

// cmdCihazlar replies with THIS device's card. In a fleet group every device's
// bot answers, so the user sees the whole fleet at once.
func (b *Bot) cmdCihazlar(ctx context.Context, chatID int64) {
	var info StatusInfo
	if b.getStatus != nil {
		info = b.getStatus()
	}
	b.reply(ctx, chatID, FormatDeviceCard(b.deviceLabel, b.osName, info))
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

// cmdSesKayit records microphone audio. Optional arg: duration in seconds
// (default 30, capped at 600). Example: "/seskayit 60".
func (b *Bot) cmdSesKayit(ctx context.Context, chatID int64, text string) {
	if b.captureAudio == nil {
		b.reply(ctx, chatID, "❌ Ses kaydı desteği bu derlemede yok.")
		return
	}

	seconds := 0 // 0 → recorder uses its default (30s)
	if arg := commandArg(text); arg != "" {
		n, err := strconv.Atoi(arg)
		if err != nil || n <= 0 {
			b.reply(ctx, chatID, "❌ Geçersiz süre. Örnek: <code>/seskayit 30</code> (saniye).")
			return
		}
		seconds = n
	}

	shown := seconds
	if shown <= 0 {
		shown = 30
	}
	if shown > 600 {
		shown = 600
	}

	b.reply(ctx, chatID, fmt.Sprintf("🎤 %d saniyelik ses kaydı alınıyor…", shown))

	data, err := b.captureAudio(ctx, seconds)
	if err != nil {
		b.reply(ctx, chatID, "❌ Ses kaydı hatası: "+safeHTML(err.Error()))
		return
	}

	caption := fmt.Sprintf("🎤 Ses kaydı (%d sn)", shown)
	if err := b.client.SendAudio(ctx, chatID, data, "kayit.mp3", caption); err != nil {
		b.reply(ctx, chatID, "❌ Ses gönderilemedi: "+safeHTML(err.Error()))
	}
}

// cmdOlaylar lists recent events. Optional arg: a number (count, capped at 50)
// or an event type to filter by (e.g. "/olaylar login_failed").
func (b *Bot) cmdOlaylar(ctx context.Context, chatID int64, text string) {
	arg := commandArg(text)
	limit := b.cfg.MaxRecentEvents()

	var events []event.Event
	var err error
	switch {
	case arg == "":
		events, err = b.store.RecentEvents(ctx, limit)
	case isNumeric(arg):
		n, _ := strconv.Atoi(arg)
		if n < 1 {
			n = 1
		}
		if n > 50 {
			n = 50
		}
		events, err = b.store.RecentEvents(ctx, n)
	default:
		events, err = b.store.QueryEvents(ctx, storage.EventFilter{Type: event.Type(arg), Limit: limit})
	}

	if err != nil {
		b.reply(ctx, chatID, "❌ Olaylar alınamadı: "+safeHTML(err.Error()))
		return
	}
	if len(events) == 0 {
		b.reply(ctx, chatID, FormatRecentEvents(events))
		return
	}
	// Each event gets a button → tap for the full CTI detail (olay:<id>).
	b.client.SendMessageWithKeyboard(ctx, chatID, FormatRecentEvents(events), eventsKeyboard(events))
}

// eventsKeyboard builds a numbered button per event for the drill-down.
func eventsKeyboard(events []event.Event) InlineKeyboardMarkup {
	var rows [][]InlineKeyboardButton
	var row []InlineKeyboardButton
	for i, ev := range events {
		row = append(row, InlineKeyboardButton{
			Text:         fmt.Sprintf("%d %s", i+1, ev.Type.Emoji()),
			CallbackData: "olay:" + strconv.FormatInt(ev.ID, 10),
		})
		if len(row) == 5 {
			rows = append(rows, row)
			row = nil
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}
	return InlineKeyboardMarkup{InlineKeyboard: rows}
}

// showEventDetail fetches one event by ID and sends its full detail.
func (b *Bot) showEventDetail(ctx context.Context, chatID int64, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return
	}
	ev, ok, err := b.store.EventByID(ctx, id)
	if err != nil {
		b.reply(ctx, chatID, "❌ Olay alınamadı: "+safeHTML(err.Error()))
		return
	}
	if !ok {
		b.reply(ctx, chatID, "❌ Olay bulunamadı.")
		return
	}
	b.reply(ctx, chatID, FormatEventDetail(ev))
}

// cmdOzet sends a per-type summary of events from the last 7 days.
func (b *Bot) cmdOzet(ctx context.Context, chatID int64) {
	const days = 7
	since := time.Now().AddDate(0, 0, -days)
	counts, total, err := b.store.EventStats(ctx, since)
	if err != nil {
		b.reply(ctx, chatID, "❌ Özet alınamadı: "+safeHTML(err.Error()))
		return
	}
	b.reply(ctx, chatID, FormatSummary(counts, total, days))
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
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

func (b *Bot) cmdGuncelle(ctx context.Context, chatID int64) {
	if b.checkUpdate == nil {
		b.reply(ctx, chatID, "❌ Güncelleme bu derlemede desteklenmiyor.")
		return
	}
	b.reply(ctx, chatID, "🔍 Güncelleme kontrol ediliyor…")
	b.reply(ctx, chatID, b.checkUpdate(ctx))
}

func (b *Bot) cmdDogrula(ctx context.Context, chatID int64) {
	b.reply(ctx, chatID, "🔍 Olay günlüğü bütünlüğü kontrol ediliyor…")
	ok, brokenAt, total, err := b.store.VerifyChain(ctx)
	if err != nil {
		b.reply(ctx, chatID, "❌ Doğrulama hatası: "+safeHTML(err.Error()))
		return
	}
	if ok {
		b.reply(ctx, chatID, fmt.Sprintf(
			"✅ <b>Bütünlük doğrulandı</b>\n%d olay · hash zinciri sağlam — kayıtlar kurcalanmamış.", total))
		return
	}
	b.reply(ctx, chatID, fmt.Sprintf(
		"⚠️ <b>KURCALAMA TESPİT EDİLDİ</b>\nOlay #%d'de hash zinciri bozuk — kayıtlar değiştirilmiş veya silinmiş olabilir. (Taranan: %d)",
		brokenAt, total))
}

func (b *Bot) cmdPing(ctx context.Context, chatID int64) {
	info := b.getStatus()
	b.reply(ctx, chatID, "🏓 Pong! Çalışma süresi: <b>"+formatDuration(info.Uptime)+"</b>")
}

// cmdYeniden sends a confirmation keyboard before actually rebooting.
func (b *Bot) cmdYeniden(ctx context.Context, chatID int64, text string) {
	if !b.checkTOTP(ctx, chatID, text, "/yeniden") {
		return
	}
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
func (b *Bot) cmdKapat(ctx context.Context, chatID int64, text string) {
	if !b.checkTOTP(ctx, chatID, text, "/kapat") {
		return
	}
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

// checkTOTP enforces two-factor confirmation for dangerous commands when a TOTP
// secret is configured. Returns true if the command may proceed.
func (b *Bot) checkTOTP(ctx context.Context, chatID int64, text, cmd string) bool {
	if !b.cfg.TOTPEnabled() {
		return true
	}
	if totp.Validate(b.cfg.TOTPSecret(), commandArg(text)) {
		return true
	}
	b.reply(ctx, chatID,
		"🔐 Bu işlem iki adımlı doğrulama gerektiriyor.\n"+
			"Komutu doğrulama kodunuzla gönderin: <code>"+cmd+" 123456</code>")
	return false
}

// commandArg returns the first argument after a command, e.g. "/kapat 123456" → "123456".
func commandArg(text string) string {
	if fields := strings.Fields(text); len(fields) >= 2 {
		return fields[1]
	}
	return ""
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
	// The delegation tree is the source of truth: a chat is authorized iff it is
	// a member. The owner and any migrated allowlist entries are seeded at boot.
	if b.acl != nil {
		return b.acl.Known(chatID)
	}
	// Fallback only if the manager was never wired (defensive).
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
	tok := text
	for i, r := range text {
		if r == ' ' || r == '\n' {
			tok = text[:i]
			break
		}
	}
	// In groups Telegram appends "@botusername" to commands — strip it.
	if at := strings.IndexByte(tok, '@'); at >= 0 {
		tok = tok[:at]
	}
	return tok
}
