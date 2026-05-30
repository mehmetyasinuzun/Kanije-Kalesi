package telegram

import (
	"sync"
	"time"
)

// cmdRateLimiter applies a per-chat token-bucket limit to incoming bot commands.
// Even an authorized chat cannot spam expensive commands (e.g. /foto) and lock
// up the camera or flood the system. Safe for concurrent use.
type cmdRateLimiter struct {
	mu       sync.Mutex
	buckets  map[int64]*cmdBucket
	capacity int
	window   time.Duration
}

type cmdBucket struct {
	tokens int
	last   time.Time
}

func newCmdRateLimiter(perMinute int) *cmdRateLimiter {
	if perMinute < 1 {
		perMinute = 20
	}
	return &cmdRateLimiter{
		buckets:  make(map[int64]*cmdBucket),
		capacity: perMinute,
		window:   time.Minute,
	}
}

// allow reports whether a command from chatID may proceed, consuming one token.
func (r *cmdRateLimiter) allow(chatID int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	b, ok := r.buckets[chatID]
	if !ok {
		r.buckets[chatID] = &cmdBucket{tokens: r.capacity - 1, last: now}
		return true
	}

	// Refill proportionally to elapsed time.
	if elapsed := now.Sub(b.last); elapsed > 0 {
		if refill := int(float64(r.capacity) * (elapsed.Seconds() / r.window.Seconds())); refill > 0 {
			b.tokens += refill
			if b.tokens > r.capacity {
				b.tokens = r.capacity
			}
			b.last = now
		}
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}
