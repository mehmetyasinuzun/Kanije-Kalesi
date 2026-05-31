//go:build windows

package fileaudit

import "testing"

func TestAccessLabel(t *testing.T) {
	cases := map[string]string{
		"0x1":     "okuma/kopya",
		"0x2":     "yazma",
		"0x10000": "silme",
		"0x10001": "silme+okuma/kopya",
		"0x3":     "yazma+okuma/kopya",
		"":        "erişim", // unparseable → generic
		"garbage": "erişim",
	}
	for in, want := range cases {
		if got := accessLabel(in); got != want {
			t.Errorf("accessLabel(%q) = %q, beklenen %q", in, got, want)
		}
	}
}
