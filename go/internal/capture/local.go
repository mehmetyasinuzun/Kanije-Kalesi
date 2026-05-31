package capture

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SaveToDisk writes capture bytes to dir as a timestamped file and returns its
// path, creating the directory if needed. prefix labels the kind ("foto",
// "ekran") and ext is the extension without a dot ("jpg"). This backs the
// camera/screenshot SaveLocal option; the caller decides whether to keep the
// file or delete it after a successful send.
func SaveToDisk(dir, prefix, ext string, data []byte) (string, error) {
	if dir == "" {
		dir = "."
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("dizin oluşturulamadı: %w", err)
	}
	name := fmt.Sprintf("%s_%s.%s", prefix, time.Now().Format("20060102_150405"), ext)
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("dosya yazılamadı: %w", err)
	}
	return path, nil
}
