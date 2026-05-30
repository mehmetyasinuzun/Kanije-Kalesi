package totp

import (
	"testing"
	"time"
)

func TestKnownRFCVector(t *testing.T) {
	// RFC 6238 test secret "12345678901234567890" → base32 below.
	// At T=59s the counter is 1; the 8-digit value is 94287082, so the
	// 6-digit truncation is 287082.
	key, err := decodeSecret("GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ")
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got := hotp(key, 1); got != "287082" {
		t.Fatalf("RFC 6238 vektörü uyuşmuyor: got=%s want=287082", got)
	}
}

func TestValidateRoundTrip(t *testing.T) {
	secret := GenerateSecret()
	key, _ := decodeSecret(secret)
	code := hotp(key, uint64(time.Now().Unix()/period))
	if !Validate(secret, code) {
		t.Fatal("o an üretilen kod doğrulanmalı")
	}
}

func TestValidateRejectsBad(t *testing.T) {
	secret := GenerateSecret()
	if Validate("", "123456") {
		t.Error("boş secret reddedilmeli")
	}
	if Validate("!!!notbase32!!!", "123456") {
		t.Error("geçersiz secret reddedilmeli")
	}
	if Validate(secret, "12345") {
		t.Error("yanlış uzunluktaki kod reddedilmeli")
	}
	if Validate(secret, "") {
		t.Error("boş kod reddedilmeli")
	}
}
