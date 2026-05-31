//go:build windows

package windows

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
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
	procPostMessageW                 = user32.NewProc("PostMessageW")
	procRegisterDeviceNotificationW  = user32.NewProc("RegisterDeviceNotificationW")
	procUnregisterDeviceNotification = user32.NewProc("UnregisterDeviceNotification")
	kernel32                         = windows.NewLazySystemDLL("kernel32.dll")
	procGetModuleHandleW             = kernel32.NewProc("GetModuleHandleW")
)

const (
	wmDeviceChange           = 0x0219
	dbtDeviceArrival         = 0x8000
	dbtDeviceRemoveCom       = 0x8004
	dbtDevTypVolume          = 0x00000002
	dbtDevTypDeviceInterface = 0x00000005

	// RegisterDeviceNotification flags.
	deviceNotifyWindowHandle        = 0x00000000
	deviceNotifyAllInterfaceClasses = 0x00000004

	wmQuit = 0x0012

	// wsExToolWindow keeps the (never-shown) window off the taskbar.
	wsExToolWindow = 0x00000080
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

type guid struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// devBroadcastDeviceInterface mirrors DEV_BROADCAST_DEVICEINTERFACE. Name is the
// first uint16 of a variable-length, null-terminated UTF-16 device path that
// follows the struct in memory.
type devBroadcastDeviceInterface struct {
	Size       uint32
	DeviceType uint32
	Reserved   uint32
	ClassGuid  guid
	Name       uint16
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
// It creates a hidden top-level window to receive Windows volume broadcasts.
type USBMonitor struct {
	hostname string
	log      *slog.Logger
	hwnd     uintptr // hidden top-level window handle

	// Debounce: one physical device fires several interface notifications.
	mu          sync.Mutex
	lastArrival time.Time
	lastRemoval time.Time
}

func NewUSBMonitor(log *slog.Logger) *USBMonitor {
	h, _ := os.Hostname()
	return &USBMonitor{hostname: h, log: log}
}

func (m *USBMonitor) Name() string { return "USBMonitor" }

// Start creates a hidden top-level window and pumps its message loop until ctx
// is canceled, translating volume WM_DEVICECHANGE messages into events.
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

		// Register for ALL device-interface classes. THIS is what lets us hear
		// about ANY device — mouse, keyboard, phone, HID, audio, network adapter —
		// not only the storage volumes that broadcast WM_DEVICECHANGE on their own.
		ifaceFilter := devBroadcastDeviceInterface{
			Size:       uint32(unsafe.Sizeof(devBroadcastDeviceInterface{})),
			DeviceType: dbtDevTypDeviceInterface,
		}
		hDevNotify, _, _ := procRegisterDeviceNotificationW.Call(
			hwnd,
			uintptr(unsafe.Pointer(&ifaceFilter)),
			deviceNotifyWindowHandle|deviceNotifyAllInterfaceClasses,
		)
		if hDevNotify != 0 {
			defer procUnregisterDeviceNotification.Call(hDevNotify)
		} else {
			m.log.Warn("RegisterDeviceNotification başarısız — yalnızca depolama USB algılanır")
		}

		// A top-level (but invisible) window receives the broadcast WM_DEVICECHANGE
		// volume arrival/removal messages automatically — no RegisterDeviceNotification
		// is needed (and DBT_DEVTYP_VOLUME isn't even a valid filter for it). A
		// message-only window would NOT receive these broadcasts, which is why the
		// window is created top-level (see createMessageWindow).
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

// handleDeviceChange routes WM_DEVICECHANGE messages by device type: storage
// volumes get rich drive metadata; everything else (HID/phone/audio/…) is
// reported as a generic device connect/disconnect.
func (m *USBMonitor) handleDeviceChange(wParam uintptr, lParam unsafe.Pointer, bus *event.Bus) {
	if lParam == nil {
		return
	}
	switch (*devBroadcastHdr)(lParam).DeviceType {
	case dbtDevTypVolume:
		m.handleVolumeChange(wParam, lParam, bus)
	case dbtDevTypDeviceInterface:
		m.handleInterfaceChange(wParam, lParam, bus)
	}
}

// handleVolumeChange handles storage-drive arrival/removal with full metadata.
func (m *USBMonitor) handleVolumeChange(wParam uintptr, lParam unsafe.Pointer, bus *event.Bus) {
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
		ev.DeviceSize, ev.DeviceFree = getDriveUsage(driveLetter)
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

// handleInterfaceChange handles ANY device-interface arrival/removal. A single
// physical device emits several interface notifications, so arrivals and removals
// are debounced within a 2-second window to avoid spamming the owner.
func (m *USBMonitor) handleInterfaceChange(wParam uintptr, lParam unsafe.Pointer, bus *event.Bus) {
	if wParam != dbtDeviceArrival && wParam != dbtDeviceRemoveCom {
		return
	}
	di := (*devBroadcastDeviceInterface)(lParam)
	path := windows.UTF16PtrToString(&di.Name)

	now := time.Now()
	m.mu.Lock()
	last := &m.lastArrival
	if wParam == dbtDeviceRemoveCom {
		last = &m.lastRemoval
	}
	suppress := !last.IsZero() && now.Sub(*last) < 2*time.Second
	*last = now
	m.mu.Unlock()
	if suppress {
		m.log.Debug("aygıt bildirimi bastırıldı (debounce)", "path", path)
		return
	}

	kind := classifyDevice(path)
	// Storage devices ALSO raise a VOLUME broadcast (handleVolumeChange) with full
	// drive metadata, so skip them here to avoid a duplicate notification.
	if kind == "Depolama aygıtı" {
		return
	}
	typ := event.TypeDeviceConnected
	verb := "bağlandı"
	if wParam == dbtDeviceRemoveCom {
		typ = event.TypeDeviceDisconnected
		verb = "çıkarıldı"
	}
	ev := event.New(typ, "USBMonitor")
	ev.Hostname = m.hostname
	ev.DevicePath = path
	ev.DeviceName = kind
	m.log.Info("aygıt "+verb, "tür", kind)
	bus.Publish(ev)
}

// classifyDevice guesses a friendly Turkish label from a device-interface path
// like `\\?\HID#VID_046D...` or `\\?\USBSTOR#Disk...`.
func classifyDevice(path string) string {
	up := strings.ToUpper(path)
	switch {
	case strings.Contains(up, "HID"):
		return "Giriş aygıtı (klavye/fare/HID)"
	case strings.Contains(up, "USBSTOR"), strings.Contains(up, "DISK"), strings.Contains(up, "STORAGE"):
		return "Depolama aygıtı"
	case strings.Contains(up, "WPD"), strings.Contains(up, "MTP"):
		return "Telefon / medya aygıtı"
	case strings.Contains(up, "BTH"), strings.Contains(up, "BLUETOOTH"):
		return "Bluetooth aygıtı"
	case strings.Contains(up, "AUDIO"), strings.Contains(up, "MEDIA"):
		return "Ses / medya aygıtı"
	case strings.Contains(up, "NET"):
		return "Ağ aygıtı"
	case strings.Contains(up, "USB"):
		return "USB aygıtı"
	default:
		return "Aygıt"
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

// getDriveUsage returns the total and free bytes of a drive (for "how full is
// this USB" analysis — metadata only, never file contents).
func getDriveUsage(letter string) (total, free int64) {
	path := letter + `:\`
	var freeBytes, totalBytes, totalFree uint64
	p, _ := windows.UTF16PtrFromString(path)
	if err := windows.GetDiskFreeSpaceEx(p, &freeBytes, &totalBytes, &totalFree); err != nil {
		return 0, 0
	}
	return int64(totalBytes), int64(freeBytes)
}

// createMessageWindow creates a hidden TOP-LEVEL window. It must be top-level
// (not message-only) to receive the broadcast WM_DEVICECHANGE volume messages
// that signal USB drive insertion/removal. It is never shown and carries the
// tool-window ex-style so it never appears on screen or in the taskbar.
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
		wsExToolWindow, // ex-style: no taskbar button
		uintptr(unsafe.Pointer(className)),
		0,
		0,          // style: not WS_VISIBLE → stays hidden
		0, 0, 0, 0, // x, y, w, h
		0,         // top-level (no parent) — required to receive volume broadcasts
		0,         // menu
		hInstance, // instance
		0,         // lpParam
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
