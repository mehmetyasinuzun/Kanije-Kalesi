//go:build !windows

package secret

import "fmt"

// DPAPI is Windows-only; on other platforms Protect falls back to storing the
// value as-is (guarded by 0600 file permissions).
func dpapiProtect(string) (string, error) {
	return "", fmt.Errorf("DPAPI yalnız Windows'ta kullanılabilir")
}

func dpapiUnprotect(string) (string, error) {
	return "", fmt.Errorf("bu yapılandırma Windows'ta şifrelenmiş; bu platformda çözülemez")
}

func dpapiSupported() bool { return false }
