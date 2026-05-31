// Package defender reads Microsoft Defender's status and recent activity so the
// owner can see whether (and how actively) the AV is watching the machine, when
// it last scanned, and what it has detected. Windows-only; other platforms report
// "unavailable".
package defender

// Status is a snapshot of Microsoft Defender's posture.
type Status struct {
	Available          bool   // false if Defender/PowerShell module isn't present
	RealTimeProtection bool   // on-access scanning active
	AntivirusEnabled   bool   // AV engine enabled
	TamperProtection   bool   // Defender's own tamper protection
	LastQuickScan      string // human-readable time of the last quick scan ("" if never)
	LastFullScan       string // last full scan ("" if never)
	SignatureVersion   string // AV signature version
	SignatureUpdated   string // when signatures were last updated
	RecentThreats      []Threat
}

// Threat is one detection from Defender's history.
type Threat struct {
	Name   string
	Time   string
	Action string
}
