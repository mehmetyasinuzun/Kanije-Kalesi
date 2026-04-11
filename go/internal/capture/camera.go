// Package capture provides camera and screenshot capture functionality.
// Camera capture uses ffmpeg as a subprocess — no CGo, no OpenCV dependency,
// fully cross-platform and embeddable.
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
	"time"
)

// CameraConfig holds camera capture settings.
type CameraConfig struct {
	FFmpegPath   string
	DeviceIndex  int
	DeviceName   string // Windows dshow: friendly name from `ffmpeg -list_devices`
	Width, Height int
	WarmupFrames int
	JPEGQuality  int
}

// Camera captures frames from a webcam using ffmpeg.
// Only one capture at a time is permitted (guarded by mu).
type Camera struct {
	cfg CameraConfig
	mu  sync.Mutex
	log *slog.Logger
}

// NewCamera creates a Camera instance with the given configuration.
func NewCamera(cfg CameraConfig, log *slog.Logger) *Camera {
	return &Camera{cfg: cfg, log: log}
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

	args := c.buildFFmpegArgs()
	c.log.Debug("kamera komutu", "args", args)

	cmd := exec.CommandContext(captureCtx, c.cfg.FFmpegPath, args...)

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

// buildFFmpegArgs constructs the ffmpeg argument list for the current platform.
func (c *Camera) buildFFmpegArgs() []string {
	var inputFormat, device string

	switch runtime.GOOS {
	case "windows":
		inputFormat = "dshow"
		if c.cfg.DeviceName != "" {
			device = "video=" + c.cfg.DeviceName
		} else {
			// fallback: use device index — requires listing first
			device = "video=0" // simplification
		}

	case "linux":
		inputFormat = "v4l2"
		device = "/dev/video" + strconv.Itoa(c.cfg.DeviceIndex)

	case "darwin":
		inputFormat = "avfoundation"
		device = strconv.Itoa(c.cfg.DeviceIndex)

	default:
		inputFormat = "v4l2"
		device = "/dev/video" + strconv.Itoa(c.cfg.DeviceIndex)
	}

	// We capture warmupFrames+1 frames and take the last one to avoid
	// dark/blurry initial frames. The select filter picks frame ≥ warmupFrames.
	selectFilter := fmt.Sprintf("select=gte(n\\,%d)", c.cfg.WarmupFrames)

	return []string{
		"-hide_banner",
		"-loglevel", "error",      // Suppress verbose output
		"-f", inputFormat,
		"-video_size", fmt.Sprintf("%dx%d", c.cfg.Width, c.cfg.Height),
		"-i", device,
		"-vframes", strconv.Itoa(c.cfg.WarmupFrames + 1),
		"-vf", selectFilter,
		"-q:v", strconv.Itoa(ffmpegQuality(c.cfg.JPEGQuality)),
		"-f", "image2",
		"-vcodec", "mjpeg",
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
	cmd.Stderr = &stderr

	cmd.Run() // Expected to fail (no input)

	return parseDeviceList(stderr.String()), nil
}

func parseDeviceList(output string) []string {
	var devices []string
	for _, line := range splitLines(output) {
		if containsAny(line, []string{"DirectShow video", "dshow"}) {
			continue
		}
		if idx := indexByte(line, '"'); idx >= 0 {
			end := indexByteFrom(line, '"', idx+1)
			if end > idx {
				devices = append(devices, line[idx+1:end])
			}
		}
	}
	return devices
}

func listLinuxDevices() []string {
	var devices []string
	for i := 0; i < 10; i++ {
		path := "/dev/video" + strconv.Itoa(i)
		if fileExists(path) {
			devices = append(devices, path)
		}
	}
	return devices
}

// ---- string helpers ----

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if len(sub) > 0 && len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func indexByteFrom(s string, b byte, from int) int {
	for i := from; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func fileExists(path string) bool {
	_, err := exec.LookPath(path)
	return err == nil
}
