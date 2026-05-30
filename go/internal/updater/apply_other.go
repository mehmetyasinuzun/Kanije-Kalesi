//go:build !windows

package updater

import "context"

// CanSelfInstall reports whether the agent can replace its own binary in place.
// On Unix the hardened systemd service runs as an unprivileged user and cannot
// overwrite the root-owned binary, so self-install is unsupported.
func CanSelfInstall() bool { return false }

// SelfInstall is unsupported off Windows; the caller notifies the user to run
// the install script instead.
func (u *Updater) SelfInstall(_ context.Context, _ *Release) error {
	return ErrSelfInstallUnsupported
}
