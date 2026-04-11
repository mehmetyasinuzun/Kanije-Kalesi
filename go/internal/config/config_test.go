package config

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestLoadAppliesEnvOverrides(t *testing.T) {
	t.Setenv("KANIJE_BOT_TOKEN", "env-token")
	t.Setenv("KANIJE_CHAT_ID", "123456")
	t.Setenv("KANIJE_LOG_LEVEL", "debug")
	t.Setenv("KANIJE_DB_PATH", "./env.db")

	cfg, err := Load(filepath.Join(t.TempDir(), "missing.toml"))
	if err != nil {
		t.Fatalf("Load() hata verdi: %v", err)
	}

	if cfg.Telegram.BotToken != "env-token" {
		t.Fatalf("BotToken env override calismadi: %q", cfg.Telegram.BotToken)
	}
	if cfg.Telegram.ChatID != 123456 {
		t.Fatalf("ChatID env override calismadi: %d", cfg.Telegram.ChatID)
	}
	if cfg.Logging.Level != "debug" {
		t.Fatalf("Log level env override calismadi: %q", cfg.Logging.Level)
	}
	if cfg.Storage.DBPath != "./env.db" {
		t.Fatalf("DB path env override calismadi: %q", cfg.Storage.DBPath)
	}
}

func TestSetFieldValidation(t *testing.T) {
	cfg := Defaults()

	if err := cfg.SetField("logging.level", "warn"); err != nil {
		t.Fatalf("logging.level set edilemedi: %v", err)
	}
	if cfg.Logging.Level != "warn" {
		t.Fatalf("logging.level beklenen degere set edilmedi: %q", cfg.Logging.Level)
	}

	if err := cfg.SetField("telegram.chat_id", "999"); err != nil {
		t.Fatalf("telegram.chat_id set edilemedi: %v", err)
	}
	if cfg.Telegram.ChatID != 999 {
		t.Fatalf("telegram.chat_id beklenen degere set edilmedi: %d", cfg.Telegram.ChatID)
	}

	if err := cfg.SetField("telegram.chat_id", "abc"); err == nil {
		t.Fatalf("gecersiz chat_id icin hata bekleniyordu")
	}
	if err := cfg.SetField("logging.level", "trace"); err == nil {
		t.Fatalf("gecersiz log seviyesi icin hata bekleniyordu")
	}
	if err := cfg.SetField("unknown.field", "x"); err == nil {
		t.Fatalf("bilinmeyen alan icin hata bekleniyordu")
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")

	cfg := Defaults()
	cfg.SetFilePath(path)
	cfg.Telegram.BotToken = "token-1234"
	cfg.Telegram.ChatID = 42

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() hata verdi: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() hata verdi: %v", err)
	}

	if loaded.Telegram.BotToken != "token-1234" {
		t.Fatalf("round-trip token beklenenden farkli: %q", loaded.Telegram.BotToken)
	}
	if loaded.Telegram.ChatID != 42 {
		t.Fatalf("round-trip chat id beklenenden farkli: %d", loaded.Telegram.ChatID)
	}
}

func TestGetSafeJSONMasksBotToken(t *testing.T) {
	cfg := Defaults()
	cfg.Telegram.BotToken = "abcd1234wxyz"
	cfg.Telegram.ChatID = 77

	jsonStr := cfg.GetSafeJSON()

	var out map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &out); err != nil {
		t.Fatalf("GetSafeJSON cikti parse edilemedi: %v", err)
	}

	tel, ok := out["telegram"].(map[string]any)
	if !ok {
		t.Fatalf("telegram alani bulunamadi")
	}
	if got, _ := tel["bot_token"].(string); got != "abcd****wxyz" {
		t.Fatalf("maskelenmis token beklenenden farkli: %q", got)
	}
}
