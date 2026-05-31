package access

import (
	"context"
	"testing"
)

func TestTransferRoot(t *testing.T) {
	ctx := context.Background()
	m := newTestMgr(t, 100) // owner = 100

	// Build a small tree: 100 → 200 → 300.
	if _, err := m.Invite(ctx, 100, 200, "A", AllCapabilities()); err != nil {
		t.Fatalf("invite 200: %v", err)
	}
	if _, err := m.Invite(ctx, 200, 300, "B", AllCapabilities()); err != nil {
		t.Fatalf("invite 300: %v", err)
	}

	if err := m.TransferRoot(ctx, 999); err != nil {
		t.Fatalf("transfer: %v", err)
	}

	// New owner is the sole root with every capability.
	if got := m.RootID(); got != 999 {
		t.Errorf("RootID = %d, beklenen 999", got)
	}
	u, ok := m.Get(999)
	if !ok {
		t.Fatal("yeni sahip ağaçta yok")
	}
	if !u.IsRoot() {
		t.Error("yeni sahip root değil")
	}
	if len(u.Caps) != len(AllCaps) {
		t.Errorf("yeni sahip %d yetkiye sahip, beklenen %d (tümü)", len(u.Caps), len(AllCaps))
	}

	// Old owner and all sub-users lost access.
	for _, old := range []int64{100, 200, 300} {
		if m.Known(old) {
			t.Errorf("eski kullanıcı %d hâlâ ağaçta", old)
		}
	}

	// Guard: zero is rejected.
	if err := m.TransferRoot(ctx, 0); err == nil {
		t.Error("geçersiz (0) sahip ID hata vermeliydi")
	}
}
