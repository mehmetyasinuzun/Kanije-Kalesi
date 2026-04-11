# 🏰 Kanije Kalesi — Bilgisayar Güvenlik İzleme Yazılımı
## Mimari Tasarım · Algoritma · Gereksinimler · Uygulama Planı

> **Amaç:** Her oturum açılışında Telegram'a bildirim gönderen, yanlış şifre girildiğinde kameradan fotoğraf çeken, ekran görüntüsü alan, **Telegram üzerinden iki yönlü komut desteği** sunan ve tüm bunları yapılandırılabilir şekilde çalıştıran hafif, kararlı bir güvenlik izleme yazılımı.

> **İsim Kökeni:** 1601'de Tiryaki Hasan Paşa, 9.000 askerle 100.000 kişilik Haçlı ordusunu 73 gün boyunca Kanije Kalesi'nde durdurdu. Önemli olan sayı değil — hazırlık, katmanlı savunma ve sınırlı kaynağı doğru kullanmaktı. Bu uygulama da aynı felsefeyle: hafif ama etkili.

---

## 📐 ÜST DÜZEY MİMARİ

```
┌───────────────────────────────────────────────────────────────────┐
│                       🏰 Kanije Kalesi                            │
│                                                                    │
│  ┌─────────────┐   ┌──────────────┐   ┌───────────────────────┐  │
│  │  Event       │   │  Action      │   │  Telegram             │  │
│  │  Listener    │──▶│  Engine      │──▶│  ├─ Notifier (→ user) │  │
│  │              │   │              │   │  └─ Bot Cmd  (← user) │  │
│  │ Win: EventLog│   │ Kamera çek   │   │                       │  │
│  │ USB: WMI     │   │ Ekran al     │   │  /status /photo       │  │
│  │              │   │ Log yaz      │   │  /screenshot /ping    │  │
│  └─────────────┘   └──────────────┘   └───────────────────────┘  │
│         │                  │                   │                   │
│  ┌──────▼──────────────────▼───────────────────▼────────────────┐ │
│  │                    Config Manager                             │ │
│  │          config.yaml — Tüm ayarlar tek dosyada                │ │
│  └──────────────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │  System Tray İkonu (🟢 aktif / 🔴 hata) + CLI (start/test) │ │
│  └──────────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
```

---

## BÖLÜM 1 — PROGRAMLAMA DİLİ SEÇİMİ

### Neden Python?

| Kriter | Python | Go | C# | Rust |
|--------|--------|----|----|------|
| Kamera erişimi | ✅ `opencv-python` (1 satır) | ⚠️ CGo gerekir | ✅ AForge | ⚠️ Binding |
| Telegram API | ✅ `python-telegram-bot` | ✅ `telebot` | ⚠️ Az kütüphane | ⚠️ |
| Windows Event Log | ✅ `pywin32` / `wmi` | ⚠️ Syscall | ✅ Native | ⚠️ |
| Linux PAM/Journal | ✅ `systemd-python` | ✅ Native | ❌ | ✅ |
| Servis olarak çalışma | ✅ `pyinstaller` + NSSM | ✅ Native binary | ✅ | ✅ |
| Kurulum kolaylığı | ★★★★★ | ★★★★☆ | ★★★☆☆ | ★★☆☆☆ |
| RAM kullanımı | ~25-40 MB | ~8-15 MB | ~30-50 MB | ~5-10 MB |
| Geliştirme hızı | ★★★★★ | ★★★☆☆ | ★★★☆☆ | ★★☆☆☆ |

**Karar:** Python 3.11+

**Gerekçe:**
```
1. Kamera (OpenCV) + Telegram (Bot API) + Event Log (pywin32) → Hepsi olgun, stabil
2. Tek codebase ile Windows + Linux çalışır (platform modülleri ayrılır)
3. PyInstaller ile tek .exe haline getirilir → Python kurulu olmayan PC'de çalışır
4. RAM: ~30 MB — arka planda sessiz çalışır, CPU kullanımı %0-1
5. Senior seviyesinde hata yönetimi ve loglama Python'da en kolay
```

> [!NOTE]
> **Performans endişesi?** Bu program sürekli hesaplama yapmıyor — sadece olay dinliyor (event-driven). Olay gelince kamera çekip Telegram'a göndermek toplam ~2 saniye. Arada CPU kullanımı %0.

---

## BÖLÜM 2 — ÖZELLİK MATRİSİ

### 2.1 — Çekirdek Özellikler (Tek Yönlü: Uygulama → Kullanıcı)

| # | Özellik | Tetikleyici | Çıktı | Ayarlanabilir mi? |
|---|---------|-------------|-------|-------------------|
| 1 | **Başarılı Giriş Bildirimi** | Windows/Linux oturum açılışı | Telegram metin mesajı | ✅ Açık/Kapalı |
| 2 | **Başarısız Giriş — Kamera** | Yanlış şifre girilmesi | Telegram'a fotoğraf | ✅ Açık/Kapalı |
| 3 | **Başarısız Giriş — Ekran Görüntüsü** | Yanlış şifre girilmesi | Telegram'a ekran görüntüsü | ✅ Açık/Kapalı |
| 4 | **Kilitleme/Kilit Açma** | Ekran kilidi açıldığında | Telegram bildirim | ✅ Açık/Kapalı |
| 5 | **Periyodik Durum Raporu** | Her X saatte bir (heartbeat) | "Sistem çalışıyor" mesajı | ✅ Saat/Kapalı |
| 6 | **USB Takma Algılama** | Yeni USB cihazı takıldığında | Telegram uyarı | ✅ Açık/Kapalı |
| 7 | **Sistem Başlatma** | Bilgisayar açıldığında | Telegram bildirim | ✅ Açık/Kapalı |

### 2.2 — Telegram İki Yönlü Komut Desteği (Kullanıcı → Uygulama)

| # | Komut | Ne Yapar | Çıktı |
|---|-------|----------|-------|
| 1 | `/status` | Sistem durumu sorgula | CPU, RAM, Disk, uptime, son olay |
| 2 | `/photo` | Kameradan anlık fotoğraf çek | Webcam fotoğrafı |
| 3 | `/screenshot` | Ekranın anlık görüntüsünü al | Masaüstü PNG |
| 4 | `/ping` | Canlılık kontrolü | "Kanije ayakta 🏰" + uptime |
| 5 | `/events` | Son 10 olayı listele | Olay listesi |
| 6 | `/config` | Aktif yapılandırmayı göster | Config (token gizlenmiş) |
| 7 | `/help` | Komut listesi | Tüm komutlar |

> **Güvenlik:** Yalnızca `config.yaml`'daki `chat_id` ile eşleşen kullanıcıdan gelen komutlar işlenir. Başka kullanıcının komutu sessizce loglanır, yanıtlanmaz.

### 2.3 — Best Practice Önerileri

```
YAPMAMASI GEREKENLER:
  ❌ Sürekli ekran görüntüsü almak (gereksiz yük + gizlilik ihlali)
  ❌ Her tuş basışını kaydetmek (keylogger → yasal sorun)
  ❌ Video kaydetmek (bant genişliği + depolama)
  ❌ Kamerayı sürekli açık tutmak (LED yanar, gizlilik)
  ❌ Bildirim spam'i yapmak (rate limiter ile engellenir)

YAPMASI GEREKENLER:
  ✅ YALNIZCA olay tetiklendiğinde çalış (event-driven)
  ✅ Kamera sadece yanlış şifrede çeksin (anlık, <1 sn)
  ✅ Ekran görüntüsü başarısız girişte alsın
  ✅ Tüm özellikler config'den açılıp kapatılabilsin
  ✅ Log dosyası tutsun (ama hassas veri yazmasın)
  ✅ Telegram bot token'ı ortam değişkeninde veya şifreli config'de sakla
  ✅ Telegram'dan komut alabilsin (iki yönlü iletişim)
  ✅ Yalnızca yetkili chat_id'ye yanıt versin
  ✅ System tray ikonu ile durumu görsel göstersin
```

---

## BÖLÜM 3 — DOSYA YAPISI

```
Kanije-Kalesi/
├── README.md                        # Proje açıklaması + Kanije Savunması hikayesi
├── LICENSE                          # MIT Lisansı
├── .gitignore                       # config.yaml + __pycache__ + log dışlama
│
├── app/                             # ← Uygulama kodu
│   ├── kanije.py                    # CLI giriş noktası (start / test / status)
│   ├── config.yaml                  # Gerçek ayarlar (GIT'E GİRMEZ — .gitignore)
│   ├── config_sample.yaml           # Token'sız örnek (GIT'E GİRER)
│   ├── requirements.txt             # Python bağımlılıkları
│   │
│   ├── core/
│   │   ├── app.py                   # Ana orkestratör — tüm modülleri yönetir
│   │   ├── config_manager.py        # YAML okuma, ENV çözümleme, doğrulama
│   │   ├── event_bus.py             # Thread-safe olay kuyruğu + rate limiter
│   │   └── logger.py                # RotatingFileHandler loglama
│   │
│   ├── listeners/
│   │   ├── windows_event.py         # Windows Security Event Log (4624/4625/4800/4801)
│   │   └── usb_monitor.py           # WMI ile USB algılama
│   │
│   ├── actions/
│   │   ├── camera.py                # OpenCV kamera çekim (thread-safe)
│   │   └── screenshot.py            # mss ekran görüntüsü
│   │
│   ├── telegram_bot/
│   │   ├── notifier.py              # Tek yönlü: mesaj + fotoğraf gönderme
│   │   ├── bot_commands.py          # İki yönlü: /status /photo /screenshot /ping
│   │   └── formatter.py             # Emoji'li mesaj şablonları
│   │
│   ├── tray/
│   │   └── tray_icon.py             # System tray (yeşil/kırmızı ikon + menü)
│   │
│   └── tests/
│       ├── test_config.py
│       ├── test_event_bus.py
│       ├── test_camera.py
│       └── test_formatter.py
│
├── docs/                            # ← Güvenlik rehberleri
│   ├── WINDOWS11_HARDENING_KALE.md
│   ├── WINDOWS10_HARDENING_KALE.md
│   ├── LINUX_HARDENING_KALE.md
│   ├── DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md
│   └── ...
│
└── HAYAT_KURTARAN_YAZILIMLAR.md     # Format sonrası yazılım rehberi
```

---

## BÖLÜM 4 — YAPILANDIRMA (config.yaml)

```yaml
# GhostGuard Yapılandırma Dosyası
# Her ayar açıklamalıdır

# ── Telegram ──
telegram:
  bot_token: "ENV:GHOSTGUARD_BOT_TOKEN"   # Ortam değişkeninden oku (güvenli)
  chat_id: "123456789"                     # Telegram chat/grup ID
  send_timeout: 10                         # Gönderim zaman aşımı (saniye)
  retry_count: 3                           # Başarısızlıkta tekrar deneme
  retry_delay: 5                           # Denemeler arası bekleme (saniye)

# ── Tetikleyiciler ──
triggers:
  login_success:
    enabled: true                          # Başarılı girişte bildirim
    message: "✅ Giriş yapıldı"
    include_username: true                 # Kullanıcı adını göster
    include_ip: true                       # IP adresini göster (RDP için)
    include_timestamp: true

  login_failed:
    enabled: true                          # Başarısız girişte bildirim
    message: "🚨 YANLIŞ ŞİFRE"
    capture_camera: true                   # Kameradan fotoğraf çek
    capture_screenshot: false              # Ekran görüntüsü al (kilit ekranında sınırlı)
    max_photos_per_minute: 3               # Spam koruması

  screen_unlock:
    enabled: true                          # Kilit açıldığında bildirim
    message: "🔓 Ekran kilidi açıldı"
    capture_camera: false

  system_boot:
    enabled: true                          # Sistem açılışında bildirim
    message: "🖥️ Sistem başlatıldı"

  usb_inserted:
    enabled: true                          # USB takıldığında bildirim
    message: "🔌 Yeni USB cihazı algılandı"

# ── Kamera ──
camera:
  device_index: 0                          # Kamera indeksi (0 = varsayılan)
  resolution: [640, 480]                   # Çözünürlük [genişlik, yükseklik]
  warmup_frames: 5                         # Kamera ısınma karesi (ilk kareler karanlık olur)
  save_local: false                        # Yerel kopyayı sakla mı?
  local_path: "./captures/"                # Yerel kayıt dizini
  jpeg_quality: 85                         # JPEG kalitesi (1-100)

# ── Ekran Görüntüsü ──
screenshot:
  save_local: false
  local_path: "./captures/"
  jpeg_quality: 70

# ── Heartbeat (Canlılık Kontrolü) ──
heartbeat:
  enabled: true
  interval_hours: 6                        # Her 6 saatte bir "çalışıyorum" mesajı
  message: "💚 GhostGuard aktif"
  include_uptime: true                     # Sistemin açık kalma süresini göster
  include_disk_usage: true                 # Disk doluluk bilgisi

# ── Loglama ──
logging:
  level: "INFO"                            # DEBUG, INFO, WARNING, ERROR
  file: "./ghostguard.log"
  max_size_mb: 10                          # Log dosyası max boyut
  backup_count: 3                          # Eski log dosyası sayısı
  console_output: false                    # Konsola da yazdır mı?

# ── Güvenlik ──
security:
  encrypt_token: false                     # Bot token'ı şifreli sakla (gelişmiş)
  delete_photos_after_send: true           # Gönderilen fotoğrafı sil
  max_events_per_minute: 10               # Flood koruması
```

---

## BÖLÜM 5 — ALGORİTMA VE AKIŞ DİYAGRAMLARI

### 5.1 — Ana Program Akışı

```
ghostguard.py başlatıldı
        │
        ▼
┌─────────────────────┐
│ config.yaml oku     │
│ Ayarları doğrula    │
│ Logger başlat       │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ OS algıla:          │
│ Windows → Win       │
│ Linux → Linux       │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ Telegram bağlantısı │
│ doğrula (test msg)  │
│ Başarısız → log +   │
│ retry               │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ Kamera doğrula      │
│ (enabled ise)       │
│ Kamera yoksa →      │
│ özelliği devre dışı │
│ bırak, log yaz      │
└─────────┬───────────┘
          │
          ▼
┌────────────────────────────────────────┐
│         ANA DÖNGÜ (Event Loop)         │
│                                         │
│  Event Listener (OS'a göre):           │
│    Windows: Event Log subscribe        │
│    Linux: journald --follow            │
│                                         │
│  Heartbeat Timer: Her X saatte tetikle │
│  USB Monitor: WMI / udev listener     │
│                                         │
│  OLAY GELDİ:                           │
│    1. Flood kontrolü (max/dk)          │
│    2. Olay tipini belirle              │
│    3. Config'e göre aksiyonları çalıştır│
│    4. Telegram'a gönder                 │
│    5. Log yaz                           │
│    6. Geçici dosyaları temizle          │
└────────────────────────────────────────┘
```

### 5.2 — Başarısız Giriş Akışı (En Kritik Senaryo)

```
Windows Event Log → Event ID 4625 (Logon Failure) algılandı
        │
        ▼
Flood kontrolü: Son 1 dakikada 3'ten fazla tetiklenme var mı?
        │
   EVET ▼ HAYIR
   Atla  │
         ▼
Config kontrolü: login_failed.capture_camera = true?
        │
   HAYIR▼ EVET
   ────▶ │
         ▼
┌─────────────────────────────┐
│ Kamera Çekim Süreci:       │
│ 1. cv2.VideoCapture(0)     │
│ 2. 5 kare oku (warmup)     │
│ 3. 6. kareyi kaydet        │
│ 4. Kamerayı serbest bırak  │
│ 5. Toplam süre: <1 saniye  │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐
│ Telegram'a gönder:          │
│ 📸 Fotoğraf                │
│ 📝 "🚨 YANLIŞ ŞİFRE"     │
│ ⏰ Tarih/Saat              │
│ 👤 Denenen kullanıcı adı   │
│ 🖥️ Bilgisayar adı         │
└──────────┬──────────────────┘
           │
           ▼
Başarılı → Log yaz → Fotoğrafı sil (config'e göre)
Başarısız → 3 kez tekrar dene → Hâlâ başarısız → Log yaz, devam et
```

### 5.3 — Telegram Mesaj Formatı

```
━━━━━━━━━━━━━━━━━━━
🚨 BAŞARISIZ GİRİŞ DENEMESİ
━━━━━━━━━━━━━━━━━━━

🖥️ Bilgisayar: YASIN-PC
👤 Denenen Kullanıcı: Administrator
⏰ Zaman: 2026-03-22 16:45:32
🌐 IP: 192.168.1.1 (yerel)
🔗 Giriş Tipi: Yerel (Konsol)
📸 Kamera görüntüsü ektedir.

─────────────────────
🏰 Kanije Kalesi v1.0.0 · Aktif: 4s 23dk
```

```
━━━━━━━━━━━━━━━━━━━
✅ BAŞARILI GİRİŞ
━━━━━━━━━━━━━━━━━━━

🖥️ Bilgisayar: YASIN-PC
👤 Kullanıcı: yasin
⏰ Zaman: 2026-03-22 08:30:15
🌐 IP: Yerel oturum
🔗 Giriş Tipi: Kilit Açma

─────────────────────
🏰 Kanije Kalesi v1.0.0
```

---

## BÖLÜM 6 — MODÜL DETAYLARI

### 6.1 — Event Listener (Windows)

```python
# listeners/windows_listener.py — Sözde Kod (Pseudocode)

"""
Windows Event Log'dan oturum olaylarını dinler.
Kullanılan Event ID'ler:

  4624 → Başarılı oturum açma (Logon)
  4625 → Başarısız oturum açma (Logon Failure)
  4634 → Oturum kapama (Logoff)
  4800 → İş istasyonu kilitlendi
  4801 → İş istasyonunun kilidi açıldı

Yöntem: win32evtlog.EvtSubscribe() ile asenkron dinleme
  → Polling DEĞİL (CPU kullanmaz)
  → Windows Event Log servisi callback çağırır
"""

import win32evtlog
import win32evtlog_constants

EVENT_MAP = {
    4624: "login_success",
    4625: "login_failed",
    4800: "screen_lock",
    4801: "screen_unlock",
}

class WindowsListener:
    def __init__(self, callback):
        self.callback = callback
    
    def start(self):
        # Security log'a subscribe ol
        # XML query ile sadece istediğimiz event ID'leri filtrele
        query = """
        <QueryList>
          <Query Id="0" Path="Security">
            <Select Path="Security">
              *[System[(EventID=4624 or EventID=4625 or
                         EventID=4800 or EventID=4801)]]
            </Select>
          </Query>
        </QueryList>
        """
        # Asenkron subscribe — CPU kullanmaz, olay gelince callback çağırır
        handle = win32evtlog.EvtSubscribe(
            "Security",
            win32evtlog.EvtSubscribeToFutureEvents,
            Query=query,
            Callback=self._on_event
        )
    
    def _on_event(self, reason, context, event):
        # Event XML'den bilgileri çıkar
        event_id = extract_event_id(event)
        username = extract_username(event)
        ip_address = extract_ip(event)
        timestamp = extract_timestamp(event)
        
        event_type = EVENT_MAP.get(event_id, "unknown")
        self.callback(event_type, {
            "username": username,
            "ip": ip_address,
            "timestamp": timestamp,
            "event_id": event_id
        })
```

### 6.2 — Event Listener (Linux)

```python
# listeners/linux_listener.py — Sözde Kod

"""
Linux'ta oturum olayları:
  1. systemd-journald → journalctl --follow ile dinle
  2. PAM modülleri → /var/log/auth.log izle
  3. utmp/wtmp → login/logout kayıtları

Yöntem: systemd journal + inotify (dosya değişikliği izleme)

İlgili log satırları:
  "session opened for user"  → Başarılı giriş
  "authentication failure"   → Başarısız giriş
  "session closed for user"  → Oturum kapandı
"""

import subprocess
import re

class LinuxListener:
    def __init__(self, callback):
        self.callback = callback
    
    def start(self):
        # journalctl --follow ile gerçek zamanlı izle
        proc = subprocess.Popen(
            ["journalctl", "--follow", "--no-pager",
             "-u", "sshd", "-u", "gdm", "-u", "lightdm",
             "-u", "systemd-logind", "--output=json"],
            stdout=subprocess.PIPE,
            text=True
        )
        
        for line in proc.stdout:
            event = self._parse_journal_line(line)
            if event:
                self.callback(event["type"], event["data"])
    
    def _parse_journal_line(self, line):
        if "authentication failure" in line:
            return {"type": "login_failed", "data": {...}}
        elif "session opened" in line:
            return {"type": "login_success", "data": {...}}
        elif "session closed" in line:
            return {"type": "screen_lock", "data": {...}}
        return None
```

### 6.3 — Kamera Modülü

```python
# actions/camera.py — Sözde Kod

"""
Kamera çekim süreci:
  1. cv2.VideoCapture(device_index) ile kamerayı aç
  2. warmup_frames kadar kare oku ve at (kamera ısınması — ilk kareler karanlık)
  3. Son kareyi kaydet
  4. Kamerayı serbest bırak
  5. Toplam süre: <1 saniye

Dikkat:
  - Kamera LED'i anlık yanar ve söner
  - Başarısız olursa (kamera yok/meşgul) hata logla, çökmesine izin verme
  - Thread-safe olmalı (aynı anda 2 kez çekilmemeli)
"""

import cv2
import threading
import tempfile
import os

class Camera:
    def __init__(self, config):
        self.device_index = config.get("device_index", 0)
        self.resolution = config.get("resolution", [640, 480])
        self.warmup = config.get("warmup_frames", 5)
        self.quality = config.get("jpeg_quality", 85)
        self._lock = threading.Lock()  # Eşzamanlı çekim engeli
    
    def capture(self) -> str | None:
        """Fotoğraf çeker, geçici dosya yolunu döndürür."""
        if not self._lock.acquire(blocking=False):
            return None  # Zaten çekim yapılıyor
        
        try:
            cap = cv2.VideoCapture(self.device_index)
            if not cap.isOpened():
                log.warning("Kamera açılamadı — atlanıyor")
                return None
            
            cap.set(cv2.CAP_PROP_FRAME_WIDTH, self.resolution[0])
            cap.set(cv2.CAP_PROP_FRAME_HEIGHT, self.resolution[1])
            
            # Warmup — ilk kareler genellikle siyah/karanlık
            for _ in range(self.warmup):
                cap.read()
            
            ret, frame = cap.read()
            cap.release()
            
            if not ret or frame is None:
                log.warning("Kamera karesi alınamadı")
                return None
            
            # Geçici dosyaya kaydet
            tmp = tempfile.NamedTemporaryFile(
                suffix=".jpg", delete=False, prefix="gg_cam_"
            )
            cv2.imwrite(tmp.name, frame,
                       [cv2.IMWRITE_JPEG_QUALITY, self.quality])
            return tmp.name
        
        except Exception as e:
            log.error(f"Kamera hatası: {e}")
            return None
        finally:
            self._lock.release()
```

### 6.4 — Ekran Görüntüsü Modülü

```python
# actions/screenshot.py — Sözde Kod

"""
Ekran görüntüsü alma:
  Windows: mss (pure Python, hızlı) veya Pillow ImageGrab
  Linux: mss (X11/Wayland)

Dikkat:
  - Kilit ekranında Windows ekran görüntüsü ALAMAZ (güvenlik kısıtı)
  - Bu nedenle ekran görüntüsü yalnızca başarılı girişte anlamlı
  - Başarısız girişte yalnızca kamera çalışır (kamera kilit ekranında da çeker)
"""

import mss
import tempfile

class Screenshot:
    def __init__(self, config):
        self.quality = config.get("jpeg_quality", 70)
    
    def capture(self) -> str | None:
        try:
            with mss.mss() as sct:
                monitor = sct.monitors[0]  # Tüm ekranlar
                img = sct.grab(monitor)
                
                tmp = tempfile.NamedTemporaryFile(
                    suffix=".png", delete=False, prefix="gg_ss_"
                )
                mss.tools.to_png(img.rgb, img.size, output=tmp.name)
                return tmp.name
        
        except Exception as e:
            log.error(f"Ekran görüntüsü hatası: {e}")
            return None
```

> [!WARNING]
> **Kilit ekranı kısıtı:** Windows kilit ekranında (`LogonUI.exe`) session 0'da çalışır; kullanıcı session'ı farklıdır. Bu yüzden kilit ekranında ekran görüntüsü alınamaz. Kamera ise donanıma doğrudan eriştiği için kilit ekranında da çalışır.

### 6.5 — Telegram Bildirim Modülü

```python
# notifiers/telegram_notifier.py — Sözde Kod

"""
Telegram Bot API:
  1. @BotFather'dan bot oluştur → Token al
  2. Bot'u gruba ekle veya direkt mesaj at
  3. chat_id: Mesaj atılacak kişi/grup ID

Retry mekanizması:
  Gönderilemezse → X saniye bekle → Tekrar dene → Max 3 deneme
  Hâlâ başarısız → Log yaz, programı DURDURMA
"""

import requests
import time

class TelegramNotifier:
    BASE_URL = "https://api.telegram.org/bot{token}"
    
    def __init__(self, config):
        token = config["bot_token"]
        if token.startswith("ENV:"):
            token = os.environ.get(token[4:], "")
        
        self.base = self.BASE_URL.format(token=token)
        self.chat_id = config["chat_id"]
        self.timeout = config.get("send_timeout", 10)
        self.retry = config.get("retry_count", 3)
        self.delay = config.get("retry_delay", 5)
    
    def send_message(self, text: str) -> bool:
        """Metin mesajı gönder."""
        for attempt in range(self.retry):
            try:
                resp = requests.post(
                    f"{self.base}/sendMessage",
                    json={"chat_id": self.chat_id, "text": text,
                          "parse_mode": "HTML"},
                    timeout=self.timeout
                )
                if resp.status_code == 200:
                    return True
                log.warning(f"Telegram {resp.status_code}: {resp.text}")
            except requests.RequestException as e:
                log.warning(f"Telegram bağlantı hatası (deneme {attempt+1}): {e}")
            
            if attempt < self.retry - 1:
                time.sleep(self.delay)
        
        log.error("Telegram mesaj gönderilemedi — tüm denemeler başarısız")
        return False
    
    def send_photo(self, photo_path: str, caption: str = "") -> bool:
        """Fotoğraf gönder."""
        for attempt in range(self.retry):
            try:
                with open(photo_path, "rb") as photo:
                    resp = requests.post(
                        f"{self.base}/sendPhoto",
                        data={"chat_id": self.chat_id, "caption": caption},
                        files={"photo": photo},
                        timeout=self.timeout
                    )
                if resp.status_code == 200:
                    return True
            except Exception as e:
                log.warning(f"Telegram fotoğraf hatası (deneme {attempt+1}): {e}")
            
            if attempt < self.retry - 1:
                time.sleep(self.delay)
        
        return False
    
    def verify_connection(self) -> bool:
        """Bot bağlantısını test et."""
        try:
            resp = requests.get(f"{self.base}/getMe", timeout=5)
            return resp.status_code == 200
        except:
            return False
```

### 6.6 — USB Monitör

```python
# actions/usb_monitor.py — Sözde Kod

"""
Windows: WMI ile USB ekleme olayını dinle
Linux: udev ile USB ekleme olayını dinle

Her iki platformda da asenkron (event-driven) — polling yok.
"""

# Windows:
import wmi

class USBMonitor:
    def __init__(self, callback):
        self.callback = callback
    
    def start_windows(self):
        c = wmi.WMI()
        watcher = c.Win32_USBHub.watch_for("creation")
        while True:
            usb = watcher()
            self.callback("usb_inserted", {
                "device": usb.Description,
                "device_id": usb.DeviceID,
                "timestamp": datetime.now()
            })

# Linux:
# pyudev kütüphanesi:
import pyudev

class USBMonitorLinux:
    def start_linux(self):
        context = pyudev.Context()
        monitor = pyudev.Monitor.from_netlink(context)
        monitor.filter_by(subsystem="usb")
        
        for device in iter(monitor.poll, None):
            if device.action == "add":
                self.callback("usb_inserted", {
                    "device": device.get("ID_MODEL", "Bilinmeyen"),
                    "vendor": device.get("ID_VENDOR", ""),
                    "timestamp": datetime.now()
                })
```

---

## BÖLÜM 7 — HATA YÖNETİMİ

### 7.1 — Asla Çökmeme Prensibi

```python
# ghostguard.py — Ana olay işleyicisi

def handle_event(event_type, data):
    """
    Senior prensibi: ASLA ÇÖKME.
    Her modül kendi hatasını yakalar.
    Telegram başarısız → log yaz, devam et.
    Kamera başarısız → log yaz, sadece mesaj gönder.
    Config hatalı → varsayılan değerlerle çalış.
    """
    try:
        # Flood kontrolü
        if rate_limiter.is_limited(event_type):
            log.debug(f"Flood limiti: {event_type} atlandı")
            return
        
        # Aksiyonları çalıştır
        config_key = event_type  # "login_failed", "login_success" vb.
        trigger_config = config["triggers"].get(config_key, {})
        
        if not trigger_config.get("enabled", False):
            return
        
        photo_path = None
        screenshot_path = None
        
        # Kamera
        if trigger_config.get("capture_camera", False):
            photo_path = camera.capture()  # Başarısızsa None döner
        
        # Ekran görüntüsü
        if trigger_config.get("capture_screenshot", False):
            screenshot_path = screenshot.capture()
        
        # Mesaj oluştur
        message = format_message(event_type, data, trigger_config)
        
        # Telegram'a gönder
        if photo_path:
            telegram.send_photo(photo_path, caption=message)
        elif screenshot_path:
            telegram.send_photo(screenshot_path, caption=message)
        else:
            telegram.send_message(message)
        
        # Temizlik
        if photo_path and config["security"]["delete_photos_after_send"]:
            os.unlink(photo_path)
        if screenshot_path and config["security"]["delete_photos_after_send"]:
            os.unlink(screenshot_path)
    
    except Exception as e:
        # EN SON SAVUNMA HATTI — burası bile hata verirse sadece logla
        log.error(f"Olay işleme hatası ({event_type}): {e}", exc_info=True)
        # Program DEVAM EDER — asla çökmez
```

### 7.2 — Rate Limiter (Flood Koruması)

```python
# Birisi 100 kez yanlış şifre girerse 100 fotoğraf çekip
# 100 Telegram mesajı göndermek istemeyiz

from collections import defaultdict
import time

class RateLimiter:
    def __init__(self, max_per_minute=10):
        self.max = max_per_minute
        self.events = defaultdict(list)
    
    def is_limited(self, event_type: str) -> bool:
        now = time.time()
        # Son 60 saniyedeki olayları filtrele
        self.events[event_type] = [
            t for t in self.events[event_type] if now - t < 60
        ]
        if len(self.events[event_type]) >= self.max:
            return True  # Limit aşıldı — bu olayı atla
        self.events[event_type].append(now)
        return False
```

---

## BÖLÜM 8 — SERVİS / DAEMON KURULUMU

### 8.1 — Windows: NSSM ile Servis

```powershell
# 1. PyInstaller ile .exe oluştur
pip install pyinstaller
pyinstaller --onefile --noconsole --name GhostGuard ghostguard.py

# 2. NSSM (Non-Sucking Service Manager) indir
# https://nssm.cc/download → nssm.exe

# 3. Windows servisi olarak kur
nssm install GhostGuard "C:\GhostGuard\GhostGuard.exe"
nssm set GhostGuard AppDirectory "C:\GhostGuard"
nssm set GhostGuard DisplayName "GhostGuard Security Monitor"
nssm set GhostGuard Description "Bilgisayar güvenlik izleme ve Telegram bildirim"
nssm set GhostGuard Start SERVICE_AUTO_START
nssm set GhostGuard ObjectName LocalSystem
nssm set GhostGuard AppStdout "C:\GhostGuard\service.log"
nssm set GhostGuard AppStderr "C:\GhostGuard\service_error.log"

# 4. Servisi başlat
nssm start GhostGuard

# 5. Doğrulama
Get-Service GhostGuard | Select-Object Status, StartType
# Status: Running, StartType: Automatic
```

> [!IMPORTANT]
> **Neden NSSM?** Windows Task Scheduler'dan farklı olarak NSSM:
> - Çökmede otomatik yeniden başlatır
> - Servis olarak çalışır → oturum açılmasa bile aktif
> - Çıkış kodlarını izler
> - Stdout/stderr yönlendirme

### 8.2 — Linux: systemd Servisi

```bash
sudo tee /etc/systemd/system/ghostguard.service <<'EOF'
[Unit]
Description=GhostGuard Security Monitor
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/ghostguard
ExecStart=/usr/bin/python3 /opt/ghostguard/ghostguard.py
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
Environment="GHOSTGUARD_BOT_TOKEN=BURAYA_TOKEN"

# Güvenlik sertleştirme
ProtectSystem=strict
ReadWritePaths=/opt/ghostguard /tmp
PrivateTmp=yes
NoNewPrivileges=yes

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now ghostguard

# Doğrulama:
sudo systemctl status ghostguard
sudo journalctl -u ghostguard -f    # Canlı log izle
```

---

## BÖLÜM 9 — TELEGRAM BOT KURULUMU (Adım Adım)

```
1. Telegram'ı aç → @BotFather'ı bul → /start

2. /newbot komutu gönder
   → Bot adı: GhostGuard Security
   → Bot kullanıcı adı: ghostguard_yasin_bot (benzersiz olmalı)
   → Bot token verilir: 123456789:ABCdefGHIjklMNO... ← BUNU KAYDET

3. Chat ID'ni öğren:
   → Botu Telegram'da bul → /start gönder
   → Tarayıcıda aç: https://api.telegram.org/bot{TOKEN}/getUpdates
   → JSON çıktısında "chat":{"id": 123456789} → Bu senin chat_id

4. Test:
   curl -s -X POST \
     "https://api.telegram.org/bot{TOKEN}/sendMessage" \
     -d "chat_id={CHAT_ID}" \
     -d "text=🔔 GhostGuard test mesajı"
   → Telegram'ına mesaj geldiyse çalışıyor ✅

5. Güvenlik:
   → Token'ı ASLA kod içine yapıştırma
   → Ortam değişkeni kullan: GHOSTGUARD_BOT_TOKEN
   → Windows: setx GHOSTGUARD_BOT_TOKEN "token_buraya"
   → Linux: export GHOSTGUARD_BOT_TOKEN="token_buraya" → .bashrc'ye ekle
```

---

## BÖLÜM 10 — BAĞIMLILIKLAR

### 10.1 — requirements.txt

```
opencv-python-headless==4.9.0.80   # Kamera (headless = GUI yok, hafif)
mss==9.0.1                          # Ekran görüntüsü (pure Python, hızlı)
requests==2.31.0                    # HTTP (Telegram API)
pyyaml==6.0.1                       # Config dosyası okuma
wmi==1.5.1; sys_platform=="win32"   # USB izleme (Windows)
pywin32==306; sys_platform=="win32" # Event Log (Windows)
pyudev==0.24.1; sys_platform=="linux" # USB izleme (Linux)
```

### 10.2 — Kurulum

```bash
# Sanal ortam oluştur (izolasyon)
python -m venv .venv

# Windows:
.venv\Scripts\activate
pip install -r requirements.txt

# Linux:
source .venv/bin/activate
pip install -r requirements.txt
```

### 10.3 — PyInstaller ile Tek .exe

```bash
# Windows — tek dosya, konsol penceresi yok:
pyinstaller --onefile --noconsole \
    --add-data "config.yaml;." \
    --name GhostGuard \
    --icon ghostguard.ico \
    ghostguard.py

# Çıktı: dist/GhostGuard.exe (~25-35 MB)
# config.yaml'ı .exe'nin yanına koy → çalıştır
```

---

## BÖLÜM 11 — TEST STRATEJİSİ

```python
# tests/test_camera.py
def test_camera_no_device():
    """Kamera yoksa None döndürmeli, çökmemeli."""
    cam = Camera({"device_index": 99})  # Olmayan kamera
    result = cam.capture()
    assert result is None  # Hata yok, None döndü

def test_camera_concurrent():
    """Aynı anda 2 çekim yapılırsa ikincisi atlanmalı."""
    cam = Camera(config)
    # Thread 1 çekerken Thread 2 None döndürmeli

# tests/test_telegram.py
def test_telegram_retry():
    """Bağlantı hatası → 3 kez tekrar denemeli."""
    # Mock requests → ilk 2 başarısız, 3. başarılı
    
def test_telegram_format():
    """Mesaj formatı doğru mu?"""
    msg = format_message("login_failed", {"username": "test"})
    assert "🚨" in msg
    assert "test" in msg

# tests/test_listener.py
def test_windows_event_parse():
    """Windows Event XML doğru parse ediliyor mu?"""
    
def test_rate_limiter():
    """60 saniyede 10'dan fazla olay → limit aktif."""
    rl = RateLimiter(max_per_minute=3)
    assert not rl.is_limited("test")  # 1
    assert not rl.is_limited("test")  # 2
    assert not rl.is_limited("test")  # 3
    assert rl.is_limited("test")      # 4 → LİMİT!
```

```bash
# Testleri çalıştır:
python -m pytest tests/ -v --tb=short

# Kapsam raporu:
python -m pytest tests/ --cov=. --cov-report=term-missing
```

---

## BÖLÜM 12 — KAYNAK KULLANIMI

```
İDLE (Olay beklerken):
  CPU:  %0.0 — Event-driven; polling yok
  RAM:  ~25-30 MB (Python interpreter + modüller)
  Disk: 0 MB/s — Hiçbir şey yazılmıyor
  Ağ:   0 KB/s — Hiçbir şey gönderilmiyor

OLAY TETİKLENDİĞİNDE (anlık, <2 sn):
  CPU:  %2-5 (kamera çekimi + JPEG encode)
  RAM:  +5-10 MB (kamera frame buffer — hemen serbest bırakılır)
  Disk: ~200 KB (geçici JPEG — hemen siliniyor)
  Ağ:   ~200-500 KB (fotoğraf Telegram'a upload)

KARŞILAŞTIRMA:
  Chrome (1 sekme):     ~150-300 MB RAM
  Discord:              ~100-200 MB RAM
  GhostGuard:           ~25-30 MB RAM ← Fark edilmez
```

---

## BÖLÜM 13 — GÜVENLİK DEĞERLENDİRMESİ

### Bot Token Güvenliği

```
❌ KÖTÜ: Token kod içinde veya config'de düz metin
   config.yaml → bot_token: "123456:ABCdef..."
   → Repo'ya yüklenirse token açığa çıkar

✅ İYİ: Ortam değişkeni
   config.yaml → bot_token: "ENV:GHOSTGUARD_BOT_TOKEN"
   → Token yalnızca sistemde, repo'da değil

✅ DAHA İYİ: Şifreli config (Windows DPAPI / Linux keyring)
   crypto.py → DPAPI ile token'ı şifrele, çalışma anında çöz
   → Başka kullanıcı dosyayı kopyalasa bile çözemez
```

### Program Güvenliği

```
✅ Minimum yetki prensibi:
   → Yalnızca Event Log okuma + Kamera + Ağ erişimi
   → Dosya yazma yalnızca geçici dizin (/tmp)
   → Kayıt defteri değiştirmez, sistem dosyalarına dokunmaz

✅ Kendi kendine savunma:
   → NSSM çökmede yeniden başlatır
   → systemd Restart=always
   → Rate limiter → Telegram API ban'ı önler (Telegram: 30 msg/sn limit)

✅ İzlenebilirlik:
   → Her olay loglanır (ama şifre/token loglanmaz)
   → Log rotasyonu (10 MB → eski silinir)
```

---

## ÖZET — İMPLEMENTASYON DURUMU

```
✅ TAMAMLANDI:
   ✅ config.yaml okuyucu (ENV çözümleme + deep merge + doğrulama)
   ✅ Logger sistemi (RotatingFileHandler)
   ✅ Telegram bağlantısı + mesaj/fotoğraf gönderme
   ✅ Windows Event Listener (4624, 4625, 4800, 4801)
   ✅ Kamera çekim modülü (OpenCV, thread-safe)
   ✅ Ekran görüntüsü modülü (mss)
   ✅ Rate limiter (flood koruması)
   ✅ USB monitör (WMI)
   ✅ Heartbeat (periyodik durum raporu)
   ✅ Hata yönetimi katmanları (asla çökmeme)
   ✅ Telegram İki Yönlü Komut Desteği (/status, /photo, /screenshot, /ping...)
   ✅ System Tray İkonu (pystray)
   ✅ CLI giriş noktası (start / test / status)
   ✅ config_sample.yaml (token'sız örnek)
   ✅ .gitignore (config.yaml korunuyor)
   ✅ MIT Lisansı

🔜 SIRADA:
   □ PyInstaller ile .exe oluşturma
   □ NSSM Windows servisi kurulum scripti
   □ Unit testler
   □ RDP (Logon Type 10) ekstra uyarı
```

---

*Son güncelleme: 2026-03-23 · Kanije Kalesi v1.0.0 · Python 3.11+ · Windows 10/11*
*Lisans: MIT · Kaynak: https://github.com/mehmetyasinuzun/Kanije-Kalesi*
