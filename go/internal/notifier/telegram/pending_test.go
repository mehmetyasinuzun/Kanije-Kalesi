package telegram

import (
	"sync"
	"testing"
)

func TestPendingActionTakeMatchingKind(t *testing.T) {
	b := &Bot{}
	b.armPendingAction("yeniden")

	if _, ok := b.takePendingAction("kapat"); ok {
		t.Fatal("yanlış kind tüketilmemeli")
	}

	state, ok := b.takePendingAction("yeniden")
	if !ok {
		t.Fatal("doğru kind tüketilmeli")
	}
	state.cancel() // süre-dolumu goroutine'ini sonlandır

	if _, ok := b.takePendingAction("yeniden"); ok {
		t.Fatal("zaten tüketilmiş işlem ikinci kez alınmamalı")
	}
}

func TestPendingActionClearConsumes(t *testing.T) {
	b := &Bot{}
	b.armPendingAction("kapat")

	state, ok := b.clearPendingAction()
	if !ok {
		t.Fatal("clear bekleyen işlemi tüketmeli")
	}
	state.cancel()

	if _, ok := b.clearPendingAction(); ok {
		t.Fatal("ikinci clear başarısız olmalı")
	}
}

func TestPendingActionReplaceCancelsOld(t *testing.T) {
	b := &Bot{}
	b.armPendingAction("yeniden")
	b.armPendingAction("kapat") // önceki 'yeniden' değiştirilmeli

	if _, ok := b.takePendingAction("yeniden"); ok {
		t.Fatal("eski işlem değiştirildiği için 'yeniden' alınmamalı")
	}
	state, ok := b.takePendingAction("kapat")
	if !ok {
		t.Fatal("yeni işlem 'kapat' alınabilmeli")
	}
	state.cancel()
}

// TestPendingActionConcurrent stresses the mutex-guarded state under the race
// detector: many goroutines arm and take simultaneously, mirroring Poll's
// per-update goroutines. It must not race or panic.
func TestPendingActionConcurrent(t *testing.T) {
	b := &Bot{}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.armPendingAction("yeniden")
			if s, ok := b.takePendingAction("yeniden"); ok {
				s.cancel()
			}
		}()
	}
	wg.Wait()
}
