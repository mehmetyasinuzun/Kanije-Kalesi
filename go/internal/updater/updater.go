// Package updater checks GitHub Releases for a newer version and (on supported
// platforms) installs it. The download is verified against the published
// SHA256SUMS.txt before it is applied.
package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ErrSelfInstallUnsupported is returned by SelfInstall on platforms where the
// agent cannot replace its own binary (e.g. a hardened systemd service running
// as an unprivileged user). The caller should fall back to a notification.
var ErrSelfInstallUnsupported = fmt.Errorf("bu platformda otomatik kurulum desteklenmiyor")

// Release describes a GitHub release relevant to the current platform.
type Release struct {
	Version  string // tag, e.g. "v1.1.2"
	AssetURL string // download URL of the platform binary
	SumsURL  string // download URL of SHA256SUMS.txt (may be empty)
	Notes    string
}

// Updater queries a GitHub repository's releases.
type Updater struct {
	owner   string
	repo    string
	current string // current version, e.g. "v1.1.1" or "dev"
	asset   string // platform asset name to look for
	http    *http.Client
	log     *slog.Logger
}

// New creates an Updater for the given repo and current version.
func New(owner, repo, currentVersion string, log *slog.Logger) *Updater {
	return &Updater{
		owner:   owner,
		repo:    repo,
		current: currentVersion,
		asset:   AssetName(),
		http:    &http.Client{Timeout: 60 * time.Second},
		log:     log,
	}
}

// AssetName returns the release asset name for the current platform/arch.
func AssetName() string {
	switch runtime.GOOS {
	case "windows":
		return "kanije-windows-amd64.exe"
	case "linux":
		switch runtime.GOARCH {
		case "arm64":
			return "kanije-linux-arm64"
		case "arm":
			return "kanije-linux-arm"
		default:
			return "kanije-linux-amd64"
		}
	default:
		return ""
	}
}

// LatestNewer fetches the latest release and reports whether it is newer than
// the current version. A "dev" build is always considered up to date so local
// builds never try to self-update.
func (u *Updater) LatestNewer(ctx context.Context) (*Release, bool, error) {
	rel, err := u.fetchLatest(ctx)
	if err != nil {
		return nil, false, err
	}
	if u.current == "dev" || u.current == "" {
		return rel, false, nil
	}
	return rel, compareVersions(rel.Version, u.current) > 0, nil
}

type ghRelease struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

func (u *Updater) fetchLatest(ctx context.Context) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", u.owner, u.repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "kanije-updater")

	resp, err := u.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("release sorgusu başarısız: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API durumu %d", resp.StatusCode)
	}

	var gh ghRelease
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&gh); err != nil {
		return nil, fmt.Errorf("release yanıtı çözümlenemedi: %w", err)
	}

	rel := &Release{Version: gh.TagName, Notes: gh.Body}
	for _, a := range gh.Assets {
		switch a.Name {
		case u.asset:
			rel.AssetURL = a.URL
		case "SHA256SUMS.txt":
			rel.SumsURL = a.URL
		}
	}
	if rel.AssetURL == "" {
		return nil, fmt.Errorf("bu platform için release varlığı bulunamadı: %s", u.asset)
	}
	return rel, nil
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

// sumFor extracts the hex digest for assetName from "sha256sum" style output.
func sumFor(sums, assetName string) (string, error) {
	for _, line := range strings.Split(sums, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == assetName {
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("%s için sağlama bulunamadı", assetName)
}

// compareVersions compares two "vX.Y.Z" strings numerically.
// Returns -1 if a<b, 0 if equal, 1 if a>b. Non-numeric or missing parts are 0.
func compareVersions(a, b string) int {
	pa := splitVersion(a)
	pb := splitVersion(b)
	for i := 0; i < len(pa) || i < len(pb); i++ {
		var x, y int
		if i < len(pa) {
			x = pa[i]
		}
		if i < len(pb) {
			y = pb[i]
		}
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
	}
	return 0
}

func splitVersion(v string) []int {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	// Drop any pre-release/build suffix after '-' or '+'.
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		v = v[:i]
	}
	parts := strings.Split(v, ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, _ := strconv.Atoi(p)
		out = append(out, n)
	}
	return out
}
