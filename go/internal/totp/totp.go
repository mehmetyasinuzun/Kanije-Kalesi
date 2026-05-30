// Package totp implements RFC 6238 time-based one-time passwords (HMAC-SHA1,
// 6 digits, 30-second period) for optional two-factor confirmation of dangerous
// commands. Secrets are standard base32 so any authenticator app works.
package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	period = 30
	digits = 6
)

var b32 = base32.StdEncoding.WithPadding(base32.NoPadding)

// GenerateSecret returns a new random base32 secret suitable for enrollment.
func GenerateSecret() string {
	b := make([]byte, 20)
	_, _ = rand.Read(b)
	return b32.EncodeToString(b)
}

// URI builds an otpauth:// URI that authenticator apps import (often via QR).
func URI(secret, issuer, account string) string {
	label := url.PathEscape(issuer + ":" + account)
	q := url.Values{}
	q.Set("secret", secret)
	q.Set("issuer", issuer)
	q.Set("digits", "6")
	q.Set("period", "30")
	q.Set("algorithm", "SHA1")
	return "otpauth://totp/" + label + "?" + q.Encode()
}

// Validate reports whether code is a valid TOTP for secret at the current time,
// tolerating ±1 step (≈30s) for clock skew between the device and the server.
func Validate(secret, code string) bool {
	key, err := decodeSecret(secret)
	if err != nil || len(key) == 0 {
		return false
	}
	code = strings.TrimSpace(code)
	if len(code) != digits {
		return false
	}

	counter := uint64(time.Now().Unix() / period)
	for _, c := range []uint64{counter - 1, counter, counter + 1} {
		if hmac.Equal([]byte(hotp(key, c)), []byte(code)) {
			return true
		}
	}
	return false
}

func decodeSecret(secret string) ([]byte, error) {
	s := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(secret), " ", ""))
	return b32.DecodeString(s)
}

// hotp computes the RFC 4226 HOTP value for a counter.
func hotp(key []byte, counter uint64) string {
	var buf [8]byte
	for i := 7; i >= 0; i-- {
		buf[i] = byte(counter & 0xff)
		counter >>= 8
	}
	mac := hmac.New(sha1.New, key)
	mac.Write(buf[:])
	sum := mac.Sum(nil)

	off := sum[len(sum)-1] & 0x0f
	val := (uint32(sum[off]&0x7f) << 24) |
		(uint32(sum[off+1]) << 16) |
		(uint32(sum[off+2]) << 8) |
		uint32(sum[off+3])

	mod := uint32(1)
	for i := 0; i < digits; i++ {
		mod *= 10
	}
	return fmt.Sprintf("%0*d", digits, val%mod)
}
