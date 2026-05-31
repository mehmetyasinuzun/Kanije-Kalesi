package telegram

import (
	"strings"
	"testing"

	"github.com/kanije-kalesi/kanije/internal/defender"
)

func TestFormatDefenderUnavailable(t *testing.T) {
	if s := FormatDefender(defender.Status{Available: false}); !strings.Contains(s, "okunamadı") {
		t.Errorf("kullanılamaz mesajı bekleniyordu: %s", s)
	}
}

func TestFormatDefenderActive(t *testing.T) {
	s := FormatDefender(defender.Status{
		Available:          true,
		RealTimeProtection: true,
		AntivirusEnabled:   true,
		LastQuickScan:      "01.01.2026 10:00",
		SignatureVersion:   "1.400.123.0",
		RecentThreats:      []defender.Threat{{Name: "Trojan:Win32/Test", Time: "dün"}},
	})
	if !strings.Contains(s, "açık") {
		t.Errorf("gerçek-zamanlı koruma 'açık' gösterilmeli: %s", s)
	}
	if !strings.Contains(s, "Trojan:Win32/Test") {
		t.Errorf("tespit edilen tehdit gösterilmeli: %s", s)
	}
}

func TestFormatDefenderClean(t *testing.T) {
	s := FormatDefender(defender.Status{Available: true, RealTimeProtection: false})
	if !strings.Contains(s, "kapalı") {
		t.Errorf("kapalı koruma gösterilmeli: %s", s)
	}
	if !strings.Contains(s, "tehdit yok") {
		t.Errorf("temiz geçmiş gösterilmeli: %s", s)
	}
}
