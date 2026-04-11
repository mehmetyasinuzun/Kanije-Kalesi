//go:build windows

package lock

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// acquire attempts to create a named Windows Mutex.
// If the mutex already exists, another instance is running.
// The named mutex is automatically released when the process exits — no cleanup needed.
func acquire(name string) (Releaser, error) {
	mutexName := "Global\\KanijeKalesi_" + name
	utf16Name, err := windows.UTF16PtrFromString(mutexName)
	if err != nil {
		return nil, fmt.Errorf("mutex adı dönüştürme hatası: %w", err)
	}

	handle, err := windows.CreateMutex(nil, false, utf16Name)
	if err != nil {
		return nil, fmt.Errorf("mutex oluşturma hatası: %w", err)
	}

	// ERROR_ALREADY_EXISTS (183) means another instance holds the mutex
	if windows.GetLastError() == windows.ERROR_ALREADY_EXISTS {
		windows.CloseHandle(handle)
		return nil, ErrAlreadyRunning
	}

	return &windowsLock{handle: handle}, nil
}

type windowsLock struct {
	handle windows.Handle
}

func (l *windowsLock) Release() error {
	return windows.CloseHandle(l.handle)
}
