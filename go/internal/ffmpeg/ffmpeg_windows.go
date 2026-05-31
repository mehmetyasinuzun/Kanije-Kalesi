//go:build windows

package ffmpeg

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// downloadURL is a static, self-contained Windows ffmpeg build (BtbN, the build
// the FFmpeg project itself links to). One file, all codecs — no install needed.
const downloadURL = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"

// Download fetches a static ffmpeg.exe into installDir and returns its path. If
// one is already there it returns immediately. ~80 MB one-time download.
func Download(installDir string) (string, error) {
	dest := filepath.Join(installDir, "ffmpeg.exe")
	if fileExists(dest) {
		return dest, nil
	}
	if err := os.MkdirAll(installDir, 0o700); err != nil {
		return "", err
	}

	tmpZip := filepath.Join(os.TempDir(), "kanije-ffmpeg.zip")
	if err := downloadFile(downloadURL, tmpZip); err != nil {
		return "", fmt.Errorf("ffmpeg indirilemedi: %w", err)
	}
	defer os.Remove(tmpZip)

	if err := extractFFmpeg(tmpZip, dest); err != nil {
		return "", fmt.Errorf("ffmpeg çıkarılamadı: %w", err)
	}
	return dest, nil
}

func downloadFile(url, dest string) error {
	client := &http.Client{Timeout: 15 * time.Minute}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "kanije")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("indirme durumu %d", resp.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(dest)
		return err
	}
	return f.Close()
}

// extractFFmpeg pulls just bin/ffmpeg.exe out of the build zip into dest.
func extractFFmpeg(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		name := strings.ToLower(strings.ReplaceAll(f.Name, "\\", "/"))
		if !strings.HasSuffix(name, "/bin/ffmpeg.exe") && name != "ffmpeg.exe" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, rc); err != nil { //nolint
			out.Close()
			os.Remove(dest)
			return err
		}
		return out.Close()
	}
	return fmt.Errorf("zip içinde ffmpeg.exe bulunamadı")
}
