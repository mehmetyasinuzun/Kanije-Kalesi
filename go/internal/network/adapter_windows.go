//go:build windows

package network

import (
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// IANA ifType values + IP Helper flags/states we care about. Defined locally to
// avoid depending on symbol names that vary across x/sys versions.
const (
	gaaFlagIncludeGateways = 0x0080
	ifOperStatusUp         = 1
	ifTypeEthernet         = 6  // also covers USB/Bluetooth bridges → use description
	ifTypePPP              = 23 // dial-up / some VPNs
	ifTypeSoftwareLoopback = 24
	ifTypeTunnel           = 131 // VPN tunnels
	ifTypeWiFi             = 71  // IEEE 802.11
	ifTypeWWANPP           = 243 // cellular (GSM)
	ifTypeWWANPP2          = 244 // cellular (CDMA)
)

// activeAdapter returns the friendly name and inferred medium of the network
// adapter currently carrying outbound traffic (the one with a default gateway).
// medium is one of: WiFi, Ethernet, USB tethering, Bluetooth, Hücresel, VPN,
// Bilinmiyor. Returns ("","") if it cannot be determined.
func activeAdapter() (friendly, medium string) {
	size := uint32(15000)
	var buf []byte
	var head *windows.IpAdapterAddresses

	for attempt := 0; attempt < 3; attempt++ {
		buf = make([]byte, size)
		head = (*windows.IpAdapterAddresses)(unsafe.Pointer(&buf[0]))
		err := windows.GetAdaptersAddresses(windows.AF_UNSPEC, gaaFlagIncludeGateways, 0, head, &size)
		if err == nil {
			break
		}
		if err == windows.ERROR_BUFFER_OVERFLOW {
			continue // size now holds the required length; retry
		}
		return "", ""
	}

	var fallback *windows.IpAdapterAddresses
	for p := head; p != nil; p = p.Next {
		if p.OperStatus != ifOperStatusUp || p.IfType == ifTypeSoftwareLoopback {
			continue
		}
		if p.FirstUnicastAddress == nil {
			continue
		}
		// The adapter with a default gateway is the one routing to the internet.
		if p.FirstGatewayAddress != nil {
			return adapterInfo(p)
		}
		if fallback == nil {
			fallback = p
		}
	}
	if fallback != nil {
		return adapterInfo(fallback)
	}
	return "", ""
}

func adapterInfo(p *windows.IpAdapterAddresses) (string, string) {
	friendly := windows.UTF16PtrToString(p.FriendlyName)
	desc := windows.UTF16PtrToString(p.Description)
	return friendly, mediumFromAdapter(p.IfType, desc)
}

// mediumFromAdapter maps an adapter's IANA ifType + description to a friendly
// medium. Description is checked first because USB tethering and Bluetooth PAN
// both report as Ethernet (ifType 6) but are distinguishable by name.
func mediumFromAdapter(ifType uint32, desc string) string {
	d := strings.ToLower(desc)
	switch {
	case strings.Contains(d, "bluetooth"):
		return "Bluetooth"
	case strings.Contains(d, "remote ndis"), strings.Contains(d, "rndis"),
		strings.Contains(d, "usb") && !strings.Contains(d, "wireless"):
		return "USB tethering"
	case strings.Contains(d, "wireless wan"), strings.Contains(d, "mobile broadband"),
		strings.Contains(d, "cellular"), strings.Contains(d, "wwan"),
		ifType == ifTypeWWANPP, ifType == ifTypeWWANPP2:
		return "Hücresel"
	case strings.Contains(d, "tap-"), strings.Contains(d, "tun"),
		strings.Contains(d, "vpn"), strings.Contains(d, "wireguard"),
		strings.Contains(d, "openvpn"), strings.Contains(d, "tailscale"),
		strings.Contains(d, "wintun"), ifType == ifTypeTunnel, ifType == ifTypePPP:
		return "VPN"
	case ifType == ifTypeWiFi:
		return "WiFi"
	case ifType == ifTypeEthernet:
		return "Ethernet"
	}
	return "Bilinmiyor"
}
