"""
Kanije Kalesi — Telegram Mesaj Formatlayıcı
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Her olay tipi için okunabilir, emoji'li Telegram mesajı oluşturur.
"""

from datetime import datetime
from core.event_bus import Event

# Uygulama versiyonu
VERSION = "1.0.0"


def format_event_message(event: Event, uptime_str: str = "") -> str:
    """
    Bir Event nesnesini Telegram mesaj metnine dönüştürür.

    Args:
        event: İşlenecek olay
        uptime_str: Uygulamanın çalışma süresi (opsiyonel)
    """
    formatters = {
        "login_success": _format_login_success,
        "login_failed": _format_login_failed,
        "screen_lock": _format_screen_lock,
        "screen_unlock": _format_screen_unlock,
        "system_boot": _format_system_boot,
        "system_sleep": _format_system_sleep,
        "system_wake": _format_system_wake,
        "usb_inserted": _format_usb_inserted,
        "usb_removed": _format_usb_removed,
        "internet_connected": _format_internet_connected,
    }

    formatter = formatters.get(event.event_type, _format_generic)
    body = formatter(event)

    # Alt bilgi
    footer = f"─────────────────────\n🏰 Kanije Kalesi v{VERSION}"
    if uptime_str:
        footer += f" · {uptime_str}"
    # Yaşam döngüsü mesajlarına /help ipucu ekle
    if event.event_type in ("system_boot", "internet_connected"):
        footer += "\n💡 Komutlar için /help"

    return f"{body}\n\n{footer}"


def format_status_message(system_info: dict) -> str:
    """
    /status komutu için sistem durum mesajı.
    """
    lines = [
        "━━━━━━━━━━━━━━━━━━━",
        "🏰 KANİJE KALESİ — DURUM",
        "━━━━━━━━━━━━━━━━━━━",
        "",
        f"🟢 Durum: Aktif",
        f"🖥️ Bilgisayar: {system_info.get('hostname', '?')}",
        f"⏱️ Uygulama Uptime: {system_info.get('app_uptime', '?')}",
        f"🔋 Sistem Uptime: {system_info.get('sys_uptime', '?')}",
        "",
        f"📊 CPU: %{system_info.get('cpu_percent', '?')}",
        f"💾 RAM: {system_info.get('ram_used', '?')} / {system_info.get('ram_total', '?')}",
        f"💽 Disk: {system_info.get('disk_used', '?')} / {system_info.get('disk_total', '?')}",
        "",
        f"📋 Toplam Olay: {system_info.get('total_events', 0)}",
        f"🚫 Atlanan (Rate Limit): {system_info.get('dropped_events', 0)}",
        f"⏰ Son Olay: {system_info.get('last_event', 'henüz yok')}",
        "",
        f"─────────────────────",
        f"🏰 Kanije Kalesi v{VERSION}",
    ]
    return "\n".join(lines)


def format_help_message() -> str:
    """/help komutu için komut listesi."""
    return (
        "🏰 *Kanije Kalesi — Komutlar*\n"
        "\n"
        "📋 *Bilgi:*\n"
        "/status — Sistem durumu (CPU, RAM, Disk, uptime)\n"
        "/events — Son 10 olay\n"
        "/config — Aktif yapılandırma\n"
        "/ping — Canlılık kontrolü\n"
        "\n"
        "📸 *Medya:*\n"
        "/photo — Anlık kamera fotoğrafı\n"
        "/screenshot — Anlık ekran görüntüsü\n"
        "\n"
        "🔧 *Kontrol:*\n"
        "/lock — Ekranı kilitle\n"
        "/restart — Bilgisayarı yeniden başlat (15sn onay)\n"
        "/shutdown — Bilgisayarı kapat (15sn onay)\n"
        "/cancel — Bekleyen restart/shutdown iptal\n"
        "\n"
        "/help — Bu mesaj"
    )


# ── Özel formatlayıcılar ──

def _format_login_success(event: Event) -> str:
    logon = event.extra.get("logon_type_str", "")
    logon_line = f"\n🔗 Giriş Tipi: {logon}" if logon else ""

    return (
        "━━━━━━━━━━━━━━━━━━━\n"
        "✅ BAŞARILI GİRİŞ\n"
        "━━━━━━━━━━━━━━━━━━━\n"
        "\n"
        f"🖥️ Bilgisayar: {event.hostname}\n"
        f"👤 Kullanıcı: {event.username}\n"
        f"⏰ Zaman: {event.timestamp:%Y-%m-%d %H:%M:%S}\n"
        f"🌐 IP: {event.ip_address or 'Yerel'}"
        f"{logon_line}"
    )


def _format_login_failed(event: Event) -> str:
    return (
        "━━━━━━━━━━━━━━━━━━━\n"
        "🚨 BAŞARISIZ GİRİŞ DENEMESİ\n"
        "━━━━━━━━━━━━━━━━━━━\n"
        "\n"
        f"🖥️ Bilgisayar: {event.hostname}\n"
        f"👤 Denenen Kullanıcı: {event.username}\n"
        f"⏰ Zaman: {event.timestamp:%Y-%m-%d %H:%M:%S}\n"
        f"🌐 IP: {event.ip_address or 'Yerel'}\n"
        f"📸 Kamera görüntüsü ektedir."
    )


def _format_screen_lock(event: Event) -> str:
    return (
        "🔒 Ekran kilitlendi\n"
        f"🖥️ {event.hostname} · ⏰ {event.timestamp:%H:%M:%S}"
    )


def _format_screen_unlock(event: Event) -> str:
    return (
        "🔓 Ekran kilidi açıldı\n"
        f"🖥️ {event.hostname} · 👤 {event.username} · ⏰ {event.timestamp:%H:%M:%S}"
    )


def _format_system_boot(event: Event) -> str:
    return (
        "━━━━━━━━━━━━━━━━━━━\n"
        "🖥️ SİSTEM BAŞLATILDI\n"
        "━━━━━━━━━━━━━━━━━━━\n"
        "\n"
        f"🖥️ Bilgisayar: {event.hostname}\n"
        f"⏰ Zaman: {event.timestamp:%Y-%m-%d %H:%M:%S}\n"
        f"🏰 Kanije Kalesi nöbete başladı."
    )


def _format_usb_inserted(event: Event) -> str:
    drive = event.extra.get("drive", "?")
    label = event.extra.get("label", "isimsiz")
    size = event.extra.get("size_gb", "?")
    fs = event.extra.get("filesystem", "?")

    return (
        "━━━━━━━━━━━━━━━━━━━\n"
        "🔌 YENİ USB CİHAZI\n"
        "━━━━━━━━━━━━━━━━━━━\n"
        "\n"
        f"💾 Sürücü: {drive}\n"
        f"🏷️ Etiket: {label}\n"
        f"📦 Boyut: {size}\n"
        f"📁 Dosya Sistemi: {fs}\n"
        f"⏰ Zaman: {datetime.now():%Y-%m-%d %H:%M:%S}"
    )


def _format_usb_removed(event: Event) -> str:
    drive = event.extra.get("drive", "?")
    label = event.extra.get("label", "bilinmiyor")

    return (
        "━━━━━━━━━━━━━━━━━━━\n"
        "⏏️ USB CİHAZI ÇIKARILDI\n"
        "━━━━━━━━━━━━━━━━━━━\n"
        "\n"
        f"💾 Sürücü: {drive}\n"
        f"🏷️ Etiket: {label}\n"
        f"⏰ Zaman: {datetime.now():%Y-%m-%d %H:%M:%S}"
    )


def _format_internet_connected(event: Event) -> str:
    pending = event.extra.get("pending_count", 0)
    pending_line = f"\n📨 Bekleyen {pending} bildirim gönderiliyor..." if pending > 0 else ""

    return (
        "━━━━━━━━━━━━━━━━━━━\n"
        "🌐 İNTERNET BAĞLANTISI KURULDU\n"
        "━━━━━━━━━━━━━━━━━━━\n"
        "\n"
        f"🖥️ Bilgisayar: {event.hostname}\n"
        f"⏰ Zaman: {event.timestamp:%Y-%m-%d %H:%M:%S}"
        f"{pending_line}"
    )


def _format_system_sleep(event: Event) -> str:
    return (
        "😴 Sistem uyku/hibernate moduna girdi\n"
        f"🖥️ {event.hostname} · ⏰ {event.timestamp:%H:%M:%S}"
    )


def _format_system_wake(event: Event) -> str:
    wake_type = event.extra.get("wake_type", "")
    wake_line = f" ({wake_type})" if wake_type else ""
    return (
        "☀️ Sistem uykudan uyandı"
        f"{wake_line}\n"
        f"🖥️ {event.hostname} · ⏰ {event.timestamp:%H:%M:%S}"
    )


def _format_generic(event: Event) -> str:
    return (
        f"📌 Olay: {event.event_type}\n"
        f"🖥️ {event.hostname} · ⏰ {event.timestamp:%H:%M:%S}"
    )
