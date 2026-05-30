// Package geoip enriches an IP address with country, city and ISP information
// using the free, key-less ipwho.is HTTPS API. Results are cached so repeated
// attempts from the same address don't hammer the service.
package geoip

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Info is the geolocation result for an IP address.
type Info struct {
	Country     string
	CountryCode string
	City        string
	ISP         string
	Flag        string // regional-indicator emoji, e.g. 🇷🇺
}

// Summary renders a compact "🇷🇺 Moscow, Russia" style string.
func (i *Info) Summary() string {
	loc := i.Country
	if i.City != "" && i.Country != "" {
		loc = i.City + ", " + i.Country
	} else if i.City != "" {
		loc = i.City
	}
	if i.Flag != "" && loc != "" {
		return i.Flag + " " + loc
	}
	if loc == "" {
		return i.Flag
	}
	return loc
}

// Resolver looks up IP addresses with an in-memory TTL cache.
type Resolver struct {
	http  *http.Client
	mu    sync.Mutex
	cache map[string]entry
	ttl   time.Duration
}

type entry struct {
	info    *Info
	expires time.Time
}

// New creates a Resolver with sensible defaults.
func New() *Resolver {
	return &Resolver{
		http:  &http.Client{Timeout: 5 * time.Second},
		cache: make(map[string]entry),
		ttl:   6 * time.Hour,
	}
}

// Lookup returns geolocation for a public IP. It returns (nil, false) for
// private/loopback/invalid addresses and on any error.
func (r *Resolver) Lookup(ctx context.Context, ip string) (*Info, bool) {
	if !IsPublic(ip) {
		return nil, false
	}
	if info, cached := r.fromCache(ip); cached {
		return info, info != nil
	}
	info := r.fetch(ctx, ip)
	r.store(ip, info)
	return info, info != nil
}

func (r *Resolver) fromCache(ip string) (*Info, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.cache[ip]
	if !ok || time.Now().After(e.expires) {
		return nil, false
	}
	return e.info, true
}

func (r *Resolver) store(ip string, info *Info) {
	r.mu.Lock()
	r.cache[ip] = entry{info: info, expires: time.Now().Add(r.ttl)}
	r.mu.Unlock()
}

type ipwhoResp struct {
	Success     bool   `json:"success"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
	Connection  struct {
		ISP string `json:"isp"`
		Org string `json:"org"`
	} `json:"connection"`
}

func (r *Resolver) fetch(ctx context.Context, ip string) *Info {
	url := "https://ipwho.is/" + ip + "?fields=success,country,country_code,city,connection"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "kanije-kalesi")

	resp, err := r.http.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var w ipwhoResp
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<16)).Decode(&w); err != nil {
		return nil
	}
	if !w.Success || w.Country == "" {
		return nil
	}

	isp := w.Connection.ISP
	if isp == "" {
		isp = w.Connection.Org
	}
	return &Info{
		Country:     w.Country,
		CountryCode: w.CountryCode,
		City:        w.City,
		ISP:         isp,
		Flag:        FlagEmoji(w.CountryCode),
	}
}

// FlagEmoji converts a two-letter ISO country code to its flag emoji.
func FlagEmoji(cc string) string {
	if len(cc) != 2 {
		return ""
	}
	cc = strings.ToUpper(cc)
	a, b := cc[0], cc[1]
	if a < 'A' || a > 'Z' || b < 'A' || b > 'Z' {
		return ""
	}
	return string(rune(0x1F1E6+int(a-'A'))) + string(rune(0x1F1E6+int(b-'A')))
}

// IsPublic reports whether ip is a routable public address worth geolocating.
func IsPublic(ip string) bool {
	p := net.ParseIP(ip)
	if p == nil {
		return false
	}
	return !(p.IsLoopback() || p.IsPrivate() || p.IsLinkLocalUnicast() ||
		p.IsLinkLocalMulticast() || p.IsUnspecified())
}
