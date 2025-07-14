# DCCPRINT

```sh
 ██████╗   ██████╗  ██████╗   ██████╗  ██████╗  ██╗ ███╗   ██╗ ████████╗
 ██╔══██╗ ██╔════╝ ██╔════╝   ██╔══██╗ ██╔══██╗ ██║ ████╗  ██║ ╚══██╔══╝
 ██║  ██║ ██║      ██║        ██████╔╝ ██████╔╝ ██║ ██╔██╗ ██║    ██║
 ██║  ██║ ██║      ██║        ██╔═══╝  ██╔══██╗ ██║ ██║╚██╗██║    ██║
 ██████╔╝ ╚██████╗ ╚██████╗   ██║      ██║  ██║ ██║ ██║ ╚████║    ██║
 ╚═════╝   ╚═════╝  ╚═════╝   ╚═╝      ╚═╝  ╚═╝ ╚═╝ ╚═╝  ╚═══╝    ╚═╝
```

DCCPRINT es una herramienta Terminal User Interface (TUI) para imprimir archivos `.pdf` en el Departamento de Ciencias de la Computación de la Universidad de Chile. App hecha en Go con la librería de BubbleTea.
Esta app fue hecha para reafirmar los conceptos de Go después de haber leído medio libro de Go y aprender cosas sobre distribución y package managers.

## Uso

Después de haber instalado, abre una terminal y ve al directorio donde se encuentra el archivo que quieres imprimir y usa el comando `dccprint`. Configura tu nombre de usuario, el lugar de impresión, el modo de impresión, y sigue los pasos que se muestran al final para ejecutar el script generado.

El flujo de uso normal sería el siguiente

Primer uso para instalar y configurar

1. Instalar
2. Usar dccprint en tu terminal
3. Configurar usuario
4. Configurar donde imprimir (salita, toqui)
5. Configurar el modo para imprimir (doble cara borde largo/corto, simple)

Luego, se reutilizará la configuración a no ser que la quieras cambiar en el menú de inicio

1. Ir a la opción de Imprimir y seleccionar pdf
2. Leer mensaje de la vista final -> Enter
3. CTRL+SHIFT+V -> Enter
4. Ingresar tu contraseña DCC para la conexión ssh
5. Esperar tu impresión

Para ver el contenido del script generado usa

```sh
cat printdcc-NOMBRE.sh
```

La configuración se guarda en `$HOME/.dccprint_config.json` y la puedes actualizar en el menú principal

## Instalación

- **Arch Linux / Manjaro (AUR)**

  ```sh
  # Usa tu AUR helper favorito
  yay -S dccprint
  paru -S dccprint
  ```

  Para actualizar

  ```sh
  yay -Syu dccprint
  paru -Syu dccprint

  ```

- **macOS (Homebrew)**

  ```sh
  brew tap fgonzalezurriola/homebrew-tap
  brew install dccprint
  ```

  Para actualizar

  ```sh
  brew update
  brew upgrade dccprint
  ```

## Snap

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

Para los package managers que usan archivos `.deb`, `.rpm` o `.apk` se puede usar el script que detecta arquitectura, OS y package manager para instalar dccprint

### Debian / Ubuntu / Fedora / openSUSE / Alpine Linux

Para Instalar o Actualizar usa la opción de curl o wget que ejecuta el script de instalación

```sh
curl -sSL https://raw.githubusercontent.com/fgonzalezurriola/dccprint/main/install.sh | bash
```

```
wget -O- https://raw.githubusercontent.com/fgonzalezurriola/dccprint/main/install.sh | bash
```

## Desinstalación

```sh
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
