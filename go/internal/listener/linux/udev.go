//go:build linux

package linux

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/kanije-kalesi/kanije/internal/event"
)

// UdevListener monitors USB device insertion and removal via the kernel's
// uevent netlink socket. This requires no external tools or CGo —
// it uses raw AF_NETLINK sockets from the syscall package.
type UdevListener struct {
	hostname string
	log      *slog.Logger
}

func NewUdevListener(log *slog.Logger) *UdevListener {
	h, _ := os.Hostname()
	return &UdevListener{hostname: h, log: log}
}

func (l *UdevListener) Name() string { return "UdevMonitor" }

// Start opens a NETLINK_KOBJECT_UEVENT socket and parses uevent messages.
// Messages look like:
//
//	add@/devices/pci0000:00/0000:00:14.0/usb1/1-1\0
//	ACTION=add\0DEVTYPE=usb_device\0SUBSYSTEM=usb\0...
func (l *UdevListener) Start(ctx context.Context, bus *event.Bus) error {
	fd, err := syscall.Socket(
		syscall.AF_NETLINK,
		syscall.SOCK_RAW,
		syscall.NETLINK_KOBJECT_UEVENT,
	)
	if err != nil {
		return fmt.Errorf("netlink socket açılamadı (root gerekli?): %w", err)
	}
	defer syscall.Close(fd)

	addr := &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Groups: 1, // UDEV_MONITOR_KERNEL
	}
	if err := syscall.Bind(fd, addr); err != nil {
		return fmt.Errorf("netlink bind hatası: %w", err)
	}

	l.log.Info("Udev USB izleme başlatıldı")

	// Context cancellation: close fd in a separate goroutine
	go func() {
		<-ctx.Done()
		syscall.Close(fd)
	}()

	buf := make([]byte, 4096)
	for {
		n, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("recvfrom hatası: %w", err)
		}

		ev, ok := l.parseUevent(buf[:n])
		if !ok {
			continue
		}
		bus.Publish(ev)
	}
}

// parseUevent converts a raw uevent message to an Event.
// The message is a sequence of NUL-terminated key=value strings.
func (l *UdevListener) parseUevent(data []byte) (event.Event, bool) {
	fields := splitUevent(data)

	action    := fields["ACTION"]
	subsystem := fields["SUBSYSTEM"]
	devtype   := fields["DEVTYPE"]
	devname   := fields["DEVNAME"]
	idVendor  := fields["ID_VENDOR"]
	idModel   := fields["ID_MODEL"]
	idFS      := fields["ID_FS_TYPE"]
	idLabel   := fields["ID_FS_LABEL"]
	idSize    := fields["ID_PART_ENTRY_SIZE"]

	// Only care about block devices of type "disk" or "partition"
	if subsystem != "block" || (devtype != "disk" && devtype != "partition") {
		return event.Event{}, false
	}

	// Skip non-USB devices
	if !strings.Contains(fields["DEVPATH"], "usb") {
		return event.Event{}, false
	}

	label := idLabel
	if label == "" {
		label = idModel
	}
	if label == "" && idVendor != "" {
		label = idVendor
	}

	switch action {
	case "add":
		ev := event.New(event.TypeUSBInserted, "UdevMonitor")
		ev.Hostname    = l.hostname
		ev.DevicePath  = "/dev/" + devname
		ev.DeviceName  = idModel
		ev.DeviceLabel = label
		ev.DeviceFS    = idFS
		if idSize != "" {
			var sectors int64
			fmt.Sscan(idSize, &sectors)
			ev.DeviceSize = sectors * 512
		}
		return ev, true

	case "remove":
		ev := event.New(event.TypeUSBRemoved, "UdevMonitor")
		ev.Hostname    = l.hostname
		ev.DevicePath  = "/dev/" + devname
		ev.DeviceName  = label
		return ev, true
	}

	return event.Event{}, false
}

// splitUevent parses a uevent message into a key→value map.
// The message consists of NUL-separated records; the first record is the
// header (e.g. "add@/path"), subsequent records are "KEY=VALUE" pairs.
func splitUevent(data []byte) map[string]string {
	fields := make(map[string]string, 16)

	parts := splitNUL(data)
	for i, p := range parts {
		if i == 0 {
			// Header: "action@devpath"
			if at := strings.IndexByte(p, '@'); at >= 0 {
				fields["ACTION"] = p[:at]
				fields["DEVPATH"] = p[at+1:]
			}
			continue
		}
		if eq := strings.IndexByte(p, '='); eq >= 0 {
			fields[p[:eq]] = p[eq+1:]
		}
	}
	return fields
}

func splitNUL(data []byte) []string {
	var result []string
	start := 0
	for i, b := range data {
		if b == 0 {
			if i > start {
				// Ensure valid UTF-8 by converting via unsafe only when needed
				s := string(data[start:i])
				result = append(result, s)
			}
			start = i + 1
		}
	}
	if start < len(data) {
		result = append(result, string(data[start:]))
	}
	return result
}

// ---- helpers to silence unused import warning ----
var _ = unsafe.Sizeof
