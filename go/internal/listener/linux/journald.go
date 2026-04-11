//go:build linux

// Package linux implements security event listeners for Linux systems.
package linux

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kanije-kalesi/sentinel/internal/event"
)

// JournaldListener monitors the systemd journal for security-relevant entries.
// It pipes `journalctl --follow --output=json` and parses JSON records.
// No CGo is needed — this is a pure subprocess approach.
type JournaldListener struct {
	hostname string
	log      *slog.Logger
}

func NewJournaldListener(log *slog.Logger) *JournaldListener {
	h, _ := os.Hostname()
	return &JournaldListener{hostname: h, log: log}
}

func (l *JournaldListener) Name() string { return "Journald" }

// journalEntry represents a subset of journald's JSON output.
type journalEntry struct {
	Message      string `json:"MESSAGE"`
	Comm         string `json:"_COMM"`
	Identifier   string `json:"SYSLOG_IDENTIFIER"`
	Priority     string `json:"PRIORITY"`
	UID          string `json:"_UID"`
	PID          string `json:"_PID"`
	RealtimeUsec string `json:"__REALTIME_TIMESTAMP"`
}

func (l *JournaldListener) Start(ctx context.Context, bus *event.Bus) error {
	cmd := exec.CommandContext(ctx,
		"journalctl",
		"--follow",
		"--output=json",
		"--since=now",
		// Only monitor auth-relevant units/identifiers
		"--identifier=sshd",
		"--identifier=sudo",
		"--identifier=systemd-logind",
		"--identifier=login",
		"--identifier=su",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("journalctl stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("journalctl başlatılamadı: %w", err)
	}

	l.log.Info("Journald izleme başlatıldı")

	go func() {
		<-ctx.Done()
		cmd.Process.Kill()
	}()

	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 1024*64), 1024*64)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "{") {
			continue
		}

		var entry journalEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			l.log.Debug("journal parse hatası", "err", err)
			continue
		}

		ev, ok := l.classifyEntry(entry)
		if !ok {
			continue
		}
		bus.Publish(ev)
	}

	if err := cmd.Wait(); err != nil && ctx.Err() == nil {
		return fmt.Errorf("journalctl beklenmedik şekilde sonlandı: %w", err)
	}
	return nil
}

// classifyEntry converts a journal entry to a security event if relevant.
func (l *JournaldListener) classifyEntry(e journalEntry) (event.Event, bool) {
	msg := strings.ToLower(e.Message)
	id  := strings.ToLower(e.Identifier)

	ts := parseJournalTime(e.RealtimeUsec)

	switch {
	// SSH login success: "Accepted password/publickey for user from IP port ..."
	case id == "sshd" && (strings.Contains(msg, "accepted password") || strings.Contains(msg, "accepted publickey")):
		user, ip := parseSSHAccepted(e.Message)
		ev := event.New(event.TypeLoginSuccess, "Journald")
		ev.Timestamp = ts
		ev.Hostname = l.hostname
		ev.Username = user
		ev.SourceIP = ip
		return ev, true

	// SSH login failed: "Failed password for user from IP ..."
	case id == "sshd" && strings.Contains(msg, "failed password"):
		user, ip := parseSSHFailed(e.Message)
		ev := event.New(event.TypeLoginFailed, "Journald")
		ev.Timestamp = ts
		ev.Hostname = l.hostname
		ev.Username = user
		ev.SourceIP = ip
		return ev, true

	// Invalid user
	case id == "sshd" && strings.Contains(msg, "invalid user"):
		user, ip := parseSSHInvalidUser(e.Message)
		ev := event.New(event.TypeLoginFailed, "Journald")
		ev.Timestamp = ts
		ev.Hostname = l.hostname
		ev.Username = user
		ev.SourceIP = ip
		ev.Extra = map[string]string{"neden": "geçersiz kullanıcı"}
		return ev, true

	// systemd-logind: screen lock
	case id == "systemd-logind" && strings.Contains(msg, "locked"):
		ev := event.New(event.TypeScreenLock, "Journald")
		ev.Timestamp = ts
		ev.Hostname = l.hostname
		return ev, true

	// systemd-logind: session opened (local login)
	case id == "systemd-logind" && strings.Contains(msg, "new session") && strings.Contains(msg, "seat"):
		user := parseLogindUser(e.Message)
		ev := event.New(event.TypeLoginSuccess, "Journald")
		ev.Timestamp = ts
		ev.Hostname = l.hostname
		ev.Username = user
		ev.LogonType = event.LogonInteractive
		return ev, true
	}

	return event.Event{}, false
}

// ---- SSH log parsers ----
// All parsing is done on UTF-8 strings — no encoding conversion needed.

func parseSSHAccepted(msg string) (user, ip string) {
	// "Accepted password for alice from 192.168.1.1 port 22 ssh2"
	parts := strings.Fields(msg)
	for i, p := range parts {
		if p == "for" && i+1 < len(parts) {
			user = parts[i+1]
		}
		if p == "from" && i+1 < len(parts) {
			ip = parts[i+1]
		}
	}
	return
}

func parseSSHFailed(msg string) (user, ip string) {
	// "Failed password for alice from 192.168.1.1 port 22 ssh2"
	// "Failed password for invalid user bob from 10.0.0.1 port ..."
	parts := strings.Fields(msg)
	for i, p := range parts {
		if p == "for" && i+1 < len(parts) {
			next := parts[i+1]
			if next == "invalid" && i+2 < len(parts) && parts[i+2] == "user" && i+3 < len(parts) {
				user = parts[i+3]
			} else {
				user = next
			}
		}
		if p == "from" && i+1 < len(parts) {
			ip = parts[i+1]
		}
	}
	return
}

func parseSSHInvalidUser(msg string) (user, ip string) {
	// "Invalid user alice from 192.168.1.1 port 22"
	parts := strings.Fields(msg)
	for i, p := range parts {
		if p == "user" && i+1 < len(parts) {
			user = parts[i+1]
		}
		if p == "from" && i+1 < len(parts) {
			ip = parts[i+1]
		}
	}
	return
}

func parseLogindUser(msg string) string {
	// "New session 3 of user alice."
	parts := strings.Fields(msg)
	for i, p := range parts {
		if p == "user" && i+1 < len(parts) {
			u := parts[i+1]
			return strings.TrimSuffix(u, ".")
		}
	}
	return ""
}

// parseJournalTime converts journald's microsecond Unix timestamp to time.Time.
// The RealtimeUsec field is a string of microseconds since the Unix epoch.
func parseJournalTime(usec string) time.Time {
	if usec == "" {
		return time.Now()
	}
	var us int64
	fmt.Sscan(usec, &us)
	return time.Unix(us/1_000_000, (us%1_000_000)*1_000)
}
