//go:build windows

package fileaudit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"os/exec"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

// fileSystemSubcategory is the GUID of the Object Access > File System audit
// subcategory — locale-independent, unlike the display name.
const fileSystemSubcategory = "{0CCE921D-69AE-11D9-BED3-505054503030}"

// runPS runs a PowerShell script (hidden window) and returns stdout, or stderr as
// the error on failure.
func runPS(script string) (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	sysproc.Hide(cmd)

	var out, errBuf strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(errBuf.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("%s", msg)
	}
	return out.String(), nil
}

// EnableAndAudit turns on File System object-access auditing and attaches a SACL
// audit rule (Everyone, read/write/delete, inherited) to path. Requires admin/
// SYSTEM rights — without them auditpol/Set-Acl fail and an error is returned.
func EnableAndAudit(path string) error {
	script := `$ErrorActionPreference='Stop'
auditpol /set /subcategory:"` + fileSystemSubcategory + `" /success:enable /failure:enable | Out-Null
$p = ` + psQuote(path) + `
if (-not (Test-Path -LiteralPath $p)) { throw 'yol bulunamadi' }
$acl = Get-Acl -LiteralPath $p -Audit
$rule = New-Object System.Security.AccessControl.FileSystemAuditRule('Everyone','ReadData,WriteData,Delete,DeleteSubdirectoriesAndFiles','ContainerInherit,ObjectInherit','None','Success,Failure')
$acl.AddAuditRule($rule)
Set-Acl -LiteralPath $p -AclObject $acl
'OK'`

	out, err := runPS(script)
	if err != nil {
		return fmt.Errorf("denetim kurulamadı (yönetici/SYSTEM yetkisi gerekir): %w", err)
	}
	if !strings.Contains(out, "OK") {
		return fmt.Errorf("denetim kurulamadı")
	}
	return nil
}

// DisableAudit removes our SACL audit rules (Everyone) from path. The global audit
// policy is left on (other tools may rely on it).
func DisableAudit(path string) error {
	script := `$ErrorActionPreference='Stop'
$p = ` + psQuote(path) + `
if (Test-Path -LiteralPath $p) {
  $acl = Get-Acl -LiteralPath $p -Audit
  $acl.PurgeAuditRules((New-Object System.Security.Principal.NTAccount('Everyone')))
  Set-Acl -LiteralPath $p -AclObject $acl
}
'OK'`

	if _, err := runPS(script); err != nil {
		return fmt.Errorf("denetim kaldırılamadı (yönetici gerekir): %w", err)
	}
	return nil
}

// RecentAccess reads recent Security-log 4663 events whose object lies under one
// of paths, newest first, capped at limit. Returns who/when/which-program/action.
func RecentAccess(paths []string, limit int) ([]AccessEvent, error) {
	if len(paths) == 0 {
		return nil, nil
	}
	if limit <= 0 {
		limit = 20
	}

	var arr strings.Builder
	for i, p := range paths {
		if i > 0 {
			arr.WriteString(",")
		}
		arr.WriteString(psQuote(p))
	}

	script := `$ErrorActionPreference='SilentlyContinue'
$paths = @(` + arr.String() + `)
$evts = Get-WinEvent -FilterHashtable @{LogName='Security'; Id=4663} -MaxEvents 1000 -ErrorAction SilentlyContinue
$out = @()
foreach ($e in $evts) {
  $x = [xml]$e.ToXml()
  $d = @{}
  foreach ($n in $x.Event.EventData.Data) { $d[$n.Name] = [string]$n.'#text' }
  $obj = $d['ObjectName']
  if (-not $obj) { continue }
  $match = $false
  foreach ($p in $paths) { if ($obj.StartsWith($p, [System.StringComparison]::OrdinalIgnoreCase)) { $match = $true; break } }
  if (-not $match) { continue }
  $out += [PSCustomObject]@{ time="$($e.TimeCreated)"; user=$d['SubjectUserName']; process=$d['ProcessName']; object=$obj; mask=$d['AccessMask'] }
  if ($out.Count -ge ` + strconv.Itoa(limit) + `) { break }
}
ConvertTo-Json -InputObject @($out) -Depth 3 -Compress`

	out, err := runPS(script)
	if err != nil {
		return nil, err
	}
	out = strings.TrimSpace(out)
	if out == "" || out == "null" {
		return nil, nil
	}

	var raw []struct {
		Time    string `json:"time"`
		User    string `json:"user"`
		Process string `json:"process"`
		Object  string `json:"object"`
		Mask    string `json:"mask"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return nil, fmt.Errorf("erişim olayları çözülemedi: %w", err)
	}

	evs := make([]AccessEvent, 0, len(raw))
	for _, r := range raw {
		evs = append(evs, AccessEvent{
			Time:    strings.TrimSpace(r.Time),
			User:    strings.TrimSpace(r.User),
			Process: strings.TrimSpace(r.Process),
			Object:  strings.TrimSpace(r.Object),
			Access:  accessLabel(r.Mask),
		})
	}
	return evs, nil
}

// accessLabel turns a 4663 AccessMask (hex) into a Turkish action label.
func accessLabel(mask string) string {
	m := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(mask)), "0x")
	n, err := strconv.ParseInt(m, 16, 64)
	if err != nil {
		return "erişim"
	}
	var parts []string
	if n&0x10000 != 0 { // DELETE
		parts = append(parts, "silme")
	}
	if n&0x2 != 0 || n&0x4 != 0 { // WriteData / AppendData
		parts = append(parts, "yazma")
	}
	if n&0x1 != 0 { // ReadData / ListDirectory
		parts = append(parts, "okuma/kopya")
	}
	if len(parts) == 0 {
		return "erişim"
	}
	return strings.Join(parts, "+")
}

// psQuote wraps s in a PowerShell single-quoted literal, escaping embedded quotes.
func psQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
