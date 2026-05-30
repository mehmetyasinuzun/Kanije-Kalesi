package secret

import "testing"

func TestProtectUnprotectRoundTrip(t *testing.T) {
	const token = "1234567890:ABCdef_test-token"

	stored := Protect(token)
	got, err := Unprotect(stored)
	if err != nil {
		t.Fatalf("Unprotect hata: %v", err)
	}
	if got != token {
		t.Fatalf("round-trip bozuk: got=%q want=%q", got, token)
	}
}

func TestUnprotectPlaintextPassthrough(t *testing.T) {
	// Önek taşımayan (ör. ortam değişkeninden gelen) değer aynen döner.
	got, err := Unprotect("plain-token")
	if err != nil || got != "plain-token" {
		t.Fatalf("passthrough bozuk: got=%q err=%v", got, err)
	}
	if IsProtected("plain-token") {
		t.Error("öneksiz değer korumalı sayılmamalı")
	}
}

func TestProtectEmpty(t *testing.T) {
	if Protect("") != "" {
		t.Error("boş değer boş kalmalı")
	}
}
