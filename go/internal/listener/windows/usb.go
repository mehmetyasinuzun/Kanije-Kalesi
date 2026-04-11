//go:build windows

package windows

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/kanije-kalesi/kanije/internal/event"
	"golang.org/x/sys/windows"
)

// ---- Windows Device Notification API bindings ----

var (
	user32                           = windows.NewLazySystemDLL("user32.dll")
	procRegisterClassExW             = user32.NewProc("RegisterClassExW")
	procCreateWindowExW              = user32.NewProc("CreateWindowExW")
	procDefWindowProcW               = user32.NewProc("DefWindowProcW")
	procGetMessageW                  = user32.NewProc("GetMessageW")
	procDispatchMessageW             = user32.NewProc("DispatchMessageW")
	procDestroyWindow                = user32.NewProc("DestroyWindow")
	procRegisterDeviceNotification   = user32.NewProc("RegisterDeviceNotification")
	procUnregisterDeviceNotification = user32.NewProc("UnregisterDeviceNotification")
	procPostMessageW                 = user32.NewProc("PostMessageW")
	kernel32                         = windows.NewLazySystemDLL("kernel32.dll")
	procGetModuleHandleW             = kernel32.NewProc("GetModuleHandleW")
)

const (
	wmDeviceChange     = 0x0219
	dbtDeviceArrival   = 0x8000
	dbtDeviceRemoveCom = 0x8004
	dbtDevTypVolume    = 0x00000002

	dbccSizeVolume           = 16
	deviceNotifyWindowHandle = 0
	deviceNotifyAllInterface = 4

	wmClose = 0x0010
	wmQuit  = 0x0012

	// Kanije WM message to signal graceful shutdown
	wmShutdownMsg = 0x0400 + 1
)

type devBroadcastHdr struct {
	Size       uint32
	DeviceType uint32
	Reserved   uint32
}

type devBroadcastVolume struct {
	Size       uint32
	DeviceType uint32
	Reserved   uint32
	UnitMask   uint32
	Flags      uint16
}

type wndClassExW struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}

type msg struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  unsafe.Pointer
	Time    uint32
	Pt      struct{ X, Y int32 }
}

// USBMonitor detects USB drive insertion and removal via WM_DEVICECHANGE.
// It creates a hidden message-only window to receive Windows device notifications.
type USBMonitor struct {
	hostname string
	log      *slog.Logger
	hwnd     uintptr // message-only window handle
}

func NewUSBMonitor(log *slog.Logger) *USBMonitor {
	h, _ := os.Hostname()
	return &USBMonitor{hostname: h, log: log}
}

func (m *USBMonitor) Name() string { return "USBMonitor" }

// Start creates a hidden window, registers for device notifications, and
// pumps messages until ctx is cancelled.
func (m *USBMonitor) Start(ctx context.Context, bus *event.Bus) error {
	// We need to run the message loop on the same OS thread that creates the window.
	type result struct{ err error }
	done := make(chan result, 1)

	go func() {
		// Lock this goroutine to the OS thread (message loop requirement)
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		hwnd, err := createMessageWindow()
		if err != nil {
			done <- result{err}
			return
		}
		m.hwnd = hwnd
		defer procDestroyWindow.Call(hwnd)

		// Register for volume change notifications
		var notifyFilter struct {
			Size       uint32
			DeviceType uint32
			Reserved   uint32
			ClassGuid  [16]byte
			Name       [1]uint16
		}
		notifyFilter.Size = uint32(unsafe.Sizeof(notifyFilter))
		notifyFilter.DeviceType = dbtDevTypVolume

		hNotify, _, _ := procRegisterDeviceNotification.Call(
			hwnd,
			uintptr(unsafe.Pointer(&notifyFilter)),
			deviceNotifyAllInterface,
		)
		if hNotify != 0 {
			defer procUnregisterDeviceNotification.Call(hNotify)
		}

		m.log.Info("USB izleme başlatıldı")

		// Context watcher: post WM_QUIT to unblock GetMessage when ctx is done
		go func() {
			<-ctx.Done()
			procPostMessageW.Call(hwnd, wmQuit, 0, 0)
		}()

		// Message pump
		var message msg
		for {
			ret, _, _ := procGetMessageW.Call(
				uintptr(unsafe.Pointer(&message)),
				0, 0, 0,
			)
			if ret == 0 || int32(ret) == -1 {
				break
			}

			if message.Message == wmDeviceChange {
				m.handleDeviceChange(message.WParam, message.LParam, bus)
			}

			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&message)))
		}

		done <- result{nil}
	}()

	return (<-done).err
}

// handleDeviceChange processes WM_DEVICECHANGE messages.
func (m *USBMonitor) handleDeviceChange(wParam uintptr, lParam unsafe.Pointer, bus *event.Bus) {
	if lParam == nil {
		return
	}
	hdr := (*devBroadcastHdr)(lParam)
	if hdr.DeviceType != dbtDevTypVolume {
		return
	}

	vol := (*devBroadcastVolume)(lParam)
	driveLetter := unitMaskToDriveLetter(vol.UnitMask)

	switch wParam {
	case dbtDeviceArrival:
		ev := event.New(event.TypeUSBInserted, "USBMonitor")
		ev.Hostname = m.hostname
		ev.DevicePath = driveLetter + `:\`
		ev.DeviceLabel = getDriveLabel(driveLetter)
		ev.DeviceName = ev.DeviceLabel
		ev.DeviceFS = getDriveFS(driveLetter)
		ev.DeviceSize = getDriveSize(driveLetter)
		m.log.Info("USB takıldı", "sürücü", driveLetter, "etiket", ev.DeviceLabel)
		bus.Publish(ev)

	case dbtDeviceRemoveCom:
		ev := event.New(event.TypeUSBRemoved, "USBMonitor")
		ev.Hostname = m.hostname
		ev.DevicePath = driveLetter + `:\`
		ev.DeviceName = driveLetter
		m.log.Info("USB çıkarıldı", "sürücü", driveLetter)
		bus.Publish(ev)
	}
}

// unitMaskToDriveLetter converts a Windows drive bitmask to a drive letter string.
// UnitMask bit 0 = A:, bit 1 = B:, bit 2 = C:, etc.
func unitMaskToDriveLetter(mask uint32) string {
	for i := 0; i < 26; i++ {
		if mask&(1<<uint(i)) != 0 {
			return string(rune('A' + i))
		}
	}
	return "?"
}

func getDriveLabel(letter string) string {
	path := letter + `:\`
	buf := make([]uint16, 256)
	err := windows.GetVolumeInformation(
		windows.StringToUTF16Ptr(path),
		&buf[0], uint32(len(buf)),
		nil, nil, nil, nil, 0,
	)
	if err != nil {
		return ""
	}
	return windows.UTF16ToString(buf)
}

func getDriveFS(letter string) string {
	path := letter + `:\`
	volBuf := make([]uint16, 256)
	fsBuf := make([]uint16, 256)
	err := windows.GetVolumeInformation(
		windows.StringToUTF16Ptr(path),
		&volBuf[0], uint32(len(volBuf)),
		nil, nil, nil,
		&fsBuf[0], uint32(len(fsBuf)),
	)
	if err != nil {
		return ""
	}
	return windows.UTF16ToString(fsBuf)
}

func getDriveSize(letter string) int64 {
	path := letter + `:\`
	var freeBytes, totalBytes, totalFree uint64
	p, _ := windows.UTF16PtrFromString(path)
	err := windows.GetDiskFreeSpaceEx(p, &freeBytes, &totalBytes, &totalFree)
	if err != nil {
		return 0
	}
	return int64(totalBytes)
}

func createMessageWindow() (uintptr, error) {
	hInstance, _, _ := procGetModuleHandleW.Call(0)
	className, _ := windows.UTF16PtrFromString("KanijeUSBWatcher")

	wc := wndClassExW{
		Size:      uint32(unsafe.Sizeof(wndClassExW{})),
		Instance:  hInstance,
		ClassName: className,
		WndProc:   windows.NewCallback(defaultWndProc),
	}

	ret, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	if ret == 0 {
		// May already be registered from a previous run — that's OK
		_ = err
	}

	hwnd, _, createErr := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		0,
		0, 0, 0, 0, 0,
		0xFFFF, // HWND_MESSAGE — message-only window
		0, hInstance, 0,
	)
	if hwnd == 0 {
		return 0, fmt.Errorf("CreateWindowEx hatası: %v", createErr)
	}
	return hwnd, nil
}

func defaultWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	ret, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return ret
}

// getDriveLetterFromPath extracts the drive letter from a volume path.
func getDriveLetterFromPath(p string) string {
	p = strings.ToUpper(filepath.VolumeName(p))
	if len(p) > 0 {
		return string(p[0])
	}
	return "?"
}

// volumeInfo holds cached info about a USB volume.
type volumeInfo struct {
	label string
	fs    string
	size  int64
	ts    time.Time
}
