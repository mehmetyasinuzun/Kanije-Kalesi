"""
Kanije Kalesi — Ana Uygulama Orkestratörü
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Tüm modülleri başlatır, event loop'u çalıştırır,
olayları işler ve graceful shutdown sağlar.

Bu dosya uygulamanın kalbidir.
"""

import signal
import socket
import threading
import time
import os
import json
from datetime import datetime
from typing import Optional
from collections import deque

from core.config_manager import ConfigManager
from core.logger import setup_logger, get_logger
from core.event_bus import EventBus, Event
from actions.camera import Camera
from actions.screenshot import Screenshot
from telegram_bot.notifier import TelegramNotifier
from telegram_bot.bot_commands import BotCommandHandler
from telegram_bot.formatter import format_event_message
from tray.tray_icon import TrayIcon


class KanijeApp:
    """
    Ana uygulama sınıfı. Tüm modülleri yönetir.

    Yaşam döngüsü:
        1. __init__: Config yükle, logger kur
        2. start(): Tüm modülleri başlat, event loop'a gir
        3. _event_loop(): Olayları al, işle, Telegram'a gönder
        4. shutdown(): Her şeyi düzgün kapat
    """

    def __init__(self, config_path: Optional[str] = None):
        self.start_time = datetime.now()
        self.recent_events: list[Event] = []  # Son 50 olayı tut
        self._running = False
        self._internet_was_up = False  # İnternet önceki durumu
        self._pending_messages: deque = deque(maxlen=100)  # İnternetsiz bekleyen bildirimler
        self._current_net_info: dict = {}  # Mevcut ağ bilgisi
        self._lock_file: Optional[str] = None

        # ── Tek Instance Kilidi ──
        self._acquire_lock()

        # ── Config ──
        self.config = ConfigManager(config_path)

        # ── Logger ──
        log_cfg = self.config.get("logging", {})
        setup_logger(
            log_file=log_cfg.get("file", "./kanije.log"),
            level=log_cfg.get("level", "INFO"),
            max_size_mb=log_cfg.get("max_size_mb", 10),
            backup_count=log_cfg.get("backup_count", 3),
            console_output=log_cfg.get("console_output", True),
        )
        self.log = get_logger("app")

        # ── Event Bus ──
        max_events = self.config.get("security.max_events_per_minute", 10)
        self.bus = EventBus(max_per_minute=max_events)

        # ── Aksiyon Modülleri ──
        self.camera = Camera(self.config.get("camera", {}))
        self.screenshot = Screenshot(self.config.get("screenshot", {}))

        # ── Telegram ──
        # Kimlik bilgileri .env'den okunur (config.yaml'da bulunmaz)
        tg_cfg = dict(self.config.get("telegram", {}) or {})
        tg_cfg["bot_token"] = self.config.get_bot_token()
        tg_cfg["chat_id"]   = self.config.get_chat_id()
        self.notifier = TelegramNotifier(tg_cfg)

        # ── Bot Komutları (iki yönlü) ──
        self.bot_commands = BotCommandHandler(
            config={**tg_cfg, **self.config.get("security", {})},
            app_context={
                "camera": self.camera,
                "screenshot": self.screenshot,
                "event_bus": self.bus,
                "config_manager": self.config,
                "start_time": self.start_time,
                "notifier": self.notifier,
                "recent_events": self.recent_events,
            },
        )

        # ── System Tray ──
        self.tray: Optional[TrayIcon] = None
        if self.config.get("tray.enabled", True):
            self.tray = TrayIcon(
                on_quit=self.shutdown,
            )

        # ── Listeners (platform'a göre) ──
        self.listeners = []

    def start(self):
        """Uygulamayı başlat."""
        self.log.info("=" * 50)
        self.log.info("🏰 Kanije Kalesi başlatılıyor...")
        self.log.info("=" * 50)

        # Config doğrulaması
        errors = self.config.validate()
        if errors:
            for err in errors:
                self.log.warning(f"Config uyarı: {err}")

        # Telegram bağlantı testi
        if self.notifier.token:
            if self.notifier.test_connection():
                self.log.info("Telegram bağlantısı başarılı ✅")
            else:
                self.log.error("Telegram bağlantısı başarısız ❌")
        else:
            self.log.warning("Telegram token ayarlanmamış — bildirimler devre dışı")

        # Kamera kontrolü
        if self.camera.is_available():
            self.log.info("Kamera erişilebilir ✅")
        else:
            self.log.warning("Kamera bulunamadı — fotoğraf çekimi devre dışı")

        # ── Sinyal yakalama (Ctrl+C + Windows shutdown) ──
        signal.signal(signal.SIGINT, lambda s, f: self.shutdown())
        signal.signal(signal.SIGTERM, lambda s, f: self.shutdown())
        try:
            import win32api
            win32api.SetConsoleCtrlHandler(self._windows_shutdown_handler, True)
        except ImportError:
            pass

        # ── Flag'i ÖNCE set et — thread'ler bu flag'e bakıyor ──
        self._running = True

        # ── Listener'ları başlat ──
        self._start_windows_listener()
        self._start_usb_monitor()
        self._start_power_monitor()
        self.bot_commands.start()

        # ── System Tray ──
        if self.tray:
            self.tray.start()

        # ── Arka plan thread'leri ──
        self._start_heartbeat()
        self._start_internet_monitor()

        # ── Başlangıç bildirimi ──
        self._send_boot_notification()

        self.log.info("🏰 Kanije Kalesi nöbete başladı!")

        # ── Ana event loop ──
        self._event_loop()


    def _event_loop(self):
        """
        Ana olay döngüsü.
        EventBus'tan olay alır → aksiyonları çalıştırır → Telegram'a gönderir.
        """
        while self._running:
            event = self.bus.get(timeout=1.0)
            if event is None:
                continue

            try:
                self._process_event(event)
            except Exception as e:
                self.log.error(f"Olay işleme hatası: {e}")

    def _process_event(self, event: Event):
        """Tek bir olayı işle: aksiyonlar + bildirim."""
        self.log.info(f"Olay işleniyor: {event}")

        # Son olaylar listesine ekle (max 50)
        self.recent_events.append(event)
        if len(self.recent_events) > 50:
            self.recent_events.pop(0)

        # Trigger config'ini kontrol et
        if not self.config.is_trigger_enabled(event.event_type):
            self.log.debug(f"Trigger devre dışı: {event.event_type}")
            return

        trigger_cfg = self.config.get_trigger(event.event_type)
        photos = []

        # Kamera çekim
        if trigger_cfg.get("capture_camera", False):
            photo = self.camera.capture()
            if photo:
                photos.append(photo)

        # Ekran görüntüsü
        if trigger_cfg.get("capture_screenshot", False):
            ss = self.screenshot.capture()
            if ss:
                photos.append(ss)

        # Uptime hesapla
        delta = datetime.now() - self.start_time
        hours, rem = divmod(int(delta.total_seconds()), 3600)
        minutes = rem // 60
        uptime_str = f"Aktif: {hours}s {minutes}dk" if hours > 0 else f"Aktif: {minutes}dk"

        # Telegram mesajı
        message = format_event_message(event, uptime_str)

        delete_after = self.config.get("security.delete_photos_after_send", True)

        if photos:
            success = self.notifier.send_photo(photos[0], caption=message, delete_after=delete_after)
            for extra_photo in photos[1:]:
                self.notifier.send_photo(extra_photo, delete_after=delete_after)
            if not success:
                self._pending_messages.append(message)
        else:
            success = self.notifier.send_message(message)
            if not success:
                self._pending_messages.append(message)

        self.log.info(f"Olay bildirildi: {event.event_type}")

    def _start_windows_listener(self):
        """Windows Event Listener'ı başlat."""
        try:
            from listeners.windows_event import WindowsEventListener
            listener = WindowsEventListener(self.bus, self.config.data)
            listener.start()
            self.listeners.append(listener)
        except ImportError:
            self.log.warning("pywin32 yüklü değil — Windows Event Listener devre dışı")
        except Exception as e:
            self.log.error(f"Windows Listener başlatılamadı: {e}")

    def _start_power_monitor(self):
        """Uyku/Hibernate ve uyanma monitörü."""
        try:
            from listeners.power_monitor import PowerMonitor
            monitor = PowerMonitor(self.bus)
            monitor.start()
            self.listeners.append(monitor)
        except ImportError:
            self.log.warning("wmi paketi yüklü değil — Güç monitörü devre dışı")
        except Exception as e:
            self.log.error(f"Güç monitörü başlatılamadı: {e}")

    def _start_usb_monitor(self):
        """USB Monitor'ü başlat."""
        usb_insert = self.config.is_trigger_enabled("usb_inserted")
        usb_remove = self.config.is_trigger_enabled("usb_removed")
        if not usb_insert and not usb_remove:
            return
        try:
            from listeners.usb_monitor import USBMonitor
            monitor = USBMonitor(self.bus, self.config.data)
            monitor.start()
            self.listeners.append(monitor)
        except ImportError:
            self.log.warning("wmi paketi yüklü değil — USB Monitor devre dışı")
        except Exception as e:
            self.log.error(f"USB Monitor başlatılamadı: {e}")

    def _start_heartbeat(self):
        """Periyodik durum raporu zamanlayıcısı."""
        if not self.config.get("heartbeat.enabled", True):
            return

        interval_hours = self.config.get("heartbeat.interval_hours", 6)
        interval_secs = interval_hours * 3600

        def heartbeat():
            while self._running:
                time.sleep(interval_secs)
                if not self._running:
                    break
                try:
                    msg = f"💚 Kanije Kalesi aktif\n🖥️ {socket.gethostname()}\n⏰ {datetime.now():%Y-%m-%d %H:%M}"
                    self.notifier.send_message(msg)
                    self.log.info("Heartbeat gönderildi")
                except Exception as e:
                    self.log.error(f"Heartbeat hatası: {e}")

        t = threading.Thread(target=heartbeat, name="Heartbeat", daemon=True)
        t.start()

    def _send_boot_notification(self):
        """Uygulama başlangıç bildirimi."""
        if not self.config.is_trigger_enabled("system_boot"):
            return

        event = Event(
            event_type="system_boot",
            hostname=socket.gethostname(),
        )
        self.bus.put(event)

    def _start_internet_monitor(self):
        """İnternet + ağ değişimi monitörü."""
        self._current_net_info = self._get_network_info()
        net_name = self._current_net_info.get("ssid") or self._current_net_info.get("type", "bilinmiyor")
        self.log.info(f"İnternet monitörü başlatıldı — mevcut ağ: {net_name} | IP: {self._current_net_info.get('ip', '?')}")

        def check_loop():
            while self._running:
                is_up = self._is_internet_available()

                # ── İnternet durumu değişimi ──
                if is_up and not self._internet_was_up:
                    self.log.info("İnternet bağlantısı kuruldu")

                    # Önce bekleyen mesajları gönder (disconnect dahil)
                    pending_count = len(self._pending_messages)
                    while self._pending_messages:
                        msg = self._pending_messages.popleft()
                        try:
                            self.notifier.send_message(msg)
                            time.sleep(0.3)
                        except Exception:
                            pass

                    # Güncel ağ bilgisini al
                    new_info = self._get_network_info()
                    new_name = new_info.get("ssid") or new_info.get("type", "bilinmiyor")

                    # İnternet bağlantı mesajını DOĞRUDAN gönder (event bus bypass)
                    connected_msg = (
                        "━━━━━━━━━━━━━━━━━━━\n"
                        "🌐 İNTERNET BAĞLANTISI KURULDU\n"
                        "━━━━━━━━━━━━━━━━━━━\n"
                        f"\n🖥️ Bilgisayar: {socket.gethostname()}\n"
                        f"📡 Ağ: {new_name}\n"
                        f"🌐 IP: {new_info.get('ip', '?')}\n"
                        f"⏰ Zaman: {datetime.now():%Y-%m-%d %H:%M:%S}\n"
                        + (f"📨 {pending_count} bekleyen bildirim gönderildi\n" if pending_count > 0 else "")
                        + f"─────────────────────\n"
                        f"🏰 Kanije Kalesi v1.0.0\n"
                        f"💡 Komutlar için /help"
                    )
                    self.notifier.send_message(connected_msg)
                    self._current_net_info = new_info

                elif not is_up and self._internet_was_up:
                    self.log.warning("İnternet bağlantısı kesildi")
                    lost_name = self._current_net_info.get("ssid") or self._current_net_info.get("type", "bilinmiyor")
                    disconnect_msg = (
                        "━━━━━━━━━━━━━━━━━━━\n"
                        "📡 İNTERNET BAĞLANTISI KESİLDİ\n"
                        "━━━━━━━━━━━━━━━━━━━\n"
                        f"\n🖥️ Bilgisayar: {socket.gethostname()}\n"
                        f"📡 Kaybedilen ağ: {lost_name}\n"
                        f"⏰ Kesilme: {datetime.now():%Y-%m-%d %H:%M:%S}\n"
                        f"─────────────────────\n"
                        f"🏰 Kanije Kalesi v1.0.0"
                    )
                    # Kuyruğa ekle — internet gelince gönderilecek
                    self._pending_messages.appendleft(disconnect_msg)

                # ── Ağ değişimi (internet varken) ──
                elif is_up and self._internet_was_up:
                    new_info = self._get_network_info()
                    old_id = self._current_net_info.get("ssid") or self._current_net_info.get("type", "")
                    new_id = new_info.get("ssid") or new_info.get("type", "")

                    if old_id and new_id and old_id != new_id:
                        self.log.info(f"Ağ değişti: {old_id} → {new_id}")
                        net_msg = (
                            "━━━━━━━━━━━━━━━━━━━\n"
                            "🔀 AĞ DEĞİŞTİRİLDİ\n"
                            "━━━━━━━━━━━━━━━━━━━\n"
                            f"\n🖥️ Bilgisayar: {socket.gethostname()}\n"
                            f"📡 Eski: {old_id}\n"
                            f"📡 Yeni: {new_id}\n"
                            f"🌐 IP: {new_info.get('ip', '?')}\n"
                            f"⏰ Zaman: {datetime.now():%Y-%m-%d %H:%M:%S}\n"
                            f"─────────────────────\n"
                            f"🏰 Kanije Kalesi v1.0.0"
                        )
                        self.notifier.send_message(net_msg)

                    self._current_net_info = new_info

                self._internet_was_up = is_up
                time.sleep(5)

        t = threading.Thread(target=check_loop, name="InternetMonitor", daemon=True)
        t.start()


    @staticmethod
    def _is_internet_available() -> bool:
        """İnternet bağlantısını kontrol et (Telegram API'ye erişim)."""
        try:
            socket.create_connection(("api.telegram.org", 443), timeout=5).close()
            return True
        except (socket.timeout, socket.error, OSError):
            return False

    @staticmethod
    def _get_network_info() -> dict:
        """Aktif ağ bilgisini döndür: {"type": "WiFi"|"Ethernet", "ssid": "...", "ip": "..."}"""
        info = {"type": "", "ssid": "", "ip": ""}

        # 1) WiFi SSID'yi al (Windows: netsh)
        try:
            import subprocess
            import re
            result = subprocess.run(
                ["netsh", "wlan", "show", "interfaces"],
                capture_output=True,   # text=True YOK — encoding'i kendimiz yönetiyoruz
                timeout=5,
                creationflags=0x08000000,  # CREATE_NO_WINDOW
            )
            if result.returncode == 0:
                # Türkçe Windows OEM encoding (cp857) → fallback utf-8 → cp1254
                raw = result.stdout
                for enc in ("utf-8", "cp1254", "cp857", "latin-1"):
                    try:
                        output = raw.decode(enc)
                        break
                    except (UnicodeDecodeError, LookupError):
                        continue
                else:
                    output = raw.decode("utf-8", errors="replace")

                ssid_match = re.search(
                    r"^\s+SSID\s*:\s*(.+)$",
                    output,
                    re.MULTILINE,
                )
                connected = bool(re.search(
                    r"(State|Durum|État)\s*:\s*(connected|Bağlandı|bağlı|Connecté)",
                    output,
                    re.IGNORECASE,
                ))
                if ssid_match and connected:
                    info["type"] = "WiFi"
                    info["ssid"] = ssid_match.group(1).strip()
        except Exception:
            pass

        # 2) IP adresini ve adaptör tipini al
        try:
            import psutil
            addrs = psutil.net_if_addrs()
            stats = psutil.net_if_stats()
            for iface, addr_list in addrs.items():
                if iface.lower() in ("loopback pseudo-interface 1", "lo"):
                    continue
                if iface.lower().startswith(("vethernet", "vmware", "virtualbox", "vbox")):
                    continue
                if iface in stats and stats[iface].isup:
                    for addr in addr_list:
                        if addr.family.name == "AF_INET" and addr.address != "127.0.0.1":
                            info["ip"] = addr.address
                            # WiFi tespit edilmediyse adaptör adına bak
                            if not info["type"]:
                                iface_l = iface.lower()
                                if "ethernet" in iface_l or "eth" in iface_l or "local area" in iface_l:
                                    info["type"] = "Ethernet"
                                elif "wi-fi" in iface_l or "wifi" in iface_l or "wireless" in iface_l:
                                    info["type"] = "WiFi"
                                else:
                                    info["type"] = iface
                            break
                    if info["ip"]:
                        break
        except Exception:
            pass

        return info

    def _acquire_lock(self):
        """Tek instance kilidi al. Eski instance çalışıyorsa öldür."""
        import tempfile
        lock_path = os.path.join(tempfile.gettempdir(), "kanije_kalesi.lock")

        # Eski PID'yi kontrol et
        if os.path.exists(lock_path):
            try:
                with open(lock_path, "r") as f:
                    old_pid = int(f.read().strip())
                # Eski süreç hâlâ çalışıyor mu?
                import psutil
                if psutil.pid_exists(old_pid):
                    try:
                        proc = psutil.Process(old_pid)
                        if "python" in proc.name().lower():
                            proc.terminate()
                            proc.wait(timeout=5)
                    except Exception:
                        pass
            except (ValueError, FileNotFoundError):
                pass

        # Kendi PID'mizi yaz
        with open(lock_path, "w") as f:
            f.write(str(os.getpid()))
        self._lock_file = lock_path

    def _release_lock(self):
        """Kilit dosyasını sil."""
        if self._lock_file and os.path.exists(self._lock_file):
            try:
                os.remove(self._lock_file)
            except OSError:
                pass

    def _windows_shutdown_handler(self, ctrl_type):
        """Windows shutdown/logoff sinyallerini yakala."""
        # CTRL_SHUTDOWN_EVENT=6, CTRL_LOGOFF_EVENT=5, CTRL_CLOSE_EVENT=2
        if ctrl_type in (2, 5, 6):
            self.shutdown()
            return True
        return False

    def shutdown(self):
        """Tüm modülleri düzgün kapat."""
        if not self._running:
            return

        self.log.info("Kanije Kalesi kapatılıyor...")
        self._running = False

        # Kapatılma bildirimi (3 sn timeout — OS shutdown'u bloklama)
        try:
            import requests
            delta = datetime.now() - self.start_time
            hours, rem = divmod(int(delta.total_seconds()), 3600)
            minutes = rem // 60
            uptime = f"{hours}s {minutes}dk" if hours > 0 else f"{minutes}dk"

            shutdown_msg = (
                "━━━━━━━━━━━━━━━━━━━\n"
                "🔴 SİSTEM KAPANIYOR\n"
                "━━━━━━━━━━━━━━━━━━━\n"
                f"\n🖥️ Bilgisayar: {socket.gethostname()}\n"
                f"⏰ Zaman: {datetime.now():%Y-%m-%d %H:%M:%S}\n"
                f"📊 Nöbet süresi: {uptime}\n"
                f"\n─────────────────────\n"
                f"🏰 Kanije Kalesi v1.0.0 · Nöbet bitti."
            )
            url = f"https://api.telegram.org/bot{self.notifier.token}/sendMessage"
            requests.post(url, data={
                "chat_id": self.notifier.chat_id,
                "text": shutdown_msg,
            }, timeout=3)
        except Exception:
            pass

        # Listener'ları durdur
        for listener in self.listeners:
            try:
                listener.stop()
            except Exception:
                pass

        # Bot komutlarını durdur
        try:
            self.bot_commands.stop()
        except Exception:
            pass

        # Tray'i kapat
        if self.tray:
            try:
                self.tray.stop()
            except Exception:
                pass

        # Kilit dosyasını sil
        self._release_lock()

        self.log.info("🏰 Kanije Kalesi kapatıldı. Nöbet bitti.")

    def test_telegram(self):
        """Telegram bağlantısını test et ve test mesajı gönder."""
        self.log.info("Telegram test başlatılıyor...")

        if not self.notifier.token:
            self.log.error("Bot token boş! config.yaml'da telegram.bot_token ayarla.")
            return False

        if not self.notifier.test_connection():
            self.log.error("Bot bağlantısı başarısız!")
            return False

        msg = (
            "━━━━━━━━━━━━━━━━━━━\n"
            "🏰 KANİJE KALESİ — TEST\n"
            "━━━━━━━━━━━━━━━━━━━\n"
            "\n"
            f"✅ Bağlantı başarılı!\n"
            f"🖥️ Bilgisayar: {socket.gethostname()}\n"
            f"⏰ Zaman: {datetime.now():%Y-%m-%d %H:%M:%S}\n"
            f"\n"
            f"Bildirimler bu sohbete gelecek.\n"
            f"Komutlar için /help yazabilirsin."
        )

        if self.notifier.send_message(msg):
            self.log.info("Test mesajı gönderildi ✅")
            return True
        else:
            self.log.error("Test mesajı gönderilemedi ❌")
            return False
