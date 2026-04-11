package capture

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"log/slog"
	"time"

	"github.com/kbinani/screenshot"
)

// Screenshotter captures the primary display using kbinani/screenshot.
// This package uses platform-native APIs (GDI on Windows, X11/XShm on Linux)
// and does not require CGo on Windows. On Linux it may require libX11.
type Screenshotter struct {
	jpegQuality int
	log         *slog.Logger
}

// NewScreenshotter creates a Screenshotter with the given JPEG quality (1-100).
func NewScreenshotter(jpegQuality int, log *slog.Logger) *Screenshotter {
	if jpegQuality <= 0 || jpegQuality > 100 {
		jpegQuality = 75
	}
	return &Screenshotter{jpegQuality: jpegQuality, log: log}
}

// Capture takes a screenshot of all displays merged into one image.
// Returns raw JPEG bytes. Respects context timeout.
func (s *Screenshotter) Capture(ctx context.Context) ([]byte, error) {
	type result struct {
		data []byte
		err  error
	}

	ch := make(chan result, 1)

	go func() {
		n := screenshot.NumActiveDisplays()
		if n == 0 {
			ch <- result{nil, fmt.Errorf("ekran bulunamadı")}
			return
		}

		// Capture primary display (index 0)
		bounds := screenshot.GetDisplayBounds(0)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			ch <- result{nil, fmt.Errorf("ekran görüntüsü hatası: %w", err)}
			return
		}

		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: s.jpegQuality}); err != nil {
			ch <- result{nil, fmt.Errorf("JPEG encode hatası: %w", err)}
			return
		}

		ch <- result{buf.Bytes(), nil}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("ekran görüntüsü iptal edildi: %w", ctx.Err())
	case r := <-ch:
		return r.data, r.err
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("ekran görüntüsü zaman aşımı")
	}
}
