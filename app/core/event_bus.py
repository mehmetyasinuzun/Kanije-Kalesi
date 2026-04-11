"""
Kanije Kalesi — Olay Yolu (Event Bus)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Thread-safe event queue. Listener'lar buraya event koyar,
ana döngü buradan alıp işler: aksiyon çalıştırır, Telegram'a gönderir.

Ayrıca rate limiter (flood koruması) içerir.
"""

import queue
import time
import threading
from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional
from collections import defaultdict


@dataclass
class Event:
    """Sistemde gerçekleşen bir olayı temsil eder."""

    event_type: str          # "login_success", "login_failed", "screen_unlock", "usb_inserted", "system_boot"
    timestamp: datetime = field(default_factory=datetime.now)
    username: str = ""
    ip_address: str = ""
    hostname: str = ""
    extra: dict = field(default_factory=dict)

    def __str__(self):
        return f"Event({self.event_type}, user={self.username}, time={self.timestamp:%H:%M:%S})"


class RateLimiter:
    """
    Flood koruması. Belirli bir sürede belirli sayıdan fazla olay gelirse engeller.
    Örnek: 1 dakikada 3'ten fazla başarısız giriş → 4. ve sonrası atlanır.
    """

    def __init__(self, max_per_minute: int = 10):
        self.max = max_per_minute
        self._events: dict[str, list[float]] = defaultdict(list)
        self._lock = threading.Lock()

    def is_limited(self, event_type: str) -> bool:
        """True dönerse bu olay atlanmalı (rate limit aşıldı)."""
        now = time.time()
        with self._lock:
            # Son 60 saniyedeki olayları tut, eskilerini sil
            self._events[event_type] = [
                t for t in self._events[event_type] if now - t < 60
            ]
            if len(self._events[event_type]) >= self.max:
                return True
            self._events[event_type].append(now)
            return False

    def get_count(self, event_type: str) -> int:
        """Son 1 dakikadaki olay sayısını döndürür."""
        now = time.time()
        with self._lock:
            return len([t for t in self._events[event_type] if now - t < 60])


class EventBus:
    """
    Thread-safe olay kuyruğu.

    Kullanım:
        bus = EventBus(max_per_minute=10)

        # Listener thread'inde:
        bus.put(Event(event_type="login_failed", username="admin"))

        # Ana döngüde:
        event = bus.get(timeout=1.0)
        if event:
            process(event)
    """

    def __init__(self, max_per_minute: int = 10):
        self._queue: queue.Queue[Event] = queue.Queue(maxsize=1000)
        self._limiter = RateLimiter(max_per_minute=max_per_minute)
        self._total_events = 0
        self._dropped_events = 0

    def put(self, event: Event) -> bool:
        """
        Olayı kuyruğa ekle. Rate limit aşıldıysa False döner (olay atlandı).
        """
        if self._limiter.is_limited(event.event_type):
            self._dropped_events += 1
            return False

        try:
            self._queue.put_nowait(event)
            self._total_events += 1
            return True
        except queue.Full:
            self._dropped_events += 1
            return False

    def get(self, timeout: float = 1.0) -> Optional[Event]:
        """
        Kuyruktan olay al. Yoksa timeout sonunda None döner.
        Bloklamaz — ana döngü bunu tekrar tekrar çağırır.
        """
        try:
            return self._queue.get(timeout=timeout)
        except queue.Empty:
            return None

    @property
    def pending(self) -> int:
        """Kuyrukta bekleyen olay sayısı."""
        return self._queue.qsize()

    @property
    def stats(self) -> dict:
        """İstatistikler."""
        return {
            "total_events": self._total_events,
            "dropped_events": self._dropped_events,
            "pending": self.pending,
        }
