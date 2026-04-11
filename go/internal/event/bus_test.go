package event

import (
	"testing"
	"time"
)

func TestBusPublishAndReceive(t *testing.T) {
	cfg := BusConfig{
		BufferSize:   16,
		MaxPerMinute: 60,
		DedupWindow:  50 * time.Millisecond,
	}
	bus := NewBus(cfg)

	ev := New(TypeLoginFailed, "test")
	ev.Username = "bob"
	ev.SourceIP = "1.2.3.4"

	if !bus.Publish(ev) {
		t.Fatal("ilk event kabul edilmeliydi")
	}

	select {
	case got := <-bus.Events():
		if got.Type != TypeLoginFailed {
			t.Fatalf("event type beklenenden farkli: %q", got.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("event kanaldan alinamadi")
	}
}

func TestBusDedupDropsDuplicates(t *testing.T) {
	cfg := BusConfig{
		BufferSize:   16,
		MaxPerMinute: 60,
		DedupWindow:  200 * time.Millisecond,
	}
	bus := NewBus(cfg)

	ev := New(TypeUSBInserted, "test")
	ev.Username = "alice"
	ev.SourceIP = "10.0.0.1"

	if !bus.Publish(ev) {
		t.Fatal("birinci publish kabul edilmeliydi")
	}
	if bus.Publish(ev) {
		t.Fatal("dedup penceresi icindeki ikinci publish reddedilmeliydi")
	}

	stats := bus.Stats()
	if stats.Deduped != 1 {
		t.Fatalf("Deduped sayaci 1 olmali: got=%d", stats.Deduped)
	}
}

func TestBusRateLimitDropsExcess(t *testing.T) {
	cfg := BusConfig{
		BufferSize:   64,
		MaxPerMinute: 2,
		DedupWindow:  0, // dedup kapali — her event farkli kullanici
	}
	bus := NewBus(cfg)

	accepted := 0
	for i := 0; i < 5; i++ {
		ev := New(TypeLoginFailed, "test")
		ev.Username = string(rune('a' + i)) // farkli kullanici → dedup es gecilir
		ev.SourceIP = "10.0.0.1"
		if bus.Publish(ev) {
			accepted++
		}
	}

	if accepted > 2 {
		t.Fatalf("rate limit 2 oldugunda en fazla 2 event gecmeli: got=%d", accepted)
	}

	stats := bus.Stats()
	if stats.Dropped == 0 {
		t.Fatal("rate limit devreye girmeli ve Dropped > 0 olmali")
	}
}

func TestBusStatsCountsCorrectly(t *testing.T) {
	cfg := BusConfig{
		BufferSize:   16,
		MaxPerMinute: 60,
		DedupWindow:  10 * time.Millisecond,
	}
	bus := NewBus(cfg)

	ev1 := New(TypeLoginSuccess, "test")
	ev1.Username = "user1"
	bus.Publish(ev1)

	ev2 := New(TypeSystemBoot, "test")
	ev2.Username = "user2"
	bus.Publish(ev2)

	// dedup icin ayni event tekrar
	bus.Publish(ev1)

	stats := bus.Stats()
	if stats.Received != 2 {
		t.Fatalf("Received=2 bekleniyor: got=%d", stats.Received)
	}
	if stats.Deduped != 1 {
		t.Fatalf("Deduped=1 bekleniyor: got=%d", stats.Deduped)
	}
	if stats.Pending != 2 {
		t.Fatalf("Pending=2 bekleniyor: got=%d", stats.Pending)
	}
}

func TestBusBufferFullDropsEvent(t *testing.T) {
	cfg := BusConfig{
		BufferSize:   2,
		MaxPerMinute: 60,
		DedupWindow:  0,
	}
	bus := NewBus(cfg)

	// 2 farklı event dolduruyor, 3. düşmeli
	for i := 0; i < 3; i++ {
		ev := New(TypeUSBRemoved, "test")
		ev.Username = string(rune('a' + i))
		bus.Publish(ev)
	}

	stats := bus.Stats()
	if stats.Dropped == 0 && stats.Deduped == 0 {
		t.Fatal("buffer dolunca en az 1 event dusmeli")
	}
}

func TestNewBusDefaultConfig(t *testing.T) {
	cfg := DefaultBusConfig()
	if cfg.BufferSize <= 0 {
		t.Fatalf("BufferSize pozitif olmali: %d", cfg.BufferSize)
	}
	if cfg.MaxPerMinute <= 0 {
		t.Fatalf("MaxPerMinute pozitif olmali: %d", cfg.MaxPerMinute)
	}
	if cfg.DedupWindow <= 0 {
		t.Fatalf("DedupWindow pozitif olmali: %v", cfg.DedupWindow)
	}
}
