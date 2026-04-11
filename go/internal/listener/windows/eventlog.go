//go:build windows

// Package windows implements security event listeners for Windows systems.
// This file handles the Windows Security Event Log via the modern EVT API
// (wevtapi.dll), available on Windows Vista / Server 2008 and later.
// No CGo is required — all calls go through golang.org/x/sys/windows LazyDLL.
package windows

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/kanije-kalesi/kanije/internal/event"
	"golang.org/x/sys/windows"
)

// ---- Windows EVT API bindings ----

var (
	wevtapi            = windows.NewLazySystemDLL("wevtapi.dll")
	procEvtSubscribe   = wevtapi.NewProc("EvtSubscribe")
	procEvtNext        = wevtapi.NewProc("EvtNext")
	procEvtRender      = wevtapi.NewProc("EvtRender")
	procEvtClose       = wevtapi.NewProc("EvtClose")
)

const (
	evtSubscribeToFutureEvents = 1
	evtRenderEventXml          = 1

	errorNoMoreItems  = 259
	errorInsufficientBuffer = 122
)

// ---- XML structures for event parsing ----

type winEvent struct {
	System    winSystem    `xml:"System"`
	EventData winEventData `xml:"EventData"`
}

type winSystem struct {
	EventID    int    `xml:"EventID"`
	TimeCreated struct {
		SystemTime string `xml:"SystemTime,attr"`
	} `xml:"TimeCreated"`
	Computer string `xml:"Computer"`
}

type winEventData struct {
	Data []winData `xml:"Data"`
}

type winData struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:",chardata"`
}

func (d winEventData) get(name string) string {
	for _, v := range d.Data {
		if v.Name == name {
			return v.Value
		}
	}
	return ""
}

// ---- EventLog Listener ----

// EventLogListener monitors the Windows Security Event Log in real time.
// It subscribes using EvtSubscribe with a signal event object (pull mode)
// to avoid Go ↔ Windows callback threading issues.
type EventLogListener struct {
	hostname string
	log      *slog.Logger
}

// NewEventLogListener creates a ready-to-use listener.
func NewEventLogListener(log *slog.Logger) *EventLogListener {
	h, _ := os.Hostname()
	return &EventLogListener{hostname: h, log: log}
}

func (l *EventLogListener) Name() string { return "EventLog" }

// Start subscribes to the Security Event Log and dispatches events until ctx
// is cancelled. The XPath query selects only the event IDs we care about.
func (l *EventLogListener) Start(ctx context.Context, bus *event.Bus) error {
	// Create a Windows manual-reset event object to signal us when new events arrive
	hSignal, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return fmt.Errorf("CreateEvent hatası: %w", err)
	}
	defer windows.CloseHandle(hSignal)

	// XPath query for the events we monitor
	query := `<QueryList>
	  <Query Id="0" Path="Security">
	    <Select Path="Security">
	      *[System[(EventID=4624 or EventID=4625 or EventID=4800 or EventID=4801 or EventID=4634)]]
	    </Select>
	  </Query>
	</QueryList>`

	channelPtr, _ := windows.UTF16PtrFromString("Security")
	queryPtr, _ := windows.UTF16PtrFromString(query)

	hSubscription, _, callErr := procEvtSubscribe.Call(
		0,                                  // session (local)
		uintptr(hSignal),                   // signal event
		uintptr(unsafe.Pointer(channelPtr)),
		uintptr(unsafe.Pointer(queryPtr)),
		0,                                  // bookmark (none)
		0,                                  // context
		0,                                  // callback (nil = signal mode)
		evtSubscribeToFutureEvents,
	)
	if hSubscription == 0 {
		return fmt.Errorf("EvtSubscribe başarısız: %v", callErr)
	}
	defer procEvtClose.Call(hSubscription)

	l.log.Info("Windows Event Log aboneliği başlatıldı",
		"sorgu", "4624/4625/4800/4801/4634")

	for {
		// Wait for signal or context cancellation (2-second poll for ctx check)
		waitResult, _ := windows.WaitForSingleObject(hSignal, 2000)
		switch waitResult {
		case windows.WAIT_OBJECT_0:
			// New events available — drain them
			l.drainEvents(hSubscription, bus)
		case uint32(windows.WAIT_TIMEOUT):
			// Check context periodically
		default:
			return fmt.Errorf("WaitForSingleObject beklenmedik sonuç: %d", waitResult)
		}

		if ctx.Err() != nil {
			return nil
		}
	}
}

// drainEvents reads all pending events from the subscription handle.
func (l *EventLogListener) drainEvents(hSub uintptr, bus *event.Bus) {
	var handles [64]uintptr
	var returned uint32

	for {
		ret, _, _ := procEvtNext.Call(
			hSub,
			64,
			uintptr(unsafe.Pointer(&handles[0])),
			1000, // timeout ms
			0,
			uintptr(unsafe.Pointer(&returned)),
		)
		if ret == 0 {
			break // No more events or error
		}

		for i := uint32(0); i < returned; i++ {
			ev, err := l.renderEvent(handles[i])
			procEvtClose.Call(handles[i])
			if err != nil {
				l.log.Debug("event render hatası", "err", err)
				continue
			}
			bus.Publish(ev)
		}
	}
}

// renderEvent converts a raw event handle to our Event struct.
func (l *EventLogListener) renderEvent(hEvent uintptr) (event.Event, error) {
	// First call: get required buffer size
	var bufferUsed, propertyCount uint32
	procEvtRender.Call(0, hEvent, evtRenderEventXml, 0, 0,
		uintptr(unsafe.Pointer(&bufferUsed)),
		uintptr(unsafe.Pointer(&propertyCount)))

	if bufferUsed == 0 {
		return event.Event{}, fmt.Errorf("boş event buffer")
	}

	// Second call: actually render
	buf := make([]uint16, bufferUsed/2+1)
	ret, _, callErr := procEvtRender.Call(
		0,
		hEvent,
		evtRenderEventXml,
		uintptr(bufferUsed),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufferUsed)),
		uintptr(unsafe.Pointer(&propertyCount)),
	)
	if ret == 0 {
		return event.Event{}, fmt.Errorf("EvtRender hatası: %v", callErr)
	}

	xmlStr := windows.UTF16ToString(buf)
	return l.parseEventXML(xmlStr)
}

// parseEventXML parses the rendered XML into our Event struct.
func (l *EventLogListener) parseEventXML(xmlStr string) (event.Event, error) {
	var we winEvent
	if err := xml.Unmarshal([]byte(xmlStr), &we); err != nil {
		return event.Event{}, fmt.Errorf("XML parse hatası: %w", err)
	}

	data := we.EventData

	// Parse timestamp
	ts := time.Now()
	if we.System.TimeCreated.SystemTime != "" {
		if parsed, err := time.Parse(time.RFC3339Nano, we.System.TimeCreated.SystemTime); err == nil {
			ts = parsed.Local()
		}
	}

	// Map logon type
	logonTypeStr := data.get("LogonType")
	var logonType event.LogonType
	if logonTypeStr != "" {
		var n int
		fmt.Sscan(logonTypeStr, &n)
		logonType = event.LogonType(n)
	}

	hostname := we.System.Computer
	if hostname == "" {
		hostname = l.hostname
	}

	username := data.get("TargetUserName")
	domain   := data.get("TargetDomainName")
	sourceIP := data.get("IpAddress")
	if sourceIP == "-" || sourceIP == "::1" || sourceIP == "127.0.0.1" {
		sourceIP = ""
	}

	switch we.System.EventID {
	case 4624:
		// Skip machine accounts (ending in $) and system accounts
		if len(username) > 0 && username[len(username)-1] == '$' {
			return event.Event{}, fmt.Errorf("makine hesabı oturumu atlanıyor")
		}
		if logonType == event.LogonService || logonType == event.LogonBatch {
			return event.Event{}, fmt.Errorf("servis/batch oturumu atlanıyor")
		}
		ev := event.New(event.TypeLoginSuccess, "EventLog")
		ev.Timestamp = ts
		ev.Hostname = hostname
		ev.Username = username
		ev.Domain   = domain
		ev.SourceIP = sourceIP
		ev.LogonType = logonType
		return ev, nil

	case 4625:
		ev := event.New(event.TypeLoginFailed, "EventLog")
		ev.Timestamp = ts
		ev.Hostname = hostname
		ev.Username = username
		ev.Domain   = domain
		ev.SourceIP = sourceIP
		ev.LogonType = logonType
		return ev, nil

	case 4800:
		ev := event.New(event.TypeScreenLock, "EventLog")
		ev.Timestamp = ts
		ev.Hostname = hostname
		ev.Username = data.get("TargetUserName")
		return ev, nil

	case 4801:
		ev := event.New(event.TypeScreenUnlock, "EventLog")
		ev.Timestamp = ts
		ev.Hostname = hostname
		ev.Username = data.get("TargetUserName")
		return ev, nil

	case 4634:
		ev := event.New(event.TypeLogoff, "EventLog")
		ev.Timestamp = ts
		ev.Hostname = hostname
		ev.Username = username
		return ev, nil

	default:
		return event.Event{}, fmt.Errorf("bilinmeyen event ID: %d", we.System.EventID)
	}
}

// ---- syscall helpers ----

func utf16PtrFromString(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}
