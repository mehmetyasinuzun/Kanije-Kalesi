"""
Kanije Kalesi — Güç/Uyku Olayları Monitörü
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
WMI Win32_PowerManagementEvent ile uyku ve uyanma olaylarını algılar.

EventType değerleri:
    4  = Suspended (uyku / hibernate)
    7  = ResumeSuspend (uykudan/hibernateden uyanma)
    18 = ResumeAutomatic (zamanlanmış uyanma — genellikle sessiz)
"""

import threading
import socket
from core.logger import get_logger
from core.event_bus import EventBus, Event

log = get_logger("power_monitor")


class PowerMonitor:
    """
    WMI ile uyku ve uyanma olaylarını izler.
    Her olayda EventBus'a bildirim gönderir.
    """

    def __init__(self, event_bus: EventBus):
        self.bus = event_bus
        self._thread = None
        self._running = False

    def start(self):
        """Güç monitörünü arka plan thread'inde başlat."""
        self._running = True
        self._thread = threading.Thread(
            target=self._watch_loop,
            name="PowerMonitor",
            daemon=True,
        )
        self._thread.start()
        log.info("Güç/Uyku monitörü başlatıldı")

    def stop(self):
        """Monitörü durdur."""
        self._running = False
        log.info("Güç/Uyku monitörü durduruldu")

    def _watch_loop(self):
        """WMI ile güç olaylarını dinle."""
        # WMI her thread'de COM'u başlatmak zorunda
        try:
            import pythoncom
            pythoncom.CoInitialize()
        except Exception:
            pass

        try:
            import wmi
            c = wmi.WMI()
            watcher = c.Win32_PowerManagementEvent.watch_for()
            log.info("WMI güç olayı izleyici aktif")

            while self._running:
                try:
                    power_event = watcher(timeout_ms=2000)

                    if not power_event:
                        continue

                    event_type_id = getattr(power_event, "EventType", 0)

                    if event_type_id == 4:
                        # Uyku veya Hibernate
                        event = Event(
                            event_type="system_sleep",
                            hostname=socket.gethostname(),
                        )
                        self.bus.put(event)
                        log.info("Sistem uyku/hibernate moduna girdi")

                    elif event_type_id == 7:
                        # Manuel uyanma (kullanıcı bir tuşa bastı)
                        event = Event(
                            event_type="system_wake",
                            hostname=socket.gethostname(),
                            extra={"wake_type": "Manuel"},
                        )
                        self.bus.put(event)
                        log.info("Sistem uykudan uyandı (manuel)")

                    elif event_type_id == 18:
                        # Otomatik uyanma (zamanlayıcı)
                        event = Event(
                            event_type="system_wake",
                            hostname=socket.gethostname(),
                            extra={"wake_type": "Otomatik"},
                        )
                        self.bus.put(event)
                        log.info("Sistem otomatik olarak uyandı")

                except Exception as timeout_err:
                    if "Timed out" not in str(timeout_err):
                        log.warning(f"Güç izleyici uyarı: {timeout_err}")

        except ImportError:
            log.error("wmi paketi yüklü değil — pip install wmi")
        except Exception as e:
            log.error(f"Güç monitörü hatası: {e}")
        finally:
            try:
                import pythoncom
                pythoncom.CoUninitialize()
            except Exception:
                pass
