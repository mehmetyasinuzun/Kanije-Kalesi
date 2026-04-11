# ================================================
# Kanije Kalesi - Otomatik Baslatma Kurulumu
# ================================================
# Task Scheduler ile bilgisayar acildiginda
# yonetici yetkisiyle otomatik baslatir.
#
# Kullanim (Yonetici PowerShell gerekli):
#   .\install.ps1          -> Kur
#   .\install.ps1 -Remove  -> Kaldir
#   .\install.ps1 -Status  -> Durum kontrol
# ================================================

param(
    [switch]$Remove,
    [switch]$Status
)

$TaskName = "KanijeKalesi"
$Description = "Kanije Kalesi Guvenlik Izleme"
$AppDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$PythonExe = (Get-Command python -ErrorAction SilentlyContinue).Source
$ScriptPath = Join-Path $AppDir "kanije.py"
$VbsPath = Join-Path $AppDir "kanije_launcher.vbs"

function Test-Admin {
    $identity = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($identity)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# === DURUM KONTROL ===
if ($Status) {
    Write-Host ""
    Write-Host "Kanije Kalesi - Durum Kontrolu" -ForegroundColor Yellow
    Write-Host ""

    $task = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
    if ($task) {
        Write-Host "  [OK] Zamanlanmis gorev kayitli: $TaskName" -ForegroundColor Green
        Write-Host "  Durum: $($task.State)" -ForegroundColor Cyan
        $info = Get-ScheduledTaskInfo -TaskName $TaskName
        Write-Host "  Son calisma: $($info.LastRunTime)" -ForegroundColor Cyan
    } else {
        Write-Host "  [X] Zamanlanmis gorev bulunamadi" -ForegroundColor Red
    }

    if (Test-Path $VbsPath) {
        Write-Host "  [OK] Launcher: $VbsPath" -ForegroundColor Green
    } else {
        Write-Host "  [X] Launcher bulunamadi" -ForegroundColor Red
    }

    if ($PythonExe) {
        Write-Host "  [OK] Python: $PythonExe" -ForegroundColor Green
    } else {
        Write-Host "  [X] Python bulunamadi" -ForegroundColor Red
    }
    Write-Host ""
    exit 0
}

# === KALDIRMA ===
if ($Remove) {
    Write-Host ""
    Write-Host "Kanije Kalesi - Kaldirma" -ForegroundColor Yellow
    Write-Host ""

    if (!(Test-Admin)) {
        Write-Host "  [X] Yonetici yetkisi gerekli!" -ForegroundColor Red
        exit 1
    }

    $task = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
    if ($task) {
        Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false
        Write-Host "  [OK] Gorev silindi: $TaskName" -ForegroundColor Green
    } else {
        Write-Host "  Zaten kayitli degil." -ForegroundColor Cyan
    }

    if (Test-Path $VbsPath) {
        Remove-Item $VbsPath -Force
        Write-Host "  [OK] Launcher silindi" -ForegroundColor Green
    }

    Write-Host ""
    Write-Host "  Kanije Kalesi otomatik baslatmadan kaldirildi." -ForegroundColor Green
    Write-Host ""
    exit 0
}

# === KURULUM ===
Write-Host ""
Write-Host "======================================" -ForegroundColor Yellow
Write-Host "  Kanije Kalesi - Kurulum" -ForegroundColor Yellow
Write-Host "======================================" -ForegroundColor Yellow
Write-Host ""

# 1. Yonetici kontrolu
if (!(Test-Admin)) {
    Write-Host "  [X] Yonetici yetkisi gerekli!" -ForegroundColor Red
    Write-Host "  PowerShell'i Yonetici olarak ac." -ForegroundColor Red
    exit 1
}
Write-Host "  [OK] Yonetici yetkisi var" -ForegroundColor Green

# 2. Python kontrolu
if (!$PythonExe) {
    Write-Host "  [X] Python bulunamadi!" -ForegroundColor Red
    exit 1
}
Write-Host "  [OK] Python: $PythonExe" -ForegroundColor Green

# 3. kanije.py kontrolu
if (!(Test-Path $ScriptPath)) {
    Write-Host "  [X] kanije.py bulunamadi: $ScriptPath" -ForegroundColor Red
    exit 1
}
Write-Host "  [OK] Script: $ScriptPath" -ForegroundColor Green

# 4. config.yaml kontrolu
$configPath = Join-Path $AppDir "config.yaml"
if (!(Test-Path $configPath)) {
    Write-Host "  [X] config.yaml bulunamadi!" -ForegroundColor Red
    exit 1
}
Write-Host "  [OK] Config: $configPath" -ForegroundColor Green

# 5. VBS launcher olustur (konsol penceresi gostermeden calistirir)
$escapedAppDir = $AppDir.Replace('\','\\')
$escapedScript = $ScriptPath.Replace('\','\\')

$vbsLines = @(
    "Set WshShell = CreateObject(""WScript.Shell"")"
    "WshShell.CurrentDirectory = ""$escapedAppDir"""
    "WshShell.Run ""pythonw """"$escapedScript"""" start"", 0, False"
)
$vbsLines | Out-File -FilePath $VbsPath -Encoding ASCII
Write-Host "  [OK] VBS launcher olusturuldu" -ForegroundColor Green

# 6. Eski gorevi sil (varsa)
$existing = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
if ($existing) {
    Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false
    Write-Host "  Eski gorev silindi, yenisi olusturuluyor..." -ForegroundColor Cyan
}

# 7. Task Scheduler gorevi olustur
$action = New-ScheduledTaskAction `
    -Execute "wscript.exe" `
    -Argument """$VbsPath""" `
    -WorkingDirectory $AppDir

$trigger = New-ScheduledTaskTrigger -AtLogOn

$principal = New-ScheduledTaskPrincipal `
    -UserId $env:USERNAME `
    -RunLevel Highest `
    -LogonType Interactive

$settings = New-ScheduledTaskSettingsSet `
    -AllowStartIfOnBatteries `
    -DontStopIfGoingOnBatteries `
    -StartWhenAvailable `
    -RestartCount 3 `
    -RestartInterval (New-TimeSpan -Minutes 1) `
    -ExecutionTimeLimit (New-TimeSpan -Days 0)

Register-ScheduledTask `
    -TaskName $TaskName `
    -Action $action `
    -Trigger $trigger `
    -Principal $principal `
    -Settings $settings `
    -Description $Description `
    -Force | Out-Null

Write-Host "  [OK] Task Scheduler gorevi olusturuldu" -ForegroundColor Green

# 8. Sonuc
Write-Host ""
Write-Host "======================================" -ForegroundColor Green
Write-Host "  Kurulum tamamlandi!" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green
Write-Host ""
Write-Host "  Gorev adi: $TaskName" -ForegroundColor Cyan
Write-Host "  Tetikleyici: Kullanici oturum actiginda" -ForegroundColor Cyan
Write-Host "  Yetki: Yonetici (en yuksek)" -ForegroundColor Cyan
Write-Host "  Cokmede: 3x otomatik yeniden baslatma" -ForegroundColor Cyan
Write-Host "  Pencere: Gizli (konsol yok)" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Bilgisayarini yeniden baslat - otomatik acilacak." -ForegroundColor Yellow
Write-Host ""
Write-Host "  Kaldirmak icin: .\install.ps1 -Remove" -ForegroundColor DarkGray
Write-Host "  Durum kontrolu: .\install.ps1 -Status" -ForegroundColor DarkGray
Write-Host ""
