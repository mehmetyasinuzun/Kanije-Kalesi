// Package storage defines the persistence interface and types.
package storage

import (
	"context"
	"time"

	"github.com/kanije-kalesi/kanije/internal/event"
)

// Storage is the interface for all persistent data operations.
// The SQLite implementation is the only production implementation;
// a no-op in-memory implementation can be used for testing.
type Storage interface {
	// SaveEvent persists a security event. ID is set by the storage layer.
	SaveEvent(ctx context.Context, ev event.Event) error

	// RecentEvents returns the last n events, newest first.
	RecentEvents(ctx context.Context, n int) ([]event.Event, error)

	// EventByID returns a single event by its ID. ok is false if not found.
	EventByID(ctx context.Context, id int64) (ev event.Event, ok bool, err error)

	// QueryEvents returns events matching the filter, newest first.
	QueryEvents(ctx context.Context, filter EventFilter) ([]event.Event, error)

	// CountEvents returns the total number of stored events.
	CountEvents(ctx context.Context) (int64, error)

	// SavePendingMessage queues a message for offline delivery.
	SavePendingMessage(ctx context.Context, text string) error

	// PendingMessages returns all queued offline messages, oldest first,
	// WITHOUT removing them. Messages are deleted only after a confirmed send
	// via DeletePendingMessages — this gives crash-safe at-least-once delivery
	// instead of dropping the queue on the first send failure.
	PendingMessages(ctx context.Context) ([]PendingMessage, error)

	// DeletePendingMessages removes the given message IDs from the queue.
	DeletePendingMessages(ctx context.Context, ids []int64) error

	// Prune removes events older than retentionDays. Called periodically.
	Prune(ctx context.Context, retentionDays int) (int64, error)

	// VerifyChain recomputes the tamper-evident hash chain over all events.
	// ok=true (brokenAt=0) means intact; otherwise brokenAt is the first
	// event ID whose stored hash no longer matches. total is the event count.
	VerifyChain(ctx context.Context) (ok bool, brokenAt int64, total int64, err error)

	// EventStats returns per-type counts for events at or after since, busiest
	// type first, plus the overall total.
	EventStats(ctx context.Context, since time.Time) (counts []TypeCount, total int64, err error)

	// Close releases all resources.
	Close() error
}

// EventFilter specifies criteria for querying events.
type EventFilter struct {
	Since  time.Time
	Until  time.Time
	Type   event.Type // Empty = all types
	Limit  int        // 0 = no limit
	Offset int
}

// PendingMessage is a queued offline message.
type PendingMessage struct {
	ID        int64
	Text      string
	CreatedAt time.Time
}

// TypeCount is a per-event-type count, used by the /ozet summary.
type TypeCount struct {
	Type  event.Type
	Count int64
}
