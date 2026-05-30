package storage

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/kanije-kalesi/kanije/internal/access"
)

func TestSQLiteUserAndAuditRoundTrip(t *testing.T) {
	ctx := context.Background()
	s, err := NewSQLite(filepath.Join(t.TempDir(), "acl.db"))
	if err != nil {
		t.Fatalf("NewSQLite: %v", err)
	}
	defer s.Close()

	// Verify it satisfies the access.Persistence contract.
	var _ access.Persistence = s

	root := access.User{ChatID: 1, Name: "Sahip", InvitedBy: 0, Caps: access.AllCapabilities(), AddedAt: time.Now()}
	child := access.User{ChatID: 2, Name: "B", InvitedBy: 1, Caps: []access.Capability{access.CapStatus, access.CapPhoto}, AddedAt: time.Now()}
	if err := s.SaveUser(ctx, root); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveUser(ctx, child); err != nil {
		t.Fatal(err)
	}

	users, err := s.LoadUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("2 kullanıcı bekleniyordu, %d", len(users))
	}

	// Update caps (upsert) and confirm the change persists.
	child.Caps = []access.Capability{access.CapStatus}
	if err := s.SaveUser(ctx, child); err != nil {
		t.Fatal(err)
	}
	got := findUser(t, s, 2)
	if len(got.Caps) != 1 || got.Caps[0] != access.CapStatus {
		t.Fatalf("güncellenen yetki kalıcı olmadı: %v", got.Caps)
	}
	if got.InvitedBy != 1 {
		t.Fatalf("invited_by kalıcı olmadı: %d", got.InvitedBy)
	}

	// Delete and confirm.
	if err := s.DeleteUser(ctx, 2); err != nil {
		t.Fatal(err)
	}
	if users, _ := s.LoadUsers(ctx); len(users) != 1 {
		t.Fatalf("silme sonrası 1 kullanıcı bekleniyordu, %d", len(users))
	}

	// Audit round-trip (newest first).
	for _, a := range []string{"foto", "ekle", "kapat"} {
		if err := s.AppendAudit(ctx, access.AuditEntry{ChatID: 1, Actor: "Sahip", Action: a, At: time.Now()}); err != nil {
			t.Fatal(err)
		}
	}
	entries, err := s.RecentAudit(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("2 denetim kaydı bekleniyordu, %d", len(entries))
	}
	if entries[0].Action != "kapat" {
		t.Fatalf("en yeni kayıt 'kapat' olmalı, %q", entries[0].Action)
	}
}

func findUser(t *testing.T, s *SQLiteStorage, id int64) access.User {
	t.Helper()
	users, err := s.LoadUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range users {
		if u.ChatID == id {
			return u
		}
	}
	t.Fatalf("kullanıcı %d bulunamadı", id)
	return access.User{}
}
