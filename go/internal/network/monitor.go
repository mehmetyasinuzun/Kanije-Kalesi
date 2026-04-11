// Package network monitors internet connectivity and network changes.
package network

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/kanije-kalesi/sentinel/internal/event"
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
}

// NewMonitor creates a Monitor with the given configuration.
func NewMonitor(cfg MonitorConfig, log *slog.Logger) *Monitor {
	h, _ := os.Hostname()
	return &Monitor{cfg: cfg, hostname: h, log: log}
}

// Run starts the monitoring loop and blocks until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context, bus *event.Bus) error {
	interval := time.Duration(m.cfg.CheckIntervalSec) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Capture initial state without publishing
	m.lastOnline = m.checkConnectivity()
	m.lastSSID, m.lastIface = m.getNetworkInfo()

	m.log.Info("Ağ izleme başlatıldı",
		"hedef", fmt.Sprintf("%s:%d", m.cfg.CheckHost, m.cfg.CheckPort),
		"aralık", interval,
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			m.tick(bus)
		}
	}
}

func (m *Monitor) tick(bus *event.Bus) {
	online := m.checkConnectivity()
	ssid, iface := m.getNetworkInfo()

	switch {
	case online && !m.lastOnline:
		ev := event.New(event.TypeNetworkUp, "NetworkMonitor")
		ev.Hostname    = m.hostname
		ev.NetworkSSID = ssid
		ev.NetworkType = inferNetworkType(iface)
		ev.LocalIP     = getLocalIP(iface)
		m.log.Info("İnternet bağlantısı kuruldu", "ağ", ssid, "ip", ev.LocalIP)
		bus.Publish(ev)

	case !online && m.lastOnline:
		ev := event.New(event.TypeNetworkDown, "NetworkMonitor")
		ev.Hostname = m.hostname
		m.log.Warn("İnternet bağlantısı kesildi")
		bus.Publish(ev)

	case online && (ssid != m.lastSSID) && m.lastSSID != "":
		ev := event.New(event.TypeNetworkChanged, "NetworkMonitor")
		ev.Hostname    = m.hostname
		ev.NetworkSSID = ssid
		ev.NetworkType = inferNetworkType(iface)
		ev.LocalIP     = getLocalIP(iface)
		m.log.Info("Ağ değişti", "önceki", m.lastSSID, "yeni", ssid)
		bus.Publish(ev)
	}

	m.lastOnline = online
	m.lastSSID   = ssid
	m.lastIface  = iface
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

// getNetworkInfo returns the active WiFi SSID and interface name.
func (m *Monitor) getNetworkInfo() (ssid, iface string) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsNetwork()
	case "linux":
		return getLinuxNetwork()
	default:
		return "", ""
	}
}

func getWindowsNetwork() (ssid, iface string) {
	// Use UTF-8 output explicitly: netsh outputs in system locale by default
	out, err := exec.Command("netsh", "wlan", "show", "interfaces").Output()
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
	// WiFi SSID via iwgetid
	if out, err := exec.Command("iwgetid", "-r").Output(); err == nil {
		ssid = strings.TrimSpace(string(out))
	}

	// Default route interface
	if out, err := exec.Command("ip", "route", "show", "default").Output(); err == nil {
		for _, field := range strings.Fields(string(out)) {
			if iface != "" {
				break
			}
		}
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
	lower := strings.ToLower(iface)
	switch {
	case strings.ContainsAny(lower, "") ||
		strings.Contains(lower, "wi-fi") ||
		strings.Contains(lower, "wifi") ||
		strings.Contains(lower, "wlan") ||
		strings.Contains(lower, "wireless"):
		return "WiFi"
	case strings.Contains(lower, "eth") ||
		strings.Contains(lower, "ethernet") ||
		strings.Contains(lower, "lan"):
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
