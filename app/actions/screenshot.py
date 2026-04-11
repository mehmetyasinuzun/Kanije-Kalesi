"""
Kanije Kalesi — Ekran Görüntüsü Modülü
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
mss kütüphanesi ile ekran görüntüsü alır.
Hafif, hızlı, cross-platform. PIL/Pillow bağımlılığı yok.
"""

import tempfile
from pathlib import Path
from datetime import datetime
from typing import Optional

from core.logger import get_logger

log = get_logger("screenshot")


class Screenshot:
    """
    Masaüstünün ekran görüntüsünü alır.

    Kullanım:
        ss = Screenshot(config)
        path = ss.capture()
        if path:
            # Telegram'a gönder
    """

    def __init__(self, config: dict):
        self.jpeg_quality = config.get("jpeg_quality", 70)
        self.save_local = config.get("save_local", False)
        self.local_path = config.get("local_path", "./captures/")

    def capture(self) -> Optional[str]:
        """
        Tüm ekranın görüntüsünü al ve PNG olarak kaydet.
        Başarılıysa dosya yolunu döndürür. Hata varsa None.
        """
        try:
            import mss
            import mss.tools
        except ImportError:
            log.error("mss yüklü değil — pip install mss")
            return None

        try:
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")

            if self.save_local:
                save_dir = Path(self.local_path)
                save_dir.mkdir(parents=True, exist_ok=True)
                filepath = str(save_dir / f"kanije_ss_{timestamp}.png")
            else:
                filepath = str(Path(tempfile.gettempdir()) / f"kanije_ss_{timestamp}.png")

            with mss.mss() as sct:
                # Tüm monitörleri kapsayan alan (multi-monitor)
                monitor = sct.monitors[0]  # 0 = tüm ekranlar birleşik
                screenshot = sct.grab(monitor)
                mss.tools.to_png(screenshot.rgb, screenshot.size, output=filepath)

            log.info(f"Ekran görüntüsü alındı: {filepath}")
            return filepath

        except Exception as e:
            log.error(f"Ekran görüntüsü hatası: {e}")
            return None
