# 🔐 GÜVENLİ KURULUM KRONOLOJİSİ
## Şifreler, Katmanlar ve Nerede Saklandıkları — Tam Dokümantasyon

> **Amaç:** Sıfırdan kurulan bir sistemde hangi şifrenin ne zaman oluşturulduğunu, nerede saklandığını, neyi koruduğunu ve SSD başka bir bilgisayara takıldığında ne işe yarayacağını soru işareti bırakmadan açıklamak.

---

## 🗺️ BÜYÜK RESİM — Şifre Haritası

```
DONANIM KATMANI (Anakart / TPM Çipi)
│
├── [1] BIOS Supervisor/Admin Şifresi     → Anakart kalıcı hafızası (NVRAM)
│       "BIOS ayarlarını ki̇tleyeti"
│
├── [2] BIOS User/Boot Şifresi            → Anakart kalıcı hafızası (NVRAM)
│       "Açılışı ki̇tler"
│
├── [3] BitLocker Başlangıç PIN'i         → TPM Çipi + SSD şifreleme zinciri
│       "Diski̇ çözer — Windows logosundan ÖNCE"
│
└── TPM Çipi                              → Anakart üzerinde fiziksel çip
        BitLocker anahtarını mühürler
        Windows Hello kimlik verisini saklar

─────────────────────────────────────────────────────────────

YAZILIM KATMANI (SSD / İşletim Sistemi)
│
├── [4] Windows Kullanıcı Şifresi         → SSD (Windows\System32\config\SAM)
│       "Masaüstüne gi̇ri̇ş"
│
├── [5] Windows Hello PIN / Biyometri     → TPM Çipi + SSD (ngc klasörü)
│       "Ki̇li̇t ekranını hızlı geç"
│
└── [6] 48 Haneli Kurtarma Anahtarı      → YOK — Sadece SENDE (kağıt/USB)
        "SSD hayat si̇gortası"

─────────────────────────────────────────────────────────────

HARİCİ DEPOLAMA KATMANI (USB / Harici Disk)
│
├── [7a] BitLocker To Go                  → USB'nin kendi depolama alanı
│        Windows native USB şifreleme
│
├── [7b] VeraCrypt Kapsayıcı              → USB içinde şifreli dosya/bölüm
│        Platformdan bağımsız, açık kaynak
│
└── [7c] Donanım Şifreli USB             → USB içindeki güvenlik denetleyicisi
         (Kingston IronKey vb.) Yazılım bağımsız
```

---

## ⏱️ KURULUM ZAMAN TÜNELİ

### AŞAMA 0 — Hazırlık (Bilgisayar henüz kapalı)

Kuruluma başlamadan önce şunlara ihtiyacın var:

| Materyal | Neden? |
|---------|--------|
| Windows 11 Pro ISO (Microsoft'tan) | Temiz, doğrulanmış kaynak |
| Rufus veya Ventoy ile hazırlanmış USB | Kurulum medyası |
| Boş bir kağıt veya fiziksel not defteri | 48 haneli anahtarı yazmak için |
| İnternet bağlantısı **OLMAYACAK** | Microsoft hesabını kurulumda atlatmak için |

---

### AŞAMA 1 — Donanım Kilidi
#### 📅 Ne Zaman: Windows kurulmadan önce, ilk açılışta

BIOS/UEFI, işletim sisteminden tamamen bağımsızdır. SSD'yi söksen, formatla­sa, yeni OS kur­san bile BIOS ayarları anakart üzerindeki **NVRAM** (Non-Volatile RAM) çipinde değişmeden kalır.

```
Bilgisayarı aç → POST ekranında F2 / Del / F10 / Esc (üreticiye göre)
                                     ↓
                              BIOS/UEFI Ekranı
```

#### 1-A: Supervisor (Admin) Şifresi

**Nedir?** BIOS menüsüne girişi kilitleyen şifre.

**Ne yapar?**
- Başka biri F2'ye basarsa bu şifre sorulur
- Şifresiz menüye giremez → Secure Boot'u kapata­maz, boot order'ı değiştiremez, TPM'i devre dışı bıra­kamaz

**Nerede saklanır?**
```
Anakart üzerindeki NVRAM / CMOS çipi
SSD'den bağımsız → SSD söküldüğünde bu şifre yerinde kalır
```

**Nasıl ayarlanır?**
```
BIOS → Security → Set Supervisor Password → Güçlü şifre gir → Kaydet (F10)
```

> [!CAUTION]
> Bu şifreyi unutursan **anakarta fiziksel müdahale** gerekir: CMOS pili çıkartmak, üretici servisine götürmek veya bazı anakartlarda jumper resetlemek. Dijital değil, fiziksel not defterine yaz.

#### 1-B: User/Boot Şifresi

**Nedir?** Bilgisayar her açıldığında, herhangi bir şey yüklenmeden önce sorulan şifre.

**Ne yapar?**
- "Hoş Geldiniz" veya Windows logosu GELMEDEN önce ekran siyah ve şifre kutusu çıkar
- Hatalı girişte sistem boot edemez, donmuş gibi durur
- USB'den bile açılmak için bu şifre aşılmalıdır (veya Supervisor şifresiyle)

**Nerede saklanır?**
```
Yine anakart NVRAM — SSD'den tamamen bağımsız
```

**Nasıl ayarlanır?**
```
BIOS → Security → Set User Password → Şifre gir → Kaydet
```

> [!NOTE]
> Bazı anakartlar Supervisor ve User şifresini birleştirmiş gösterir; bazıları ayrı ayrı sunar. Her ikisini de farklı şifrelerle kur.

#### 1-C: Diğer Kritik BIOS Ayarları (Şifre Değil Ama Yapılacak)

```
Secure Boot     → Enabled (Standard mod)
TPM             → Enabled (Firmware TPM / PTT / fTPM)
VT-d / IOMMU   → Enabled (Kernel DMA Koruması için)
VT-x / AMD-V   → Enabled (HVCI ve VBS için)
Boot Order      → Sadece iç disk; USB boot Disabled veya en alta al
```

---

### AŞAMA 2 — İşletim Sistemi Kurulumu
#### 📅 Ne Zaman: USB'yi taktıktan, boot ettikten sonra

#### 2-A: İnternetsiz Kurulum (Kritik)

```
Windows Kurulum Sihirbazı → Ağ bağlantısı ekranı
                                   ↓
                    "İnternet bağlantım yok" seçeneğini tıkla
                    (Yoksa "Sınırlı kurulum seçenekleri" bağlantısı)
                                   ↓
                         Yerel Hesap Oluştur
```

**Neden internetsiz?**
- Microsoft hesabıyla kurulum yapılırsa BitLocker kurtarma anahtarı **otomatik olarak Microsoft sunucularına yüklenir**
- Biz anahtarın yalnızca bizde olmasını istiyoruz

#### 2-B: Kullanıcı Adı ve Şifre Oluşturma

**Kullanıcı şifresi nedir?**
Masaüstüne geçmek için kullandığın standart Windows şifresi.

**Nerede saklanır?**
```
C:\Windows\System32\config\SAM (Security Account Manager)
Bu dosya hash formatında saklanır — düz metin değil
Sistem çalışırken bu dosyayı kopyalayamazsın (kilitli)
Ama sistem kapalıyken ve disk şifrelenmemişse çalınabilir
→ Bu yüzden BitLocker şart!
```

**Şifre hash formatı (NTLM):**
```
Yazdığın şifre: "MerhAbA123!"
SAM dosyasındaki hali: 8846f7eaee8fb117ad06bdd830b7586c
(Tersine çevrilemez — ama hash tablosuyla kırılabilir)
→ BitLocker + LSA PPL bu riski ortadan kaldırır
```

> [!TIP]
> Kurulumda **iki hesap** oluştur stratejisi:
> 1. **Admin hesabı** (kurulumda oluşturulan ilk hesap) — yalnızca sistem değişiklikleri için
> 2. Kurulum bittikten sonra ayrıca **Standart Kullanıcı hesabı** aç — günlük kullanım için

---

### AŞAMA 3 — Kaleyi Kilitleme: BitLocker
#### 📅 Ne Zaman: Masaüstüne ilk girişin hemen ardından, `gpedit.msc` ayarlarından SONRA

> [!IMPORTANT]
> Önce **Grup İlkesi ayarını** yapmazsan BitLocker seni sadece TPM ile (PIN olmadan) şifrelemeye yönlendirir. Sıralama kritik.

#### 3-A: Önce Grup İlkesi (gpedit.msc)

```
Win + R → gpedit.msc

Bilgisayar Yapılandırması
  └─ Yönetim Şablonları
       └─ Windows Bileşenleri
            └─ BitLocker Sürücü Şifrelemesi
                 └─ İşletim Sistemi Sürücüleri
                      └─ "Başlangıçta ek kimlik doğrulaması iste"
                           → Etkin
                           → TPM ile başlangıç PIN'ini zorunlu kıl ✓
```

```powershell
# Home sürümü için (gpedit.msc yok) Registry alternatifi:
$Path = "HKLM:\SOFTWARE\Policies\Microsoft\FVE"
New-Item -Path $Path -Force | Out-Null
Set-ItemProperty -Path $Path -Name "UseAdvancedStartup" -Value 1
Set-ItemProperty -Path $Path -Name "UseTPMPIN"          -Value 1  # TPM+PIN zorunlu
Set-ItemProperty -Path $Path -Name "OSEncryptionType"   -Value 1  # Tam şifreleme
```

#### 3-B: BitLocker PIN Nedir ve Nasıl Çalışır?

```
Bilgisayar Güç Düğmesi
        ↓
BIOS POST (User/Boot şifresi burada sorulabilir)
        ↓
TPM Ölçümleri (PCR değerleri doğrulanır)
        ↓
┌─────────────────────────────────────────────────┐
│         BİTLOCKER PRE-BOOT PIN EKRANI           │
│                                                  │
│   Enter PIN: ████████                           │
│                                                  │
│   (Windows logosu henüz gelmedi)                │
│   (Disk hâlâ şifreli — PIN olmadan açılmıyor)  │
└─────────────────────────────────────────────────┘
        ↓ PIN doğru girildi
TPM, PIN ile birleşerek şifre çözme anahtarını serbest bırakır
        ↓
SSD şifresi çözülür → Windows yüklenir → Windows Giriş Ekranı
        ↓
Windows Kullanıcı Şifresi / Windows Hello PIN
        ↓
Masaüstü
```

**BitLocker PIN nerede saklanır?**

```
PIN'in kendisi hiçbir yerde düz metin saklanmaz.

Şu şekilde çalışır:
  1. Senin girdiğin PIN → Hash fonksiyonundan geçer
  2. Bu hash → TPM içindeki Volume Master Key (VMK) ile birleşir
  3. VMK → SSD üzerindeki Full Volume Encryption Key (FVEK)'i çözer
  4. FVEK → SSD'deki veriyi gerçek zamanlı şifreler/çözer

Anahtarlar zinciri:
PIN (senin aklında) + TPM (anakart) = VMK (SSD metadata) → FVEK (SSD şifreleyici)

PIN'i unutursan: 48 haneli kurtarma anahtarı bu zinciri atlayarak
doğrudan VMK'ya erişir.
```

#### 3-C: BitLocker'ı Etkinleştirme

```powershell
# Yönetici PowerShell — gpedit.msc ayarları uygulandıktan SONRA

# 1. TPM durumunu doğrula
Get-Tpm
# TpmPresent: True
# TpmEnabled: True

# 2. BitLocker'ı TPM+PIN ile başlat
$pin = Read-Host -AsSecureString "BitLocker Pre-Boot PIN oluşturun (min. 8 karakter)"
Enable-BitLocker -MountPoint "C:" `
    -EncryptionMethod XtsAes256 `
    -TPMandPinProtector `
    -Pin $pin `
    -UsedSpaceOnly:$false

# 3. Kurtarma anahtarını YAZDIR ve KASAYA KALDIR
$rv = (Get-BitLockerVolume -MountPoint "C:").KeyProtector |
      Where-Object {$_.KeyProtectorType -eq "RecoveryPassword"}
Write-Host "48 HANELİ KURTARMA ANAHTARIN:" -ForegroundColor Yellow
Write-Host $rv.RecoveryPassword -ForegroundColor Red
# Bu anahtarı dijital bir yerde SAKLA --- yalnızca fiziksel (kağıt/laminat/kasa)

# 4. İkinci bir USB yedek anahtarı ekle (opsiyonel ama güçlü)
# Add-BitLockerKeyProtector -MountPoint "C:" -StartupKeyProtector -StartupKeyPath "E:\"
```

#### 3-D: 48 Haneli Kurtarma Anahtarı

**Nedir?**
BitLocker'ın otomatik oluşturduğu, 48 rakamdan oluşan acil durum anahtarı.

**Nasıl görünür?**
```
123456-789012-345678-901234-567890-123456-789012-345678
(6'lı gruplar halinde 8 grup = 48 rakam)
```

**Nerede saklanır?**
```
SENDE olmalı — başka hiçbir yerde:
  ✅ Lamine edilmiş kağıt → fiziksel kasa
  ✅ Ayrı bir offline USB (internete bağlı bilgisayara takılmamış)
  ✅ Güvenli fiziksel konum (banka kasası, kilitli çekmece)
  ❌ Microsoft hesabı (buluta çıkar)
  ❌ OneDrive veya Google Drive
  ❌ E-posta
  ❌ Telefonunun galeri veya not uygulaması
```

**Hangi Durumda Kullanılır?**

| Durum | Sonuç |
|-------|-------|
| BitLocker PIN'ini 5 yanlış girdiysen | Kurtarma anahtarı istenir |
| Anakartı değiştirdinse (TPM sıfırlandı) | Kurtarma anahtarı istenir |
| BIOS güncellemesi TPM state'i değiştirdi | Kurtarma anahtarı istenir |
| Secure Boot ayarı değiştirildiyse | Kurtarma anahtarı istenir |
| SSD başka bilgisayara takıldıysa | Kurtarma anahtarı + PIN gerekir |
| PIN'i unuttuysan | **Sadece** kurtarma anahtarıyla açılır |

> [!CAUTION]
> 48 haneli anahtarı kaybedersen ve PIN'i de unutursan: diskteki tüm veri **matematiksel olarak** sonsuza kadar erişilemez. Bu tasarım gereği böyledir.

---

### AŞAMA 4 — Günlük Güvenlik: Windows Hello
#### 📅 Ne Zaman: Her şey tamamlandıktan sonra, Ayarlar menüsünden

#### Windows Hello PIN ile BitLocker PIN Farkı

Bu iki PIN **tamamen farklı şeylerdir:**

```
┌─────────────────────────────────────────────────────────────────┐
│  BitLocker Pre-Boot PIN                                          │
│  ─────────────────────────────────────────────────────────────  │
│  • Windows logosu GELMEDEN ÖNCE sorulur                         │
│  • Diski çözer                                                   │
│  • Yanlış girilirse sistem boot edemez                          │
│  • TPM + bu PIN = Disk şifrelenmiş → açık                      │
│  • Yalnızca rakamlardan oluşur (önerilen: 8+ hane)             │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│  Windows Hello PIN                                               │
│  ─────────────────────────────────────────────────────────────  │
│  • Windows YÜKLENDİKTEN SONRA kilit ekranında sorulur          │
│  • Kullanıcı oturumunu açar (disk zaten çözülmüş)              │
│  • Rakam + harf + özel karakter içerebilir                     │
│  • Parmak izi veya yüz tanıma da bu aşamada çalışır           │
│  • TPM içinde şifreli olarak saklanır                          │
└─────────────────────────────────────────────────────────────────┘
```

**Windows Hello Nerede Saklanır?**
```
C:\Windows\ServiceProfiles\LocalService\AppData\Local\Microsoft\Ngc\
  → Bu klasör kullanıcıya ait şifreli NGC (Next Generation Credential) verisini tutar
  → TPM yardımıyla korunur; başka bir sisteme taşınamaz
  → HVCI ve Credential Guard aktifse bu klasörün içeriği
    sanallaştırılmış LsaIso.exe sürecinde korunur
```

**Nasıl Ayarlanır?**
```
Ayarlar → Hesaplar → Oturum Açma Seçenekleri
  → Windows Hello PIN → Ekle
  → Windows Hello Parmak İzi (varsa sensör) → Ekle
  → Windows Hello Yüz Tanıma (varsa kamera) → Ekle
```

---

## 🔑 ŞİFRE ANA TABLOSU

| # | Şifre Adı | Nerede Oluşur | Nerede Saklanır | SSD Sökülünce? | Unuutulunca? |
|---|-----------|--------------|-----------------|----------------|--------------|
| 1 | BIOS Admin Şifresi | BIOS menüsü | Anakart NVRAM | **Hâlâ geçerli** | Fiziksel müdahale |
| 2 | BIOS Boot Şifresi | BIOS menüsü | Anakart NVRAM | **Hâlâ geçerli** | Fiziksel müdahale |
| 3 | BitLocker PIN | Windows kurulumu sonrası | TPM + SSD metadata | Yeni PC'de kurtarma anahtarı ister | 48 haneli anahtar |
| 4 | 48H Kurtarma Anahtarı | BitLocker etkinleşince | **Yalnızca sende** | Gerekir | Veri kalıcı kayıp |
| 5 | Windows Kullanıcı Şifresi | Windows kurulumu | SSD (SAM dosyası) | İşlevsiz (yeni PC tanımaz) | Hesap sıfırlama |
| 6 | Windows Hello PIN | Ayarlar menüsü | TPM + SSD (Ngc) | İşlevsiz (TPM bağımlı) | Yeniden oluştur |

---

## 💡 KRİTİK SENARYO ANALİZİ

### "SSD'yi Söküp Başka Bilgisayara Taktım — Ne Olur?"

```
SSD başka PC'ye takıldı
          ↓
BIOS şifresi? → Yeni bilgisayarda YOK (eski anakarta aitti)
          ↓
BitLocker? → AKTİF → Disk şifreli → Kurtarma PIN'i istenir
          ↓
Kurtarma anahtarı girilmeden açılamaz
          ↓
Kurtarma anahtarı da yoksa → VERİ SONSUZA KADAR ERİŞİLEMEZ

Sonuç: SSD fiziksel olarak çalınsa bile içindeki veriye ulaşmak
mevcut teknoloji ile matematiksel olarak imkansızdır.
(XTS-AES 256-bit kırmak → Evrenin ömrü × milyonla çarpımı yıl sürer)
```

### "BitLocker PIN'ini Unuttum — Ne Olur?"

```
BitLocker PIN ekranı → Yanlış PIN
          ↓
5 başarısız deneme → Sistem kurtarma moduna geçer
          ↓
"48 haneli kurtarma anahtarınızı girin" ekranı
          ↓
Anahtarı gir → Disk açılır → Yeni PIN oluştur
          ↓
Anahtarın da yoksa → VERİ SONSUZA KADAR KAYIP
```

### "BIOS Admin Şifresini Unuttum — Ne Olur?"

```
BIOS Admin şifresi unutuldu
          ↓
Seçenek 1: Anakart üreticisinin reset komutu (bazı masaüstülerde jumper)
          ↓
Seçenek 2: CMOS pilini çıkar, 30 dakika bekle, tak → BIOS sıfırlanır
           UYARI: Bu yöntem bazı modern anakartlarda ÇALIŞMAZ
           (NVRAM kalıcı, pil bağımlı değil)
          ↓
Seçenek 3: Üretici servisine götür
          ↓
SONUÇ: BitLocker veri kaybı OLMAZ (SSD hâlâ şifreli, sadece BIOS ayarları gitti)
```

---

## 🔌 USB VE HARİCİ DİSK ŞİFRELEME

### Neden USB'yi Şifrelemek Zorundasın?

```
Şifresiz USB 16 GB:
  ├─ İçinde şifreler.txt, özel fotoğraflar, iş belgeleri
  ├─ USB kayboldu veya çalındı
  └─ Bulan kişi → Tak, aç, oku. Bitti.

Şifreli USB:
  ├─ Aynı senaryo
  └─ Bulan kişi → Tak → "Bu sürücü şifreli" → İçerik OKUNAMAZ
```

---

### Yöntem A — BitLocker To Go (Windows Native, En Kolay)

**Ne gerekir:** Windows 10/11 Pro veya Enterprise. Home'da okuma yapılabilir ama şifreleme için Pro gerekir.

**Nasıl uygulanır:**
```
USB'yi tak → Dosya Gezgini → USB sağ tık → "BitLocker'ı Aç"
           → "Parola kullanarak sürücünün kilidini aç"
           → Şifre gir (güçlü)
           → Kurtarma anahtarını kaydet (yine fiziksel)
           → "Tüm sürücüyü şifrele" seç
           → "Yeni şifreleme modu" seç
           → Şifrelemeyi başlat
```

```powershell
# PowerShell ile:
$USBPin = Read-Host -AsSecureString "USB Şifresi"
Enable-BitLocker -MountPoint "E:" `
    -EncryptionMethod XtsAes256 `
    -PasswordProtector `
    -Password $USBPin `
    -UsedSpaceOnly:$false

# Kurtarma anahtarını al:
(Get-BitLockerVolume -MountPoint "E:").KeyProtector |
    Where-Object {$_.KeyProtectorType -eq "RecoveryPassword"} |
    Select-Object RecoveryPassword
```

**USB başka bilgisayarda açılır mı?**
```
Windows bilgisayar → USB tak → Şifre sor → Gir → Açılır ✅
Linux bilgisayar   → USB tak → dislocker komutu ile açılabilir (şifre gerekir) ✅
macOS              → Üçüncü parti araç gerekir (BitLocker Anywhere) ⚠️
```

**BitLocker To Go Sınırlamaları:**
```
❌ macOS'ta native açılamaz
❌ Eski Windows sürümlerinde (XP, Vista) çalışmaz
❌ TPM bağımlı değil — şifre brute-force'a açık (güçlü şifre seç)
```

---

### Yöntem B — VeraCrypt (Platformdan Bağımsız, Açık Kaynak)

**VeraCrypt nedir?**
Eski TrueCrypt'in güvenli devamı. Windows, Linux ve macOS'ta aynı kapsayıcıyı açabilirsin. NSA, FBI ve GCHQ gibi kurumların incelemesinden geçmiş, bağımsız denetlenmiş.

**İki modu var:**

```
Mod 1: Şifreli Dosya Kapsayıcısı
  USB üzerinde tek bir şifreli ".vc" dosyası oluşturursun.
  Bu dosya içine klasörler, belgeler koyarsın.
  Sadece VeraCrypt ile açılır.
  USB'de hem şifreli dosya hem de şifresiz dosyalar bir arada durabilir.

Mod 2: Tam Bölüm Şifreleme
  USB'nin tamamı VeraCrypt ile şifrelenir.
  USB'yi takan herkes "Biçimlendirilmemiş sürücü" görür.
  Daha güçlü ama VeraCrypt kurulu olmayan PC'de açılamaz.
```

**Kurulum ve Kullanım:**
```
1. https://www.veracrypt.fr → İndir, kur
2. VeraCrypt → "Create Volume"
3. "Create a file container" seç
4. "Standard VeraCrypt volume" seç
5. Dosyanın konumunu ve adını seç (örn: E:\arsiv.vc)
6. Şifreleme algoritması: AES-256 + SHA-512 (varsayılan, yeterli)
7. Boyut belirle (örn: 10 GB)
8. Şifre gir (güçlü, 20+ karakter)
9. "Format" → Kapsayıcı oluşturulur
10. "Mount" → Sürücü harfi seç → Şifre → Aç
11. Sanal sürücü masaüstünde beliririr → Dosyaları oraya kopyala
12. "Dismount" → Kapsayıcı kapanır, içerik şifreli
```

**VeraCrypt Kriptografi Özellikleri:**
```
Şifreleme: AES-256-XTS
Hash:       SHA-512 (anahtar türetme)
PBKDF2 iterasyon: 500.000 (brute-force'u yavaşlatır)
Gizli Volüm (Hidden Volume): Zorlama altında farklı şifre, farklı içerik
```

**Gizli Volüm (Plausible Deniability) Özelliği:**
```
Dışarı: Dolu görünen VC kapsayıcısı — "E:\arsiv.vc"

İçeride iki katman:
  ├── Dış volüm (Şifre A): Masum içenik — belgeler, fotoğraflar
  └── İç gizli volüm (Şifre B): Gerçek gizli veriler

Baskı altındaysan Şifre A'yı ver → Masum içerik görünür
Şifre B'nin varlığından kriptografik olarak ispat edilemez
```

```powershell
# VeraCrypt komut satırı (kurulduktan sonra):
# Kapsayıcıyı mount et
& "C:\Program Files\VeraCrypt\VeraCrypt.exe" /v "E:\arsiv.vc" /l Z /p "ŞifrenBurada" /q

# Dismount et
& "C:\Program Files\VeraCrypt\VeraCrypt.exe" /d Z /q
```

---

### Yöntem C — Donanım Şifreli USB (En Güçlü)

Yazılıma hiç bağımlı olmayan, kendi içinde şifreleme yapan USB bellekler.

| Model | Şifreleme | Özellik | Fiyat Aralığı |
|-------|-----------|---------|---------------|
| **Kingston IronKey D500S** | AES-256 XTS (FIPS 140-3 Lvl 3) | 10 yanlış → kendi kendini imha | Yüksek |
| **Kingston IronKey Keypad 200** | AES-256 XTS | Fiziksel tuş takımı (PC gerekmiyor) | Orta-Yüksek |
| **Apricorn Aegis Secure Key 3NXC** | AES-256 XTS (FIPS 140-2 Lvl 3) | PIN tuş takımı, brute-force kilidi | Orta-Yüksek |
| **Verbatim Keypad Secure** | AES-256 | Tuş takımı, giriş | Orta |

**Nasıl Çalışır?**
```
USB'yi tak
    ↓
USB ÜZERİNDEKİ tuş takımına PIN gir (PC klavyesine değil)
    ↓
PIN doğruysa USB kendisi şifreyi çözer → PC "normal USB" görür
    ↓
10 yanlış PIN → Firmware tüm veriyi şifreler ve anahtarı imha eder
    ↓
USB tahrip edilip NAND chip'i başka okuyucuya takılsa bile
şifreleme donanım içinde olduğu için veri OKUNAMAZ
```

**Avantajları:**
```
✅ Keylogger'a karşı güvenli (PIN klavyeye değil USB'ye giriliyor)
✅ Yazılım bağımsız — her OS, her PC'de çalışır
✅ Brute-force sonrası otomatik imha
✅ Donanım FIPS 140-3 sertifikalı (askeri standart)
✅ Hiçbir sürücü kurulumu gerektirmez
```

**Dezavantajları:**
```
❌ Pahalı (normal USB'nin 5-15 katı)
❌ PIN forgotten = Veri kalıcı kayıp (tasarım gereği)
```

---

### USB Şifreleme Yöntemleri Karşılaştırması

| Özellik | BitLocker To Go | VeraCrypt | Donanım Şifreli |
|---------|----------------|-----------|-----------------|
| Kurulum zorluğu | Kolay | Orta | Çok Kolay |
| Platform desteği | Windows native | Win+Linux+Mac | Tüm platformlar |
| Keylogger koruması | ❌ | ❌ | ✅ |
| Brute-force kilidi | ❌ | PBKDF2 yavaşlatır | ✅ (otomatik imha) |
| Gizli volüm | ❌ | ✅ | ❌ |
| Maliyet | Ücretsiz | Ücretsiz | Yüksek |
| FIPS Sertifikası | ❌ | ❌ | ✅ (bazı modeller) |
| Şifre unutulursa | Kurtarma anahtarı | Yedek yoksa kayıp | Veri imha |
| Güvenlik seviyesi | ★★★☆☆ | ★★★★☆ | ★★★★★ |

---

## 📋 NİHAİ KONTROL LİSTESİ

Kurulum tamamlandığında elinde şunların olması gerekiyor:

### Fiziksel Kasada Olması Gerekenler (Kağıt/Laminat)

```
┌─────────────────────────────────────────────────────────┐
│  GÜVENLİ KURULUM BELGE KASASI                          │
│                                                          │
│  [1] BIOS Admin Şifresi: _____________________________ │
│      (Anakart: _______________, Tarih: _____________)  │
│                                                          │
│  [2] BIOS Boot Şifresi: ______________________________ │
│                                                          │
│  [3] BitLocker Pre-Boot PIN: _________________________ │
│      (C: sürücüsü için)                                │
│                                                          │
│  [4] BitLocker 48H Kurtarma Anahtarı:                  │
│      ______-______-______-______-______-______-______-_│
│      (Kurtarma Kimliği: _____________________________)  │
│                                                          │
│  [5] Windows Kullanıcı Şifresi: _____________________ │
│      (Kullanıcı Adı: _______________________________)  │
│                                                          │
│  [6] USB Şifresi (BitLocker/VeraCrypt): ______________ │
│      (Hangi USB: __________________________________)    │
│                                                          │
│  Oluşturma Tarihi: ________________________________     │
└─────────────────────────────────────────────────────────┘
```

### Dijital Ortamda OLMAMASI Gerekenler

```
❌ Bu belgeler şu yerlerde BULUNMAMALI:
   → Microsoft hesabı / OneDrive
   → Google Drive / Dropbox / iCloud
   → E-posta gelen kutusu veya taslaklar
   → Telefonun not uygulaması
   → Şifresiz bir USB
   → Ekran fotoğrafı olarak galeri
   → Bilgisayar masaüstünde "şifreler.txt"
```

---

## ✅ HIZLI ADIM REFERANSI

```
AŞAMA 1 — BIOS (Windows yok)
  □ BIOS Admin Şifresi koy
  □ BIOS Boot Şifresi koy
  □ Secure Boot → Enabled
  □ TPM → Enabled
  □ VT-x + VT-d → Enabled
  □ Boot order'ı kilitle

AŞAMA 2 — Windows Kurulum
  □ İnternetsiz kur
  □ Yerel hesap oluştur
  □ Tüm telemetri seçeneklerini kapat

AŞAMA 3 — BitLocker
  □ gpedit.msc → TPM+PIN zorunlu kıl
  □ BitLocker'ı XTS-AES 256 ile etkinleştir
  □ Pre-Boot PIN oluştur
  □ 48 haneli anahtarı yaz → kasaya kaldır
  □ manage-bde -status C: ile doğrula

AŞAMA 4 — Windows Hello
  □ Windows Hello PIN oluştur (BitLocker PIN'inden farklı)
  □ Biyometri varsa ekle

AŞAMA 5 — USB Şifreleme
  □ BitLocker To Go veya VeraCrypt seç
  □ USB'yi şifrele
  □ USB kurtarma anahtarını da kasaya ekle

AŞAMA 6 — Doğrulama
  □ Yeniden başlat → BitLocker PIN soruldu mu?
  □ Masaüstü açıldı mı?
  □ Windows Hello çalışıyor mu?
  □ USB şifreli açılıyor mu?
```

---

*Son güncelleme: 2026-03-22 · Windows 11 Pro / Windows 10 Pro için hazırlanmıştır.*
