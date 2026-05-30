package webhook

import (
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/kanije-kalesi/kanije/internal/event"
)

func TestDiscordPayload(t *testing.T) {
	ev := event.New(event.TypeLoginFailed, "test")
	ev.Hostname = "PC1"
	ev.Username = "alice"
	ev.SourceIP = "1.2.3.4"

	var m discordMessage
	if err := json.Unmarshal(buildPayload("discord", ev), &m); err != nil {
		t.Fatalf("discord payload geçersiz JSON: %v", err)
	}
	if len(m.Embeds) != 1 {
		t.Fatalf("1 embed bekleniyor: %d", len(m.Embeds))
	}
	if !strings.Contains(m.Embeds[0].Title, "Başarısız") {
		t.Errorf("başlık etiketi yanlış: %q", m.Embeds[0].Title)
	}
	if m.Embeds[0].Color != 0xFF6B00 {
		t.Errorf("alert rengi yanlış: %x", m.Embeds[0].Color)
	}
}

func TestJSONPayload(t *testing.T) {
	ev := event.New(event.TypeUSBInserted, "test")
	ev.Hostname = "PC1"

	var m map[string]any
	if err := json.Unmarshal(buildPayload("json", ev), &m); err != nil {
		t.Fatalf("json payload geçersiz: %v", err)
	}
	if m["type"] != "usb_inserted" {
		t.Errorf("type yanlış: %v", m["type"])
	}
	if m["hostname"] != "PC1" {
		t.Errorf("hostname yanlış: %v", m["hostname"])
	}
}

func TestNewDropsInvalidTargets(t *testing.T) {
	s := New([]Target{
		{Name: "good", URL: "https://example.com/hook", Format: "discord"},
		{Name: "bos", URL: "", Format: "json"},
		{Name: "gecersiz", URL: "not-a-url", Format: "json"},
	}, slog.Default())

	if !s.Enabled() {
		t.Fatal("geçerli hedef olmalı")
	}
	if len(s.targets) != 1 {
		t.Fatalf("yalnız 1 geçerli hedef kalmalı: %d", len(s.targets))
	}
}
