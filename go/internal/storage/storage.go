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

	// QueryEvents returns events matching the filter, newest first.
	QueryEvents(ctx context.Context, filter EventFilter) ([]event.Event, error)

	// CountEvents returns the total number of stored events.
	CountEvents(ctx context.Context) (int64, error)

	// SavePendingMessage queues a message for offline delivery.
	SavePendingMessage(ctx context.Context, text string) error

	// PopPendingMessages atomically retrieves and deletes all pending messages.
	PopPendingMessages(ctx context.Context) ([]PendingMessage, error)

	// Prune removes events older than retentionDays. Called periodically.
	Prune(ctx context.Context, retentionDays int) (int64, error)

	// Close releases all resources.
	Close() error
}

// EventFilter specifies criteria for querying events.
type EventFilter struct {
	Since    time.Time
	Until    time.Time
	Type     event.Type // Empty = all types
	Limit    int        // 0 = no limit
	Offset   int
}

// PendingMessage is a queued offline message.
type PendingMessage struct {
	ID        int64
	Text      string
	CreatedAt time.Time
}
