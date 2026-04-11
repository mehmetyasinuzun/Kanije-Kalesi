"""
Kanije Kalesi — Loglama Sistemi
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Dosyaya + konsola loglama. Otomatik rotasyon (boyut aşılınca eski dosya yedeklenir).
"""

import logging
import sys
from logging.handlers import RotatingFileHandler
from pathlib import Path


def setup_logger(
    name: str = "kanije",
    log_file: str = "./kanije.log",
    level: str = "INFO",
    max_size_mb: int = 10,
    backup_count: int = 3,
    console_output: bool = True,
) -> logging.Logger:
    """
    Uygulamanın merkezi logger'ını oluşturur.

    Args:
        name: Logger adı
        log_file: Log dosyası yolu
        level: Log seviyesi (DEBUG, INFO, WARNING, ERROR)
        max_size_mb: Log dosyası max boyut (MB)
        backup_count: Tutulacak yedek log dosyası sayısı
        console_output: Konsola da yazdır mı?

    Returns:
        logging.Logger nesnesi
    """
    logger = logging.getLogger(name)

    # Zaten handler varsa tekrar ekleme (multiple init koruması)
    if logger.handlers:
        return logger

    logger.setLevel(getattr(logging, level.upper(), logging.INFO))

    # Format: zaman | seviye | modül | mesaj
    formatter = logging.Formatter(
        fmt="%(asctime)s | %(levelname)-7s | %(name)-12s | %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )

    # ── Dosya Handler (Rotating) ──
    log_path = Path(log_file)
    log_path.parent.mkdir(parents=True, exist_ok=True)

    file_handler = RotatingFileHandler(
        filename=str(log_path),
        maxBytes=max_size_mb * 1024 * 1024,
        backupCount=backup_count,
        encoding="utf-8",
    )
    file_handler.setFormatter(formatter)
    logger.addHandler(file_handler)

    # ── Konsol Handler ──
    if console_output:
        console_handler = logging.StreamHandler(sys.stdout)
        console_handler.setFormatter(formatter)
        logger.addHandler(console_handler)

    return logger


def get_logger(module_name: str = None) -> logging.Logger:
    """
    Alt modüller için child logger döndürür.

    Kullanım:
        from core.logger import get_logger
        log = get_logger("camera")
        log.info("Fotoğraf çekildi")

    Log çıktısı: 2026-03-22 16:45:32 | INFO    | kanije.camera | Fotoğraf çekildi
    """
    if module_name:
        return logging.getLogger(f"kanije.{module_name}")
    return logging.getLogger("kanije")
