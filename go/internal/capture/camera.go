// Package capture provides camera and screenshot capture functionality.
// Camera capture uses ffmpeg as a subprocess — no CGo, no OpenCV dependency,
// fully cross-platform and embeddable.
package capture

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

// CameraConfig holds camera capture settings.
type CameraConfig struct {
	FFmpegPath    string
	DeviceIndex   int
	DeviceName    string // Windows dshow: friendly name from `ffmpeg -list_devices`
	Width, Height int
	WarmupFrames  int
	JPEGQuality   int
}

// Camera captures frames from a webcam using ffmpeg.
// Only one capture at a time is permitted (guarded by mu).
type Camera struct {
	cfg            CameraConfig
	mu             sync.Mutex
	log            *slog.Logger
	ffmpegOverride atomic.Value // string — runtime-provisioned ffmpeg path (auto-download)
}

// NewCamera creates a Camera instance with the given configuration.
func NewCamera(cfg CameraConfig, log *slog.Logger) *Camera {
	return &Camera{cfg: cfg, log: log}
}

// SetFFmpegPath updates the ffmpeg binary path at runtime (e.g. after the agent
// auto-downloads ffmpeg). Safe for concurrent use.
func (c *Camera) SetFFmpegPath(path string) {
	if path != "" {
		c.ffmpegOverride.Store(path)
	}
}

// ffmpegPath returns the runtime override if set, else the configured path, else
// bare "ffmpeg" (PATH lookup).
func (c *Camera) ffmpegPath() string {
	if v := c.ffmpegOverride.Load(); v != nil {
		if s, _ := v.(string); s != "" {
			return s
		}
	}
	if c.cfg.FFmpegPath != "" {
		return c.cfg.FFmpegPath
	}
	return "ffmpeg"
}

// Capture takes a single photo and returns raw JPEG bytes.
// Returns an error if ffmpeg is unavailable or the device is busy.
func (c *Camera) Capture(ctx context.Context) ([]byte, error) {
	if !c.mu.TryLock() {
		return nil, fmt.Errorf("kamera meşgul — lütfen bekleyin")
	}
	defer c.mu.Unlock()

	captureCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	device, err := c.resolveDevice()
	if err != nil {
		return nil, err
	}
	args := c.buildFFmpegArgs(device)
	c.log.Debug("kamera komutu", "args", args)

	cmd := exec.CommandContext(captureCtx, c.ffmpegPath(), args...)
	sysproc.Hide(cmd) // GUI binary'de ffmpeg konsol penceresi parlatmasın

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg hatası: %w\nstderr: %s", err, truncate(stderr.String(), 300))
	}

	if stdout.Len() == 0 {
		return nil, fmt.Errorf("ffmpeg boş çıktı verdi")
	}

	c.log.Debug("kamera başarılı", "boyut", stdout.Len())
	return stdout.Bytes(), nil
}

// resolveDevice picks the ffmpeg input-device string for this platform. On
// Windows (dshow) ffmpeg needs the camera's FRIENDLY NAME, not a numeric index —
// "video=0" never works (that was the bug). If the user hasn't configured a name
// we auto-detect the first available camera via -list_devices.
func (c *Camera) resolveDevice() (string, error) {
	switch runtime.GOOS {
	case "windows":
		name := c.cfg.DeviceName
		if name == "" {
			devs, _ := ListDevices(c.ffmpegPath())
			if len(devs) == 0 {
				return "", fmt.Errorf("kamera bulunamadı — kamera bağlı ve açık mı? Değilse /kurulum → 📷 Kamera'dan cihaz adını elle girin")
			}
			name = devs[0]
			c.log.Info("kamera otomatik seçildi", "cihaz", name)
		}
		return "video=" + name, nil
	case "darwin":
		return strconv.Itoa(c.cfg.DeviceIndex), nil
	default: // linux + others
		return "/dev/video" + strconv.Itoa(c.cfg.DeviceIndex), nil
	}
}

// buildFFmpegArgs constructs the ffmpeg argument list for the current platform.
func (c *Camera) buildFFmpegArgs(device string) []string {
	var inputFormat string

	switch runtime.GOOS {
	case "windows":
		inputFormat = "dshow"
	case "darwin":
		inputFormat = "avfoundation"
	default: // linux + others
		inputFormat = "v4l2"
	}

	// Drop the first WarmupFrames frames (dark/blurry sensor warm-up) then emit
	// EXACTLY ONE frame. select=gte(n,warmup) passes frames from index warmup on;
	// -frames:v 1 then takes the first of those.
	selectFilter := fmt.Sprintf("select=gte(n\\,%d)", c.cfg.WarmupFrames)

	return []string{
		"-hide_banner",
		"-loglevel", "error", // Suppress verbose output
		"-f", inputFormat,
		"-video_size", fmt.Sprintf("%dx%d", c.cfg.Width, c.cfg.Height),
		"-i", device,
		"-vf", selectFilter,
		"-frames:v", "1", // exactly one output frame (was -vframes warmup+1)
		"-q:v", strconv.Itoa(ffmpegQuality(c.cfg.JPEGQuality)),
		"-c:v", "mjpeg",
		// image2pipe streams a single JPEG to stdout. The plain "image2" muxer is
		// a FILE writer and refuses a pipe with >1 frame ("Cannot write more than
		// one file with the same name") — that was the bug.
		"-f", "image2pipe",
		"pipe:1", // Output to stdout
	}
}

// ffmpegQuality converts JPEG quality (1-100) to ffmpeg's q:v scale (1-31).
// Higher JPEG quality = lower ffmpeg q:v value.
func ffmpegQuality(jpegQuality int) int {
	if jpegQuality <= 0 {
		jpegQuality = 85
	}
	if jpegQuality > 100 {
		jpegQuality = 100
	}
	// Map 100→1, 1→31
	q := 31 - int(float64(jpegQuality-1)/99.0*30)
	if q < 1 {
		q = 1
	}
	return q
}

// ListDevices returns the list of available camera devices on this system.
// Useful for the Telegram setup wizard to let the user pick a camera.
func ListDevices(ffmpegPath string) ([]string, error) {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	var args []string
	switch runtime.GOOS {
	case "windows":
		args = []string{"-hide_banner", "-f", "dshow", "-list_devices", "true", "-i", "dummy"}
	case "linux":
		// List /dev/video* devices
		return listLinuxDevices(), nil
	default:
		return nil, fmt.Errorf("cihaz listeleme bu platformda desteklenmiyor")
	}

	var stderr bytes.Buffer
	cmd := exec.Command(ffmpegPath, args...)
	sysproc.Hide(cmd)
	cmd.Stderr = &stderr

	cmd.Run() // Expected to fail (no input)

	return parseDshowDevices(stderr.String(), "video"), nil
}

func listLinuxDevices() []string {
	var devices []string
	for i := 0; i < 10; i++ {
		path := "/dev/video" + strconv.Itoa(i)
		if _, err := os.Stat(path); err == nil {
			devices = append(devices, path)
		}
	}
	return devices
}

// truncate shortens s to at most n bytes, trimming on a valid UTF-8 boundary so
// a multi-byte rune is never split. Used to cap ffmpeg stderr in error messages.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return strings.ToValidUTF8(s[:n], "") + "…"
}
