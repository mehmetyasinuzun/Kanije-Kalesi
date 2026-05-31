package access

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"
)

// User is a member of the delegation tree.
type User struct {
	ChatID    int64
	Name      string
	InvitedBy int64 // 0 = root/owner
	Caps      []Capability
	AddedAt   time.Time
}

// IsRoot reports whether u is the owner node.
func (u User) IsRoot() bool { return u.InvitedBy == 0 }

// AuditEntry is a record of a privileged action, for the /loglar command.
type AuditEntry struct {
	ID     int64
	ChatID int64 // actor chat ID
	Actor  string
	Action string // short verb: "foto", "ekle", "cikar", "duzenle", "kapat"...
	Target string // affected user / extra detail (optional)
	At     time.Time
}

// Persistence is the storage contract the Manager needs. The SQLite storage
// implements it structurally (no import cycle: access defines the interface).
type Persistence interface {
	LoadUsers(ctx context.Context) ([]User, error)
	SaveUser(ctx context.Context, u User) error
	DeleteUser(ctx context.Context, chatID int64) error
	AppendAudit(ctx context.Context, e AuditEntry) error
	RecentAudit(ctx context.Context, n int) ([]AuditEntry, error)
}

// ReparentMode decides where a removed user's children are re-attached.
// Removal never deletes the subtree — orphans are always re-homed.
type ReparentMode int

const (
	// ReparentRoot attaches orphans to the owner (top). This is the default.
	ReparentRoot ReparentMode = iota
	// ReparentGrandparent attaches orphans to the removed user's parent.
	ReparentGrandparent
	// ReparentRemover attaches orphans to whoever performed the removal.
	ReparentRemover
)

// Manager holds the delegation tree in memory and persists changes. The hot
// path (Can) is read-locked and never does I/O; mutations are rare and
// serialized under the write lock.
type Manager struct {
	mu     sync.RWMutex
	users  map[int64]*User
	rootID int64
	store  Persistence
	log    *slog.Logger
}

// NewManager creates an empty manager. Call Bootstrap before use.
func NewManager(store Persistence, log *slog.Logger) *Manager {
	return &Manager{
		users: make(map[int64]*User),
		store: store,
		log:   log,
	}
}

// Bootstrap loads users from storage. If none exist (fresh install), it seeds
// the owner as root (all capabilities) and migrates any pre-existing allowlist
// entries as full-capability children of root, preserving their prior access.
func (m *Manager) Bootstrap(ctx context.Context, ownerChatID int64, allowed []int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	loaded, err := m.store.LoadUsers(ctx)
	if err != nil {
		return fmt.Errorf("kullanıcılar yüklenemedi: %w", err)
	}

	if len(loaded) > 0 {
		for i := range loaded {
			u := loaded[i]
			u.Caps = normalizeCaps(u.Caps)
			m.users[u.ChatID] = &u
			if u.IsRoot() {
				m.rootID = u.ChatID
			}
		}
		// Defensive: if no explicit root was stored, adopt the owner.
		if m.rootID == 0 {
			m.rootID = ownerChatID
		}
		// Self-heal: the owner must always hold EVERY capability, including any
		// added by a newer version (e.g. terminal). The stored set can predate
		// new caps, which is why /terminal could report "yetki yok" after an update.
		m.grantOwnerFullCapsLocked(ctx, ownerChatID)
		return nil
	}

	// Fresh install — seed root + migrate allowlist.
	now := time.Now().UTC()
	root := &User{ChatID: ownerChatID, Name: "Sahip", InvitedBy: 0, Caps: AllCapabilities(), AddedAt: now}
	m.users[ownerChatID] = root
	m.rootID = ownerChatID
	if err := m.store.SaveUser(ctx, *root); err != nil {
		return fmt.Errorf("kök kullanıcı kaydedilemedi: %w", err)
	}

	for _, id := range allowed {
		if id == 0 || id == ownerChatID {
			continue
		}
		if _, exists := m.users[id]; exists {
			continue
		}
		u := &User{ChatID: id, Name: fmt.Sprintf("Kullanıcı %d", id), InvitedBy: ownerChatID, Caps: AllCapabilities(), AddedAt: now}
		m.users[id] = u
		if err := m.store.SaveUser(ctx, *u); err != nil {
			return fmt.Errorf("izinli kullanıcı taşınamadı: %w", err)
		}
	}
	return nil
}

// grantOwnerFullCapsLocked ensures the owner row exists and holds every current
// capability. Must hold the write lock. Idempotent — only writes on change.
func (m *Manager) grantOwnerFullCapsLocked(ctx context.Context, ownerChatID int64) {
	root, ok := m.users[m.rootID]
	if !ok {
		root = &User{ChatID: ownerChatID, Name: "Sahip", InvitedBy: 0, Caps: AllCapabilities(), AddedAt: time.Now().UTC()}
		m.users[ownerChatID] = root
		m.rootID = ownerChatID
		_ = m.store.SaveUser(ctx, *root)
		return
	}
	full := AllCapabilities()
	if !sameCaps(root.Caps, full) {
		root.Caps = full
		_ = m.store.SaveUser(ctx, *root)
	}
}

// RootID returns the owner chat ID.
func (m *Manager) RootID() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rootID
}

// Known reports whether chatID is a member of the tree at all.
func (m *Manager) Known(chatID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.users[chatID]
	return ok
}

// Can reports whether chatID holds capability c. Unknown users have no caps.
func (m *Manager) Can(chatID int64, c Capability) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[chatID]
	return ok && hasCap(u.Caps, c)
}

// Get returns a copy of a user.
func (m *Manager) Get(chatID int64) (User, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[chatID]
	if !ok {
		return User{}, false
	}
	return cloneUser(u), true
}

// ListEntry is one row of a flattened subtree, with indentation depth.
type ListEntry struct {
	User  User
	Depth int
}

// Subtree returns viewer's descendants (excluding viewer), depth-first and
// sorted, for the management UI. Use IncludeSelf via Get for the viewer itself.
func (m *Manager) Subtree(viewer int64) []ListEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []ListEntry
	m.walkLocked(viewer, 0, &out)
	return out
}

func (m *Manager) walkLocked(parent int64, depth int, out *[]ListEntry) {
	kids := m.childrenLocked(parent)
	sort.Slice(kids, func(i, j int) bool {
		if !kids[i].AddedAt.Equal(kids[j].AddedAt) {
			return kids[i].AddedAt.Before(kids[j].AddedAt)
		}
		return kids[i].ChatID < kids[j].ChatID
	})
	for _, k := range kids {
		*out = append(*out, ListEntry{User: cloneUser(k), Depth: depth})
		m.walkLocked(k.ChatID, depth+1, out)
	}
}

func (m *Manager) childrenLocked(parent int64) []*User {
	var kids []*User
	for _, u := range m.users {
		if u.InvitedBy == parent && u.ChatID != parent {
			kids = append(kids, u)
		}
	}
	return kids
}

// Visible reports whether target is within viewer's subtree (so viewer may
// see/manage it). A user can always see itself.
func (m *Manager) Visible(viewer, target int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if viewer == target {
		return true
	}
	return m.isDescendantLocked(target, viewer)
}

// isDescendantLocked reports whether target is strictly below ancestor.
func (m *Manager) isDescendantLocked(target, ancestor int64) bool {
	if target == ancestor {
		return false
	}
	cur, ok := m.users[target]
	if !ok {
		return false
	}
	// Walk up parent links; a tree has no cycles, but guard with a hop cap.
	for hops := 0; cur != nil && hops <= len(m.users)+1; hops++ {
		if cur.InvitedBy == ancestor {
			return true
		}
		if cur.InvitedBy == 0 {
			return false
		}
		cur = m.users[cur.InvitedBy]
	}
	return false
}

// Invite adds newChat under inviter, granting requested ∩ inviter.Caps
// (no escalation). Requires the inviter to hold CapInvite.
func (m *Manager) Invite(ctx context.Context, inviter, newChat int64, name string, requested []Capability) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	inv, ok := m.users[inviter]
	if !ok {
		return User{}, fmt.Errorf("davet eden tanımlı değil")
	}
	if !hasCap(inv.Caps, CapInvite) {
		return User{}, fmt.Errorf("davet etme yetkin yok")
	}
	if newChat == 0 {
		return User{}, fmt.Errorf("geçersiz chat ID")
	}
	if newChat == inviter {
		return User{}, fmt.Errorf("kendini ekleyemezsin")
	}
	if _, exists := m.users[newChat]; exists {
		return User{}, fmt.Errorf("bu kişi zaten ekli")
	}

	granted := intersect(requested, inv.Caps) // hard cap: never exceeds inviter
	if name == "" {
		name = fmt.Sprintf("Kullanıcı %d", newChat)
	}
	u := &User{ChatID: newChat, Name: name, InvitedBy: inviter, Caps: granted, AddedAt: time.Now().UTC()}
	m.users[newChat] = u

	if err := m.store.SaveUser(ctx, *u); err != nil {
		delete(m.users, newChat) // roll back in-memory on persist failure
		return User{}, fmt.Errorf("kaydedilemedi: %w", err)
	}
	m.auditLocked(ctx, inviter, "ekle", fmt.Sprintf("%s (%d)", name, newChat))
	return cloneUser(u), nil
}

// SetCaps edits target's capabilities to requested ∩ editor.Caps, then
// cascade-clamps every descendant so none exceeds its (possibly reduced)
// parent. Requires editor to hold CapManage and target to be in editor's subtree.
func (m *Manager) SetCaps(ctx context.Context, editor, target int64, requested []Capability) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ed, ok := m.users[editor]
	if !ok {
		return fmt.Errorf("düzenleyen tanımlı değil")
	}
	t, ok := m.users[target]
	if !ok {
		return fmt.Errorf("kişi bulunamadı")
	}
	if t.IsRoot() {
		return fmt.Errorf("kök kullanıcı düzenlenemez")
	}
	if !hasCap(ed.Caps, CapInvite) {
		return fmt.Errorf("düzenleme (davet) yetkin yok")
	}
	if !m.isDescendantLocked(target, editor) {
		return fmt.Errorf("yalnızca kendi alt ağacındaki kişileri düzenleyebilirsin")
	}

	t.Caps = intersect(requested, ed.Caps)
	changed := []*User{t}
	m.clampDescendantsLocked(target, &changed)

	for _, u := range changed {
		if err := m.store.SaveUser(ctx, *u); err != nil {
			return fmt.Errorf("kaydedilemedi: %w", err)
		}
	}
	m.auditLocked(ctx, editor, "duzenle", fmt.Sprintf("%s (%d)", t.Name, target))
	return nil
}

// clampDescendantsLocked walks below parent, intersecting each child's caps
// with its parent's, recursively. Appends every modified user to changed.
func (m *Manager) clampDescendantsLocked(parent int64, changed *[]*User) {
	pcaps := m.users[parent].Caps
	for _, c := range m.childrenLocked(parent) {
		next := intersect(c.Caps, pcaps)
		if !sameCaps(next, c.Caps) {
			c.Caps = next
			*changed = append(*changed, c)
		}
		m.clampDescendantsLocked(c.ChatID, changed)
	}
}

// Remove deletes target and re-homes its direct children per mode (default
// root). Orphans (and their subtrees) are cap-clamped to the new parent.
// Requires remover to hold CapManage and target to be in remover's subtree.
func (m *Manager) Remove(ctx context.Context, remover, target int64, mode ReparentMode) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	rm, ok := m.users[remover]
	if !ok {
		return fmt.Errorf("çıkaran tanımlı değil")
	}
	t, ok := m.users[target]
	if !ok {
		return fmt.Errorf("kişi bulunamadı")
	}
	if t.IsRoot() {
		return fmt.Errorf("kök kullanıcı çıkarılamaz")
	}
	if target == remover {
		return fmt.Errorf("kendini çıkaramazsın")
	}
	if !hasCap(rm.Caps, CapManage) {
		return fmt.Errorf("çıkarma yetkin yok")
	}
	if !m.isDescendantLocked(target, remover) {
		return fmt.Errorf("yalnızca kendi alt ağacındaki kişileri çıkarabilirsin")
	}

	var newParent int64
	switch mode {
	case ReparentGrandparent:
		newParent = t.InvitedBy
	case ReparentRemover:
		newParent = remover
	default:
		newParent = m.rootID
	}
	if _, ok := m.users[newParent]; !ok {
		newParent = m.rootID
	}

	var changed []*User
	for _, c := range m.childrenLocked(target) {
		c.InvitedBy = newParent
		c.Caps = intersect(c.Caps, m.users[newParent].Caps)
		changed = append(changed, c)
		m.clampDescendantsLocked(c.ChatID, &changed)
	}

	delete(m.users, target)
	if err := m.store.DeleteUser(ctx, target); err != nil {
		return fmt.Errorf("silinemedi: %w", err)
	}
	for _, u := range changed {
		if err := m.store.SaveUser(ctx, *u); err != nil {
			return fmt.Errorf("yeniden bağlama kaydedilemedi: %w", err)
		}
	}
	m.auditLocked(ctx, remover, "cikar", fmt.Sprintf("%s (%d)", t.Name, target))
	return nil
}

// TransferRoot hands the whole device to a new owner: it wipes the entire
// delegation tree (the old owner and every sub-user lose access) and re-seeds a
// single new root holding all capabilities. Used by /aktar. Idempotent if
// newRootID is already the sole root.
func (m *Manager) TransferRoot(ctx context.Context, newRootID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if newRootID == 0 {
		return fmt.Errorf("geçersiz yeni sahip ID")
	}

	// Remove every existing user from storage, then reset in-memory state.
	for id := range m.users {
		if err := m.store.DeleteUser(ctx, id); err != nil {
			return fmt.Errorf("eski erişim ağacı silinemedi: %w", err)
		}
	}
	m.users = make(map[int64]*User)

	root := &User{ChatID: newRootID, Name: "Sahip", InvitedBy: 0, Caps: AllCapabilities(), AddedAt: time.Now().UTC()}
	m.users[newRootID] = root
	m.rootID = newRootID
	if err := m.store.SaveUser(ctx, *root); err != nil {
		return fmt.Errorf("yeni sahip kaydedilemedi: %w", err)
	}
	m.auditLocked(ctx, newRootID, "aktar", fmt.Sprintf("yeni sahip %d", newRootID))
	return nil
}

// HasChildren reports whether target has any direct children (used by the UI to
// decide whether to ask about re-parenting on removal).
func (m *Manager) HasChildren(target int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.childrenLocked(target)) > 0
}

// Log records a privileged action performed by chatID (best-effort).
func (m *Manager) Log(ctx context.Context, chatID int64, action, target string) {
	m.mu.RLock()
	name := ""
	if u, ok := m.users[chatID]; ok {
		name = u.Name
	}
	m.mu.RUnlock()
	m.appendAudit(ctx, AuditEntry{ChatID: chatID, Actor: name, Action: action, Target: target, At: time.Now().UTC()})
}

// RecentAudit returns the latest n audit entries, newest first.
func (m *Manager) RecentAudit(ctx context.Context, n int) ([]AuditEntry, error) {
	return m.store.RecentAudit(ctx, n)
}

// auditLocked records an action while the write lock is held (name already known).
func (m *Manager) auditLocked(ctx context.Context, chatID int64, action, target string) {
	name := ""
	if u, ok := m.users[chatID]; ok {
		name = u.Name
	}
	m.appendAudit(ctx, AuditEntry{ChatID: chatID, Actor: name, Action: action, Target: target, At: time.Now().UTC()})
}

func (m *Manager) appendAudit(ctx context.Context, e AuditEntry) {
	if err := m.store.AppendAudit(ctx, e); err != nil {
		m.log.Warn("denetim günlüğü yazılamadı", "hata", err)
	}
}

func cloneUser(u *User) User {
	cp := *u
	cp.Caps = append([]Capability(nil), u.Caps...)
	return cp
}

func sameCaps(a, b []Capability) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
