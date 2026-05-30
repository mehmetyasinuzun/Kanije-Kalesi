//go:build windows

package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

// CanSelfInstall reports whether the agent can replace its own binary in place.
func CanSelfInstall() bool { return true }

// SelfInstall downloads (and verifies) the new binary, then spawns a detached
// helper that waits for this process to exit, swaps the executable, and restarts
// the scheduled task. The caller must shut the app down shortly after this
// returns so the helper can replace the now-unlocked exe.
func (u *Updater) SelfInstall(ctx context.Context, rel *Release) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	newPath := exe + ".new"
	if err := u.downloadVerified(ctx, rel, newPath); err != nil {
		return err
	}

	scriptPath := filepath.Join(os.TempDir(), "kanije-swap.ps1")
	if err := os.WriteFile(scriptPath, []byte(buildSwapScript(os.Getpid(), newPath, exe)), 0o600); err != nil {
		os.Remove(newPath)
		return err
	}

	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	sysproc.Hide(cmd)
	if err := cmd.Start(); err != nil {
		os.Remove(newPath)
		os.Remove(scriptPath)
		return err
	}
	u.log.Info("güncelleme yardımcısı başlatıldı", "sürüm", rel.Version)
	return nil
}

// buildSwapScript waits for the agent (pid) to exit, replaces exePath with
// newPath, and restarts the scheduled task.
func buildSwapScript(pid int, newPath, exePath string) string {
	var b strings.Builder
	b.WriteString("$ErrorActionPreference='SilentlyContinue'\n")
	b.WriteString("$p=" + strconv.Itoa(pid) + "\n")
	b.WriteString("while (Get-Process -Id $p -ErrorAction SilentlyContinue) { Start-Sleep -Milliseconds 300 }\n")
	b.WriteString("Start-Sleep -Milliseconds 800\n")
	b.WriteString("Move-Item -Force " + psQuote(newPath) + " " + psQuote(exePath) + "\n")
	b.WriteString("Start-ScheduledTask -TaskName KanijeKalesi\n")
	return b.String()
}

// psQuote wraps s in a PowerShell single-quoted literal, escaping embedded quotes.
func psQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// ---- download + verify (only the Windows self-installer needs these) ----

// downloadVerified downloads the platform asset to destPath and verifies its
// SHA-256 against the release's SHA256SUMS.txt. If sums are unavailable the
// download still succeeds (best effort) but is logged.
func (u *Updater) downloadVerified(ctx context.Context, rel *Release, destPath string) error {
	got, err := u.download(ctx, rel.AssetURL, destPath)
	if err != nil {
		return fmt.Errorf("binary indirilemedi: %w", err)
	}

	if rel.SumsURL == "" {
		u.log.Warn("SHA256SUMS bulunamadı, doğrulama atlanıyor")
		return nil
	}

	want, err := u.expectedSum(ctx, rel.SumsURL)
	if err != nil {
		os.Remove(destPath)
		return fmt.Errorf("sağlama listesi alınamadı: %w", err)
	}
	if !strings.EqualFold(got, want) {
		os.Remove(destPath)
		return fmt.Errorf("SHA-256 uyuşmuyor — indirme bozuk veya kurcalanmış olabilir")
	}
	return nil
}

// download writes url to destPath and returns the file's SHA-256 hex digest.
func (u *Updater) download(ctx context.Context, url, destPath string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "kanije-updater")

	resp, err := u.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("indirme durumu %d", resp.StatusCode)
	}

	f, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(f, h), resp.Body); err != nil {
		f.Close()
		os.Remove(destPath)
		return "", err
	}
	if err := f.Close(); err != nil {
		os.Remove(destPath)
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// expectedSum fetches SHA256SUMS.txt and returns the hash for this platform asset.
func (u *Updater) expectedSum(ctx context.Context, sumsURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sumsURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "kanije-updater")
	resp, err := u.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	if err != nil {
		return "", err
	}
	return sumFor(string(body), u.asset)
}
