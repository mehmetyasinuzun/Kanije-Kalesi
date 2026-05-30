package updater

import (
	"runtime"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"v1.1.2", "v1.1.1", 1},
		{"v1.1.1", "v1.1.2", -1},
		{"v1.1.1", "v1.1.1", 0},
		{"1.2.0", "v1.10.0", -1}, // sayısal karşılaştırma: 2 < 10 (sözlüksel değil)
		{"v2.0.0", "v1.9.9", 1},
		{"v1.1.1", "v1.1.1-beta", 0}, // önek/sonek soyulur
		{"v1.1", "v1.1.0", 0},        // eksik parça 0 sayılır
	}
	for _, c := range cases {
		if got := compareVersions(c.a, c.b); got != c.want {
			t.Errorf("compareVersions(%q,%q)=%d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestSumFor(t *testing.T) {
	sums := "abc123  kanije-linux-amd64\ndef456  kanije-windows-amd64.exe\n"

	got, err := sumFor(sums, "kanije-windows-amd64.exe")
	if err != nil || got != "def456" {
		t.Fatalf("sumFor=%q err=%v, want def456", got, err)
	}
	if _, err := sumFor(sums, "olmayan-dosya"); err == nil {
		t.Fatal("eksik varlık için hata beklenir")
	}
}

func TestAssetNameForPlatform(t *testing.T) {
	name := AssetName()
	switch runtime.GOOS {
	case "windows":
		if name != "kanije-windows-amd64.exe" {
			t.Errorf("windows asset adı yanlış: %q", name)
		}
	case "linux":
		if name == "" {
			t.Error("linux asset adı boş olmamalı")
		}
	}
}
