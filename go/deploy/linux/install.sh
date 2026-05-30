#!/bin/bash
# Kanije Kalesi — Linux / Raspberry Pi kurulum
# Yerel binary varsa kullanır, yoksa GitHub Release'ten indirir.
#
# Tek satır (önerilen):
#   curl -fsSL https://raw.githubusercontent.com/mehmetyasinuzun/Kanije-Kalesi/master/go/deploy/linux/get.sh | sudo bash -s -- --token "TOKEN" --chat "CHATID"
#
# Repo içinden:
#   sudo bash install.sh --token "TOKEN" --chat "CHATID"
#   sudo bash install.sh --remove | --status

set -euo pipefail

REPO="mehmetyasinuzun/Kanije-Kalesi"
RAW="https://raw.githubusercontent.com/$REPO/master/go"
BINARY_NAME="kanije"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/kanije"
DATA_DIR="/var/lib/kanije"
SERVICE_FILE="/etc/systemd/system/kanije.service"
SERVICE_NAME="kanije-kalesi"
KANIJE_USER="kanije"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
info()  { echo -e "${CYAN}ℹ ${NC}$*"; }
ok()    { echo -e "${GREEN}✅${NC} $*"; }
warn()  { echo -e "${YELLOW}⚠ ${NC}$*"; }
error() { echo -e "${RED}❌${NC} $*" >&2; exit 1; }

# ---- argümanlar ----
REMOVE=false; STATUS=false; TOKEN=""; CHAT=""
while [[ $# -gt 0 ]]; do
    case "$1" in
        --remove) REMOVE=true; shift ;;
        --status) STATUS=true; shift ;;
        --token)  TOKEN="${2:-}"; shift 2 ;;
        --chat)   CHAT="${2:-}"; shift 2 ;;
        *) shift ;;
    esac
done

if $STATUS; then
    systemctl status "$SERVICE_NAME" --no-pager -l 2>/dev/null || warn "Servis kurulu değil"
    exit 0
fi

[[ $EUID -ne 0 ]] && error "Bu script root gerektirir. 'sudo' ile çalıştırın."

if $REMOVE; then
    systemctl stop "$SERVICE_NAME" 2>/dev/null || true
    systemctl disable "$SERVICE_NAME" 2>/dev/null || true
    rm -f "$SERVICE_FILE"; systemctl daemon-reload
    rm -f "$INSTALL_DIR/$BINARY_NAME"
    ok "Kaldırıldı. Config ve veriler korundu: $CONFIG_DIR, $DATA_DIR"
    exit 0
fi

# ---- binary'yi temin et ----
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH=""
for c in \
    "$(dirname "$SCRIPT_DIR")/kanije" \
    "./kanije" \
    "./dist/kanije-linux-amd64" \
    "./dist/kanije-linux-arm64" \
    "./dist/kanije-linux-arm"; do
    [[ -f "$c" ]] && { BINARY_PATH="$(realpath "$c")"; break; }
done

if [[ -z "$BINARY_PATH" ]]; then
    case "$(uname -m)" in
        x86_64|amd64)       ASSET="kanije-linux-amd64" ;;
        aarch64|arm64)      ASSET="kanije-linux-arm64" ;;
        armv7l|armv6l|arm)  ASSET="kanije-linux-arm" ;;
        *) error "Desteklenmeyen mimari: $(uname -m)" ;;
    esac
    info "Hazır binary indiriliyor ($ASSET)..."
    URL="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep -o "https://[^\"]*$ASSET" | head -1)"
    [[ -z "$URL" ]] && error "Release bulunamadı. Bir sürüm yayınlayın veya 'make build-linux' ile derleyin."
    TMPBIN="$(mktemp)"
    curl -fsSL "$URL" -o "$TMPBIN"
    BINARY_PATH="$TMPBIN"
    ok "İndirildi"
fi
info "Binary: $BINARY_PATH"

# ---- kullanıcı + dizinler ----
id "$KANIJE_USER" &>/dev/null || useradd -r -s /sbin/nologin -d "$DATA_DIR" -c "Kanije Kalesi daemon" "$KANIJE_USER"
mkdir -p "$CONFIG_DIR" "$DATA_DIR"
chown "$KANIJE_USER:$KANIJE_USER" "$DATA_DIR"
chmod 750 "$DATA_DIR" "$CONFIG_DIR"

# ---- binary kur ----
install -m 755 -o root -g root "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
ok "Binary kuruldu: $INSTALL_DIR/$BINARY_NAME"

# ---- config ----
if [[ ! -f "$CONFIG_DIR/config.toml" ]]; then
    EXAMPLE="$(dirname "$SCRIPT_DIR")/config.example.toml"
    if [[ -f "$EXAMPLE" ]]; then
        install -m 640 -o root -g "$KANIJE_USER" "$EXAMPLE" "$CONFIG_DIR/config.toml"
    elif curl -fsSL "$RAW/config.example.toml" -o "$CONFIG_DIR/config.toml" 2>/dev/null; then
        chmod 640 "$CONFIG_DIR/config.toml"; chown "root:$KANIJE_USER" "$CONFIG_DIR/config.toml"
    fi
fi

# ---- systemd service ----
SERVICE_SRC="$SCRIPT_DIR/kanije.service"
if [[ -f "$SERVICE_SRC" ]]; then
    install -m 644 "$SERVICE_SRC" "$SERVICE_FILE"
else
    curl -fsSL "$RAW/deploy/linux/kanije.service" -o "$SERVICE_FILE"
fi

# ---- secrets ----
SECRETS_FILE="$CONFIG_DIR/secrets.env"
if [[ -n "$TOKEN" && -n "$CHAT" ]]; then
    printf 'KANIJE_BOT_TOKEN=%s\nKANIJE_CHAT_ID=%s\n' "$TOKEN" "$CHAT" > "$SECRETS_FILE"
    ok "Telegram bilgileri kaydedildi"
elif [[ ! -f "$SECRETS_FILE" ]]; then
    cat > "$SECRETS_FILE" <<'EOF'
KANIJE_BOT_TOKEN=buraya_bot_token_yaz
KANIJE_CHAT_ID=buraya_chat_id_yaz
EOF
    warn "Telegram bilgilerini girin: $SECRETS_FILE"
fi
chmod 640 "$SECRETS_FILE"; chown "root:$KANIJE_USER" "$SECRETS_FILE"

# EnvironmentFile satırını etkinleştir
sed -i 's|^# EnvironmentFile=|EnvironmentFile=|' "$SERVICE_FILE"

systemctl daemon-reload
systemctl enable "$SERVICE_NAME" >/dev/null 2>&1
ok "Systemd servisi etkinleştirildi"

# ---- başlat ----
if [[ -n "$TOKEN" && -n "$CHAT" ]]; then
    systemctl restart "$SERVICE_NAME"
    ok "Servis başlatıldı — arka planda çalışıyor"
    echo ""
    ok "Kurulum tamamlandı! 🏰  Telegram'da /kurulum yazarak ayarlayın."
else
    echo ""
    ok "Kurulum tamamlandı!"
    echo "  1. Telegram bilgilerini girin:  sudo nano $SECRETS_FILE"
    echo "  2. Başlatın:                    sudo systemctl start $SERVICE_NAME"
fi
