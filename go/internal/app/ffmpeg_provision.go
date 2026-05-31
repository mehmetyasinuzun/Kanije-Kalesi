package app

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/kanije-kalesi/kanije/internal/config"
	"github.com/kanije-kalesi/kanije/internal/ffmpeg"
	"github.com/kanije-kalesi/kanije/internal/notifier/telegram"
)

// appBaseDir is the agent's working directory — the config's folder, or the
// executable's folder if config is path-less. ffmpeg.exe, schedules.json and
// markers all live here.
func appBaseDir(cfg *config.Config) string {
	dir := filepath.Dir(cfg.FilePath())
	if dir == "" || dir == "." {
		if exe, err := os.Executable(); err == nil {
			dir = filepath.Dir(exe)
		}
	}
	return dir
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// provisionFFmpeg makes camera/audio "just work" without a manual ffmpeg install.
// If ffmpeg is already resolvable it's a no-op; otherwise it downloads a static
// build (Windows) in the background, wires it into the capture modules, and tells
// the owner. On failure it explains how to install ffmpeg manually.
func (a *App) provisionFFmpeg(ctx context.Context) error {
	baseDir := appBaseDir(a.cfg)
	if ffmpeg.Resolve(a.cfg.Camera.FFmpegPath, baseDir) != "" {
		return nil // already available — nothing to do
	}
	if !a.cfg.IsConfigured() {
		<-ctx.Done()
		return nil
	}

	a.notify(ctx, "🎥 Kamera/ses için <b>ffmpeg indiriliyor</b> (~80 MB, tek seferlik)… Birazdan hazır olacak.")

	path, err := ffmpeg.Download(baseDir)
	if err != nil {
		a.log.Warn("ffmpeg indirilemedi", "err", err)
		a.notify(ctx, "⚠️ ffmpeg otomatik indirilemedi: "+telegram.SafeText(err.Error())+
			"\n\nKamera/ses için ffmpeg gerekir. Elle kurmak için: <code>winget install ffmpeg</code> ya da <code>choco install ffmpeg</code> (sonra /yeniden).")
		return nil
	}

	a.camera.SetFFmpegPath(path)
	if a.audioRec != nil {
		a.audioRec.SetFFmpegPath(path)
	}
	a.log.Info("ffmpeg sağlandı", "yol", path)
	a.notify(ctx, "✅ <b>ffmpeg hazır</b> — kamera (/foto, /panik) ve ses (/seskayit) artık çalışıyor.")
	return nil
}

// notify sends a one-off owner message with its own short timeout (so a slow
// send never blocks the caller).
func (a *App) notify(ctx context.Context, msg string) {
	if !a.cfg.IsConfigured() {
		return
	}
	sendCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	a.bot.SendMessage(sendCtx, msg)
}
