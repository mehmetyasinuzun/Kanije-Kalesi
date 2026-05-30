// Package webhook fans out security events to outgoing HTTP webhook
// destinations (Discord embeds or a generic JSON payload). Delivery is
// best-effort — failures are logged and never block event processing.
package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kanije-kalesi/kanije/internal/event"
)

// Target is one outgoing webhook destination.
type Target struct {
	Name   string
	URL    string
	Format string // "discord" | "json"
}

// Sender delivers events to a set of targets.
type Sender struct {
	targets []Target
	http    *http.Client
	log     *slog.Logger
}

// New builds a Sender. Targets without a URL are dropped.
func New(targets []Target, log *slog.Logger) *Sender {
	valid := make([]Target, 0, len(targets))
	for _, t := range targets {
		if strings.HasPrefix(t.URL, "http") {
			valid = append(valid, t)
		}
	}
	return &Sender{
		targets: valid,
		http:    &http.Client{Timeout: 10 * time.Second},
		log:     log,
	}
}

// Enabled reports whether any valid targets are configured.
func (s *Sender) Enabled() bool { return len(s.targets) > 0 }

// Send delivers ev to every target. Each POST is independent and best-effort.
func (s *Sender) Send(ctx context.Context, ev event.Event) {
	for _, t := range s.targets {
		body := buildPayload(t.Format, ev)
		if err := s.post(ctx, t.URL, body); err != nil {
			s.log.Warn("webhook gönderilemedi", "hedef", t.Name, "err", err)
		}
	}
}

func (s *Sender) post(ctx context.Context, url string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

// buildPayload returns the JSON body for the given format ("discord" or "json").
func buildPayload(format string, ev event.Event) []byte {
	if strings.EqualFold(format, "discord") {
		b, _ := json.Marshal(discordPayload(ev))
		return b
	}
	b, _ := json.Marshal(jsonPayload(ev))
	return b
}

// ---- Discord ----

type discordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type discordEmbed struct {
	Title     string         `json:"title"`
	Color     int            `json:"color"`
	Fields    []discordField `json:"fields,omitempty"`
	Timestamp string         `json:"timestamp"`
}

type discordMessage struct {
	Embeds []discordEmbed `json:"embeds"`
}

func discordPayload(ev event.Event) discordMessage {
	emb := discordEmbed{
		Title:     ev.Type.Emoji() + " " + ev.Type.Label(),
		Color:     severityColor(ev.Severity),
		Timestamp: ev.Timestamp.UTC().Format(time.RFC3339),
	}
	add := func(name, val string) {
		if val != "" {
			emb.Fields = append(emb.Fields, discordField{Name: name, Value: val, Inline: true})
		}
	}
	add("Bilgisayar", ev.Hostname)
	add("Kullanıcı", ev.Username)
	add("IP", ev.SourceIP)
	if loc := ev.Extra["📍 Konum"]; loc != "" {
		add("Konum", loc)
	}
	add("Cihaz", ev.DeviceLabel)
	return discordMessage{Embeds: []discordEmbed{emb}}
}

func severityColor(s event.Severity) int {
	switch s {
	case event.SeverityCritical:
		return 0xE01E1E
	case event.SeverityAlert:
		return 0xFF6B00
	case event.SeverityWarning:
		return 0xF1C40F
	default:
		return 0x3498DB
	}
}

// ---- Generic JSON ----

func jsonPayload(ev event.Event) map[string]any {
	m := map[string]any{
		"type":      string(ev.Type),
		"label":     ev.Type.Label(),
		"severity":  ev.Severity.String(),
		"timestamp": ev.Timestamp.UTC().Format(time.RFC3339),
	}
	put := func(k, v string) {
		if v != "" {
			m[k] = v
		}
	}
	put("hostname", ev.Hostname)
	put("username", ev.Username)
	put("source_ip", ev.SourceIP)
	put("device", ev.DeviceLabel)
	if len(ev.Extra) > 0 {
		m["extra"] = ev.Extra
	}
	return m
}
