"""
Kanije Kalesi — Kamera Modülü
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
OpenCV ile kameradan fotoğraf çeker.
Thread-safe: aynı anda sadece 1 çekim yapılır.
Kamera yoksa hata vermez, None döner.
"""

import threading
import tempfile
from pathlib import Path
from datetime import datetime
from typing import Optional

from core.logger import get_logger

log = get_logger("camera")


class Camera:
    """
    Webcam'den tek kare fotoğraf çeker.

    Kullanım:
        cam = Camera(config)
        photo_path = cam.capture()
        if photo_path:
            # Telegram'a gönder
    """

    def __init__(self, config: dict):
        self.device_index = config.get("device_index", 0)
        self.resolution = config.get("resolution", [640, 480])
        self.warmup_frames = config.get("warmup_frames", 5)
        self.jpeg_quality = config.get("jpeg_quality", 85)
        self.save_local = config.get("save_local", False)
        self.local_path = config.get("local_path", "./captures/")
        self._lock = threading.Lock()

    def capture(self) -> Optional[str]:
        """
        Kameradan tek kare çek ve JPEG olarak kaydet.
        Başarılıysa dosya yolunu döndürür. Hata varsa None döner.

        Thread-safe: aynı anda 2 çekim yapılmaz.
        """
        # Başka bir thread zaten çekiyorsa atla
        if not self._lock.acquire(blocking=False):
            log.warning("Kamera zaten kullanımda — çekim atlandı")
            return None

        try:
            return self._do_capture()
        finally:
            self._lock.release()

    def _do_capture(self) -> Optional[str]:
        """Kamera çekim işlemi (iç fonksiyon)."""
        try:
            import cv2
        except ImportError:
            log.error("opencv-python-headless yüklü değil — pip install opencv-python-headless")
            return None

        cap = None
        try:
            cap = cv2.VideoCapture(self.device_index, cv2.CAP_DSHOW)

            if not cap.isOpened():
                log.warning(f"Kamera açılamadı (index={self.device_index})")
                return None

            # Çözünürlük ayarla
            cap.set(cv2.CAP_PROP_FRAME_WIDTH, self.resolution[0])
            cap.set(cv2.CAP_PROP_FRAME_HEIGHT, self.resolution[1])

            # Warmup: ilk kareler genellikle karanlık/bulanık olur
            for _ in range(self.warmup_frames):
                cap.read()

            # Asıl kareyi çek
            ret, frame = cap.read()

            if not ret or frame is None:
                log.warning("Kamera karesi okunamadı")
                return None

            # Dosya yolu belirle
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")

            if self.save_local:
                save_dir = Path(self.local_path)
                save_dir.mkdir(parents=True, exist_ok=True)
                filepath = str(save_dir / f"kanije_cam_{timestamp}.jpg")
            else:
                filepath = str(Path(tempfile.gettempdir()) / f"kanije_cam_{timestamp}.jpg")

            # JPEG olarak kaydet
            params = [cv2.IMWRITE_JPEG_QUALITY, self.jpeg_quality]
            cv2.imwrite(filepath, frame, params)

            log.info(f"Fotoğraf çekildi: {filepath}")
            return filepath

        except Exception as e:
            log.error(f"Kamera hatası: {e}")
            return None

        finally:
            if cap is not None:
                cap.release()

    def is_available(self) -> bool:
        """Kameranın erişilebilir olup olmadığını test et."""
        try:
            import cv2
            cap = cv2.VideoCapture(self.device_index, cv2.CAP_DSHOW)
            ok = cap.isOpened()
            cap.release()
            return ok
        except Exception:
            return False
