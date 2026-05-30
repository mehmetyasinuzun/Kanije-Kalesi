package event

import "testing"

// TestNewBusGuardsZeroValues ensures an unvalidated config (e.g. MaxPerMinute=0)
// cannot panic the bus via a divide-by-zero on the refill interval.
func TestNewBusGuardsZeroValues(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("NewBus sıfır değerlerle panik etti: %v", r)
		}
	}()

	bus := NewBus(BusConfig{BufferSize: 0, MaxPerMinute: 0, DedupWindow: 0})

	ev := New(TypeLoginFailed, "test")
	ev.Username = "x"
	// Publish de panik etmemeli ve en az bir event geçmeli (buffer >= 1).
	if !bus.Publish(ev) {
		t.Fatal("guard sonrası en az bir event yayınlanabilmeli")
	}
}
