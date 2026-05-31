package telegram

import "testing"

func TestLevenshtein(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"kitten", "sitting", 3},
		{"abc", "abc", 0},
		{"", "abc", 3},
		{"abc", "", 3},
		{"/foto", "/foton", 1},
	}
	for _, c := range cases {
		if got := levenshtein(c.a, c.b); got != c.want {
			t.Errorf("levenshtein(%q,%q) = %d, beklenen %d", c.a, c.b, got, c.want)
		}
	}
}

func TestSuggestCommand(t *testing.T) {
	cases := map[string]string{
		"/foton":      "/foto",   // one extra letter
		"/ekrn":       "/ekran",  // one missing letter
		"/statu":      "/status", // one missing letter
		"/xyzqwerty":  "",        // nothing close
		"merhaba dlk": "",        // not even close to a command
	}
	for in, want := range cases {
		if got := suggestCommand(in); got != want {
			t.Errorf("suggestCommand(%q) = %q, beklenen %q", in, got, want)
		}
	}
}
