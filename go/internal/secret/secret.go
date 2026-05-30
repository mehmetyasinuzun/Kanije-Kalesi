// Package secret protects sensitive config values (the Telegram bot token) at
// rest. On Windows it uses DPAPI (CryptProtectData), binding the ciphertext to
// the current user account so a stolen config.toml is useless on another machine
// or account. On other platforms it transparently stores the value as-is and
// relies on the 0600 file permissions.
package secret

import "strings"

const dpapiPrefix = "dpapi:"

// Protect returns the at-rest representation of plaintext. On Windows this is a
// DPAPI-encrypted, "dpapi:"-prefixed string; elsewhere it is the value itself.
// If encryption fails for any reason, the plaintext is stored as a safe fallback.
func Protect(plaintext string) string {
	if plaintext == "" {
		return ""
	}
	enc, err := dpapiProtect(plaintext)
	if err != nil {
		return plaintext
	}
	return dpapiPrefix + enc
}

// Unprotect reverses Protect. A value without the "dpapi:" prefix is returned
// unchanged (plaintext or env-provided). A prefixed value is decrypted.
func Unprotect(stored string) (string, error) {
	if !strings.HasPrefix(stored, dpapiPrefix) {
		return stored, nil
	}
	return dpapiUnprotect(strings.TrimPrefix(stored, dpapiPrefix))
}

// Supported reports whether real at-rest encryption is available on this platform.
func Supported() bool { return dpapiSupported() }

// IsProtected reports whether stored is an encrypted value.
func IsProtected(stored string) bool { return strings.HasPrefix(stored, dpapiPrefix) }
