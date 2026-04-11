// Package lock provides a cross-platform single-instance guard.
// On Windows, it uses a named Mutex. On Unix, it uses flock on a temp file.
package lock

import "errors"

// ErrAlreadyRunning is returned when another instance is already running.
var ErrAlreadyRunning = errors.New("kanije zaten çalışıyor")

// Releaser releases the lock when the process is done.
type Releaser interface {
	Release() error
}

// Acquire attempts to claim the singleton lock for the given name.
// Returns ErrAlreadyRunning if another instance holds the lock.
func Acquire(name string) (Releaser, error) {
	return acquire(name)
}
