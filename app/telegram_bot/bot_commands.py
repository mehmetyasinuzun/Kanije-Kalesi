"""
Kanije Kalesi — Telegram Bot Komut İşleyici
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
İki yönlü iletişim: Kullanıcı Telegram'dan komut gönderir,
bot cevap verir. Polling tabanlı (webhook gerektirmez).

Desteklenen Komutlar:
    /status      — Sistem durumu (CPU, RAM, Disk, uptime)
    /photo       — Anlık kamera fotoğrafı çek ve gönder
    /screenshot  — Anlık ekran görüntüsü al ve gönder
    /ping        — Canlılık kontrolü (hızlı yanıt)
    /events      — Son 10 olay listesi
    /config      — Aktif yapılandırma (token gizli)
    /lock        — Ekranı kilitle
    /restart     — Bilgisayarı yeniden başlat (onay gerekir)
    /shutdown    — Bilgisayarı kapat (onay gerekir)
    /cancel      — Bekleyen restart/shutdown iptal
    /help        — Komut listesi
"""

import threading
import subprocess
import socket
import time
from datetime import datetime, timedelta
from typing import Callable, Optional

import requests
import psutil

from core.logger import get_logger
from core.event_bus import Event
from telegram_bot.formatter import format_status_message, format_help_message

log = get_logger("bot_cmd")

API_BASE = "https://api.telegram.org/bot{token}"


class BotCommandHandler:
    """
    Telegram Bot API polling ile komutları dinler ve işler.
    Yalnızca izinli chat_id'lerden gelen komutları kabul eder.

    Kullanım:
        handler = BotCommandHandler(config, app_context)
        handler.start()  # Arka plan thread'inde polling başlar
    """

    def __init__(self, config: dict, app_context: dict):
        """
        Args:
            config: Telegram config bölümü
            app_context: Uygulamanın paylaşılan durumu
                - "camera": Camera nesnesi
                - "screenshot": Screenshot nesnesi
                - "event_bus": EventBus nesnesi
                - "config_manager": ConfigManager nesnesi
                - "start_time": Uygulama başlangıç zamanı
                - "notifier": TelegramNotifier nesnesi
                - "recent_events": Son olayların listesi
        """
        self.token = config.get("bot_token", "")
        self.allowed_chat_ids = self._build_allowed_ids(config)
        self.ctx = app_context
        self._base_url = API_BASE.format(token=self.token)
        self._thread = None
        self._running = False
        self._last_update_id = 0
        self._pending_action = None  # {"type": "shutdown"|"restart", "timer": Timer, "chat_id": str}

    def _build_allowed_ids(self, config: dict) -> set:
        """İzinli chat ID'lerini set olarak oluştur."""
        ids = set()
        main_id = str(config.get("chat_id", ""))
        if main_id:
            ids.add(main_id)
        extra = config.get("allowed_chat_ids", [])
        for eid in extra:
            ids.add(str(eid))
        return ids

    def start(self):
        """Polling thread'ini başlat."""
        if not self.token:
            log.warning("Bot token boş — komut dinleme devre dışı")
            return

        self._running = True
        self._thread = threading.Thread(
            target=self._polling_loop,
            name="BotCommands",
            daemon=True,
        )
        self._thread.start()
        log.info("Telegram Bot komut dinleyici başlatıldı (polling)")

    def stop(self):
        """Polling'i durdur."""
        self._running = False
        log.info("Telegram Bot komut dinleyici durduruldu")

    def _polling_loop(self):
        """
        Telegram getUpdates API'sini düzenli aralıklarla çağırarak
        yeni mesajları alır ve işler.
        Bağlantı hatasında exponential backoff uygular.
        """
        log.info("Polling döngüsü başladı")
        sleep_interval = 2       # Başlangıç bekleme (saniye)
        max_interval = 30        # Maksimum bekleme (saniye)
        base_interval = 2

        while self._running:
            try:
                updates = self._get_updates()
                for update in updates:
                    self._process_update(update)
                sleep_interval = base_interval  # Başarılı → normal aralık
            except requests.ConnectionError:
                # İnternet yok — gürültülü loglama değil, sessizce bekle
                log.debug(f"Bağlantı yok, {sleep_interval}sn bekleniyor")
                sleep_interval = min(sleep_interval * 2, max_interval)
            except requests.Timeout:
                log.debug(f"Timeout, {sleep_interval}sn bekleniyor")
                sleep_interval = min(sleep_interval * 2, max_interval)
            except requests.RequestException as e:
                log.warning(f"Polling hatası: {e}")
                sleep_interval = min(sleep_interval * 2, max_interval)
            except Exception as e:
                log.error(f"Polling beklenmeyen hata: {e}")
                sleep_interval = base_interval

            time.sleep(sleep_interval)


    def _get_updates(self) -> list:
        """Telegram'dan yeni mesajları al."""
        url = f"{self._base_url}/getUpdates"
        params = {
            "offset": self._last_update_id + 1,
            "timeout": 10,  # Long polling (10 sn bekle, yeni mesaj gelirse hemen dön)
            "allowed_updates": ["message"],
        }

        resp = requests.get(url, params=params, timeout=15)
        data = resp.json()

        if not data.get("ok"):
            return []

        results = data.get("result", [])

        # Son update_id'yi güncelle
        if results:
            self._last_update_id = results[-1]["update_id"]

        return results

    def _process_update(self, update: dict):
        """Tek bir update mesajını işle."""
        message = update.get("message", {})
        text = message.get("text", "").strip()
        chat = message.get("chat", {})
        chat_id = str(chat.get("id", ""))
        from_user = message.get("from", {})
        username = from_user.get("username", from_user.get("first_name", "?"))

        # Güvenlik: yalnızca izinli chat_id'lerden gelen komutları işle
        if chat_id not in self.allowed_chat_ids:
            log.warning(f"Yetkisiz komut girişimi: chat_id={chat_id}, user={username}, text={text}")
            return

        if not text.startswith("/"):
            return  # Komut değil — atla

        command = text.split()[0].lower()
        # /command@botname formatını temizle
        if "@" in command:
            command = command.split("@")[0]

        log.info(f"Komut alındı: {command} (user={username})")

        # Komut yönlendirme
        handlers = {
            "/status": self._cmd_status,
            "/photo": self._cmd_photo,
            "/screenshot": self._cmd_screenshot,
            "/ping": self._cmd_ping,
            "/events": self._cmd_events,
            "/config": self._cmd_config,
            "/lock": self._cmd_lock,
            "/restart": self._cmd_restart,
            "/shutdown": self._cmd_shutdown,
            "/cancel": self._cmd_cancel,
            "/help": self._cmd_help,
            "/start": self._cmd_help,
        }

        handler = handlers.get(command)
        if handler:
            try:
                handler(chat_id)
            except Exception as e:
                log.error(f"Komut işleme hatası ({command}): {e}")
                self._reply(chat_id, f"❌ Komut hatası: {e}")
        else:
            self._reply(chat_id, f"❓ Bilinmeyen komut: {command}\n/help yazarak komut listesini görebilirsin.")

    # ── Komut İşleyicileri ──

    def _cmd_status(self, chat_id: str):
        """Sistem durumu"""
        info = {
            "hostname": socket.gethostname(),
            "app_uptime": self._get_app_uptime(),
            "sys_uptime": self._get_sys_uptime(),
            "cpu_percent": psutil.cpu_percent(interval=1),
            "ram_used": f"{psutil.virtual_memory().used / (1024**3):.1f} GB",
            "ram_total": f"{psutil.virtual_memory().total / (1024**3):.1f} GB",
            "disk_used": f"{psutil.disk_usage('C:/').used / (1024**3):.0f} GB",
            "disk_total": f"{psutil.disk_usage('C:/').total / (1024**3):.0f} GB",
            "total_events": self.ctx.get("event_bus").stats["total_events"] if self.ctx.get("event_bus") else 0,
            "dropped_events": self.ctx.get("event_bus").stats["dropped_events"] if self.ctx.get("event_bus") else 0,
            "last_event": self._get_last_event_str(),
        }
        msg = format_status_message(info)
        self._reply(chat_id, msg)

    def _cmd_photo(self, chat_id: str):
        """Anlık kamera fotoğrafı"""
        camera = self.ctx.get("camera")
        if not camera:
            self._reply(chat_id, "📷 Kamera modülü devre dışı.")
            return

        self._reply(chat_id, "📸 Fotoğraf çekiliyor...")
        path = camera.capture()
        if path:
            notifier = self.ctx.get("notifier")
            if notifier:
                notifier.send_photo(path, caption="📸 Anlık kamera görüntüsü")
            else:
                self._reply(chat_id, "❌ Notifier bulunamadı.")
        else:
            self._reply(chat_id, "❌ Kamera fotoğrafı alınamadı.")

    def _cmd_screenshot(self, chat_id: str):
        """Anlık ekran görüntüsü"""
        ss = self.ctx.get("screenshot")
        if not ss:
            self._reply(chat_id, "🖥️ Screenshot modülü devre dışı.")
            return

        self._reply(chat_id, "🖥️ Ekran görüntüsü alınıyor...")
        path = ss.capture()
        if path:
            notifier = self.ctx.get("notifier")
            if notifier:
                notifier.send_photo(path, caption="🖥️ Anlık ekran görüntüsü")
            else:
                self._reply(chat_id, "❌ Notifier bulunamadı.")
        else:
            self._reply(chat_id, "❌ Ekran görüntüsü alınamadı.")

    def _cmd_ping(self, chat_id: str):
        """Hızlı canlılık kontrolü"""
        self._reply(chat_id, f"🏰 Kanije ayakta! Uptime: {self._get_app_uptime()}")

    def _cmd_events(self, chat_id: str):
        """Son olaylar"""
        recent = self.ctx.get("recent_events", [])
        if not recent:
            self._reply(chat_id, "📋 Henüz kayıtlı olay yok.")
            return

        lines = ["📋 *Son Olaylar:*\n"]
        for i, ev in enumerate(recent[-10:], 1):
            icon = {"login_success": "✅", "login_failed": "🚨", "screen_unlock": "🔓",
                    "usb_inserted": "🔌", "system_boot": "🖥️"}.get(ev.event_type, "📌")
            lines.append(f"{i}. {icon} {ev.event_type} — {ev.username or '-'} — {ev.timestamp:%H:%M:%S}")

        self._reply(chat_id, "\n".join(lines), parse_mode="Markdown")

    def _cmd_config(self, chat_id: str):
        """Aktif yapılandırmayı insan okunur formatta göster."""
        cm = self.ctx.get("config_manager")
        if not cm:
            self._reply(chat_id, "⚙️ Config bulunamadı.")
            return

        def state(enabled: bool) -> str:
            return "✅ Açık" if enabled else "⛔ Kapalı"

        # ── Telegram ──
        token = cm.get("telegram.bot_token", "")
        token_display = f"{token[:6]}...{token[-4:]}" if len(token) > 10 else "⚠️ Ayarlanmamış"
        chat_id_cfg = cm.get("telegram.chat_id", "⚠️ Ayarlanmamış")

        # ── Tetikleyiciler ──
        triggers = {
            "login_success":  ("✅ Başarılı giriş",     "login_success"),
            "login_failed":   ("🚨 Başarısız giriş",    "login_failed"),
            "screen_lock":    ("🔒 Ekran kilitleme",     "screen_lock"),
            "screen_unlock":  ("🔓 Kilit açma",         "screen_unlock"),
            "system_boot":    ("🖥️ Sistem başlangıcı",  "system_boot"),
            "system_sleep":   ("😴 Uyku/Hibernate",     "system_sleep"),
            "system_wake":    ("☀️ Uykudan uyanma",     "system_wake"),
            "usb_inserted":   ("🔌 USB takma",          "usb_inserted"),
            "usb_removed":    ("⏏️ USB çıkarma",        "usb_removed"),
        }
        trigger_lines = []
        for key, (label, trigger_key) in triggers.items():
            enabled = cm.is_trigger_enabled(trigger_key)
            cam = cm.get(f"triggers.{trigger_key}.capture_camera", False)
            cam_icon = " 📷" if cam else ""
            trigger_lines.append(f"  {'✅' if enabled else '⛔'} {label}{cam_icon}")

        # ── Kamera ──
        cam_idx = cm.get("camera.device_index", 0)
        cam_res = cm.get("camera.resolution", [640, 480])

        # ── Heartbeat ──
        hb_enabled = cm.get("heartbeat.enabled", True)
        hb_hours = cm.get("heartbeat.interval_hours", 6)

        # ── Güvenlik ──
        del_photos = cm.get("security.delete_photos_after_send", True)
        max_ev = cm.get("security.max_events_per_minute", 10)

        lines = [
            "━━━━━━━━━━━━━━━━━━━",
            "⚙️ KANİJE KALESİ — AYARLAR",
            "━━━━━━━━━━━━━━━━━━━",
            "",
            "📡 *Telegram*",
            f"  🔑 Token: `{token_display}`",
            f"  💬 Chat ID: `{chat_id_cfg}`",
            "",
            "🎯 *Tetikleyiciler*",
        ] + trigger_lines + [
            "",
            "📷 *Kamera*",
            f"  🎥 Kamera: #{cam_idx}  |  📐 Çözünürlük: {cam_res[0]}×{cam_res[1]}",
            "",
            "💚 *Heartbeat*",
            f"  {state(hb_enabled)}  |  ⏱️ Her {hb_hours} saatte bir",
            "",
            "🔐 *Güvenlik*",
            f"  🗑️ Fotoğraf gönderim sonrası sil: {state(del_photos)}",
            f"  🚦 Rate limit: {max_ev} olay/dk",
            "",
            "─────────────────────",
            "🏰 Kanije Kalesi v1.0.0",
        ]

        self._reply(chat_id, "\n".join(lines), parse_mode="Markdown")


    def _cmd_lock(self, chat_id: str):
        """Ekranı kilitle"""
        try:
            subprocess.run(["rundll32.exe", "user32.dll,LockWorkStation"], check=True)
            self._reply(chat_id, "🔒 Ekran kilitlendi.")
            log.info("Ekran uzaktan kilitlendi")
        except Exception as e:
            self._reply(chat_id, f"❌ Kilitleme hatası: {e}")

    def _cmd_restart(self, chat_id: str):
        """Bilgisayarı yeniden başlat (onay gerekir)"""
        if self._pending_action:
            self._reply(chat_id, "⚠️ Zaten bekleyen bir işlem var. Önce /cancel yaz.")
            return

        self._reply(chat_id, "🔄 Bilgisayar 15 saniye içinde yeniden başlatılacak.\n❌ İptal için /cancel yaz.")

        def do_restart():
            self._pending_action = None
            self._reply(chat_id, "🔄 Yeniden başlatılıyor...")
            log.info("Bilgisayar uzaktan yeniden başlatılıyor")
            subprocess.run(["shutdown", "/r", "/t", "3"], check=True)

        timer = threading.Timer(15.0, do_restart)
        timer.start()
        self._pending_action = {"type": "restart", "timer": timer, "chat_id": chat_id}

    def _cmd_shutdown(self, chat_id: str):
        """Bilgisayarı kapat (onay gerekir)"""
        if self._pending_action:
            self._reply(chat_id, "⚠️ Zaten bekleyen bir işlem var. Önce /cancel yaz.")
            return

        self._reply(chat_id, "⏻ Bilgisayar 15 saniye içinde kapanacak.\n❌ İptal için /cancel yaz.")

        def do_shutdown():
            self._pending_action = None
            self._reply(chat_id, "⏻ Kapatılıyor...")
            log.info("Bilgisayar uzaktan kapatılıyor")
            subprocess.run(["shutdown", "/s", "/t", "3"], check=True)

        timer = threading.Timer(15.0, do_shutdown)
        timer.start()
        self._pending_action = {"type": "shutdown", "timer": timer, "chat_id": chat_id}

    def _cmd_cancel(self, chat_id: str):
        """Bekleyen restart/shutdown iptal"""
        if not self._pending_action:
            self._reply(chat_id, "✅ Bekleyen işlem yok.")
            return

        action_type = self._pending_action["type"]
        self._pending_action["timer"].cancel()
        self._pending_action = None
        name = "Yeniden başlatma" if action_type == "restart" else "Kapatma"
        self._reply(chat_id, f"✅ {name} iptal edildi.")
        log.info(f"{name} uzaktan iptal edildi")

    def _cmd_help(self, chat_id: str):
        """Yardım mesajı"""
        self._reply(chat_id, format_help_message(), parse_mode="Markdown")

    # ── Yardımcılar ──

    def _reply(self, chat_id: str, text: str, parse_mode: str = None):
        """Belirtilen chat'e yanıt gönder."""
        url = f"{self._base_url}/sendMessage"
        payload = {"chat_id": chat_id, "text": text}
        if parse_mode:
            payload["parse_mode"] = parse_mode
        try:
            requests.post(url, data=payload, timeout=10)
        except requests.RequestException as e:
            log.error(f"Yanıt gönderilemedi: {e}")

    def _get_app_uptime(self) -> str:
        """Uygulama çalışma süresini insan okunur formatta döndür."""
        start = self.ctx.get("start_time")
        if not start:
            return "?"
        delta = datetime.now() - start
        hours, remainder = divmod(int(delta.total_seconds()), 3600)
        minutes, seconds = divmod(remainder, 60)
        if hours > 0:
            return f"{hours}s {minutes}dk"
        return f"{minutes}dk {seconds}sn"

    def _get_sys_uptime(self) -> str:
        """Sistem uptime'ını döndür."""
        try:
            boot = datetime.fromtimestamp(psutil.boot_time())
            delta = datetime.now() - boot
            days = delta.days
            hours, remainder = divmod(delta.seconds, 3600)
            minutes = remainder // 60
            parts = []
            if days > 0:
                parts.append(f"{days}g")
            parts.append(f"{hours}s {minutes}dk")
            return " ".join(parts)
        except Exception:
            return "?"

    def _get_last_event_str(self) -> str:
        """Son olayın açıklamasını döndür."""
        recent = self.ctx.get("recent_events", [])
        if not recent:
            return "henüz yok"
        last = recent[-1]
        return f"{last.event_type} — {last.timestamp:%H:%M:%S}"
