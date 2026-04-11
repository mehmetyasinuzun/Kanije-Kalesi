//go:build !windows

package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func acquire(name string) (Releaser, error) {
	dir := os.TempDir()
	path := filepath.Join(dir, "kanije_"+name+".lock")

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("lock dosyası açılamadı: %w", err)
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		// Read PID of the existing holder for a better error message
		pid := readPID(path)
		return nil, fmt.Errorf("%w (PID %d zaten çalışıyor)", ErrAlreadyRunning, pid)
	}

	// Write our PID so users can identify the running instance
	f.Truncate(0)
	f.WriteString(strconv.Itoa(os.Getpid()))

	return &unixLock{f: f, path: path}, nil
}

type unixLock struct {
	f    *os.File
	path string
}

func (l *unixLock) Release() error {
	l.f.Close()
	os.Remove(l.path)
	return nil
}

func readPID(path string) int {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(b)))
	return n
}
