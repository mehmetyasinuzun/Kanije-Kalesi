#!/bin/bash
# Kanije Kalesi — Tek satır web kurulum (Linux / Raspberry Pi)
#
#   curl -fsSL https://raw.githubusercontent.com/mehmetyasinuzun/Kanije-Kalesi/master/go/deploy/linux/get.sh | sudo bash -s -- --token "TOKEN" --chat "CHATID"
#
# install.sh + service dosyasını indirir, kurulumu çalıştırır (binary'yi de o indirir).
set -euo pipefail

REPO="mehmetyasinuzun/Kanije-Kalesi"
RAW="https://raw.githubusercontent.com/$REPO/master/go/deploy/linux"
TMP="$(mktemp -d)"

echo "🏰 Kanije Kalesi kurulum betikleri indiriliyor..."
curl -fsSL "$RAW/install.sh"     -o "$TMP/install.sh"
curl -fsSL "$RAW/kanije.service" -o "$TMP/kanije.service"

bash "$TMP/install.sh" "$@"
