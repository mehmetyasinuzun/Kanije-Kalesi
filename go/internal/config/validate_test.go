package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateClampsUnsafeValues(t *testing.T) {
	c := Defaults()
	c.Security.MaxEventsPerMinute = 0
	c.Heartbeat.IntervalHours = -5
	c.Camera.JPEGQuality = 999
	c.Screenshot.JPEGQuality = 0
	c.Network.CheckPort = 0
	c.Network.CheckHost = ""
	c.Logging.Level = "verbose"
	c.Storage.MaxRecentEvents = 0
	c.Storage.EventRetentionDays = -1

	c.validate()

	if c.Security.MaxEventsPerMinute < 1 {
		t.Errorf("MaxEventsPerMinute clamp edilmeli: %d", c.Security.MaxEventsPerMinute)
	}
	if c.Heartbeat.IntervalHours < 1 {
		t.Errorf("IntervalHours clamp edilmeli: %d", c.Heartbeat.IntervalHours)
	}
	if c.Camera.JPEGQuality < 1 || c.Camera.JPEGQuality > 100 {
		t.Errorf("Camera.JPEGQuality 1..100 olmalı: %d", c.Camera.JPEGQuality)
	}
	if c.Screenshot.JPEGQuality < 1 || c.Screenshot.JPEGQuality > 100 {
		t.Errorf("Screenshot.JPEGQuality 1..100 olmalı: %d", c.Screenshot.JPEGQuality)
	}
	if c.Network.CheckPort != 443 {
		t.Errorf("CheckPort default'a düşmeli: %d", c.Network.CheckPort)
	}
	if c.Network.CheckHost == "" {
		t.Error("CheckHost boş kalmamalı")
	}
	if c.Logging.Level != "info" {
		t.Errorf("geçersiz log seviyesi info'ya düşmeli: %q", c.Logging.Level)
	}
	if c.Storage.MaxRecentEvents < 1 {
		t.Errorf("MaxRecentEvents clamp edilmeli: %d", c.Storage.MaxRecentEvents)
	}
	if c.Storage.EventRetentionDays < 1 {
		t.Errorf("EventRetentionDays clamp edilmeli: %d", c.Storage.EventRetentionDays)
	}
}

func TestDefaultsAreSelfConsistentAfterValidate(t *testing.T) {
	c := Defaults()
	before := c.Storage.MaxRecentEvents
	c.validate()
	// Varsayılanlar zaten geçerli olmalı — validate onları değiştirmemeli.
	if c.Storage.MaxRecentEvents != before {
		t.Errorf("varsayılan MaxRecentEvents validate ile değişmemeli: %d -> %d",
			before, c.Storage.MaxRecentEvents)
	}
}

func TestGetBoolAndToggleRoundTrip(t *testing.T) {
	c := Defaults()
	c.SetFilePath("") // bellek-içi: persist no-op

	before, err := c.GetBool("heartbeat.enabled")
	if err != nil {
		t.Fatalf("GetBool hata: %v", err)
	}
	flip := "false"
	if !before {
		flip = "true"
	}
	if err := c.SetField("heartbeat.enabled", flip); err != nil {
		t.Fatalf("SetField hata: %v", err)
	}
	after, _ := c.GetBool("heartbeat.enabled")
	if after == before {
		t.Fatal("toggle değeri gerçekten değiştirmeli")
	}
}

func TestGetBoolUnknownKeyErrors(t *testing.T) {
	c := Defaults()
	if _, err := c.GetBool("olmayan.anahtar"); err == nil {
		t.Fatal("bilinmeyen anahtar hata dönmeli")
	}
}

func TestThreadSafeGetters(t *testing.T) {
	c := Defaults()
	c.Telegram.ChatID = 42
	c.Telegram.AllowedChatIDs = []int64{7, 8}
	c.Heartbeat.IntervalHours = 3

	if c.ChatID() != 42 {
		t.Errorf("ChatID() = %d, want 42", c.ChatID())
	}
	if c.HeartbeatInterval() != 3*time.Hour {
		t.Errorf("HeartbeatInterval() = %v, want 3h", c.HeartbeatInterval())
	}
	ids := c.AllowedChatIDs()
	if len(ids) != 2 {
		t.Fatalf("AllowedChatIDs() len = %d, want 2", len(ids))
	}
	// Dönen slice bir kopya olmalı — değiştirmek config'i etkilememeli.
	ids[0] = 999
	if c.AllowedChatIDs()[0] == 999 {
		t.Error("AllowedChatIDs() kopya dönmeli, iç durum korunmalı")
	}
}

func TestNetworkTriggersMigratedOnLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "old.toml")
	// network_up/down içermeyen eski bir config dosyası.
	if err := os.WriteFile(path, []byte("[triggers.login_failed]\nenabled = true\n"), 0o600); err != nil {
		t.Fatalf("dosya yazma hatası: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load hata: %v", err)
	}
	if _, ok := cfg.GetTrigger("network_up"); !ok {
		t.Error("network_up trigger Load'da migrate edilmeli")
	}
	if _, ok := cfg.GetTrigger("network_down"); !ok {
		t.Error("network_down trigger Load'da migrate edilmeli")
	}
}
