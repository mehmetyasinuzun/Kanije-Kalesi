package storage

import (
	"context"
	"strings"
	"time"

	"github.com/kanije-kalesi/kanije/internal/access"
)

// This file makes *SQLiteStorage satisfy access.Persistence (structurally —
// no interface declaration needed). Capabilities are stored as a comma-joined
// string in a single column; the set is small and fixed.

func capsToString(caps []access.Capability) string {
	if len(caps) == 0 {
		return ""
	}
	parts := make([]string, len(caps))
	for i, c := range caps {
		parts[i] = string(c)
	}
	return strings.Join(parts, ",")
}

func capsFromString(s string) []access.Capability {
	if s == "" {
		return nil
	}
	raw := strings.Split(s, ",")
	out := make([]access.Capability, 0, len(raw))
	for _, p := range raw {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, access.Capability(p))
		}
	}
	return out
}

// LoadUsers returns every user in the delegation tree.
func (s *SQLiteStorage) LoadUsers(ctx context.Context) ([]access.User, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT chat_id, name, invited_by, caps, added_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []access.User
	for rows.Next() {
		var (
			u       access.User
			capsStr string
			addedAt string
		)
		if err := rows.Scan(&u.ChatID, &u.Name, &u.InvitedBy, &capsStr, &addedAt); err != nil {
			return nil, err
		}
		u.Caps = capsFromString(capsStr)
		if t, err := time.Parse(time.RFC3339Nano, addedAt); err == nil {
			u.AddedAt = t
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// SaveUser inserts or updates a user (upsert by chat_id).
func (s *SQLiteStorage) SaveUser(ctx context.Context, u access.User) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO users (chat_id, name, invited_by, caps, added_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(chat_id) DO UPDATE SET
			name       = excluded.name,
			invited_by = excluded.invited_by,
			caps       = excluded.caps`,
		u.ChatID, u.Name, u.InvitedBy, capsToString(u.Caps),
		u.AddedAt.UTC().Format(time.RFC3339Nano))
	return err
}

// DeleteUser removes a single user by chat ID.
func (s *SQLiteStorage) DeleteUser(ctx context.Context, chatID int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE chat_id = ?`, chatID)
	return err
}

// AppendAudit records a privileged action.
func (s *SQLiteStorage) AppendAudit(ctx context.Context, e access.AuditEntry) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO audit_log (chat_id, actor, action, target, at) VALUES (?, ?, ?, ?, ?)`,
		e.ChatID, e.Actor, e.Action, e.Target, e.At.UTC().Format(time.RFC3339Nano))
	return err
}

// RecentAudit returns the latest n audit entries, newest first.
func (s *SQLiteStorage) RecentAudit(ctx context.Context, n int) ([]access.AuditEntry, error) {
	if n <= 0 {
		n = 20
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, chat_id, actor, action, target, at FROM audit_log ORDER BY id DESC LIMIT ?`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []access.AuditEntry
	for rows.Next() {
		var (
			e  access.AuditEntry
			at string
		)
		if err := rows.Scan(&e.ID, &e.ChatID, &e.Actor, &e.Action, &e.Target, &at); err != nil {
			return nil, err
		}
		if t, err := time.Parse(time.RFC3339Nano, at); err == nil {
			e.At = t
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
