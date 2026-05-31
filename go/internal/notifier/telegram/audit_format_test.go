package telegram

import (
	"strings"
	"testing"

	"github.com/kanije-kalesi/kanije/internal/fileaudit"
)

func TestFormatAccessEventsEmpty(t *testing.T) {
	if s := FormatAccessEvents(nil); !strings.Contains(s, "Kayıt yok") {
		t.Errorf("boş liste mesajı bekleniyordu: %s", s)
	}
}

func TestFormatAccessEvents(t *testing.T) {
	evs := []fileaudit.AccessEvent{{
		Time:    "01.01.2026 12:00",
		User:    "Yasin",
		Process: `C:\Windows\explorer.exe`,
		Object:  `C:\Users\Yasin\Music\gizli.txt`,
		Access:  "okuma/kopya",
	}}
	s := FormatAccessEvents(evs)
	if !strings.Contains(s, "Yasin") {
		t.Errorf("kullanıcı gösterilmeli: %s", s)
	}
	if !strings.Contains(s, "okuma/kopya") {
		t.Errorf("erişim türü gösterilmeli: %s", s)
	}
}

func TestShortPath(t *testing.T) {
	if got := shortPath("kısa"); got != "kısa" {
		t.Errorf("kısa yol değişmemeli: %q", got)
	}
	long := strings.Repeat("a", 80)
	got := shortPath(long)
	if len([]rune(got)) > 52 || !strings.HasPrefix(got, "…") {
		t.Errorf("uzun yol kısaltılmalı (…önekiyle): %q (len %d)", got, len([]rune(got)))
	}
}
