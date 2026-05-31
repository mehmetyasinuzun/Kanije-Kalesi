package capture

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"math"
)

// FrameDiff decodes two JPEG frames and returns the mean absolute per-pixel
// luminance difference (0-255 scale); higher means more changed. The frames must
// be the same size. Pixels are sampled every 4th row/column for speed. This backs
// the camera motion detector — a result above a threshold means "something moved".
func FrameDiff(a, b []byte) (float64, error) {
	ia, err := jpeg.Decode(bytes.NewReader(a))
	if err != nil {
		return 0, fmt.Errorf("kare A çözülemedi: %w", err)
	}
	ib, err := jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		return 0, fmt.Errorf("kare B çözülemedi: %w", err)
	}

	ba, bb := ia.Bounds(), ib.Bounds()
	if ba.Dx() != bb.Dx() || ba.Dy() != bb.Dy() {
		return 0, fmt.Errorf("kare boyutları farklı")
	}

	const step = 4
	var total float64
	var count int
	for y := ba.Min.Y; y < ba.Max.Y; y += step {
		for x := ba.Min.X; x < ba.Max.X; x += step {
			total += math.Abs(luminance(ia, x, y) - luminance(ib, x, y))
			count++
		}
	}
	if count == 0 {
		return 0, nil
	}
	return total / float64(count), nil
}

// luminance returns the Rec.601 luma of a pixel on a 0-255 scale.
func luminance(img image.Image, x, y int) float64 {
	r, g, b, _ := img.At(x, y).RGBA() // each 0-65535
	return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 257.0
}
