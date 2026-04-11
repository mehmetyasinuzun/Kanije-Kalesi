package event

import (
	"testing"
	"time"
)

func TestDefaultSeverity(t *testing.T) {
	tests := []struct {
		typeIn Type
		want   Severity
	}{
		{typeIn: TypeLoginFailed, want: SeverityAlert},
		{typeIn: TypeUSBInserted, want: SeverityWarning},
		{typeIn: TypeSystemBoot, want: SeverityInfo},
		{typeIn: TypeSystemShutdown, want: SeverityWarning},
		{typeIn: TypeError, want: SeverityAlert},
		{typeIn: TypeLoginSuccess, want: SeverityInfo},
	}

	for _, tt := range tests {
		if got := DefaultSeverity(tt.typeIn); got != tt.want {
			t.Fatalf("DefaultSeverity(%q) = %v, want %v", tt.typeIn, got, tt.want)
		}
	}
}

func TestSeverityUnknownRepresentations(t *testing.T) {
	s := Severity(999)
	if got := s.String(); got != "bilinmiyor" {
		t.Fatalf("String() = %q, want bilinmiyor", got)
	}
	if got := s.Emoji(); got != "❓" {
		t.Fatalf("Emoji() = %q, want ❓", got)
	}
}

func TestNewSetsExpectedDefaults(t *testing.T) {
	ev := New(TypeLoginFailed, "UnitTest")

	if ev.Type != TypeLoginFailed {
		t.Fatalf("Type = %q, want %q", ev.Type, TypeLoginFailed)
	}
	if ev.Source != "UnitTest" {
		t.Fatalf("Source = %q, want UnitTest", ev.Source)
	}
	if ev.Severity != SeverityAlert {
		t.Fatalf("Severity = %v, want %v", ev.Severity, SeverityAlert)
	}
	if ev.Timestamp.IsZero() {
		t.Fatalf("Timestamp sifir olamaz")
	}
	if dt := time.Since(ev.Timestamp); dt < 0 || dt > 2*time.Second {
		t.Fatalf("Timestamp beklenen aralikta degil: dt=%v", dt)
	}
}
