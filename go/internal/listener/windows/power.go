//go:build windows

package windows

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"unsafe"

	"github.com/kanije-kalesi/kanije/internal/event"
	"golang.org/x/sys/windows"
)

const (
	// WM_POWERBROADCAST power event codes
	pbtApmSuspend         = 0x0004 // System is suspending
	pbtApmResumeSuspend   = 0x0007 // Manual resume from suspend
	pbtApmResumeAutomatic = 0x0012 // Automatic resume (timer/network)
	pbtPowersettingChange = 0x8013 // Power setting change

	wmPowerBroadcast = 0x0218
)

// PowerMonitor detects system sleep and wake events via WM_POWERBROADCAST.
// It uses a message-only window — same pattern as USBMonitor.
type PowerMonitor struct {
	hostname string
	log      *slog.Logger
}

func NewPowerMonitor(log *slog.Logger) *PowerMonitor {
	h, _ := os.Hostname()
	return &PowerMonitor{hostname: h, log: log}
}

func (m *PowerMonitor) Name() string { return "PowerMonitor" }

func (m *PowerMonitor) Start(ctx context.Context, bus *event.Bus) error {
	type result struct{ err error }
	done := make(chan result, 1)

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		hwnd, err := createPowerWindow()
		if err != nil {
			done <- result{fmt.Errorf("güç izleme penceresi oluşturulamadı: %w", err)}
			return
		}
		defer procDestroyWindow.Call(hwnd)

		m.log.Info("Güç olayı izleme başlatıldı")

		go func() {
			<-ctx.Done()
			procPostMessageW.Call(hwnd, wmQuit, 0, 0)
		}()

		var message msg
		for {
			ret, _, _ := procGetMessageW.Call(
				uintptr(unsafe.Pointer(&message)),
				0, 0, 0,
			)
			if ret == 0 || int32(ret) == -1 {
				break
			}

			if message.Message == wmPowerBroadcast {
				m.handlePower(message.WParam, bus)
			}

			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&message)))
		}

		done <- result{nil}
	}()

	return (<-done).err
}

func (m *PowerMonitor) handlePower(wParam uintptr, bus *event.Bus) {
	switch wParam {
	case pbtApmSuspend:
		ev := event.New(event.TypeSystemSleep, "PowerMonitor")
		ev.Hostname = m.hostname
		m.log.Info("Sistem uyku moduna giriyor")
		bus.Publish(ev)

	case pbtApmResumeSuspend:
		ev := event.New(event.TypeSystemWake, "PowerMonitor")
		ev.Hostname = m.hostname
		ev.WakeType = "manuel"
		m.log.Info("Sistem uykudan uyandı (manuel)")
		bus.Publish(ev)

	case pbtApmResumeAutomatic:
		ev := event.New(event.TypeSystemWake, "PowerMonitor")
		ev.Hostname = m.hostname
		ev.WakeType = "otomatik"
		m.log.Info("Sistem uykudan uyandı (otomatik)")
		bus.Publish(ev)
	}
}

func createPowerWindow() (uintptr, error) {
	hInstance, _, _ := procGetModuleHandleW.Call(0)
	className, _ := windows.UTF16PtrFromString("KanijePowerWatcher")

	wc := wndClassExW{
		Size:      uint32(unsafe.Sizeof(wndClassExW{})),
		Instance:  hInstance,
		ClassName: className,
		WndProc:   windows.NewCallback(defaultWndProc),
	}

	procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	hwnd, _, createErr := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		0,
		0, 0, 0, 0, 0,
		0xFFFF, // HWND_MESSAGE
		0, hInstance, 0,
	)
	if hwnd == 0 {
		return 0, fmt.Errorf("CreateWindowEx hatası: %v", createErr)
	}
	return hwnd, nil
}
