package telegram

import (
	"strings"
	"testing"

	"github.com/kanije-kalesi/kanije/internal/event"
)

func TestEscapeHTMLEscapesMarkupChars(t *testing.T) {
	cases := map[string]string{
		"a & b":     "a &amp; b",
		"<script>":  "&lt;script&gt;",
		"x > y < z": "x &gt; y &lt; z",
		"normal":    "normal",
		"":          "",
	}
	for in, want := range cases {
		if got := EscapeHTML(in); got != want {
			t.Errorf("EscapeHTML(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEscapeHTMLSinglePassNoDoubleEscape(t *testing.T) {
	// "&" must be escaped exactly once — a naive multi-replace would turn
	// "&amp;" into "&amp;amp;amp;". The Replacer guarantees a single pass.
	if got := EscapeHTML("&amp;"); got != "&amp;amp;" {
		t.Fatalf("tek geçişli kaçış bekleniyor: got=%q", got)
	}
}

func TestFormatEventEscapesAttackerControlledFields(t *testing.T) {
	ev := event.New(event.TypeUSBInserted, "test")
	ev.DeviceLabel = `USB <b>& "evil"</b>`
	ev.DevicePath = `D:\<x>`

	out := FormatEvent(ev)

	// Saldırgan kontrollü markup ham geçmemeli.
	if strings.Contains(out, "<b>&") {
		t.Fatalf("saldırgan markup kaçışsız geçti:\n%s", out)
	}
	// Beklenen entity'ler üretilmiş olmalı.
	if !strings.Contains(out, "&amp;") || !strings.Contains(out, "&lt;") {
		t.Fatalf("beklenen HTML entity'leri yok:\n%s", out)
	}
	// Bizim statik <b> iskelemiz korunmalı (yalnızca dinamik alan kaçışlanır).
	if !strings.Contains(out, "💾 Cihaz: <b>") {
		t.Fatalf("statik iskele korunmalı:\n%s", out)
	}
}

func TestFormatEventEscapesUsernameAndSSID(t *testing.T) {
	ev := event.New(event.TypeNetworkChanged, "test")
	ev.NetworkSSID = `Ev & Ofis <5G>`
	out := FormatEvent(ev)
	if strings.Contains(out, "<5G>") || !strings.Contains(out, "&lt;5G&gt;") {
		t.Fatalf("SSID kaçışlanmalı:\n%s", out)
	}
}

func TestSafeTextLeavesValidTextUnchanged(t *testing.T) {
	in := "geçerli Türkçe metin çğüöşı"
	if got := SafeText(in); got != in {
		t.Fatalf("geçerli metin değişmemeli: got=%q", got)
	}
}
