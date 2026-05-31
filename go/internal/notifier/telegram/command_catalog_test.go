package telegram

import (
	"strings"
	"testing"
)

// TestFormatHelpFromCatalog checks /yardim is generated and covers key commands.
func TestFormatHelpFromCatalog(t *testing.T) {
	h := FormatHelp()
	for _, want := range []string{"/tetikkamera", "/tetikses", "/koruma", "/rehber", "Tetik Modları", "Komutlar"} {
		if !strings.Contains(h, want) {
			t.Errorf("/yardim çıktısında %q yok", want)
		}
	}
}

// TestMenuCommandsConsistent ensures the "/" menu is well-formed and in sync with
// the catalog (no leading slash, no blank labels, reasonable size).
func TestMenuCommandsConsistent(t *testing.T) {
	m := menuCommands()
	if len(m) < 15 {
		t.Fatalf("menü beklenenden kısa: %d komut", len(m))
	}
	seen := map[string]bool{}
	for _, c := range m {
		if strings.HasPrefix(c.Command, "/") {
			t.Errorf("menü komutu '/' ile başlamamalı: %q", c.Command)
		}
		if c.Command == "" || c.Description == "" {
			t.Errorf("eksik menü girdisi: %+v", c)
		}
		if seen[c.Command] {
			t.Errorf("menüde mükerrer komut: %q", c.Command)
		}
		seen[c.Command] = true
	}
}

// TestCatalogCommandsHaveLeadingSlash guards the catalog format.
func TestCatalogCommandsHaveLeadingSlash(t *testing.T) {
	for _, cat := range commandCatalog {
		for _, it := range cat.items {
			if !strings.HasPrefix(it.cmd, "/") {
				t.Errorf("katalog komutu '/' ile başlamalı: %q (%s)", it.cmd, cat.header)
			}
			if it.help == "" {
				t.Errorf("katalog komutunun açıklaması boş: %q", it.cmd)
			}
		}
	}
}
