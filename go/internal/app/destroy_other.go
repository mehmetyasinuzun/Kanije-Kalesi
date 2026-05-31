//go:build !windows

package app

import "fmt"

// systemDestroy is implemented only on Windows (secure wipe + systemreset factory
// reset). On other platforms /imha reports it is unsupported.
func systemDestroy(_ destroyPlan) error {
	return fmt.Errorf("imha yalnızca Windows'ta destekleniyor")
}
