// Package app is the application orchestrator. It wires all modules together,
// manages their lifecycle, and handles graceful shutdown.
package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanije-kalesi/kanije/internal/capture"
	"github.com/kanije-kalesi/kanije/internal/config"
	"github.com/kanije-kalesi/kanije/internal/event"
	"github.com/kanije-kalesi/kanije/internal/listener"
	applock "github.com/kanije-kalesi/kanije/internal/lock"
	"github.com/kanije-kalesi/kanije/internal/network"
	"github.com/kanije-kalesi/kanije/internal/notifier/telegram"
	"github.com/kanije-kalesi/kanije/internal/storage"
	"github.com/kanije-kalesi/kanije/internal/sysinfo"
	"golang.org/x/sync/errgroup"
)

// App is the top-level application object. Create with New(), run with Run().
type App struct {
	cfg     *config.Config
	log     *slog.Logger
	store   storage.Storage
	bus     *event.Bus
	bot     *telegram.Bot
	manager *listener.Manager
	netMon  *network.Monitor
	camera  *capture.Camera
	screen  *capture.Screenshotter
	instLock applock.Releaser

	// Metrics
	startedAt time.Time
	lastEvent *event.Event
}

// New creates and wires the application. Call Run() to start.
func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	// Single-instance guard (handled in main, but stored here for Release on shutdown)
	_ = applock.ErrAlreadyRunning // ensure import used

	// Storage
	store, err := storage.NewSQLite(cfg.Storage.DBPath)
	if err != nil {
		return nil, fmt.Errorf("veritabanı açılamadı: %w", err)
	}

	// Event bus
	bus := event.NewBus(event.BusConfig{
		BufferSize:   512,
		MaxPerMinute: cfg.Security.MaxEventsPerMinute,
		DedupWindow:  time.Duration(cfg.Security.DedupWindowSec) * time.Second,
	})

	// Telegram client
	tgClient := telegram.NewClient(
		cfg.Telegram.BotToken,
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

	// Screenshot
	screen := capture.NewScreenshotter(cfg.Screenshot.JPEGQuality, log.With("module", "screenshot"))

	app := &App{
		cfg:    cfg,
		log:    log,
		store:  store,
		bus:    bus,
		camera: cam,
		screen: screen,
	}

	// Bot
	bot := telegram.NewBot(telegram.BotConfig{
		Config:  cfg,
		Client:  tgClient,
		Wizard:  wizard,
		Store:   store,
		Log:     log.With("module", "bot"),
		LockScreen:    lockScreen,
		CapturePhoto:  cam.Capture,
		CaptureScreen: screen.Capture,
		GetStatus:     app.collectStatus,
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

	// Root context — cancelled on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

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

	// Announce boot
	a.announceBoot()

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

	// Track last event for /status
	a.lastEvent = &ev

	// Check if this trigger is enabled
	trig, ok := a.cfg.GetTrigger(string(ev.Type))
	if !ok || !trig.Enabled {
		return
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

	// Send notification
	if !a.cfg.IsConfigured() {
		return
	}

	sendCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := a.bot.SendEvent(sendCtx, ev); err != nil {
		a.log.Warn("bildirim gönderilemedi, kuyruğa alındı", "err", err)
		// Queue for offline delivery
		if qErr := a.store.SavePendingMessage(ctx, telegram.FormatEvent(ev)); qErr != nil {
			a.log.Error("kuyruk yazma hatası", "err", qErr)
		}
	}
}

// heartbeat sends periodic status messages.
func (a *App) heartbeat(ctx context.Context) error {
	if !a.cfg.Heartbeat.Enabled {
		<-ctx.Done()
		return nil
	}

	interval := time.Duration(a.cfg.Heartbeat.IntervalHours) * time.Hour
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			a.sendHeartbeat(ctx)
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
		diskFree  = info.Disks[0].Free
		diskTotal = info.Disks[0].Total
	}

	total, _ := a.store.CountEvents(ctx)
	text := telegram.FormatHeartbeat(uptime, diskFree, diskTotal, total, info.Platform)

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := a.bot.SendMessage(timeoutCtx, text); err != nil {
		a.log.Warn("heartbeat gönderilemedi", "err", err)
	}
}

// flushOfflineQueue periodically tries to send queued messages.
func (a *App) flushOfflineQueue(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			msgs, err := a.store.PopPendingMessages(ctx)
			if err != nil || len(msgs) == 0 {
				continue
			}

			a.log.Info("bekleyen mesajlar gönderiliyor", "adet", len(msgs))

			sendCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			for _, m := range msgs {
				if err := a.bot.SendMessage(sendCtx, m.Text); err != nil {
					// Put back on failure
					a.store.SavePendingMessage(ctx, m.Text)
					break
				}
			}
			cancel()
		}
	}
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

	if a.cfg.IsConfigured() {
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

	if a.instLock != nil {
		a.instLock.Release()
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
		LastEvent:   a.lastEvent,
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
