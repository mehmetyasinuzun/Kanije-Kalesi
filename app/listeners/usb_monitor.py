"""
Kanije Kalesi — USB Takma/Çıkarma Algılayıcı
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
WMI ile USB cihazı takma VE çıkarma olaylarını algılar.
Her USB olayında EventBus'a bildirim gönderir.

EventType değerleri:
    1 = Configuration Changed
    2 = Device Arrival (takma)
    3 = Device Removal (çıkarma)
    4 = Docking
"""

import threading
from core.logger import get_logger
from core.event_bus import EventBus, Event

log = get_logger("usb_monitor")


class USBMonitor:
    """
    WMI (Windows Management Instrumentation) kullanarak USB olaylarını izler.
    Hem takma hem çıkarma algılar.
    """

    def __init__(self, event_bus: EventBus, config: dict):
        self.bus = event_bus
        self.config = config
        self._threads = []
        self._running = False
        self._known_drives = {}  # drive_name → label (çıkarmada bilgi için)

    def start(self):
        """USB monitörü arka plan thread'lerinde başlat."""
        self._running = True

        # Takma watcher
        t_insert = threading.Thread(
            target=self._watch_loop,
            args=(2, "usb_inserted"),
            name="USBMonitor-Insert",
            daemon=True,
        )
        t_insert.start()
        self._threads.append(t_insert)

        # Çıkarma watcher
        t_remove = threading.Thread(
            target=self._watch_loop,
            args=(3, "usb_removed"),
            name="USBMonitor-Remove",
            daemon=True,
        )
        t_remove.start()
        self._threads.append(t_remove)

        log.info("USB Monitor başlatıldı (takma + çıkarma)")

    def stop(self):
        """USB monitörü durdur."""
        self._running = False
        log.info("USB Monitor durduruldu")

    def _watch_loop(self, event_type_id: int, event_name: str):
        """WMI ile USB olaylarını dinle."""
        try:
            # WMI her thread'de COM'u başlatmak zorunda
            import pythoncom
            pythoncom.CoInitialize()
        except Exception:
            pass

        try:
            import wmi

            c = wmi.WMI()
            watcher = c.Win32_VolumeChangeEvent.watch_for(EventType=event_type_id)

            action = "takma" if event_type_id == 2 else "çıkarma"
            log.info(f"WMI USB watcher aktif ({action})")

            while self._running:
                try:
                    usb_event = watcher(timeout_ms=2000)

                    if usb_event:
                        drive_name = getattr(usb_event, "DriveName", "bilinmiyor")

                        if event_type_id == 2:  # Takma
                            extra_info = self._get_usb_details(c, drive_name)
                            # Çıkarmada bilgi vermek için kaydet
                            self._known_drives[drive_name] = extra_info.get("label", "")

                            event = Event(
                                event_type=event_name,
                                extra={
                                    "drive": drive_name,
                                    "label": extra_info.get("label", ""),
                                    "size_gb": extra_info.get("size_gb", ""),
                                    "filesystem": extra_info.get("filesystem", ""),
                                },
                            )
                            self.bus.put(event)
                            log.info(f"USB takıldı: {drive_name} — {extra_info.get('label', 'isimsiz')}")

                        else:  # Çıkarma
                            label = self._known_drives.pop(drive_name, "bilinmiyor")
                            event = Event(
                                event_type=event_name,
                                extra={
                                    "drive": drive_name,
                                    "label": label,
                                },
                            )
                            self.bus.put(event)
                            log.info(f"USB çıkarıldı: {drive_name} — {label}")

                except Exception as timeout_err:
                    if "Timed out" not in str(timeout_err):
                        log.warning(f"USB watcher uyarı: {timeout_err}")

        except ImportError:
            log.error("wmi paketi yüklü değil — pip install wmi")
        except Exception as e:
            log.error(f"USB Monitor hatası: {e}")

    def _get_usb_details(self, wmi_conn, drive_name: str) -> dict:
        """Takılan USB hakkında ek bilgi topla."""
        details = {"label": "", "size_gb": "", "filesystem": ""}
        try:
            for disk in wmi_conn.Win32_LogicalDisk():
                if disk.DeviceID and disk.DeviceID.upper() == drive_name.upper():
                    details["label"] = disk.VolumeName or "isimsiz"
                    details["filesystem"] = disk.FileSystem or "bilinmiyor"
                    if disk.Size:
                        size_gb = int(disk.Size) / (1024 ** 3)
                        details["size_gb"] = f"{size_gb:.1f} GB"
                    break
        except Exception:
            pass
        return details
