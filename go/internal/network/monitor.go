// Package network monitors internet connectivity and network changes.
package network

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kanije-kalesi/kanije/internal/event"
	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

// MonitorConfig holds network monitoring settings.
type MonitorConfig struct {
	CheckIntervalSec int
	CheckHost        string
	CheckPort        int
}

// Monitor watches internet connectivity and local network changes.
// It publishes events when connectivity is lost/gained or the network changes.
type Monitor struct {
	cfg      MonitorConfig
	hostname string
	log      *slog.Logger

	lastOnline bool
	lastSSID   string
	lastIface  string
	lastMedium string

	// Cached public (external) IP — refreshed at most every few minutes.
	pubMu   sync.Mutex
	pubIP   string
	pubTime time.Time
}

// NewMonitor creates a Monitor with the given configuration.
func NewMonitor(cfg MonitorConfig, log *slog.Logger) *Monitor {
	h, _ := os.Hostname()
	return &Monitor{cfg: cfg, hostname: h, log: log}
}

// Run starts the monitoring loop and blocks until ctx is canceled.
func (m *Monitor) Run(ctx context.Context, bus *event.Bus) error {
	interval := time.Duration(m.cfg.CheckIntervalSec) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Capture initial state without publishing
	m.lastOnline = m.checkConnectivity()
	m.lastSSID, m.lastIface, m.lastMedium = m.getNetworkInfo()

	m.log.Info("Ağ izleme başlatıldı",
		"hedef", fmt.Sprintf("%s:%d", m.cfg.CheckHost, m.cfg.CheckPort),
		"aralık", interval,
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			m.tick(ctx, bus)
		}
	}
}

func (m *Monitor) tick(ctx context.Context, bus *event.Bus) {
	online := m.checkConnectivity()
	ssid, iface, medium := m.getNetworkInfo()

	switch {
	case online && !m.lastOnline:
		ev := event.New(event.TypeNetworkUp, "NetworkMonitor")
		m.fillNetwork(ctx, &ev, ssid, iface, medium)
		m.log.Info("İnternet bağlantısı kuruldu", "medium", medium, "ağ", ssid, "ip", ev.LocalIP)
		bus.Publish(ev)

	case !online && m.lastOnline:
		ev := event.New(event.TypeNetworkDown, "NetworkMonitor")
		ev.Hostname = m.hostname
		if m.lastMedium != "" {
			ev.Extra = map[string]string{"📡 Kesilen": networkLabel(m.lastMedium, m.lastSSID)}
		}
		m.log.Warn("İnternet bağlantısı kesildi", "önceki", m.lastMedium)
		bus.Publish(ev)

	case online && m.connectionChanged(ssid, iface, medium):
		ev := event.New(event.TypeNetworkChanged, "NetworkMonitor")
		m.fillNetwork(ctx, &ev, ssid, iface, medium)
		ev.Extra["🔀 Değişim"] = networkLabel(m.lastMedium, m.lastSSID) + " → " + networkLabel(medium, ssid)
		m.log.Info("Ağ değişti", "önceki", networkLabel(m.lastMedium, m.lastSSID), "yeni", networkLabel(medium, ssid))
		bus.Publish(ev)
	}

	m.lastOnline = online
	m.lastSSID = ssid
	m.lastIface = iface
	m.lastMedium = medium
}

// connectionChanged reports whether the active connection's medium, SSID, or
// interface differs from the last observed state. Requires a prior observation
// so the first online sample doesn't spuriously fire.
func (m *Monitor) connectionChanged(ssid, iface, medium string) bool {
	if m.lastIface == "" && m.lastMedium == "" {
		return false
	}
	return ssid != m.lastSSID || iface != m.lastIface || medium != m.lastMedium
}

// fillNetwork populates the common network fields (medium, SSID, IPs, interface)
// on a network event.
func (m *Monitor) fillNetwork(ctx context.Context, ev *event.Event, ssid, iface, medium string) {
	ev.Hostname = m.hostname
	ev.NetworkSSID = ssid
	ev.NetworkType = medium
	ev.LocalIP = getLocalIP(iface)
	ev.PublicIP = m.publicIP(ctx)
	ev.Extra = map[string]string{}
	if iface != "" {
		ev.Extra["🔌 Arayüz"] = iface
	}
}

// networkLabel renders a compact "medium (SSID)" description for logs/events.
func networkLabel(medium, ssid string) string {
	switch {
	case medium == "" && ssid == "":
		return "bilinmiyor"
	case ssid != "" && medium != "":
		return medium + " (" + ssid + ")"
	case ssid != "":
		return ssid
	default:
		return medium
	}
}

// PublicIP returns the device's external/public IP (cached up to 5 minutes).
// Exposed for /status and /cihazlar.
func (m *Monitor) PublicIP(ctx context.Context) string {
	return m.publicIP(ctx)
}

// publicIP returns the cached external IP, refreshing it at most every 5 minutes.
func (m *Monitor) publicIP(ctx context.Context) string {
	m.pubMu.Lock()
	if m.pubIP != "" && time.Since(m.pubTime) < 5*time.Minute {
		ip := m.pubIP
		m.pubMu.Unlock()
		return ip
	}
	m.pubMu.Unlock()

	ip := fetchPublicIP(ctx)
	if ip != "" {
		m.pubMu.Lock()
		m.pubIP, m.pubTime = ip, time.Now()
		m.pubMu.Unlock()
	}
	return ip
}

// fetchPublicIP queries api.ipify.org for the external IP. Returns "" on error.
func fetchPublicIP(ctx context.Context) string {
	reqCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, "https://api.ipify.org", nil)
	if err != nil {
		return ""
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return ""
	}
	ip := strings.TrimSpace(string(body))
	if net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}

// checkConnectivity tests TCP connectivity to the configured target.
func (m *Monitor) checkConnectivity() bool {
	host := m.cfg.CheckHost
	port := m.cfg.CheckPort
	if host == "" {
		host = "api.telegram.org"
	}
	if port == 0 {
		port = 443
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 3*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// getNetworkInfo returns the active connection's WiFi SSID, interface name, and
// medium (WiFi/Ethernet/USB tethering/Bluetooth/Hücresel/VPN).
func (m *Monitor) getNetworkInfo() (ssid, iface, medium string) {
	switch runtime.GOOS {
	case "windows":
		// IP Helper gives the precise active adapter + medium (incl. USB
		// tethering / Bluetooth / cellular). netsh gives the WiFi SSID.
		if friendly, med := activeAdapter(); friendly != "" {
			ssid, _ = getWindowsNetwork()
			return ssid, friendly, med
		}
		ssid, iface = getWindowsNetwork()
		return ssid, iface, inferNetworkType(iface)
	case "linux":
		ssid, iface = getLinuxNetwork()
		return ssid, iface, inferNetworkType(iface)
	default:
		return "", "", ""
	}
}

func getWindowsNetwork() (ssid, iface string) {
	// This runs every poll on a -H=windowsgui binary, so the console window must
	// be suppressed (sysproc.Hide) or netsh would flash a window each time.
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	sysproc.Hide(cmd)
	out, err := cmd.Output()
	if err != nil {
		return "", "Ethernet"
	}

	// netsh output is usually in system locale, try to detect SSID line
	text := string(out)
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		// Match "SSID" but not "BSSID"
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if strings.EqualFold(key, "SSID") && !strings.EqualFold(key, "BSSID") {
			ssid = val
		}
		// "Name" in English, "Ad" in Turkish
		if strings.EqualFold(key, "Name") || key == "Ad" {
			iface = val
		}
	}
	return
}

func getLinuxNetwork() (ssid, iface string) {
	// WiFi SSID via iwgetid (sysproc.Hide is a no-op on Linux).
	wlanCmd := exec.Command("iwgetid", "-r")
	sysproc.Hide(wlanCmd)
	if out, err := wlanCmd.Output(); err == nil {
		ssid = strings.TrimSpace(string(out))
	}

	// Default route interface
	routeCmd := exec.Command("ip", "route", "show", "default")
	sysproc.Hide(routeCmd)
	if out, err := routeCmd.Output(); err == nil {
		parts := strings.Fields(string(out))
		for i, p := range parts {
			if p == "dev" && i+1 < len(parts) {
				iface = parts[i+1]
				break
			}
		}
	}
	return
}

func inferNetworkType(iface string) string {
	l := strings.ToLower(iface)
	switch {
	case strings.Contains(l, "wi-fi"), strings.Contains(l, "wifi"),
		strings.HasPrefix(l, "wlan"), strings.HasPrefix(l, "wlp"),
		strings.Contains(l, "wireless"):
		return "WiFi"
	case strings.HasPrefix(l, "bnep"), strings.Contains(l, "bluetooth"):
		return "Bluetooth"
	case strings.HasPrefix(l, "usb"), strings.Contains(l, "rndis"), strings.Contains(l, "ncm"):
		return "USB tethering"
	case strings.HasPrefix(l, "wwan"), strings.HasPrefix(l, "ppp"),
		strings.Contains(l, "cellular"), strings.Contains(l, "mobile"):
		return "Hücresel"
	case strings.HasPrefix(l, "tun"), strings.HasPrefix(l, "tap"),
		strings.HasPrefix(l, "wg"), strings.Contains(l, "vpn"),
		strings.Contains(l, "wireguard"), strings.Contains(l, "tailscale"):
		return "VPN"
	case strings.HasPrefix(l, "eth"), strings.HasPrefix(l, "en"),
		strings.Contains(l, "ethernet"), strings.Contains(l, "lan"):
		return "Ethernet"
	}
	return "Bilinmiyor"
}

func getLocalIP(ifaceName string) string {
	if ifaceName == "" {
		return ""
	}
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return ""
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}
		if ip.To4() != nil {
			return ip.String()
		}
	}
	return ""
}
