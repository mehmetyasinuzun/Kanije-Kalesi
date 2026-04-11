package telegram

import (
	"strings"
	"testing"
)

func TestEscapeMarkdownEscapesSpecialCharacters(t *testing.T) {
	in := "Hello_*[x](y)!"
	got := EscapeMarkdown(in)
	want := "Hello\\_\\*\\[x\\]\\(y\\)\\!"

	if got != want {
		t.Fatalf("EscapeMarkdown() = %q, want %q", got, want)
	}
}

func TestSafeTextReplacesInvalidUTF8(t *testing.T) {
	in := string([]byte{0xff, 'A'})
	got := SafeText(in)

	if strings.ContainsRune(got, '\uFFFD') {
		t.Fatalf("SafeText gecersiz rune birakmamali: %q", got)
	}
	if got != "?A" {
		t.Fatalf("SafeText() = %q, want %q", got, "?A")
	}
}

func TestSafeTextValidUTF8Unchanged(t *testing.T) {
	in := "Kanije Kalesi"
	if got := SafeText(in); got != in {
		t.Fatalf("SafeText gecerli metni degistirmemeli: got=%q", got)
	}
}
