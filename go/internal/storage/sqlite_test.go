package storage

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/kanije-kalesi/kanije/internal/event"
)

func TestSQLiteEventAndPendingMessageFlow(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "kanije-test.db")

	s, err := NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("NewSQLite() hata verdi: %v", err)
	}
	defer s.Close()

	ev := event.New(event.TypeLoginFailed, "unit")
	ev.Username = "alice"
	ev.Timestamp = time.Now().Add(-48 * time.Hour)
	ev.Extra = map[string]string{"k": "v"}

	if err := s.SaveEvent(ctx, ev); err != nil {
		t.Fatalf("SaveEvent() hata verdi: %v", err)
	}

	count, err := s.CountEvents(ctx)
	if err != nil {
		t.Fatalf("CountEvents() hata verdi: %v", err)
	}
	if count != 1 {
		t.Fatalf("event sayisi beklenenden farkli: got=%d want=1", count)
	}

	recent, err := s.RecentEvents(ctx, 10)
	if err != nil {
		t.Fatalf("RecentEvents() hata verdi: %v", err)
	}
	if len(recent) != 1 {
		t.Fatalf("recent event sayisi beklenenden farkli: %d", len(recent))
	}
	if recent[0].Username != "alice" {
		t.Fatalf("recent username beklenenden farkli: %q", recent[0].Username)
	}

	filtered, err := s.QueryEvents(ctx, EventFilter{Type: event.TypeLoginFailed, Limit: 5})
	if err != nil {
		t.Fatalf("QueryEvents() hata verdi: %v", err)
	}
	if len(filtered) != 1 {
		t.Fatalf("filtered event sayisi beklenenden farkli: %d", len(filtered))
	}

	if err := s.SavePendingMessage(ctx, "offline message"); err != nil {
		t.Fatalf("SavePendingMessage() hata verdi: %v", err)
	}

	pending, err := s.PopPendingMessages(ctx)
	if err != nil {
		t.Fatalf("PopPendingMessages() hata verdi: %v", err)
	}
	if len(pending) != 1 {
		t.Fatalf("pending mesaj sayisi beklenenden farkli: %d", len(pending))
	}

	secondPop, err := s.PopPendingMessages(ctx)
	if err != nil {
		t.Fatalf("ikinci PopPendingMessages() hata verdi: %v", err)
	}
	if len(secondPop) != 0 {
		t.Fatalf("ikinci pop bos olmaliydi: %d", len(secondPop))
	}

	deleted, err := s.Prune(ctx, 1)
	if err != nil {
		t.Fatalf("Prune() hata verdi: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("Prune() silinen satir beklenenden farkli: got=%d want=1", deleted)
	}
}
