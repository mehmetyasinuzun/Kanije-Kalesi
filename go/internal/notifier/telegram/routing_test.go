package telegram

import "testing"

func TestExtractCommandStripsBotMention(t *testing.T) {
	cases := map[string]string{
		"/foto":              "/foto",
		"/foto dizustu":      "/foto",
		"/foto@kanijebot":    "/foto",
		"/foto@kanijebot ev": "/foto",
		"/seskayit@bot 30":   "/seskayit",
		"merhaba":            "",
		"/CihazLar@bot":      "/CihazLar",
	}
	for in, want := range cases {
		if got := extractCommand(in); got != want {
			t.Errorf("extractCommand(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRouteGroupCommand(t *testing.T) {
	b := &Bot{deviceLabel: "dizustu"}

	tests := []struct {
		name        string
		cmd, text   string
		wantHandled bool
		wantText    string
	}{
		{"broadcast cihazlar", "/cihazlar", "/cihazlar", true, "/cihazlar"},
		{"broadcast yardim", "/yardim", "/yardim", true, "/yardim"},
		{"target match", "/foto", "/foto dizustu", true, "/foto"},
		{"target match case-insensitive", "/foto", "/foto DIZUSTU", true, "/foto"},
		{"target match with args", "/seskayit", "/seskayit dizustu 30", true, "/seskayit 30"},
		{"target other device", "/foto", "/foto masaustu", false, "/foto masaustu"},
		{"no target in group", "/foto", "/foto", false, "/foto"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handled, stripped := b.routeGroupCommand(tc.cmd, tc.text)
			if handled != tc.wantHandled {
				t.Fatalf("handled = %v, want %v", handled, tc.wantHandled)
			}
			if handled && stripped != tc.wantText {
				t.Errorf("stripped = %q, want %q", stripped, tc.wantText)
			}
		})
	}
}
