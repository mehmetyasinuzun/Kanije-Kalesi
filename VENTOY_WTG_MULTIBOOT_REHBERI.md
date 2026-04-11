# 🧰 VENTOY + WINDOWS TO GO + LINUX — TAM MULTIBOOT REHBERİ
## Tek USB SSD'de Birden Fazla OS · VHD · rEFInd · Persistence · Güvenlik

> **Amaç:** Tek bir harici USB SSD üzerinde Windows To Go, Ubuntu, Kali ve diğer ISO'ları Ventoy ile bir arada çalıştırmak. Her senaryonun adım adım kurulumu, sorun giderme ve güvenlik analizi.

---

## 📐 BÜYÜK RESİM — Mimari Seçenekleri

```
┌─────────────────────────────────────────────────────────────────┐
│  YOL A: Ventoy + WTG (VHD) + ISO'lar                           │
│  ─────────────────────────────────────────────────────────────  │
│  USB SSD → Ventoy yapısı → WTG.vhd + Ubuntu.iso + Kali.iso    │
│  Tek boot menüsü, kurulum kolay, WTG kalıcı, ISO'lar live      │
│  Güvenlik: ★★★☆☆                                                │
├─────────────────────────────────────────────────────────────────┤
│  YOL B: rEFInd + WTG Bölümü + Ubuntu Kurulu Bölümü             │
│  ─────────────────────────────────────────────────────────────  │
│  USB SSD → Bölümlü yapı → Her OS kendi bölümünde tam kurulu    │
│  En esnek, her OS kalıcı, BitLocker + LUKS2 tam güvenlik       │
│  Güvenlik: ★★★★☆                                                │
├─────────────────────────────────────────────────────────────────┤
│  YOL C: Saf Ventoy — Her Şey ISO + Ubuntu Persistence          │
│  ─────────────────────────────────────────────────────────────  │
│  USB SSD → Ventoy → ISO'lar + persistence dosyaları            │
│  Kurulum en kolay, WTG yok (live Windows), Ubuntu yarı kalıcı  │
│  Güvenlik: ★★☆☆☆                                                │
└─────────────────────────────────────────────────────────────────┘
```

---

## BÖLÜM 0 — VENTOY NEDİR? (Temel Kavramlar)

### 0.1 — Ventoy Nasıl Çalışır?

```
Normal USB Kurulum Medyası:
  ISO → Rufus → USB'ye yaz → Yalnızca o ISO'yu açar
  Yeni ISO için USB'yi SİL, yeniden yaz

Ventoy:
  Bir kez kur → USB'ye istediğin kadar ISO at → Hepsi menüde gösterir
  ISO silmek/eklemek için USB'yi formatlamak gerekmez

Bölüm Yapısı (Ventoy kurulumu sonrası):
  ├─ Bölüm 1: FAT32 (32 MB) — Ventoy bootloader ve meta veriler
  └─ Bölüm 2: exFAT/NTFS (kalan tüm alan) — ISO, VHD, dosyalar buraya
```

### 0.2 — Ventoy'un Desteklediği Dosya Türleri

| Dosya Türü | Örnek | Kaydedilir mi? |
|-----------|-------|----------------|
| `.iso` | ubuntu.iso, kali.iso | ❌ Live (RAM'de) |
| `.img` | disk.img | ❌ Live |
| `.vhd` / `.vhdx` | windows_wtg.vhd | ✅ Kalıcı (sanal disk) |
| `.efi` | bootx64.efi | Doğrudan çalıştır |
| Persistence `.dat` | ubuntu-rw.dat | ✅ ISO değişiklikleri kaydeder |

### 0.3 — Gerekli Donanım

```
✅ USB SSD Kutusu (UASP destekli, USB 3.0+)
✅ M.2 NVMe SSD (en az 256 GB önerilir)
✅ Hedef bilgisayar USB 3.0+ port
✅ Yönetici yetkili Windows bilgisayar (kurulum için)
```

---

## BÖLÜM 1 — YOL A: VENTOY + WTG (VHD) + ISO'LAR

### 1.1 — Genel Mimari

```
USB SSD (256 GB örnek)
├── [Ventoy Bölüm 1 — FAT32 32MB]
└── [Ventoy Bölüm 2 — exFAT/NTFS ~255 GB]
    ├── windows11_wtg.vhd      (80 GB) ← WTG kurulu, kalıcı
    ├── ubuntu-24.04.iso       (5 GB)  ← Live, değişmez
    ├── kali-2024-live.iso     (4 GB)  ← Live, değişmez
    ├── persistence/                   ← Ubuntu değişiklik kayıt dizini
    │   └── ubuntu-24.04-persistence.dat (8 GB) ← Ubuntu kalıcı veriler
    └── tools/                         ← Yardımcı araçlar
        ├── malwarebytes.exe
        └── ventoy.json                ← Yapılandırma
```

### 1.2 — Adım 1: Ventoy'u USB SSD'ye Kur

```
1. https://github.com/ventoy/Ventoy/releases → En son sürümü indir
2. ventoy-x.x.x-windows.zip → Aç → ventoy2disk.exe → Yönetici olarak çalıştır
3. "Device" → USB SSD'ni seç (DİKKAT: Yanlış disk = veri kaybı)
4. "Partition Style" → GPT (UEFI için şart)
5. "Secure Boot Support" → ✅ Etkinleştir (önerilir)
6. "Install" → Onayla → Bekle (~1 dk)
7. Kurulum tamamlandı: USB'de iki bölüm oluştu
```

```powershell
# PowerShell ile Ventoy doğrulama:
Get-Disk | Select-Object Number, FriendlyName, Size, BusType
Get-Partition | Where-Object {$_.DiskNumber -eq X} |  # X = USB disk numarası
    Select-Object PartitionNumber, Size, Type, DriveLetter
# 2 bölüm görünmeli: FAT32 (küçük) + exFAT/NTFS (büyük)
```

### 1.3 — Adım 2: WTG için VHD Oluştur ve Kur

#### VHD Oluşturma (Diskpart)

```
Yönetici Komut İstemi (cmd) veya PowerShell:
```

```cmd
diskpart

:: VHD oluştur
create vdisk file="E:\windows11_wtg.vhd" maximum=81920 type=expandable
:: (81920 MB = 80 GB, expandable = dinamik büyüyen)

:: VHD'yi bağla
select vdisk file="E:\windows11_wtg.vhd"
attach vdisk

:: Bölümleri oluştur (GPT)
convert gpt
create partition efi size=260
format quick fs=fat32 label="System"
assign letter=S

create partition msr size=16

create partition primary
format quick fs=ntfs label="Windows"
assign letter=W

exit
```

#### WTG'yi VHD'ye Kur (WinToUSB)

```
1. WinToUSB → Yönetici olarak aç
2. "ISO veya DVD'den kur" → Windows 11 Pro ISO seç
3. Hedef disk: "W:" harfli VHD diski seç
4. Kurulum modu: "Windows To Go (VHD)"
5. UEFI + GPT seç
6. "Sonraki" → Kurulum ~20-40 dk sürer
7. Bitti → VHD'yi ayır (Disk Yönetimi → Sağ tık → VHD Ayır)
```

#### VHD'yi Ventoy'a Taşı

```powershell
# VHD dosyasını Ventoy'un ikinci bölümüne taşı (E: = Ventoy bölümü)
Move-Item "C:\Kullanici\windows11_wtg.vhd" "E:\windows11_wtg.vhd"

# Kontrol: Dosya boyutu doğru mu?
Get-Item "E:\windows11_wtg.vhd" | Select-Object Name, Length
```

### 1.4 — Adım 3: Linux ISO'larını Ekle

```cmd
:: ISO'ları Ventoy bölümüne kopyala (E: = Ventoy)
copy ubuntu-24.04-desktop-amd64.iso E:\
copy kali-linux-2024.2-live-amd64.iso E:\
:: Bitti. Ventoy menüde otomatik gösterir.
```

### 1.5 — Adım 4: Ubuntu Persistence (Değişiklikler Kaydedilsin)

Persistence olmadan Ubuntu live ISO her açılışta sıfırlanır. Persistence ile kurduğun programlar ve dosyalar kalır.

#### Persistence Dosyasını Oluştur (Linux'ta)

```bash
# Bir Linux sisteminde (veya WSL2) çalıştır:

# 8 GB persistence dosyası oluştur
dd if=/dev/zero of=/mnt/ventoy/persistence/ubuntu-24.04-desktop-persistence.dat \
    bs=1M count=8192 status=progress

# ext4 formatla ve etiketi "casper-rw" yap (Ubuntu persistence etiket adı)
mkfs.ext4 -L casper-rw \
    /mnt/ventoy/persistence/ubuntu-24.04-desktop-persistence.dat
```

#### Persistence Dosyasını Oluştur (Windows'ta — WSL2 ile)

```powershell
# WSL2 açık değilse: wsl --install
wsl

# WSL içinde:
dd if=/dev/zero of=/mnt/e/persistence/ubuntu-24.04-desktop-persistence.dat bs=1M count=8192
mkfs.ext4 -L casper-rw /mnt/e/persistence/ubuntu-24.04-desktop-persistence.dat
exit
```

#### ventoy.json ile Persistence Bağla

```json
// E:\ventoy\ventoy.json dosyasını oluştur:
{
    "persistence": [
        {
            "image": "/ubuntu-24.04-desktop-amd64.iso",
            "backend": "/persistence/ubuntu-24.04-desktop-persistence.dat"
        }
    ],
    "control": [
        {
            "VTOY_DEFAULT_SEARCH_ROOT": "/"
        }
    ]
}
```

### 1.6 — Adım 5: Ventoy Menüsünü Özelleştir

```json
// E:\ventoy\ventoy.json — Gelişmiş yapılandırma:
{
    "theme": {
        "display_mode": "GUI",
        "serial_param": "",
        "fontsz": 24
    },
    "menu_alias": [
        {
            "image": "/windows11_wtg.vhd",
            "alias": "🪟 Windows 11 Pro (WTG — Kalıcı)"
        },
        {
            "image": "/ubuntu-24.04-desktop-amd64.iso",
            "alias": "🐧 Ubuntu 24.04 LTS (Persistence)"
        },
        {
            "image": "/kali-linux-2024.2-live-amd64.iso",
            "alias": "🐉 Kali Linux 2024 (Live)"
        }
    ],
    "persistence": [
        {
            "image": "/ubuntu-24.04-desktop-amd64.iso",
            "backend": "/persistence/ubuntu-24.04-desktop-persistence.dat"
        }
    ],
    "control": [
        {"VTOY_DEFAULT_SEARCH_ROOT": "/"},
        {"VTOY_WIN11_BYPASS_CHECK": "1"},
        {"VTOY_VHD_NO_WARNING": "1"}
    ]
}
```

### 1.7 — Boot Süreci (Yol A)

```
Bilgisayarı aç → BIOS Boot Menu (F12) → USB SSD seç
                                ↓
               Ventoy Grafik Menüsü açılır
    ┌──────────────────────────────────────────────┐
    │  🪟 Windows 11 Pro (WTG — Kalıcı)           │
    │  🐧 Ubuntu 24.04 LTS (Persistence)          │
    │  🐉 Kali Linux 2024 (Live)                  │
    └──────────────────────────────────────────────┘
                         ↓
    Windows seçilince: VHD mount → WTG açılır (kalıcı)
    Ubuntu seçilince:  ISO + persistence → değişiklikler kaydedilir
    Kali seçilince:    Live ISO → RAM'de çalışır, kapanışta sıfırlanır
```

---

## BÖLÜM 2 — YOL B: rEFInd + WTG BÖLÜMÜ + UBUNTU KURULU BÖLÜMÜ

Bu yol en güçlüdür. Her OS kendi bölümünde tam kurulu ve şifreli çalışır.

### 2.1 — Disk Bölüm Planı

```
USB SSD (256 GB örnek):
┌──────────┬──────────────┬──────────────┬──────────────┐
│  EFI     │  Windows 11  │  Ubuntu 24   │  Ortak Veri  │
│  FAT32   │  NTFS        │  ext4+LUKS2  │  exFAT       │
│  500 MB  │  100 GB      │  80 GB       │  Kalan       │
│  rEFInd  │  BitLocker   │  LUKS2 şifreli│             │
└──────────┴──────────────┴──────────────┴──────────────┘
     ↑                                         ↑
  rEFInd her iki OS'u tarar        Her iki OS buraya yazar
  ve menüde gösterir
```

### 2.2 — Disk Hazırlığı (Diskpart)

```cmd
diskpart

list disk
select disk X    :: X = USB SSD numarası — DİKKAT!

clean            :: TÜM VERİ SİLİNİR
convert gpt

:: EFI bölümü
create partition efi size=500
format quick fs=fat32 label="EFI"
assign letter=S

:: Windows (WTG) bölümü
create partition primary size=102400   :: 100 GB
format quick fs=ntfs label="Win11WTG"
assign letter=W

:: Ubuntu bölümü (LUKS sonrası formatlanacak, şimdi sadece alan ayır)
create partition primary size=81920    :: 80 GB
:: Formatlamıyoruz — Ubuntu kurulumunda LUKS2 ile şifrelenecek

:: Ortak veri bölümü
create partition primary
format quick fs=exfat label="OrtakVeri"
assign letter=D

exit
```

### 2.3 — Windows WTG Bölümüne Kur

```powershell
# WinToUSB ile WTG kurulumu:
# Hedef: W: → Windows 11 Pro ISO → Windows To Go → GPT/UEFI
# (BÖLÜM 1'deki WinToUSB adımlarını burada uygula, hedef fiziksel bölüm)
```

### 2.4 — Ubuntu'yu LUKS2 ile Kur

```
1. Ubuntu 24.04 ISO'yu ayrı bir USB belleğe yaz (Rufus ile)
2. Bilgisayarı Ubuntu USB'sinden başlat
3. Kurulum Sihirbazı → "Başka bir şey" (gelişmiş bölümlendirme) seç
4. USB SSD'deki 3. bölümü seç (80 GB'lık)
5. "Bu bölümü şifrele (LUKS)" seçeneğini işaretle
6. Parola gir (güçlü — 20+ karakter)
7. / (root) olarak bağla → Devam
8. Önyükleyici: USB SSD'nin EFI bölümü (S:)
9. Kur → Bekle (~10-20 dk)
```

```bash
# Ubuntu kurulumu sonrası LUKS doğrulama (Terminal):
cryptsetup luksDump /dev/sdX3   # X = USB disk

# LUKS header yedekle (kritik!):
sudo cryptsetup luksHeaderBackup /dev/sdX3 \
    --header-backup-file /media/kullanici/OrtakVeri/luks_ubuntu_header.img
```

### 2.5 — rEFInd Kurulumu

rEFInd, EFI bölümüne kurulur ve hem Windows WTG hem Ubuntu'yu otomatik tanır.

#### Windows Tarafından:

```powershell
# rEFInd indir: https://www.rodsbooks.com/refind/
# refind-bin-x.xx.x.zip → Aç → refind → içindeki refind-install.exe

# EFI bölümünü bağla
mountvol S: /S   # S: EFI bölümü

# rEFInd dosyalarını EFI'ye kopyala
New-Item -ItemType Directory -Force -Path "S:\EFI\refind"
Copy-Item ".\refind-bin-x.xx.x\refind\*" "S:\EFI\refind\" -Recurse

# UEFI'ye rEFInd girdisi ekle
bcdedit /set "{bootmgr}" path \EFI\refind\refind_x64.efi
```

#### Ubuntu Tarafından (Daha Kolay):

```bash
# Ubuntu açıkken:
sudo apt install refind
sudo refind-install
# rEFInd otomatik olarak EFI'ye kurulur ve tüm OS'ları tarar
```

#### rEFInd Yapılandırması:

```bash
# /boot/efi/EFI/refind/refind.conf
sudo nano /boot/efi/EFI/refind/refind.conf
```

```
# refind.conf içeriği (ekle/değiştir):
timeout 10
use_graphics_for osx,linux,windows
showtools shell,shutdown,reboot,memtest,about,hidden_tags,firmware

# Tema (opsiyonel — güzel görünüm)
include themes/refind-theme-regular/theme.conf

# Tarama sırası
scanfor internal,external,optical,manual

# Windows WTG girdisi (eğer otomatik algılamıyorsa):
menuentry "Windows 11 WTG" {
    icon     /EFI/refind/icons/os_win8.png
    volume   "Win11WTG"
    loader   \EFI\Microsoft\Boot\bootmgfw.efi
}

# Ubuntu girdisi (bunu otomatik ekler ama elle de tanımlayabilirsin):
menuentry "Ubuntu 24.04" {
    icon     /EFI/refind/icons/os_linux.png
    volume   "ubuntu"
    loader   /boot/vmlinuz
    initrd   /boot/initrd.img
    options  "root=/dev/mapper/ubuntu--vg-ubuntu--lv ro quiet splash"
}
```

### 2.6 — rEFInd Menüsü Boot Süreci (Yol B)

```
Bilgisayarı aç → BIOS → USB SSD seç → rEFInd açılır
    ┌─────────────────────────────────────────────────┐
    │   [🪟]           [🐧]          [🔧]            │
    │  Windows 11   Ubuntu 24.04    Shell/Tools       │
    │    WTG                                          │
    └─────────────────────────────────────────────────┘
                    ↓
    Windows seçildi: BitLocker parola → WTG açılır
    Ubuntu seçildi:  LUKS parola → Ubuntu açılır
    (Her OS tamamen bağımsız, kalıcı, şifreli)
```

---

## BÖLÜM 3 — YOL C: SAF VENTOY — HER ŞEY ISO + PERSISTENCE

En basit senaryo. WTG yok, sadece live/persistence ISO'lar.

### 3.1 — Kurulum

```
1. Ventoy'u USB SSD'ye kur (Bölüm 1.2)
2. ISO'ları kopyala:
   ubuntu-24.04.iso → E:\
   kali-linux.iso   → E:\
   win11.iso        → E:\ (live — kurulum medyası olarak kullanılır)

3. Ubuntu persistence ekle (Bölüm 1.5)
4. Bitti — USB açık, menü hazır
```

### 3.2 — Kali Persistence

```bash
# Kali için persistence dosyası:
# (Kali persistence etiketi "persistence" — Ubuntu'dan farklı)
dd if=/dev/zero of=/mnt/ventoy/persistence/kali-persistence.dat bs=1M count=4096
mkfs.ext4 -L persistence /mnt/ventoy/persistence/kali-persistence.dat

# Kali persistence için ayrı yapılandırma dosyası:
# persistence.conf root dizinine yaz:
# mount ile dosyayı bağla, içine persistence.conf ekle:
mkdir /tmp/kali_p
mount /mnt/ventoy/persistence/kali-persistence.dat /tmp/kali_p
echo "/ union" | tee /tmp/kali_p/persistence.conf
umount /tmp/kali_p
```

```json
// ventoy.json'a Kali persistence ekle:
{
    "persistence": [
        {
            "image": "/ubuntu-24.04-desktop-amd64.iso",
            "backend": "/persistence/ubuntu-24.04-desktop-persistence.dat"
        },
        {
            "image": "/kali-linux-2024.2-live-amd64.iso",
            "backend": "/persistence/kali-persistence.dat"
        }
    ]
}
```

---

## BÖLÜM 4 — VHD NEDİR? + TEKNİK DERİNLEME

### 4.0 — VHD Nedir? (Sıfırdan Açıklama)

**VHD (Virtual Hard Disk)** — tek bir dosyanın içinde tam bir fiziksel disk gibi davranan sanal disk formatı. Microsoft tarafından geliştirilmiş, Hyper-V ve Windows'un native desteklediği bir standarttır.

**Somut Analoji:**
```
Gerçek Disk:  Fiziksel metal/silikon — bilgisayara takılır, bölümleri var, OS kurarsın
VHD Dosyası:  windows_wtg.vhd → USB SSD'de duran TEK BİR DOSYA
              Ama içi tıpatıp fiziksel disk gibi davranır:
              EFI bölümü var, NTFS bölümü var, Windows kurulu

Bir zip arşivi nasıl içinde onlarca dosya taşıyorsa,
VHD de içinde tam bir işletim sistemi taşır — ama açıldığında
gerçek bir disk gibi görünür ve çalışır.
```

**Bağlanınca ne olur?**
```
windows_wtg.vhd dosyası (USB SSD'de duruyor)
        ↓ Windows veya Ventoy bu dosyayı "bağlar" (mount)
        ↓
Sanal Disk beliriyor → Disk Yönetimi'nde C: veya D: olarak
        ↓
İçindeki Windows sanki gerçek diske kurulu gibi çalışıyor
        ↓
Kapatınca → sanal disk ayrılır → yeniden tek .vhd dosyasına döner
```

**Neden WTG için VHD kullanıyoruz?**
```
Ventoy, ISO dosyalarını açabilir → ama ISO read-only, kalıcı değil
Ventoy, VHD dosyalarını da açabilir → ama VHD read-write, içindeki Windows kalıcı
Böylece: Ventoy menüsünden "kurulu, kalıcı" bir Windows başlatabiliyorsun
ve aynı USB'de Ubuntu ISO, Kali ISO ile yan yana durabiliyor
```

**Sabit (Fixed) vs Dinamik (Dynamic) VHD:**
```
Sabit VHD:    80 GB diyorsun → dosya hemen 80 GB yer kaplar (boş bile olsa)
              Avantaj: Daha hızlı (bloklar önceden ayrılmış)
              Dezavantaj: Yer israfı

Dinamik VHD:  80 GB diyorsun → dosya başta küçük, doldukça büyür
              Avantaj: Yer tasarrufu (30 GB dolduysa dosya ~30 GB)
              Dezavantaj: Az performans kaybı (blok tahsisi runtime'da)

Önerim: Dinamik başla, performans sorunun olursa sabite geç
```

### 4.1 — VHD vs VHDX Farkı

| Özellik | VHD | VHDX |
|---------|-----|------|
| Max boyut | 2 TB | 64 TB |
| Block boyutu | 512 byte | 1-256 MB (ayarlanabilir) |
| Güç kesintisi koruması | ❌ | ✅ |
| SSD trim desteği | ❌ | ✅ |
| Ventoy uyumu | ✅ Tam | ✅ Tam |
| Performans | Normal | Daha iyi |
| **Önerim** | Uyumluluk için | **Performans için** |

```powershell
# VHDX oluştur (önerilir):
New-VHD -Path "E:\windows11_wtg.vhdx" -SizeBytes 80GB -Dynamic -BlockSizeBytes 1MB
```

### 4.2 — VHD Boot Zinciri (Ventoy'da VHD nasıl açılır?)

```
Ventoy bootloader başlar
    ↓
VHD/VHDX dosyasını algılar
    ↓
FIMDISK veya libVHD ile sanal disk belleğe bağlanır
    ↓
Sanal diskin EFI bölümü taranır: \EFI\Microsoft\Boot\bootmgfw.efi bulunur
    ↓
Boot yöneticisi kontrolü VHD içindeki Windows'a verir
    ↓
Windows, "ben sanal diskte çalışıyorum" görür → Native Boot VHD modu
    ↓
WTG açılır (Host disk ve VHD aynı anda erişilebilir)
```

### 4.3 — VHD Boyutu Yönetimi

```powershell
# VHD ne kadar doluyor?
$vhd = "E:\windows11_wtg.vhdx"
$size = (Get-Item $vhd).Length / 1GB
Write-Host "VHD dosya boyutu: $([math]::Round($size, 2)) GB"

# VHD'yi bağlayıp içini gör:
Mount-VHD -Path $vhd -ReadOnly
Get-Disk | Select-Object Number, FriendlyName, Size   # yeni disk görünür
Dismount-VHD -Path $vhd

# Dinamik VHD'yi sıkıştır (kullanılmayan alanı geri al):
Optimize-VHD -Path $vhd -Mode Full
```

---

## BÖLÜM 5 — GÜVENLİK: YOL KARŞILAŞTIRMASI

### 5.1 — Her Yolun Güvenlik Profili

#### Yol A (Ventoy + VHD):

```
Avantaj:
  ✅ WTG sanal disk içinde → Host OS dosya sistemi ayrı
  ✅ VHD BitLocker ile şifrelenebilir
  ✅ ISO'lar sadece okunabilir — değiştirilemiyor

Dezavantaj:
  ❌ Ventoy bölümü (exFAT) şifrelenmez → ISO'lar görünebilir
  ❌ VHD şifrelenmediyse direk kopyalanıp başka cihazda açılır
  ❌ Ubuntu live → persistence dosyası şifrelenmez (ext4 düz)
```

**VHD'yi BitLocker ile Şifrele:**
```powershell
# VHD'yi bağla:
Mount-VHD -Path "E:\windows11_wtg.vhdx"

# Bağlanan diski bul:
$disk = Get-Disk | Where-Object {$_.FriendlyName -like "*virtual*"} |
        Select-Object -Last 1
$letter = (Get-Partition -DiskNumber $disk.Number |
           Where-Object {$_.DriveLetter}).DriveLetter

# BitLocker uygula (parola tabanlı — TPM yok çünkü VHD):
Enable-BitLocker -MountPoint "$letter`:" `
    -EncryptionMethod XtsAes256 `
    -PasswordProtector `
    -Password (Read-Host -AsSecureString "VHD BitLocker Parolası") `
    -UsedSpaceOnly:$false

# Dismount — artık VHD şifreli:
Dismount-VHD -Path "E:\windows11_wtg.vhdx"
```

**Ubuntu Persistence Dosyasını Şifrele (LUKS):**
```bash
# Persistence dosyasını LUKS ile şifrele (henüz format yapmadan önce):
sudo cryptsetup luksFormat --type luks2 \
    /mnt/ventoy/persistence/ubuntu-persistence.dat

sudo cryptsetup open \
    /mnt/ventoy/persistence/ubuntu-persistence.dat ubuntu_persist

sudo mkfs.ext4 -L casper-rw /dev/mapper/ubuntu_persist
sudo cryptsetup close ubuntu_persist
```

> [!WARNING]
> Şifreli persistence dosyasını ventoy.json'a tanıtmak karmaşıktır. Gelişmiş kullanıcı için uygundur; basit kurulum istiyorsan şifresiz bırak, kritik veri depolama.

#### Yol B (rEFInd + Bölümler):

```
Avantaj:
  ✅ Windows bölümü → BitLocker XTS-AES 256 (parola tabanlı)
  ✅ Ubuntu bölümü → LUKS2 AES-256-XTS
  ✅ Her OS tamamen yalıtılmış bölümde
  ✅ rEFInd UEFI Secure Boot destekler
  ✅ Ortak veri bölümü de şifrelenebilir (VeraCrypt)

Dezavantaj:
  ❌ Kurulum en karmaşık
  ❌ Bölüm boyutlarını sonradan değiştirmek zor
  ❌ EFI bölümü şifrelenmez (her zaman böyle — standart)
```

**Ortak Veri Bölümünü VeraCrypt ile Şifrele:**
```
VeraCrypt → Volume oluştur → Cihaz tabanlı → exFAT bölümünü seç
          → AES-256 + SHA-512 → Parola → Biçimlendir
Her iki OS'tan (Windows ve Ubuntu) VeraCrypt ile aç
```

#### Yol C (Saf Ventoy):

```
Avantaj:
  ✅ En kolay kurulum
  ✅ ISO'lar değiştirilemez (salt okunur mount)

Dezavantaj:
  ❌ Persistence şifrelenemez (Ventoy kısıtı)
  ❌ WTG yok → Windows her açılışta temiz (kayıp olmaz ama kalıcı değil)
  ❌ Veri yalıtımı en zayıf
```

### 5.2 — Çapraz OS Erişim Riski

```
YOL A (Ventoy+VHD):
  Ubuntu Live açıkken → E:\windows11_wtg.vhdx dosyasına erişebilir
  VHD BitLocker şifrelenmediyse → dosyayı kopyalayabilir
  ÇÖZÜM: VHD'yi BitLocker ile şifrele (yukarıdaki adım)

YOL B (rEFInd+Bölümler):
  Ubuntu açıkken → lsblk → Windows WTG NTFS bölümü görünür
  lsblk → /dev/sdX2   Windows NTFS bölümü
  sudo mount /dev/sdX2 /mnt → BitLocker aktifse "Unknown FS" → Erişilmez ✅
  BitLocker yoksa → sam, system, security çekilebilir → RİSK

  Windows açıkken → Ubuntu LUKS2 bölümü görünür
  Disk Yönetimi → "Biçimlendirilmemiş" → Açılamaz ✅ (LUKS anlaşılmaz)
```

---

## BÖLÜM 6 — VENTOY İPUCU VE GELİŞMİŞ ÖZELLİKLER

### 6.1 — Ventoy Secure Boot Uyumu

```
Ventoy'un Secure Boot zinciri:
  BIOS → shim.efi (Microsoft imzalı) → grubx64.efi (MOK imzalı) → Ventoy

Secure Boot ile Ventoy kullanmak için:
  1. Ventoy kurulumunda "Secure Boot Support" → Etkin
  2. İlk açılışta "Enroll Key" ekranı gelir
  3. "Enroll Key from Disk" → ENROLL_THIS_KEY_IN_MOKMANAGER.cer seç
  4. Onayla → Yeniden başlat → Artık Secure Boot ile açılır

Sorun çıkarsa:
  BIOS → Secure Boot → "MokManager" → Ventoy sertifikasını manuel ekle
```

### 6.2 — Ventoy Temalar ve Özelleştirme

```bash
# Tema indirme: https://github.com/ventoy/ventoy-plugin
# E:\ventoy\ klasörüne temayı kopyala

# ventoy.json'a ekle:
{
    "theme": {
        "file": "/ventoy/themes/ventoy-theme/theme.txt",
        "display_mode": "GUI",
        "fontsz": 24,
        "gfxmode": "1920x1080"
    }
}
```

### 6.3 — Ventoy Gizli ISO Özelliği

```json
// Bazı ISO'ları menüden gizle (parola ile aç):
{
    "control": [
        {"VTOY_MENU_TIMEOUT": "10"},
        {"VTOY_DEFAULT_MENU_MODE": "0"},
        {"VTOY_FILE_FLT": "iso|img|vhd|vhdx"}
    ],
    "password": [
        {
            "image": "/gizli_tools.iso",
            "passwd": "sifre123"
        }
    ]
}
```

### 6.4 — Ventoy Güncelleme (ISO'ları Silmeden)

```
Ventoy yeni sürüm çıktında MEVCUT ISO VE VERİLERİ SİLMEDEN güncelleme:

ventoy2disk.exe → "Upgrade" sekmesi → Cihazı seç → Upgrade
Yalnızca Ventoy bootloader güncellenir, veriler dokunulmaz
```

---

## BÖLÜM 7 — SORUN GİDERME

### 7.1 — "VHD açılmıyor" Hataları

| Hata | Neden | Çözüm |
|------|-------|-------|
| `Error: Can't mount VHD` | VHD bozuk veya çok büyük | VHD'yi defrag et, NTFS'yi kontrol et |
| `0xc000000e` | VHD'deki WTG USB sürücüsünü tanımıyor | WinToUSB ile WTG modunda yeniden kur |
| `0xc0000225` | BCD bozuk | VHD'yi bağla → `bcdboot W:\Windows /s S: /f UEFI` |
| "Disk not found" | Ventoy VHD'yi bulamıyor | İsimde Türkçe karakter olmamalı (`wtg.vhd`, `windows.vhd`) |
| Siyah ekran | UEFI/Legacy karışıklığı | BIOS → CSM → Disabled |

### 7.2 — "Ubuntu persistence çalışmıyor"

```bash
# Sorun: Persistence dosyası bağlanmıyor
# Kontrol:
sudo dmesg | grep casper
sudo blkid | grep casper-rw
# "casper-rw" etiketi görünmeli

# Etiket doğru mu?
sudo e2label /dev/sdXY   # casper-rw çıkmalı

# Yanlışsa düzelt:
sudo e2label /dev/sdXY casper-rw
```

### 7.3 — "rEFInd Windows'u görmüyor"

```bash
# Windows WTG EFI girişini kontrol et:
sudo efibootmgr -v | grep -i windows

# rEFInd'in Windows'u taraması için:
sudo nano /boot/efi/EFI/refind/refind.conf
# Şu satırı bul ve düzenle:
# scanfor internal,external,optical,manual
# external ekli olmalı (USB diskler için)

# Manuel Windows girdisi ekle:
menuentry "Windows 11 WTG" {
    icon /EFI/refind/icons/os_win8.png
    volume "Win11WTG"
    loader \EFI\Microsoft\Boot\bootmgfw.efi
    disabled     # Bu satırı sil
}
```

### 7.4 — "Ventoy menüsü gelmiyor, direkt diske gidiyor"

```
BIOS → Boot Priority → USB SSD'yi en üste al
KESİNLİKLE: İç diski değil, USB SSD'yi seç

Boot Order:
  1. USB SSD (Ventoy)        ← En üst
  2. İç NVMe                 ← İkinci
  3. Diğerleri               ← Altta

Ventoy menüsü aşılıp direkt iç diske gidiyorsa:
  → USB SSD düzgün algılanmıyor → Farklı USB portuna tak
  → BIOS'ta "USB Boot" seçeneği kapalı olabilir → Aç
```

---

## BÖLÜM 8 — SSD vs HDD: MULTIBOOT İÇİN FARKLAR

### 8.1 — Bu Sistem HDD'de Birebir Aynı Çalışır mı?

Kısa cevap: **Teknik olarak evet, ama pratikte ciddi farklar var.**

```
SSD (NVMe/SATA):
  ├─ Ventoy boot süresi: 3-5 saniye
  ├─ VHD mount süresi: 2-4 saniye
  ├─ WTG Windows açılış: 15-30 saniye
  ├─ TRIM desteği: Evet → SSD ömrü korunur
  ├─ Parçalanma sorunu: Yok (SSD'de anlamlı değil)
  └─ Titreşim/darbe dayanımı: Yüksek (hareketli parça yok)

HDD (5400/7200 RPM):
  ├─ Ventoy boot süresi: 8-15 saniye
  ├─ VHD mount süresi: 10-30 saniye ← KRİTİK YAVAŞLAMA
  ├─ WTG Windows açılış: 60-120+ saniye
  ├─ TRIM desteği: Yok
  ├─ Parçalanma sorunu: VAR → VHD parçalanırsa çok yavaşlar
  └─ Titreşim/darbe dayanımı: Düşük (taşınabilir kullanımda risk)
```

### 8.2 — VHD Parçalanma Sorunu (Yalnızca HDD)

> [!WARNING]
> Dinamik VHD dosyası HDD'de zamanla **parçalanır** (fragmentation). Bu, VHD'nin farklı fiziksel sektörlere dağılması demektir. SSD'de bu sorun yoktur çünkü erişim süresi konumdan bağımsızdır.

```powershell
# HDD'de VHD parçalanmasını kontrol et:
defrag E: /a /v   # /a = analiz, /v = detaylı rapor
# "Dosya parçalanması" yüzdesini oku — %10 üzeriyse defrag et

# Sadece VHD dosyasını defrag et (tüm diski değil):
defrag E: /o      # /o = SSD ve HDD'ye göre otomatik en iyileştirme

# Alternatif — VHD'yi Contig ile defrag et (Sysinternals):
# contig.exe E:\windows11_wtg.vhd
```

### 8.3 — HDD İçin Özel Yapılandırma Önerileri

```
1. VHD Türü: SABİT (Fixed) kullan — dinamik değil
   → Sabit VHD'nin blokları bitişik olur → parçalanma olmaz
   → Dezavantaj: 80 GB VHD hemen 80 GB yer kaplar

2. Ventoy Dosya Sistemi: NTFS kullan (exFAT yerine)
   → NTFS defrag yapılabilir, exFAT yapılamaz
   → Ventoy kurulumunda "Partition Style: GPT" + "FS: NTFS" seç

3. Persistence Dosyası: Küçük tut (4 GB max)
   → HDD'de büyük persistence dosyası → rastgele yazma → çok yavaş

4. ISO Dosyaları: Diskin başına koy
   → Kopyalama sırasında disk boşsa → dosyalar bitişik yerleşir
   → ISO → VHD → Persistence sırasıyla kopyala (bu sıra önemli)
```

```powershell
# HDD'de SABİT VHD oluştur:
New-VHD -Path "E:\windows11_wtg.vhd" -SizeBytes 80GB -Fixed
# "Fixed" = sabit boyut → parçalanma yok → HDD'de çok daha hızlı
```

### 8.4 — SSD İçin Özel Yapılandırma Önerileri

```
1. VHD Türü: DİNAMİK (Dynamic) kullan
   → SSD'de parçalanma sorunu yok → dinamik yer tasarrufu sağlar
   → VHDX tercih et (TRIM desteği, güç kesintisi koruması)

2. Ventoy Dosya Sistemi: exFAT veya NTFS — ikisi de iyi
   → exFAT: Linux ve Windows uyumlu, basit
   → NTFS: Daha güçlü (ACL, journal), ama Linux'ta tam yazma için ntfs-3g gerekir

3. TRIM Aktif Olmalı:
   → Windows'ta: fsutil behavior query DisableDeleteNotify → 0 olmalı
   → USB SSD kutusu UASP destekliyse TRIM USB üzerinden geçer
   → Desteklemiyorsa (BOT modu): TRIM çalışmaz → SSD zamanla yavaşlar

4. Wear Leveling Bilinci:
   → Sık VHD yazma-okuma SSD ömrünü etkiler (ama modern SSD'ler 600+ TBW)
   → Persistence boyutunu makul tut (8-16 GB)
```

```powershell
# SSD'de DİNAMİK VHDX oluştur:
New-VHD -Path "E:\windows11_wtg.vhdx" -SizeBytes 80GB -Dynamic -BlockSizeBytes 1MB
# BlockSizeBytes 1MB = Daha granüler büyüme, daha az israf

# SSD TRIM destekleniyor mu kontrol et:
fsutil behavior query DisableDeleteNotify
# 0 = TRIM aktif ✅

# USB kutusu UASP kullanıyor mu?
Get-PnpDevice -Class "SCSIAdapter" |
    Where-Object {$_.FriendlyName -like "*UAS*"}
# Sonuç varsa = UASP aktif → TRIM USB'den geçer
```

### 8.5 — SSD vs HDD Multiboot Özet Tablosu

| Yapılandırma | HDD | SSD |
|-------------|-----|-----|
| **Ventoy Dosya Sistemi** | NTFS (defrag için) | exFAT veya NTFS |
| **VHD Türü** | ✅ Sabit (Fixed) | ✅ Dinamik (Dynamic) |
| **VHD Formatı** | VHD (uyumluluk) | VHDX (TRIM + güç koruması) |
| **Parçalanma Riski** | ⚠️ Yüksek (defrag şart) | ❌ Yok |
| **Boot Süresi** | 60-120+ sn | 15-30 sn |
| **Persistence Boyutu** | Max 4 GB (yavaşlık) | 8-16 GB sorunsuz |
| **Taşınabilirlik** | ❌ Darbe riski (plaka hasarı) | ✅ Dayanıklı |
| **TRIM** | ❌ Yok | ✅ (UASP kutu gerekir) |
| **Güvenli Silme** | DBAN / 7 geçiş | ATA Secure Erase |
| **BitLocker Performansı** | %15-25 yavaşlama | <%1 |
| **Genel Önerim** | ⚠️ Çalışır ama yavaş | ✅ Önerilen |

> [!IMPORTANT]
> **Harici çıkarılabilir disk olarak HDD kullanacaksan:** Mutlaka darbeye dayanıklı bir kutu (Silicon Power, LaCie Rugged vb.) kullan. HDD'nin manyetik plakaları fiziksel darbe ile kalıcı hasar görür — taşınabilir kullanımda SSD kesinlikle tercih edilmeli.

---

## BÖLÜM 9 — NİHAİ KARAR MATRİSİ

### Senaryona Göre Yol Seçimi

| Senaryo | Önerilen Yol |
|---------|-------------|
| "Hızlı başlayayım, WTG + birkaç ISO yeterli" | **Yol A** |
| "Her şey tam kurulu, kalıcı, güvenli olsun" | **Yol B** |
| "Sadece live ISO'larım olsun, basit kalsın" | **Yol C** |
| "Maksimum güvenlik, karmaşıklık sorun değil" | **Yol B + LUKS2 + BitLocker** |
| "En esnek, her şeyi denemek istiyorum" | **Yol A** başla, sonra **Yol B**'ye geç |

### Güvenlik Skoru Karşılaştırması

| Özellik | Yol A | Yol B | Yol C |
|---------|-------|-------|-------|
| Windows şifreleme | ⚠️ VHD BitLocker (opsiyonel) | ✅ Tam BitLocker | ❌ Yok |
| Linux şifreleme | ❌ Persistence şifresiz | ✅ LUKS2 | ❌ Yok |
| Ortak bölüm | ❌ exFAT açık | ⚠️ VeraCrypt ile | ❌ exFAT açık |
| Kurulum kolaylığı | ★★★★☆ | ★★☆☆☆ | ★★★★★ |
| Esneklik | ★★★★☆ | ★★★☆☆ | ★★★★☆ |
| Güvenlik | ★★★☆☆ | ★★★★☆ | ★★☆☆☆ |

### Gerekli Araçlar Özeti

| Araç | Amaç | İndirme |
|------|------|---------|
| **Ventoy** | Çoklu ISO boot | github.com/ventoy/Ventoy |
| **WinToUSB** | WTG kurulumu | easyuefi.com/wintousb |
| **Rufus** | WTG VHD alternatifi | rufus.ie |
| **rEFInd** | EFI boot yöneticisi | rodsbooks.com/refind |
| **VeraCrypt** | Ortak bölüm şifreleme | veracrypt.fr |
| **GParted** (Live) | Disk bölümlendirme | gparted.org |

---

## KONTROL LİSTESİ

```
YOL A — VENTOY + VHD:
  □ Ventoy kuruldu (GPT + Secure Boot)
  □ VHD oluşturuldu (VHDX önerilir)
  □ WinToUSB ile WTG VHD'ye kuruldu
  □ VHD Ventoy bölümüne taşındı
  □ Linux ISO'ları kopyalandı
  □ Ubuntu persistence dosyası oluşturuldu (ext4, casper-rw)
  □ ventoy.json yapılandırıldı
  □ Fast Startup KAPATILDI
  □ VHD BitLocker ile şifrelendi (opsiyonel ama önerilir)
  □ Boot testi yapıldı: WTG, Ubuntu, Kali

YOL B — rEFInd + BÖLÜMLER:
  □ Diskpart ile bölümler oluşturuldu
  □ WinToUSB ile WTG NTFS bölümüne kuruldu
  □ Ubuntu LUKS2 ile kuruldu (ayrı ISO USB'den)
  □ LUKS header yedeklendi (kasaya)
  □ rEFInd EFI bölümüne kuruldu
  □ refind.conf düzenlendi (her iki OS tanımlı)
  □ WTG bölümü BitLocker (parola) ile şifrelendi
  □ Ortak exFAT bölümü VeraCrypt ile şifrelendi
  □ Fast Startup KAPATILDI
  □ Boot testi yapıldı: rEFInd menüsü → her iki OS

YOL C — SAF VENTOY:
  □ Ventoy kuruldu
  □ ISO'lar kopyalandı
  □ Ubuntu + Kali persistence dosyaları oluşturuldu
  □ ventoy.json yapılandırıldı
  □ Boot testi yapıldı
```

---

*Son güncelleme: 2026-03-22 · Ventoy 1.0.x · Windows 11 Pro · Ubuntu 24.04 · Kali 2024.x için hazırlanmıştır.*
