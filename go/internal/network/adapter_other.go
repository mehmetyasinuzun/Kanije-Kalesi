//go:build !windows

package network

// activeAdapter is Windows-only (IP Helper API). On other platforms the medium
// is inferred from the interface name in monitor.go instead.
func activeAdapter() (friendly, medium string) {
	return "", ""
}
