//go:build windows

package network

import "testing"

func TestMediumFromAdapter(t *testing.T) {
	cases := []struct {
		ifType uint32
		desc   string
		want   string
	}{
		{6, "Realtek PCIe GbE Family Controller", "Ethernet"},
		{71, "Intel(R) Wi-Fi 6 AX201 160MHz", "WiFi"},
		{6, "Remote NDIS based Internet Sharing Device", "USB tethering"},
		{6, "Bluetooth Device (Personal Area Network)", "Bluetooth"},
		{243, "Mobile Broadband Device", "Hücresel"},
		{6, "TAP-Windows Adapter V9", "VPN"},
		{131, "WireGuard Tunnel", "VPN"},
		{6, "Some Unknown Adapter", "Ethernet"},
		{999, "Totally Unknown", "Bilinmiyor"},
	}
	for _, c := range cases {
		if got := mediumFromAdapter(c.ifType, c.desc); got != c.want {
			t.Errorf("mediumFromAdapter(%d,%q) = %q, want %q", c.ifType, c.desc, got, c.want)
		}
	}
}
