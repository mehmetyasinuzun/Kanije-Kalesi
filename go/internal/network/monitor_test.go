package network

import "testing"

func TestInferNetworkType(t *testing.T) {
	cases := map[string]string{
		"Wi-Fi":                        "WiFi",
		"wlan0":                        "WiFi",
		"wlp3s0":                       "WiFi",
		"Ethernet 2":                   "Ethernet",
		"eth0":                         "Ethernet",
		"enp0s3":                       "Ethernet",
		"bnep0":                        "Bluetooth",
		"Bluetooth Network Connection": "Bluetooth",
		"usb0":                         "USB tethering",
		"wwan0":                        "Hücresel",
		"ppp0":                         "Hücresel",
		"tun0":                         "VPN",
		"wg0":                          "VPN",
		"tailscale0":                   "VPN",
		"":                             "Bilinmiyor",
		"weirdadapter":                 "Bilinmiyor",
	}
	for in, want := range cases {
		if got := inferNetworkType(in); got != want {
			t.Errorf("inferNetworkType(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNetworkLabel(t *testing.T) {
	cases := []struct {
		medium, ssid, want string
	}{
		{"WiFi", "EvAğı", "WiFi (EvAğı)"},
		{"Ethernet", "", "Ethernet"},
		{"", "GizliSSID", "GizliSSID"},
		{"USB tethering", "", "USB tethering"},
		{"", "", "bilinmiyor"},
	}
	for _, c := range cases {
		if got := networkLabel(c.medium, c.ssid); got != c.want {
			t.Errorf("networkLabel(%q,%q) = %q, want %q", c.medium, c.ssid, got, c.want)
		}
	}
}
