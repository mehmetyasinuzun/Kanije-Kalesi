"""
Kanije Kalesi — CLI Giriş Noktası
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Komut satırı arayüzü.

Kullanım:
    python kanije.py start          — Uygulamayı başlat
    python kanije.py test           — Telegram bağlantısını test et
    python kanije.py status         — Çalışma durumunu kontrol et
    python kanije.py --config X     — Özel config dosyası kullan
    python kanije.py --help         — Yardım
"""

import sys
import os
import argparse

# Modül yolunu ayarla (app/ dışından çalıştırılabilirlik)
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))


def main():
    parser = argparse.ArgumentParser(
        prog="kanije",
        description="🏰 Kanije Kalesi — Bilgisayar Güvenlik İzleme Yazılımı",
        epilog="Detaylı bilgi: https://github.com/mehmetyasinuzun/Kanije-Kalesi",
    )

    parser.add_argument(
        "command",
        choices=["start", "test", "status"],
        help="start: Uygulamayı başlat | test: Telegram testi | status: Durum kontrolü",
    )
    parser.add_argument(
        "--config", "-c",
        default=None,
        help="Özel config.yaml yolu (varsayılan: çalışma dizini)",
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Ayrıntılı çıktı (DEBUG seviyesi)",
    )

    args = parser.parse_args()

    # Verbose modda config override
    config_override = args.config

    if args.command == "start":
        _cmd_start(config_override, args.verbose)
    elif args.command == "test":
        _cmd_test(config_override)
    elif args.command == "status":
        _cmd_status()


def _cmd_start(config_path, verbose):
    """Uygulamayı başlat."""

    # Yönetici yetkisi kontrolü (Windows)
    if sys.platform == "win32":
        import ctypes
        if not ctypes.windll.shell32.IsUserAnAdmin():
            print("⚠️  UYARI: Yönetici yetkisi olmadan çalışıyorsun!")
            print("   → Security Event Log okunamaz (login/kilit olayları algılanmaz)")
            print("   → PowerShell'i SAĞ TIK → 'Yönetici olarak çalıştır' ile aç")
            print()
            answer = input("   Yine de devam etmek istiyor musun? (e/h): ").strip().lower()
            if answer != "e":
                print("   Çıkılıyor. Yönetici PowerShell ile tekrar dene.")
                sys.exit(1)
            print()

    # .env yoksa oluştur, token eksikse kısa uyarı ver
    try:
        from core.setup_wizard import ensure_env, get_credentials
        ensure_env()
        token, chat_id = get_credentials()
        if not token or not chat_id:
            print("⚠️  Token/Chat ID eksik — tray ikonundan ⚙️ Telegram Ayarları ile gir.")
            print()
    except Exception:
        pass

    print("🏰 Kanije Kalesi başlatılıyor...")
    print("   Durdurmak için Ctrl+C")
    print()

    from core.app import KanijeApp

    app = KanijeApp(config_path=config_path)

    if verbose:
        import logging
        logging.getLogger("kanije").setLevel(logging.DEBUG)

    try:
        app.start()
    except KeyboardInterrupt:
        app.shutdown()
    except Exception as e:
        print(f"\n❌ Kritik hata: {e}")
        app.shutdown()
        sys.exit(1)


def _cmd_test(config_path):
    """Telegram bağlantısını test et."""
    print("🔌 Telegram bağlantısı test ediliyor...")
    print()

    from core.app import KanijeApp

    app = KanijeApp(config_path=config_path)
    success = app.test_telegram()

    if success:
        print("\n✅ Her şey çalışıyor! Telegram'ını kontrol et.")
    else:
        print("\n❌ Test başarısız. Kontrol et:")
        print("   1. config.yaml'da bot_token doğru mu?")
        print("   2. config.yaml'da chat_id doğru mu?")
        print("   3. İnternet bağlantın var mı?")
        print("   4. Bot'u Telegram'da /start ile başlattın mı?")

    sys.exit(0 if success else 1)


def _cmd_status():
    """Çalışma durumunu kontrol et."""
    print("🏰 Kanije Kalesi — Durum Kontrolü")
    print()

    # Config var mı?
    from pathlib import Path
    config_path = Path("config.yaml")
    if config_path.exists():
        print(f"   ✅ Config dosyası: {config_path.absolute()}")
    else:
        print(f"   ❌ Config dosyası bulunamadı: {config_path.absolute()}")
        print(f"      config.yaml oluştur ve token/chat_id bilgilerini gir.")

    # Log var mı?
    log_path = Path("kanije.log")
    if log_path.exists():
        size_kb = log_path.stat().st_size / 1024
        print(f"   📋 Log dosyası: {log_path.absolute()} ({size_kb:.1f} KB)")
    else:
        print(f"   📋 Log dosyası: henüz oluşturulmamış")

    # Python versiyon
    print(f"   🐍 Python: {sys.version.split()[0]}")

    # Bağımlılık kontrolü
    deps = {
        "cv2": "opencv-python-headless",
        "mss": "mss",
        "requests": "requests",
        "yaml": "pyyaml",
        "psutil": "psutil",
        "pystray": "pystray",
        "PIL": "Pillow",
    }

    print(f"\n   Bağımlılıklar:")
    for module, package in deps.items():
        try:
            __import__(module)
            print(f"   ✅ {package}")
        except ImportError:
            print(f"   ❌ {package} — pip install {package}")

    # Windows-only bağımlılıklar
    if sys.platform == "win32":
        win_deps = {"win32evtlog": "pywin32", "wmi": "wmi"}
        for module, package in win_deps.items():
            try:
                __import__(module)
                print(f"   ✅ {package}")
            except ImportError:
                print(f"   ❌ {package} — pip install {package}")


if __name__ == "__main__":
    main()
