#!/bin/bash
# Kanije Kalesi Sentinel — Linux/Raspberry Pi Kurulum Scripti
# Desteklenen: Ubuntu, Debian, Raspberry Pi OS, Arch Linux
#
# Kullanım:
#   sudo bash install.sh
#   sudo bash install.sh --remove
#   sudo bash install.sh --status

set -euo pipefail

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

# ---- Argument parsing ----
REMOVE=false
STATUS=false
for arg in "$@"; do
    case $arg in
        --remove) REMOVE=true ;;
        --status) STATUS=true ;;
    esac
done

# ---- Status check ----
if $STATUS; then
    if systemctl is-active --quiet "$SERVICE_NAME" 2>/dev/null; then
        ok "Servis çalışıyor: $SERVICE_NAME"
        systemctl status "$SERVICE_NAME" --no-pager -l
    else
        warn "Servis çalışmıyor: $SERVICE_NAME"
    fi
    exit 0
fi

# ---- Removal ----
if $REMOVE; then
    info "Kanije Kalesi kaldırılıyor..."
    systemctl stop "$SERVICE_NAME" 2>/dev/null || true
    systemctl disable "$SERVICE_NAME" 2>/dev/null || true
    rm -f "$SERVICE_FILE"
    systemctl daemon-reload
    rm -f "$INSTALL_DIR/$BINARY_NAME"
    ok "Kaldırıldı. Config ve veriler korundu: $CONFIG_DIR, $DATA_DIR"
    exit 0
fi

# ---- Require root ----
[[ $EUID -ne 0 ]] && error "Bu script root hakları gerektiriyor. sudo ile çalıştırın."

# ---- Detect binary ----
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH=""

# Look for the binary relative to the script or in common paths
for candidate in \
    "$(dirname "$SCRIPT_DIR")/kanije" \
    "$SCRIPT_DIR/../../kanije" \
    "./kanije" \
    "./dist/kanije-linux-amd64" \
    "./dist/kanije-linux-arm64" \
    "./dist/kanije-linux-arm"; do
    if [[ -f "$candidate" ]]; then
        BINARY_PATH="$(realpath "$candidate")"
        break
    fi
done

[[ -z "$BINARY_PATH" ]] && error "kanije binary bulunamadı. Önce 'make build-linux' çalıştırın."

info "Binary bulundu: $BINARY_PATH"

# ---- Architecture check ----
BINARY_ARCH=$(file "$BINARY_PATH" | grep -oP '(x86-64|ARM|aarch64)' | head -1)
SYSTEM_ARCH=$(uname -m)
info "Sistem mimarisi: $SYSTEM_ARCH | Binary: $BINARY_ARCH"

# ---- Create user ----
if ! id "$KANIJE_USER" &>/dev/null; then
    useradd -r -s /sbin/nologin -d "$DATA_DIR" -c "Kanije Kalesi daemon" "$KANIJE_USER"
    ok "Kullanıcı oluşturuldu: $KANIJE_USER"
else
    info "Kullanıcı zaten var: $KANIJE_USER"
fi

# ---- Create directories ----
mkdir -p "$CONFIG_DIR" "$DATA_DIR"
chown "$KANIJE_USER:$KANIJE_USER" "$DATA_DIR"
chmod 750 "$DATA_DIR" "$CONFIG_DIR"
ok "Dizinler hazırlandı"

# ---- Install binary ----
install -m 755 -o root -g root "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
ok "Binary kuruldu: $INSTALL_DIR/$BINARY_NAME"

# ---- Install config (if not present) ----
if [[ ! -f "$CONFIG_DIR/config.toml" ]]; then
    EXAMPLE="$(dirname "$SCRIPT_DIR")/config.example.toml"
    if [[ -f "$EXAMPLE" ]]; then
        install -m 640 -o root -g "$KANIJE_USER" "$EXAMPLE" "$CONFIG_DIR/config.toml"
        ok "Örnek config kuruldu: $CONFIG_DIR/config.toml"
    else
        warn "Örnek config bulunamadı. config.toml'u kendiniz oluşturun."
    fi
fi

# ---- Install systemd service ----
SERVICE_SRC="$(dirname "${BASH_SOURCE[0]}")/kanije.service"
if [[ -f "$SERVICE_SRC" ]]; then
    install -m 644 "$SERVICE_SRC" "$SERVICE_FILE"
else
    error "kanije.service dosyası bulunamadı: $SERVICE_SRC"
fi

systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
ok "Systemd servisi etkinleştirildi"

# ---- Secrets file hint ----
SECRETS_FILE="$CONFIG_DIR/secrets.env"
if [[ ! -f "$SECRETS_FILE" ]]; then
    cat > "$SECRETS_FILE" << 'EOF'
# Telegram kimlik bilgileri — bu dosya sistem loglarında görünmez
# kanije.service EnvironmentFile direktifi ile okunur
KANIJE_BOT_TOKEN=buraya_bot_token_yaz
KANIJE_CHAT_ID=buraya_chat_id_yaz
EOF
    chmod 640 "$SECRETS_FILE"
    chown "root:$KANIJE_USER" "$SECRETS_FILE"
    warn "Telegram kimlik bilgilerini düzenleyin: $SECRETS_FILE"
fi

# ---- Enable EnvironmentFile in service ----
# Uncomment the EnvironmentFile line in the service
sed -i 's|^# EnvironmentFile=|EnvironmentFile=|' "$SERVICE_FILE"
systemctl daemon-reload

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
ok "Kanije Kalesi Sentinel kurulumu tamamlandı!"
echo ""
echo "Sıradaki adımlar:"
echo ""
echo "  1. Telegram kimlik bilgilerini düzenleyin:"
echo "     sudo nano $SECRETS_FILE"
echo ""
echo "  2. Servisi başlatın:"
echo "     sudo systemctl start $SERVICE_NAME"
echo ""
echo "  3. Durumu kontrol edin:"
echo "     sudo systemctl status $SERVICE_NAME"
echo "     sudo journalctl -u $SERVICE_NAME -f"
echo ""
echo "  4. Telegram botunuza /kurulum yazarak ayarları yapın."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
