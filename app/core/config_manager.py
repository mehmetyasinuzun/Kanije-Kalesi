"""
Kanije Kalesi — Yapılandırma Yöneticisi
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
config.yaml dosyasını okur, doğrular ve ortam değişkenlerini çözümler.
Telegram kimlik bilgileri .env dosyasından yüklenir (config.yaml'a dokunmana gerek yok).
"""

import os
import sys
import yaml
from pathlib import Path
from copy import deepcopy

# .env otomatik yükle (varsa)
try:
    from dotenv import load_dotenv as _load_dotenv
    _env_path = Path(__file__).parent.parent / ".env"
    _load_dotenv(_env_path, override=False)  # Zaten set edilmişleri ezmez
except ImportError:
    pass  # python-dotenv yoksa ENV: mekanizması yine de çalışır

# ── Varsayılan yapılandırma (config.yaml yoksa veya eksik alan varsa) ──
DEFAULTS = {
    "telegram": {
        # bot_token ve chat_id buraya yazma — app/.env dosyasından okunuyor
        "send_timeout": 10,
        "retry_count": 3,
        "retry_delay": 5,
    },
    "triggers": {
        "login_success": {
            "enabled": True,
            "capture_camera": False,
            "capture_screenshot": False,
        },
        "login_failed": {
            "enabled": True,
            "capture_camera": True,
            "capture_screenshot": False,
            "max_photos_per_minute": 3,
        },
        "screen_lock": {
            "enabled": True,
        },
        "screen_unlock": {
            "enabled": True,
            "capture_camera": False,
        },
        "system_boot": {
            "enabled": True,
        },
        "usb_inserted": {
            "enabled": True,
        },
        "usb_removed": {
            "enabled": True,
        },
        "system_sleep": {
            "enabled": True,
        },
        "system_wake": {
            "enabled": True,
        },
    },
    "camera": {
        "device_index": 0,
        "resolution": [640, 480],
        "warmup_frames": 5,
        "jpeg_quality": 85,
        "save_local": False,
        "local_path": "./captures/",
    },
    "screenshot": {
        "jpeg_quality": 70,
        "save_local": False,
        "local_path": "./captures/",
    },
    "heartbeat": {
        "enabled": True,
        "interval_hours": 6,
        "include_uptime": True,
        "include_disk_usage": True,
    },
    "logging": {
        "level": "INFO",
        "file": "./kanije.log",
        "max_size_mb": 10,
        "backup_count": 3,
        "console_output": True,
    },
    "security": {
        "delete_photos_after_send": True,
        "max_events_per_minute": 10,
        "allowed_chat_ids": [],
    },
    "tray": {
        "enabled": True,
        "show_notifications": True,
    },
}


def _resolve_env(value: str) -> str:
    """
    'ENV:VARIABLE_NAME' formatındaki değeri ortam değişkeninden çözümle.
    Örnek: 'ENV:GHOSTGUARD_BOT_TOKEN' → os.environ['GHOSTGUARD_BOT_TOKEN']
    """
    if isinstance(value, str) and value.startswith("ENV:"):
        env_key = value[4:]
        env_val = os.environ.get(env_key, "")
        if not env_val:
            print(f"[UYARI] Ortam değişkeni bulunamadı: {env_key}")
        return env_val
    return value


def _deep_merge(base: dict, override: dict) -> dict:
    """
    İki dict'i derin birleştir. override'daki değerler base'i ezer.
    Eksik anahtarlar base'den gelir.
    """
    result = deepcopy(base)
    for key, value in override.items():
        if key in result and isinstance(result[key], dict) and isinstance(value, dict):
            result[key] = _deep_merge(result[key], value)
        else:
            result[key] = deepcopy(value)
    return result


def _find_config_path() -> Path:
    """
    Config dosyasını şu sırayla arar:
    1. Çalışma dizini
    2. Exe/script'in yanı
    3. Kullanıcının ev dizini (.kanije/)
    """
    candidates = [
        Path.cwd() / "config.yaml",
        Path(sys.argv[0]).parent / "config.yaml",
        Path.home() / ".kanije" / "config.yaml",
    ]
    for path in candidates:
        if path.exists():
            return path
    return candidates[0]  # Varsayılan: çalışma dizini


class ConfigManager:
    """
    YAML yapılandırma dosyasını yönetir.

    Kullanım:
        config = ConfigManager()
        token = config.get("telegram.bot_token")
        camera_on = config.get("triggers.login_failed.capture_camera")
    """

    def __init__(self, config_path: str = None):
        if config_path:
            self.path = Path(config_path)
        else:
            self.path = _find_config_path()

        self.data = deepcopy(DEFAULTS)

        if self.path.exists():
            self._load()
        else:
            print(f"[BİLGİ] Config dosyası bulunamadı: {self.path}")
            print(f"[BİLGİ] Varsayılan ayarlar kullanılıyor.")

        # Ortam değişkenlerini çözümle
        self._resolve_all_env(self.data)

    def _load(self):
        """YAML dosyasını oku ve varsayılanlarla birleştir."""
        try:
            with open(self.path, "r", encoding="utf-8") as f:
                user_config = yaml.safe_load(f) or {}
            self.data = _deep_merge(DEFAULTS, user_config)
        except yaml.YAMLError as e:
            print(f"[HATA] Config dosyası okunamadı: {e}")
            print(f"[BİLGİ] Varsayılan ayarlar kullanılıyor.")

    def _resolve_all_env(self, d: dict):
        """Dict içindeki tüm ENV: değerlerini çözümle (recursive)."""
        for key, value in d.items():
            if isinstance(value, dict):
                self._resolve_all_env(value)
            elif isinstance(value, str):
                d[key] = _resolve_env(value)

    def get(self, dotted_key: str, default=None):
        """
        Noktalı anahtar ile değer oku.
        Örnek: config.get("telegram.bot_token")
        """
        keys = dotted_key.split(".")
        current = self.data
        for k in keys:
            if isinstance(current, dict) and k in current:
                current = current[k]
            else:
                return default
        return current

    def get_bot_token(self) -> str:
        """
        Telegram bot token'u döndür.
        Öncelik: .env (KANIJE_BOT_TOKEN) → config.yaml telegram.bot_token
        """
        return (
            os.environ.get("KANIJE_BOT_TOKEN", "").strip()
            or self.get("telegram.bot_token", "")
        )

    def get_chat_id(self) -> str:
        """
        Telegram chat ID'döndür.
        Öncelik: .env (KANIJE_CHAT_ID) → config.yaml telegram.chat_id
        """
        return (
            os.environ.get("KANIJE_CHAT_ID", "").strip()
            or self.get("telegram.chat_id", "")
        )

    def validate(self) -> list:
        """
        Kritik ayarları doğrula. Hata listesi döndürür (boşsa sorun yok).
        """
        errors = []

        if not self.get_bot_token():
            errors.append("Bot Token eksik — app/.env dosyasına KANIJE_BOT_TOKEN ekle")

        if not self.get_chat_id():
            errors.append("Chat ID eksik — app/.env dosyasına KANIJE_CHAT_ID ekle")

        return errors

    def get_trigger(self, trigger_name: str) -> dict:
        """Bir tetikleyicinin tüm ayarlarını döndürür."""
        return self.get(f"triggers.{trigger_name}", {})

    def is_trigger_enabled(self, trigger_name: str) -> bool:
        """Bir tetikleyicinin açık olup olmadığını kontrol et."""
        return self.get(f"triggers.{trigger_name}.enabled", False)

    def dump_safe(self) -> dict:
        """Config'i token'lar gizlenmiş şekilde döndürür (loglama/debug için)."""
        import copy
        safe = copy.deepcopy(self.data)
        if "telegram" in safe and "bot_token" in safe["telegram"]:
            token = safe["telegram"]["bot_token"]
            if token and len(token) > 10:
                safe["telegram"]["bot_token"] = token[:6] + "..." + token[-4:]
            else:
                safe["telegram"]["bot_token"] = "***"
        return safe

    def __repr__(self):
        return f"ConfigManager(path={self.path}, keys={list(self.data.keys())})"
