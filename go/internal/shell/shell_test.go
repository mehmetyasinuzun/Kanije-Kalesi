package shell

import (
	"context"
	"strings"
	"testing"
)

func TestParseCd(t *testing.T) {
	cases := []struct {
		in     string
		target string
		isCd   bool
	}{
		{"cd C:\\foo", "C:\\foo", true},
		{"cd ..", "..", true},
		{"CD /tmp", "/tmp", true},
		{"cd", "", true},
		{"whoami", "", false},
		{"echo cd", "", false},
	}
	for _, c := range cases {
		target, isCd := parseCd(c.in)
		if isCd != c.isCd || target != c.target {
			t.Errorf("parseCd(%q) = (%q,%v), want (%q,%v)", c.in, target, isCd, c.target, c.isCd)
		}
	}
}

func TestTruncateOutput(t *testing.T) {
	if got := truncateOutput("kısa"); got != "kısa" {
		t.Errorf("kısa çıktı değişmemeli: %q", got)
	}
	long := strings.Repeat("a", maxOutput+500)
	got := truncateOutput(long)
	if len(got) <= maxOutput || !strings.Contains(got, "kısaltıldı") {
		t.Errorf("uzun çıktı kısaltılmalı + not içermeli")
	}
}

func TestRunEcho(t *testing.T) {
	r := New()
	out := r.Run(context.Background(), 1, "echo kanije123")
	if !strings.Contains(out, "kanije123") {
		t.Errorf("echo çıktısı 'kanije123' içermeli, got %q", out)
	}
}

func TestRunCdPersists(t *testing.T) {
	r := New()
	tmp := t.TempDir()
	out := r.Run(context.Background(), 7, "cd "+tmp)
	if !strings.Contains(out, "📂") {
		t.Fatalf("cd başarılı olmalı, got %q", out)
	}
	if cwd := r.Cwd(7); cwd != tmp {
		// filepath.Clean may normalize separators; compare loosely.
		if !strings.EqualFold(strings.TrimRight(cwd, "\\/"), strings.TrimRight(tmp, "\\/")) {
			t.Errorf("cwd kalıcı olmalı: got %q want %q", cwd, tmp)
		}
	}
}

func TestRunCdNonexistent(t *testing.T) {
	r := New()
	out := r.Run(context.Background(), 9, "cd /kanije-yok-boyle-bir-dizin-12345")
	if !strings.Contains(out, "bulunamadı") {
		t.Errorf("olmayan dizin reddedilmeli, got %q", out)
	}
}
