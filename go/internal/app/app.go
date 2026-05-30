// Package app is the application orchestrator. It wires all modules together,
// manages their lifecycle, and handles graceful shutdown.
package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/kanije-kalesi/kanije/internal/access"
	"github.com/kanije-kalesi/kanije/internal/capture"
	"github.com/kanije-kalesi/kanije/internal/config"
	"github.com/kanije-kalesi/kanije/internal/event"
	"github.com/kanije-kalesi/kanije/internal/geoip"
	"github.com/kanije-kalesi/kanije/internal/listener"
	"github.com/kanije-kalesi/kanije/internal/network"
	"github.com/kanije-kalesi/kanije/internal/notifier/telegram"
	"github.com/kanije-kalesi/kanije/internal/notifier/webhook"
	"github.com/kanije-kalesi/kanije/internal/shell"
	"github.com/kanije-kalesi/kanije/internal/storage"
	"github.com/kanije-kalesi/kanije/internal/sysinfo"
	"github.com/kanije-kalesi/kanije/internal/updater"
	"golang.org/x/sync/errgroup"
)

const (
	repoOwner = "mehmetyasinuzun"
	repoName  = "Kanije-Kalesi"
)

// App is the top-level application object. Create with New(), run with Run().
type App struct {
	cfg      *config.Config
	log      *slog.Logger
	store    storage.Storage
	bus      *event.Bus
	bot      *telegram.Bot
	manager  *listener.Manager
	netMon   *network.Monitor
	camera   *capture.Camera
	screen   *capture.Screenshotter
	updater  *updater.Updater
	geo      *geoip.Resolver
	webhooks *webhook.Sender

	version  string
	osName   string             // cached friendly OS name (e.g. "Windows 11")
	cancel   context.CancelFunc // triggers graceful shutdown (used by self-update)
	updating atomic.Bool        // true while a self-update restart is in progress

	// Metrics
	startedAt time.Time
	lastEvent atomic.Pointer[event.Event]
}

// New creates and wires the application. Call Run() to start.
// version is the build version (e.g. "v1.1.2"), used for self-update checks.
func New(cfg *config.Config, log *slog.Logger, version string) (*App, error) {
	// Note: the single-instance guard is acquired and released in main.cmdStart.

	// Storage
	store, err := storage.NewSQLite(cfg.Storage.DBPath)
	if err != nil {
		return nil, fmt.Errorf("veritabanı açılamadı: %w", err)
	}

	// Access control — capability-based delegation tree for multi-user control.
	// Seeds the owner as root (all caps) and migrates any prior allowlist on first run.
	acl := access.NewManager(store, log.With("module", "access"))
	if err := acl.Bootstrap(context.Background(), cfg.ChatID(), cfg.AllowedChatIDs()); err != nil {
		store.Close()
		return nil, fmt.Errorf("erişim yöneticisi başlatılamadı: %w", err)
	}

	// Event bus
	bus := event.NewBus(event.BusConfig{
		BufferSize:   512,
		MaxPerMinute: cfg.Security.MaxEventsPerMinute,
		DedupWindow:  time.Duration(cfg.Security.DedupWindowSec) * time.Second,
	})

	// Telegram client
	tgClient := telegram.NewClient(
		cfg.BotToken(),
		cfg.Telegram.SendTimeoutSec,
		log.With("module", "telegram"),
	)

	// Setup wizard
	wizard := telegram.NewSetupWizard(cfg, tgClient, log.With("module", "wizard"))

	// Camera
	cam := capture.NewCamera(capture.CameraConfig{
		FFmpegPath:   cfg.Camera.FFmpegPath,
		DeviceIndex:  cfg.Camera.DeviceIndex,
		DeviceName:   cfg.Camera.DeviceName,
		Width:        cfg.Camera.Width,
		Height:       cfg.Camera.Height,
		WarmupFrames: cfg.Camera.WarmupFrames,
		JPEGQuality:  cfg.Camera.JPEGQuality,
	}, log.With("module", "camera"))

	// Audio recorder (microphone via ffmpeg) — on-demand /seskayit
	audioRec := capture.NewAudioRecorder(capture.AudioConfig{
		FFmpegPath: cfg.Audio.FFmpegPath,
		DeviceName: cfg.Audio.DeviceName,
		Bitrate:    cfg.Audio.Bitrate,
	}, log.With("module", "audio"))

	// Screenshot
	screen := capture.NewScreenshotter(cfg.Screenshot.JPEGQuality, log.With("module", "screenshot"))

	// Webhook fan-out targets (Discord / generic JSON).
	var hooks []webhook.Target
	for _, w := range cfg.WebhookTargets() {
		hooks = append(hooks, webhook.Target{Name: w.Name, URL: w.URL, Format: w.Format})
	}

	app := &App{
		cfg:      cfg,
		log:      log,
		store:    store,
		bus:      bus,
		camera:   cam,
		screen:   screen,
		version:  version,
		osName:   sysinfo.OSName(),
		updater:  updater.New(repoOwner, repoName, version, log.With("module", "updater")),
		geo:      geoip.New(),
		webhooks: webhook.New(hooks, log.With("module", "webhook")),
	}

	// Fleet identity — friendly device label (falls back to hostname).
	deviceLabel := cfg.DeviceLabel()
	if deviceLabel == "" {
		if h, err := os.Hostname(); err == nil && h != "" {
			deviceLabel = h
		} else {
			deviceLabel = "cihaz"
		}
	}

	// Bot
	bot := telegram.NewBot(telegram.BotConfig{
		Config:        cfg,
		Client:        tgClient,
		Wizard:        wizard,
		Store:         store,
		Access:        acl,
		Shell:         shell.New(),
		Log:           log.With("module", "bot"),
		DeviceLabel:   deviceLabel,
		OSName:        app.osName,
		LockScreen:    lockScreen,
		CapturePhoto:  cam.Capture,
		CaptureScreen: screen.Capture,
		CaptureAudio:  audioRec.Record,
		GetStatus:     app.collectStatus,
		CheckUpdate:   func(ctx context.Context) string { return app.checkForUpdate(ctx, true) },
	})
	app.bot = bot

	// Listener manager (platform-specific listeners injected by main.go)
	app.manager = listener.NewManager(log.With("module", "listeners"))

	// Network monitor
	app.netMon = network.NewMonitor(network.MonitorConfig{
		CheckIntervalSec: cfg.Network.CheckIntervalSec,
		CheckHost:        cfg.Network.CheckHost,
		CheckPort:        cfg.Network.CheckPort,
	}, log.With("module", "network"))

	return app, nil
}

// SetListeners injects platform-specific listeners.
// Must be called before Run().
func (a *App) SetListeners(listeners ...listener.Listener) {
	a.manager = listener.NewManager(a.log.With("module", "listeners"), listeners...)
}

// Run starts all subsystems and blocks until a shutdown signal is received.
// Returns nil on clean shutdown, error on fatal startup failure.
func (a *App) Run() error {
	a.startedAt = time.Now()

	// Pre-flight checks
	if !a.cfg.IsConfigured() {
		a.log.Warn("Telegram yapılandırması eksik — bot komutları çalışmayacak",
			"öneri", "KANIJE_BOT_TOKEN ve KANIJE_CHAT_ID ortam değişkenlerini ayarlayın")
	} else {
		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := a.bot.TestConnection(ctx)
		cancel()
		if err != nil {
			a.log.Warn("Telegram bağlantısı kurulamadı", "err", err)
		} else {
			a.log.Info("Telegram bağlantısı doğrulandı")
		}
	}

	// Root context — canceled on SIGINT/SIGTERM or programmatically (self-update).
	sigCtx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()
	rootCtx, cancel := context.WithCancel(sigCtx)
	a.cancel = cancel
	defer cancel()

	g, ctx := errgroup.WithContext(rootCtx)

	// Listener supervisor
	g.Go(func() error {
		return a.manager.Run(ctx, a.bus)
	})

	// Event dispatcher
	g.Go(func() error {
		return a.dispatch(ctx)
	})

	// Telegram bot polling
	g.Go(func() error {
		if !a.cfg.IsConfigured() {
			<-ctx.Done()
			return nil
		}
		return a.bot.Poll(ctx)
	})

	// Network monitor
	g.Go(func() error {
		return a.netMon.Run(ctx, a.bus)
	})

	// Heartbeat
	g.Go(func() error {
		return a.heartbeat(ctx)
	})

	// Offline queue flusher
	g.Go(func() error {
		return a.flushOfflineQueue(ctx)
	})

	// Daily pruning
	g.Go(func() error {
		return a.prune(ctx)
	})

	// Periodic update checker
	g.Go(func() error {
		return a.updateChecker(ctx)
	})

	// Optional Prometheus /metrics endpoint
	g.Go(func() error {
		return a.metricsServer(ctx)
	})

	// Announce boot (and report if we just self-updated)
	a.announceBoot()
	a.notifyIfUpdated()

	a.log.Info("🏰 Kanije Kalesi hazır — bekçi göreve başladı")

	if err := g.Wait(); err != nil && err != context.Canceled {
		return err
	}

	// Graceful shutdown
	a.shutdown()
	return nil
}

// dispatch reads events from the bus, stores them, and sends notifications.
func (a *App) dispatch(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case ev, ok := <-a.bus.Events():
			if !ok {
				return nil
			}
			a.handleEvent(ctx, ev)
		}
	}
}

func (a *App) handleEvent(ctx context.Context, ev event.Event) {
	// Store event
	if err := a.store.SaveEvent(ctx, ev); err != nil {
		a.log.Warn("olay kaydedilemedi", "err", err, "type", ev.Type)
	}

	// Stamp the OS for the notification (e.g. "💻 DESKTOP-X · Windows 11").
	if ev.OS == "" {
		ev.OS = a.osName
	}

	// Track last event for /status
	a.lastEvent.Store(&ev)

	// Check if this trigger is enabled
	trig, ok := a.cfg.GetTrigger(string(ev.Type))
	if !ok || !trig.Enabled {
		return
	}

	// Quiet hours: during the configured window only alert/critical events are
	// pushed; everything else is stored silently (no notify, no media capture).
	if ev.Severity < event.SeverityAlert && a.cfg.InQuietHours(ev.Timestamp) {
		return
	}

	// GeoIP enrichment: annotate events carrying a public source IP with the
	// country/city/ISP so a failed login reads e.g. "🇷🇺 Moscow, Russia".
	if ev.SourceIP != "" && a.cfg.GeoIPEnabled() {
		gctx, gcancel := context.WithTimeout(ctx, 4*time.Second)
		if info, ok := a.geo.Lookup(gctx, ev.SourceIP); ok {
			if ev.Extra == nil {
				ev.Extra = make(map[string]string, 2)
			}
			ev.Extra["📍 Konum"] = info.Summary()
			if info.ISP != "" {
				ev.Extra["🌐 Sağlayıcı"] = info.ISP
			}
		}
		gcancel()
	}

	// Capture media if configured
	if trig.CaptureCamera {
		if data, err := a.camera.Capture(ctx); err == nil {
			ev.Attachments = append(ev.Attachments, event.Attachment{
				Type:    event.AttachmentPhoto,
				Data:    data,
				Caption: "📷 Otomatik çekim — " + ev.Timestamp.Format("15:04:05"),
			})
		} else {
			a.log.Debug("kamera çekimi başarısız", "err", err)
		}
	}

	if trig.CaptureScreenshot {
		if data, err := a.screen.Capture(ctx); err == nil {
			ev.Attachments = append(ev.Attachments, event.Attachment{
				Type:    event.AttachmentScreenshot,
				Data:    data,
				Caption: "🖥️ Otomatik ekran görüntüsü — " + ev.Timestamp.Format("15:04:05"),
			})
		} else {
			a.log.Debug("ekran görüntüsü başarısız", "err", err)
		}
	}

	// Send notifications. Telegram and webhooks are independent — webhooks fire
	// even if Telegram isn't configured.
	sendCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if a.cfg.IsConfigured() {
		if err := a.bot.SendEvent(sendCtx, ev); err != nil {
			a.log.Warn("bildirim gönderilemedi, kuyruğa alındı", "err", err)
			// Queue for offline delivery
			if qErr := a.store.SavePendingMessage(ctx, telegram.FormatEvent(ev)); qErr != nil {
				a.log.Error("kuyruk yazma hatası", "err", qErr)
			}
		}
	}

	if a.webhooks.Enabled() {
		a.webhooks.Send(sendCtx, ev)
	}
}

// heartbeat sends periodic status messages. The enabled flag is re-checked on
// every tick so changes made at runtime via /kurulum take effect without a restart.
func (a *App) heartbeat(ctx context.Context) error {
	interval := a.cfg.HeartbeatInterval()
	if interval <= 0 {
		interval = 6 * time.Hour
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if a.cfg.HeartbeatEnabled() && !a.cfg.InQuietHours(time.Now()) {
				a.sendHeartbeat(ctx)
			}
		}
	}
}

func (a *App) sendHeartbeat(ctx context.Context) {
	if !a.cfg.IsConfigured() {
		return
	}

	info := sysinfo.Collect()
	uptime := time.Since(a.startedAt)

	var diskFree, diskTotal uint64
	if len(info.Disks) > 0 {
		diskFree = info.Disks[0].Free
		diskTotal = info.Disks[0].Total
	}

	total, _ := a.store.CountEvents(ctx)
	hostname, _ := os.Hostname()
	text := telegram.FormatHeartbeat(uptime, diskFree, diskTotal, total, hostname, a.osName, info.Platform)

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := a.bot.SendMessage(timeoutCtx, text); err != nil {
		a.log.Warn("heartbeat gönderilemedi", "err", err)
	}
}

// flushOfflineQueue periodically tries to deliver queued messages. Messages are
// removed only after a confirmed send, so a failure mid-batch leaves the unsent
// remainder queued for the next attempt — at-least-once delivery in FIFO order,
// with no data loss on transient errors or a crash.
func (a *App) flushOfflineQueue(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			msgs, err := a.store.PendingMessages(ctx)
			if err != nil || len(msgs) == 0 {
				continue
			}

			a.log.Info("bekleyen mesajlar gönderiliyor", "adet", len(msgs))

			sendCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			sent := flushPending(sendCtx, msgs, a.bot.SendMessage)
			cancel()

			if len(sent) < len(msgs) {
				a.log.Warn("bazı bekleyen mesajlar gönderilemedi, kuyrukta korunuyor",
					"gönderilen", len(sent), "toplam", len(msgs))
			}
			if len(sent) > 0 {
				if err := a.store.DeletePendingMessages(ctx, sent); err != nil {
					a.log.Error("gönderilen mesajlar kuyruktan silinemedi", "err", err)
				}
			}
		}
	}
}

// flushPending sends each queued message via send, in order, and returns the IDs
// that were successfully delivered. It stops at the first failure so the unsent
// remainder is preserved by the caller — giving at-least-once, in-order delivery.
func flushPending(ctx context.Context, msgs []storage.PendingMessage, send func(context.Context, string) error) []int64 {
	var sent []int64
	for _, m := range msgs {
		if err := send(ctx, m.Text); err != nil {
			break
		}
		sent = append(sent, m.ID)
	}
	return sent
}

// prune removes old events daily.
func (a *App) prune(ctx context.Context) error {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			n, err := a.store.Prune(ctx, a.cfg.Storage.EventRetentionDays)
			if err != nil {
				a.log.Warn("temizleme hatası", "err", err)
			} else if n > 0 {
				a.log.Info("eski olaylar temizlendi", "adet", n)
			}
		}
	}
}

// updateChecker periodically checks GitHub Releases for a newer version. On a
// match it either self-installs (Windows, when auto_install) or notifies.
func (a *App) updateChecker(ctx context.Context) error {
	if !a.cfg.UpdateEnabled() {
		<-ctx.Done()
		return nil
	}
	interval := a.cfg.UpdateInterval()
	if interval <= 0 {
		interval = 24 * time.Hour
	}

	// Small initial delay so the check doesn't run at the exact moment of boot.
	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			if msg, newer := a.runUpdateCheck(ctx, false); newer && a.cfg.IsConfigured() {
				sendCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
				a.bot.SendMessage(sendCtx, msg)
				cancel()
			}
			timer.Reset(interval)
		}
	}
}

// checkForUpdate runs an on-demand check (e.g. the /guncelle command) and
// returns a user-facing message.
func (a *App) checkForUpdate(ctx context.Context, force bool) string {
	msg, _ := a.runUpdateCheck(ctx, force)
	return msg
}

// runUpdateCheck performs one check. When force (or auto_install) is set and the
// platform supports it, the newer version is downloaded and applied — the app
// then shuts down so the swap helper can restart it. Returns (message, isNewer).
func (a *App) runUpdateCheck(ctx context.Context, force bool) (string, bool) {
	rel, newer, err := a.updater.LatestNewer(ctx)
	if err != nil {
		return "❌ Güncelleme kontrolü başarısız: " + telegram.SafeText(err.Error()), false
	}
	if !newer {
		return "✅ Zaten en güncel sürümü kullanıyorsunuz (" + telegram.SafeText(a.version) + ").", false
	}

	if (force || a.cfg.UpdateAutoInstall()) && updater.CanSelfInstall() {
		if err := a.updater.SelfInstall(ctx, rel); err != nil {
			return "❌ Güncelleme indirilemedi: " + telegram.SafeText(err.Error()), true
		}
		a.updating.Store(true) // shutdown() bunu görüp "kapatılıyor" bildirimi atmaz
		// Let the reply send, then trigger graceful shutdown so the helper can
		// replace the binary and restart the scheduled task.
		go func() {
			time.Sleep(3 * time.Second)
			if a.cancel != nil {
				a.cancel()
			}
		}()
		return "🔄 <b>" + telegram.SafeText(rel.Version) + "</b> indirildi. Güncelleniyor — birazdan yeni sürümle döneceğim…", true
	}

	hint := "Kurmak için /guncelle yazın."
	if !updater.CanSelfInstall() {
		hint = "Güncellemek için sunucuda kurulum betiğini yeniden çalıştırın."
	}
	return "🆕 Yeni sürüm mevcut: <b>" + telegram.SafeText(rel.Version) + "</b>\n" + hint, true
}

// notifyIfUpdated compares the running version with the last recorded one and,
// if it changed, reports a successful update. The marker lives next to the config.
func (a *App) notifyIfUpdated() {
	if a.version == "dev" || a.version == "" {
		return
	}
	dir := filepath.Dir(a.cfg.FilePath())
	if dir == "" {
		dir = "."
	}
	marker := filepath.Join(dir, "version.txt")

	prev, _ := os.ReadFile(marker)
	prevVer := strings.TrimSpace(string(prev))

	if prevVer != "" && prevVer != a.version && a.cfg.IsConfigured() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		a.bot.SendMessage(ctx, "✅ Güncellendi: <b>"+telegram.SafeText(prevVer)+"</b> → <b>"+telegram.SafeText(a.version)+"</b>")
		cancel()
	}
	_ = os.WriteFile(marker, []byte(a.version), 0o600)
}

// metricsServer serves the optional Prometheus /metrics endpoint until ctx ends.
func (a *App) metricsServer(ctx context.Context) error {
	if !a.cfg.MetricsEnabled() {
		<-ctx.Done()
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", a.metricsHandler)
	srv := &http.Server{
		Addr:              a.cfg.MetricsAddr(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		sc, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_ = srv.Shutdown(sc)
		cancel()
	}()

	a.log.Info("Prometheus /metrics yayında", "adres", a.cfg.MetricsAddr())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Warn("metrics sunucu hatası", "err", err)
	}
	return nil
}

// metricsHandler writes the Prometheus text exposition format.
func (a *App) metricsHandler(w http.ResponseWriter, r *http.Request) {
	st := a.bus.Stats()
	uptime := int64(time.Since(a.startedAt).Seconds())
	total, _ := a.store.CountEvents(r.Context())

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	fmt.Fprintf(w, "# HELP kanije_events_received_total Bus'a kabul edilen toplam olay\n")
	fmt.Fprintf(w, "# TYPE kanije_events_received_total counter\nkanije_events_received_total %d\n", st.Received)
	fmt.Fprintf(w, "# HELP kanije_events_dropped_total Hiz limiti veya dolu buffer nedeniyle dusen olaylar\n")
	fmt.Fprintf(w, "# TYPE kanije_events_dropped_total counter\nkanije_events_dropped_total %d\n", st.Dropped)
	fmt.Fprintf(w, "# HELP kanije_events_deduped_total Yinelenme nedeniyle atlanan olaylar\n")
	fmt.Fprintf(w, "# TYPE kanije_events_deduped_total counter\nkanije_events_deduped_total %d\n", st.Deduped)
	fmt.Fprintf(w, "# HELP kanije_bus_pending Bus buffer'inda bekleyen olay\n")
	fmt.Fprintf(w, "# TYPE kanije_bus_pending gauge\nkanije_bus_pending %d\n", st.Pending)
	fmt.Fprintf(w, "# HELP kanije_events_stored_total Veritabanindaki toplam olay\n")
	fmt.Fprintf(w, "# TYPE kanije_events_stored_total gauge\nkanije_events_stored_total %d\n", total)
	fmt.Fprintf(w, "# HELP kanije_uptime_seconds Calisma suresi\n")
	fmt.Fprintf(w, "# TYPE kanije_uptime_seconds gauge\nkanije_uptime_seconds %d\n", uptime)
	fmt.Fprintf(w, "# HELP kanije_build_info Surum bilgisi\n")
	fmt.Fprintf(w, "# TYPE kanije_build_info gauge\nkanije_build_info{version=%q} 1\n", a.version)
}

// announceBoot sends a system boot notification.
func (a *App) announceBoot() {
	if !a.cfg.IsConfigured() {
		return
	}

	hostname, _ := os.Hostname()
	ev := event.New(event.TypeSystemBoot, "App")
	ev.Hostname = hostname

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.bot.SendEvent(ctx, ev)
}

// shutdown is called after the errgroup exits — performs cleanup.
func (a *App) shutdown() {
	a.log.Info("Kapatılıyor…")

	if a.cfg.IsConfigured() && !a.updating.Load() {
		hostname, _ := os.Hostname()
		ev := event.New(event.TypeSystemShutdown, "App")
		ev.Hostname = hostname

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		a.bot.SendEvent(ctx, ev)
	}

	if err := a.store.Close(); err != nil {
		a.log.Warn("veritabanı kapatma hatası", "err", err)
	}

	a.log.Info("Kanije Kalesi durduruldu. Hoşça kalın! 🏰")
}

// collectStatus returns a snapshot for the /status Telegram command.
func (a *App) collectStatus() telegram.StatusInfo {
	si := sysinfo.Collect()
	busStats := a.bus.Stats()

	info := telegram.StatusInfo{
		CPUPercent:  si.CPUPercent,
		MemPercent:  si.MemPercent,
		MemUsed:     si.MemUsed,
		MemTotal:    si.MemTotal,
		Uptime:      time.Since(a.startedAt),
		BusReceived: busStats.Received,
		BusDropped:  busStats.Dropped,
		LastEvent:   a.lastEvent.Load(),
	}

	for _, d := range si.Disks {
		info.Disks = append(info.Disks, telegram.DiskInfo{
			Path:  d.Path,
			Free:  d.Free,
			Total: d.Total,
		})
	}

	return info
}
