# 🏰 Kanije Kalesi — Teknik Walkthrough

## Proje Özeti

**14 Python dosyası**, **6 modül**, **~1200 satır kod**. Windows'ta arka planda çalışan, olay tabanlı güvenlik izleme uygulaması. Oturum giriş/çıkışlarını, başarısız giriş denemelerini, USB olaylarını izler; kameradan fotoğraf çeker; Telegram'a bildirim gönderir; Telegram'dan komut alır.

---

## 📁 Dosya Yapısı

```
Kanije-Kalesi/app/
├── kanije.py                        ← CLI giriş noktası
├── config.yaml                      ← Tüm ayarlar (kullanıcı düzenler)
├── requirements.txt                 ← Python bağımlılıkları
│
├── core/                            ← Çekirdek altyapı
│   ├── app.py                       ← Ana orkestratör (tüm modülleri yönetir)
│   ├── config_manager.py            ← YAML okuma, ENV çözümleme, doğrulama
│   ├── event_bus.py                 ← Thread-safe olay kuyruğu + rate limiter
│   └── logger.py                    ← RotatingFileHandler loglama
│
├── listeners/                       ← Olay dinleyicileri
│   ├── windows_event.py             ← Windows Security Event Log (4624/4625/4800/4801)
│   └── usb_monitor.py               ← WMI ile USB algılama
│
├── actions/                         ← Aksiyon modülleri
│   ├── camera.py                    ← OpenCV kamera çekim
│   └── screenshot.py                ← mss ekran görüntüsü
│
├── telegram_bot/                    ← Telegram iletişim
│   ├── notifier.py                  ← Tek yönlü: mesaj/fotoğraf gönderme
│   ├── bot_commands.py              ← İki yönlü: komut alma ve yanıtlama
│   └── formatter.py                 ← Mesaj şablonları (emoji formatları)
│
└── tray/                            ← Sistem tepsisi
    └── tray_icon.py                 ← Yeşil/kırmızı ikon + sağ tık menü
```

---

## 🔄 Veri Akış Diyagramı

```
Windows Event Log ──┐                          ┌── Telegram Mesaj
(4624/4625/4800)    │                          │
                    ├──→ EventBus ──→ app.py ──┤── Telegram Fotoğraf
USB WMI Watcher ────┘    (queue)    (process)  │
                                               └── Log Dosyası

Telegram Komutları (/status, /photo...) ──→ bot_commands.py ──→ Telegram Yanıt
```

---

## 🔌 Kullanılan API ve Endpoint'ler

### Telegram Bot API

| Endpoint | Metod | Ne İçin | Dosya |
|----------|-------|---------|-------|
| `/bot{TOKEN}/getMe` | GET | Bot doğrulama | [notifier.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/telegram_bot/notifier.py) |
| `/bot{TOKEN}/sendMessage` | POST | Metin mesajı gönder | [notifier.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/telegram_bot/notifier.py) |
| `/bot{TOKEN}/sendPhoto` | POST | Fotoğraf gönder (multipart/form-data) | [notifier.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/telegram_bot/notifier.py) |
| `/bot{TOKEN}/sendDocument` | POST | Dosya gönder | [notifier.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/telegram_bot/notifier.py) |
| `/bot{TOKEN}/getUpdates` | GET | Komutları al (long polling) | [bot_commands.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/telegram_bot/bot_commands.py) |

**Rate Limit:** Telegram saniyede 30 mesaj limitine sahip. [RateLimiter](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/core/event_bus.py#34-63) sınıfı bunu koruyor.

### Windows API

| API | Ne İçin | Dosya |
|-----|---------|-------|
| `win32evtlog.EvtSubscribe` | Security Event Log'a abone ol (async callback) | [windows_event.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/listeners/windows_event.py) |
| `win32evtlog.EvtRender` | Event XML çıktısını al | [windows_event.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/listeners/windows_event.py) |
| `WMI Win32_VolumeChangeEvent` | USB disk takma algılama | [usb_monitor.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/listeners/usb_monitor.py) |
| `WMI Win32_LogicalDisk` | Disk bilgisi (etiket, boyut) | [usb_monitor.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/listeners/usb_monitor.py) |

### Yerel API

| API | Ne İçin | Dosya |
|-----|---------|-------|
| `cv2.VideoCapture(0, cv2.CAP_DSHOW)` | Kamera açma (DirectShow) | [camera.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/actions/camera.py) |
| `mss.mss().grab()` | Ekran görüntüsü alma | [screenshot.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/actions/screenshot.py) |
| `psutil.cpu_percent()` | CPU kullanım yüzdesi | [bot_commands.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/telegram_bot/bot_commands.py) |
| `psutil.virtual_memory()` | RAM bilgisi | [bot_commands.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/telegram_bot/bot_commands.py) |
| `psutil.disk_usage()` | Disk doluluk bilgisi | [bot_commands.py](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/telegram_bot/bot_commands.py) |

---

## 🧵 Thread Mimarisi

```
Ana Thread (kanije.py → app.py → _event_loop)
│
├── [Daemon] WinEventListener thread
│   └── win32evtlog.EvtSubscribe callback
│       └── Event algılanınca → EventBus.put()
│
├── [Daemon] USBMonitor thread
│   └── WMI watcher (2sn timeout döngü)
│       └── USB algılanınca → EventBus.put()
│
├── [Daemon] BotCommands thread
│   └── Telegram getUpdates polling (2sn aralık)
│       └── Komut gelince → işle, yanıtla
│
├── [Daemon] Heartbeat thread
│   └── X saatte bir heartbeat mesajı gönder
│
└── [Daemon] TrayIcon thread
    └── pystray ana döngüsü (menü, ikon)
```

Tüm thread'ler `daemon=True` — ana süreç kapanınca hepsi otomatik sonlanır.

---

## 📱 TELEGRAM BOT KURULUMU — Sıfırdan Adım Adım

### Adım 1: Bot Oluştur

```
1. Telefonunda veya bilgisayarında Telegram'ı aç

2. Üst arama çubuğuna @BotFather yaz ve tıkla
   (Mavi tik işareti olan resmi Telegram botu)

3. /start yaz veya START butonuna bas

4. /newbot yaz

5. BotFather sana soracak:
   "Alright, a new bot. How are we going to call it?"
   → Bot'un GÖRÜNEN ADINI yaz: Kanije Kalesi

6. BotFather tekrar soracak:
   "Good. Now let's choose a username for your bot."
   → Bot'un KULLANICI ADINI yaz: kanije_kalesi_bot
     (sonunda _bot olmalı, benzersiz olmalı, Türkçe karakter yok)

7. BotFather sana bir TOKEN verecek:
   ╔═══════════════════════════════════════════╗
   ║ Done! Congratulations...                  ║
   ║ Use this token to access the HTTP API:    ║
   ║                                           ║
   ║ 7123456789:AAHxyz-YourSecretTokenHere     ║
   ║                                           ║
   ╚═══════════════════════════════════════════╝

   → BU TOKEN'I KOPYALA (tamamını, : dahil)
```

### Adım 2: Chat ID'ni Öğren

```
1. Telegram'da arama çubuğuna kendi bot adını yaz: @kanije_kalesi_bot

2. Bot'u aç ve /start yaz

3. Şimdi tarayıcıyı aç ve şu adresi gir
   (TOKEN kısmına kendi token'ını koy):

   https://api.telegram.org/bot7123456789:AAHxyz-YourSecretTokenHere/getUpdates

4. Tarayıcıda JSON çıktısı göreceksin:
   {
     "ok": true,
     "result": [{
       "message": {
         "chat": {
           "id": 987654321,     ← BU SENİN CHAT ID'N
           "first_name": "Yasin",
           "type": "private"
         }
       }
     }]
   }

   → 987654321 sayısını kopyala (bu senin chat_id'n)
```

### Adım 3: config.yaml'a Yaz

```yaml
# config.yaml dosyasını aç ve şunları doldur:
telegram:
  bot_token: "7123456789:AAHxyz-YourSecretTokenHere"    # ← Adım 1'deki token
  chat_id: "987654321"                                   # ← Adım 2'deki chat ID
```

### Adım 4: Test Et

```powershell
cd c:\Users\Yasin\Downloads\Kanije-Kalesi\app
python kanije.py test
```

Telegram'ına mesaj geldi mi? Geldiyse **her şey hazır** 🎉

---

## 🚀 Kurulum ve Çalıştırma

### 1. Bağımlılıkları Kur

```powershell
cd c:\Users\Yasin\Downloads\Kanije-Kalesi\app
pip install -r requirements.txt
```

### 2. Config'i Düzenle

```powershell
notepad config.yaml
# bot_token ve chat_id alanlarını doldur (yukarıdaki Telegram rehberine bak)
```

### 3. Durum Kontrolü

```powershell
python kanije.py status
# Tüm bağımlılıklar ✅ olmalı
```

### 4. Telegram Testi

```powershell
python kanije.py test
# Telegram'ına test mesajı gelmeli
```

### 5. Uygulamayı Başlat

```powershell
python kanije.py start
# Ctrl+C ile durdur
```

---

## 📋 Telegram Bot Komutları

| Komut | Açıklama | Yanıt |
|-------|----------|-------|
| `/status` | Sistem durumu | CPU, RAM, Disk, uptime, son olay |
| `/photo` | Anlık kamera fotoğrafı | Webcam'den çekilen fotoğraf |
| `/screenshot` | Anlık ekran görüntüsü | Masaüstünün ekran görüntüsü |
| `/ping` | Canlılık kontrolü | "Kanije ayakta!" + uptime |
| `/events` | Son 10 olay | Olay listesi |
| `/config` | Aktif yapılandırma | Config (token gizli) |
| `/help` | Komut listesi | Tüm komutlar |

---

## ⚙️ Yapılandırılabilir Özellikler

Her özellik [config.yaml](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/config.yaml) üzerinden açılıp kapatılabilir:

| Özellik | Config Yolu | Varsayılan |
|---------|-------------|------------|
| Başarılı giriş bildirimi | `triggers.login_success.enabled` | `true` |
| Başarısız girişte kamera | `triggers.login_failed.capture_camera` | `true` |
| Ekran kilidi bildirimi | `triggers.screen_unlock.enabled` | `true` |
| USB algılama | `triggers.usb_inserted.enabled` | `true` |
| Heartbeat | `heartbeat.enabled` | `true` |
| System Tray | `tray.enabled` | `true` |
| Konsol çıktısı | `logging.console_output` | `true` |
| Fotoğraf otomatik silme | `security.delete_photos_after_send` | `true` |
| Flood koruması | `security.max_events_per_minute` | `10` |

---

## 🔒 Güvenlik Notları

- **Bot Token** hiçbir zaman log dosyasına yazılmaz ([dump_safe](file:///c:/Users/Yasin/Downloads/Kanije-Kalesi/app/core/config_manager.py#213-224) metodu gizler)
- **Komut güvenliği:** Yalnızca `chat_id`'de tanımlı kullanıcının komutlarına yanıt verilir
- **Yetkisiz erişim** sessizce loglanır ama yanıtlanmaz
- **Fotoğraflar** gönderildikten sonra geçici dizinden silinir
- **Rate limiter** dakikada 10'dan fazla olayı engeller (Telegram ban koruması)

---

*Oluşturulma: 2026-03-22 · Python 3.11+ · Windows 10/11*
