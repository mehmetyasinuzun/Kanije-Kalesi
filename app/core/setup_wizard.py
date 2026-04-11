"""
Kanije Kalesi — .env Yönetimi ve Telegram Kurulum Penceresi
"""

import os
import sys
import subprocess
from pathlib import Path

ENV_PATH = Path(__file__).parent.parent / ".env"


# ─── .env okuma / yazma ──────────────────────────────────────────

def ensure_env():
    if not ENV_PATH.exists():
        ENV_PATH.write_text(
            "# Kanije Kalesi — Telegram Kimlik Bilgileri\n"
            "KANIJE_BOT_TOKEN=\n"
            "KANIJE_CHAT_ID=\n",
            encoding="utf-8",
        )
    try:
        from dotenv import load_dotenv
        load_dotenv(ENV_PATH, override=False)
    except ImportError:
        _manual_load()


def _manual_load():
    if not ENV_PATH.exists():
        return
    for line in ENV_PATH.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        k, _, v = line.partition("=")
        k, v = k.strip(), v.strip()
        if k and v and not os.environ.get(k):
            os.environ[k] = v


def get_credentials() -> tuple:
    ensure_env()
    return (
        os.environ.get("KANIJE_BOT_TOKEN", "").strip(),
        os.environ.get("KANIJE_CHAT_ID", "").strip(),
    )


def save_credentials(token: str, chat_id: str):
    ENV_PATH.write_text(
        "# Kanije Kalesi — Telegram Kimlik Bilgileri\n"
        f"KANIJE_BOT_TOKEN={token.strip()}\n"
        f"KANIJE_CHAT_ID={chat_id.strip()}\n",
        encoding="utf-8",
    )
    os.environ["KANIJE_BOT_TOKEN"] = token.strip()
    os.environ["KANIJE_CHAT_ID"] = chat_id.strip()


# ─── Tkinter kurulum penceresi ───────────────────────────────────

def open_settings_dialog():
    """Kompakt, sabit boyutlu Telegram ayarları kartı."""
    import threading
    try:
        import tkinter as tk
        from tkinter import ttk
    except ImportError:
        print("[HATA] tkinter bulunamadı.")
        return

    token_now, chat_now = get_credentials()

    # ── Palet ──
    BG      = "#1c1c2a"
    CARD    = "#252535"
    ACCENT  = "#6b8cff"
    TEXT    = "#d0d6f0"
    MUTED   = "#5a5f7a"
    OK      = "#7ecfa0"
    ERR     = "#e07a7a"
    EBORDER = "#3a3a55"
    FONT    = ("Segoe UI", 9)
    FONT_B  = ("Segoe UI", 9, "bold")

    W, H = 370, 210
    root = tk.Tk()
    root.title("Telegram Ayarları")
    root.configure(bg=BG)
    root.resizable(False, False)
    root.attributes("-topmost", True)

    root.update_idletasks()
    sx, sy = root.winfo_screenwidth(), root.winfo_screenheight()
    root.geometry(f"{W}x{H}+{(sx-W)//2}+{(sy-H)//2}")

    # ── ttk stil ──
    style = ttk.Style(root)
    style.theme_use("clam")
    style.configure("E.TEntry",
        fieldbackground=CARD, foreground=TEXT,
        insertcolor=TEXT, borderwidth=1,
        relief="solid", padding=(6, 4), font=FONT)
    style.configure("A.TButton",
        background=ACCENT, foreground="#ffffff",
        font=FONT_B, borderwidth=0, relief="flat", padding=(10, 5))
    style.map("A.TButton",
        background=[("active", "#5272df")])
    style.configure("F.TButton",
        background=CARD, foreground=MUTED,
        font=FONT, borderwidth=0, relief="flat", padding=(10, 5))
    style.map("F.TButton",
        background=[("active", "#2e2e45")],
        foreground=[("active", TEXT)])

    # ── Kart çerçevesi ──
    card = tk.Frame(root, bg=CARD, bd=0)
    card.place(x=16, y=16, width=W-32, height=H-32)

    # Başlık çizgisi
    tk.Frame(card, bg=ACCENT, height=2).pack(fill="x", side="top")

    body = tk.Frame(card, bg=CARD, padx=18, pady=12)
    body.pack(fill="both", expand=True)
    body.columnconfigure(1, weight=1)

    def row(label, var, r, show=""):
        tk.Label(body, text=label, bg=CARD, fg=MUTED,
                 font=FONT, anchor="w", width=9
                 ).grid(row=r, column=0, sticky="w", pady=(0, 10))
        e = ttk.Entry(body, textvariable=var, style="E.TEntry", show=show)
        e.grid(row=r, column=1, sticky="ew", pady=(0, 10))
        return e

    t_var = tk.StringVar(value=token_now)
    c_var = tk.StringVar(value=chat_now)
    first = row("Token", t_var, 0)
    row("Chat ID", c_var, 1)

    # Durum satırı
    status_var = tk.StringVar()
    status_lbl = tk.Label(body, textvariable=status_var,
                          bg=CARD, fg=MUTED, font=("Segoe UI", 8), anchor="w")
    status_lbl.grid(row=2, column=0, columnspan=2, sticky="w")

    def _set(msg, color=MUTED):
        status_var.set(msg)
        status_lbl.configure(fg=color)

    # ── Test (non-blocking) ──
    def _test():
        t = t_var.get().strip()
        if not t:
            _set("Token boş!", ERR); return
        _set("Bağlanıyor...", MUTED)

        def _run():
            try:
                import requests
                r = requests.get(
                    f"https://api.telegram.org/bot{t}/getMe", timeout=8)
                d = r.json()
                if d.get("ok"):
                    name = d["result"].get("username", "?")
                    root.after(0, _set, f"✓  @{name}", OK)
                else:
                    root.after(0, _set, d.get("description", "Hata"), ERR)
            except Exception as e:
                root.after(0, _set, str(e)[:50], ERR)

        threading.Thread(target=_run, daemon=True).start()

    # ── Kaydet ──
    def _save():
        t, c = t_var.get().strip(), c_var.get().strip()
        if not t or not c:
            _set("Her iki alan doldurulmalı!", ERR); return
        save_credentials(t, c)
        _set("Kaydedildi — yeniden başlatınca aktif olur", OK)

    # ── Butonlar ──
    btn = tk.Frame(body, bg=CARD)
    btn.grid(row=3, column=0, columnspan=2, sticky="e", pady=(8, 0))

    ttk.Button(btn, text="Test",  command=_test, style="F.TButton").pack(side="left", padx=(0, 6))
    ttk.Button(btn, text="Kaydet", command=_save, style="A.TButton").pack(side="left", padx=(0, 6))
    ttk.Button(btn, text="Kapat",  command=root.destroy, style="F.TButton").pack(side="left")

    first.focus()
    root.mainloop()


# ─── Tray'den güvenle çağır ──────────────────────────────────────

def open_settings_subprocess():
    app_dir = str(Path(__file__).parent.parent)
    subprocess.Popen(
        [sys.executable, "-c",
         f"import sys; sys.path.insert(0, {repr(app_dir)}); "
         f"from core.setup_wizard import open_settings_dialog; open_settings_dialog()"],
        creationflags=0x00000008,
    )


if __name__ == "__main__":
    open_settings_dialog()
