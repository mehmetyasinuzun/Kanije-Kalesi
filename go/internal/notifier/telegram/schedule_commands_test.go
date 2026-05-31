package telegram

import "testing"

func TestParseInterval(t *testing.T) {
	ok := map[string]int{
		"60sn": 60, "30dk": 1800, "2sa": 7200, "1g": 86400,
		"90s": 90, "5m": 300, "1h": 3600, "1d": 86400,
	}
	for in, want := range ok {
		got, err := parseInterval(in)
		if err != nil || got != want {
			t.Errorf("parseInterval(%q) = %d, err %v; beklenen %d", in, got, err, want)
		}
	}

	// 45sn < 60 floor, and malformed inputs must all error.
	for _, bad := range []string{"45sn", "abc", "30", "sa", "-5dk", ""} {
		if _, err := parseInterval(bad); err == nil {
			t.Errorf("parseInterval(%q) hata vermeliydi", bad)
		}
	}
}

func TestFormatInterval(t *testing.T) {
	cases := map[int]string{60: "1 dk", 1800: "30 dk", 3600: "1 saat", 7200: "2 saat", 86400: "1 gün"}
	for sec, want := range cases {
		if got := formatInterval(sec); got != want {
			t.Errorf("formatInterval(%d) = %q, beklenen %q", sec, got, want)
		}
	}
}
