package capture

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func grayJPEG(t *testing.T, level uint8, w, h int) []byte {
	t.Helper()
	img := image.NewGray(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetGray(x, y, color.Gray{Y: level})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95}); err != nil {
		t.Fatalf("jpeg encode: %v", err)
	}
	return buf.Bytes()
}

func TestFrameDiffIdentical(t *testing.T) {
	a := grayJPEG(t, 100, 64, 48)
	d, err := FrameDiff(a, a)
	if err != nil {
		t.Fatal(err)
	}
	if d > 2 { // same frame → ~0 (small JPEG tolerance)
		t.Errorf("aynı kare farkı düşük olmalı, alınan %.2f", d)
	}
}

func TestFrameDiffDifferent(t *testing.T) {
	black := grayJPEG(t, 0, 64, 48)
	white := grayJPEG(t, 255, 64, 48)
	d, err := FrameDiff(black, white)
	if err != nil {
		t.Fatal(err)
	}
	if d < 200 { // black vs white → near-max difference
		t.Errorf("siyah-beyaz farkı yüksek olmalı, alınan %.2f", d)
	}
}

func TestFrameDiffSizeMismatch(t *testing.T) {
	a := grayJPEG(t, 100, 64, 48)
	b := grayJPEG(t, 100, 32, 24)
	if _, err := FrameDiff(a, b); err == nil {
		t.Error("farklı boyutlu kareler hata vermeli")
	}
}
