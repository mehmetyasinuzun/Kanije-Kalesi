"""
Kanije Kalesi — System Tray İkonu
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Sistem tepsisinde küçük ikon gösterir.
Yeşil = çalışıyor, Kırmızı = hata.
Sağ tık menüsü: Durum, Log Aç, Çıkış.
"""

import threading
from pathlib import Path
from typing import Callable, Optional

from core.logger import get_logger

log = get_logger("tray")


def _create_icon_image(color: str = "green"):
    """
    Programatik olarak basit bir renkli daire ikonu oluşturur.
    PIL/Pillow gerektirir (pystray zaten bağımlılık olarak çeker).
    """
    try:
        from PIL import Image, ImageDraw

        size = 64
        img = Image.new("RGBA", (size, size), (0, 0, 0, 0))
        draw = ImageDraw.Draw(img)

        colors = {
            "green": (76, 175, 80, 255),
            "red": (244, 67, 54, 255),
            "yellow": (255, 193, 7, 255),
        }
        fill = colors.get(color, colors["green"])

        # Daire çiz
        draw.ellipse([4, 4, size - 4, size - 4], fill=fill)

        # İçine K harfi (Kanije)
        try:
            draw.text((20, 14), "K", fill=(255, 255, 255, 255))
        except Exception:
            pass

        return img
    except ImportError:
        log.warning("Pillow yüklü değil — tray ikonu oluşturulamadı")
        return None


class TrayIcon:
    """
    Sistem tepsisi ikonu.

    Kullanım:
        tray = TrayIcon(on_quit=app.shutdown)
        tray.start()
        tray.set_status("running")  # yeşil
        tray.set_status("error")    # kırmızı
    """

    def __init__(self, on_quit: Callable = None, on_open_log: Callable = None,
                 on_setup: Callable = None):
        self.on_quit = on_quit
        self.on_open_log = on_open_log
        self.on_setup = on_setup
        self._icon = None
        self._thread = None

    def start(self):
        """Tray ikonunu arka plan thread'inde başlat."""
        self._thread = threading.Thread(
            target=self._run,
            name="TrayIcon",
            daemon=True,
        )
        self._thread.start()
        log.info("System Tray ikonu başlatıldı")

    def _run(self):
        """pystray ana döngüsü."""
        try:
            import pystray
            from pystray import MenuItem

            image = _create_icon_image("green")
            if image is None:
                log.warning("Tray ikonu oluşturulamadı — tray devre dışı")
                return

            menu = pystray.Menu(
                MenuItem("🏰 Kanije Kalesi — Aktif", None, enabled=False),
                pystray.Menu.SEPARATOR,
                MenuItem("⚙️ Telegram Ayarları", self._action_setup),
                MenuItem("📋 Log Dosyasını Aç", self._action_open_log),
                pystray.Menu.SEPARATOR,
                MenuItem("❌ Çıkış", self._action_quit),
            )

            self._icon = pystray.Icon(
                name="KanijeKalesi",
                icon=image,
                title="Kanije Kalesi — Güvenlik İzleme",
                menu=menu,
            )

            self._icon.run()

        except ImportError:
            log.warning("pystray yüklü değil — tray devre dışı. pip install pystray")
        except Exception as e:
            log.error(f"Tray hatası: {e}")

    def set_status(self, status: str):
        """
        İkon rengini değiştir.
        status: "running" (yeşil), "error" (kırmızı), "warning" (sarı)
        """
        if self._icon is None:
            return

        color_map = {
            "running": "green",
            "error": "red",
            "warning": "yellow",
        }
        color = color_map.get(status, "green")
        new_img = _create_icon_image(color)
        if new_img:
            self._icon.icon = new_img

    def stop(self):
        """Tray ikonunu kapat."""
        if self._icon:
            try:
                self._icon.stop()
            except Exception:
                pass

    def _action_quit(self, icon, item):
        """Çıkış menü aksiyonu."""
        log.info("Kullanıcı tray'den çıkış istedi")
        if self.on_quit:
            self.on_quit()
        self.stop()

    def _action_open_log(self, icon, item):
        """Log dosyasını aç."""
        if self.on_open_log:
            self.on_open_log()
        else:
            import subprocess
            log_path = Path("kanije.log")
            if log_path.exists():
                subprocess.Popen(["notepad.exe", str(log_path)])

    def _action_setup(self, icon, item):
        """Telegram kimlik bilgileri ayar penceresi (ayrı süreç)."""
        if self.on_setup:
            self.on_setup()
        else:
            try:
                from core.setup_wizard import open_settings_subprocess
                open_settings_subprocess()
            except Exception as e:
                log.error(f"Setup penceresi açılamadı: {e}")
