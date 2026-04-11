package event

import (
	"fmt"
	"sync/atomic"
	"time"
)

// BusConfig holds tuneable parameters for the event bus.
type BusConfig struct {
	BufferSize      int           // Channel buffer depth
	MaxPerMinute    int           // Rate limit: max events per type per minute
	DedupWindow     time.Duration // Events with same key within this window are dropped
}

// DefaultBusConfig returns safe production defaults.
func DefaultBusConfig() BusConfig {
	return BusConfig{
		BufferSize:   512,
		MaxPerMinute: 10,
		DedupWindow:  3 * time.Second,
	}
}

// Bus is the central event pipeline. Listeners publish events; the dispatcher
// reads them. All methods are safe for concurrent use.
type Bus struct {
	ch      chan Event
	limiter *tokenBucket
	dedup   *dedupCache

	totalReceived atomic.Int64
	totalDropped  atomic.Int64
	totalDeduped  atomic.Int64
}

// NewBus creates a new event bus with the given configuration.
func NewBus(cfg BusConfig) *Bus {
	return &Bus{
		ch:      make(chan Event, cfg.BufferSize),
		limiter: newTokenBucket(cfg.MaxPerMinute, time.Minute/time.Duration(cfg.MaxPerMinute)),
		dedup:   newDedupCache(cfg.DedupWindow),
	}
}

// Publish attempts to add an event to the bus.
// Returns false (without blocking) if the event is rate-limited, deduplicated,
// or the buffer is full.
func (b *Bus) Publish(ev Event) bool {
	// Dedup check: same type + username within window
	dedupKey := fmt.Sprintf("%s:%s:%s", ev.Type, ev.Username, ev.SourceIP)
	if b.dedup.IsDuplicate(dedupKey) {
		b.totalDeduped.Add(1)
		return false
	}

	// Rate limit check per event type
	if !b.limiter.Allow(string(ev.Type)) {
		b.totalDropped.Add(1)
		return false
	}

	// Non-blocking send
	select {
	case b.ch <- ev:
		b.totalReceived.Add(1)
		return true
	default:
		// Buffer full
		b.totalDropped.Add(1)
		return false
	}
}

// Events returns the read-only channel for the dispatcher.
func (b *Bus) Events() <-chan Event {
	return b.ch
}

// Stats returns current counters for the /status command.
func (b *Bus) Stats() BusStats {
	return BusStats{
		Received: b.totalReceived.Load(),
		Dropped:  b.totalDropped.Load(),
		Deduped:  b.totalDeduped.Load(),
		Pending:  int64(len(b.ch)),
	}
}

// BusStats contains bus metrics.
type BusStats struct {
	Received int64
	Dropped  int64
	Deduped  int64
	Pending  int64
}
