// Package listener defines the Listener interface and common utilities.
// Platform-specific implementations live in sub-packages (windows/, linux/).
package listener

import (
	"context"

	"github.com/kanije-kalesi/sentinel/internal/event"
)

// Listener watches a source of security events and publishes them to the bus.
// Implementations must be goroutine-safe and honour context cancellation.
type Listener interface {
	// Name returns a short human-readable identifier for logging.
	// Example: "EventLog", "USBMonitor", "PowerMonitor"
	Name() string

	// Start begins monitoring and blocks until ctx is cancelled or a fatal
	// error occurs. It must not leak goroutines after returning.
	// Transient errors should be handled internally with back-off;
	// only unrecoverable errors should be returned.
	Start(ctx context.Context, bus *event.Bus) error
}
