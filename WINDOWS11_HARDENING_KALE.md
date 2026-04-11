# 🏰 WINDOWS 11 — "KALE" MİMARİSİ: TAM SPEKTRUM SERTLEŞTİRME REHBERİ

> **Hedef:** Donanım silikon katmanından uygulama belleğine kadar her katmanı kilitleyen, fiziksel çalınma dahil tüm senaryolara karşı dayanıklı bir Windows 11 sistemi inşa etmek.

---

## 📐 MİMARİ HARITA — Saldırı Yüzeyi ve Savunma Katmanları

```
┌─────────────────────────────────────────────────────────┐
│  KATMAN 0 — Fiziksel / Donanım                          │
│  BIOS Admin Şifresi · Secure Boot · TPM 2.0 · DMA      │
├─────────────────────────────────────────────────────────┤
│  KATMAN 1 — Önyükleme (Boot)                            │
│  BitLocker Pre-Boot PIN · UEFI Boot Order · Secure Boot │
├─────────────────────────────────────────────────────────┤
│  KATMAN 2 — İşletim Sistemi Çekirdeği                   │
│  VBS · HVCI · Credential Guard · WDAC · LSA Protection  │
├─────────────────────────────────────────────────────────┤
│  KATMAN 3 — Bellek (RAM)                                │
│  HVCI · Pagefile Temizleme · Hibernate Devre Dışı       │
├─────────────────────────────────────────────────────────┤
│  KATMAN 4 — Süreç / Uygulama                            │
│  Exploit Guard · ASLR · DEP · CFG · ACG · CIG          │
├─────────────────────────────────────────────────────────┤
│  KATMAN 5 — Ağ                                          │
│  DNS-over-HTTPS · SMB Hardening · Firewall · NTP Auth   │
├─────────────────────────────────────────────────────────┤
│  KATMAN 6 — Kimlik / Hesap                              │
│  Windows Hello · FIDO2 · MFA · Standart Hesap · LSA     │
├─────────────────────────────────────────────────────────┤
│  KATMAN 7 — Denetim / İzleme                            │
│  Gelişmiş Denetim · Sysmon · WDAC Günlükleri           │
└─────────────────────────────────────────────────────────┘
```

---

## KATMAN 0 — Donanım ve BIOS/UEFI Kilitleme

> [!IMPORTANT]
> Bu adım, işletim sistemi kurulmadan **önce** yapılmalıdır. BIOS'u kilitlemeyen bir sistem, işletim sistemi seviyesindeki tüm önlemleri atlatılabilir hale getirir.

### 0.1 — BIOS/UEFI Şifreleri

Bilgisayar açılır açılmaz **F2 / Del / F10** (üreticiye göre değişir) ile BIOS'a gir.

| Ayar | Değer | Neden? |
|------|-------|--------|
| **Supervisor / Admin Password** | Güçlü bir şifre belirle | BIOS ayarlarına yetkisiz erişimi engeller |
| **User / Boot Password** | Güçlü bir şifre belirle | POST sonrası disk önyüklemesini kilitler |
| **HDD Password** (varsa) | Güçlü bir şifre belirle | Disk, başka bir bilgisayara takılsa bile açılmaz |

> [!CAUTION]
> BIOS şifreni unutursan pil sıfırlaması veya servis gerekebilir. Şifreyi fiziksel olarak güvenli bir yerde sakla (dijital değil).

### 0.2 — Secure Boot

```
BIOS → Security / Boot → Secure Boot → Enabled
```

- Modunu **"Standard"** veya **"Windows UEFI"** olarak ayarla (Custom değil).
- **Platform Key (PK)**, **Key Exchange Key (KEK)** ve **Signature Database (db)** alanlarını sıfırlama; Microsoft'un varsayılan anahtarları bu kötü yazılımlara karşı korumalı.

**Ne sağlar?** Yalnızca Microsoft'un kriptografik imzasını taşıyan önyükleyicilerin (bootloader) çalışmasına izin verir. Bootkitler, rootkitler ve imzasız işletim sistemleri boot edemez.

### 0.3 — TPM 2.0

```
BIOS → Security → TPM → Enabled (veya Firmware TPM: Intel PTT / AMD fTPM)
```

- TPM versiyonunun **2.0** olduğunu teyit et (1.2 yetersizdir).
- **PCR (Platform Configuration Register) değerleri** doğrulanmış sistem durumuna kilitlenir; bu sayede disk başka bir makinede açılamaz.

### 0.4 — Sanallaştırma (Virtualization)

```
BIOS → CPU / Advanced → Intel VT-x (veya AMD-V) → Enabled
BIOS → CPU / Advanced → VT-d / AMD IOMMU → Enabled
```

> [!NOTE]
> VT-d / AMD IOMMU donanım düzeyinde DMA koruması sağlar. Yalnızca VT-x açmak yeterli değildir.

### 0.5 — I/O Port Yönetimi

| Port | Önerim |
|------|--------|
| Thunderbolt/USB4 | Güvenilir cihazlara "Authorized Only" koy (BIOS'ta varsa) |
| Seri Port (COM) | Disable |
| Paralel Port (LPT) | Disable |
| Wake on LAN / PXE Boot | Disable (kullanmıyorsan) |
| Boot from USB/CD | Disable veya şifreli bırak |

---

## KATMAN 1 — Önyükleme (Boot) Kilitleme

### 1.1 — BitLocker Pre-Boot PIN (En Kritik Adım)

Varsayılan BitLocker yalnızca TPM kullanır. Bu durumda bilgisayar **açılır, hiç şifre sormaz** ve disk şifreli olsa bile RAM üzerinde anahtar açık durur. **TPM + PIN zorunludur.**

#### Grup İlkesi ile PIN'i Zorunlu Kılma

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
| *Uyumlu TPM olmadan BitLocker'a izin ver* | İşaretli olabilir (esneklik), ama gerçek koruma TPM+PIN'den gelir |
| *Başlangıçta TPM ile...* | **"TPM ile başlangıç PIN'ini zorunlu kıl"** seç |
| **BitLocker şifrelemesi için sürücü şifreleme türü** | **Etkin → Tam şifreleme** |
| **İşletim Sistemi için XTS-AES şifre gücünü yapılandır** | **Etkin → XTS-AES 256-bit** |

#### BitLocker'ı Komut Satırından Aktifleştirme (Tercih Edilen)

```powershell
# Yönetici PowerShell
# Uyarı: Önce yukarıdaki Grup İlkesi uygulanmış olmalı

# Sistemin TPM durumunu doğrula
Get-Tpm

# BitLocker'ı TPM+PIN ile başlat (PIN'i güvenli gir)
$pin = Read-Host -AsSecureString "BitLocker PIN girin"
Enable-BitLocker -MountPoint "C:" `
    -EncryptionMethod XtsAes256 `
    -TPMandPinProtector `
    -Pin $pin `
    -UsedSpaceOnly:$false   # Tüm disk (boş alanı da şifrele)

# Kurtarma anahtarını al ve KASAYA YEDEKLE
(Get-BitLockerVolume -MountPoint "C:").KeyProtector |
    Where-Object { $_.KeyProtectorType -eq "RecoveryPassword" }
```

> [!CAUTION]
> Kurtarma anahtarını Microsoft hesabına veya buluta yükleme. Çıktısını al, laminele, fiziksel kasada sakla. Ya da yerel bir USB'ye yaz ve o USB'yi ayrı bir yerde kilitle.

#### BitLocker Durumunu Doğrula

```powershell
Get-BitLockerVolume -MountPoint "C:" | Select-Object *
# EncryptionPercentage: 100 olmalı
# VolumeStatus: FullyEncrypted olmalı
# KeyProtector: {Tpm, TpmPin, RecoveryPassword} görünmeli
```

---

## KATMAN 2 — Çekirdek Yalıtımı (Kernel Hardening)

### 2.1 — VBS (Virtualization-Based Security)

VBS, güvenilir bir "güvenli dünya" (Secure World) oluşturmak için CPU sanallaştırmasını kullanır. Kimlik bilgileri ve çekirdek bütünlük kodu bu yalıtılmış bölgede çalışır.

**Kontrol:**
```powershell
# VBS durumu
(Get-CimInstance -ClassName Win32_DeviceGuard -Namespace root\Microsoft\Windows\DeviceGuard).VirtualizationBasedSecurityStatus
# 2 = Çalışıyor olmalı
```

**Grup İlkesi ile Zorunlu Kılma:**
```
Bilgisayar Yapılandırması
  └─ Yönetim Şablonları
       └─ Sistem
            └─ Device Guard
                 └─ Sanallaştırma Tabanlı Güvenliği Aç → Etkin
                    Platform Güvenlik Düzeyi: Secure Boot ve DMA Koruması
                    Virtualization Based Protection of Code Integrity: UEFI Kilidiyle Etkin
                    Credential Guard Yapılandırması: UEFI Kilidiyle Etkin
```

### 2.2 — HVCI (Hypervisor-Protected Code Integrity / Bellek Bütünlüğü)

```
Ayarlar → Gizlilik ve Güvenlik → Windows Güvenliği
  → Cihaz Güvenliği → Çekirdek Yalıtımı Ayrıntıları
     → Bellek Bütünlüğü: AÇIK
```

**Ne sağlar?** Bir sürücü veya rootkit, çekirdek belleğine değiştirilmiş kod enjekte etmeye çalışırsa **hypervisor katmanı bu girişimi bloke eder.** Ring-0 seviyesi saldırıları etkisiz hale gelir.

**Doğrulama:**
```powershell
$DevGuard = Get-CimInstance -ClassName Win32_DeviceGuard `
    -Namespace root\Microsoft\Windows\DeviceGuard
$DevGuard.SecurityServicesRunning
# 2 içermeli = HVCI aktif
```

### 2.3 — Credential Guard

Windows'un NTLM hash ve Kerberos biletlerini LSA (Local Security Authority) sürecinin içinde tutması yerine, VBS içinde yalıtılmış **LSAIso** sürecine taşımasını sağlar. Pass-the-Hash / Pass-the-Ticket saldırıları işe yaramaz.

```powershell
# Credential Guard durumu
(Get-CimInstance -ClassName Win32_DeviceGuard -Namespace root\Microsoft\Windows\DeviceGuard).SecurityServicesRunning
# 1 içermeli = Credential Guard aktif
```

**Grup İlkesi yolu:** Yukarıdaki Device Guard politikasından **"UEFI Kilidiyle Etkin"** seçilmişse otomatik gelir.

### 2.4 — LSA Koruması (RunAsPPL)

LSA sürecini **Protected Process Light (PPL)** olarak işaretle. Mimikatz ve benzeri araçlar `lsass.exe`'ye enjeksiyon yapamaz.

```powershell
# Regedit veya PowerShell:
Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Lsa" `
    -Name "RunAsPPL" -Value 2 -Type DWord
# Değer 1 = PPL, Değer 2 = PPL (UEFI Kilidi ile — daha güçlü)

# Doğrulama (yeniden başlatma sonrası):
Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Lsa" -Name "RunAsPPL"
```

> [!NOTE]
> Değer 2 ile ayarlandığında bu ayarı geri almak için UEFI'ye girilmesi gerekir. Bilerek yapıyorsun, devam et.

### 2.5 — WDAC (Windows Defender Application Control)

Yalnızca izin verilen uygulamaların çalışmasına izin veren beyaz liste (allowlist) politikasıdır. İmzasız veya tanınmayan çalıştırılabilir dosyalar başlatılamaz.

```powershell
# Temel politika oluştur (denetim modunda başla, önce test et)
$PolicyPath = "$env:USERPROFILE\Desktop\WDACPolicy.xml"
New-CIPolicy -FilePath $PolicyPath `
    -Level Publisher `
    -Fallback Hash `
    -UserPEs `
    -MultiplePolicyFormat

# Politikayı ikili formata dönüştür
$BinPath = "$env:USERPROFILE\Desktop\WDACPolicy.bin"
ConvertFrom-CIPolicy -XmlFilePath $PolicyPath -BinaryFilePath $BinPath

# Sisteme uygula
Copy-Item $BinPath "$env:windir\System32\CodeIntegrity\SIPolicy.p7b"
```

> [!WARNING]
> WDAC'ı hemen zorunlu moda (enforce) almak sistemi kilitleyebilir. Önce **denetim modunda** günlükleri izle, ardından enforce et.

---

## KATMAN 3 — RAM Güvenliği (Cold Boot ve Bellek Sızıntısı Koruması)

> [!IMPORTANT]
> RAM, bilgisayar kapatıldıktan sonra saniyeler (bazen dakikalarca — soğutulursa saatler) boyunca veriyi fiziksel olarak tutar. Bu "Cold Boot Attack" olarak adlandırılır. Aşağıdaki önlemler bu tehdit modeline karşıdır.

### 3.1 — Sayfa Dosyasını (Pagefile) Kapanışta Temizle

Windows, RAM dolduğunda bazı içerikleri `pagefile.sys` dosyasına yazar. Kapatılırken bu dosya silinmeden bırakılırsa içindeki şifreler, anahtarlar ve oturum verileri diskte kalır.

```powershell
# Yönetici olarak çalıştır
Set-ItemProperty `
    -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management" `
    -Name "ClearPageFileAtShutdown" `
    -Value 1 -Type DWord

# Doğrulama:
Get-ItemProperty `
    -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management" `
    -Name "ClearPageFileAtShutdown"
```

> [!NOTE]
> Bu ayar sistemi kapanışta birkaç saniye yavaşlatır. Güvenlik için makul bir bedel.

### 3.2 — Hibernate (Hazırda Bekletme) Dosyasını Devre Dışı Bırak

`hiberfil.sys` dosyası RAM'in tamamını diske yazar. BitLocker açık olsa bile içinde oturum anahtarları bulunabilir.

```powershell
# Yönetici PowerShell
powercfg /hibernate off

# Doğrulama — hiberfil.sys artık olmamalı:
Test-Path "C:\hiberfil.sys"  # False dönmeli
```

### 3.3 — Uyku Modu Yerine Kapatma Politikası

```powershell
# Uyku modunu devre dışı bırak (uyku = RAM kapatılmaz = sızıntı riski)
powercfg /change standby-timeout-ac 0
powercfg /change standby-timeout-dc 0

# Modern Standby (S0ix) yerine S3 sleep tercih edilmeli
# Modern Standby aktifken sistem "kapanmış" gibi görünse de arka planda uyanık
```

### 3.4 — Bellek Bütünlüğü / HVCI (RAM İçi Süreç Koruması)

Daha önce Katman 2'de etkinleştirilen HVCI, aynı zamanda RAM içinde çalışan kernel-mode kodları için **cryptographic integrity** doğrulaması yapar. Bir sürücü RAM'de kendini değiştiremez.

### 3.5 — Çekirdek DMA Koruması

Thunderbolt / USB4 bağlantı noktaları üzerinden DMA (Direct Memory Access) saldırısı yapılabilir — fiziksel erişimi olan biri saniyeler içinde RAM içeriğini çekebilir.

```
Ayarlar → Gizlilik ve Güvenlik → Windows Güvenliği
  → Cihaz Güvenliği → "Çekirdek DMA Koruması: Açık" görünmeli
```

```powershell
# Doğrulama:
(Get-CimInstance -ClassName Win32_DeviceGuard `
    -Namespace root\Microsoft\Windows\DeviceGuard).KernelDmaProtectionPolicy
# 1 = Etkin olmalı
```

**BIOS'ta VT-d / AMD IOMMU etkinleştirilmiş olmalı** (bkz. Katman 0.4).

---

## KATMAN 4 — Süreç ve Uygulama Belleği Koruması

### 4.1 — Exploit Protection Tam Yapılandırması

```
Windows Güvenliği → Uygulama ve Tarayıcı Denetimi
  → Exploit Protection Ayarları → Sistem Ayarları
```

| Koruma | Ayar | Açıklama |
|--------|------|----------|
| **CFG** (Control Flow Guard) | Açık (varsayılan) | Fonksiyon çağrı akışını doğrular |
| **DEP** (Data Execution Prevention) | Açık (varsayılan) | Veri bölgelerinde kod çalıştırmayı engeller |
| **SEHOP** | Açık | Yapılandırılmış istisna yönetimi manipülasyonunu engeller |
| **Heap Randomization** | Açık | Yığın alanını rastgele düzenler |
| **Zorunlu ASLR** | Açık | Tüm modülleri rastgele adrese yükler (ASLR desteklemeyenleri de zorlar) |

**Zorunlu ASLR PowerShell ile:**
```powershell
# Zorunlu ASLR'ı etkinleştir (bazı eski uygulamaları kırabilir, test et)
Set-ProcessMitigation -System -Enable ForceRelocateImages
```

**Exploit Protection politikasını dışa aktar ve yedekle:**
```powershell
Get-ProcessMitigation -RegistryConfigFilePath "C:\Güvenlik\ExploitProtection.xml"
# Bu dosyayı başka makinelere de uygulayabilirsin
```

### 4.2 — Uygulama Bazlı Ek Korumalar

Tarayıcı gibi kritik uygulamalara ek hafifletmeler uygula:

```powershell
# Örnek: Chrome için ek korumalar
Set-ProcessMitigation -Name chrome.exe `
    -Enable DEP, SEHOP, ForceRelocateImages, BottomUp, `
            StrictHandle, DisableWin32kSystemCalls, `
            AuditSystemCall, BlockNonMicrosoftFonts, `
            DisableExtensionPoints, ProhibitDynamicCode
```

> [!WARNING]
> `ProhibitDynamicCode` (ACG) ve `DisableWin32kSystemCalls` bazı uygulamaları bozabilir. Her uygulamada test et.

### 4.3 — Attack Surface Reduction (ASR) Kuralları

ASR, Defender'ın içinde bulunan ama çoğu kişinin açmadığı güçlü bir katmandır.

```powershell
# Tüm temel ASR kuralları — Yönetici PowerShell

$ASRRules = @(
    "BE9BA2D9-53EA-4CDC-84E5-9B1EEEE46550",  # Office makrolarının Win32 API çağırmasını engelle
    "D4F940AB-401B-4EFC-AADC-AD5F3C50688A",  # Office'in alt süreç oluşturmasını engelle
    "3B576869-A4EC-4529-8536-B80A7769E899",  # Office'in yürütülebilir içerik oluşturmasını engelle
    "75668C1F-73B5-4CF0-BB93-3ECF5CB7CC84",  # Office uygulamalarının süreçlere enjeksiyonunu engelle
    "D3E037E1-3EB8-44C8-A917-57927947596D",  # JS/VBS'nin indirilen içerik çalıştırmasını engelle
    "5BEB7EFE-FD9A-4556-801D-275E5FFC04CC",  # Karartılmış betiklerin yürütülmesini engelle
    "92E97FA1-2EDF-4476-BDD6-9DD0B4DDDC7B",  # Win32 API'sini Office makrolarından çağırmayı engelle
    "01443614-CD74-433A-B99E-2ECDC07BFC25",  # Yaygın olmayan veya fidye yazılımı belirtisi taşıyan süreçleri engelle
    "C1DB55AB-C21A-4637-BB3F-A12568109D35",  # Fidye yazılımı davranışını engelle
    "9E6C4E1F-7D60-472F-BA1A-A39EF669E4B2",  # LSA kimlik bilgilerini çalmayı engelle
    "D1E49AAC-8F56-4280-B9BA-993A6D77406C",  # PsExec ve WMI'dan süreç oluşturulmasını engelle
    "B2B3F03D-6A65-4F7B-A9C7-1C7EF74A9BA4",  # E-posta ve web postasından yürütülebilir içeriği engelle
    "26190899-1602-49E8-8B27-EB1D0A1CE869",  # Office uygulamalarının iletişim uygulamaları oluşturmasını engelle
    "7674BA52-37EB-4A4F-A9A1-F0F9A1619A2C",  # Adobe Reader alt süreç oluşturmasını engelle
    "E6DB77E5-3DF2-4CF1-B95A-636979351E5B"   # USB'den çalışan güvenilmez ve imzasız süreçleri engelle
)

foreach ($Rule in $ASRRules) {
    Add-MpPreference -AttackSurfaceReductionRules_Ids $Rule `
                     -AttackSurfaceReductionRules_Actions Enabled
}

# Doğrulama:
Get-MpPreference | Select-Object AttackSurfaceReductionRules_Ids, AttackSurfaceReductionRules_Actions
```

---

## KATMAN 5 — Ağ Güvenliği

### 5.1 — DNS over HTTPS (DoH)

```
Ayarlar → Ağ ve İnternet → Ethernet / Wi-Fi → [Bağlantı Adı]
  → DNS sunucu ataması → Düzenle
  → "Manuel" → "Şifreli yalnızca (HTTPS üzerinden DNS)"
```

**Önerilen DNS sunucular:**

| Sağlayıcı | IPv4 | IPv6 | DoH URL |
|-----------|------|------|---------|
| Cloudflare | 1.1.1.1 | 2606:4700:4700::1111 | https://cloudflare-dns.com/dns-query |
| Quad9 (Kötücül engeli dahil) | 9.9.9.9 | 2620:fe::fe | https://dns.quad9.net/dns-query |

### 5.2 — SMB Sertleştirme

```powershell
# SMB 1.0 kaldır (WannaCry, EternalBlue saldırı vektörü)
Disable-WindowsOptionalFeature -Online -FeatureName "SMB1Protocol" -NoRestart

# SMB İmzalama Zorunlu Kıl (MitM saldırısına karşı)
Set-SmbServerConfiguration -RequireSecuritySignature $true -Force
Set-SmbClientConfiguration -RequireSecuritySignature $true -Force

# SMB şifrelemeyi etkinleştir (yerel ağda bile veri şifreli gider)
Set-SmbServerConfiguration -EncryptData $true -Force

# SMB Null Session'ları devre dışı bırak
Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\LanManServer\Parameters" `
    -Name "RestrictNullSessAccess" -Value 1

# Doğrulama:
Get-SmbServerConfiguration | Select-Object RequireSecuritySignature, EncryptData
```

### 5.3 — Windows Güvenlik Duvarı (Whitelist Yaklaşımı)

```powershell
# Tüm gelen bağlantıları varsayılan engelle
Set-NetFirewallProfile -Profile Domain,Public,Private -DefaultInboundAction Block

# Tüm giden bağlantıları varsayılan engelle (opsiyonel — ağır ama güçlü)
# Set-NetFirewallProfile -Profile Domain,Public,Private -DefaultOutboundAction Block

# Yalnızca ihtiyaç duyulan giden bağlantılara izin ver
New-NetFirewallRule -DisplayName "İzin: DNS" -Direction Outbound `
    -Protocol UDP -RemotePort 53 -Action Allow
New-NetFirewallRule -DisplayName "İzin: HTTPS" -Direction Outbound `
    -Protocol TCP -RemotePort 443 -Action Allow

# NetBIOS ve LLMNR binding kapat
Set-NetAdapterBinding -Name "*" -ComponentID ms_msclient -Enabled $false

# LLMNR Grup İlkesi:
# Bilgisayar Yapılandırması → Yönetim Şablonları → Ağ → DNS İstemcisi
# → Çok Noktaya Yayın Ad Çözümlemesini Kapat → Etkin
```

### 5.4 — LLDP, NetBIOS ve WPAD Kapatma

```powershell
# WPAD (Web Proxy Auto-Discovery) - kimlik bilgisi hırsızlığı vektörü
Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" `
    -Name "AutoDetect" -Value 0

# NetBIOS over TCP/IP devre dışı (tüm adaptörler)
$adapters = Get-WmiObject Win32_NetworkAdapterConfiguration | Where-Object { $_.IPEnabled }
foreach ($adapter in $adapters) {
    $adapter.SetTcpipNetbios(2)  # 2 = Disable NetBIOS over TCP/IP
}
```

### 5.5 — NTP Zaman Sunucusu Güvenliği

Zaman manipülasyonu Kerberos saldırılarına yol açabilir.

```powershell
# Güvenilir NTP sunucuya kilitle
w32tm /config /manualpeerlist:"time.cloudflare.com" /syncfromflags:manual /reliable:YES /update
Restart-Service w32tm
w32tm /resync
```

---

## KATMAN 6 — Kimlik ve Hesap Güvenliği

### 6.1 — Hesap Yapısı: Ayrıcalık Ayrımı

```
┌─────────────────────────────────────────────────────┐
│  Yönetici Hesabı (Admin)                            │
│  • Sadece yazılım kurulumu / sistem değişikliği     │
│  • Günlük işlerde ASLA kullanılmaz                  │
├─────────────────────────────────────────────────────┤
│  Standart Kullanıcı Hesabı (Günlük)                 │
│  • Web tarama, ofis uygulamaları, medya             │
│  • Tüm zararlı kodlar bu hesapla çalışırsa sınırlı  │
└─────────────────────────────────────────────────────┘
```

```powershell
# Yeni Standart kullanıcı oluştur
$Sifre = Read-Host -AsSecureString "Yeni kullanıcı şifresi"
New-LocalUser -Name "Kullanici" -Password $Sifre -FullName "Günlük Hesap"
Add-LocalGroupMember -Group "Users" -Member "Kullanici"
# "Administrators" grubuna EKLEME

# Yerleşik Administrator hesabını devre dışı bırak
Disable-LocalUser -Name "Administrator"
```

### 6.2 — UAC'ı Maksimum Seviyeye Getir

```
Ayarlar → Hesaplar → Kullanıcı Hesabı Denetimi
  → En üst seviye: "Uygulamalar sistem ayarlarını değiştirmeye çalıştığında her zaman uyar"
```

```powershell
# Regedit ile:
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" `
    -Name "ConsentPromptBehaviorAdmin" -Value 2    # Her zaman sor
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" `
    -Name "ConsentPromptBehaviorUser" -Value 1     # Kimlik bilgisi iste
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" `
    -Name "PromptOnSecureDesktop" -Value 1         # Güvenli masaüstünde sor
```

### 6.3 — Windows Hello / FIDO2 Anahtar ile Giriş

```powershell
# PIN'siz hesap politikasını engelle — PIN zorunlu olsun
Set-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows\System" `
    -Name "AllowDomainPINLogon" -Value 1
```

- **Fiziksel FIDO2 donanım anahtarı** (YubiKey, Google Titan) kullanmak PIN'den de güvenlidir.
- Windows Hello for Business PIN'ini 6 haneli, karmaşık yap.

### 6.4 — Hesap Kilitleme Politikası

```powershell
# 5 başarısız denemeden sonra hesabı kilitle
net accounts /lockoutthreshold:5 /lockoutduration:30 /lockoutwindow:30

# Doğrulama:
net accounts
```

### 6.5 — Parola Politikası

```powershell
net accounts /minpwlen:16 /maxpwage:90 /uniquepw:10
```

---

## KATMAN 7 — Denetim, Kayıt ve İzleme

### 7.1 — Gelişmiş Denetim Politikası

```powershell
# Tüm kritik kategorileri etkinleştir
auditpol /set /category:"Logon/Logoff" /success:enable /failure:enable
auditpol /set /category:"Account Logon" /success:enable /failure:enable
auditpol /set /category:"Account Management" /success:enable /failure:enable
auditpol /set /category:"Object Access" /success:enable /failure:enable
auditpol /set /category:"Policy Change" /success:enable /failure:enable
auditpol /set /category:"Privilege Use" /success:enable /failure:enable
auditpol /set /category:"System" /success:enable /failure:enable
auditpol /set /category:"Detailed Tracking" /success:enable /failure:enable

# Mevcut durumu görüntüle:
auditpol /get /category:*
```

**İzlenecek kritik Event ID'leri:**

| Event ID | Anlam | Önem |
|----------|-------|------|
| 4624 | Başarılı oturum açma | Orta |
| 4625 | Başarısız oturum açma | **Yüksek** |
| 4648 | Açık kimlik bilgileriyle oturum açma | **Yüksek** |
| 4719 | Denetim politikası değişti | **Kritik** |
| 4732 | Yöneticiler grubuna üye eklendi | **Kritik** |
| 4768 | Kerberos TGT talebi | Orta |
| 4776 | Kimlik bilgisi doğrulama girişimi | **Yüksek** |
| 7045 | Yeni servis kuruldu | **Kritik** |
| 4697 | Sistem servisine yeni servis yüklendi | **Kritik** |
| 1102 | Denetim günlüğü temizlendi | **KRİTİK** |

### 7.2 — Sysmon Kurulumu (Derinlemesine Süreç İzleme)

Sysmon, varsayılan Windows günlüklerinden çok daha derin telemetri toplar.

```powershell
# Sysmon indir (https://learn.microsoft.com/sysinternals/sysmon)
# SwiftOnSecurity yapılandırması ile kur (kapsamlı, test edilmiş kural seti)
Invoke-WebRequest `
    -Uri "https://raw.githubusercontent.com/SwiftOnSecurity/sysmon-config/master/sysmonconfig-export.xml" `
    -OutFile "C:\Tools\sysmon-config.xml"

# Kur:
.\Sysmon64.exe -accepteula -i C:\Tools\sysmon-config.xml

# Doğrulama:
Get-Service Sysmon64
# Günlükler: Event Viewer → Applications and Services Logs → Microsoft → Windows → Sysmon
```

### 7.3 — PowerShell Günlükleme

```powershell
# Script Block Logging — her PowerShell komutunu kaydet
Set-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows\PowerShell\ScriptBlockLogging" `
    -Name "EnableScriptBlockLogging" -Value 1

# Module Logging
Set-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows\PowerShell\ModuleLogging" `
    -Name "EnableModuleLogging" -Value 1

# Transcription (her oturumun tam metni)
$TranscriptPath = "C:\PowerShellLogs"
New-Item -ItemType Directory -Force -Path $TranscriptPath | Out-Null
Set-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows\PowerShell\Transcription" `
    -Name "EnableTranscripting" -Value 1
Set-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows\PowerShell\Transcription" `
    -Name "OutputDirectory" -Value $TranscriptPath
```

---

## BÜTÜNLÜK DOĞRULAMA — Kalenin Sağlamlığını Test Et

### Tek Seferlik Durum Raporu

```powershell
Write-Host "========== KALE DURUM RAPORU ==========" -ForegroundColor Cyan

# BitLocker
$BL = Get-BitLockerVolume -MountPoint "C:"
Write-Host "BitLocker: $($BL.VolumeStatus) / Metot: $($BL.EncryptionMethod)" `
    -ForegroundColor $(if ($BL.VolumeStatus -eq "FullyEncrypted") {"Green"} else {"Red"})

# TPM
$TPM = Get-Tpm
Write-Host "TPM Mevcut: $($TPM.TpmPresent) / Etkin: $($TPM.TpmEnabled) / v$($TPM.ManufacturerVersion)" `
    -ForegroundColor $(if ($TPM.TpmEnabled) {"Green"} else {"Red"})

# VBS
$DG = Get-CimInstance -ClassName Win32_DeviceGuard -Namespace root\Microsoft\Windows\DeviceGuard
$VBSStatus = @{0="Devre Dışı"; 1="Etkin Değil"; 2="Etkin"; 3="Etkin (UEFI Kilitli)"}
Write-Host "VBS: $($VBSStatus[$DG.VirtualizationBasedSecurityStatus])" `
    -ForegroundColor $(if ($DG.VirtualizationBasedSecurityStatus -ge 2) {"Green"} else {"Red"})

# HVCI
$HVCIOn = $DG.SecurityServicesRunning -contains 2
Write-Host "HVCI (Bellek Bütünlüğü): $(if ($HVCIOn) {'AÇIK'} else {'KAPALI'})" `
    -ForegroundColor $(if ($HVCIOn) {"Green"} else {"Red"})

# Credential Guard
$CGOn = $DG.SecurityServicesRunning -contains 1
Write-Host "Credential Guard: $(if ($CGOn) {'AÇIK'} else {'KAPALI'})" `
    -ForegroundColor $(if ($CGOn) {"Green"} else {"Red"})

# Hibernate
$HibOn = (Get-Item C:\hiberfil.sys -ErrorAction SilentlyContinue)
Write-Host "Hibernate (Kapalı olmalı): $(if ($HibOn) {'AÇIK — RİSK!'} else {'KAPALI — GÜVENLI'})" `
    -ForegroundColor $(if ($HibOn) {"Red"} else {"Green"})

# LSA PPL
$LSA = (Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\Lsa").RunAsPPL
Write-Host "LSA PPL: $LSA (2 ideal)" `
    -ForegroundColor $(if ($LSA -ge 1) {"Green"} else {"Red"})

# PageFile Temizleme
$PF = (Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management").ClearPageFileAtShutdown
Write-Host "Pagefile Temizleme: $PF (1 olmalı)" `
    -ForegroundColor $(if ($PF -eq 1) {"Green"} else {"Red"})

# SMB1
$SMB1 = (Get-WindowsOptionalFeature -Online -FeatureName "SMB1Protocol").State
Write-Host "SMB1: $SMB1 (Disabled olmalı)" `
    -ForegroundColor $(if ($SMB1 -eq "Disabled") {"Green"} else {"Red"})

Write-Host "=======================================" -ForegroundColor Cyan
```

---

## ÖNEMLİ HATIRLATMALAR

> [!CAUTION]
> **Kurtarma anahtarı yönetimi:** BitLocker kurtarma anahtarını çıktı alıp fiziksel kasada sakla. Microsoft'a veya OneDrive'a gönderme. Bunu yapmazsan disk şifresini kırma ihtimali sıfıra yaklaşır — kurtarma anahtarı olmadan kendin de içine giremezsin.

> [!WARNING]
> **UEFI Kilidi:** LSA PPL ve Credential Guard'ı UEFI kilidi ile etkinleştirirsen bu ayarları geri almak için BIOS'a girman gerekir. Bilerek yapıyorsun, ama test ortamında dene.

> [!TIP]
> **Güncelleme disiplini:** Windows Update'i erteletme. Her "Patch Tuesday" (Salı) kritik güvenlik yamaları gelir. Çekirdek katmanındaki güvenlik açıkları yama yapılmadan bu kılavuzdaki hiçbir şey tam koruma sağlamaz.

> [!NOTE]
> **Tehdit modeli:** Bu kılavuz fiziksel erişim, Cold Boot, disk çıkarma, Thunderbolt DMA, Pass-the-Hash, process injection ve birçok yazılım tabanlı saldırıya karşı koruma sağlar. Hiçbir sistem yüzde yüz güvenli değildir. Saldırı yüzeyini küçültmek ve saldırıyı pahalıya mal etmek hedeflenir.

---

## ÖZET — Hangi Katman Neyi Koruyor?

| Tehdit | Koruma Katmanı | Tedbir |
|--------|----------------|--------|
| Disk Çalındı | Katman 1 | BitLocker XTS-AES 256 + Başlangıç PIN |
| Cold Boot (RAM Dondurma) | Katman 3 | Hibernate kapalı · Pagefile temizleme · HVCI |
| Thunderbolt DMA Saldırısı | Katman 3 | Kernel DMA Koruması · VT-d/IOMMU |
| Bootkit / Rootkit | Katman 0, 1 | Secure Boot · BIOS Admin Şifresi |
| Kernel Exploit | Katman 2 | HVCI · VBS · WDAC |
| Kimlik Bilgisi Hırsızlığı | Katman 2, 6 | Credential Guard · LSA PPL |
| Pass-the-Hash | Katman 2 | Credential Guard |
| Makro / Ofis Saldırısı | Katman 4 | ASR Kuralları |
| DNS Ele Geçirme | Katman 5 | DNS over HTTPS |
| SMB Yayılma | Katman 5 | SMB1 Kaldırma · SMB İmzalama |
| Ayrıcalık Yükseltme | Katman 6 | UAC maks. · Standart Hesap · WDAC |
| Saldırı Görünmezliği | Katman 7 | Sysmon · Gelişmiş Denetim · PS Günlük |

---

*Son güncelleme: 2026-03-22 · Windows 11 Pro için hazırlanmıştır.*
