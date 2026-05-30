package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kanije-kalesi/kanije/internal/event"
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
    extra       TEXT,
    hash        TEXT,
    prev_hash   TEXT
);

CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_events_type      ON events(type);

CREATE TABLE IF NOT EXISTS pending_messages (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    text       TEXT NOT NULL,
    created_at TEXT NOT NULL
);

-- Multi-user delegation tree (capability-based access control).
CREATE TABLE IF NOT EXISTS users (
    chat_id    INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    invited_by INTEGER NOT NULL DEFAULT 0,
    caps       TEXT NOT NULL DEFAULT '',
    added_at   TEXT NOT NULL
);

-- Privileged-action audit trail (/loglar).
CREATE TABLE IF NOT EXISTS audit_log (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id INTEGER NOT NULL,
    actor   TEXT,
    action  TEXT NOT NULL,
    target  TEXT,
    at      TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_audit_at ON audit_log(at DESC);

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
	// The modernc.org/sqlite driver uses "sqlite" as the driver name.
	// WAL + a busy timeout keep writes durable while tolerating brief lock
	// contention instead of failing immediately with "database is locked".
	dsn := dbPath + "?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("SQLite açılamadı (%s): %w", dbPath, err)
	}

	// Serialize access through a single connection. Event volume is low and this
	// avoids writer/writer lock contention entirely; reads are sub-millisecond.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("şema oluşturma hatası: %w", err)
	}

	// Migrate older databases that predate the tamper-evident hash chain.
	if err := ensureColumns(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("şema geçişi hatası: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

// ensureColumns adds columns introduced after the initial schema to existing
// databases (SQLite has no "ADD COLUMN IF NOT EXISTS", so we probe first).
func ensureColumns(db *sql.DB) error {
	rows, err := db.Query("PRAGMA table_info(events)")
	if err != nil {
		return err
	}
	defer rows.Close()

	have := make(map[string]bool)
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return err
		}
		have[name] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, col := range []string{"hash", "prev_hash"} {
		if !have[col] {
			if _, err := db.Exec("ALTER TABLE events ADD COLUMN " + col + " TEXT"); err != nil {
				return err
			}
		}
	}
	return nil
}

// SaveEvent writes an event to the database, linking it into a tamper-evident
// hash chain: each row's hash = SHA256(previous-row-hash · canonical-fields).
// Modifying or deleting any stored row breaks the chain, which VerifyChain
// detects. All string values are stored as UTF-8 (SQLite's native encoding).
func (s *SQLiteStorage) SaveEvent(ctx context.Context, ev event.Event) error {
	var extraJSON []byte
	if len(ev.Extra) > 0 {
		if b, err := json.Marshal(ev.Extra); err == nil {
			extraJSON = b
		}
	}

	tsStr := ev.Timestamp.UTC().Format(time.RFC3339Nano)
	canonical := canonicalRow(
		string(ev.Type), int(ev.Severity), tsStr, ev.Source, ev.Hostname, ev.Username,
		ev.SourceIP, ev.NetworkSSID, ev.NetworkType, ev.LocalIP, int(ev.LogonType), ev.Domain,
		ev.DeviceName, ev.DeviceLabel, ev.DeviceSize, ev.DeviceFS, ev.DevicePath, ev.WakeType,
		string(extraJSON),
	)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Chain head: the most recent row's hash (empty for the genesis event).
	var prevHash sql.NullString
	_ = tx.QueryRowContext(ctx, "SELECT hash FROM events ORDER BY id DESC LIMIT 1").Scan(&prevHash)

	thisHash := chainHash(prevHash.String, canonical)

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO events (
			type, severity, timestamp, source, hostname, username,
			source_ip, network_ssid, network_type, local_ip,
			logon_type, domain,
			device_name, device_label, device_size, device_fs, device_path,
			wake_type, extra, hash, prev_hash
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?
		)`,
		string(ev.Type), int(ev.Severity), tsStr, ev.Source, ev.Hostname, ev.Username,
		ev.SourceIP, ev.NetworkSSID, ev.NetworkType, ev.LocalIP, int(ev.LogonType), ev.Domain,
		ev.DeviceName, ev.DeviceLabel, ev.DeviceSize, ev.DeviceFS, ev.DevicePath, ev.WakeType,
		string(extraJSON), thisHash, prevHash.String,
	); err != nil {
		return err
	}
	return tx.Commit()
}

// canonicalRow builds the deterministic field representation that is hashed.
// The order and separator must never change, or old hashes stop verifying.
func canonicalRow(typ string, severity int, ts, source, hostname, username, sourceIP, ssid, ntype, localIP string, logonType int, domain, devName, devLabel string, devSize int64, devFS, devPath, wakeType, extra string) string {
	return strings.Join([]string{
		typ, strconv.Itoa(severity), ts, source, hostname, username,
		sourceIP, ssid, ntype, localIP, strconv.Itoa(logonType), domain,
		devName, devLabel, strconv.FormatInt(devSize, 10), devFS, devPath, wakeType, extra,
	}, "\x1f")
}

// chainHash returns the hex SHA-256 of (prevHash · canonical).
func chainHash(prevHash, canonical string) string {
	sum := sha256.Sum256([]byte(prevHash + "\x1f" + canonical))
	return hex.EncodeToString(sum[:])
}

// VerifyChain recomputes the hash chain and reports the first row whose stored
// hash doesn't match — evidence of tampering or deletion. ok=true with
// brokenAt=0 means the chain is intact. Rows predating the chain (NULL hash)
// are counted but skipped.
func (s *SQLiteStorage) VerifyChain(ctx context.Context) (ok bool, brokenAt int64, total int64, err error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, type, severity, timestamp, source, hostname, username,
		       source_ip, network_ssid, network_type, local_ip,
		       logon_type, domain,
		       device_name, device_label, device_size, device_fs, device_path,
		       wake_type, extra, hash, prev_hash
		FROM events ORDER BY id ASC`)
	if err != nil {
		return false, 0, 0, err
	}
	defer rows.Close()

	prev := ""
	first := true // the first hashed row's predecessor may have been pruned
	for rows.Next() {
		var id, devSize int64
		var severity, logonType int
		var typ, ts, source, hostname, username, sourceIP, ssid, ntype, localIP string
		var domain, devName, devLabel, devFS, devPath, wakeType, extra string
		var storedHash, storedPrev sql.NullString

		if err := rows.Scan(&id, &typ, &severity, &ts, &source, &hostname, &username,
			&sourceIP, &ssid, &ntype, &localIP, &logonType, &domain,
			&devName, &devLabel, &devSize, &devFS, &devPath, &wakeType, &extra,
			&storedHash, &storedPrev); err != nil {
			return false, 0, 0, err
		}

		total++
		if !storedHash.Valid || storedHash.String == "" {
			prev, first = "", true // pre-chain row; the next hashed row re-anchors
			continue
		}

		canonical := canonicalRow(typ, severity, ts, source, hostname, username,
			sourceIP, ssid, ntype, localIP, logonType, domain,
			devName, devLabel, devSize, devFS, devPath, wakeType, extra)

		// 1) Internal consistency: hash = H(stored prev · fields). Catches any
		//    modification of a row's fields or its hash.
		if chainHash(storedPrev.String, canonical) != storedHash.String {
			return false, id, total, nil
		}
		// 2) Linkage: each row links to the actual previous hash. Skipped for the
		//    first hashed row, so pruning the oldest events stays valid while a
		//    deletion in the middle still breaks the link.
		if !first && storedPrev.String != prev {
			return false, id, total, nil
		}
		prev, first = storedHash.String, false
	}
	return true, 0, total, rows.Err()
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

// EventByID returns one event by ID for the drill-down detail view.
func (s *SQLiteStorage) EventByID(ctx context.Context, id int64) (event.Event, bool, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, type, severity, timestamp, source, hostname, username,
		       source_ip, network_ssid, network_type, local_ip,
		       logon_type, domain,
		       device_name, device_label, device_size, device_fs, device_path,
		       wake_type, extra
		FROM events WHERE id = ? LIMIT 1`, id)
	if err != nil {
		return event.Event{}, false, err
	}
	defer rows.Close()

	evs, err := scanEvents(rows)
	if err != nil {
		return event.Event{}, false, err
	}
	if len(evs) == 0 {
		return event.Event{}, false, nil
	}
	return evs[0], true, nil
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

// EventStats returns per-type counts for events at or after since, busiest first.
func (s *SQLiteStorage) EventStats(ctx context.Context, since time.Time) ([]TypeCount, int64, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT type, COUNT(*) AS n FROM events
		WHERE timestamp >= ?
		GROUP BY type ORDER BY n DESC`,
		since.UTC().Format(time.RFC3339Nano))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var counts []TypeCount
	var total int64
	for rows.Next() {
		var tc TypeCount
		if err := rows.Scan((*string)(&tc.Type), &tc.Count); err != nil {
			return nil, 0, err
		}
		counts = append(counts, tc)
		total += tc.Count
	}
	return counts, total, rows.Err()
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

// PendingMessages returns all queued offline messages, oldest first, without
// deleting them. The caller deletes only successfully-sent messages via
// DeletePendingMessages, so a crash or send failure never loses the queue.
func (s *SQLiteStorage) PendingMessages(ctx context.Context) ([]PendingMessage, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, text, created_at FROM pending_messages ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []PendingMessage
	for rows.Next() {
		var m PendingMessage
		var tsStr string
		if err := rows.Scan(&m.ID, &m.Text, &tsStr); err != nil {
			return nil, err
		}
		m.CreatedAt, _ = time.Parse(time.RFC3339, tsStr)
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

// DeletePendingMessages removes the given message IDs from the queue.
func (s *SQLiteStorage) DeletePendingMessages(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := "DELETE FROM pending_messages WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
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
