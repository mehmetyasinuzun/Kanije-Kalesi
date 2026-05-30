// Package access implements a capability-based delegation tree for multi-user
// control of the bot. Each user holds a set of capabilities; a user may invite
// others but can never grant a capability they do not themselves hold
// (no privilege escalation), and may only see/manage users within their own
// invitation subtree. All checks are enforced server-side — the only client
// input ever trusted is the Telegram-authenticated chat ID.
package access

import "sort"

// Capability is a single permission, usually mapping to one bot command.
type Capability string

const (
	CapStatus   Capability = "status"   // /status
	CapEvents   Capability = "events"   // /olaylar, /ozet
	CapPhoto    Capability = "photo"    // /foto
	CapScreen   Capability = "screen"   // /ekran
	CapAudio    Capability = "audio"    // /seskayit
	CapLock     Capability = "lock"     // /kilitle
	CapRestart  Capability = "restart"  // /yeniden
	CapShutdown Capability = "shutdown" // /kapat
	CapConfig   Capability = "config"   // /kurulum, /ayarlar
	CapUpdate   Capability = "update"   // /guncelle
	CapVerify   Capability = "verify"   // /dogrula
	CapInvite   Capability = "invite"   // add sub-users (/ekle)
	CapManage   Capability = "manage"   // edit/remove sub-users (/yonetim)
)

// CapMeta describes a capability for UI rendering.
type CapMeta struct {
	Cap   Capability
	Label string // Turkish, human-readable
	Emoji string
}

// AllCaps is the canonical, ordered capability catalog. Order drives the
// management UI layout and serialization.
var AllCaps = []CapMeta{
	{CapStatus, "Durum", "📊"},
	{CapEvents, "Olaylar", "📋"},
	{CapPhoto, "Kamera", "📷"},
	{CapScreen, "Ekran", "🖥️"},
	{CapAudio, "Ses kaydı", "🎤"},
	{CapLock, "Kilitle", "🔒"},
	{CapRestart, "Yeniden başlat", "🔄"},
	{CapShutdown, "Kapat", "⏻"},
	{CapConfig, "Ayarlar", "⚙️"},
	{CapUpdate, "Güncelle", "⬆️"},
	{CapVerify, "Günlük doğrula", "🛡️"},
	{CapInvite, "Kişi ekle/düzenle", "➕"},
	{CapManage, "Kişi çıkar", "🗑️"},
}

// capLabels indexes AllCaps by capability for O(1) lookup.
var capLabels = func() map[Capability]CapMeta {
	m := make(map[Capability]CapMeta, len(AllCaps))
	for _, c := range AllCaps {
		m[c.Cap] = c
	}
	return m
}()

// Label returns the human-readable Turkish label for a capability.
func (c Capability) Label() string {
	if m, ok := capLabels[c]; ok {
		return m.Label
	}
	return string(c)
}

// Emoji returns the icon for a capability.
func (c Capability) Emoji() string {
	if m, ok := capLabels[c]; ok {
		return m.Emoji
	}
	return "•"
}

// IsKnown reports whether c is a recognized capability.
func (c Capability) IsKnown() bool {
	_, ok := capLabels[c]
	return ok
}

// allCapList is every capability as a flat slice (the owner's full set).
var allCapList = func() []Capability {
	out := make([]Capability, len(AllCaps))
	for i, c := range AllCaps {
		out[i] = c.Cap
	}
	return out
}()

// AllCapabilities returns a fresh copy of every capability (the root/owner set).
func AllCapabilities() []Capability {
	return normalizeCaps(allCapList)
}

// Roles are preset capability bundles offered in the invite UI. Users start
// from a role and then fine-tune individual capabilities.
var Roles = []struct {
	Key   string
	Label string
	Caps  []Capability
}{
	{"izleyici", "👁️ İzleyici (sadece görüntüle)", []Capability{CapStatus, CapEvents}},
	{"operator", "🛠️ Operatör (+medya, kilit)", []Capability{
		CapStatus, CapEvents, CapPhoto, CapScreen, CapAudio, CapLock, CapVerify,
	}},
	{"yonetici", "👑 Yönetici (tüm yetkiler)", allCapList},
}

// RoleCaps returns the capability set for a role key (clamped/normalized), or
// nil if the role is unknown.
func RoleCaps(key string) []Capability {
	for _, r := range Roles {
		if r.Key == key {
			return normalizeCaps(r.Caps)
		}
	}
	return nil
}

// ---- capability set helpers (operate on normalized slices) ----

// normalizeCaps returns a sorted, de-duplicated slice containing only known
// capabilities. nil/empty in → empty (non-nil) out for stable comparisons.
func normalizeCaps(caps []Capability) []Capability {
	seen := make(map[Capability]struct{}, len(caps))
	out := make([]Capability, 0, len(caps))
	for _, c := range caps {
		if !c.IsKnown() {
			continue
		}
		if _, dup := seen[c]; dup {
			continue
		}
		seen[c] = struct{}{}
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// hasCap reports whether caps contains c.
func hasCap(caps []Capability, c Capability) bool {
	for _, x := range caps {
		if x == c {
			return true
		}
	}
	return false
}

// intersect returns the capabilities present in BOTH a and b (normalized).
// This is the no-escalation primitive: granted = requested ∩ granter.
func intersect(a, b []Capability) []Capability {
	bs := make(map[Capability]struct{}, len(b))
	for _, c := range b {
		bs[c] = struct{}{}
	}
	out := make([]Capability, 0, len(a))
	for _, c := range a {
		if _, ok := bs[c]; ok {
			out = append(out, c)
		}
	}
	return normalizeCaps(out)
}
