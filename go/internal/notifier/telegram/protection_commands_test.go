package telegram

import (
	"strings"
	"testing"
)

func TestParseAction(t *testing.T) {
	ok := map[string]string{
		"kilit": "lock", "lock": "lock",
		"alarm": "lock_alert", "alert": "lock_alert",
		"sil": "lock_alert_wipe", "wipe": "lock_alert_wipe", "imha": "lock_alert_wipe",
		"kilitmodu": "alert_lockdown", "lockdown": "alert_lockdown",
	}
	for in, want := range ok {
		got, valid := parseAction(in)
		if !valid || got != want {
			t.Errorf("parseAction(%q) = %q,%v; beklenen %q", in, got, valid, want)
		}
	}
	if _, valid := parseAction("saçma"); valid {
		t.Error("geçersiz aksiyon valid=true döndü")
	}
}

func TestActionTR(t *testing.T) {
	if got := actionTR("lock"); got != "Kilitle" {
		t.Errorf("actionTR(lock) = %q", got)
	}
	full := actionTR("lock_alert_wipe")
	for _, want := range []string{"Kilitle", "Alarm", "Güvenli Sil"} {
		if !strings.Contains(full, want) {
			t.Errorf("actionTR(lock_alert_wipe)=%q, %q içermeli", full, want)
		}
	}
}

func TestOnOffTR(t *testing.T) {
	if !strings.Contains(onOffTR(true), "açık") || !strings.Contains(onOffTR(false), "kapalı") {
		t.Error("onOffTR yanlış")
	}
}
