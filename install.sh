#!/bin/bash

# Script for installing dccprint
# Support for
# Archs:
#   amd64
#   arm64
#   i386


set -e

# Arch
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH=amd64 ;;
    aarch64) ARCH=arm64 ;;
    armv7l) ARCH=arm ;;
    i386|i686) ARCH=386 ;;
    *) echo "Arquitectura no soportada: $ARCH"; exit 1 ;;
esac

# Operative system
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
PKG=""
INSTALL=""

if [ "$OS" = "linux" ]; then
    if command -v apt-get >/dev/null; then
        PKG="deb"
        INSTALL="sudo dpkg -i"
    elif command -v dnf >/dev/null; then
        PKG="rpm"
        INSTALL="sudo dnf install -y"
    elif command -v apk >/dev/null; then
        PKG="apk"
        INSTALL="sudo apk add --allow-untrusted"
    else
        echo "No se detectó un gestor de paquetes soportado (dpkg, dnf, apk, apt)"; exit 1
    fi
elif [ "$OS" = "darwin" ]; then
    if command -v brew >/dev/null; then
        echo "Usa las instrucciones para MacOS con Homebew"; exit 0
    else
        echo "Instalación manual en macOS no soportada. Usa Homebrew."; exit 1
    fi
else
    echo "Sistema operativo no soportado: $OS"; exit 1
fi

REPO="https://github.com/fgonzalezurriola/dccprint/releases/latest/download"
FILE="dccprint_*_${OS}_${ARCH}.${PKG}"
URL="$REPO/$FILE"

# Download package
if command -v wget >/dev/null; then
    wget "$URL" -O "$FILE"
elif command -v curl >/dev/null; then
    curl -L "$URL" -o "$FILE"
else
    echo "wget o curl no están instalados"; exit 1
fi

$INSTALL "$FILE"
rm "$FILE"

echo "dccprint instalado correctamente."
