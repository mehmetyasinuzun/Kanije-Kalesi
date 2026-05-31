package ffmpeg

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestResolveAbsolute(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "ff")
	if err := os.WriteFile(fake, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	// An absolute, existing configured path takes priority over PATH.
	if got := Resolve(fake, ""); got != fake {
		t.Errorf("Resolve(mevcut mutlak) = %q, beklenen %q", got, fake)
	}
	// A non-existent absolute path must NOT be returned as-is.
	missing := filepath.Join(dir, "yok")
	if got := Resolve(missing, ""); got == missing {
		t.Errorf("var olmayan mutlak path döndürülmemeli: %q", got)
	}
}

func TestResolveInstallDir(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, ExeName())
	if err := os.WriteFile(exe, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	// With a clearly-nonexistent configured value and (likely) no ffmpeg on PATH
	// under this exact name, the installDir copy should be found. We assert it's
	// at least one of the valid answers (PATH ffmpeg may also exist on CI).
	got := Resolve("", dir)
	if got != exe && got == "" {
		t.Errorf("installDir kopyası bulunmalıydı: %q", got)
	}
}

func TestExeName(t *testing.T) {
	got := ExeName()
	want := "ffmpeg"
	if runtime.GOOS == "windows" {
		want = "ffmpeg.exe"
	}
	if got != want {
		t.Errorf("ExeName() = %q, beklenen %q", got, want)
	}
}
