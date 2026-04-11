# 🏰 WINDOWS 10 — "KALE" MİMARİSİ: TAM SPEKTRUM SERTLEŞTİRME REHBERİ

> **Hedef:** Donanım silikon katmanından uygulama belleğine kadar her katmanı kilitleyen, fiziksel çalınma dahil tüm senaryolara karşı dayanıklı bir Windows 10 sistemi inşa etmek.

> [!NOTE]
> Bu rehber **Windows 10 Pro / Enterprise (22H2)** için hazırlanmıştır. Home sürümünde `gpedit.msc` (Grup İlkesi Düzenleyici) yoktur. Home kullanıcıları için alternatif Registry yolları ilgili bölümlerde ayrıca belirtilmiştir.

---

## 📐 MİMARİ HARITA — Saldırı Yüzeyi ve Savunma Katmanları

```
┌─────────────────────────────────────────────────────────┐
│  KATMAN 0 — Fiziksel / Donanım                          │
│  BIOS Admin Şifresi · Secure Boot · TPM · DMA          │
├─────────────────────────────────────────────────────────┤
│  KATMAN 1 — Önyükleme (Boot)                            │
│  BitLocker Pre-Boot PIN · UEFI Boot Order              │
├─────────────────────────────────────────────────────────┤
│  KATMAN 2 — İşletim Sistemi Çekirdeği                   │
│  VBS · HVCI · Credential Guard · WDAC · LSA PPL        │
├─────────────────────────────────────────────────────────┤
│  KATMAN 3 — Bellek (RAM)                                │
│  Pagefile Temizleme · Hibernate Kapatma · DMA Koruması  │
├─────────────────────────────────────────────────────────┤
│  KATMAN 4 — Süreç / Uygulama                            │
│  Exploit Guard · ASLR · DEP · CFG · ASR Kuralları      │
├─────────────────────────────────────────────────────────┤
│  KATMAN 5 — Ağ                                          │
│  SMB Hardening · Firewall · WPAD/LLMNR/NetBIOS Kapatma │
├─────────────────────────────────────────────────────────┤
│  KATMAN 6 — Kimlik / Hesap                              │
│  Windows Hello · Standart Hesap · UAC · LSA            │
├─────────────────────────────────────────────────────────┤
│  KATMAN 7 — Denetim / İzleme + Telemetri Kapatma       │
│  Gelişmiş Denetim · Sysmon · PS Günlük · Telemetri     │
└─────────────────────────────────────────────────────────┘
```

> [!IMPORTANT]
> Windows 10, Windows 11'e kıyasla bazı önemli farklılıklar taşır:
> - **Credential Guard** yalnızca **Enterprise** lisansında varsayılan desteklidir; Pro'da manuel etkinleştirilir.
> - **DNS over HTTPS (DoH)** Windows 10'da işletim sistemi seviyesinde yerleşik ayar olarak gelmez (sadece Win 11'de var); 3. taraf araç veya router seviyesinde sağlanır.
> - **Smart App Control** Win 10'da yoktur; karşılığı WDAC'tır.
> - **Telemetri** Win 10'da çok daha agresiftir; kapatılması ayrı bir bölüm gerektirir.

---

## KATMAN 0 — Donanım ve BIOS/UEFI Kilitleme

> [!IMPORTANT]
> Bu adım, işletim sistemi kurulmadan **önce** yapılmalıdır. BIOS'u kilitlemeyen bir sistem, OS seviyesindeki tüm önlemleri atlatılabilir hale getirir.

### 0.1 — BIOS/UEFI Şifreleri

Bilgisayar açılır açılmaz **F2 / Del / F10** (üreticiye göre değişir) tuşuna basarak BIOS'a gir.

| Ayar | Değer | Neden? |
|------|-------|--------|
| **Supervisor / Admin Password** | Güçlü bir şifre belirle | BIOS ayarlarını değiştirmeyi engeller |
| **User / Boot Password** | Güçlü bir şifre belirle | POST sonrası disk önyüklemesini kilitler |
| **HDD Password** (destekleniyorsa) | Güçlü bir şifre belirle | Disk başka makineye takılsa bile açılmaz |

> [!CAUTION]
> BIOS şifreni kaybedersen pil sıfırlaması veya servis müdahalesi gerekebilir. Şifreyi fiziksel, dijital olmayan bir yerde sakla.

### 0.2 — Secure Boot

```
BIOS → Security / Boot → Secure Boot → Enabled
Mod → "Standard" veya "Windows UEFI" (Custom değil)
```

Secure Boot, imzasız önyükleyicilerin, bootkitlerin ve rootkitlerin başlamasını engeller. Windows 10 için imzalı güncel sürücüler kullanıyorsan hiçbir şey bozulmaz.

### 0.3 — TPM

```
BIOS → Security → TPM → Enabled
```

Windows 10 kurulum için TPM zorunlu değildir (bu Windows 11'e özgü gereksinim) ama BitLocker ve VBS için **TPM 2.0** şiddetle önerilir.

- Firmwareyı kontrol et: `tpm.msc` → TPM versiyonunu gör.
- **TPM 1.2** varsa BitLocker çalışır fakat Pre-Boot PIN ek yapılandırma ister; TPM 2.0 tercih edilir.

### 0.4 — Sanallaştırma (HVCI ve VBS için Zorunlu)

```
BIOS → CPU / Advanced → Intel VT-x (veya AMD-V) → Enabled
BIOS → CPU / Advanced → VT-d / AMD IOMMU         → Enabled
```

> [!NOTE]
> VT-d / AMD IOMMU olmadan Kernel DMA Koruması çalışmaz. Sadece VT-x açmak yeterli değildir.

### 0.5 — I/O Port Yönetimi

| Port | Önerim |
|------|--------|
| Thunderbolt / USB | BIOS'ta "Authorized Only" veya "Disabled" |
| Seri / Paralel Port | Disabled |
| Wake on LAN / PXE Boot | Disabled (kullanmıyorsan) |
| Boot from USB/CD | Disabled veya şifreli bırak |

---

## KATMAN 1 — Boot Kilitleme: BitLocker Pre-Boot PIN

### 1.1 — Neden Yalnızca TPM Yetmez?

Varsayılan BitLocker = **yalnızca TPM** → Bilgisayar açılır, hiç şifre sormaz, disk anahtarı RAM'de açıkta durur.
**TPM + PIN kombinasyonu zorunludur.**

### 1.2 — Grup İlkesi ile PIN'i Zorunlu Kılma

**`Win + R` → `gpedit.msc`**

```
Bilgisayar Yapılandırması
  └─ Yönetim Şablonları
       └─ Windows Bileşenleri
            └─ BitLocker Sürücü Şifrelemesi
                 └─ İşletim Sistemi Sürücüleri
```

| Politika | Ayar |
|----------|------|
| **Başlangıçta ek kimlik doğrulaması iste** | **Etkin** |
| *Uyumlu TPM olmadan BitLocker'a izin ver* | İşaretli |
| *Başlangıçta TPM ile...* | **"TPM ile başlangıç PIN'ini zorunlu kıl"** |
| **BitLocker şifrelemesi için sürücü şifreleme türü** | **Etkin → Tam şifreleme** |
| **İşletim Sistemi için şifre gücünü yapılandır** | **Etkin → XTS-AES 256-bit** |

> [!NOTE]
> **Home sürümü için (gpedit yoksa) Registry alternatifi:**
> ```powershell
> # TPM+PIN zorunlu kıl
> $Path = "HKLM:\SOFTWARE\Policies\Microsoft\FVE"
> New-Item -Path $Path -Force | Out-Null
> Set-ItemProperty -Path $Path -Name "UseAdvancedStartup"            -Value 1
> Set-ItemProperty -Path $Path -Name "EnableBDEWithNoTPM"            -Value 0
> Set-ItemProperty -Path $Path -Name "UseTPMPIN"                     -Value 1  # TPM+PIN zorunlu
> Set-ItemProperty -Path $Path -Name "EncryptionMethodWithXtsFdvDropDown" -Value 7  # XTS-AES 256
> ```

### 1.3 — BitLocker'ı Aktifleştirme

```powershell
# Yönetici PowerShell

# TPM durumu
Get-Tpm

# BitLocker'ı TPM+PIN ile etkinleştir
$pin = Read-Host -AsSecureString "BitLocker Pre-Boot PIN girin"
Enable-BitLocker -MountPoint "C:" `
    -EncryptionMethod XtsAes256 `
    -TPMandPinProtector `
    -Pin $pin `
    -UsedSpaceOnly:$false    # Tüm disk (boş alanı da şifrele)

# Kurtarma anahtarını ekle ve görüntüle — KASAYA YEDEKLE
Add-BitLockerKeyProtector -MountPoint "C:" -RecoveryPasswordProtector
(Get-BitLockerVolume -MountPoint "C:").KeyProtector |
    Where-Object { $_.KeyProtectorType -eq "RecoveryPassword" }
```

> [!CAUTION]
> Kurtarma anahtarını Microsoft hesabına veya OneDrive'a gönderme. Çıktı al, laminele, fiziksel kasada sakla.

### 1.4 — Doğrulama

```powershell
Get-BitLockerVolume -MountPoint "C:" |
    Select-Object VolumeStatus, EncryptionMethod, EncryptionPercentage, KeyProtector

# EncryptionPercentage  : 100
# VolumeStatus          : FullyEncrypted
# EncryptionMethod      : XtsAes256
# KeyProtector          : {Tpm, TpmPin, RecoveryPassword}
```

---

## KATMAN 2 — Çekirdek Yalıtımı (Kernel Hardening)

### 2.1 — VBS (Virtualization-Based Security)

VBS, CPU sanallaştırmasını kullanarak yalıtılmış bir "Güvenli Dünya" oluşturur. Kimlik bilgileri ve çekirdek bütünlük kodu bu ayrı bölgede korunur.

**Grup İlkesi:**
```
Bilgisayar Yapılandırması
  └─ Yönetim Şablonları
       └─ Sistem
            └─ Device Guard
                 └─ Windows Defender Credential Guard'ı Aç
                    → Etkin
                    → Seçenek: "UEFI Kilidiyle Etkin"
```

**Registry alternatifi:**
```powershell
$DGPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows\DeviceGuard"
New-Item -Path $DGPath -Force | Out-Null
Set-ItemProperty -Path $DGPath -Name "EnableVirtualizationBasedSecurity"     -Value 1
Set-ItemProperty -Path $DGPath -Name "RequirePlatformSecurityFeatures"        -Value 3  # Secure Boot + DMA
Set-ItemProperty -Path $DGPath -Name "HypervisorEnforcedCodeIntegrity"        -Value 1  # HVCI
Set-ItemProperty -Path $DGPath -Name "HVCIMATRequired"                        -Value 0
Set-ItemProperty -Path $DGPath -Name "LsaCfgFlags"                            -Value 1  # Credential Guard
```

**Doğrulama (yeniden başlatma sonrası):**
```powershell
$DG = Get-CimInstance -ClassName Win32_DeviceGuard `
      -Namespace root\Microsoft\Windows\DeviceGuard
$DG.VirtualizationBasedSecurityStatus
# 2 = Etkin (çalışıyor)
```

### 2.2 — HVCI (Bellek Bütünlüğü)

Windows 10'da Windows Güvenliği → Cihaz Güvenliği → Çekirdek Yalıtımı menüsü **1803 (April 2018 Update)** ve sonrasında mevcuttur.

```
Başlat → Windows Güvenliği → Cihaz Güvenliği
  → Çekirdek Yalıtımı Ayrıntıları
     → Bellek Bütünlüğü: AÇIK
```

**Registry (veya otomatikleştirmek için):**
```powershell
Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\DeviceGuard\Scenarios\HypervisorEnforcedCodeIntegrity" `
    -Name "Enabled" -Value 1 -Type DWord
```

**Ne sağlar?** Ring-0 kernel exploit'leri, değiştirilmiş sürücü enjeksiyonları ve rootkitler hypervisor katmanında engellenir.

**Doğrulama:**
```powershell
$DG = Get-CimInstance -ClassName Win32_DeviceGuard `
      -Namespace root\Microsoft\Windows\DeviceGuard
$DG.SecurityServicesRunning
# 2 içermeli = HVCI aktif
```

### 2.3 — Credential Guard

> [!WARNING]
> **Windows 10 Home:** Credential Guard **desteklenmez.**
> **Windows 10 Pro:** Manuel etkinleştirme gerekir (aşağıdaki adımlar).
> **Windows 10 Enterprise:** Varsayılan etkin olabilir; teyit et.

**Ön koşullar (Pro için):**
- UEFI + Secure Boot açık olmalı
- VBS etkin olmalı
- TPM 2.0 (önerilir; 1.2 ile de çalışır ama daha zayıf)
- 64-bit işlemci

**Etkinleştirme (Group Policy):**
```
Bilgisayar Yapılandırması → Yönetim Şablonları → Sistem → Device Guard
  → "Windows Defender Credential Guard'ı Aç" → Etkin → UEFI Kilidiyle Etkin
```

**Doğrulama:**
```powershell
$DG = Get-CimInstance -ClassName Win32_DeviceGuard `
      -Namespace root\Microsoft\Windows\DeviceGuard
$DG.SecurityServicesRunning
# 1 içermeli = Credential Guard aktif
```

**Ne sağlar?** NTLM hash ve Kerberos biletleri, `lsass.exe` yerine sanal yalıtılmış `LsaIso.exe` sürecinde tutulur. Mimikatz / Pass-the-Hash / Pass-the-Ticket saldırıları etkisiz kalır.

### 2.4 — LSA Koruması (RunAsPPL)

`lsass.exe` sürecini **Protected Process Light (PPL)** olarak işaretle.

```powershell
# Regedit veya PowerShell (yeniden başlatma gerekir)
Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Lsa" `
    -Name "RunAsPPL" -Value 1 -Type DWord
# Değer 1 = PPL
# Değer 2 = PPL + UEFI Kilidi (Win 11'e özgü; Win 10'da 1 kullan)

# Doğrulama:
Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\Lsa" -Name RunAsPPL
```

**Ek LSA sertleştirme:**
```powershell
# WDigest kimlik bilgilerini bellekte tutmayı kapat
# (Mimikatz'ın düz metin şifre çalmasını engeller)
Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\SecurityProviders\WDigest" `
    -Name "UseLogonCredential" -Value 0

# Doğrulama:
Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\SecurityProviders\WDigest" `
    -Name UseLogonCredential
# 0 olmalı
```

### 2.5 — WDAC (Windows Defender Application Control)

Windows 10'da WDAC **1709 (Fall Creators Update)** ile geldi. Uygulama beyaz listesi — yalnızca izin verilen uygulamalar çalışır.

```powershell
# Denetim modunda başla — önce ne kırılıyor gör

$PolicyPath = "$env:TEMP\WDACPolicy.xml"
New-CIPolicy -FilePath $PolicyPath `
    -Level Publisher `
    -Fallback Hash `
    -UserPEs `
    -MultiplePolicyFormat

# XML'i binary'e çevir
$BinPath = "$env:TEMP\WDACPolicy.bin"
ConvertFrom-CIPolicy -XmlFilePath $PolicyPath -BinaryFilePath $BinPath

# Sisteme uygula
Copy-Item $BinPath "$env:windir\System32\CodeIntegrity\SIPolicy.p7b"
```

> [!WARNING]
> WDAC'ı hemen **zorunlu moda** (enforce) alma. Önce denetim modunda (audit) günlükleri **birkaç gün** izle:
> ```
> Event Viewer → Applications and Services Logs
>   → Microsoft → Windows → CodeIntegrity → Operational
> ```
> Engellenecek uygulamaları whitelist'e ekledikten sonra enforce et.

---

## KATMAN 3 — RAM Güvenliği

> [!IMPORTANT]
> RAM, güç kesildikten **saniyeler** (soğutulursa dakikalar) sonra veriyi fiziksel olarak tutar — buna **Cold Boot Attack** denir. Bu katman bu tehdit modeline karşıdır.

### 3.1 — Pagefile.sys Kapanışta Temizle

```powershell
Set-ItemProperty `
    -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management" `
    -Name "ClearPageFileAtShutdown" `
    -Value 1 -Type DWord

# Doğrulama:
Get-ItemProperty `
    -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management" `
    -Name "ClearPageFileAtShutdown"
# 1 olmalı
```

> [!NOTE]
> Bu, kapatma süresini 10–30 saniye uzatır. Tamamen kapatmak (pagefile = none) yerine bu yol daha güvenli — pagefile tamamen kalkınca RAM dolu olduğunda uygulamalar çöker.

### 3.2 — Hiberfil.sys Kapat (Hazırda Bekletme)

`hiberfil.sys`, RAM içeriğinin **bire bir kopyasını** diske yazar — BitLocker açık olsa bile oturum anahtarları bu dosyada bulunabilir.

```powershell
# Yönetici PowerShell
powercfg /hibernate off

# Doğrulama:
Test-Path "C:\hiberfil.sys"
# False dönmeli
```

> [!WARNING]
> Hibernate kapatılınca **Hızlı Başlatma (Fast Startup)** da devre dışı kalır. Açılış 2–4 saniye uzar. Kabul et — güvenlik kazanımı buna değer.

### 3.3 — Uyku (Sleep) Modunu Kısıtla

```powershell
# Uyku modunda RAM güç almaya devam eder — Cold Boot riski düşüktür ama sıfır değil
# En güvenli: Uyku yerine "Kapat" tercih et
powercfg /change standby-timeout-ac 0
powercfg /change standby-timeout-dc 0
```

### 3.4 — Çekirdek DMA Koruması

Thunderbolt / PCIe üzerinden DMA saldırısıyla RAM içeriği fiziksel erişimle saniyeler içinde çekilebilir.

Windows 10'da Kernel DMA Koruması **1803** ile geldi fakat yalnızca **UEFI + VT-d/IOMMU** açık sistemlerde çalışır.

```powershell
# Doğrulama:
$DG = Get-CimInstance -ClassName Win32_DeviceGuard `
      -Namespace root\Microsoft\Windows\DeviceGuard
$DG.KernelDmaProtectionPolicy
# 1 = Etkin
```

Eğer `0` dönüyorsa: BIOS'ta **VT-d (Intel)** veya **AMD IOMMU** açılmalı.

### 3.5 — Bellekten WDigest Kimlik Bilgilerini Kaldır

Yukarıdaki Katman 2.4'te yapıldı (`UseLogonCredential = 0`). Mimikatz'ın RAM'den düz metin şifre okumasını engeller.

---

## KATMAN 4 — Süreç ve Uygulama Koruması

### 4.1 — Exploit Protection Yapılandırması

```
Başlat → Windows Güvenliği → Uygulama ve Tarayıcı Denetimi
  → Exploit Protection Ayarları → Sistem Ayarları
```

| Koruma | Ayar | Ne Sağlar? |
|--------|------|------------|
| **CFG** (Control Flow Guard) | Açık (varsayılan) | Fonksiyon çağrı akışını doğrular |
| **DEP** (Data Execution Prevention) | Açık (varsayılan) | Veri bölgelerinde kod çalıştırmayı engeller |
| **SEHOP** | Açık | SEH istisna zinciri manipülasyonunu engeller |
| **Heap Randomization** | Açık | Heap adreslerini rastgele düzenler |
| **Zorunlu ASLR** | **Açık yap** | ASLR desteklemeyen modülleri de rastgele adrese koyar |

**Zorunlu ASLR PowerShell ile:**
```powershell
Set-ProcessMitigation -System -Enable ForceRelocateImages

# Tarayıcı gibi kritik uygulamaya ek koruma:
Set-ProcessMitigation -Name msedge.exe `
    -Enable DEP, SEHOP, ForceRelocateImages, BottomUp, StrictHandle, `
            DisableWin32kSystemCalls, BlockNonMicrosoftFonts, DisableExtensionPoints

# Politikayı yedekle:
Get-ProcessMitigation -RegistryConfigFilePath "C:\Guvenlik\ExploitProtection_W10.xml"
```

> [!WARNING]
> `ProhibitDynamicCode` eski uygulamaları kırabilir. Her uygulamada önce test et.

### 4.2 — Attack Surface Reduction (ASR) Kuralları

ASR, Windows Defender içinde gömülü ama çoğu kişinin açmadığı güçlü bir katmandır. Windows 10 **1709** itibarıyla desteklenir.

```powershell
# Yönetici PowerShell — Tüm temel ASR kurallarını etkinleştir

$ASRRules = @(
    "BE9BA2D9-53EA-4CDC-84E5-9B1EEEE46550",  # Office makrolarının Win32 API çağırmasını engelle
    "D4F940AB-401B-4EFC-AADC-AD5F3C50688A",  # Office'in alt süreç oluşturmasını engelle
    "3B576869-A4EC-4529-8536-B80A7769E899",  # Office'in yürütülebilir içerik oluşturmasını engelle
    "75668C1F-73B5-4CF0-BB93-3ECF5CB7CC84",  # Office uygulamalarının süreçlere enjeksiyonunu engelle
    "D3E037E1-3EB8-44C8-A917-57927947596D",  # JS/VBS'nin indirilen içeriği yürütmesini engelle
    "5BEB7EFE-FD9A-4556-801D-275E5FFC04CC",  # Karartılmış betiklerin yürütülmesini engelle
    "92E97FA1-2EDF-4476-BDD6-9DD0B4DDDC7B",  # Win32 API çağrılarını Office makrolarından engelle
    "01443614-CD74-433A-B99E-2ECDC07BFC25",  # Fidye yazılımı belirtisi taşıyan süreçleri engelle
    "C1DB55AB-C21A-4637-BB3F-A12568109D35",  # Fidye yazılımı davranışını engelle
    "9E6C4E1F-7D60-472F-BA1A-A39EF669E4B2",  # LSA kimlik bilgisi çalmayı engelle
    "D1E49AAC-8F56-4280-B9BA-993A6D77406C",  # PsExec ve WMI'dan süreç oluşturulmasını engelle
    "B2B3F03D-6A65-4F7B-A9C7-1C7EF74A9BA4",  # Eposta/web postasından yürütülebilir içeriği engelle
    "26190899-1602-49E8-8B27-EB1D0A1CE869",  # Office uygulamalarından iletişim uygulamaları oluşturmayı engelle
    "7674BA52-37EB-4A4F-A9A1-F0F9A1619A2C",  # Adobe Reader alt süreç oluşturmasını engelle
    "E6DB77E5-3DF2-4CF1-B95A-636979351E5B"   # USB'den güvenilmez/imzasız süreçleri engelle
)

foreach ($Rule in $ASRRules) {
    Add-MpPreference -AttackSurfaceReductionRules_Ids $Rule `
                     -AttackSurfaceReductionRules_Actions Enabled
}

# Doğrulama:
Get-MpPreference | Select-Object AttackSurfaceReductionRules_Ids, AttackSurfaceReductionRules_Actions
```

### 4.3 — PowerShell Kısıtlama Politikası

```powershell
# PowerShell yürütme politikasını kısıtla
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope LocalMachine
# RemoteSigned: yerel scriptler çalışır, internetten gelenler imzalı olmalı

# Constrained Language Mode (CLM) — daha agresif kısıtlama
# (WDAC aktifken otomatik devreye girer)
# Manuel test:
$ExecutionContext.SessionState.LanguageMode
# "FullLanguage" yerine "ConstrainedLanguage" görünmeli (WDAC aktifse)
```

---

## KATMAN 5 — Ağ Güvenliği

### 5.1 — DNS Güvenliği

> [!WARNING]
> Windows 10, işletim sistemi düzeyinde yerleşik **DNS over HTTPS (DoH)** sunmaz (bu Win 11 özelliğidir). Alternatifler:

**Seçenek A — Tarayıcı Seviyesi DoH (Hızlı):**
```
Edge: Ayarlar → Gizlilik, Arama ve Hizmetler → Güvenlik → Güvenli DNS kullan
Chrome: Ayarlar → Gizlilik ve Güvenlik → Güvenlik → Güvenli DNS kullan
```

**Seçenek B — NextDNS veya Cloudflared (Sistem Geneli):**
```powershell
# cloudflared Windows istemcisini kur (https://developers.cloudflare.com/cloudflare-one/connections/connect-devices/warp/)
# Veya NextDNS uygulamasını yükle — sistem geneli DoH sağlar
```

**Seçenek C — Router Seviyesi DoH:**
Router'ın DNS ayarından Cloudflare (1.1.1.1) veya Quad9 (9.9.9.9) kullan. Tüm cihazlara uygulanır.

### 5.2 — SMB Sertleştirme

```powershell
# SMB 1.0 kaldır (WannaCry, EternalBlue saldırı vektörü)
Disable-WindowsOptionalFeature -Online -FeatureName "SMB1Protocol" -NoRestart

# SMB İmzalama Zorunlu Kıl (MitM saldırısına karşı)
Set-SmbServerConfiguration -RequireSecuritySignature $true -Force
Set-SmbClientConfiguration -RequireSecuritySignature $true -Force

# SMB Şifreleme (isteğe bağlı — performans maliyeti var ama güçlü)
Set-SmbServerConfiguration -EncryptData $true -Force

# SMB Null Session kapat
Set-ItemProperty `
    -Path "HKLM:\SYSTEM\CurrentControlSet\Services\LanManServer\Parameters" `
    -Name "RestrictNullSessAccess" -Value 1

# Doğrulama:
Get-SmbServerConfiguration | Select-Object RequireSecuritySignature, EncryptData
Get-WindowsOptionalFeature -Online -FeatureName "SMB1Protocol" | Select-Object State
```

### 5.3 — Windows Güvenlik Duvarı Sertleştirme

```powershell
# Tüm profillerde varsayılan gelen trafiği engelle
Set-NetFirewallProfile -Profile Domain,Public,Private -DefaultInboundAction Block -Enabled True

# Günlüğü etkinleştir
Set-NetFirewallProfile -Profile Domain,Public,Private `
    -LogFileName "%SystemRoot%\System32\LogFiles\Firewall\pfirewall.log" `
    -LogMaxSizeKilobytes 32768 `
    -LogAllowed True `
    -LogBlocked True

# Temel giden kurallara izin ver (beyaz liste yaklaşımı)
New-NetFirewallRule -DisplayName "İzin: DNS UDP"  -Direction Outbound -Protocol UDP -RemotePort 53  -Action Allow
New-NetFirewallRule -DisplayName "İzin: DNS TCP"  -Direction Outbound -Protocol TCP -RemotePort 53  -Action Allow
New-NetFirewallRule -DisplayName "İzin: HTTPS"    -Direction Outbound -Protocol TCP -RemotePort 443 -Action Allow
New-NetFirewallRule -DisplayName "İzin: HTTP"     -Direction Outbound -Protocol TCP -RemotePort 80  -Action Allow
New-NetFirewallRule -DisplayName "İzin: NTP"      -Direction Outbound -Protocol UDP -RemotePort 123 -Action Allow
```

### 5.4 — LLMNR, NetBIOS ve WPAD Kapatma

Bu üçü saldırganların ağ trafiğini yönlendirmek için kullandığı isim çözümleme protokolleridir.

```powershell
# LLMNR kapat (Grup İlkesi):
# Bilgisayar Yapılandırması → Yönetim Şablonları → Ağ → DNS İstemcisi
# → Çok Noktaya Yayın Ad Çözümlemesini Kapat → Etkin

# LLMNR Registry (Home dahil):
New-Item -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows NT\DNSClient" -Force | Out-Null
Set-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows NT\DNSClient" `
    -Name "EnableMulticast" -Value 0

# NetBIOS over TCP/IP devre dışı (tüm adaptörler):
$adapters = Get-WmiObject Win32_NetworkAdapterConfiguration | Where-Object { $_.IPEnabled }
foreach ($adapter in $adapters) {
    $adapter.SetTcpipNetbios(2)   # 2 = Disable
}

# WPAD (Web Proxy Auto-Discovery) kapat:
Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" `
    -Name "AutoDetect" -Value 0
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings\WinHttp" `
    -Name "DisableWpad" -Value 1 -Type DWord
```

### 5.5 — NTP Zaman Sunucusu

Kerberos saldırıları zaman manipülasyonuna dayanabilir.

```powershell
w32tm /config /manualpeerlist:"time.cloudflare.com,0x8 pool.ntp.org,0x8" `
              /syncfromflags:manual /reliable:YES /update
Restart-Service w32tm
w32tm /resync /force
```

---

## KATMAN 6 — Kimlik ve Hesap Güvenliği

### 6.1 — Ayrıcalık Ayrımı: İki Hesap Yapısı

```
┌─────────────────────────────────────────────────────┐
│  Yönetici Hesabı (Admin)                            │
│  • Sadece kurulum / sistem değişikliği              │
│  • Günlük işlerde ASLA kullanılmaz                  │
├─────────────────────────────────────────────────────┤
│  Standart Kullanıcı Hesabı (Günlük)                 │
│  • Web tarama, yazışma, medya, geliştirme           │
│  • Zararlı yazılım bu hesapla çalışırsa kısıtlı     │
└─────────────────────────────────────────────────────┘
```

```powershell
# Yeni standart kullanıcı oluştur
$Sifre = Read-Host -AsSecureString "Yeni kullanıcı şifresi"
New-LocalUser -Name "Kullanici" -Password $Sifre -FullName "Günlük Hesap"
Add-LocalGroupMember -Group "Users" -Member "Kullanici"
# Administrators grubuna EKLEME

# Yerleşik Administrator hesabını devre dışı bırak
Disable-LocalUser -Name "Administrator"

# Kullanıcı listesi:
Get-LocalUser | Select-Object Name, Enabled, LastLogon
```

### 6.2 — UAC Maksimum Seviye

```powershell
$UACPath = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System"
Set-ItemProperty -Path $UACPath -Name "ConsentPromptBehaviorAdmin" -Value 2  # Her zaman sor
Set-ItemProperty -Path $UACPath -Name "ConsentPromptBehaviorUser"  -Value 1  # Kimlik bilgisi iste
Set-ItemProperty -Path $UACPath -Name "PromptOnSecureDesktop"      -Value 1  # Güvenli masaüstünde sor
Set-ItemProperty -Path $UACPath -Name "EnableLUA"                  -Value 1  # UAC aktif
```

### 6.3 — Windows Hello (Biyometrik / PIN)

```
Ayarlar → Hesaplar → Oturum Açma Seçenekleri
  → Windows Hello PIN → Ekle (en az 6 haneli sayısal veya alfanümerik)
  → Parola yerine Windows Hello kullan
```

Fiziksel **FIDO2 donanım anahtarı** (YubiKey, Google Titan Key) PIN'den de güçlüdür.

### 6.4 — Hesap Kilitleme Politikası

```powershell
net accounts /lockoutthreshold:5 /lockoutduration:30 /lockoutwindow:30
net accounts /minpwlen:16 /maxpwage:90 /uniquepw:10

# Doğrulama:
net accounts
```

### 6.5 — Yerel Güvenlik İlkesi Ek Sertleştirme

**secpol.msc** veya Registry:

```powershell
$SecPath = "HKLM:\SYSTEM\CurrentControlSet\Control\Lsa"

# Anonim hesap ve paylaşım erişimini kısıtla
Set-ItemProperty -Path $SecPath -Name "RestrictAnonymous"             -Value 1
Set-ItemProperty -Path $SecPath -Name "RestrictAnonymousSAM"          -Value 1
Set-ItemProperty -Path $SecPath -Name "EveryoneIncludesAnonymous"     -Value 0

# LM ve NTLMv1 kimlik doğrulamayı devre dışı bırak — yalnızca NTLMv2 kabul et
Set-ItemProperty -Path $SecPath -Name "LmCompatibilityLevel"          -Value 5
# 5 = Yalnızca NTLMv2 gönder ve kabul et; LM ve NTLMv1 reddedilir
```

---

## KATMAN 7 — Denetim, İzleme ve Telemetri Kapatma

### 7.1 — Telemetri ve Gizlilik (Windows 10'a Özgü)

> [!IMPORTANT]
> Windows 10, Windows 11'e göre çok daha agresif telemetri gönderir. Aşağıdaki adımlar bu veriyi minimuma indirir.

**Grup İlkesi (Pro/Enterprise):**
```
Bilgisayar Yapılandırması → Yönetim Şablonları → Windows Bileşenleri
  → Veri Toplama ve Önizleme Derlemeleri
     → "Telemetri İzin Ver" → Etkin → Değer: 0 (Security) veya 1 (Basic)
```

**Registry (Home dahil):**
```powershell
# Telemetri seviyesini minimum yap
$TelPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows\DataCollection"
New-Item -Path $TelPath -Force | Out-Null
Set-ItemProperty -Path $TelPath -Name "AllowTelemetry"         -Value 0   # 0 = Security (Enterprise), 1 = Basic (Pro/Home minimum)
Set-ItemProperty -Path $TelPath -Name "DisableOneSettingsDownloads" -Value 1
Set-ItemProperty -Path $TelPath -Name "DoNotShowFeedbackNotifications" -Value 1

# Reklam kimliğini kapat
Set-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows\AdvertisingInfo" `
    -Name "DisabledByGroupPolicy" -Value 1

# Cortana'yı kısıtla
$CortPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows\Windows Search"
New-Item -Path $CortPath -Force | Out-Null
Set-ItemProperty -Path $CortPath -Name "AllowCortana" -Value 0

# Aktivite geçmişi kaydını kapat
Set-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows\System" `
    -Name "PublishUserActivities" -Value 0
Set-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows\System" `
    -Name "EnableActivityFeed" -Value 0

# Bağlantılı kullanıcı deneyimlerini kapat
$DiagPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows\DataCollection"
Set-ItemProperty -Path $DiagPath -Name "AllowDeviceNameInTelemetry" -Value 0

# DiagTrack servisini durdur ve devre dışı bırak
Stop-Service -Name DiagTrack -Force
Set-Service  -Name DiagTrack -StartupType Disabled
```

### 7.2 — Gelişmiş Denetim Politikası

```powershell
auditpol /set /category:"Logon/Logoff"     /success:enable /failure:enable
auditpol /set /category:"Account Logon"    /success:enable /failure:enable
auditpol /set /category:"Account Management" /success:enable /failure:enable
auditpol /set /category:"Object Access"    /success:enable /failure:enable
auditpol /set /category:"Policy Change"    /success:enable /failure:enable
auditpol /set /category:"Privilege Use"    /success:enable /failure:enable
auditpol /set /category:"System"           /success:enable /failure:enable
auditpol /set /category:"Detailed Tracking" /success:enable /failure:enable

# Mevcut durumu görüntüle:
auditpol /get /category:*
```

**İzlenecek Kritik Event ID'ler:**

| Event ID | Anlam | Önem |
|----------|-------|------|
| 4624 | Başarılı oturum açma | Orta |
| 4625 | Başarısız oturum açma | **Yüksek** |
| 4648 | Açık kimlik bilgileriyle oturum açma | **Yüksek** |
| 4672 | Yönetici ayrıcalıklarıyla oturum açma | **Yüksek** |
| 4698 | Zamanlanmış görev oluşturuldu | **Kritik** |
| 4719 | Denetim politikası değişti | **Kritik** |
| 4732 | Yöneticiler grubuna üye eklendi | **Kritik** |
| 4768 | Kerberos TGT talebi | Orta |
| 4776 | Kimlik bilgisi doğrulama girişimi | **Yüksek** |
| 7045 | Yeni servis kuruldu | **Kritik** |
| 4697 | Sisteme yeni servis yüklendi | **Kritik** |
| 1102 | Denetim günlüğü temizlendi | **KRİTİK** |
| 4688 | Yeni süreç oluşturuldu | Orta |

### 7.3 — Sysmon Kurulumu

```powershell
# Araç dizini oluştur
New-Item -ItemType Directory -Force -Path "C:\Tools" | Out-Null

# SwiftOnSecurity sysmon yapılandırmasını indir
Invoke-WebRequest `
    -Uri "https://raw.githubusercontent.com/SwiftOnSecurity/sysmon-config/master/sysmonconfig-export.xml" `
    -OutFile "C:\Tools\sysmon-config.xml"

# Sysmon64.exe'yi Sysinternals'tan indir ve kur
# https://learn.microsoft.com/en-us/sysinternals/downloads/sysmon
.\Sysmon64.exe -accepteula -i C:\Tools\sysmon-config.xml

# Servis durumu:
Get-Service Sysmon64

# Günlükler:
# Event Viewer → Applications and Services Logs → Microsoft → Windows → Sysmon → Operational
```

### 7.4 — PowerShell Günlükleme

```powershell
# Script Block Logging
$PSPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows\PowerShell"

New-Item "$PSPath\ScriptBlockLogging" -Force | Out-Null
Set-ItemProperty "$PSPath\ScriptBlockLogging" -Name "EnableScriptBlockLogging" -Value 1

# Module Logging
New-Item "$PSPath\ModuleLogging" -Force | Out-Null
Set-ItemProperty "$PSPath\ModuleLogging" -Name "EnableModuleLogging" -Value 1

# Transcription (her oturumun tam metin kaydı)
$TranscriptPath = "C:\PSLogs"
New-Item -ItemType Directory -Force -Path $TranscriptPath | Out-Null
New-Item "$PSPath\Transcription" -Force | Out-Null
Set-ItemProperty "$PSPath\Transcription" -Name "EnableTranscripting" -Value 1
Set-ItemProperty "$PSPath\Transcription" -Name "OutputDirectory"     -Value $TranscriptPath
Set-ItemProperty "$PSPath\Transcription" -Name "EnableInvocationHeader" -Value 1
```

### 7.5 — Otomatik Güncellemeleri Yapılandır

```powershell
# Windows Update'i "Güvenlik güncellemeleri otomatik, diğerleri manuel" yap
$WUPath = "HKLM:\SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate\AU"
New-Item -Path $WUPath -Force | Out-Null
Set-ItemProperty -Path $WUPath -Name "AUOptions"              -Value 3   # Otomatik indir, kurulumu sor
Set-ItemProperty -Path $WUPath -Name "AutoInstallMinorUpdates" -Value 1
Set-ItemProperty -Path $WUPath -Name "NoAutoRebootWithLoggedOnUsers" -Value 1
```

---

## BÜTÜNLÜK DOĞRULAMA — Kalenin Sağlamlığını Test Et

```powershell
Write-Host "`n========== WINDOWS 10 KALE DURUM RAPORU ==========" -ForegroundColor Cyan

# 1. BitLocker
try {
    $BL = Get-BitLockerVolume -MountPoint "C:" -ErrorAction Stop
    $BLColor = if ($BL.VolumeStatus -eq "FullyEncrypted" -and $BL.EncryptionMethod -match "256") {"Green"} else {"Red"}
    Write-Host "[BitLocker   ] $($BL.VolumeStatus) / $($BL.EncryptionMethod)" -ForegroundColor $BLColor
} catch { Write-Host "[BitLocker   ] Sorgu başarısız" -ForegroundColor Red }

# 2. TPM
$TPM = Get-Tpm
$TPMColor = if ($TPM.TpmPresent -and $TPM.TpmEnabled) {"Green"} else {"Red"}
Write-Host "[TPM         ] Mevcut: $($TPM.TpmPresent) / Etkin: $($TPM.TpmEnabled)" -ForegroundColor $TPMColor

# 3. VBS
$DG = Get-CimInstance -ClassName Win32_DeviceGuard -Namespace root\Microsoft\Windows\DeviceGuard
$VBSMap = @{0="Devre Dışı"; 1="Etkin Değil (donanım sorunu?)"; 2="Etkin"; 3="Etkin+UEFI Kilitli"}
$VBSColor = if ($DG.VirtualizationBasedSecurityStatus -ge 2) {"Green"} else {"Red"}
Write-Host "[VBS         ] $($VBSMap[$DG.VirtualizationBasedSecurityStatus])" -ForegroundColor $VBSColor

# 4. HVCI
$HVCIOn = $DG.SecurityServicesRunning -contains 2
Write-Host "[HVCI        ] $(if ($HVCIOn) {'AÇIK'} else {'KAPALI'})" -ForegroundColor $(if ($HVCIOn) {"Green"} else {"Red"})

# 5. Credential Guard
$CGOn = $DG.SecurityServicesRunning -contains 1
Write-Host "[Cred. Guard ] $(if ($CGOn) {'AÇIK'} else {'KAPALI'})" -ForegroundColor $(if ($CGOn) {"Green"} else {"Yellow"})

# 6. LSA PPL
$LSA = (Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\Lsa" -ErrorAction SilentlyContinue).RunAsPPL
Write-Host "[LSA PPL     ] $LSA (1 veya 2 olmalı)" -ForegroundColor $(if ($LSA -ge 1) {"Green"} else {"Red"})

# 7. WDigest
$WD = (Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\SecurityProviders\WDigest" -ErrorAction SilentlyContinue).UseLogonCredential
Write-Host "[WDigest     ] $WD (0 olmalı)" -ForegroundColor $(if ($WD -eq 0) {"Green"} else {"Red"})

# 8. Hibernate
$HibOn = (Get-Item "C:\hiberfil.sys" -ErrorAction SilentlyContinue)
Write-Host "[Hibernate   ] $(if ($HibOn) {'AÇIK — RİSK!'} else {'KAPALI — GÜVENLI'})" -ForegroundColor $(if ($HibOn) {"Red"} else {"Green"})

# 9. Pagefile Temizleme
$PF = (Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management" -ErrorAction SilentlyContinue).ClearPageFileAtShutdown
Write-Host "[Pagefile    ] Temizleme: $PF (1 olmalı)" -ForegroundColor $(if ($PF -eq 1) {"Green"} else {"Red"})

# 10. SMB1
$SMB1 = (Get-WindowsOptionalFeature -Online -FeatureName "SMB1Protocol" -ErrorAction SilentlyContinue).State
Write-Host "[SMB1        ] $SMB1 (Disabled olmalı)" -ForegroundColor $(if ($SMB1 -eq "Disabled") {"Green"} else {"Red"})

# 11. Kernel DMA
$DMAOn = $DG.KernelDmaProtectionPolicy
Write-Host "[DMA Koruma  ] $DMAOn (1 olmalı)" -ForegroundColor $(if ($DMAOn -eq 1) {"Green"} else {"Yellow"})

# 12. DiagTrack (Telemetri servisi)
$Diag = Get-Service DiagTrack -ErrorAction SilentlyContinue
Write-Host "[Telemetri   ] Durum: $($Diag.Status) / Başlangıç: $($Diag.StartType) (Disabled olmalı)" `
    -ForegroundColor $(if ($Diag.StartType -eq "Disabled") {"Green"} else {"Yellow"})

Write-Host "===================================================" -ForegroundColor Cyan
```

---

## ÖNEMLİ HATIRLATMALAR

> [!CAUTION]
> **Kurtarma Anahtarı:** BitLocker kurtarma anahtarını Microsoft'a veya buluta gönderme. Çıktı al, fiziksel kasada sakla. Anahtar olmadan kendi disketin de erişilmez olur.

> [!WARNING]
> **Home Sürümü Kısıtlamaları:** gpedit.msc, Credential Guard ve bazı WDAC özellikleri Home'da çalışmaz. Kritik güvenlik için Pro sürümüne geçiş değerlendirilebilir.

> [!WARNING]
> **HVCI ve Eski Sürücüler:** HVCI açıldığında imzasız veya eski sürücüler (bazı oyun anti-cheat sistemi, donanım vendor araçları) çalışmayı durdurabilir. Açmadan önce sürücülerin güncel olduğundan emin ol.

> [!TIP]
> **Güncelleme Disiplini:** Windows 10 desteği **Ekim 2025** sonunda sona erecek. Mümkünse Windows 11'e geçiş planlanmalı. Güncelleme gelmesi durduğunda hiçbir sertleştirme yeterli olmaz.

> [!NOTE]
> **Tehdit Modeli:** Bu rehber fiziksel erişim, Cold Boot, DMA saldırısı, Pass-the-Hash, credential dump, makro saldırısı, ağ dinleme ve telemetri sızıntısına karşı koruma sağlar. Hiçbir sistem yüzde yüz güvenli değildir; amaç saldırı maliyetini maksimize etmektir.

---

## ÖZET — Hangi Katman Neyi Koruyor?

| Tehdit | Koruma Katmanı | Tedbir |
|--------|----------------|--------|
| Disk Çalındı | Katman 1 | BitLocker XTS-AES 256 + Pre-Boot PIN |
| Cold Boot | Katman 3 | Hibernate kapalı · Pagefile temizleme |
| Thunderbolt DMA | Katman 3 | Kernel DMA Koruması · VT-d/IOMMU |
| Bootkit / Rootkit | Katman 0 | Secure Boot · BIOS Admin Şifresi |
| Kernel Exploit | Katman 2 | HVCI · VBS · WDAC |
| Mimikatz / Credential Dump | Katman 2 | Credential Guard · LSA PPL · WDigest=0 |
| Pass-the-Hash | Katman 2 | Credential Guard |
| Makro / Office Saldırısı | Katman 4 | ASR Kuralları |
| DNS Zehirleme | Katman 5 | DoH (tarayıcı/router üzerinden) |
| SMB Yayılma (WannaCry) | Katman 5 | SMB1 kaldırma · SMB İmzalama |
| LLMNR / NBT-NS Spoofing | Katman 5 | LLMNR/NetBIOS kapatma |
| Ayrıcalık Yükseltme | Katman 6 | UAC maks. · Standart Hesap · WDAC |
| Telemetri Sızıntısı | Katman 7 | DiagTrack durdurma · Policy kısıtlama |
| Saldırı Görünmezliği | Katman 7 | Sysmon · Gelişmiş Denetim · PS Günlük |

---

## Windows 10 vs Windows 11 — Güvenlik Farkları

| Özellik | Windows 10 | Windows 11 |
|---------|-----------|-----------|
| DoH (yerleşik OS) | ❌ Yok — tarayıcı/3. taraf | ✅ Yerleşik |
| Credential Guard (Pro) | ⚠️ Manuel etkinleştirme | ✅ Varsayılan etkin |
| Smart App Control | ❌ Yok | ✅ Var |
| TPM 2.0 Zorunluluğu | ❌ Yok | ✅ Kurulum gereksinimi |
| HVCI varsayılan | ❌ Kapalı | ✅ Yeni cihazlarda açık |
| Telemetri agresifliği | ⚠️ Daha yüksek | ✅ Daha düşük |
| Destek sonu | ⚠️ Ekim 2025 | ✅ 2031+ |

---

*Son güncelleme: 2026-03-22 · Windows 10 Pro/Enterprise (22H2) için hazırlanmıştır.*
