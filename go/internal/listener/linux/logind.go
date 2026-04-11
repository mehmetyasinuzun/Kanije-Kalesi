//go:build linux

package linux

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kanije-kalesi/sentinel/internal/event"
)

// LogindListener monitors screen lock/unlock and sleep/wake events from
// systemd-logind via D-Bus. It uses the `dbus-monitor` subprocess to avoid
// CGo — no additional Go packages required.
type LogindListener struct {
	hostname string
	log      *slog.Logger
}

func NewLogindListener(log *slog.Logger) *LogindListener {
	h, _ := os.Hostname()
	return &LogindListener{hostname: h, log: log}
}

func (l *LogindListener) Name() string { return "LogindMonitor" }

// Start subscribes to relevant D-Bus signals via dbus-monitor and parses the
// text output. The signals we monitor:
//
//	org.freedesktop.login1.Session → Lock / Unlock
//	org.freedesktop.login1.Manager → PrepareForSleep
func (l *LogindListener) Start(ctx context.Context, bus *event.Bus) error {
	cmd := exec.CommandContext(ctx,
		"dbus-monitor",
		"--system",
		"type='signal',interface='org.freedesktop.login1.Session'",
		"type='signal',interface='org.freedesktop.login1.Manager',member='PrepareForSleep'",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("dbus-monitor stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		// dbus-monitor may not be installed — log and return gracefully
		l.log.Warn("dbus-monitor başlatılamadı, logind izleme devre dışı",
			"err", err,
			"öneri", "sudo apt install dbus-x11")
		return nil
	}

	l.log.Info("Logind izleme başlatıldı (dbus-monitor)")

	go func() {
		<-ctx.Done()
		cmd.Process.Kill()
	}()

	parser := newDbusTextParser(l.log)
	if err := parser.Run(ctx, stdout, func(ev event.Event) {
		ev.Hostname = l.hostname
		bus.Publish(ev)
	}); err != nil && ctx.Err() == nil {
		return fmt.Errorf("dbus-monitor beklenmedik sonlanma: %w", err)
	}

	cmd.Wait()
	return nil
}

// dbusTextParser parses the line-oriented text output of dbus-monitor.
// A signal block looks like:
//
//	signal time=1234.567 sender=:1.10 -> destination=(null destination) serial=5
//	   path=/org/freedesktop/login1/session/_31; interface=org.freedesktop.login1.Session; member=Lock
type dbusTextParser struct {
	log *slog.Logger
}

func newDbusTextParser(log *slog.Logger) *dbusTextParser {
	return &dbusTextParser{log: log}
}

func (p *dbusTextParser) Run(ctx context.Context, r interface{ ReadString(byte) (string, error) }, emit func(event.Event)) error {
	// We'll use a simple line reader approach
	type lineReader interface {
		ReadString(byte) (string, error)
	}
	lr, ok := r.(lineReader)
	if !ok {
		return fmt.Errorf("geçersiz reader")
	}

	var currentMember string

	for {
		if ctx.Err() != nil {
			return nil
		}

		line, err := lr.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)

		// Parse signal header: extract member
		if strings.HasPrefix(line, "signal") {
			currentMember = ""
			continue
		}

		// Parse signal body line: "path=...; interface=...; member=..."
		if strings.Contains(line, "member=") {
			parts := strings.Split(line, ";")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "member=") {
					currentMember = strings.TrimPrefix(part, "member=")
				}
			}

			switch currentMember {
			case "Lock":
				ev := event.New(event.TypeScreenLock, "LogindMonitor")
				ev.Timestamp = time.Now()
				emit(ev)
			case "Unlock":
				ev := event.New(event.TypeScreenUnlock, "LogindMonitor")
				ev.Timestamp = time.Now()
				emit(ev)
			}
			continue
		}

		// PrepareForSleep boolean argument
		if currentMember == "PrepareForSleep" {
			if strings.TrimSpace(line) == "boolean true" {
				ev := event.New(event.TypeSystemSleep, "LogindMonitor")
				ev.Timestamp = time.Now()
				emit(ev)
			} else if strings.TrimSpace(line) == "boolean false" {
				ev := event.New(event.TypeSystemWake, "LogindMonitor")
				ev.Timestamp = time.Now()
				ev.WakeType = "manuel"
				emit(ev)
			}
		}
	}
}
