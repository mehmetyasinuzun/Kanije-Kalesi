package schedule

import (
	"path/filepath"
	"testing"
	"time"
)

func newStore(t *testing.T) *Store {
	t.Helper()
	s, err := NewStore(filepath.Join(t.TempDir(), "s.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestAddListRemove(t *testing.T) {
	s := newStore(t)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	t1, err := s.Add("/foto", 1800, now)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if t1.ID != 1 {
		t.Errorf("ilk id = %d, beklenen 1", t1.ID)
	}
	if !t1.NextRun.Equal(now.Add(1800 * time.Second)) {
		t.Errorf("NextRun yanlış: %v", t1.NextRun)
	}

	if _, err := s.Add("/ekran", 3600, now); err != nil {
		t.Fatal(err)
	}
	if len(s.List()) != 2 {
		t.Fatalf("List = %d, beklenen 2", len(s.List()))
	}

	ok, _ := s.Remove(1)
	if !ok || len(s.List()) != 1 {
		t.Errorf("Remove(1) başarısız: ok=%v, kalan=%d", ok, len(s.List()))
	}
	if ok, _ := s.Remove(999); ok {
		t.Error("olmayan id silinince true döndü")
	}
}

func TestPersistAndSeq(t *testing.T) {
	path := filepath.Join(t.TempDir(), "s.json")
	s, _ := NewStore(path)
	now := time.Now()
	s.Add("/foto", 1800, now)

	s2, err := NewStore(path) // reopen
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if len(s2.List()) != 1 {
		t.Fatalf("kalıcılık: List = %d, beklenen 1", len(s2.List()))
	}
	// Seq must survive so IDs never collide after a restart.
	t2, _ := s2.Add("/ekran", 1800, now)
	if t2.ID != 2 {
		t.Errorf("yeniden açılışta id = %d, beklenen 2", t2.ID)
	}
}

func TestDue(t *testing.T) {
	s := newStore(t)
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s.Add("/foto", 60, base) // NextRun = base+60

	if d := s.Due(base.Add(30 * time.Second)); len(d) != 0 {
		t.Errorf("erken: due = %d, beklenen 0", len(d))
	}

	d := s.Due(base.Add(90 * time.Second))
	if len(d) != 1 {
		t.Fatalf("due = %d, beklenen 1", len(d))
	}
	if !s.List()[0].NextRun.After(base.Add(90 * time.Second)) {
		t.Error("Due sonrası NextRun ilerlemeli")
	}
}

func TestDueSkipsMissedSlots(t *testing.T) {
	s := newStore(t)
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	s.Add("/foto", 60, base) // every 60s

	// App was "down" 10 minutes → 10 slots passed, but it must fire ONCE, not 10×.
	d := s.Due(base.Add(10 * time.Minute))
	if len(d) != 1 {
		t.Errorf("kaçan slotlar tek tetiklenmeli, alınan: %d", len(d))
	}
}
