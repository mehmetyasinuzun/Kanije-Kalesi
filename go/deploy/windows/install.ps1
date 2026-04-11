#Requires -Version 5.1
<#
.SYNOPSIS
    Kanije Kalesi - Windows Kurulum Scripti

.DESCRIPTION
    Go binary'sini Windows Task Scheduler'a kaydeder.
    Sisteme giriş yapıldığında otomatik olarak arka planda başlar.
    Konsol penceresi açmaz (-H=windowsgui ile derlenmiştir).

.PARAMETER Remove
    Zamanlanmış görevi kaldırır.

.PARAMETER Status
    Görev durumunu gösterir.

.EXAMPLE
    .\install.ps1              # Kur
    .\install.ps1 -Remove      # Kaldır
    .\install.ps1 -Status      # Durum

.NOTES
    Yönetici hakları gerekmez (kullanıcı Task Scheduler kullanılır).
    Binary'nin aynı dizinde olması veya PATH'de bulunması gerekir.
#>

param(
    [switch]$Remove,
    [switch]$Status
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$TaskName   = "KanijeKalesi"
$TaskDesc   = "Kanije Kalesi güvenlik izleme aracı — otomatik başlatma"
$BinaryName = "kanije.exe"

# Locate the binary
$ScriptDir  = Split-Path -Parent $MyInvocation.MyCommand.Path
$BinaryPath = Join-Path (Split-Path -Parent $ScriptDir) $BinaryName

if (-not (Test-Path $BinaryPath)) {
    # Try PATH
    $found = Get-Command $BinaryName -ErrorAction SilentlyContinue
    if ($found) {
        $BinaryPath = $found.Source
    } else {
        Write-Error "kanije.exe bulunamadı. Binary'nin PATH'de veya proje dizininde olduğundan emin olun."
        exit 1
    }
}

$ConfigPath = Join-Path (Split-Path -Parent $BinaryPath) "config.toml"

if ($Status) {
    $task = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
    if ($task) {
        Write-Host "✅ Görev kayıtlı: $TaskName" -ForegroundColor Green
        Write-Host "   Durum : $($task.State)"
        Write-Host "   Binary: $BinaryPath"
    } else {
        Write-Host "❌ Görev bulunamadı: $TaskName" -ForegroundColor Red
    }
    exit 0
}

if ($Remove) {
    Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false -ErrorAction SilentlyContinue
    Write-Host "✅ Görev kaldırıldı: $TaskName" -ForegroundColor Green
    exit 0
}

# ---- Install ----

Write-Host "🏰 Kanije Kalesi kurulum başlatılıyor..." -ForegroundColor Cyan
Write-Host "   Binary : $BinaryPath"
Write-Host "   Config : $ConfigPath"

# Build task arguments
$Arguments = "start"
if (Test-Path $ConfigPath) {
    $Arguments += " --config `"$ConfigPath`""
}

# Create the task action
$Action  = New-ScheduledTaskAction -Execute $BinaryPath -Argument $Arguments -WorkingDirectory (Split-Path $BinaryPath)

# Trigger: At user logon
$Trigger = New-ScheduledTaskTrigger -AtLogOn

# Settings: restart on failure, allow running on battery, no UI
$Settings = New-ScheduledTaskSettingsSet `
    -ExecutionTimeLimit (New-TimeSpan -Hours 0) `
    -RestartCount 3 `
    -RestartInterval (New-TimeSpan -Minutes 1) `
    -StartWhenAvailable `
    -RunOnlyIfNetworkAvailable:$false `
    -DisallowHardTerminate:$false

# Principal: current user, no elevation needed
$Principal = New-ScheduledTaskPrincipal `
    -UserId $env:USERNAME `
    -LogonType Interactive `
    -RunLevel Limited

# Remove existing task if present
Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false -ErrorAction SilentlyContinue

# Register
Register-ScheduledTask `
    -TaskName   $TaskName `
    -Action     $Action `
    -Trigger    $Trigger `
    -Settings   $Settings `
    -Principal  $Principal `
    -Description $TaskDesc | Out-Null

Write-Host ""
Write-Host "✅ Kurulum tamamlandı!" -ForegroundColor Green
Write-Host ""
Write-Host "Yapılandırma (henüz yapılmadıysa):"
Write-Host "  kanije.exe setup --token <BOT_TOKEN> --chat <CHAT_ID>"
Write-Host ""
Write-Host "Şimdi başlatmak için:"
Write-Host "  Start-ScheduledTask -TaskName $TaskName"
Write-Host ""
Write-Host "Sistemi yeniden başlattığınızda otomatik olarak çalışacak."
