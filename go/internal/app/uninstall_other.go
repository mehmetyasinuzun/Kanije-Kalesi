//go:build !windows

package app

import "fmt"

// systemUninstall is implemented only on Windows (Task Scheduler + detached
// PowerShell helper). On other platforms /kaldir reports it is unsupported.
func systemUninstall(_ removalPlan) error {
	return fmt.Errorf("kaldırma yalnızca Windows'ta destekleniyor")
}
