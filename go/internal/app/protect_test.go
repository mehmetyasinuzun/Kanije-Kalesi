package app

import (
	"strings"
	"testing"

	"github.com/kanije-kalesi/kanije/internal/event"
)

func TestUsbMatches(t *testing.T) {
	// Empty filter → any USB removal matches.
	if !usbMatches("", event.Event{}) {
		t.Error("boş filtre her USB'yi eşleştirmeli")
	}
	if !usbMatches("Kingston", event.Event{DeviceName: "Kingston"}) {
		t.Error("aygıt adı eşleşmeli")
	}
	if !usbMatches("USB DISK", event.Event{DeviceLabel: "usb disk"}) {
		t.Error("etiket büyük/küçük harf duyarsız eşleşmeli")
	}
	if usbMatches("Kingston", event.Event{DeviceName: "SanDisk"}) {
		t.Error("eşleşmeyen aygıt false dönmeli")
	}
}

func TestActionLabel(t *testing.T) {
	if got := actionLabel("lock"); got != "Kilitle" {
		t.Errorf("actionLabel(lock) = %q", got)
	}
	full := actionLabel("lock_alert_wipe")
	for _, want := range []string{"Kilitle", "Alarm", "Güvenli Sil"} {
		if !strings.Contains(full, want) {
			t.Errorf("actionLabel(lock_alert_wipe)=%q, %q içermeli", full, want)
		}
	}
}
