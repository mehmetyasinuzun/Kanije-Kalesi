//go:build !windows

package defender

// GetStatus is implemented only on Windows (Microsoft Defender). Elsewhere it
// reports the AV as unavailable.
func GetStatus() Status {
	return Status{}
}
