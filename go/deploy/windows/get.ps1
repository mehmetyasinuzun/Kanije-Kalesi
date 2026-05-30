# ============================================================
#  Kanije Kalesi - Tek satir web kurulum (Windows)
#
#  Argumansiz (token/chat sorulur):
#    irm https://raw.githubusercontent.com/mehmetyasinuzun/Kanije-Kalesi/master/go/deploy/windows/get.ps1 | iex
#
#  Argumanli (token/chat onceden verilir, hicbir sey sorulmaz):
#    & ([scriptblock]::Create((irm https://raw.githubusercontent.com/mehmetyasinuzun/Kanije-Kalesi/master/go/deploy/windows/get.ps1))) -Token "BOTTOKEN" -Chat "CHATID"
#
#  install.ps1'i indirir, yonetici olarak calistirir; o da hazir
#  kanije.exe'yi GitHub Release'ten indirip gizli kurar.
# ============================================================
param(
    [string]$Token,
    [string]$Chat
)

$ErrorActionPreference = "Stop"
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$base = "https://raw.githubusercontent.com/mehmetyasinuzun/Kanije-Kalesi/master/go/deploy/windows/install.ps1"
$tmp  = Join-Path $env:TEMP "kanije-install.ps1"

Write-Host "Kanije Kalesi kurulum betigi indiriliyor..." -ForegroundColor Cyan
Invoke-WebRequest -UseBasicParsing -Uri $base -OutFile $tmp

$argList = @("-NoProfile", "-ExecutionPolicy", "Bypass", "-File", "`"$tmp`"")
if ($Token) { $argList += @("-Token", "`"$Token`"") }
if ($Chat)  { $argList += @("-Chat", "`"$Chat`"") }

Write-Host "Yonetici olarak calistiriliyor (tek seferlik onay)..." -ForegroundColor Yellow
Start-Process powershell.exe -Verb RunAs -ArgumentList $argList
