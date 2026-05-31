//go:build !windows

// Package clipboard reads the system clipboard. Implemented only on Windows.
package clipboard

import "fmt"

// ReadText is unsupported on non-Windows platforms.
func ReadText() (string, error) {
	return "", fmt.Errorf("pano okuma yalnızca Windows'ta destekleniyor")
}
