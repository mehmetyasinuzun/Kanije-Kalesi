package event

import (
	"sync"
	"time"
)

// tokenBucket implements a per-key token bucket rate limiter.
// Thread-safe.
type tokenBucket struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	capacity int           // max tokens
	refill   time.Duration // time between refills
}

type bucket struct {
	tokens   int
	lastSeen time.Time
}

func newTokenBucket(capacity int, refillInterval time.Duration) *tokenBucket {
	return &tokenBucket{
		buckets:  make(map[string]*bucket),
		capacity: capacity,
		refill:   refillInterval,
	}
}

// Allow returns true if the given key is below the rate limit.
func (tb *tokenBucket) Allow(key string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	b, ok := tb.buckets[key]
	if !ok {
		tb.buckets[key] = &bucket{tokens: tb.capacity - 1, lastSeen: now}
		return true
	}

	// Refill based on elapsed time
	elapsed := now.Sub(b.lastSeen)
	refills := int(elapsed / tb.refill)
	if refills > 0 {
		b.tokens += refills
		if b.tokens > tb.capacity {
			b.tokens = tb.capacity
		}
		b.lastSeen = now
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// dedupCache prevents identical events from flooding within a time window.
// Thread-safe.
type dedupCache struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
}

func newDedupCache(window time.Duration) *dedupCache {
	dc := &dedupCache{
		seen:   make(map[string]time.Time),
		window: window,
	}
	go dc.cleanup()
	return dc
}

// IsDuplicate returns true if this key was seen within the dedup window.
func (dc *dedupCache) IsDuplicate(key string) bool {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if t, ok := dc.seen[key]; ok {
		if time.Since(t) < dc.window {
			return true
		}
	}
	dc.seen[key] = time.Now()
	return false
}

// cleanup periodically removes stale entries to prevent memory leak.
func (dc *dedupCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		dc.mu.Lock()
		now := time.Now()
		for k, t := range dc.seen {
			if now.Sub(t) > dc.window*10 {
				delete(dc.seen, k)
			}
		}
		dc.mu.Unlock()
	}
}
