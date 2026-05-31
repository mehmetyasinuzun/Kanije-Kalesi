package capture

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

const (
	// DefaultRecordSeconds is used when /seskayit is called with no argument.
	DefaultRecordSeconds = 30
	// MaxRecordSeconds caps a single recording (10 minutes) so a runaway or
	// malicious request can't fill the disk / hold the mic indefinitely.
	MaxRecordSeconds = 600
)

// AudioConfig holds microphone capture settings.
type AudioConfig struct {
	FFmpegPath string // ffmpeg binary; empty = "ffmpeg" (PATH)
	DeviceName string // Windows dshow mic name; empty = auto-detect first mic
	Bitrate    string // MP3 bitrate, e.g. "96k"; empty = "96k"
}

// AudioRecorder records short microphone clips on demand via ffmpeg (no CGo,
// cross-platform; the subprocess console window is hidden on Windows).
// Only one recording runs at a time (guarded by mu).
type AudioRecorder struct {
	cfg            AudioConfig
	mu             sync.Mutex
	log            *slog.Logger
	ffmpegOverride atomic.Value // string — runtime-provisioned ffmpeg path
}

// NewAudioRecorder creates an AudioRecorder with the given configuration.
func NewAudioRecorder(cfg AudioConfig, log *slog.Logger) *AudioRecorder {
	return &AudioRecorder{cfg: cfg, log: log}
}

// SetFFmpegPath updates the ffmpeg path at runtime (after auto-download).
func (a *AudioRecorder) SetFFmpegPath(path string) {
	if path != "" {
		a.ffmpegOverride.Store(path)
	}
}

// Record captures `seconds` of microphone audio and returns MP3 bytes.
// seconds is clamped to [1, MaxRecordSeconds]; 0 or negative means the default.
func (a *AudioRecorder) Record(ctx context.Context, seconds int) ([]byte, error) {
	if seconds <= 0 {
		seconds = DefaultRecordSeconds
	}
	if seconds > MaxRecordSeconds {
		seconds = MaxRecordSeconds
	}

	if !a.mu.TryLock() {
		return nil, fmt.Errorf("ses kaydı meşgul — lütfen bekleyin")
	}
	defer a.mu.Unlock()

	device, err := a.resolveDevice()
	if err != nil {
		return nil, err
	}

	// Allow the recording duration plus a generous startup/encode margin.
	recCtx, cancel := context.WithTimeout(ctx, time.Duration(seconds+20)*time.Second)
	defer cancel()

	args := a.buildArgs(device, seconds)
	a.log.Debug("ses kayıt komutu", "args", args)

	cmd := exec.CommandContext(recCtx, a.ffmpegPath(), args...)
	sysproc.Hide(cmd) // GUI binary'de ffmpeg konsol penceresi parlatmasın

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg ses hatası: %w\nstderr: %s", err, truncate(stderr.String(), 300))
	}
	if stdout.Len() == 0 {
		return nil, fmt.Errorf("ffmpeg boş ses çıktısı verdi (mikrofon erişimi/izni yok olabilir)")
	}

	a.log.Debug("ses kaydı başarılı", "boyut", stdout.Len(), "saniye", seconds)
	return stdout.Bytes(), nil
}

func (a *AudioRecorder) ffmpegPath() string {
	if v := a.ffmpegOverride.Load(); v != nil {
		if s, _ := v.(string); s != "" {
			return s
		}
	}
	if a.cfg.FFmpegPath != "" {
		return a.cfg.FFmpegPath
	}
	return "ffmpeg"
}

// buildArgs constructs the ffmpeg argument list for the current platform.
func (a *AudioRecorder) buildArgs(device string, seconds int) []string {
	var inputFormat, input string

	switch runtime.GOOS {
	case "windows":
		inputFormat = "dshow"
		input = "audio=" + device
	case "darwin":
		inputFormat = "avfoundation"
		input = ":" + device // avfoundation audio-only: ":<index>"
	default: // linux & others
		inputFormat = "alsa"
		input = device
	}

	bitrate := a.cfg.Bitrate
	if bitrate == "" {
		bitrate = "96k"
	}

	return []string{
		"-hide_banner",
		"-loglevel", "error",
		"-f", inputFormat,
		"-i", input,
		"-t", strconv.Itoa(seconds),
		"-ac", "1", // mono — smaller files, fine for voice
		"-ar", "44100",
		"-b:a", bitrate,
		"-f", "mp3",
		"pipe:1",
	}
}

// resolveDevice returns the input device to record from.
func (a *AudioRecorder) resolveDevice() (string, error) {
	if a.cfg.DeviceName != "" {
		return a.cfg.DeviceName, nil
	}
	switch runtime.GOOS {
	case "windows":
		devs := ListAudioDevices(a.ffmpegPath())
		if len(devs) == 0 {
			return "", fmt.Errorf("mikrofon bulunamadı — /kurulum → 🎤 Ses'ten cihaz adını girin")
		}
		return devs[0], nil
	case "darwin":
		return "0", nil // default audio input index
	default:
		return "default", nil // ALSA default device
	}
}

// ListAudioDevices returns the Windows dshow audio (microphone) device names.
// On non-Windows platforms it returns nil (devices are referenced by a fixed
// name like "default"). Used by the setup wizard and auto-detection.
func ListAudioDevices(ffmpegPath string) []string {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	if runtime.GOOS != "windows" {
		return nil
	}

	cmd := exec.Command(ffmpegPath, "-hide_banner", "-f", "dshow", "-list_devices", "true", "-i", "dummy")
	sysproc.Hide(cmd)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Run() // Expected to "fail" (dummy input) — device list goes to stderr.

	return parseDshowDevices(stderr.String(), "audio")
}
