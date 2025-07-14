#!/usr/bin/env bash

set -e

echo "Detectando arquitectura y sistema operativo..."

# Detectar arquitectura
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH=amd64 ;;
    aarch64) ARCH=arm64 ;;
    armv7l) ARCH=arm ;;
    i386|i686) ARCH=386 ;;
    *) echo "Arquitectura no soportada: $ARCH"; exit 1 ;;
esac

# Detectar sistema operativo y gestor de paquetes
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
PKG=""
INSTALL=""

if [ "$OS" = "linux" ]; then
    if command -v apt-get >/dev/null; then
        PKG="deb"
        INSTALL="sudo dpkg -i"
    elif command -v dnf >/dev/null || command -v zypper >/dev/null; then
        PKG="rpm"
        if command -v zypper >/dev/null; then
            INSTALL="sudo zypper install"
        else
            INSTALL="sudo dnf install -y"
        fi
    elif command -v apk >/dev/null; then
        PKG="apk"
        INSTALL="sudo apk add --allow-untrusted"
    else
        echo "No se detectó un gestor de paquetes soportado (dpkg, dnf, zypper, apk, apt)"; exit 1
    fi
elif [ "$OS" = "darwin" ]; then
    if command -v brew >/dev/null; then
        echo "Usa las instrucciones para MacOS con Homebrew"; exit 0
    else
        echo "Instalación manual en macOS no soportada. Usa Homebrew."; exit 1
    fi
else
    echo "Sistema operativo no soportado: $OS"; exit 1
fi

echo "Obteniendo la última versión de dccprint..."
VERSION=$(curl -s https://api.github.com/repos/fgonzalezurriola/dccprint/releases/latest | grep tag_name | cut -d '"' -f 4)
if [ -z "$VERSION" ]; then
    echo "No se pudo obtener la última versión. Revisa tu conexión a internet."; exit 1
fi

FILE="dccprint_${VERSION#v}_linux_${ARCH}.${PKG}"
URL="https://github.com/fgonzalezurriola/dccprint/releases/download/${VERSION}/${FILE}"

echo "Descargando $FILE desde $URL ..."
if command -v wget >/dev/null; then
    wget "$URL" -O "$FILE"
elif command -v curl >/dev/null; then
    curl -L "$URL" -o "$FILE"
else
    echo "wget o curl no están instalados"; exit 1
fi

echo "Instalando dccprint (puede requerir tu contraseña de sudo)..."
$INSTALL "$FILE"

echo "Limpiando archivos temporales..."
rm "$FILE"

echo "dccprint instalado o actualizado correctamente."
