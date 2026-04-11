package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kanije-kalesi/sentinel/internal/event"
	_ "modernc.org/sqlite" // Pure-Go SQLite driver (no CGo — works in cross-compiled binaries)
)

const schema = `
CREATE TABLE IF NOT EXISTS events (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    type        TEXT NOT NULL,
    severity    INTEGER NOT NULL DEFAULT 0,
    timestamp   TEXT NOT NULL,
    source      TEXT,
    hostname    TEXT,
    username    TEXT,
    source_ip   TEXT,
    network_ssid TEXT,
    network_type TEXT,
    local_ip    TEXT,
    logon_type  INTEGER,
    domain      TEXT,
    device_name TEXT,
    device_label TEXT,
    device_size INTEGER,
    device_fs   TEXT,
    device_path TEXT,
    wake_type   TEXT,
    extra       TEXT
);

CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_events_type      ON events(type);

CREATE TABLE IF NOT EXISTS pending_messages (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    text       TEXT NOT NULL,
    created_at TEXT NOT NULL
);

-- Enable WAL mode for better concurrent read/write performance
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
PRAGMA foreign_keys=ON;
PRAGMA encoding='UTF-8';
`

// SQLiteStorage is a Storage implementation backed by a local SQLite database.
// It is safe for concurrent use from multiple goroutines.
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLite opens (or creates) a SQLite database at dbPath.
func NewSQLite(dbPath string) (*SQLiteStorage, error) {
	// The modernc.org/sqlite driver uses "sqlite" as the driver name
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)")
	if err != nil {
		return nil, fmt.Errorf("SQLite açılamadı (%s): %w", dbPath, err)
	}

	// Connection pool: single writer, multiple readers
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("şema oluşturma hatası: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

// SaveEvent writes an event to the database.
// All string values are stored as UTF-8 (SQLite's native encoding).
func (s *SQLiteStorage) SaveEvent(ctx context.Context, ev event.Event) error {
	var extraJSON []byte
	if len(ev.Extra) > 0 {
		var err error
		extraJSON, err = json.Marshal(ev.Extra)
		if err != nil {
			extraJSON = nil
		}
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO events (
			type, severity, timestamp, source, hostname, username,
			source_ip, network_ssid, network_type, local_ip,
			logon_type, domain,
			device_name, device_label, device_size, device_fs, device_path,
			wake_type, extra
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?,
			?, ?, ?, ?, ?,
			?, ?
		)`,
		string(ev.Type),
		int(ev.Severity),
		ev.Timestamp.UTC().Format(time.RFC3339Nano),
		ev.Source,
		ev.Hostname,
		ev.Username,
		ev.SourceIP,
		ev.NetworkSSID,
		ev.NetworkType,
		ev.LocalIP,
		int(ev.LogonType),
		ev.Domain,
		ev.DeviceName,
		ev.DeviceLabel,
		ev.DeviceSize,
		ev.DeviceFS,
		ev.DevicePath,
		ev.WakeType,
		string(extraJSON),
	)
	return err
}

// RecentEvents returns the last n events, newest first.
func (s *SQLiteStorage) RecentEvents(ctx context.Context, n int) ([]event.Event, error) {
	if n <= 0 {
		n = 10
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, type, severity, timestamp, source, hostname, username,
		       source_ip, network_ssid, network_type, local_ip,
		       logon_type, domain,
		       device_name, device_label, device_size, device_fs, device_path,
		       wake_type, extra
		FROM events
		ORDER BY timestamp DESC
		LIMIT ?`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvents(rows)
}

// QueryEvents returns events matching the filter.
func (s *SQLiteStorage) QueryEvents(ctx context.Context, filter EventFilter) ([]event.Event, error) {
	query := `SELECT id, type, severity, timestamp, source, hostname, username,
		       source_ip, network_ssid, network_type, local_ip,
		       logon_type, domain,
		       device_name, device_label, device_size, device_fs, device_path,
		       wake_type, extra
		FROM events WHERE 1=1`
	var args []any

	if !filter.Since.IsZero() {
		query += " AND timestamp >= ?"
		args = append(args, filter.Since.UTC().Format(time.RFC3339Nano))
	}
	if !filter.Until.IsZero() {
		query += " AND timestamp <= ?"
		args = append(args, filter.Until.UTC().Format(time.RFC3339Nano))
	}
	if filter.Type != "" {
		query += " AND type = ?"
		args = append(args, string(filter.Type))
	}
	query += " ORDER BY timestamp DESC"
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", filter.Offset)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvents(rows)
}

// CountEvents returns the total event count.
func (s *SQLiteStorage) CountEvents(ctx context.Context) (int64, error) {
	var n int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM events").Scan(&n)
	return n, err
}

// SavePendingMessage queues a message for offline delivery.
func (s *SQLiteStorage) SavePendingMessage(ctx context.Context, text string) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO pending_messages (text, created_at) VALUES (?, ?)",
		text, time.Now().UTC().Format(time.RFC3339))
	return err
}

// PopPendingMessages atomically retrieves and deletes all pending messages.
func (s *SQLiteStorage) PopPendingMessages(ctx context.Context) ([]PendingMessage, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, "SELECT id, text, created_at FROM pending_messages ORDER BY id ASC")
	if err != nil {
		return nil, err
	}

	var msgs []PendingMessage
	for rows.Next() {
		var m PendingMessage
		var tsStr string
		if err := rows.Scan(&m.ID, &m.Text, &tsStr); err != nil {
			rows.Close()
			return nil, err
		}
		m.CreatedAt, _ = time.Parse(time.RFC3339, tsStr)
		msgs = append(msgs, m)
	}
	rows.Close()

	if len(msgs) > 0 {
		if _, err := tx.ExecContext(ctx, "DELETE FROM pending_messages"); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return msgs, nil
}

// Prune removes events older than retentionDays.
func (s *SQLiteStorage) Prune(ctx context.Context, retentionDays int) (int64, error) {
	cutoff := time.Now().UTC().AddDate(0, 0, -retentionDays).Format(time.RFC3339Nano)
	result, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE timestamp < ?", cutoff)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Close shuts down the database connection cleanly.
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// ---- Row scanning ----

func scanEvents(rows *sql.Rows) ([]event.Event, error) {
	var events []event.Event
	for rows.Next() {
		var ev event.Event
		var tsStr, extraStr string
		var logonType int
		var severity int

		err := rows.Scan(
			&ev.ID,
			(*string)(&ev.Type),
			&severity,
			&tsStr,
			&ev.Source,
			&ev.Hostname,
			&ev.Username,
			&ev.SourceIP,
			&ev.NetworkSSID,
			&ev.NetworkType,
			&ev.LocalIP,
			&logonType,
			&ev.Domain,
			&ev.DeviceName,
			&ev.DeviceLabel,
			&ev.DeviceSize,
			&ev.DeviceFS,
			&ev.DevicePath,
			&ev.WakeType,
			&extraStr,
		)
		if err != nil {
			return nil, fmt.Errorf("satır okuma hatası: %w", err)
		}

		ev.Severity = event.Severity(severity)
		ev.LogonType = event.LogonType(logonType)
		ev.Timestamp, _ = time.Parse(time.RFC3339Nano, tsStr)
		ev.Timestamp = ev.Timestamp.Local()

		if extraStr != "" && extraStr != "null" {
			json.Unmarshal([]byte(extraStr), &ev.Extra)
		}

		events = append(events, ev)
	}
	return events, rows.Err()
}
