package telegram

import (
	"strings"
	"testing"
	"time"

	"github.com/kanije-kalesi/kanije/internal/sysinfo"
)

func TestFormatBatteryAbsent(t *testing.T) {
	s := FormatBattery(sysinfo.Battery{Present: false})
	if !strings.Contains(s, "algılanmadı") {
		t.Errorf("pil yok mesajı bekleniyordu: %s", s)
	}
}

func TestFormatBatteryCharging(t *testing.T) {
	s := FormatBattery(sysinfo.Battery{
		Present: true, Percent: 55, Charging: true,
		Remaining: 90 * time.Minute,
	})
	if !strings.Contains(s, "%55") {
		t.Errorf("yüzde gösterilmedi: %s", s)
	}
	if !strings.Contains(s, "Şarj") {
		t.Errorf("şarj durumu gösterilmedi: %s", s)
	}
}

func TestFormatBatteryDischarging(t *testing.T) {
	s := FormatBattery(sysinfo.Battery{Present: true, Percent: 12, Charging: false})
	if !strings.Contains(s, "%12") || !strings.Contains(s, "Pilden") {
		t.Errorf("boşalma durumu beklenmedik: %s", s)
	}
}
