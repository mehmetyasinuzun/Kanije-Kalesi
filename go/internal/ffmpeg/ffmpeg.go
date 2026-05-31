// Package ffmpeg locates — and on Windows, auto-provisions — the ffmpeg binary
// that camera/audio capture depend on, so the agent works without the user
// installing anything manually.
package ffmpeg

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ExeName is the ffmpeg executable file name for this platform.
func ExeName() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	}
	return "ffmpeg"
}

// Resolve returns a usable ffmpeg path, or "" if none is found:
//  1. an explicit configured path (absolute & existing, or resolvable on PATH)
//  2. "ffmpeg" on the system PATH
//  3. a previously auto-provisioned copy in installDir
func Resolve(configured, installDir string) string {
	if configured != "" && configured != "ffmpeg" {
		if filepath.IsAbs(configured) {
			if fileExists(configured) {
				return configured
			}
		} else if p, err := exec.LookPath(configured); err == nil {
			return p
		}
	}
	if p, err := exec.LookPath("ffmpeg"); err == nil {
		return p
	}
	if installDir != "" {
		if local := filepath.Join(installDir, ExeName()); fileExists(local) {
			return local
		}
	}
	return ""
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}
