package event

import (
	"strings"
	"testing"
)

func TestAllEventTypesHaveSeverity(t *testing.T) {
	all := []Type{
		TypeLoginSuccess, TypeLoginFailed, TypeLogoff,
		TypeScreenLock, TypeScreenUnlock,
		TypeSystemBoot, TypeSystemShutdown, TypeSystemSleep, TypeSystemWake,
		TypeUSBInserted, TypeUSBRemoved,
		TypeNetworkUp, TypeNetworkDown, TypeNetworkChanged,
		TypeHeartbeat, TypeError,
	}
	for _, tp := range all {
		s := DefaultSeverity(tp)
		if s.String() == "" {
			t.Fatalf("DefaultSeverity(%q) bos string dondu", tp)
		}
		if s.Emoji() == "" {
			t.Fatalf("DefaultSeverity(%q).Emoji() bos dondu", tp)
		}
	}
}

func TestEventStringContainsType(t *testing.T) {
	ev := New(TypeUSBInserted, "TestSource")
	ev.Username = "testuser"
	s := ev.String()

	if !strings.Contains(s, string(TypeUSBInserted)) {
		t.Fatalf("String() event type icermeli: %q", s)
	}
	if !strings.Contains(s, "TestSource") {
		t.Fatalf("String() kaynak icermeli: %q", s)
	}
}

func TestLogonTypeStrings(t *testing.T) {
	cases := []struct {
		lt   LogonType
		want string
	}{
		{LogonInteractive, "Etkileşimli"},
		{LogonNetwork, "Ağ"},
		{LogonBatch, "Batch"},
		{LogonService, "Servis"},
		{LogonUnlock, "Kilit Açma"},
		{LogonNetworkCleartext, "Ağ (Açık Metin)"},
		{LogonNewCredentials, "Yeni Kimlik"},
		{LogonRemoteInteractive, "Uzak Masaüstü"},
		{LogonCachedInteractive, "Önbellekli Etkileşimli"},
		{LogonType(99), "Tip-99"},
	}
	for _, tc := range cases {
		if got := tc.lt.String(); got != tc.want {
			t.Fatalf("LogonType(%d).String() = %q, want %q", int(tc.lt), got, tc.want)
		}
	}
}

func TestSeverityValues(t *testing.T) {
	if SeverityInfo >= SeverityWarning {
		t.Fatal("SeverityInfo < SeverityWarning olmali")
	}
	if SeverityWarning >= SeverityAlert {
		t.Fatal("SeverityWarning < SeverityAlert olmali")
	}
	if SeverityAlert >= SeverityCritical {
		t.Fatal("SeverityAlert < SeverityCritical olmali")
	}
}

func TestEventExtraFieldsPreserved(t *testing.T) {
	ev := New(TypeError, "unit")
	ev.Extra = map[string]string{
		"key1": "val1",
		"key2": "değer2",
	}

	if ev.Extra["key1"] != "val1" {
		t.Fatalf("Extra[key1] beklenenden farkli: %q", ev.Extra["key1"])
	}
	if ev.Extra["key2"] != "değer2" {
		t.Fatalf("Extra[key2] UTF-8 korunmali: %q", ev.Extra["key2"])
	}
}
