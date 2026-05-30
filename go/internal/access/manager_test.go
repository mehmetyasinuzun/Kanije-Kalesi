package access

import (
	"context"
	"io"
	"log/slog"
	"sort"
	"testing"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// fakeStore is an in-memory Persistence for tests.
type fakeStore struct {
	users map[int64]User
	audit []AuditEntry
}

func newFakeStore() *fakeStore { return &fakeStore{users: map[int64]User{}} }

func (f *fakeStore) LoadUsers(_ context.Context) ([]User, error) {
	out := make([]User, 0, len(f.users))
	for _, u := range f.users {
		out = append(out, u)
	}
	return out, nil
}
func (f *fakeStore) SaveUser(_ context.Context, u User) error { f.users[u.ChatID] = u; return nil }
func (f *fakeStore) DeleteUser(_ context.Context, id int64) error {
	delete(f.users, id)
	return nil
}
func (f *fakeStore) AppendAudit(_ context.Context, e AuditEntry) error {
	f.audit = append(f.audit, e)
	return nil
}
func (f *fakeStore) RecentAudit(_ context.Context, n int) ([]AuditEntry, error) {
	if n > len(f.audit) {
		n = len(f.audit)
	}
	return f.audit[len(f.audit)-n:], nil
}

func newTestMgr(t *testing.T, owner int64) *Manager {
	t.Helper()
	m := NewManager(newFakeStore(), discardLogger())
	if err := m.Bootstrap(context.Background(), owner, nil); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	return m
}

func caps(u User) []Capability { return u.Caps }

func sortedKeys(cs []Capability) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = string(c)
	}
	sort.Strings(out)
	return out
}

func mustGet(t *testing.T, m *Manager, id int64) User {
	t.Helper()
	u, ok := m.Get(id)
	if !ok {
		t.Fatalf("kullanıcı %d yok", id)
	}
	return u
}

func TestBootstrapSeedsRootWithAllCaps(t *testing.T) {
	m := newTestMgr(t, 1)
	root := mustGet(t, m, 1)
	if !root.IsRoot() {
		t.Fatal("owner kök olmalı")
	}
	if len(caps(root)) != len(AllCaps) {
		t.Fatalf("kök tüm yetkilere sahip olmalı: %v", caps(root))
	}
	if m.RootID() != 1 {
		t.Fatalf("RootID 1 olmalı, %d", m.RootID())
	}
}

func TestBootstrapMigratesAllowlist(t *testing.T) {
	m := NewManager(newFakeStore(), discardLogger())
	if err := m.Bootstrap(context.Background(), 1, []int64{2, 3, 1 /*dup owner*/}); err != nil {
		t.Fatal(err)
	}
	if !m.Known(2) || !m.Known(3) {
		t.Fatal("izinli kullanıcılar taşınmalı")
	}
	u2 := mustGet(t, m, 2)
	if u2.InvitedBy != 1 {
		t.Fatalf("taşınan kullanıcı köke bağlı olmalı, %d", u2.InvitedBy)
	}
}

func TestBootstrapHealsOwnerCaps(t *testing.T) {
	store := newFakeStore()
	// Simulate a DB written by an older version: owner predates newer caps.
	store.users[1] = User{ChatID: 1, Name: "Sahip", InvitedBy: 0, Caps: []Capability{CapStatus, CapEvents}}

	m := NewManager(store, discardLogger())
	if err := m.Bootstrap(context.Background(), 1, nil); err != nil {
		t.Fatal(err)
	}
	if !m.Can(1, CapTerminal) {
		t.Fatal("owner güncellemeden sonra yeni yetkileri (terminal) otomatik almalı")
	}
	if root := mustGet(t, m, 1); len(root.Caps) != len(AllCaps) {
		t.Fatalf("owner tüm yetkilere sahip olmalı, %d/%d", len(root.Caps), len(AllCaps))
	}
}

func TestInviteNoEscalation(t *testing.T) {
	m := newTestMgr(t, 1)
	// A(=1, root) invites B(=2) with a subset that includes invite+manage.
	if _, err := m.Invite(context.Background(), 1, 2, "B", []Capability{
		CapStatus, CapEvents, CapPhoto, CapInvite, CapManage,
	}); err != nil {
		t.Fatal(err)
	}
	// B tries to grant C a capability B does NOT have (shutdown) → dropped.
	c, err := m.Invite(context.Background(), 2, 3, "C", []Capability{CapStatus, CapShutdown})
	if err != nil {
		t.Fatal(err)
	}
	if hasCap(c.Caps, CapShutdown) {
		t.Fatalf("yükseltme olmamalı: C shutdown almamalı, %v", sortedKeys(c.Caps))
	}
	if !hasCap(c.Caps, CapStatus) {
		t.Fatalf("C status almalı, %v", sortedKeys(c.Caps))
	}
}

func TestInviteRequiresInviteCap(t *testing.T) {
	m := newTestMgr(t, 1)
	// B has NO invite capability.
	if _, err := m.Invite(context.Background(), 1, 2, "B", []Capability{CapStatus, CapManage}); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Invite(context.Background(), 2, 3, "C", []Capability{CapStatus}); err == nil {
		t.Fatal("davet yetkisi olmayan kişi davet edememeli")
	}
}

func TestVisibilityIsSubtreeScoped(t *testing.T) {
	m := newTestMgr(t, 1) // A = root
	mustInvite(t, m, 1, 2, "B", CapStatus, CapInvite, CapManage)
	mustInvite(t, m, 2, 3, "C", CapStatus)

	if !m.Visible(1, 2) || !m.Visible(1, 3) {
		t.Fatal("A hem B'yi hem C'yi görmeli")
	}
	if !m.Visible(2, 3) {
		t.Fatal("B, C'yi görmeli")
	}
	if m.Visible(2, 1) {
		t.Fatal("B, A'yı (üstünü) görmemeli")
	}
	if m.Visible(3, 2) {
		t.Fatal("C, B'yi görmemeli")
	}
	if got := len(m.Subtree(1)); got != 2 {
		t.Fatalf("A alt ağacı 2 olmalı, %d", got)
	}
	if got := len(m.Subtree(3)); got != 0 {
		t.Fatalf("C alt ağacı 0 olmalı, %d", got)
	}
}

func TestSetCapsClampsToEditorAndCascades(t *testing.T) {
	m := newTestMgr(t, 1)
	mustInvite(t, m, 1, 2, "B", CapStatus, CapEvents, CapPhoto, CapInvite, CapManage)
	mustInvite(t, m, 2, 3, "C", CapStatus, CapEvents, CapPhoto, CapInvite, CapManage)
	mustInvite(t, m, 3, 4, "D", CapStatus, CapEvents)

	// B reduces C to only status. C's child D must lose events too (cascade).
	if err := m.SetCaps(context.Background(), 2, 3, []Capability{CapStatus, CapShutdown}); err != nil {
		t.Fatal(err)
	}
	c := mustGet(t, m, 3)
	if hasCap(c.Caps, CapShutdown) || hasCap(c.Caps, CapEvents) {
		t.Fatalf("C yalnız status olmalı (shutdown sızmamalı), %v", sortedKeys(c.Caps))
	}
	d := mustGet(t, m, 4)
	if hasCap(d.Caps, CapEvents) {
		t.Fatalf("kademeli kırpma: D events kaybetmeli, %v", sortedKeys(d.Caps))
	}
}

func TestRemoveReparentsToRootByDefault(t *testing.T) {
	m := newTestMgr(t, 1)
	mustInvite(t, m, 1, 2, "A2", CapStatus, CapInvite, CapManage)
	mustInvite(t, m, 2, 3, "B", CapStatus, CapInvite, CapManage)
	mustInvite(t, m, 3, 4, "C", CapStatus)

	// root removes B(=3); C(=4) must survive, re-homed to root by default.
	if err := m.Remove(context.Background(), 1, 3, ReparentRoot); err != nil {
		t.Fatal(err)
	}
	if m.Known(3) {
		t.Fatal("B silinmeli")
	}
	if !m.Known(4) {
		t.Fatal("C silinmemeli (kademeli silme yok)")
	}
	if c := mustGet(t, m, 4); c.InvitedBy != 1 {
		t.Fatalf("C köke bağlanmalı (varsayılan), bağlı: %d", c.InvitedBy)
	}
}

func TestRemoveReparentsToGrandparent(t *testing.T) {
	m := newTestMgr(t, 1)
	mustInvite(t, m, 1, 2, "A2", CapStatus, CapInvite, CapManage)
	mustInvite(t, m, 2, 3, "B", CapStatus, CapInvite, CapManage)
	mustInvite(t, m, 3, 4, "C", CapStatus)

	// root removes B(=3) with grandparent mode → C goes to B's parent A2(=2).
	if err := m.Remove(context.Background(), 1, 3, ReparentGrandparent); err != nil {
		t.Fatal(err)
	}
	if c := mustGet(t, m, 4); c.InvitedBy != 2 {
		t.Fatalf("C bir üste (A2=2) bağlanmalı, bağlı: %d", c.InvitedBy)
	}
}

func TestRemoveGuards(t *testing.T) {
	m := newTestMgr(t, 1)
	mustInvite(t, m, 1, 2, "B", CapStatus, CapInvite, CapManage) // can invite + manage
	mustInvite(t, m, 1, 3, "X", CapStatus)                       // sibling, no manage
	mustInvite(t, m, 2, 4, "Bchild", CapStatus)                  // in B's subtree

	if err := m.Remove(context.Background(), 1, 1, ReparentRoot); err == nil {
		t.Fatal("kök çıkarılamamalı")
	}
	if err := m.Remove(context.Background(), 2, 2, ReparentRoot); err == nil {
		t.Fatal("kendini çıkaramamalı")
	}
	// X has no manage → cannot remove anyone.
	if err := m.Remove(context.Background(), 3, 4, ReparentRoot); err == nil {
		t.Fatal("manage yetkisi olmayan çıkaramamalı")
	}
	// B cannot remove X (not in B's subtree).
	if err := m.Remove(context.Background(), 2, 3, ReparentRoot); err == nil {
		t.Fatal("alt ağaç dışındaki kişi çıkarılamamalı")
	}
}

func TestCanRejectsUnknownAndMissingCap(t *testing.T) {
	m := newTestMgr(t, 1)
	mustInvite(t, m, 1, 2, "B", CapStatus)
	if m.Can(2, CapShutdown) {
		t.Fatal("B shutdown yapamamalı")
	}
	if !m.Can(2, CapStatus) {
		t.Fatal("B status yapabilmeli")
	}
	if m.Can(999, CapStatus) {
		t.Fatal("bilinmeyen kullanıcı hiçbir şey yapamamalı")
	}
}

func mustInvite(t *testing.T, m *Manager, inviter, newChat int64, name string, cs ...Capability) {
	t.Helper()
	if _, err := m.Invite(context.Background(), inviter, newChat, name, cs); err != nil {
		t.Fatalf("invite %d→%d: %v", inviter, newChat, err)
	}
}
