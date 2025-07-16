# DCCPRINT

DCCPRINT es una Terminal User Interface (TUI) escrita en Go, junto a la librería de BubbleTea, para que estudiantes puedan imprimir archivos `.pdf` dentro del Departamento de Ciencias de la Computación de la Universidad de Chile.

## Requisitos

Tener una cuenta DCC, esta cuenta se puede pedir al obtener el código DCC. Las instrucciones para pedirla se encuentran en [la página de sistemas](https://sistemas.dcc.uchile.cl/)

Tener una máquina con sistema operativo macOS o Linux.

> [!WARNING]
> Esto no ha sido testeado con WSL. En teoría debería funcionar, de ser así hazmelo saber en mi telegram @fgonzalezurriola o en una issue de github.

## Uso

El flujo de uso esperado es el siguiente

1. Abrir una terminal, ir al directorio donde está el pdf, usar `dccprint`
2. Escribir usuario DCC
3. Seleccionar Salita o Toqui
4. Seleccionar entre los 3 modos de impresión

> [!TIP]
> Borde largo es para anillarlo tipo libro
>
> Borde corto es para anillarlo tipo croquera

Luego saldrá el menú de inicio, que usará la configuración guardada para imprimir sigue estos pasos

1. Ir al menú de **Imprimir PDF**,
2. Seleccionar PDF con **Enter**
3. Leer el mensaje y salir con **Enter**
4. **CTRL+SHIFT+V** + **Enter** (Esto ejecutará el script generado con ./printdcc-script.sh)
5. Ingresar tu contraseña de usuario DCC

Con esto se envía el comando para imprimir con la configuración guardada en `$HOME/.dccprint_config.json`. Esta configuración la puedes actualizar en el menú principal

> [!TIP]
> El archivo `.sh` generado se autoelimina en el uso.
> Puedes usar `cat` para ver su contenido antes de ejecutarlo

## Instalación

### Arch Linux / Manjaro (AUR)

Para instalar

```sh
# Usa tu AUR helper favorito, puedes usar paru con la mismas flags
yay -S dccprint
```

Para actualizar

```sh
yay -Syu dccprint
```

---

### macOS (Homebrew)

```sh
brew tap fgonzalezurriola/homebrew-tap
brew install dccprint
```

Para actualizar

```sh
brew update
brew upgrade dccprint
```

---

### Snap

Snap es una herramienta que viene preinstalada en Ubuntu y funciona en muchas otras distros

![snapper](/assets/snapper.gif)

Gif del legendario Terry Davis

Para instalar

```sh
snap install dccprint
```

Para actualizar

```sh
sudo snap refresh
```

## Otros

También está contemplada la distribución en formato `.deb`, `.rpm` o `.apk` usando el script de `install.sh` puedes instalar dccprint

### Debian / Ubuntu / Fedora / openSUSE / Alpine Linux

Para **Instalar** o **Actualizar** usa curl o wget para descargar y ejecutar el script de instalación

```sh
curl -sSL https://raw.githubusercontent.com/fgonzalezurriola/dccprint/main/install.sh | bash
```

```sh
wget -O- https://raw.githubusercontent.com/fgonzalezurriola/dccprint/main/install.sh | bash
```

## Desinstalación

```sh
# Eliminar el archivo de configuración primero
# Luego, desinstala según como la instalaste
rm $HOME/.dccprint_config.json

# Arch Linux / Manjaro (AUR)
yay -R dccprint
paru -R dccprint

# Homebrew (macOS/Linux)
brew uninstall dccprint

# Snap
snap remove dccprint

# Debian/Ubuntu
sudo apt-get remove dccprint

# Fedora/openSUSE
sudo dnf remove dccprint

# Alpine Linux
sudo apk del dccprint

# Si instalaste el binario manualmente
sudo rm /usr/local/bin/dccprint
```
