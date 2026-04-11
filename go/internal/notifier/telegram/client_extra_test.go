package telegram

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestEscapeMarkdownAllSpecialChars(t *testing.T) {
	// Telegram MarkdownV2 spec: \ _ * [ ] ( ) ~ ` > # + - = | { } . !
	specials := `\_*[]()~` + "`" + `>#+-=|{}.!`
	got := EscapeMarkdown(specials)

	// Her özel karakter önünde \ olmalı
	for _, r := range specials {
		escaped := `\` + string(r)
		if !strings.Contains(got, escaped) {
			t.Fatalf("karakter %q kaçırılmamış: got=%q", string(r), got)
		}
	}
}

func TestEscapeMarkdownPreservesRegularText(t *testing.T) {
	in := "Kanije Kalesi 123 ABC çşğüöı"
	got := EscapeMarkdown(in)
	// Backslash eklenmeli sadece özel karakterlere — bu metinde yoklar
	if strings.Contains(got, `\K`) || strings.Contains(got, `\1`) {
		t.Fatalf("normal karakterler kacırılmamali: got=%q", got)
	}
}

func TestSafeTextOutputIsValidUTF8(t *testing.T) {
	// Çeşitli geçersiz sekanslar
	invalidInputs := [][]byte{
		{0xff, 0xfe},
		{0x80},
		{0xC0, 0x80},
		{'h', 'i', 0xff, '!'},
	}

	for _, b := range invalidInputs {
		got := SafeText(string(b))
		if !utf8.ValidString(got) {
			t.Fatalf("SafeText() ciktisi gecerli UTF-8 olmali: input=%x got=%q", b, got)
		}
	}
}

func TestSafeTextEmptyString(t *testing.T) {
	if got := SafeText(""); got != "" {
		t.Fatalf("bos string icin bos string donmeli: got=%q", got)
	}
}

func TestEscapeMarkdownEmptyString(t *testing.T) {
	if got := EscapeMarkdown(""); got != "" {
		t.Fatalf("bos string icin bos string donmeli: got=%q", got)
	}
}

func TestSafeTextTurkishCharacters(t *testing.T) {
	in := "Güvenlik İzleme Aracı — Türkçe Test: çşğüöı ÇŞĞÜÖİ"
	got := SafeText(in)
	if got != in {
		t.Fatalf("gecerli Turkce UTF-8 degismemeli:\n  in =%q\n  got=%q", in, got)
	}
}
