package network

import "testing"

func TestInferNetworkType(t *testing.T) {
	tests := []struct {
		name  string
		iface string
		want  string
	}{
		{name: "windows wifi", iface: "Wi-Fi", want: "WiFi"},
		{name: "linux wifi", iface: "wlan0", want: "WiFi"},
		{name: "ethernet", iface: "Ethernet 2", want: "Ethernet"},
		{name: "eth short", iface: "eth0", want: "Ethernet"},
		{name: "unknown", iface: "tun0", want: "Bilinmiyor"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inferNetworkType(tt.iface); got != tt.want {
				t.Fatalf("inferNetworkType(%q) = %q, want %q", tt.iface, got, tt.want)
			}
		})
	}
}
