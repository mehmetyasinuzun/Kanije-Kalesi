package config

import (
	"path/filepath"
	"testing"
)

func TestDefaultsHasSensibleValues(t *testing.T) {
	cfg := Defaults()

	if cfg.Telegram.SendTimeoutSec <= 0 {
		t.Fatalf("SendTimeoutSec pozitif olmali: %d", cfg.Telegram.SendTimeoutSec)
	}
	if cfg.Heartbeat.IntervalHours <= 0 {
		t.Fatalf("HeartbeatInterval pozitif olmali: %d", cfg.Heartbeat.IntervalHours)
	}
	if cfg.Storage.EventRetentionDays <= 0 {
		t.Fatalf("EventRetentionDays pozitif olmali: %d", cfg.Storage.EventRetentionDays)
	}
	if cfg.Security.MaxEventsPerMinute <= 0 {
		t.Fatalf("MaxEventsPerMinute pozitif olmali: %d", cfg.Security.MaxEventsPerMinute)
	}
	if len(cfg.Triggers) == 0 {
		t.Fatal("varsayilan tetikleyiciler bos olamaz")
	}
}

func TestIsConfigured(t *testing.T) {
	cfg := Defaults()

	if cfg.IsConfigured() {
		t.Fatal("bos config yapilandirilmis sayilmamali")
	}

	cfg.Telegram.BotToken = "tok"
	if cfg.IsConfigured() {
		t.Fatal("sadece token varken yapilandirilmis sayilmamali")
	}

	cfg.Telegram.ChatID = 99
	if !cfg.IsConfigured() {
		t.Fatal("token + chatID varken yapilandirilmis sayilmali")
	}
}

func TestSetTriggerPersistsAndReloads(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")

	cfg := Defaults()
	cfg.SetFilePath(path)
	cfg.Telegram.BotToken = "x"
	cfg.Telegram.ChatID = 1

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() hata verdi: %v", err)
	}

	trig := TriggerConfig{
		Enabled:           true,
		CaptureCamera:     true,
		CaptureScreenshot: false,
		MaxPhotosPerMin:   5,
	}
	if err := cfg.SetTrigger("usb_inserted", trig); err != nil {
		t.Fatalf("SetTrigger() hata verdi: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() hata verdi: %v", err)
	}

	got, ok := loaded.GetTrigger("usb_inserted")
	if !ok {
		t.Fatal("usb_inserted trigger bulunamadi")
	}
	if !got.CaptureCamera {
		t.Fatal("CaptureCamera true olmali")
	}
	if got.MaxPhotosPerMin != 5 {
		t.Fatalf("MaxPhotosPerMin=5 bekleniyor: got=%d", got.MaxPhotosPerMin)
	}
}

func TestSetFieldHeartbeatHours(t *testing.T) {
	cfg := Defaults()
	cfg.SetFilePath("") // in-memory mod

	if err := cfg.SetField("heartbeat.interval_hours", "12"); err != nil {
		t.Fatalf("heartbeat.interval_hours set edilemedi: %v", err)
	}
	if cfg.Heartbeat.IntervalHours != 12 {
		t.Fatalf("IntervalHours=12 bekleniyor: got=%d", cfg.Heartbeat.IntervalHours)
	}

	// Gecersiz deger
	if err := cfg.SetField("heartbeat.interval_hours", "0"); err == nil {
		t.Fatal("0 icin hata bekleniyordu")
	}
	if err := cfg.SetField("heartbeat.interval_hours", "abc"); err == nil {
		t.Fatal("abc icin hata bekleniyordu")
	}
}

func TestSetFieldMaxEventsPerMinute(t *testing.T) {
	cfg := Defaults()
	cfg.SetFilePath("")

	if err := cfg.SetField("security.max_events_per_minute", "30"); err != nil {
		t.Fatalf("max_events_per_minute set edilemedi: %v", err)
	}
	if cfg.Security.MaxEventsPerMinute != 30 {
		t.Fatalf("MaxEventsPerMinute=30 bekleniyor: got=%d", cfg.Security.MaxEventsPerMinute)
	}
}

func TestSetFieldCameraIndex(t *testing.T) {
	cfg := Defaults()
	cfg.SetFilePath("")

	if err := cfg.SetField("camera.device_index", "1"); err != nil {
		t.Fatalf("camera.device_index set edilemedi: %v", err)
	}
	if cfg.Camera.DeviceIndex != 1 {
		t.Fatalf("DeviceIndex=1 bekleniyor: got=%d", cfg.Camera.DeviceIndex)
	}

	if err := cfg.SetField("camera.device_index", "-1"); err == nil {
		t.Fatal("negatif indeks icin hata bekleniyordu")
	}
}

func TestSetFieldCameraDeviceName(t *testing.T) {
	cfg := Defaults()
	cfg.SetFilePath("")

	if err := cfg.SetField("camera.device_name", "HD Webcam"); err != nil {
		t.Fatalf("camera.device_name set edilemedi: %v", err)
	}
	if cfg.Camera.DeviceName != "HD Webcam" {
		t.Fatalf("DeviceName beklenenle eslesmiyor: %q", cfg.Camera.DeviceName)
	}
}

func TestParseBoolVariants(t *testing.T) {
	trues := []string{"true", "1", "evet", "açık", "on", "TRUE", "Evet"}
	for _, s := range trues {
		if !parseBool(s) {
			t.Fatalf("parseBool(%q) false dondu, true bekleniyor", s)
		}
	}

	falses := []string{"false", "0", "hayir", "off", "False"}
	for _, s := range falses {
		if parseBool(s) {
			t.Fatalf("parseBool(%q) true dondu, false bekleniyor", s)
		}
	}
}

func TestLoadMissingFileFallsBackToDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.toml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("eksik dosya icin hata olmamali: %v", err)
	}
	if cfg.Telegram.SendTimeoutSec != 15 {
		t.Fatalf("varsayilan SendTimeoutSec=15 bekleniyor: got=%d", cfg.Telegram.SendTimeoutSec)
	}
}
