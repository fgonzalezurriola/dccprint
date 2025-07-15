package scripts

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fgonzalezurriola/dccprint/internal/config"
)

// Func to remove al Shell Scripts that starts with "printdcc_"
// RemoveGeneratedScripts deletes all generated shell scripts in the given directory that start with "dccprint_" and end with ".sh".
func RemoveGeneratedScripts(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "dccprint_") && strings.HasSuffix(entry.Name(), ".sh") {
			fullPath := filepath.Join(dir, entry.Name())
			err := os.Remove(fullPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Func to retrieve all pdfs in the current dir
func GetPDFFiles() []string {
	var PDFs []string

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pdf") {
			PDFs = append(PDFs, entry.Name())
		}
	}

	return PDFs
}

// ValidatePSWithGhostscript runs Ghostscript on the .ps file with a 2-second timeout
// Returns an error, nil if it's valid
// If the timeout it's reached or the ghoscript command passes, return nil
// If ghostcript returns an error if the file is invalid
func ValidatePSWithGhostscript(psPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "gs", "-sDEVICE=nullpage", "-dBATCH", "-dNOPAUSE", psPath)
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil
		}
		return fmt.Errorf("Ghostscript could not validate the .ps file. Suggestion: flatten the PDF or use repair tools. Error: %w", err)
	}
	return nil
}

// Func to create the main feature in order to print
func CreateScript(filename string) (string, error) {
	escapedName := EscapeFilename(filename)
	basename := strings.TrimSuffix(escapedName, filepath.Ext(escapedName))
	cfg := config.Load()
	username := cfg.Account
	printer := cfg.Printer
	mode := cfg.Mode

	// Create a temporal copy pdf to use in SCP and SSH if the escaped name is different than the filename
	useTempCopy := escapedName != filename
	if useTempCopy {
		input, err := os.ReadFile(filename)
		if err != nil {
			return "", fmt.Errorf("Error leyendo el archivo seleccionado: %w", err)
		}
		err = os.WriteFile(escapedName, input, 0644)
		if err != nil {
			return "", fmt.Errorf("Error creando copia temporal: %w", err)
		}
	}

	// If the script uses the escaped copy, it's deleted at the end
	// The script rm itself after use
	scriptContent := `#!/usr/bin/env bash
ORANGE='\033[38;5;208m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${ORANGE}"
echo '
  ██████╗   ██████╗  ██████╗      ██████╗  ██████╗  ██╗ ███╗   ██╗ ████████╗
  ██╔══██╗ ██╔════╝ ██╔════╝      ██╔══██╗ ██╔══██╗ ██║ ████╗  ██║ ╚══██╔══╝
  ██║  ██║ ██║      ██║           ██████╔╝ ██████╔╝ ██║ ██╔██╗ ██║    ██║
  ██║  ██║ ██║      ██║           ██╔═══╝  ██╔══██╗ ██║ ██║╚██╗██║    ██║
  ██████╔╝ ╚██████╗ ╚██████╗      ██║      ██║  ██║ ██║ ██║ ╚████║    ██║
  ╚═════╝   ╚═════╝  ╚═════╝      ╚═╝      ╚═╝  ╚═╝ ╚═╝ ╚═╝  ╚═══╝    ╚═╝
'
echo -e "${NC}"

echo '==============================================================='
echo 'DCC PRINT - SCRIPT GENERADO'
echo 'Este script:'	
echo '1. Con SCP sube el archivo PDF a Anakena'
echo '2. Con SSH se conecta y ejecuta comandos de impresión'
echo '3. Usa el comando para imprimir según la configuración seleccionada'
echo 'A continuación ingresa tu contraseña para subir el archivo'
echo '==============================================================='

`
	// SCP step
	scriptContent += fmt.Sprintf("echo -e \"${GREEN}1. Subiendo archivo %s...${NC}\"\n", escapedName)
	scriptContent += fmt.Sprintf("scp %q %s@anakena.dcc.uchile.cl:~\n", escapedName, username)
	scriptContent += "if [ $? -ne 0 ]; then\n"
	scriptContent += "  echo -e \"${RED}ERROR: Falló la subida del archivo a anakena. Verifica tu conexión y vuelve a intentar.${NC}\"\n"
	scriptContent += "  exit 1\nfi\n\n"

	// SSH step
	scriptContent += "echo -e \"${GREEN}2. Conectando a anakena y ejecutando comandos...${NC}\"\n"
	scriptContent += fmt.Sprintf("ssh %s@anakena.dcc.uchile.cl << 'EOF'\n", username)

	var printCommand string
	pdfname := filepath.Base(escapedName)
	psname := strings.TrimSuffix(pdfname, filepath.Ext(pdfname)) + ".ps"

	// Validate .ps generated
	if err := ValidatePSWithGhostscript(psname); err != nil {
		return "", err
	}

	switch printer {
	case "Toqui":
		switch mode {
		case "Simple (Reverso en blanco)":
			printCommand = fmt.Sprintf("pdf2ps %s %s && lpr %s", pdfname, psname, psname)
		case "Doble cara, Borde largo (Recomendado)":
			printCommand = fmt.Sprintf("pdf2ps %s %s && duplex %s|lpr", pdfname, psname, psname)
		case "Doble cara, Borde corto":
			printCommand = fmt.Sprintf("pdf2ps %s %s && duplex -l %s|lpr", pdfname, psname, psname)
		}
	case "Salita":
		switch mode {
		case "Simple (Reverso en blanco)":
			printCommand = fmt.Sprintf("pdf2ps %s %s && lpr -P hp-335 %s", pdfname, psname, psname)
		case "Doble cara, Borde largo (Recomendado)":
			printCommand = fmt.Sprintf("pdf2ps %s %s && duplex %s|lpr -P hp-335", pdfname, psname, psname)
		case "Doble cara, Borde corto":
			printCommand = fmt.Sprintf("pdf2ps %s %s && duplex -l %s|lpr -P hp-335", pdfname, psname, psname)
		}
	}

	scriptContent += fmt.Sprintf("if [ ! -f \"%s\" ]; then\n", pdfname)
	scriptContent += "  echo \"ERROR: El archivo PDF no se encontró en el directorio home\"\n"
	scriptContent += "  exit 1\nfi\n\n"
	scriptContent += printCommand + "\n\n"

	switch printer {
	case "Salita":
		scriptContent += "lpq -P hp-335\n"
	case "Toqui":
		scriptContent += "lpq\n"
	}

	scriptContent += "papel\n"
	scriptContent += "EOF\n\n"

	scriptContent += "if [ $? -ne 0 ]; then\n"
	scriptContent += "  echo -e \"${RED}ERROR: Falló la conexión o ejecución de comandos en anakena.${NC}\"\n"
	scriptContent += "  exit 1\nfi\n\n"

	scriptContent += "echo -e \"${GREEN}¡IMPRESIÓN COMPLETADA!${NC}\"\n"
	scriptContent += "echo -e \"Recuerda: usa 'ssh usuario@anakena.dcc.uchile.cl' y el comando 'papel' para ver impresiones restantes.\"\n"

	scriptPath := "print-" + basename + ".sh"
	scriptContent += `rm -- "$0"`

	if useTempCopy {
		scriptContent += fmt.Sprintf("\nrm -- %q\n", escapedName)
	}

	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		return "", fmt.Errorf("error escribiendo archivo: %w", err)
	}

	return scriptPath, nil
}
func CopyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "clip")
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux", "freebsd", "openbsd", "netbsd":
		// Try X11, Wayland and an ancient clipboard
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd = exec.Command("wl-copy")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("ninguna herramienta de portapapeles encontrada (se necesita xclip, wl-copy o xsel)")
		}
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func PrintSuccessMessage(scriptName string) tea.Cmd {
	return func() tea.Msg {
		fmt.Print("\n\n")
		fmt.Println("Script generado exitosamente!")
		fmt.Printf("Comando copiado: ./%s\n", scriptName)
		fmt.Print("\n")
		fmt.Println("Instrucciones:")
		fmt.Println("1. Ctrl+Shift+V para pegar")
		fmt.Println("2. Enter para ejecutar")
		fmt.Println("3. Ingresa contraseña SSH")
		fmt.Println("4. Archivo se enviará a impresora")
		fmt.Print("\n")
		fmt.Println("Listo para imprimir!")
		fmt.Print("\n")
		return nil
	}
}
