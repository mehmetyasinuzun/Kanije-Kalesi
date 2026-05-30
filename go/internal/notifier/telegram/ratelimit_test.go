package telegram

import "testing"

func TestCmdRateLimiterBlocksFlood(t *testing.T) {
	r := newCmdRateLimiter(5)

	allowed := 0
	for i := 0; i < 20; i++ {
		if r.allow(1) {
			allowed++
		}
	}
	if allowed != 5 {
		t.Fatalf("ilk pencerede en fazla 5 komut geçmeli: got=%d", allowed)
	}
}

func TestCmdRateLimiterIsolatesChats(t *testing.T) {
	r := newCmdRateLimiter(2)

	// chat 1'i tüket
	r.allow(1)
	r.allow(1)
	if r.allow(1) {
		t.Fatal("chat 1 limiti aşmamalı")
	}
	// chat 2 etkilenmemeli
	if !r.allow(2) {
		t.Fatal("farklı chat kendi limitine sahip olmalı")
	}
}

func TestCmdRateLimiterDefaultsOnInvalid(t *testing.T) {
	r := newCmdRateLimiter(0) // geçersiz → varsayılan 20
	allowed := 0
	for i := 0; i < 25; i++ {
		if r.allow(9) {
			allowed++
		}
	}
	if allowed != 20 {
		t.Fatalf("geçersiz değerde varsayılan 20 olmalı: got=%d", allowed)
	}
}
