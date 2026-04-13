<div align="center">

<br>

```
      ╔══════════════════════════════════════════════════╗
      ║          🏰  KANİJE KALESİ  🏰                  ║
      ║   Katmanlı Siber Savunma · Güvenlik İzleme       ║
      ╚══════════════════════════════════════════════════╝
```

*9.000 asker, 100.000 kişilik orduyu 73 gün boyunca durdurdu.*  
*Biz de aynısını yapıyoruz — sadece savaş alanı değişti.*

<br>

[![Go](https://img.shields.io/badge/Kanije_Kalesi-Go_1.21+-00ADD8?style=flat-square&logo=go)](go/)
[![Python](https://img.shields.io/badge/Legacy-Python_3.11+-3776AB?style=flat-square&logo=python)](app/)
[![Platform](https://img.shields.io/badge/Platform-Windows_%7C_Linux_%7C_Raspberry_Pi-555?style=flat-square)](#)
[![License](https://img.shields.io/badge/License-MIT-22c55e?style=flat-square)](LICENSE)

<br>

</div>

---

## Bu Repo Ne İçerir?

Bu depo iki bölümden oluşur:

| Bölüm | Klasör | Ne Yapar? |
|-------|--------|-----------|
| **Kanije Kalesi Uygulaması** | [`go/`](go/) | Gerçek zamanlı güvenlik izleme + Telegram bot |
| **Sertleştirme Rehberleri** | Kök dizin `.md` dosyaları | Windows / Linux güvenlik sertleştirme kılavuzları |

> Her bölüm **tamamen bağımsızdır.** Kanije Kalesi olmadan rehberleri okuyabilir, rehberlere bakmadan Kanije Kalesi'i kurabilirsiniz.

---

<br>

# 🏰 Bölüm I — Kanije Kalesi Güvenlik İzleme Aracı

Sisteminize bağlı bir bekçi gibi çalışır. Şüpheli bir şey olduğunda — birisi yanlış şifre girdiğinde, USB takıldığında, bilgisayar uykudan uyandığında — **anında Telegram'a bildirir.**

Tüm ayarlar Telegram üzerinden yapılır. Config dosyasına hiç dokunmanıza gerek yok.

<br>

## İzlenen Olaylar

| Olay | Bildirim | Otomatik Kamera |
|------|----------|-----------------|
| ✅ Başarılı oturum açma | ✓ | İsteğe bağlı |
| 🚨 Başarısız giriş denemesi | ✓ | **Evet** |
| 🔒 Ekran kilitlendi / kilidi açıldı | ✓ | — |
| 🖥️ Sistem başlatıldı / kapandı | ✓ | — |
| 😴 Uyku / uyanma | ✓ | — |
| 🔌 USB takıldı / çıkarıldı | ✓ | — |
| 🌐 İnternet bağlantısı değişti | ✓ | — |
| 💓 Periyodik durum raporu | ✓ | — |

<br>

## Telegram Komutları

```
/status    →  CPU, RAM, disk, çalışma süresi
/foto      →  Anlık kamera fotoğrafı
/ekran     →  Ekran görüntüsü
/olaylar   →  Son 10 güvenlik olayı
/kilitle   →  Ekranı kilitle
/yeniden   →  Sistemi yeniden başlat (onay gerekli)
/kapat     →  Sistemi kapat (onay gerekli)
/kurulum   →  ⚙️ Etkileşimli ayar menüsü
/yardim    →  Tüm komutlar
```

<br>

## Kurulum

### Ön Gereksinimler

| | Gereksinim | Kontrol |
|-|-----------|---------|
| ✅ | **Go 1.21+** | `go version` |
| ✅ | **ffmpeg** | `ffmpeg -version` |
| ✅ | **Telegram Bot Token** | [@BotFather](https://t.me/BotFather) → `/newbot` |
| ✅ | **Telegram Chat ID** | [@userinfobot](https://t.me/userinfobot) → `/start` |

<br>

### Adım 1 — İndirin

```bash
git clone https://github.com/mehmetyasinuzun/Kanije-Kalesi.git
cd Kanije-Kalesi/go
```

### Adım 2 — Bağımlılıkları kurun

```bash
go mod tidy
```

### Adım 3 — Telegram bilgilerini kaydedin

```bash
go run ./cmd/kanije/ setup \
  --token "1234567890:AABBccDDee..." \
  --chat  "123456789"
```

Bu komut `config.toml` dosyasını oluşturur. Başka hiçbir şey yapmanıza gerek yok.

### Adım 4 — Test edin

```bash
go run ./cmd/kanije/ test
```

```
✅ Bağlantı başarılı!
   Bot adı : @benim_botum
   Chat ID : 123456789
```

### Adım 5 — Başlatın

```bash
go run ./cmd/kanije/ start
```

**Telegram botunuza `/kurulum` yazın → her şeyi oradan ayarlayın.**

<br>

## Telegram Kurulum Menüsü

`/kurulum` komutu şu etkileşimli menüyü açar:

```
⚙️ Kanije Kalesi — Kurulum Menüsü

  [🎯 Tetikleyiciler]
  [📷 Kamera Ayarları]
  [💓 Heartbeat]
  [🔐 Güvenlik]
  [📋 Loglama]
  [✅ Tamamlandı]
```

Her kategori kendi alt menüsünü açar. Toggle'lar tek tıkla değişir. Sayı gerektiren ayarlarda bot sizden yazmanızı ister. Tüm değişiklikler anında `config.toml`'a kaydedilir.

<br>

## Otomatik Başlatma

### Windows — Task Scheduler

```powershell
# Derle (konsol penceresi açmayan sessiz binary)
go build -ldflags="-s -w -H=windowsgui" -o kanije.exe ./cmd/kanije/

# Task Scheduler'a kaydet
.\deploy\windows\install.ps1

# Kontrol
.\deploy\windows\install.ps1 -Status

# Kaldır
.\deploy\windows\install.ps1 -Remove
```

### Linux — systemd

```bash
# Derle
go build -o kanije ./cmd/kanije/

# Kur
sudo bash deploy/linux/install.sh

# Telegram bilgilerini yaz
sudo nano /etc/kanije/secrets.env

# Başlat ve logları izle
sudo systemctl start kanije-kalesi
sudo journalctl -u kanije-kalesi -f
```

### Raspberry Pi

```bash
# Pi 4 / Pi 5 için cross-compile (CGo gerektirmez)
make build-arm64

# Pi 3 için
make build-arm

# Pi'ye kopyala
scp dist/kanije-linux-arm64 pi@raspberrypi:/usr/local/bin/kanije

# Pi'de kur
ssh pi@raspberrypi "sudo bash /path/to/deploy/linux/install.sh"
```

Pi'de ek paket:
```bash
sudo apt install ffmpeg dbus-x11 iw
```

<br>

## Derleme Hedefleri

```bash
make build          # Mevcut platform
make build-windows  # Windows AMD64
make build-linux    # Linux AMD64
make build-arm64    # Raspberry Pi 4 / 5
make build-arm      # Raspberry Pi 3
make build-all      # Tüm platformlar aynı anda
```

<br>

## Ortam Değişkenleri

Config dosyası oluşturmadan da çalıştırabilirsiniz:

```bash
export KANIJE_BOT_TOKEN="1234567890:AABBcc..."
export KANIJE_CHAT_ID="123456789"
export KANIJE_LOG_LEVEL="info"   # debug | info | warn | error

./kanije start
```

<br>

---

<br>

# 📚 Bölüm II — Güvenlik Sertleştirme Rehberleri

<br>

## ⚔️ Neden "Kanije Kalesi"?

**1601, Kanije (bugünkü Nagykanizsa, Macaristan).**

73 yaşındaki Tiryaki Hasan Paşa, elinde 9.000 asker ve 100 küçük topla, Habsburg Arşidükü II. Ferdinand'ın komutasındaki **100.000 kişilik Haçlı ordusuna** karşı Kanije Kalesi'ni 73 gün savundu.

Cephane bitti — kale içinde baruthane kurdu. Erzak tükendi — sahte mektuplarla "padişahın ordusu yolda" dedirtti. Duvarlar yıkıldı — gece onarıp sabah sapasağlam gösterdi.

73. gece, "artık bitti" denildiğinde, gece baskınıyla Arşidük'ün karargahını bastı. 100.000 kişilik ordu kaçtı. 9.000 kişi kazandı.

**Çünkü önemli olan sayı değil, katmanlı savunmaydı.**

<br>

## Rehberler

### İşletim Sistemi Sertleştirme

| Rehber | İçerik |
|--------|--------|
| [WINDOWS11_HARDENING_KALE.md](WINDOWS11_HARDENING_KALE.md) | Windows 11 — BIOS'tan Sysmon'a 7 katman |
| [WINDOWS10_HARDENING_KALE.md](WINDOWS10_HARDENING_KALE.md) | Windows 10 — 7 katman + telemetri kapatma |
| [LINUX_HARDENING_KALE.md](LINUX_HARDENING_KALE.md) | Kali Linux / Ubuntu — LUKS2, AppArmor, SSH |

### Disk Güvenliği ve Şifreleme

| Rehber | İçerik |
|--------|--------|
| [DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md](DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md) | Dual boot saldırı yüzeyleri, ayrı disk mimarisi |
| [SIFRE_KRONOLOJISI_VE_USB_SIFRELEME.md](SIFRE_KRONOLOJISI_VE_USB_SIFRELEME.md) | BitLocker, VeraCrypt, donanım şifreli USB |

### Harici Disk ve Multiboot

| Rehber | İçerik |
|--------|--------|
| [HARICI_USB_SSD_BOOT_REHBERI.md](HARICI_USB_SSD_BOOT_REHBERI.md) | Windows To Go, USB SSD'den boot |
| [VENTOY_WTG_MULTIBOOT_REHBERI.md](VENTOY_WTG_MULTIBOOT_REHBERI.md) | Ventoy + WTG + Linux multiboot, VHD, rEFInd |

### Araçlar

| Rehber | İçerik |
|--------|--------|
| [HAYAT_KURTARAN_YAZILIMLAR.md](HAYAT_KURTARAN_YAZILIMLAR.md) | Format sonrası kurulacak yazılımlar — winget script |

<br>

## Katmanlı Savunma Modeli

```
KATMAN 0 — Donanım / BIOS-UEFI
  BIOS şifresi · Secure Boot · TPM 2.0
  Düşman kapıya bile gelemez.

KATMAN 1 — Önyükleme / Disk Şifreleme
  BitLocker Pre-Boot PIN · LUKS2
  Disk çalınsa bile içi okunamaz.

KATMAN 2 — Çekirdek
  VBS · HVCI · Credential Guard · LSA PP
  Ring-0 saldırıları duvara toslar.

KATMAN 3 — Bellek
  Pagefile temizleme · hibernate kapatma · DMA koruması
  Cold boot saldırısına karşı.

KATMAN 4 — Süreç / Uygulama
  Exploit Guard · ASLR · DEP · ASR kuralları

KATMAN 5 — Ağ
  SMB sertleştirme · DNS-over-HTTPS · LLMNR/WPAD kapatma

KATMAN 6 — Kimlik / Hesap
  Standart hesap · UAC max · FIDO2 · hesap kilitleme

KATMAN 7 — Denetim / İzleme  ← Kanije Kalesi burada devreye girer
  Sysmon · gelişmiş denetim · PowerShell günlükleme
  Gözcü kulesinden her hareket görünür.
```

<br>

## Hızlı Başlangıç — Hangi Rehberi Okuyayım?

```
Yeni bilgisayar mı aldın?
  → SIFRE_KRONOLOJISI_VE_USB_SIFRELEME.md
  → WINDOWS11_HARDENING_KALE.md

Linux mu kuruyorsun?
  → LINUX_HARDENING_KALE.md

Dual boot mu yapacaksın?
  → DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md  (riskleri önce bil)

Harici diskten boot mu edeceksin?
  → HARICI_USB_SSD_BOOT_REHBERI.md
  → VENTOY_WTG_MULTIBOOT_REHBERI.md

Format mı attın?
  → HAYAT_KURTARAN_YAZILIMLAR.md  (tek scriptle her şeyi kur)

Sistemi izlemek mi istiyorsun?
  → go/  (Kanije Kalesi — Telegram bildirim + kamera)
```

<br>

---

## Güvenlik Notu

Bu rehberler kendi sistemlerini güçlendirmek isteyen bireyler için hazırlandı. Buradaki bilgiler savunma amaçlıdır. Başkalarının sistemlerine yetkisiz erişim yasa dışıdır.

Tiryaki Hasan Paşa saldırmadı — savundu.

---

## Lisans

[MIT](LICENSE) © 2026 Kanije Kalesi

<br>

<div align="center">
<sub><i>"Kale düşmez, kalenin içindekiler düşer."</i></sub>
<br>
<sub>🏰</sub>
</div>
