package telegram

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// The undo window is set to time.Hour in these tests so the real timer never
// fires during the test; state transitions are driven by calling confirm/fire/
// abort/expire directly. This keeps every assertion deterministic.

func TestDangerConfirmThenFireRuns(t *testing.T) {
	m := newDangerManager()
	var ran atomic.Bool
	op := m.add(100, "kapat", "Sistemi kapat", time.Hour, func(context.Context) { ran.Store(true) })

	if _, ok := m.confirm(op.id); !ok {
		t.Fatal("onay başarılı olmalı")
	}
	m.fire(op.id)

	if !ran.Load() {
		t.Fatal("geri-alma penceresi dolunca işlem çalışmalı")
	}
	if _, ok := m.abort(op.id); ok {
		t.Fatal("çalışan işlem map'te kalmamalı")
	}
}

func TestDangerAbortBeforeFireSkipsRun(t *testing.T) {
	m := newDangerManager()
	var ran atomic.Bool
	op := m.add(100, "kapat", "Sistemi kapat", time.Hour, func(context.Context) { ran.Store(true) })

	m.confirm(op.id)
	if _, ok := m.abort(op.id); !ok {
		t.Fatal("geri-alma penceresinde abort başarılı olmalı")
	}
	m.fire(op.id) // op artık yok

	if ran.Load() {
		t.Fatal("geri alınan işlem ASLA çalışmamalı")
	}
}

func TestDangerExpireDropsUnconfirmed(t *testing.T) {
	m := newDangerManager()
	var expired atomic.Bool
	m.onExpire = func(*dangerOp) { expired.Store(true) }
	op := m.add(100, "kapat", "Sistemi kapat", time.Hour, func(context.Context) {})

	m.expire(op.id)

	if !expired.Load() {
		t.Fatal("onaylanmamış op süresi dolunca onExpire çağrılmalı")
	}
	if _, ok := m.confirm(op.id); ok {
		t.Fatal("süresi dolmuş op onaylanamamalı")
	}
}

func TestDangerExpireAfterConfirmIsNoop(t *testing.T) {
	m := newDangerManager()
	var ran, expired atomic.Bool
	m.onExpire = func(*dangerOp) { expired.Store(true) }
	op := m.add(100, "kapat", "Sistemi kapat", time.Hour, func(context.Context) { ran.Store(true) })

	m.confirm(op.id)
	m.expire(op.id) // onaylandıktan sonra confirm-TTL süresi dolsa bile etkisiz olmalı

	if expired.Load() {
		t.Fatal("onaylanmış op için onExpire çağrılmamalı")
	}
	m.fire(op.id)
	if !ran.Load() {
		t.Fatal("onaylı op hâlâ çalışabilmeli")
	}
}

func TestDangerDoubleConfirmRejected(t *testing.T) {
	m := newDangerManager()
	op := m.add(100, "kapat", "Sistemi kapat", time.Hour, func(context.Context) {})

	if _, ok := m.confirm(op.id); !ok {
		t.Fatal("ilk onay başarılı olmalı")
	}
	if _, ok := m.confirm(op.id); ok {
		t.Fatal("ikinci onay (çift tıklama) reddedilmeli")
	}
}

func TestDangerSecondOpReplacesSameChat(t *testing.T) {
	m := newDangerManager()
	op1 := m.add(100, "kapat", "Sistemi kapat", time.Hour, func(context.Context) {})
	op2 := m.add(100, "yeniden", "Sistemi yeniden başlat", time.Hour, func(context.Context) {})

	if _, ok := m.confirm(op1.id); ok {
		t.Fatal("aynı sohbetin eski op'u değiştirildiği için onaylanmamalı")
	}
	if _, ok := m.confirm(op2.id); !ok {
		t.Fatal("yeni op onaylanmalı")
	}
}

func TestDangerCancelChatScopesByChat(t *testing.T) {
	m := newDangerManager()
	opA := m.add(100, "kapat", "A", time.Hour, func(context.Context) {})
	opB := m.add(200, "kapat", "B", time.Hour, func(context.Context) {})

	got, ok := m.cancelChat(100)
	if !ok || got.id != opA.id {
		t.Fatalf("sohbet 100'ün op'u iptal edilmeli (got=%v ok=%v)", got, ok)
	}
	if _, ok := m.confirm(opB.id); !ok {
		t.Fatal("başka sohbetin op'u /iptal'den etkilenmemeli")
	}
}

func TestDangerCancelChatEmpty(t *testing.T) {
	m := newDangerManager()
	if _, ok := m.cancelChat(999); ok {
		t.Fatal("bekleyen op yokken cancelChat false dönmeli")
	}
}

// TestDangerConcurrent stresses the manager under the race detector: many
// goroutines add/confirm/fire/abort across a few chats simultaneously, mirroring
// Poll's per-update goroutines. It must not race or panic.
func TestDangerConcurrent(t *testing.T) {
	m := newDangerManager()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(chat int64) {
			defer wg.Done()
			op := m.add(chat, "kapat", "t", time.Hour, func(context.Context) {})
			if _, ok := m.confirm(op.id); ok {
				m.fire(op.id)
			} else {
				m.abort(op.id)
			}
			m.cancelChat(chat)
		}(int64(i % 5))
	}
	wg.Wait()
}

func TestOpSeqNum(t *testing.T) {
	if got := opSeqNum("op42"); got != 42 {
		t.Fatalf("op42 → 42 beklenir, %d alındı", got)
	}
	if got := opSeqNum("bozuk"); got != 0 {
		t.Fatalf("geçersiz ID → 0 beklenir, %d alındı", got)
	}
}

func TestUndoLabel(t *testing.T) {
	if got := undoLabel(undoWindowSystem); got != "15 saniye" {
		t.Fatalf("'15 saniye' beklenir, %q alındı", got)
	}
}
