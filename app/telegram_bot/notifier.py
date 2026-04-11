"""
Kanije Kalesi — Telegram Bildirim Gönderici
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Tek yönlü: Uygulama → Telegram.
Metin mesajı ve fotoğraf gönderme. Retry mekanizması dahil.
"""

import os
import time
from pathlib import Path
from typing import Optional

import requests

from core.logger import get_logger

log = get_logger("telegram")

# Telegram Bot API base URL
API_BASE = "https://api.telegram.org/bot{token}"


class TelegramNotifier:
    """
    Telegram Bot API üzerinden mesaj ve fotoğraf gönderir.

    Kullanım:
        notifier = TelegramNotifier(config)
        notifier.send_message("Test mesajı")
        notifier.send_photo("/path/to/photo.jpg", caption="Kamera görüntüsü")
    """

    def __init__(self, config: dict):
        self.token = config.get("bot_token", "")
        self.chat_id = str(config.get("chat_id", ""))
        self.timeout = config.get("send_timeout", 10)
        self.retry_count = config.get("retry_count", 3)
        self.retry_delay = config.get("retry_delay", 5)
        self._base_url = API_BASE.format(token=self.token)

    def send_message(self, text: str, parse_mode: str = None) -> bool:
        if not self.token or not self.chat_id:
            log.debug("Token/Chat ID eksik — bildirim atlandı")
            return False
        url = f"{self._base_url}/sendMessage"
        payload = {"chat_id": self.chat_id, "text": text}
        if parse_mode:
            payload["parse_mode"] = parse_mode
        return self._send_with_retry(url, data=payload)

    def send_photo(self, photo_path: str, caption: str = "", delete_after: bool = True) -> bool:
        """
        Telegram'a fotoğraf gönder.

        Args:
            photo_path: JPEG dosyasının yolu
            caption: Fotoğraf altı yazısı
            delete_after: Gönderdikten sonra dosyayı sil mi?

        Returns:
            True = başarılı, False = başarısız
        """
        if not os.path.exists(photo_path):
            log.error(f"Fotoğraf bulunamadı: {photo_path}")
            return False

        url = f"{self._base_url}/sendPhoto"
        payload = {
            "chat_id": self.chat_id,
        }
        if caption:
            payload["caption"] = caption

        try:
            with open(photo_path, "rb") as photo_file:
                files = {"photo": photo_file}
                success = self._send_with_retry(url, data=payload, files=files)
        except IOError as e:
            log.error(f"Fotoğraf dosyası okunamadı: {e}")
            return False

        # Başarılı gönderimden sonra geçici dosyayı sil
        if success and delete_after:
            try:
                os.remove(photo_path)
                log.debug(f"Geçici dosya silindi: {photo_path}")
            except OSError:
                pass

        return success

    def send_document(self, doc_path: str, caption: str = "") -> bool:
        """Telegram'a dosya (dokuman) gönder."""
        if not os.path.exists(doc_path):
            log.error(f"Dosya bulunamadı: {doc_path}")
            return False

        url = f"{self._base_url}/sendDocument"
        payload = {"chat_id": self.chat_id}
        if caption:
            payload["caption"] = caption

        try:
            with open(doc_path, "rb") as f:
                files = {"document": f}
                return self._send_with_retry(url, data=payload, files=files)
        except IOError as e:
            log.error(f"Dosya okunamadı: {e}")
            return False

    def test_connection(self) -> bool:
        """
        Bot token'ı doğrula ve bağlantıyı test et.
        getMe API çağrısı yapar.
        """
        if not self.token:
            log.error("Telegram bot token boş!")
            return False

        try:
            url = f"{self._base_url}/getMe"
            resp = requests.get(url, timeout=self.timeout)
            data = resp.json()

            if data.get("ok"):
                bot_name = data["result"].get("username", "?")
                log.info(f"Telegram bağlantısı başarılı — Bot: @{bot_name}")
                return True
            else:
                log.error(f"Telegram doğrulama hatası: {data.get('description', 'bilinmiyor')}")
                return False
        except requests.RequestException as e:
            log.error(f"Telegram bağlantı hatası: {e}")
            return False

    def _send_with_retry(self, url: str, data: dict = None, files: dict = None) -> bool:
        """Retry mekanizmalı HTTP POST gönderimi."""
        if not self.token or not self.chat_id:
            return False
        for attempt in range(1, self.retry_count + 1):
            try:
                resp = requests.post(url, data=data, files=files, timeout=self.timeout)
                result = resp.json()
                if result.get("ok"):
                    log.debug(f"Telegram gönderim başarılı (deneme {attempt})")
                    return True
                else:
                    desc = result.get("description", "bilinmiyor")
                    log.warning(f"Telegram API hatası (deneme {attempt}): {desc}")
            except requests.Timeout:
                log.warning(f"Telegram timeout (deneme {attempt}/{self.retry_count})")
            except requests.ConnectionError:
                log.warning(f"Telegram bağlantı hatası (deneme {attempt}/{self.retry_count})")
            except requests.RequestException as e:
                log.warning(f"Telegram istek hatası (deneme {attempt}): {e}")
            if attempt < self.retry_count:
                time.sleep(self.retry_delay)
        log.error(f"Telegram gönderimi {self.retry_count} denemeden sonra başarısız")
        return False
