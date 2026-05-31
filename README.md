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

[![CI](https://github.com/mehmetyasinuzun/Kanije-Kalesi/actions/workflows/ci.yml/badge.svg)](https://github.com/mehmetyasinuzun/Kanije-Kalesi/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Kanije_Kalesi-Go_1.21+-00ADD8?style=flat-square&logo=go)](go/)
[![Python](https://img.shields.io/badge/Python-deprecated-9ca3af?style=flat-square&logo=python)](app/DEPRECATED.md)
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
/pil       →  🔋 Pil durumu (yüzde, şarj, kalan süre)
/foto      →  Anlık kamera fotoğrafı
/ekran     →  Ekran görüntüsü
/seskayit  →  Mikrofon kaydı (/seskayit 30 → 30 sn, varsayılan 30, en çok 600)
/pano      →  📋 Panodaki metni getir
/panik     →  🆘 Tek komutla kanıt topla: foto+ekran+ses+dış IP (/panik kilit → ekranı da kilitler)
/hareket   →  🎥 Kamera hareket algılama (/hareket ac · kapat · esik <n>)
/dinle     →  🎧 Canlı dinleme (/dinle 10 → 10 dk, parça parça ses · /dinle kapat)
/olaylar   →  Son olaylar (/olaylar <tip> veya <sayı> ile filtrele)
/ozet      →  Son 7 günün olay özeti (tip bazlı)
/dogrula   →  Olay günlüğü bütünlüğünü doğrula (hash-chain)
/defender  →  🛡️ Microsoft Defender durumu + son taramalar + tespitler
/guncelle  →  Yeni sürümü kontrol et ve kur
/dosya     →  📁 Dosya gez/indir (/dosya <yol> · /dosya al <yol>)
/erisim    →  👁️ Dosyana kim erişti/kopyaladı (/erisim kur <yol> · yönetici gerekir)
/zamanla   →  ⏰ Komut zamanla (/zamanla 30dk /foto · liste · sil <id>)
/kilitle   →  Ekranı kilitle
/yeniden   →  Sistemi yeniden başlat (onay + 15 sn geri-al)
/kapat     →  Sistemi kapat (onay + 15 sn geri-al)
/cihazlar  →  🛰️ Tüm cihazları listele (fleet)
/terminal  →  💻 Uzak komut çalıştır (cd kalıcı; /terminalix yönetici)
/ekle      →  👥 Kişi ekle (/ekle <chat_id> [isim])
/yonetim   →  👥 Eklediklerini gör/düzenle/çıkar
/loglar    →  🧾 Son işlemler (kim ne yaptı)
/kurulum   →  ⚙️ Etkileşimli ayar menüsü
/yardim    →  Tüm komutlar

# 🛡️ Sahip — geri dönüşsüz (yalnız cihaz sahibi · çift onay + 15 sn geri-al)
/kaldir    →  🧹 Kanije'yi izsiz kaldır (görev + dosyalar + eski sürümler)
/aktar     →  🔁 Botu yeni sahibe devret (/aktar <chat_id> [token])
/imha      →  💥 Kanije verisi + Müzik klasörü içi güvenli silinir (/imha ONAYLA · OS korunur)
/koruma    →  🛡️ Fiziksel tehdit koruması: dead-man switch · USB dead-man · yanlış-giriş · RAM-only
/kilit tam →  🔒 Tam kilit (lockdown): ekran açılınca anında tekrar kilitlenir (/kilit tam kapat ile aç)
/tuzak     →  🍯 Honeypot tuzak dosyalar: dokunan yakalanır → alarm/kilit (/tuzak kur)
```

<br>

## Gelişmiş Özellikler

| Özellik | Açıklama |
|---------|----------|
| 🔐 **Token şifreleme** | Windows'ta bot token DPAPI ile şifrelenir — config çalınsa bile başka hesapta açılmaz |
| 🌍 **IP konum + dış IP** | Başarısız girişte kaynak IP ülke/şehir/ISP ile zenginleşir (🇷🇺); ağ olaylarında dış IP gösterilir |
| 🔗 **Kurcalama tespiti** | Olay günlüğü hash-chain ile imzalı; `/dogrula` ile bütünlük kontrolü |
| 🔄 **Otomatik güncelleme** | `/guncelle` ya da günlük kontrol → indirir, gizli kurar, yeniden başlar (Windows) |
| 📣 **Çoklu hedef** | Telegram + Discord/webhook'a aynı anda bildirim (`[[webhooks]]`) |
| 🔇 **Sessiz saatler** | Gece penceresinde yalnız kritik olaylar (`[quiet_hours]`) |
| 🛡️ **Anti-abuse** | Komut hız sınırı + opsiyonel 2FA (TOTP) hassas komutlarda |
| 📊 **Prometheus** | Opsiyonel `/metrics` endpoint (homelab izleme) |
| 🎤 **Ses kaydı** | `/seskayit [saniye]` — mikrofondan kayıt, Telegram'a MP3 (varsayılan 30 sn) |
| 👥 **Çok kullanıcı** | Yetki devri ağacı: kişi ekle/çıkar, rol + komut-bazlı yetki, yükseltme yok, herkes kendi alt-ağacını görür, `/loglar` denetim izi |
| 🛰️ **Fleet (çok cihaz)** | Tek grup, N cihaz: her cihaza ayrı bot, komutlar cihaz adıyla yönlenir (`/foto dizustu`), `/cihazlar` listeler |
| 💻 **Uzak terminal** | `/terminal <komut>` etkileşimli shell (cwd kalıcı), `/terminalix` yönetici; `terminal` yetkisiyle korunur, denetlenir |
| 📋 **Derin olay detayı (CTI)** | Her olayda ağ medium'u (WiFi/USB tethering/Bluetooth/Hücresel/VPN), SSID, iç/dış IP; `/olaylar`'da olay başına butonla tam detay |
| 🛑 **Kurcalama alarmı** | Watchdog: exe/config/DB silinmesi, otomatik-başlatma görevinin kapatılması ve **beklenmedik kapanma** (kill/çökme) → sahibe kritik alarm + kurcalayanın fotoğrafı |
| 🧹 **İzsiz kaldırma** | `/kaldir` — görev, dosyalar, **eski sürümler** (.new/.old), config, DB, log, captures, temp betikleri ve kendini siler |
| 🔁 **Devretme** | `/aktar <chat_id> [token]` — botu yeni sahibe devreder, erişim ağacını sıfırlar, yeni kimlikle yeniden başlar |
| 💥 **Uzaktan imha** | `/imha ONAYLA` — Kanije verisi (token/geçmiş) + **Müzik klasörü içi** + ek hedefleri güvenli üzerine-yazıp siler; **OS korunur** (fabrika sıfırlama yok, düşük AV yüzeyi). Hedefler: `/imha hedef ekle/sil` |
| 🆘 **Panik modu** | `/panik` — tek komutla foto + ekran + ses + dış IP toplar; `/panik kilit` ekranı da kilitler |
| 📁 **Dosya erişimi** | `/dosya` ile uzaktan dizin gez, `/dosya al <yol>` ile dosya indir (≤45 MB) |
| 📋 **Pano & pil** | `/pano` panodaki metni getirir, `/pil` batarya durumunu (yüzde/şarj/kalan) gösterir |
| 🗑️ **Gönder-sil** | `SaveLocal` açıkken yerel kopya Telegram'a **başarıyla** gittikten sonra diskten silinir (offline'da korunur) |
| ⏰ **Zamanlama** | `/zamanla 30dk /foto` — periyodik komut; restart'a dayanıklı (JSON), kaçan slotlar tek tetiklenir |
| 🎥 **Hareket algılama** | `/hareket ac` — kamera kare farkı (luma) eşiği aşınca otomatik foto + olay |
| 🎧 **Canlı dinleme** | `/dinle 10` — mikrofonu 20 sn'lik parçalar halinde akıtır, `/dinle kapat` durdurur |
| 🔒 **Binary bütünlüğü** | Watchdog exe SHA-256'sını izler; runtime'da değiştirilirse (kod enjeksiyonu / truva'lı kopya) kritik alarm |
| 🛡️ **Defender izleme** | `/defender` — gerçek-zamanlı koruma durumu, son hızlı/tam tarama zamanı, imza sürümü ve son tespitler (AV seni ne kadar/ne zaman izliyor) |
| 🛡️ **Fiziksel tehdit koruması** | `/koruma` — yapılandırılabilir politika motoru: **dead-man switch** (X saat check-in yoksa), **USB dead-man** (BusKill), **yanlış-giriş sayacı** → her biri kendi aksiyonuyla (kilitle / alarm+foto / güvenli sil). Wipe opt-in + 60 sn geri-alma. **RAM-only mod**: DB hiç diske yazılmaz |
| 🔒 **Tam kilit (lockdown)** | `/kilit tam` — "Lost Mode": ekran her açıldığında ~2 sn içinde tekrar kilitlenir, doğru şifre girilse bile cihaz kullanılamaz; yalnızca `/kilit tam kapat` ile açılır |
| 🍯 **Honeypot tuzak** | `/tuzak kur` — cazip **sahte** dosyalar (Şifreler.txt vb.) + SACL denetimi; biri **dokununca/kopyalayınca** → alarm + foto + (opsiyonel) lockdown. Saldırganın verisini **kopyalamaz**, sadece tuzağa dokunduğunu yakalar |
| 👁️ **Dosya erişim denetimi** | `/erisim kur <yol>` — Windows SACL denetimiyle bir klasöre **kim/ne zaman/hangi programla** eriştiğini (okuma/kopyalama/yazma/silme) Security log'dan okur (`/erisim`). Yönetici/SYSTEM gerekir |

<br>

## 👥 Çok Kullanıcı (Yetki Devri)

Aynı botu güvenle başkalarıyla paylaş. Sahip (kurulumdaki kişi) tüm yetkilere sahiptir; başkalarını **kısıtlı** ekleyebilir:

- `/ekle <chat_id> [isim]` → kişiyi ekler, ardından **rol** (İzleyici/Operatör/Yönetici) seç veya **tek tek** yetki aç/kapat.
- **Yükseltme yok:** Birine yalnızca **sende olan** yetkileri verebilirsin (sunucu tarafında zorlanır — istek manipülasyonu işe yaramaz).
- **Görünürlük:** Herkes yalnızca **kendi eklediklerini** (ve onların eklediklerini) görür/yönetir.
- **Çıkarma kişiyi silmez-zinciri koparmaz:** B çıkarılırsa altındakiler silinmez; **en tepeye** (veya bir üste / sana) yeniden bağlanır, yeni ebeveynin yetkisine kırpılır.
- `/yonetim` ile ağacı gör, düzenle, çıkar · `/loglar` ile kim ne yaptı.

## 🛰️ Çoklu Cihaz (Fleet)

Birden fazla bilgisayarı **tek Telegram grubundan** yönet — sunucu/port gerekmez:

1. Her cihaz için ayrı bir bot oluştur (@BotFather) ve **gizlilik modunu kapat** (`/mybots → Bot Settings → Group Privacy → Turn off`).
2. Bir Telegram **grubu** kur, tüm bot'ları gruba ekle.
3. Her cihazda `/kurulum → 🛰️ Fleet / Cihaz`: **cihaz adı** (örn. `dizustu`) ve **grup ID**'sini gir.
4. Artık olaylar gruba düşer (hangi cihaz olduğu etiketli). Komutları cihaz adıyla yönlendir:

```
/cihazlar           →  tüm cihazlar kendini listeler
/foto dizustu       →  yalnızca "dizustu" fotoğraf çeker
/seskayit ev-pc 60  →  "ev-pc" 60 sn ses kaydı alır
```

<br>

## Kurulum

### ⚡ Hızlı Kurulum — ~10 saniye, Go gerekmez (önerilen)

Hazır binary indirilir, **gizli** (masaüstünde ikon/pencere yok) ve **yükseltilmiş**
(giriş denemelerini görebilmek için) bir otomatik-başlatma görevi olarak kurulur. Tek onay.

**Windows** — PowerShell'de tek satır. **Önerilen (indir → çalıştır, token/chat sorulur):**
```powershell
irm "https://raw.githubusercontent.com/mehmetyasinuzun/Kanije-Kalesi/master/go/deploy/windows/install.ps1" -OutFile "$env:TEMP\setup.ps1"; & "$env:TEMP\setup.ps1"
```

Token/chat'i **baştan vermek** istersen sona ekle (sorulmaz):
```powershell
irm "https://raw.githubusercontent.com/mehmetyasinuzun/Kanije-Kalesi/master/go/deploy/windows/install.ps1" -OutFile "$env:TEMP\setup.ps1"; & "$env:TEMP\setup.ps1" -Token "BOT_TOKEN" -Chat "CHAT_ID"
```

En kısa (dosya bırakmaz; token/chat sorulur):
```powershell
irm https://raw.githubusercontent.com/mehmetyasinuzun/Kanije-Kalesi/master/go/deploy/windows/get.ps1 | iex
```

Alternatif: [Releases](../../releases)'tan `kanije-windows-amd64.exe` + `install.bat`'ı indirip
**install.bat'a çift tıklayın**.

**Linux / Raspberry Pi** — tek satır:
```bash
curl -fsSL https://raw.githubusercontent.com/mehmetyasinuzun/Kanije-Kalesi/master/go/deploy/linux/get.sh | sudo bash -s -- --token "BOT_TOKEN" --chat "CHAT_ID"
```

Kurulumdan sonra Telegram botunuza **/kurulum** yazın — gerisi oradan.

> 🔇 **Gizli çalışma:** Ajan masaüstünde hiçbir iz bırakmaz (konsol penceresi yok,
> sistem tepsisi ikonu yok). Bilgisayara izinsiz erişmek isteyen biri uygulamayı göremez.

<br>

### Manuel / Geliştirici Kurulumu (Go ile)

#### Ön Gereksinimler

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

## Geliştirme & Test

Derleme, test ve lint kapıları GitHub Actions ile her push/PR'da **Windows + Linux
(amd64 / arm64 / arm)** matrisinde çalışır — yarış dedektörü ve `golangci-lint` dahil.

```bash
cd go
go test ./...            # birim testleri
go test -race ./...      # yarış dedektörü (Linux + CGo)
gofmt -l .               # biçim — çıktı boş olmalı
go vet ./...
golangci-lint run        # statik analiz (go/.golangci.yml)
```

Katkı rehberi: [CONTRIBUTING.md](CONTRIBUTING.md) · Güvenlik politikası: [SECURITY.md](SECURITY.md)

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
