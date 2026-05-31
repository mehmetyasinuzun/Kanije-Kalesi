//go:build windows

package defender

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

// statusScript collects Defender posture + recent detections and emits compact
// JSON. Everything is best-effort (SilentlyContinue) so a missing cmdlet on an
// older SKU yields an empty object rather than an error.
const statusScript = `$ErrorActionPreference='SilentlyContinue'
$s = Get-MpComputerStatus
if (-not $s) { '{}'; return }
$names = @{}
Get-MpThreat | ForEach-Object { $names[$_.ThreatID] = $_.ThreatName }
$threats = Get-MpThreatDetection | Sort-Object InitialDetectionTime -Descending | Select-Object -First 5 | ForEach-Object {
  $n = if ($names.ContainsKey($_.ThreatID)) { $names[$_.ThreatID] } else { "ID:$($_.ThreatID)" }
  [PSCustomObject]@{ name = "$n"; time = "$($_.InitialDetectionTime)" }
}
[PSCustomObject]@{
  rtp        = [bool]$s.RealTimeProtectionEnabled
  av         = [bool]$s.AntivirusEnabled
  tamper     = [bool]$s.IsTamperProtected
  quickScan  = "$($s.QuickScanEndTime)"
  fullScan   = "$($s.FullScanEndTime)"
  sigVer     = "$($s.AntivirusSignatureVersion)"
  sigUpdated = "$($s.AntivirusSignatureLastUpdated)"
  threats    = @($threats)
} | ConvertTo-Json -Depth 4 -Compress`

// GetStatus queries Microsoft Defender via PowerShell. On any failure it returns
// a zero Status with Available=false.
func GetStatus() Status {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", statusScript)
	sysproc.Hide(cmd)

	out, err := cmd.Output()
	if err != nil {
		return Status{}
	}

	var raw struct {
		RTP        bool   `json:"rtp"`
		AV         bool   `json:"av"`
		Tamper     bool   `json:"tamper"`
		QuickScan  string `json:"quickScan"`
		FullScan   string `json:"fullScan"`
		SigVer     string `json:"sigVer"`
		SigUpdated string `json:"sigUpdated"`
		Threats    []struct {
			Name string `json:"name"`
			Time string `json:"time"`
		} `json:"threats"`
	}
	if err := json.Unmarshal(out, &raw); err != nil {
		return Status{}
	}

	st := Status{
		Available:          true,
		RealTimeProtection: raw.RTP,
		AntivirusEnabled:   raw.AV,
		TamperProtection:   raw.Tamper,
		LastQuickScan:      strings.TrimSpace(raw.QuickScan),
		LastFullScan:       strings.TrimSpace(raw.FullScan),
		SignatureVersion:   strings.TrimSpace(raw.SigVer),
		SignatureUpdated:   strings.TrimSpace(raw.SigUpdated),
	}
	for _, t := range raw.Threats {
		name := strings.TrimSpace(t.Name)
		if name == "" {
			continue
		}
		st.RecentThreats = append(st.RecentThreats, Threat{Name: name, Time: strings.TrimSpace(t.Time)})
	}
	return st
}
