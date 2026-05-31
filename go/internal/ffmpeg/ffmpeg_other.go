//go:build !windows

package ffmpeg

import "fmt"

// Download is Windows-only auto-provisioning. On Linux/macOS the user installs
// ffmpeg via the package manager (apt/dnf/brew).
func Download(_ string) (string, error) {
	return "", fmt.Errorf("otomatik ffmpeg indirme yalnızca Windows'ta; ffmpeg'i paket yöneticinizle kurun (ör. apt install ffmpeg / brew install ffmpeg)")
}
