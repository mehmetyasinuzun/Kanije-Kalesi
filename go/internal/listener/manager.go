package listener

import (
	"context"
	"log/slog"
	"time"

	"github.com/kanije-kalesi/sentinel/internal/event"
	"golang.org/x/sync/errgroup"
)

const (
	restartDelay    = 5 * time.Second
	maxRestartDelay = 2 * time.Minute
)

// Manager supervises a set of Listeners. If a listener exits with an error,
// the manager restarts it with exponential back-off — unless ctx is done.
type Manager struct {
	listeners []Listener
	log       *slog.Logger
}

// NewManager creates a manager with the given listeners.
func NewManager(log *slog.Logger, listeners ...Listener) *Manager {
	return &Manager{
		listeners: listeners,
		log:       log,
	}
}

// Run starts all listeners concurrently and blocks until ctx is cancelled.
// It returns when all supervised goroutines have exited.
func (m *Manager) Run(ctx context.Context, bus *event.Bus) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, l := range m.listeners {
		l := l // capture loop variable
		g.Go(func() error {
			return m.supervise(ctx, l, bus)
		})
	}

	return g.Wait()
}

// supervise runs a single listener with automatic restart on error.
func (m *Manager) supervise(ctx context.Context, l Listener, bus *event.Bus) error {
	delay := restartDelay

	for {
		m.log.Info("listener başlatılıyor", "listener", l.Name())
		err := l.Start(ctx, bus)

		if ctx.Err() != nil {
			// Context cancelled — normal shutdown, not an error
			m.log.Info("listener durduruldu", "listener", l.Name())
			return nil
		}

		if err != nil {
			m.log.Error("listener hata verdi, yeniden başlatılıyor",
				"listener", l.Name(),
				"err", err,
				"delay", delay,
			)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(delay):
			// Exponential back-off
			delay *= 2
			if delay > maxRestartDelay {
				delay = maxRestartDelay
			}
		}
	}
}
