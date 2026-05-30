#Requires -Version 5.1
<#
.SYNOPSIS
    Kanije Kalesi - Windows tek-tik, gizli kurulum.

.DESCRIPTION
    Tek bir UAC onayiyla her seyi yapar:
      1. Kendini yonetici olarak yukseltir (gereken tum izinler bastan alinir).
      2. kanije.exe'yi bulur -> yoksa GitHub Release'ten indirir -> o da yoksa Go ile derler.
      3. Telegram token + chat ID'yi alir (parametre verilmezse sorar).
      4. Gizli (pencere/ikon yok) ve YUKSELTILMIS bir zamanlanmis gorev kurar.
         Yukseltilmis gorev sayesinde Guvenlik olay gunlugu (giris denemeleri)
         okunabilir ve bir daha UAC sorulmaz.
      5. Ajani hemen baslatir.

    Kurulduktan sonra masaustunde HICBIR sey gorunmez. Yonetim Telegram'dan yapilir.

.PARAMETER Token   Telegram bot token (verilmezse sorulur)
.PARAMETER Chat    Telegram chat ID (verilmezse sorulur)
.PARAMETER Remove  Kurulumu kaldirir (config korunur)
.PARAMETER Status  Gorev durumunu gosterir
.PARAMETER NoStart Kurar ama hemen baslatmaz

.EXAMPLE
    .\install.ps1
    .\install.ps1 -Token "123:ABC" -Chat "987654321"
    .\install.ps1 -Remove
#>
param(
    [string]$Token,
    [string]$Chat,
    [switch]$Remove,
    [switch]$Status,
    [switch]$NoStart
)

$ErrorActionPreference = "Stop"
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

# Turkce karakterler konsolda dogru gorunsun.
try { [Console]::OutputEncoding = [System.Text.Encoding]::UTF8 } catch { }

# ---- Sabitler ----
$TaskName   = "KanijeKalesi"
$AppName    = "KanijeKalesi"
$InstallDir = Join-Path $env:LOCALAPPDATA $AppName
$BinaryPath = Join-Path $InstallDir "kanije.exe"
$ConfigPath = Join-Path $InstallDir "config.toml"
$RepoOwner  = "mehmetyasinuzun"
$RepoName   = "Kanije-Kalesi"
$AssetName  = "kanije-windows-amd64.exe"

function Write-Step($msg)    { Write-Host "  > $msg" -ForegroundColor Cyan }
function Write-Ok($msg)      { Write-Host "[OK] $msg" -ForegroundColor Green }
function Write-WarnMsg($msg) { Write-Host "[!]  $msg" -ForegroundColor Yellow }

# ---- Yonetici kontrolu ve kendini yukseltme ----
function Test-Admin {
    $id = [Security.Principal.WindowsIdentity]::GetCurrent()
    (New-Object Security.Principal.WindowsPrincipal($id)).IsInRole(
        [Security.Principal.WindowsBuiltInRole]::Administrator)
}

if (-not (Test-Admin)) {
    if (-not $PSCommandPath) {
        Write-WarnMsg "Bu betik bir dosyadan calistirilmali. Yonetici PowerShell'de tekrar deneyin."
        exit 1
    }
    Write-Host "Yonetici haklari isteniyor (tek seferlik)..." -ForegroundColor Yellow
    $argList = @("-NoProfile", "-ExecutionPolicy", "Bypass", "-File", "`"$PSCommandPath`"")
    if ($Token)   { $argList += @("-Token", "`"$Token`"") }
    if ($Chat)    { $argList += @("-Chat", "`"$Chat`"") }
    if ($Remove)  { $argList += "-Remove" }
    if ($Status)  { $argList += "-Status" }
    if ($NoStart) { $argList += "-NoStart" }
    Start-Process powershell.exe -Verb RunAs -ArgumentList $argList
    exit
}

# ---- Durum ----
if ($Status) {
    $task = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
    if ($task) {
        Write-Ok "Kurulu: $TaskName ($($task.State))"
        Write-Host "   Binary: $BinaryPath"
        Write-Host "   Config: $ConfigPath"
    } else {
        Write-WarnMsg "Kurulu degil."
    }
    Read-Host "`nKapatmak icin Enter"
    exit 0
}

# ---- Kaldirma ----
if ($Remove) {
    Stop-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
    Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false -ErrorAction SilentlyContinue
    Get-Process -Name "kanije" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
    Write-Ok "Kaldirildi. Yapilandirma korundu: $ConfigPath"
    Read-Host "`nKapatmak icin Enter"
    exit 0
}

Write-Host ""
Write-Host "=== Kanije Kalesi - Kurulum ===" -ForegroundColor Cyan
Write-Host ""

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

# Mevcut kurulumu durdur — hem binary kilidini açar (guncelleme icin) hem temiz baslar.
Stop-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue | Out-Null
Get-Process -Name "kanije" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Milliseconds 500

# ---- 1) Binary'yi temin et ----
# Oncelik: repodaki yerel derleme -> en son GitHub Release (guncelleme) -> mevcut kurulu (cevrimdisi) -> Go ile derle.
function Resolve-Binary {
    $scriptDir = Split-Path -Parent $PSCommandPath
    $goRoot    = Split-Path -Parent (Split-Path -Parent $scriptDir)  # .../go

    # a) Repoda yerel derleme (gelistirici / cift-tik yolu)
    $candidates = @(
        (Join-Path $goRoot "kanije.exe"),
        (Join-Path $goRoot "dist\$AssetName"),
        (Join-Path $scriptDir "kanije.exe")
    )
    foreach ($c in $candidates) {
        if (Test-Path $c) {
            Copy-Item $c $BinaryPath -Force
            Write-Ok "Yerel binary kullanildi: $c"
            return $true
        }
    }

    # b) En son surumu indir (tek-satir yolu + her calistirmada guncelleme)
    Write-Step "En son binary indiriliyor (GitHub Release)..."
    try {
        $rel = Invoke-RestMethod -Uri "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest" `
            -Headers @{ "User-Agent" = "kanije-installer" }
        $asset = $rel.assets | Where-Object { $_.name -eq $AssetName } | Select-Object -First 1
        if ($asset) {
            Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $BinaryPath -UseBasicParsing
            Write-Ok "Indirildi: $($rel.tag_name)"
            return $true
        }
        Write-WarnMsg "Release varligi bulunamadi: $AssetName"
    } catch {
        Write-WarnMsg "Indirme basarisiz: $($_.Exception.Message)"
    }

    # c) Internet yoksa, zaten kurulu olani kullan
    if (Test-Path $BinaryPath) {
        Write-Ok "Mevcut binary kullaniliyor (cevrimdisi yedek)"
        return $true
    }

    # d) Go ile derle
    if (Get-Command go -ErrorAction SilentlyContinue) {
        if (Test-Path (Join-Path $goRoot "cmd\kanije")) {
            Write-Step "Go ile derleniyor..."
            Push-Location $goRoot
            try {
                & go build -ldflags="-s -w -H=windowsgui" -o $BinaryPath ./cmd/kanije/
                if (Test-Path $BinaryPath) { Write-Ok "Derlendi"; return $true }
            } finally { Pop-Location }
        }
    }

    return $false
}

if (-not (Resolve-Binary)) {
    Write-WarnMsg "kanije.exe temin edilemedi (indirme ve derleme basarisiz)."
    Write-Host "   Internet baglantinizi kontrol edin veya Go kurun: https://go.dev/dl/"
    Read-Host "`nKapatmak icin Enter"
    exit 1
}

# ---- 2) Telegram bilgilerini al ----
if (-not (Test-Path $ConfigPath)) {
    if (-not $Token) {
        Write-Host ""
        Write-Host "Telegram bot bilgileri gerekli:" -ForegroundColor Yellow
        Write-Host "  - Token: @BotFather -> /newbot"
        Write-Host "  - Chat ID: @userinfobot -> /start"
        Write-Host ""
        $Token = Read-Host "Bot Token"
    }
    if (-not $Chat) {
        $Chat = Read-Host "Chat ID"
    }

    Write-Step "Yapilandirma kaydediliyor..."
    & $BinaryPath setup --token "$Token" --chat "$Chat" --config "$ConfigPath" | Out-Null
    if (-not (Test-Path $ConfigPath)) {
        Write-WarnMsg "Yapilandirma olusturulamadi (token/chat gecersiz olabilir)."
        Read-Host "`nKapatmak icin Enter"
        exit 1
    }

    Write-Step "Telegram baglantisi test ediliyor..."
    & $BinaryPath test --config "$ConfigPath"
} else {
    Write-Ok "Mevcut yapilandirma kullaniliyor: $ConfigPath"
}

# ---- 3) Gizli + yukseltilmis zamanlanmis gorev ----
Write-Step "Otomatik baslatma gorevi kuruluyor (gizli, yukseltilmis)..."

$Action = New-ScheduledTaskAction -Execute $BinaryPath `
    -Argument "start --config `"$ConfigPath`"" -WorkingDirectory $InstallDir

$Trigger = New-ScheduledTaskTrigger -AtLogOn

$Settings = New-ScheduledTaskSettingsSet `
    -AllowStartIfOnBatteries `
    -DontStopIfGoingOnBatteries `
    -StartWhenAvailable `
    -RestartCount 3 `
    -RestartInterval (New-TimeSpan -Minutes 1) `
    -ExecutionTimeLimit (New-TimeSpan -Seconds 0) `
    -Hidden

# RunLevel Highest -> yukseltilmis calisir, calisma zamaninda UAC sorulmaz.
$Principal = New-ScheduledTaskPrincipal `
    -UserId "$env:USERDOMAIN\$env:USERNAME" `
    -LogonType Interactive `
    -RunLevel Highest

Register-ScheduledTask -TaskName $TaskName -Action $Action -Trigger $Trigger `
    -Settings $Settings -Principal $Principal `
    -Description "Kanije Kalesi - gizli guvenlik izleme" -Force | Out-Null

Write-Ok "Gorev kuruldu (oturum acilisinda otomatik, gizli)"

# ---- 4) Baslat ----
if (-not $NoStart) {
    Start-ScheduledTask -TaskName $TaskName
    Write-Ok "Ajan baslatildi - arka planda, gorunmez calisiyor"
}

Write-Host ""
Write-Host "===============================================" -ForegroundColor Green
Write-Ok "Kurulum tamamlandi!"
Write-Host ""
Write-Host "  - Masaustunde hicbir ikon/pencere gorunmez."
Write-Host "  - Telegram botunuza /kurulum yazarak ayarlayin."
Write-Host "  - Kaldirmak icin: install.ps1 -Remove"
Write-Host "===============================================" -ForegroundColor Green
Start-Sleep -Seconds 1
