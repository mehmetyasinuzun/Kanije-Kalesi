# 💾 DUAL BOOT GÜVENLİK ANALİZİ & DEPOLAMA SERTLEŞTİRME REHBERİ
## HDD · SSD · Dual Boot Saldırı Yüzeyleri · Ayrı Disk Mimarisi

> **Kapsam:** Tüm dual boot kombinasyonlarının güvenlik analizi (Win10+Win10, Win11+Win11, Win10+Win11, Windows+Linux) ve HDD/SSD için depolama katmanı güvenlik sertleştirmesi.

---

## 📐 TEMEL MİMARİ: Tek Disk vs Ayrı Diskler

```
❌ RİSKLİ — Tek Disk, Çoklu Bölüm:
┌────────────────────────────────────────────────────────────┐
│  Tek Fiziksel Disk                                         │
│ ┌──────────┬────────────┬────────────┬────────────────┐   │
│ │   ESP    │  OS 1      │  OS 2      │  Ortak/Veri    │   │
│ │(EFI Part)│(Win 10)    │(Win 11)    │                │   │
│ │ ORTAK    │ BitLocker  │ BitLocker  │                │   │
│ │ HEDEF ↗  │ (bağımsız) │ (bağımsız) │                │   │
│ └──────────┴────────────┴────────────┴────────────────┘   │
│  OS-1 açıkken OS-2 bölümüne ERIŞIM MÜMKÜN                 │
│  EFI bölümü her iki OS'tan değiştirilebilir               │
└────────────────────────────────────────────────────────────┘

✅ GÜVENLİ — İki Ayrı Fiziksel Disk:
┌──────────────────────┐    ┌──────────────────────────────┐
│  Disk 1 (NVMe/SSD)   │    │  Disk 2 (SSD/HDD)            │
│ ┌──────┬───────────┐ │    │ ┌──────┬──────────────────┐  │
│ │ ESP1 │  Windows  │ │    │ │ ESP2 │  Windows 10 /    │  │
│ │      │    11     │ │    │ │      │  Linux / Win11   │  │
│ │      │ BitLocker │ │    │ │      │  BitLocker / LUKS│  │
│ └──────┴───────────┘ │    │ └──────┴──────────────────┘  │
│  Kendi TPM PCR'ı     │    │  Kendi TPM PCR'ı             │
│  Tamamen bağımsız    │    │  Tamamen bağımsız             │
└──────────────────────┘    └──────────────────────────────┘
  BIOS'ta boot order kilitli → sadece birini seç
  Diğer disk çalışırken ikincisi FIZIKSEL OLARAK ULAŞILAMAZ
```

> [!IMPORTANT]
> **Ayrı fiziksel diskler** en yüksek güvenlik düzeyini sağlar. Bir OS çalışırken diğer diskin içeriğine donanım katmanında erişim yoktur. Tüm EFI, TPM ve BitLocker yapılandırmaları tamamen bağımsız çalışır.

---

## BÖLÜM 1 — DUAL BOOT SALDIRI YÜZEYLERİ

### 1.1 — Ortak Altyapı: Neden Her Kombinasyon Riskli?

Dual boot'ta hangi kombinasyon olursa olsun şu katmanlar **her zaman ortaktır:**

| Ortak Katman | Risk | Etki |
|-------------|------|------|
| **EFI Sistem Bölümü (ESP)** | Bootloader değişikliği | OS-1'den OS-2'nin EFI girdisi değiştirilebilir; bootkit yerleştirilebilir |
| **TPM PCR Ölçümleri** | Boot zinciri karışması | Her OS geçişi TPM ölçümlerini değiştirir → BitLocker kurtarma anahtarı ister |
| **BCD (Boot Config Data)** | Merkezi önyükleme yönetimi | Tek bir EFI değişkeni değişikliği her iki OS'u etkiler |
| **Fiziksel NAND/Plaka** | Wear leveling / manyetik kalıntı | Bir OS'tan silinen veri, diğer OS üzerinden forensic araçlarla okunabilir |

---

### 1.2 — Senaryo A: Windows 10 + Windows 10 (İki Farklı Kurulum)

```
Sistem: Tek disk
OS-A: Windows 10 (22H2) — C:
OS-B: Windows 10 (21H2) — D:
```

**Saldırı Yüzeyleri:**

| # | Saldırı | Mekanizma | Sonuç |
|---|---------|-----------|-------|
| 1 | **Eski Çekirdek Exploit** | OS-B eski sürüm Windows 10 → yamansız güvenlik açığı | OS-B'den başlayan exploit, ortak EFI'yi hedefler |
| 2 | **Çapraz NTFS Okuma** | OS-A çalışırken OS-B'nin NTFS bölümü `diskpart` ile erişilebilir | OS-B'nin `SAM`, `SECURITY`, `SYSTEM` hive dosyaları çekilebilir → hash dump |
| 3 | **SAM Hash Çapraz Senkronizasyon** | İki farklı SAM veritabanı — aynı şifre kullanılıyorsa | OS-B'nin hash'i çalınır, OS-A'ya Pass-the-Hash uygulanır |
| 4 | **BCD Manipülasyonu** | OS-A'dan `bcdedit` komutu OS-B'nin boot girdisini değiştirir | OS-B yerine kötücül imaj yüklenir (evil maid saldırısı) |
| 5 | **Payload Kalıcılığı** | OS-B'ye zararlı yazılım yerleştirilir; OS-A aktifken arka planda yürütülmez ama disk üzerinde durur | OS-B açıldığında payload çalışır; OS-A'dan temizleme mümkün değil |

**Spesifik Tehdit: SAM Dosyası Çapraz Erişim**
```powershell
# OS-A çalışırken OS-B'nin locked SAM dosyasına erişim (örnek saldırı):
# reg save HKLM\SAM C:\stolen_sam.hiv  ← OS-A'dan çalışmaz (kendi sistemi)
# Ama OS-B'nin SAM FIZIKSEL DOSYASI OS-A altında erişilebilir:
# D:\Windows\System32\config\SAM ← Bu dosya D: sürücüsüne erişimle okunabilir
# Sonra: impacket-secretsdump LOCAL -sam stolen_sam.hiv -system stolen_system.hiv

# SAVUNMA: OS-B bölümünü BitLocker ile şifrele + ayrı kurtarma anahtarı
```

**Azaltma:**
- Her iki bölüm de BitLocker XTS-AES 256 + PIN ile şifreli olmalı
- Ayrı güçlü PIN'ler (aynı PIN'i iki OS için kullanma)
- Kullanılmayan OS'un bölümünü `manage-bde -lock D:` ile kilitle

---

### 1.3 — Senaryo B: Windows 11 + Windows 11 (İki Farklı Kurulum)

```
Sistem: Tek disk
OS-A: Windows 11 Pro (23H2) — C:
OS-B: Windows 11 Home veya test kurulumu — D:
```

**Saldırı Yüzeyleri:**

| # | Saldırı | Mekanizma | Sonuç |
|---|---------|-----------|-------|
| 1 | **VBS/HVCI Çakışması** | İki ayrı VBS yapılandırması ortak hypervisor PCR ölçümleri üzerinden çakışabilir | Credential Guard beklenmedik şekilde devre dışı kalabilir |
| 2 | **TPM PCR-7 (Secure Boot Durumu) Değişimi** | OS-B'nin önyükleyicisi Secure Boot'u farklı ölçer | OS-A'nın BitLocker'ı kurtarma moduna düşer |
| 3 | **Daha Zayıf OS-B'den Pivot** | OS-B Home → Credential Guard yok → NTLM hashleri RAM'de açık | OS-B açıkken dump alınır; aynı kullanıcı OS-A'ya giriş yaparsa hash kullanılır |
| 4 | **Smart App Control Devre Dışı** | Win11'de SAC yalnızca temiz kurulumda devrede; OS-B farklı SAC durumunda | SAC bypass'ı OS-B üzerinden test edilip OS-A'ya taşınabilir |
| 5 | **UEFI Firmware Saldırısı** | İki OS ortak UEFI değişkenlerine yazar | OS-B'den `efi-updatevar` komutu DB/DBX değişkenlerini değiştirir → Secure Boot zayıflar |

**Spesifik Tehdit: UEFI Değişken Manipülasyonu**
```
Win11 OS-B (daha az güvenli) açıkken:
  → efivar veya PowerShell EFI API ile UEFI DB tablosuna kötücül imza eklenir
  → Secure Boot hâlâ "aktif" görünür ama artık imzalı sayan listesi genişledi
  → OS-A açıldığında bootloader imzayı "geçerli" bulur → Bootkit yüklenir

SAVUNMA: BIOS'ta "Secure Boot key management" erişimini Admin şifresiyle kilitle
```

---

### 1.4 — Senaryo C: Windows 10 + Windows 11 (Karma) — En Yaygın ve En Riskli

```
Sistem: Tek disk (veya iki disk)
OS-A: Windows 11 Pro — C:
OS-B: Windows 10 Pro — D:
```

Bu senaryo en yaygın ve güvenlik açısından en sorunlu kombinasyondur.

**Saldırı Yüzeyleri:**

| # | Saldırı | Mekanizma | Sonuç |
|---|---------|-----------|-------|
| 1 | **En Zayıf Halka Prensibi** | Win10 pek çok Win11 güvenlik özelliğinden yoksun | Saldırgan Win10'dan başlar, Win11 dosyalarına ulaşır |
| 2 | **Credential Guard Yokluğu (Win10 Pro)** | Win10'da Credential Guard varsayılan kapalı | NTLM hash dump → Pass-the-Hash → Win11 oturumuna geçiş |
| 3 | **Telemetri Farkı** | Win10 çok daha fazla telemetri gönderir | Saldırgan win10 ağ trafiği üzerinden oturum bilgisine ulaşabilir |
| 4 | **BitLocker Pre-Boot PIN Yorgunluğu** | Her OS geçişinde PIN → kullanıcı PIN'i kaldırır | TPM-only mod → Cold Boot koruması biter |
| 5 | **Patch Seviyesi Farkı** | Win10 her zaman Win11'den daha az yamaya sahip olur | Win10'da açık N-day exploit → privilege escalation → disk erişimi |
| 6 | **ASR Kural Farkı** | Win10'da bazı ASR kuralları farklı davranır | Saldırgan Win10 altında test eder, Win11'de çalışır mı dener |

**Bu Senaryo için En İyi Yapılandırma:**
```
Eğer kesinlikle tek diskte yapılacaksa:

1. Her iki OS da BitLocker XTS-AES 256 + ayrı PIN
2. Win10 bölümünü Win11 açıkken KİLİTLE:
   manage-bde -lock D: -ForceDismount

3. Win11 açıkken Win10 bölümüne otomatik kilitleme:
   Görev Zamanlayıcısı → Kullanıcı oturumu açıldığında:
   manage-bde -lock D: -ForceDismount

4. Win10 bölümünü sadece gerektiğinde aç:
   manage-bde -unlock D: -RecoveryPassword <48-haneli-anahtar>
```

---

### 1.5 — Senaryo D: Windows + Linux (En Karmaşık Saldırı Yüzeyi)

```
Sistem: Tek disk
OS-A: Windows 10 veya 11 — C:
OS-B: Ubuntu / Fedora / Arch — /dev/sda2
Bootloader: GRUB2 (ortak)
```

**Windows+Linux dual boot, en geniş saldırı yüzeyine sahip kombinasyondur.**

**Saldırı Yüzeyleri:**

| # | Saldırı | Mekanizma | Sonuç |
|---|---------|-----------|-------|
| 1 | **GRUB Bootloader Saldırısı** | GRUB, UEFI'ye göre çok daha az kısıtlı bir ortam sunar | `grub rescue` shell'den Windows NTFS bölümüne erişim |
| 2 | **NTFS-3G Erişimi** | Linux'tan `sudo mount /dev/sda1 /mnt` → Windows dosya sistemi tam erişim | SAM, SECURITY, SYSTEM kayıt dosyaları kopyalanabilir |
| 3 | **BitLocker Atlama (Kim Açıktaysa)** | Linux'tan `dislocker` aracıyla şifrelenmemiş NTFS bölümüne erişim | BitLocker aktif değilse Windows bölümü Linux altında tamamen okunabilir |
| 4 | **Mimikatz Linux Versiyonu** | `pypykatz` Python aracı Linux altında çalışır | Linux'tan Windows SAM/LSASS dosyası analizi |
| 5 | **Sudo Zafiyeti (PrivEsc)** | Linux tarafında `sudo` açığı veya SUID bit exploit | Root erişimi → LUKS şifreleme bypass → Windows bölümüne erişim |
| 6 | **Farklı Güvenlik Duvarı Politikası** | Linux'ta `iptables` yapılandırması, Windows Defender Firewall'dan farklı | Saldırgan Linux tarafından ağ keşfi yapar, Windows tarafına pivot eder |
| 7 | **Ortak Ev Dizini/Veri Bölümü** | Ortak NTFS veya exFAT veri bölümü kullanıyorsa | Bir OS'ta oluşturulan zararlı dosya diğer OS'u etkiler |
| 8 | **Linux Kernel Exploit** | Daha az sertleştirilmiş kernel (AppArmor/SELinux kapalıysa) | Kernel exploit → fiziksel bellek erişimi → Windows bölümü açılır |
| 9 | **Secure Boot Bypass (Eski Linux)** | BootHole (CVE-2020-10713) gibi GRUB açıkları | Secure Boot aktif olsa bile kötücül GRUB imajı yüklenir |

**Windows+Linux için Kritik Savunma Yapılandırması:**

```bash
# Linux tarafında LUKS2 ile tam disk şifreleme
# (BitLocker'ın Linux karşılığı)

# Yeni kurulumda:
cryptsetup luksFormat --type luks2 \
    --cipher aes-xts-plain64 \
    --key-size 512 \
    --hash sha512 \
    --iter-time 5000 \
    /dev/sda2

# LUKS header yedekle (BitLocker kurtarma anahtarı gibi):
cryptsetup luksHeaderBackup /dev/sda2 --header-backup-file /mnt/usb/luks_header.img

# Mevcut LUKS şifreleme doğrulama:
cryptsetup luksDump /dev/sda2
```

```bash
# Linux tarafında Windows bölümünü otomatik BAĞLAMA (mount etmekten kaçın)
# /etc/fstab'a Windows bölümünü EKLEME
# Gerektiğinde manuel mount:
sudo mount -o ro,noexec /dev/sda1 /mnt/windows   # Yalnızca oku, çalıştırma
```

```bash
# AppArmor (Ubuntu/Debian) veya SELinux (Fedora/RHEL) aktif olmalı
sudo aa-status      # AppArmor aktif mi?
sestatus            # SELinux aktif mi?

# Sudo zaman aşımını kısa tut:
echo "Defaults timestamp_timeout=1" | sudo tee -a /etc/sudoers.d/security
```

```bash
# Linux Güvenlik Duvarı (ufw):
sudo ufw enable
sudo ufw default deny incoming
sudo ufw default deny outgoing
sudo ufw allow out 443   # HTTPS
sudo ufw allow out 53    # DNS
```

**Windows Tarafından Alınması Gereken Ek Önlem:**
```powershell
# Windows tarafında GRUB'un EFI girdisini bloke etmek için Secure Boot DB:
# (ileri seviye — BIOS'ta Secure Boot Custom Mode gerekir)

# Daha pratik yol: BIOS boot order'ı kilitle
# BIOS → Boot Priority → GRUB'u listeden kaldır veya Windows'u ilk sıraya al
# BIOS Admin şifresi → GRUB'un boot önceliği değiştirilememesi için

# Linux bölümünü Windows altında erişilmez yap:
# ext4 dosya sistemi Windows'ta varsayılan okunamaz — avantaj
# Ama üçüncü parti araçlar (Ext2Fsd, WSL) okuyabilir
# Bu araçları KALDIR veya KURMA:
Get-WindowsOptionalFeature -Online | Where-Object {$_.FeatureName -like "*WSL*"}
```

---

## BÖLÜM 2 — DUAL BOOT'TA BİTLOCKER VE TPM SORUNU

### 2.1 — TPM PCR Değerleri ve Dual Boot Çakışması

BitLocker, diskin şifrelerini çözmeden önce TPM'deki **PCR (Platform Configuration Register)** değerlerini ölçer. Bu ölçümler her açılışta doğrulanır.

```
TPM PCR Değerleri ve Ne Ölçtükleri:
PCR 0 → UEFI Firmware
PCR 2 → UEFI Option ROM
PCR 4 → EFI Boot Manager (hangi OS seçildi)  ← DUAL BOOT'TA DEĞİŞİR
PCR 7 → Secure Boot durumu
PCR 11 → Windows Boot Manager ölçümleri      ← WIN10/WIN11 ARASINDA FARKLI

Dual boot'ta her OS geçişi PCR-4 ve PCR-11'i değiştirir.
BitLocker bu değişikliği "yetkisiz müdahale" sayabilir → Kurtarma anahtarı ister.
```

**BitLocker'ın Dual Boot ile Uyumlu Çalışması için:**

```powershell
# 1. PCR profilini kontrol et — hangi PCR'lar kullanılıyor?
manage-bde -protectors -get C:
# "PCR Validation Profile: 0, 2, 4, 11" gibi bir çıktı göreceksin

# 2. Dual boot için PCR-4'ü ÇIKAR (önyükleyici değişimini yoksay):
# (Güvenliği biraz gevşetir ama dual boot'u kullanışlı kılar)
manage-bde -protectors -delete C: -type TPM
manage-bde -protectors -add C: -tpm
# Alternatif: TPM+PIN kullan — PIN sabit olduğu için PCR değişimi önemli olmaz

# EN SAĞLIKLI YÖNTEM (dual boot için):
# TPM+PIN kullan → PCR değişse de PIN doğruluğu garantiler
# Her OS için ayrı PIN → güvenlik seviyesi korunur
```

---

### 2.2 — Dual Boot'ta Bölüm Kilitleme (Güvenlik Rutini)

```powershell
# Windows 11 açılışında Windows 10 bölümünü otomatik kilitle
# Görev Zamanlayıcısı → Tetikleyici: Kullanıcı oturumu açıldığında
# Eylem: Aşağıdaki komutu çalıştır

manage-bde -lock D: -ForceDismount

# Veya startup script olarak:
# C:\Windows\System32\GroupPolicy\Machine\Scripts\Startup\ altına ekle
```

---

## BÖLÜM 3 — AYRIAYRFIZIKSEL DİSK MİMARİSİ (En Güvenli Yapı)

### 3.1 — Neden Ayrı Fiziksel Diskler?

```
Güvenlik Karşılaştırması:

Tek Disk — Çok Bölüm:
  ├─ EFI paylaşımlı → saldırı yüzeyi ORTAK
  ├─ NTFS erişimi çapraz → veri izolasyonu YOK
  ├─ TPM PCR'ları karışık → BitLocker bozuluyor
  ├─ Bir OS etkilenirse diğeri DOĞRUDAN RİSKTE
  └─ Güvenlik notu: ★★☆☆☆

İki Fiziksel Disk:
  ├─ Her diskin kendi EFI'si → BAĞIMSIZ önyükleme
  ├─ Fiziksel erişim yok → çapraz okuma İMKANSIZ
  ├─ TPM PCR'ları ayrı → BitLocker kararlı çalışır
  ├─ Bir disk etkilenirse diğeri DOKUNULMAZ
  └─ Güvenlik notu: ★★★★★
```

### 3.2 — Ayrı Disk Mimarisi Kurulum Kılavuzu

```
Donanım önerisi:
  Disk 1: NVMe SSD (ana OS) — PCIe 4.0 / 3.0
  Disk 2: SATA SSD veya HDD (ikincil OS) — 2.5" veya 3.5"
```

**BIOS Kurulumu:**
```
1. Her iki diski tak
2. BIOS → Boot Priority → Disk 1'i EN ÜSTE al
3. BIOS → Secure Boot → Enabled
4. Her iki diski ayrı ayrı kur (kurulum sırasında sadece bir disk takılı olsun — karışıklık önlemi)
5. BIOS Admin Şifresi belirle → Boot order değiştirilemesin
6. OS geçişi için: BIOS'a gir → Boot Priority manuel değiştir → Seç
```

**PowerShell — Mevcut Diskler ve Bölümleri Görüntüle:**
```powershell
# Hangi disk hangi OS'u barındırıyor?
Get-Disk | Select-Object Number, FriendlyName, Size, BusType, PartitionStyle

# Detaylı bölüm bilgisi:
Get-Partition | Select-Object DiskNumber, PartitionNumber, DriveLetter, Size, Type

# Disk sağlık durumu:
Get-PhysicalDisk | Select-Object FriendlyName, MediaType, HealthStatus, OperationalStatus
```

---

## BÖLÜM 4 — HDD GÜVENLİK YAPILANDIRMASI

### 4.1 — HDD Tehdit Modeli

```
HDD Fiziksel Özellikleri → Güvenlik Sonuçları:

Manyetik Plaka:
  ├─ "Silinen" veri: Disk üzerinde manyetik iz KALIR
  ├─ Üzerine yazılmadan önce forensic araçlarla GERİ ALINAB̈İLİR
  ├─ BitLocker olmadan disk çıkarılırsa TÜM VERİ OKUNAB̈İLİR
  └─ Cold Boot: Manyetik plaka güç kesildikten sonra uzun süre veri tutar (RAM değil HDD)

Platter Hızı (RPM):
  ├─ 5400 RPM: Yazılım BitLocker ile belirgin yavaşlama (%15-25)
  ├─ 7200 RPM: Daha iyi performans (%8-15 yavaşlama)
  └─ AES-NI destekli CPU ile: <%3 yavaşlama
```

### 4.2 — HDD İçin BitLocker Yapılandırması

```powershell
# 1. AES-NI desteğini kontrol et
(Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management" `
    -Name "FeatureSettingsOverride" -ErrorAction SilentlyContinue) |
    Out-Null

# CPU AES desteği:
[System.Security.Cryptography.Aes]::Create().KeySize  # 256 destekleniyor mu?

# 2. HDD için BitLocker — TÜUM DİSK şifrele (UsedSpaceOnly ASLA kullanma HDD'de)
$pin = Read-Host -AsSecureString "BitLocker Pre-Boot PIN"
Enable-BitLocker -MountPoint "C:" `
    -EncryptionMethod XtsAes256 `
    -TPMandPinProtector -Pin $pin `
    -UsedSpaceOnly:$false    # HDD'de ESKİ VERİLER DE şifrelensin

# 3. Şifreleme durumu izle (uzun sürer — büyük HDD'lerde saatler)
while ((Get-BitLockerVolume C:).EncryptionPercentage -lt 100) {
    $pct = (Get-BitLockerVolume C:).EncryptionPercentage
    Write-Host "Şifreleme: %$pct" -NoNewline
    Start-Sleep 10
    Write-Host "`r" -NoNewline
}
Write-Host "Şifreleme tamamlandı."
```

### 4.3 — HDD Pagefile Yapılandırması

```powershell
# HDD'de pagefile parçalanır — sabit boyut belirle
# Bu hem güvenlik hem performans açısından önemlidir

# Otomatik yönetimi kapat, sabit boyut ayarla:
$CS = Get-WmiObject -Class Win32_ComputerSystem
$CS.AutomaticManagedPageFile = $false
$CS.Put()

# RAM boyutuna göre optimal pagefile: RAM/2 ile RAM*1.5 arası
$RAM_GB = [math]::Round((Get-WmiObject Win32_ComputerSystem).TotalPhysicalMemory / 1GB)
$PF_MB  = $RAM_GB * 1024   # RAM'e eşit

$PF = Get-WmiObject -Class Win32_PageFileSetting -Filter "Name='C:\\pagefile.sys'"
if ($PF) {
    $PF.InitialSize = $PF_MB
    $PF.MaximumSize = $PF_MB
    $PF.Put()
}

# Kapanışta temizle — ZORUNLU:
Set-ItemProperty `
    -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management" `
    -Name "ClearPageFileAtShutdown" -Value 1
```

### 4.4 — HDD Güvenli Silme (Disk Kullanımdan Kalkıyorsa)

```
Yöntem 1: Windows Format (en basit ama forensic dayanıklı değil)
  format C: /P:7   → 7 geçişli üzerine yazma (DoD 5220.22-M standardı)

Yöntem 2: DBAN (Darik's Boot and Nuke) — En Güvenli
  → https://dban.org
  → USB'den boot et
  → "DoD Short" (3 geçiş) veya "Gutmann" (35 geçiş) seç
  → Tüm manyetik iz silinir

Yöntem 3: PowerShell ile HDD üzerine yazma (BitLocker açıksa yeterli)
  # BitLocker şifreliyse diski formatlayıp anahtarı imha etmek yeterli
  # Anahtar olmadan içerik matematiksel olarak geri alınamaz
```

### 4.5 — HDD Fiziksel Güvenlik Ek Adımları

```powershell
# HDD S.M.A.R.T. durumunu izle — fiziksel sorun erken tespiti
# Disk arızalanmadan önce güvenli yedek al:
Get-Disk | ForEach-Object {
    $d = $_
    Get-StorageReliabilityCounter -Disk $d |
        Select-Object @{N="Disk";E={$d.FriendlyName}},
                      Temperature, ReadErrorsTotal, WriteErrorsTotal,
                      Wear, PowerOnHours
}

# Otomatik S.M.A.R.T. izleme servisi (CrystalDiskInfo veya HWiNFO64 önerilir)
```

---

## BÖLÜM 5 — SSD GÜVENLİK YAPILANDIRMASI

### 5.1 — SSD Tehdit Modeli

```
SSD NAND Fiziksel Özellikleri → Güvenlik Sonuçları:

Wear Leveling:
  ├─ Controller, yazma ömrünü dengelemek için verileri farklı hücrelere dağıtır
  ├─ "Silinen" bloklar HEMEN SİLİNMEZ — boş olarak işaretlenir, fiziksel veri durur
  ├─ Üstüne yazma komutu aynı fiziksel hücreye yazamaz
  └─ Forensic araçlar NAND chip'i çıkarıp doğrudan okuyabilir (chip-off forensics)

eDrive (Hardware Encryption):
  ├─ Bazı SSD'ler kendi içinde şifreleme yapar
  ├─ Samsung 840/850/860 EVO: Donanım şifreleme KIRILI (CVE-2018-12037/38)
  ├─ BitLocker eDrive'ı algılarsa SW şifreleme YAPMAZ (güvenilmez)
  └─ ÇÖZÜM: eDrive'ı devre dışı bırak, yazılım şifreleme kullan

TRIM:
  ├─ OS'un SSD'ye "bu blok artık kullanılmıyor" demesi
  ├─ Şifreli SSD'de TRIM metadata sızdırabilir (hangi bloklar boş — düşük risk)
  └─ Güvenlik vs performans dengesi: Çoğu durumda TRIM açık bırakılabilir
```

### 5.2 — SSD İçin BitLocker Yapılandırması

```powershell
# 1. KRITIK: Donanım şifrelemeyi KAPAT — yazılım şifrelemeyi ZORLA
$FVEPath = "HKLM:\SOFTWARE\Policies\Microsoft\FVE"
New-Item -Path $FVEPath -Force | Out-Null
Set-ItemProperty -Path $FVEPath -Name "OSHardwareEncryption"                 -Value 0
Set-ItemProperty -Path $FVEPath -Name "FDVHardwareEncryption"                -Value 0
Set-ItemProperty -Path $FVEPath -Name "RDVHardwareEncryption"                -Value 0
Set-ItemProperty -Path $FVEPath -Name "OSAllowSoftwareEncryptionFailover"    -Value 1

# Politikayı uygula:
gpupdate /force

# 2. SSD için BitLocker etkinleştir
$pin = Read-Host -AsSecureString "BitLocker Pre-Boot PIN"
Enable-BitLocker -MountPoint "C:" `
    -EncryptionMethod XtsAes256 `
    -TPMandPinProtector -Pin $pin `
    -UsedSpaceOnly:$false    # SSD'de de tam şifreleme (wear leveling nedeniyle)

# 3. Şifreleme yöntemini doğrula — "XtsAes256" VE "Software" görünmeli
manage-bde -status C:
# "Encryption Method: XTS-AES 256" → doğru
# "Hardware Encryption: None" veya "Software-based" → doğru (eDrive kullanılmıyor)
```

### 5.3 — TRIM Yapılandırması

```powershell
# TRIM durumunu kontrol et:
fsutil behavior query DisableDeleteNotify
# 0 = TRIM aktif (DOĞRU)
# 1 = TRIM devre dışı

# Eğer kapalıysa aç:
fsutil behavior set DisableDeleteNotify 0

# Manuel TRIM çalıştır (haftada bir öneririm):
Optimize-Volume -DriveLetter C -ReTrim -Verbose

# Otomatik optimizasyon zamanlaması doğrula:
Get-ScheduledTask -TaskName "ScheduledDefrag" |
    Select-Object TaskName, State
# SSD için: Windows TRIM çalıştırır, disk birleştirme YAPMAZ (otomatik)
```

### 5.4 — SSD Güvenli Silme (Disk Kullanımdan Kalkıyorsa)

```
Yöntem 1: ATA Secure Erase (En Güvenli — Donanım Seviyesi)
  ├─ SSD controller'ın tüm hücreleri sıfırlamasını sağlar
  ├─ Wear-leveled bloklardaki "eski" verileri de siler
  └─ Araç: hdparm (Linux) veya üretici yazılımı

Yöntem 2: Üretici Yazılımı (En Kolay)
  Samsung SSD → Samsung Magician → "Secure Erase"
  Crucial SSD → Crucial Storage Executive → "Sanitize"
  WD SSD      → WD Dashboard → "Erase"
  Intel SSD   → Intel MAS (Memory and Storage Tool)

Yöntem 3: NVMe Format (NVMe SSD için)
  Linux bootable USB → nvme format /dev/nvme0 --ses=1
  → ses=1: Cryptographic erase (şifreleme anahtarını imha eder)
  → ses=2: User data erase (fiziksel üzerine yazma)

Yöntem 4: BitLocker Kullanıyorsan (En Hızlı)
  → BitLocker şifreleme anahtarını imha et:
  manage-bde -delete C: -type RecoveryPassword
  manage-bde -delete C: -type TPM
  → Disk şifreli ama anahtar yok → içerik matematiksel olarak geri alınamaz
  → Format at, yeni OS kur
```

### 5.5 — SSD'ye Özgü Ek Güvenlik

```powershell
# Defrag'ı KALDIRMA sadece DOĞRULA — Win10/11 SSD'yi tanır ve defrag yapmaz:
$OptService = Get-Service -Name "defragsvc"
Write-Host "Optimize Disk Service: $($OptService.Status)"
# Servis çalışıyor olabilir ama SSD'ye TRIM yapar, defrag yapmaz

# SSD NVMe güvenlik durumu (varsa):
# NVMe Sanitize komutu desteği:
# Get-PhysicalDisk | ForEach-Object { $_.FriendlyName }

# SSD sıcaklık izleme (aşırı ısı veri bütünlüğünü etkiler):
Get-PhysicalDisk | ForEach-Object {
    Get-StorageReliabilityCounter -Disk $_ |
        Select-Object @{N="Disk";E={$_.FriendlyName}}, Temperature, Wear
}
```

---

## BÖLÜM 6 — HDD vs SSD GÜVENLİK KARŞILAŞTIRMA TABLOSU

| Özellik | HDD | SSD |
|---------|-----|-----|
| **Veri Kalıntısı** | Manyetik iz — silinen veri KALIR | NAND hücre — wear leveling nedeniyle KALIR |
| **Forensic Riski** | Orta (platter görüntüleme) | Yüksek (chip-off — maliyet yüksek) |
| **BitLocker Şifreleme** | Yazılım (AES-NI ile yeterli) | Yazılım ZORLA (eDrive güvenilmez) |
| **`UsedSpaceOnly:$false`** | ✅ **Zorunlu** | ✅ **Zorunlu** (Wear leveling'e karşı) |
| **TRIM** | ❌ Yok | ✅ Aktif olmalı |
| **Güvenli Silme** | DBAN / 7 geçiş format | ATA Secure Erase / Üretici aracı |
| **Pagefile boyutu** | Sabit belirle (parçalanma önlemi) | Varsayılan (Windows yönetimine bırak) |
| **BitLocker Performans** | 5400RPM: %15-25, 7200RPM: %8-15 yavaşlama | AES-NI ile <%1 fark |
| **Cold Boot Riski** | RAM'deki kadar değil; plaka güçsüz kalınca hızlı bozulur | NAND güçsüz kalınca hücreler korunur (daha uzun) |
| **Defrag** | ✅ Gerekli | ❌ KAPAT (Windows otomatik) |
| **S.M.A.R.T. İzleme** | ✅ Önemli (mekanik arıza erken uyarı) | ✅ Önemli (Wear değeri izle) |

---

## BÖLÜM 7 — TAM ÖZET VE TAVSİYE MATRİSİ

### Dual Boot Kombinasyonları — Güvenlik Skoru

| Kombinasyon | Disk | Güvenlik Skoru | Öneri |
|-------------|------|----------------|-------|
| Win11 + Win11 | Tek disk | ★★☆☆☆ | Kaçın veya her iki bölümü BitLocker+PIN ile kilitle |
| Win10 + Win10 | Tek disk | ★★☆☆☆ | Kaçın; SAM çapraz erişimi kritik risk |
| Win10 + Win11 | Tek disk | ★☆☆☆☆ | **Kaçın** — en zayıf kombinasyon |
| Win + Linux | Tek disk | ★☆☆☆☆ | **Kaçın** — GRUB, NTFS-3G, pypykatz riskleri |
| Win11 + Win11 | Ayrı disk | ★★★★☆ | Kabul edilebilir |
| Win10 + Win11 | Ayrı disk | ★★★☆☆ | Makul — ayrı PIN zorunlu |
| Win + Linux | Ayrı disk | ★★★☆☆ | Makul — LUKS2 + AppArmor/SELinux zorunlu |
| Tek OS | Tek disk | ★★★★★ | **En güvenli** |

### Temel Kurallar (Tüm Senaryolar İçin)

```
1. Ayrı fiziksel disk → en yüksek güvenlik
2. Her OS için ayrı BitLocker PIN (aynı PIN'i kullanma)
3. Kullanılmayan OS bölümünü kilitle: manage-bde -lock X: -ForceDismount
4. BIOS Admin şifresi → boot order kilitle
5. HDD: UsedSpaceOnly:$false + Sabit pagefile + DBAN ile silme
6. SSD: eDrive kapat + TRIM aktif + Üretici aracıyla silme
7. Win+Linux: LUKS2 + AppArmor/SELinux + sudo timeout kısa
8. Dual boot'ta GRUB yerine Windows Boot Manager tercih et (Secure Boot uyumluluğu)
```

---

*Son güncelleme: 2026-03-22 · Windows 10/11 Pro · Ubuntu/Fedora Linux için hazırlanmıştır.*
