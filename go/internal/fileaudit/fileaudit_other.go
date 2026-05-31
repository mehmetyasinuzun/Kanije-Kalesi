//go:build !windows

package fileaudit

import "fmt"

var errUnsupported = fmt.Errorf("dosya erişim denetimi yalnızca Windows'ta destekleniyor")

// EnableAndAudit is Windows-only (SACL + audit policy).
func EnableAndAudit(_ string) error { return errUnsupported }

// DisableAudit is Windows-only.
func DisableAudit(_ string) error { return errUnsupported }

// RecentAccess is Windows-only (Security event log 4663).
func RecentAccess(_ []string, _ int) ([]AccessEvent, error) { return nil, errUnsupported }
