package app

import (
	"context"
	"errors"
	"testing"

	"github.com/kanije-kalesi/kanije/internal/storage"
)

func TestFlushPendingDeliversAllOnSuccess(t *testing.T) {
	msgs := []storage.PendingMessage{{ID: 1, Text: "a"}, {ID: 2, Text: "b"}, {ID: 3, Text: "c"}}

	var got []string
	sent := flushPending(context.Background(), msgs, func(_ context.Context, text string) error {
		got = append(got, text)
		return nil
	})

	if len(sent) != 3 || sent[0] != 1 || sent[2] != 3 {
		t.Fatalf("tüm mesajlar teslim edilmeli (FIFO id): got=%v", sent)
	}
	if len(got) != 3 || got[0] != "a" || got[2] != "c" {
		t.Fatalf("FIFO gönderim sırası korunmalı: got=%v", got)
	}
}

func TestFlushPendingStopsAtFirstFailureWithoutLoss(t *testing.T) {
	msgs := []storage.PendingMessage{{ID: 1, Text: "a"}, {ID: 2, Text: "b"}, {ID: 3, Text: "c"}}

	calls := 0
	sent := flushPending(context.Background(), msgs, func(_ context.Context, _ string) error {
		calls++
		if calls == 2 {
			return errors.New("ağ hatası")
		}
		return nil
	})

	// Yalnızca ilk mesaj teslim edilmiş sayılmalı; 2. başarısız, 3. hiç denenmemeli.
	// Böylece caller 2 ve 3'ü kuyrukta tutar (veri kaybı yok).
	if len(sent) != 1 || sent[0] != 1 {
		t.Fatalf("hata öncesi mesaj(lar) teslim edilmeli, sonrası değil: got=%v", sent)
	}
	if calls != 2 {
		t.Fatalf("ilk hatadan sonra gönderim denenmemeli: deneme=%d", calls)
	}
}

func TestFlushPendingEmptyDoesNotSend(t *testing.T) {
	sent := flushPending(context.Background(), nil, func(_ context.Context, _ string) error {
		t.Fatal("boş kuyrukta send çağrılmamalı")
		return nil
	})
	if sent != nil {
		t.Fatalf("boş kuyruk nil dönmeli: got=%v", sent)
	}
}
