"""
Kanije Kalesi — Windows Event Log Listener
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Windows'un Security Event Log'unu dinler.
Oturum açma, yanlış şifre, kilit açma gibi olayları yakalar.

Event ID'ler:
    4624 → Başarılı oturum açma
    4625 → Başarısız oturum denemesi
    4800 → Ekran kilitlendi
    4801 → Ekran kilidi açıldı
"""

import threading
import xml.etree.ElementTree as ET
from datetime import datetime

from core.logger import get_logger
from core.event_bus import EventBus, Event

log = get_logger("win_event")

# Windows Event ID → Kanije olay tipi eşleştirmesi
EVENT_MAP = {
    4624: "login_success",
    4625: "login_failed",
    4800: "screen_lock",
    4801: "screen_unlock",
}

# Logon Type açıklamaları (4624/4625 olaylarında ek bilgi)
LOGON_TYPES = {
    2: "Yerel (Konsol)",
    3: "Ağ (SMB/RPC)",
    7: "Kilit Açma",
    10: "Uzak Masaüstü (RDP)",
    11: "Önbellekli Kimlik",
}


def _parse_event_xml(xml_str: str) -> dict:
    """
    Windows Event Log XML çıktısını parse eder.
    Kullanıcı adı, IP, logon type gibi bilgileri çıkarır.
    """
    info = {
        "username": "",
        "ip_address": "",
        "logon_type": 0,
        "logon_type_str": "",
        "domain": "",
        "event_id": 0,
    }

    try:
        root = ET.fromstring(xml_str)
        ns = {"e": "http://schemas.microsoft.com/win/2004/08/events/event"}

        # Event ID
        system = root.find("e:System", ns)
        if system is not None:
            eid_el = system.find("e:EventID", ns)
            if eid_el is not None and eid_el.text:
                info["event_id"] = int(eid_el.text)

        # EventData alanları
        event_data = root.find("e:EventData", ns)
        if event_data is not None:
            for data in event_data.findall("e:Data", ns):
                name = data.get("Name", "")
                value = data.text or ""

                if name == "TargetUserName":
                    info["username"] = value
                elif name == "IpAddress":
                    info["ip_address"] = value if value != "-" else "Yerel"
                elif name == "LogonType":
                    try:
                        lt = int(value)
                        info["logon_type"] = lt
                        info["logon_type_str"] = LOGON_TYPES.get(lt, f"Tip {lt}")
                    except ValueError:
                        pass
                elif name == "TargetDomainName":
                    info["domain"] = value
    except ET.ParseError:
        log.warning("Event XML parse hatası")

    return info


class WindowsEventListener:
    """
    Windows Security Event Log'u asenkron dinler.
    Tespit edilen olayları EventBus'a gönderir.
    """

    def __init__(self, event_bus: EventBus, config: dict):
        self.bus = event_bus
        self.config = config
        self._thread = None
        self._running = False
        self._hostname = ""
        # Çift olay engelleme: {("event_type", "username"): timestamp}
        self._last_events = {}
        self._dedup_window = 3  # Aynı olay 3 saniye içinde tekrar gelirse atla

    def start(self):
        """Listener'ı arka plan thread'inde başlat."""
        import socket
        self._hostname = socket.gethostname()
        self._running = True
        self._thread = threading.Thread(
            target=self._listen_loop,
            name="WinEventListener",
            daemon=True,
        )
        self._thread.start()
        log.info("Windows Event Listener başlatıldı")

    def stop(self):
        """Listener'ı durdur."""
        self._running = False
        log.info("Windows Event Listener durduruldu")

    def _listen_loop(self):
        """
        Ana dinleme döngüsü.
        win32evtlog.EvtSubscribe ile asenkron event alır.
        """
        try:
            import win32evtlog

            # Security log'a abone ol — yalnızca ilgili Event ID'leri filtrele
            query = (
                "<QueryList>"
                "  <Query Id='0' Path='Security'>"
                "    <Select Path='Security'>"
                "      *[System[(EventID=4624 or EventID=4625 or EventID=4800 or EventID=4801)]]"
                "    </Select>"
                "  </Query>"
                "</QueryList>"
            )

            # Yeni olayları dinle (eski olayları atla)
            handle = win32evtlog.EvtSubscribe(
                "Security",
                win32evtlog.EvtSubscribeToFutureEvents,
                Query=query,
                Callback=self._on_event,
            )

            log.info("Security Event Log aboneliği aktif")

            # Thread canlı kalsın
            while self._running:
                import time
                time.sleep(1)

        except ImportError:
            log.error("pywin32 yüklü değil — pip install pywin32")
        except Exception as e:
            log.error(f"Event Log dinleme hatası: {e}")

    def _is_duplicate(self, event_type: str, username: str) -> bool:
        """Aynı olay kısa sürede tekrar geldiyse True döner."""
        key = (event_type, username.upper())
        now = datetime.now()

        last_time = self._last_events.get(key)
        if last_time and (now - last_time).total_seconds() < self._dedup_window:
            log.debug(f"Çift olay engellendi: {event_type}/{username} ({self._dedup_window}sn pencere)")
            return True

        self._last_events[key] = now

        # Eski kayıtları temizle (bellek sızıntısı engeli)
        expired = [k for k, v in self._last_events.items() if (now - v).total_seconds() > 60]
        for k in expired:
            del self._last_events[k]

        return False

    def _on_event(self, action, context, event_handle):
        """
        Windows Event Log callback'i.
        Her yeni olay geldiğinde bu fonksiyon çağrılır.
        """
        import win32evtlog

        if action != win32evtlog.EvtSubscribeActionDeliver:
            return

        try:
            xml_str = win32evtlog.EvtRender(event_handle, win32evtlog.EvtRenderEventXml)
            info = _parse_event_xml(xml_str)

            event_id = info["event_id"]
            event_type = EVENT_MAP.get(event_id)

            if not event_type:
                return

            # Sistem hesaplarını filtrele (SYSTEM, DWM, UMFD gibi)
            username = info["username"]
            skip_users = {"SYSTEM", "DWM-1", "DWM-2", "DWM-3", "UMFD-0", "UMFD-1",
                          "UMFD-2", "UMFD-3", "LOCAL SERVICE", "NETWORK SERVICE",
                          "ANONYMOUS LOGON", ""}
            if username.upper() in skip_users:
                return

            # Gürültülü logon type'ları filtrele:
            # 0=bilinmiyor, 5=servis, 11=önbellekli (çift tetikleme kaynağı)
            if info["logon_type"] in (0, 5, 11):
                return

            # Çift olay kontrolü (aynı kullanıcı + aynı olay 10sn içinde)
            if self._is_duplicate(event_type, username):
                return

            event = Event(
                event_type=event_type,
                username=username,
                ip_address=info["ip_address"],
                hostname=self._hostname,
                extra={
                    "event_id": event_id,
                    "logon_type": info["logon_type"],
                    "logon_type_str": info["logon_type_str"],
                    "domain": info["domain"],
                },
            )

            self.bus.put(event)
            log.info(f"Olay yakalandı: {event}")

        except Exception as e:
            log.error(f"Event işleme hatası: {e}")
