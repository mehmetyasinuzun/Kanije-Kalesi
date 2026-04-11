package storage

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/kanije-kalesi/kanije/internal/event"
)

func TestSQLiteQueryEventsWithFilters(t *testing.T) {
	ctx := context.Background()
	s, err := NewSQLite(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("NewSQLite() hata verdi: %v", err)
	}
	defer s.Close()

	// 3 farklı tip event ekle
	types := []event.Type{event.TypeLoginFailed, event.TypeUSBInserted, event.TypeLoginFailed}
	for i, tp := range types {
		ev := event.New(tp, "unit")
		ev.Username = "user" + string(rune('a'+i))
		if err := s.SaveEvent(ctx, ev); err != nil {
			t.Fatalf("SaveEvent() hata verdi: %v", err)
		}
	}

	// Tipe göre filtrele
	filtered, err := s.QueryEvents(ctx, EventFilter{Type: event.TypeLoginFailed, Limit: 10})
	if err != nil {
		t.Fatalf("QueryEvents() hata verdi: %v", err)
	}
	if len(filtered) != 2 {
		t.Fatalf("TypeLoginFailed için 2 event bekleniyor: got=%d", len(filtered))
	}

	// Limit testi
	limited, err := s.QueryEvents(ctx, EventFilter{Limit: 1})
	if err != nil {
		t.Fatalf("QueryEvents(limit=1) hata verdi: %v", err)
	}
	if len(limited) != 1 {
		t.Fatalf("limit=1 için 1 event bekleniyor: got=%d", len(limited))
	}
}

func TestSQLitePruneKeepsRecentEvents(t *testing.T) {
	ctx := context.Background()
	s, err := NewSQLite(filepath.Join(t.TempDir(), "prune.db"))
	if err != nil {
		t.Fatalf("NewSQLite() hata verdi: %v", err)
	}
	defer s.Close()

	// Eski (silinmeli)
	old := event.New(event.TypeSystemBoot, "unit")
	old.Timestamp = time.Now().Add(-10 * 24 * time.Hour) // 10 gün önce
	if err := s.SaveEvent(ctx, old); err != nil {
		t.Fatalf("SaveEvent(eski) hata verdi: %v", err)
	}

	// Yeni (korunmalı)
	recent := event.New(event.TypeSystemBoot, "unit")
	recent.Timestamp = time.Now()
	if err := s.SaveEvent(ctx, recent); err != nil {
		t.Fatalf("SaveEvent(yeni) hata verdi: %v", err)
	}

	deleted, err := s.Prune(ctx, 7) // 7 günden eskiyi sil
	if err != nil {
		t.Fatalf("Prune() hata verdi: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("Prune() 1 event silmeli: got=%d", deleted)
	}

	count, err := s.CountEvents(ctx)
	if err != nil {
		t.Fatalf("CountEvents() hata verdi: %v", err)
	}
	if count != 1 {
		t.Fatalf("Prune sonrasi 1 event kalmali: got=%d", count)
	}
}

func TestSQLiteMultiplePendingMessages(t *testing.T) {
	ctx := context.Background()
	s, err := NewSQLite(filepath.Join(t.TempDir(), "pending.db"))
	if err != nil {
		t.Fatalf("NewSQLite() hata verdi: %v", err)
	}
	defer s.Close()

	messages := []string{"mesaj 1", "mesaj 2 — Türkçe: çşğüöı", "mesaj 3"}
	for _, m := range messages {
		if err := s.SavePendingMessage(ctx, m); err != nil {
			t.Fatalf("SavePendingMessage(%q) hata verdi: %v", m, err)
		}
	}

	popped, err := s.PopPendingMessages(ctx)
	if err != nil {
		t.Fatalf("PopPendingMessages() hata verdi: %v", err)
	}
	if len(popped) != len(messages) {
		t.Fatalf("3 mesaj bekleniyor: got=%d", len(popped))
	}

	// Türkçe mesaj korunmalı
	found := false
	for _, p := range popped {
		if p.Text == messages[1] {
			found = true
		}
	}
	if !found {
		t.Fatal("Turkce mesaj korunmali")
	}

	// İkinci pop boş olmalı
	second, err := s.PopPendingMessages(ctx)
	if err != nil {
		t.Fatalf("ikinci PopPendingMessages() hata verdi: %v", err)
	}
	if len(second) != 0 {
		t.Fatalf("ikinci pop bos olmali: got=%d", len(second))
	}
}

func TestSQLiteEventExtraFieldsRoundTrip(t *testing.T) {
	ctx := context.Background()
	s, err := NewSQLite(filepath.Join(t.TempDir(), "extra.db"))
	if err != nil {
		t.Fatalf("NewSQLite() hata verdi: %v", err)
	}
	defer s.Close()

	ev := event.New(event.TypeLoginFailed, "unit")
	ev.Username = "testuser"
	ev.Extra = map[string]string{
		"reason":  "wrong password",
		"turkish": "çşğüöı",
	}

	if err := s.SaveEvent(ctx, ev); err != nil {
		t.Fatalf("SaveEvent() hata verdi: %v", err)
	}

	recent, err := s.RecentEvents(ctx, 1)
	if err != nil {
		t.Fatalf("RecentEvents() hata verdi: %v", err)
	}
	if len(recent) != 1 {
		t.Fatalf("1 event bekleniyor: got=%d", len(recent))
	}

	got := recent[0]
	if got.Extra["reason"] != "wrong password" {
		t.Fatalf("Extra[reason] korunmali: %q", got.Extra["reason"])
	}
	if got.Extra["turkish"] != "çşğüöı" {
		t.Fatalf("Extra UTF-8 korunmali: %q", got.Extra["turkish"])
	}
}

func TestSQLiteQueryEventsSince(t *testing.T) {
	ctx := context.Background()
	s, err := NewSQLite(filepath.Join(t.TempDir(), "since.db"))
	if err != nil {
		t.Fatalf("NewSQLite() hata verdi: %v", err)
	}
	defer s.Close()

	cutoff := time.Now()

	old := event.New(event.TypeLoginFailed, "unit")
	old.Timestamp = cutoff.Add(-1 * time.Hour)
	s.SaveEvent(ctx, old)

	fresh := event.New(event.TypeLoginFailed, "unit")
	fresh.Timestamp = cutoff.Add(1 * time.Second)
	s.SaveEvent(ctx, fresh)

	results, err := s.QueryEvents(ctx, EventFilter{Since: cutoff, Limit: 10})
	if err != nil {
		t.Fatalf("QueryEvents(since) hata verdi: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("since filtresiyle 1 event bekleniyor: got=%d", len(results))
	}
}
