# 💾 HARİCİ USB SSD'DEN WINDOWS BOOT REHBERİ
## Windows To Go · Rufus · WinToUSB · Güvenlik Analizi

> **Amaç:** Windows 10/11'i harici USB SSD kutusundan başlatmak, hem iç hem dış disk olarak kullanmak ve bu yapının güvenliğini tam anlamıyla anlamak.

---

## 🔍 NEDEN STANDART WINDOWS HARİCİ DİSKTEN BOOT ETMİYOR?

Senin yaşadığın tam olarak bu sorun. İç SSD'yi kutusuna koyup USB ile bağladın, BIOS'tan seçtin ama boot olmadı. Bunun tek bir nedeni var:

```
Normal Windows Kurulumu → Yalnızca SATA/NVMe sürücülerine özelleştirilmiş
                                          ↓
  Başlangıç sürücüsü: storahci.sys (SATA) veya stornvme.sys (NVMe)
                                          ↓
  Bu sürücüler USB denetleyicisini (USB Mass Storage / UAS) TANIMAZ
                                          ↓
  Boot loader diskten okumaya çalışır → USB sürücüsü yüklü değil → HATA
                                          ↓
  0xc000000e / BSOD veya sessiz başarısızlık
```

**Kısaca:** Standart kurulumda Windows, "ben USB üzerinden çalışacağım" diye yapılandırılmamıştır. Bunun için özel yöntemler gerekir.

---

## 📐 MİMARİ: 3 YÖNTEM KARŞILAŞTIRMASI

```
┌────────────────────────────────────────────────────────┐
│  YÖNTEM 1: Windows To Go (Resmi)                       │
│  Win 10 Enterprise lisansı gerekir                     │
│  Sertifikalı USB 3.0 flash/SSD gerekir                │
│  Güvenlik: ★★★★☆                                       │
├────────────────────────────────────────────────────────┤
│  YÖNTEM 2: Rufus "Windows To Go" Modu                  │
│  Tüm Win 10/11 sürümleri (Pro dahil)                  │
│  Herhangi USB 3.x SSD kutusu                          │
│  Güvenlik: ★★★★☆                                       │
├────────────────────────────────────────────────────────┤
│  YÖNTEM 3: WinToUSB (Üçüncü Taraf)                    │
│  GUI arayüzü, mevcut kurulumu klonlayabilir           │
│  Güvenlik: ★★★☆☆ (lisans doğrulama riski)            │
└────────────────────────────────────────────────────────┘
```

---

## BÖLÜM 1 — DONANIM GEREKSİNİMLERİ

### 1.1 — USB SSD Kutusu (Enclosure) Seçimi

Harici disk kutusunun kalitesi boot başarısını doğrudan etkiler.

| Özellik | Minimum | Önerilen | Neden Önemli? |
|---------|---------|---------|---------------|
| **USB Sürümü** | USB 3.0 (5 Gbps) | USB 3.2 Gen 2 (10 Gbps) | Daha hızlı = sistem daha akıcı |
| **Protokol** | USB-MSC | **UASP** (USB Attached SCSI Protocol) | UASP gecikmeyi %50 azaltır |
| **Bağlantı Tipi** | USB-A | **USB-C** | Hem bilgisayar hem telefon/tablet |
| **Çip seti** | JMicron, ASMedia | **ASMedia ASM2362** (NVMe için) | Kararlılık ve uyumluluk |
| **NVMe/SATA** | M.2 SATA | **M.2 NVMe** | NVMe çok daha hızlı |

**Önerilen Kutular:**
```
NVMe için:
  ├─ ORICO M2PAC3-G20 (USB4 / 20Gbps — en hızlı)
  ├─ Sabrent EC-SNVE (USB 3.2 Gen 2 / 10Gbps)
  └─ UGREEN CM400 (USB-C 3.2 Gen 2)

SATA M.2 için:
  ├─ ORICO M2PV-C3 (USB 3.2 Gen 2)
  └─ Inateck FE2011 (USB 3.0, bütçe dostu)
```

**UASP Desteğini Kontrol Et:**
```powershell
# Bağlı USB depolama cihazlarını ve protokolü göster
Get-PnpDevice -Class "USB" | Where-Object {$_.FriendlyName -like "*UAS*"}

# Disk yöneticisinde detay
Get-Disk | Select-Object Number, FriendlyName, BusType
# BusType: USB → bağlı
```

### 1.2 — Bilgisayarın USB Port Gereksinimleri

```
USB 2.0 (480 Mbps) → Boot çalışır ama çok YAVAŞ — kullanılamaz
USB 3.0 (5 Gbps)   → Minimum kabul edilebilir
USB 3.1 Gen 2 (10 Gbps) → İyi
USB 3.2 Gen 2x2 (20 Gbps) → Çok iyi
USB4 / Thunderbolt 3/4 (40 Gbps) → Mükemmel (iç disk hızına yakın)

Port rengine bak:
Mavi port → USB 3.0 (5 Gbps) ✓
Kırmızı/Sarı port → Genellikle "always on" veya 3.1 ✓
Normal siyah → USB 2.0 ✗
```

> [!TIP]
> **Her zaman aynı USB portunu kullan.** Windows, cihaz sürücü profilini port bazlı kaydeder. Farklı port → farklı donanım profili → sürücü yeniden yükleme → ilk açılış yavaşlığı ve nadiren BSOD. En hızlı (USB 3.x mavi/kırmızı) portu belirle ve o portu kalıcı kullan.

### 1.3 — "Hem İç Hem Dış Disk" Yapısı

Kullanıcının istediği tam bu yapı — SSD'yi çıkarmadan hem iç hem USB olarak kullanmak:

```
İSTENEN SENARYO:
┌─────────────────────────────────────────────────────────┐
│  Bilgisayar                                             │
│                                                          │
│  M.2 Yuvası: [SSD takılı] ──────── İÇ DISK olarak      │
│                                     çalışır             │
│        │                                                 │
│        └── Ayrıca USB kablosu ile ── DIŞ DISK olarak   │
│            kasanın dışına çıkar     da çalışır ?        │
└─────────────────────────────────────────────────────────┘

GERÇEK:
Bir disk aynı anda hem SATA/NVMe hem USB olarak
bağlanamaz — fiziksel olarak imkansız.

ÇÖZÜM:
İki ayrı SSD kullan:
  Disk 1 (M.2 iç): Ana OS — her zaman içerde
  Disk 2 (Kutu + USB): Taşınabilir OS — dışardan bağlan
```

---

## BÖLÜM 2 — YÖNTEM A: RUFUS ile Windows To Go (Önerilen)

**Rufus**, Microsoft'un resmi Windows To Go özelliğini tüm Windows sürümlerine (Pro dahil) ücretsiz açar.

### 2.1 — Gereksinimler

```
✅ Rufus 3.x veya üzeri (https://rufus.ie)
✅ Windows 10/11 ISO (Microsoft'tan temiz)
✅ Hedef USB SSD — en az 64 GB, tercihen 256 GB+
✅ USB 3.0+ port
```

### 2.2 — Adım Adım Kurulum

```
1. Rufus'u yönetici olarak aç

2. Cihaz: Listeden USB SSD'ni seç
   (DİKKAT: Yanlış diski seçme — tüm veri silinir)

3. Önyükleme seçimi düğmesine tıkla → "Disk or ISO image"
   → ISO dosyasını seç (Win10 veya Win11)

4. Bölüm düzeni: GPT
   Hedef sistem: UEFI (CSM olmayan)

5. ÖNEMLİ ADIM:
   "Görüntü seçeneği" açılır menüsü:
   → "Windows To Go" seçeneğini seç
      (Bu seçenek USB sürücülerini önyükleme imajına gömer)

6. Dosya sistemi: NTFS
   Küme boyutu: Varsayılan

7. "BAŞLAT" → Uyarıyı onayla → Bekle (20-40 dk, SSD hızına göre)
```

> [!IMPORTANT]
> "Windows To Go" seçeneği yalnızca USB SSD seçildiğinde görünür. Normal USB bellek seçilirse bu seçenek gelmez (boyut/hız yetersiz olarak algılanır).

### 2.3 — İlk Açılış Yapılandırması

```
1. BIOS → Boot Priority → USB SSD'ni seçili yap
2. Bilgisayarı başlat → Windows OOBE (kurulum sihirbazı) açılır
3. Bölge, dil seç
4. KRİTİK: İnternete BAĞLANMA — yerel hesap oluştur
5. Gizlilik ayarlarını kapat
6. İlk masaüstüne geç

Sonra Katman 2-7 güvenlik adımlarını uygula (WINDOWS11_HARDENING_KALE.md)
```

### 2.4 — Rufus Windows To Go — Teknik Arka Plan

Rufus, standart ISO kurulumundan şu farkı yapar:

```
Normal Kurulum:
  setup.exe → yalnızca SATA/NVMe için storahci.sys/stornvme.sys kur

Rufus WTG Modu:
  1. Tüm Windows dosyalarını USB'ye kopyala
  2. USB sürücüsünü (usbstor.sys, UASPSTOR.sys) boot başlangıç
     sürücüleri listesine (CRITICAL_DEVICE_DATABASE) ekle
  3. BootCampUSB benzeri bir yapı oluşturarak USB → NTFS boot 
     zinciri kur (UEFI: \EFI\Microsoft\Boot\bootmgfw.efi)
  4. Registry: HKLM\SYSTEM\CurrentControlSet\Services\usbstor
     → Start = 0 (BOOT_START) yap
     Normal kurulumda bu değer 3 (DEMAND_START) — USB boot'ta geç yükleniyor
```

---

## BÖLÜM 3 — YÖNTEM B: WinToUSB (GUI Arayüzlü)

### 3.1 — WinToUSB Nedir?

EasyUEFI'nin aracı. Hem yeni kurulum hem de mevcut iç disk kurulumunu harici diske **klonlayabilir**.

```
https://www.easyuefi.com/wintousb/
Ücretsiz sürüm: Windows 10/11 Home/Pro destekler
```

### 3.2 — Mevcut Kurulumu Klonlama

Bu senin asıl istediğin senaryo:

```
1. WinToUSB'yi yönetici olarak aç
2. Sol panel: "Windows'u bu bilgisayardan klonla"
3. Kaynak: C: (mevcut Windows kurulumu)
4. Hedef: USB SSD'yi seç
5. Kurulum modu:
   → "UEFI" → GPT bölüm yapısı
6. Başlat → Klonlama tamamlanır (1-3 saat, boyuta göre)
7. USB SSD'den boot et → Klonlanmış sistem açılır
```

> [!WARNING]
> Klonlama sonrası yeni sistem farklı bir donanım profili algılar → Windows aktivasyonu sorunu çıkabilir. Dijital lisans Microsoft hesabına bağlıysa çözülür; OEM lisans bağlanmayabilir.

### 3.3 — WinToUSB ile Temiz Kurulum

```
1. WinToUSB → "ISO veya DVD'den kur"
2. ISO seç → Windows sürümünü seç
3. Hedef USB SSD'yi seç → Bölüm düzeni: GPT/UEFI
4. Windows To Go modu seç
5. Yükle → İlk açılışta OOBE gelir
```

---

## BÖLÜM 4 — YÖNTEM C: Resmi Windows To Go (Win 10 Enterprise)

Microsoft'un resmi özelliği — **yalnızca Enterprise lisans** ile çalışır.

```powershell
# Windows 10 Enterprise — dahili araç
# Denetim Masası → Windows To Go

# PowerShell ile kontrol
Get-WindowsOptionalFeature -Online -FeatureName "EnterpriseWTG"
```

**Sertifikalı Donanım (Resmi WTG için):**
- Kingston DataTraveler Workspace
- IronKey W700/W710
- Super Talent Express RC8/RC4

> [!NOTE]
> Windows 11 Enterprise de Windows To Go özelliğini **kaldırdı** (21H2 ile deprecated). Rufus yöntemi artık Microsoft'un önerdiği alternatif.

---

## BÖLÜM 5 — BOOT PERFORMANSI

### 5.1 — Hız Karşılaştırması

```
Depolama → USB Bağlantısı → Gerçek Okuma Hızı → Açılış Süresi*

İç NVMe SSD (PCIe 4.0)        → 7000 MB/s → ~5-8 sn
Thunderbolt 4 Harici NVMe      → 3200 MB/s → ~8-12 sn
USB4 (40Gbps) Harici NVMe      → 3000 MB/s → ~10-14 sn
USB 3.2 Gen 2 (10Gbps) NVMe   → 900 MB/s  → ~15-25 sn  ← Minimum kabul edilebilir
USB 3.2 Gen 1 (5Gbps) NVMe    → 400 MB/s  → ~30-45 sn
USB 3.0 SATA SSD               → 400 MB/s  → ~35-50 sn
USB 3.0 USB Flash Bellek        → 100 MB/s  → ~2-5 dk   ← Kullanılamaz
USB 2.0                         → 40 MB/s   → ~10+ dk   ← Kesinlikle kullanma

*Cold boot, SSD varış sonrası masaüstüne kadar süre
```

### 5.2 — UASP vs BOT (Bulk-Only Transfer) Farkı

```powershell
# UASP destekleniyor mu?
Get-PnpDevice -Class "SCSIAdapter" |
    Where-Object {$_.FriendlyName -like "*UAS*"} |
    Select-Object FriendlyName, Status

# Disk gecikmesi ölç (CrystalDiskMark alternatifi):
$disk = "C:"  # veya USB sürücü harfi
$testFile = "$disk\speedtest.tmp"
$size = 512MB
$sw = [System.Diagnostics.Stopwatch]::StartNew()
[System.IO.File]::WriteAllBytes($testFile, (New-Object byte[] $size))
$sw.Stop()
Remove-Item $testFile -Force
Write-Host "Yazma hızı: $([math]::Round($size / $sw.Elapsed.TotalSeconds / 1MB, 1)) MB/s"
```

---

## BÖLÜM 6 — GÜVENLİK YAPILANDIRMASI

### 6.1 — BitLocker ile Harici Disk Şifreleme

> [!IMPORTANT]
> Windows To Go diskini mutlaka BitLocker ile şifrele. Disk kaybolursa veya çalınırsa içindeki tüm sistem erişilebilir olacak.

```powershell
# Windows To Go diskine (D: harfi verildiyse) BitLocker uygula
# NOT: WTG disk hem sistem hem de veri diski — sistem diski olarak şifreleme:

# Grup İlkesi ayarı (hedef disk üzerinde)
# Halihazırda bu disk Windows To Go olarak çalışıyorken yapılır:
$pin = Read-Host -AsSecureString "BitLocker Pre-Boot PIN (USB Disk)"
Enable-BitLocker -MountPoint "C:" `
    -EncryptionMethod XtsAes256 `
    -TPMandPinProtector `
    -Pin $pin `
    -UsedSpaceOnly:$false

# SORUN: TPM farklı bilgisayarda çalışmaz!
# Çözüm: TPM yerine yalnızca şifre (parola) tabanlı koruma:
Enable-BitLocker -MountPoint "C:" `
    -EncryptionMethod XtsAes256 `
    -PasswordProtector `
    -Password (Read-Host -AsSecureString "Disk Şifresi") `
    -UsedSpaceOnly:$false

# Her iki yöntemi birlikte ekle:
Add-BitLockerKeyProtector -MountPoint "C:" -RecoveryPasswordProtector
```

> [!NOTE]
> Windows To Go, **birden fazla bilgisayarda** kullanıldığı için TPM bağlı koruma çalışmaz — her bilgisayarın TPM çipi farklıdır. Bu nedenle **Parola tabanlı BitLocker** kullanılır. Güvenlik açısından TPM+PIN'den daha zayıftır; bunu göz önünde bulundur.

### 6.2 — Windows To Go Politikası — Veri Sızıntısı Önleme

Windows To Go aktifken ana bilgisayarın disklerine erişimi kısıtla:

```
Grup İlkesi (gpedit.msc):
Bilgisayar Yapılandırması → Yönetim Şablonları → Windows Bileşenleri
  → Taşınabilir İşletim Sistemi
     → "Windows To Go için Ana Bilgisayar Disklerinin Bağlanmasını Engelle"
     → ETKİN

Bu ayar aktifken:
  → Ana bilgisayarın C:, D: vb. diskleri WTG içinden görünmez
  → Veri sızıntısı ve çapraz disk erişimi engellenir
```

```powershell
# Registry ile (Home sürümü dahil):
$WTGPath = "HKLM:\SOFTWARE\Policies\Microsoft\PortableOperatingSystem"
New-Item -Path $WTGPath -Force | Out-Null
Set-ItemProperty -Path $WTGPath -Name "MountVolumes" -Value 0
# 0 = Ana bilgisayar diskleri bağlanmasın

# PortableOperatingSystem = 1 — Sistemi "taşınabilir" modda işaretle
# Sık sök-tak yapacaksan bu değeri ayarla:
# Windows bu değeri 1 görünce her açılışta donanım keşfi yapar,
# farklı bilgisayarlarda sürücü kurulumunu otomatikleştirir
Set-ItemProperty -Path $WTGPath -Name "PortableOperatingSystem" -Value 1
# Doğrulama:
Get-ItemProperty $WTGPath | Select-Object MountVolumes, PortableOperatingSystem
```

> [!NOTE]
> `PortableOperatingSystem = 1` değeri Windows'a "bu sistem farklı donanımlarda çalışacak" sinyali verir. Sürücü uyum sorunlarını azaltır ve farklı bilgisayarlarda ilk açılışı hızlandırır. Tek bir bilgisayarda kullanıyorsan bu değer gerekli değildir.

### 6.3 — Hızlı Başlatma (Fast Startup) — MUTLAKA KAPAT

> [!IMPORTANT]
> Bu adım sık sök-tak yapanlar için **kritiktir ve belgenin en önemli pratikal detayıdır.** Fast Startup aktifken Windows tam kapanış yapmaz; RAM durumunu hiberfil.sys'e yazar. USB diski bu durumdayken farklı bir bağlantıyla (veya farklı PC'de) açmaya çalışırsan tutarsız disk durumu → `INACCESSIBLE_BOOT_DEVICE` BSOD veya dosya sistemi hatası.

```
Denetim Masası → Güç Seçenekleri
  → "Güç düğmelerinin yapacaklarını seçin"
  → "Şu anda kullanılamayan ayarları değiştirin" tıkla
  → "Hızlı başlatmayı aç (önerilen)" kutusunu KALDIR
  → Değişiklikleri Kaydet
```

```powershell
# PowerShell ile (yönetici):
# Fast Startup'ı kapat:
Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Power" `
    -Name "HiberbootEnabled" -Value 0

# Hibernate'i de kapat (WTG disklerde kesinlikle kapalı olmalı):
powercfg /hibernate off

# Doğrulama:
(Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Power").HiberbootEnabled
# 0 olmalı

Test-Path "C:\hiberfil.sys"   # False olmalı
```

**Teknik neden:**
```
Fast Startup AÇIK:
  Kapatma → RAM snapshot → hiberfil.sys → Disk "kilitli" durumda bekler
  Başka PC'de USB SSD'yi açmaya çalış:
  → Disk önceki sistemin snapshot'ını görmek ister
  → Snapshot'ı bulamaz → BSOD / Dosya sistemi hatası

Fast Startup KAPALI:
  Kapatma → Tam temiz kapanış → Disk bağımsız ve taşınabilir durumda
  Başka PC → Sorunsuz açılır
```

### 6.4 — Windows To Go için Ağ Yalıtımı

Farklı bilgisayarlarda kullandığında ağ kimlik bilgilerinin karışmaması için:

```powershell
# Her açılışta DNS önbelleğini temizle
$Task = @"
<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2">
  <Triggers><BootTrigger><Enabled>true</Enabled></BootTrigger></Triggers>
  <Actions>
    <Exec><Command>ipconfig</Command><Arguments>/flushdns</Arguments></Exec>
  </Actions>
</Task>
"@
Register-ScheduledTask -TaskName "FlushDNS_Boot" -Xml $Task -Force

# Windows Credential Manager'ı temizle (oturum kapanışında)
$LogoffTask = @"
<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2">
  <Triggers><SessionStateChangeTrigger><StateChange>RemoteDisconnect</StateChange></SessionStateChangeTrigger></Triggers>
  <Actions>
    <Exec><Command>cmdkey</Command><Arguments>/list | ForEach {cmdkey /delete:$_}</Arguments></Exec>
  </Actions>
</Task>
"@
```

---

## BÖLÜM 7 — İÇ DİSK vs DIŞ DİSK: GÜVENLİK KARŞILAŞTIRMASI

Bu bölüm en kritik kısım:

### 7.1 — TPM Bağlaması Sorunu

```
İÇ DİSK (NVMe/SATA):
  TPM çipi → PCR ölçümleri → BitLocker anahtarını mühürler
  Disk başka bilgisayara takılsa bile TPM eşleşmez → Açılmaz
  Güvenlik: ★★★★★

DIŞ DİSK (USB):
  TPM bağlaması MÜMKÜN DEĞİL (farklı bilgisayarlarda açılmalı)
  BitLocker yalnızca PAROLA tabanlı çalışır
  Disk çalınırsa → Parola kırılabilir (brute-force)
  Güvenlik: ★★★☆☆
```

### 7.2 — Fiziksel Güvenlik

```
İÇ DİSK:
  ├─ Kasanın içinde, fiziksel erişim zor
  ├─ Hırsızlık = tüm bilgisayarı çalmak demek
  └─ Risk: Düşük

DIŞ DİSK (USB Kablo ile):
  ├─ Kablo takılıyken fiziksel olarak çıkarılabilir
  ├─ Yalnızca USB SSD kutusunu çalmak yeterli
  ├─ Kafede, ofiste, ulaşımda — risk her zaman yüksek
  └─ Risk: YÜKSEK — BitLocker şifreleme zorunlu
```

### 7.3 — Hız ve Gecikme Farkı

```
İÇ NVMe (PCIe 4.0) → 7000 MB/s okuma → Sistem akıcı
USB 3.2 Gen 2 NVMe → 900 MB/s okuma  → Fark hissedilir ama kullanılabilir
USB 3.0 SATA SSD   → 400 MB/s okuma  → Minimum kabul edilebilir
```

### 7.4 — Kapsamlı Karşılaştırma Tablosu

| Güvenlik Özelliği | İç Disk | Harici USB Disk |
|------------------|---------|-----------------|
| **TPM + PIN** | ✅ Tam destek | ❌ Desteklenmez (çok PC) |
| **BitLocker** | ✅ XTS-AES 256 + PIN | ⚠️ Yalnızca Parola |
| **Pre-Boot Doğrulama** | ✅ PIN (Windows gelmeden önce) | ⚠️ Parola (daha zayıf) |
| **Fiziksel Hırsızlık** | Düşük risk | **Yüksek risk** |
| **Sahte Disk Takma** | Zor (case açmak gerekir) | Kolay (USB çek-tak) |
| **Veri Yalıtımı** | Tam (iç disk — başka OS göremez) | Politika ile sağlanabilir |
| **Farklı PC'de Kullanım** | ❌ TPM uyuşmazlığı | ✅ Esnek |
| **Hız** | ★★★★★ | ★★★☆☆ (USB sınırı) |
| **Cold Boot Riski** | Normal | Daha yüksek (taşınabilir) |
| **BIOS Kilidi Etkisi** | ✅ BIOS şifresi korur | ⚠️ Başka PC'de BIOS şifresi yok |

### 7.5 — Hangi Senaryo İçin Hangi Çözüm?

| Kullanım Amacı | Önerim |
|---------------|--------|
| Masaüstü — Asla yerinden kalkmıyor | İç disk — tam güvenlik |
| Laptop — Evde/ofiste sabit | İç disk |
| Laptop — Sık seyahat, farklı PC'de kullanım | Harici disk + güçlü BitLocker parolası |
| Güvenli taşınabilir OS (farklı bilgisayarlar) | WTG + Parola BitLocker + Veri izolasyonu |
| Maksimum güvenlik + taşınabilirlik | Thunderbolt 4 SSD + FIDO2 anahtar |

---

## BÖLÜM 8 — BIOS AYARLARI: USB BOOT

### 8.1 — BIOS'ta USB Boot Etkinleştirme

```
BIOS → Boot → Boot Priority:
  1. USB SSD'yi listede en üste sürükle (veya + tuşu ile yukarı al)
  2. Internal SSD'yi ikinci sıraya bırak
  Kaydet (F10) → Çık

Alternatif — her seferinde BIOS'a girmeden:
  Bilgisayar başlarken F12 / F11 / Esc / F8 (üreticiye göre)
  → "Boot Menu" açılır → USB SSD seç
  → Bu BIOS ayarlarını değiştirmez, tek seferlik seçim
```

### 8.2 — CSM (Compatibility Support Module) Sorunu

```
UEFI + Secure Boot: CSM KAPALI olmalı
  → GPT disk + UEFI önyükleyici → Sorunsuz çalışır

CSM AÇIKSA ve Legacy boot seçiliyse:
  → MBR gerektir → Rufus WTG işlemi GPT kullanır → ÇAKIŞMA
  → Çözüm: BIOS → CSM → Disabled

Kontrol:
  msinfo32 → BIOS Modu: UEFI görünmeli (Legacy değil)
```

```powershell
# UEFI mi Legacy mi?
$FW = (Get-ItemProperty "HKLM:\System\CurrentControlSet\Control" -Name "PEFirmwareType").PEFirmwareType
if ($FW -eq 2) { Write-Host "UEFI ✓" -ForegroundColor Green }
else           { Write-Host "Legacy BIOS — UEFI'ye geç!" -ForegroundColor Red }
```

### 8.3 — Secure Boot ile USB Boot

```
Secure Boot AÇIKKEN:
  Rufus WTG → Microsoft imzalı önyükleyici → Secure Boot geçer ✅
  Kendi oluşturduğun imzasız imaj → Secure Boot reddeder ❌

Secure Boot ile USB'den boot için:
  1. BIOS → Secure Boot → Standard (varsayılan anahtarlar)
  2. Rufus WTG modu kullan → İmzalı önyükleyici otomatik eklenir
  3. Boot et → Sorunsuz çalışır

Sorun çıkarsa:
  msinfo32 → Güvenli Önyükleme Durumu: Açık görünmeli
  Değilse: BIOS → Secure Boot → Restore Factory Keys → Kaydet
```

---

## BÖLÜM 9 — SORUN GİDERME

### "USB'den boot et, ama açılmıyor" — Hata Listesi

| Hata | Olası Neden | Çözüm |
|------|-------------|-------|
| `0xc000000e` | USB sürücüsü boot başlangıcında yüklü değil | Rufus WTG modu ile yeniden oluştur |
| `0xc0000225` | BCD (Boot Config) bozuk | `bcdboot X:\Windows /s X: /f UEFI` |
| Siyah ekran | UEFI/Legacy uyumsuzluğu | BIOS → CSM → Disabled |
| Secure Boot hatası | İmzasız önyükleyici | Rufus WTG kullan, CSM kapat |
| Yavaş / donuyor | USB 2.0 kullanılıyor | USB 3.0+ porta tak |
| Activation hatası | Donanım değişti | Microsoft hesabı ile yeniden bağla |

### BCD Onarma (Rufus Sonrası Hata)

```powershell
# Windows PE veya kurulum USB'sinden:
# Shift+F10 → Komut İstemi

# Hangi disk Windows? (harfleri bul)
diskpart
list vol
exit

# BCD yeniden oluştur (X harfi USB disk)
bcdboot X:\Windows /s X: /f UEFI /l tr-TR

# BCD doğrula
bcdedit /store X:\EFI\Microsoft\Boot\BCD /enum
```

---

## BÖLÜM 10 — NİHAİ KONTROL LİSTESİ

```
DONANIM HAZIRLIĞI:
  □ USB 3.0+ kutu — UASP destekli
  □ M.2 NVMe SSD (en az 128 GB, önerilen 256 GB+)
  □ USB 3.0+ port bilgisayarda mevcut

RUFUS WTG KURULUMU:
  □ Windows 10/11 Pro ISO indirildi (Microsoft'tan)
  □ Rufus 3.x+ indirildi (rufus.ie)
  □ Hedef disk seçildi (doğru disk!)
  □ "Windows To Go" modu seçildi
  □ GPT + UEFI seçili
  □ Kurulum tamamlandı (~20-40 dk)

BIOS AYARLARI:
  □ CSM → Disabled
  □ Secure Boot → Enabled (Standard)
  □ Boot Priority → USB SSD en üstte
  □ BIOS Admin Şifresi → Koyuldu

İLK AÇILIŞ:
  □ İnternetsiz kurulum
  □ Yerel hesap oluşturuldu
  □ Gizlilik ayarları kapatıldı

HİBRİT KULLANIM (Sık Sök-Tak):
  □ Fast Startup (Hızlı Başlatma) → KAPATILDI
  □ Hibernate → KAPATILDI (powercfg /hibernate off)
  □ PortableOperatingSystem = 1 → Ayarlandı
  □ Her zaman aynı USB 3.x portu → Belirlendi ve işaretlendi

GÜVENLİK ADIMI:
  □ BitLocker → Parola tabanlı (TPM yok) → XTS-AES 256
  □ 48 haneli kurtarma anahtarı → Kasaya kaldırıldı
  □ Ana bilgisayar disk erişimi → Grup İlkesi ile kapatıldı
  □ Windows 11 Hardening adımları uygulandı
     (WINDOWS11_HARDENING_KALE.md'e bak)
```

---

## ÖZET

| Konu | Sonuç |
|------|-------|
| Neden daha önce boot olmadı? | Standart Windows USB sürücüsünü boot anında tanımaz |
| Çözüm | Rufus → Windows To Go modu → USB sürücüleri başlangıca gömülür |
| Güvenlik farkı (iç vs dış) | İç disk TPM bağlaması → çok güçlü; Dış disk yalnızca parola → daha zayıf |
| Minimum donanım | USB 3.0 kutu (UASP) + NVMe SSD |
| Önerilen donanım | USB 3.2 Gen 2 / USB4 kutu + PCIe 4.0 NVMe |
| BitLocker yöntemi | Parola tabanlı (TPM birden fazla PC'de çalışmaz) |
| Hız kaybı | USB 3.2 Gen 2 ile %70-80 iç disk hızı alınabilir |

---

*Son güncelleme: 2026-03-22 · Windows 10 22H2 / Windows 11 23H2 için hazırlanmıştır.*
