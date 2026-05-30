package geoip

import "testing"

func TestFlagEmoji(t *testing.T) {
	if got := FlagEmoji("RU"); got != "🇷🇺" {
		t.Errorf("RU bayrağı yanlış: %q", got)
	}
	if got := FlagEmoji("us"); got != "🇺🇸" {
		t.Errorf("küçük harf de çalışmalı: %q", got)
	}
	if FlagEmoji("X") != "" {
		t.Error("tek harf boş dönmeli")
	}
	if FlagEmoji("12") != "" {
		t.Error("rakam boş dönmeli")
	}
}

func TestIsPublic(t *testing.T) {
	public := []string{"8.8.8.8", "1.1.1.1", "203.0.113.7"}
	private := []string{"192.168.1.1", "10.0.0.1", "127.0.0.1", "::1", "169.254.1.1", "0.0.0.0", "", "notanip"}

	for _, ip := range public {
		if !IsPublic(ip) {
			t.Errorf("%s public olmalı", ip)
		}
	}
	for _, ip := range private {
		if IsPublic(ip) {
			t.Errorf("%s public OLMAMALI", ip)
		}
	}
}

func TestSummary(t *testing.T) {
	i := &Info{Country: "Russia", City: "Moscow", Flag: "🇷🇺"}
	if got := i.Summary(); got != "🇷🇺 Moscow, Russia" {
		t.Errorf("summary = %q", got)
	}
	only := &Info{Country: "Germany"}
	if got := only.Summary(); got != "Germany" {
		t.Errorf("yalnız ülke = %q", got)
	}
}
