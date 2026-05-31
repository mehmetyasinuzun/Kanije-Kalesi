package capture

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveToDisk(t *testing.T) {
	dir := t.TempDir()

	p, err := SaveToDisk(dir, "foto", "jpg", []byte("hello"))
	if err != nil {
		t.Fatalf("SaveToDisk: %v", err)
	}
	if got := filepath.Dir(p); got != dir {
		t.Errorf("dizin = %q, beklenen %q", got, dir)
	}
	base := filepath.Base(p)
	if !strings.HasPrefix(base, "foto_") || !strings.HasSuffix(base, ".jpg") {
		t.Errorf("dosya adı beklenmedik: %q", base)
	}

	b, err := os.ReadFile(p)
	if err != nil || string(b) != "hello" {
		t.Errorf("içerik = %q (err %v), beklenen \"hello\"", b, err)
	}
}

func TestSaveToDiskCreatesDir(t *testing.T) {
	sub := filepath.Join(t.TempDir(), "yeni", "alt")
	if _, err := SaveToDisk(sub, "ekran", "jpg", []byte("x")); err != nil {
		t.Fatalf("eksik dizin oluşturulmalı: %v", err)
	}
	if _, err := os.Stat(sub); err != nil {
		t.Errorf("dizin oluşmadı: %v", err)
	}
}
