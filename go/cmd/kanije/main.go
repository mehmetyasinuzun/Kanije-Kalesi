// Command kanije is the main entry point for Kanije Kalesi Security Sentinel.
//
// Usage:
//
//	kanije start [--config path]     Start the security monitor
//	kanije test  [--config path]     Test Telegram connection
//	kanije setup --token T --chat C  Bootstrap minimal configuration
//	kanije version                   Print version and exit
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strconv"

	"github.com/kanije-kalesi/sentinel/internal/config"
	"github.com/kanije-kalesi/sentinel/internal/lock"
	"github.com/kanije-kalesi/sentinel/internal/notifier/telegram"
	"github.com/lmittmann/tint"
)

// Set at build time: go build -ldflags "-X main.Version=1.0.0"
var (
	Version   = "dev"
	BuildDate = "unknown"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "HATA: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "start":
		return cmdStart(rest)
	case "test":
		return cmdTest(rest)
	case "setup":
		return cmdSetup(rest)
	case "version", "--version", "-v":
		fmt.Printf("Kanije Kalesi Sentinel %s (%s) — %s/%s\n",
			Version, BuildDate, runtime.GOOS, runtime.GOARCH)
		return nil
	case "help", "--help", "-h":
		printUsage()
		return nil
	default:
		return fmt.Errorf("bilinmeyen komut: %q\nKullanım için: kanije help", cmd)
	}
}

// cmdStart is the main operational command.
func cmdStart(args []string) error {
	cfgPath := parseFlag(args, "--config", defaultConfigPath())

	log := buildLogger("info")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("yapılandırma yüklenemedi: %w", err)
	}

	log = buildLogger(cfg.Logging.Level)

	log.Info("🏰 Kanije Kalesi Sentinel başlatılıyor",
		"versiyon", Version,
		"platform", runtime.GOOS+"/"+runtime.GOARCH,
		"config", cfgPath,
	)

	// Single-instance guard is handled inside app.New(), but we check here
	// to give a clear error message before importing the heavy app package.
	if cfg.Security.SingleInstance {
		l, err := lock.Acquire("main")
		if err != nil {
			if errors.Is(err, lock.ErrAlreadyRunning) {
				return fmt.Errorf("kanije zaten çalışıyor")
			}
			return err
		}
		defer l.Release()
	}

	// Import app here to keep the binary small when running other subcommands
	return startApp(cfg, log)
}

// cmdTest validates the Telegram configuration by calling getMe.
func cmdTest(args []string) error {
	cfgPath := parseFlag(args, "--config", defaultConfigPath())

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("yapılandırma yüklenemedi: %w", err)
	}

	if !cfg.IsConfigured() {
		return fmt.Errorf(
			"Telegram token veya chat ID yapılandırılmamış.\n" +
				"  kanije setup --token <TOKEN> --chat <CHAT_ID>\n" +
				"veya ortam değişkenleri: KANIJE_BOT_TOKEN, KANIJE_CHAT_ID")
	}

	fmt.Println("Telegram bağlantısı test ediliyor…")

	log := buildLogger("error")
	client := telegram.NewClient(cfg.Telegram.BotToken, cfg.Telegram.SendTimeoutSec, log)

	user, err := client.GetMe(context.Background())
	if err != nil {
		return fmt.Errorf("bağlantı başarısız: %w", err)
	}

	fmt.Printf("✅ Bağlantı başarılı!\n")
	fmt.Printf("   Bot adı : @%s\n", user.Username)
	fmt.Printf("   Bot ID  : %d\n", user.ID)
	fmt.Printf("   Chat ID : %d\n", cfg.Telegram.ChatID)
	return nil
}

// cmdSetup writes minimal credentials to a config file.
// After this, users configure everything else via /kurulum in Telegram.
func cmdSetup(args []string) error {
	token  := parseFlag(args, "--token", os.Getenv("KANIJE_BOT_TOKEN"))
	chatID := parseFlag(args, "--chat", os.Getenv("KANIJE_CHAT_ID"))
	cfgPath := parseFlag(args, "--config", "config.toml")

	if token == "" {
		return fmt.Errorf(
			"--token gerekli.\n" +
				"BotFather'dan aldığınız token'ı girin:\n" +
				"  kanije setup --token <TOKEN> --chat <CHAT_ID>")
	}
	if chatID == "" {
		return fmt.Errorf(
			"--chat gerekli.\n" +
				"@userinfobot veya @RawDataBot ile chat ID'nizi öğrenin:\n" +
				"  kanije setup --token <TOKEN> --chat <CHAT_ID>")
	}

	id, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return fmt.Errorf("geçersiz chat ID: %q — sayı olmalı", chatID)
	}

	cfg := config.Defaults()
	cfg.Telegram.BotToken = token
	cfg.Telegram.ChatID   = id
	cfg.SetFilePath(cfgPath)

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("yapılandırma kaydedilemedi: %w", err)
	}

	fmt.Printf("✅ Yapılandırma kaydedildi: %s\n\n", cfgPath)
	fmt.Printf("Bağlantıyı test edin:\n")
	fmt.Printf("  kanije test\n\n")
	fmt.Printf("İzlemeyi başlatın:\n")
	fmt.Printf("  kanije start\n\n")
	fmt.Printf("Telegram botunuza /kurulum yazarak tüm ayarları yapılandırabilirsiniz.\n")
	return nil
}

// ---- Helpers ----

func printUsage() {
	fmt.Print(`
🏰 Kanije Kalesi Sentinel — Güvenlik İzleme Aracı

Kullanım:
  kanije start [--config <yol>]               İzlemeyi başlat
  kanije test  [--config <yol>]               Telegram bağlantısını test et
  kanije setup --token <T> --chat <C>         Telegram yapılandırmasını kaydet
  kanije version                              Versiyon bilgisi

Hızlı başlangıç:
  1. kanije setup --token <BOT_TOKEN> --chat <CHAT_ID>
  2. kanije test
  3. kanije start
  4. Telegram'da botunuza /kurulum yazarak ayarları yapın

Ortam değişkenleri:
  KANIJE_BOT_TOKEN   Telegram bot token'ı (BotFather'dan)
  KANIJE_CHAT_ID     Telegram chat ID'si
  KANIJE_LOG_LEVEL   Log seviyesi: debug|info|warn|error
  KANIJE_DB_PATH     SQLite veritabanı yolu

`)
}

func parseFlag(args []string, flag, defaultVal string) string {
	for i, arg := range args {
		if arg == flag && i+1 < len(args) {
			return args[i+1]
		}
		if len(arg) > len(flag)+1 && arg[:len(flag)+1] == flag+"=" {
			return arg[len(flag)+1:]
		}
	}
	return defaultVal
}

func defaultConfigPath() string {
	if _, err := os.Stat("config.toml"); err == nil {
		return "config.toml"
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.toml"
	}
	return home + "/.kanije/config.toml"
}

func buildLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      lvl,
		TimeFormat: "15:04:05",
	}))
}
