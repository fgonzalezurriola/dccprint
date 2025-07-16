package scripts

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
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
		if !entry.IsDir() && ((strings.HasPrefix(entry.Name(), "dccprint_") && strings.HasSuffix(entry.Name(), ".sh")) || strings.HasPrefix(entry.Name(), "dccprint-temp-")) {
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

// ValidatePDFWithGhostscript validates a PDF file using Ghostscript with a 3-second timeout.
// It uses fast flags (-o /dev/null -sDEVICE=nullpage) to avoid disk IO and speed up validation.
// If Ghostscript does not finish in time, the process group is killed to avoid zombies (works on Mac and Linux).
// Returns nil if the file is valid or if timeout is reached. Returns an error if Ghostscript detects a fatal error in the file.
func ValidatePDFWithGhostscript(pdfPath string) error {
	if _, err := exec.LookPath("gs"); err != nil {
		return fmt.Errorf("Ghostscript (gs) is not installed: %w", err)
	}

	cmd := exec.Command("gs", "-o", "/dev/null", "-sDEVICE=nullpage", pdfPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var output strings.Builder
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Could not start Ghostscript: %w", err)
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case err := <-done:
		outStr := output.String()
		if err != nil {
			if strings.Contains(outStr, "Error") || strings.Contains(outStr, "FATAL") || strings.Contains(outStr, "Unrecoverable error") {
				return fmt.Errorf("Ghostscript detected fatal error in PDF file. Output: %s", outStr)
			}
			return nil
		}
		return nil
	case <-time.After(3 * time.Second):
		killProcessGroup(cmd)
		return nil
	}
}

// killProcessGroup forcefully kills the process group for the given command.
// This ensures that all child processes are terminated, preventing zombies.
func killProcessGroup(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	pgid := cmd.Process.Pid
	if cmd.SysProcAttr != nil && cmd.SysProcAttr.Setpgid {
		pgid = cmd.Process.Pid
	}
	_ = syscall.Kill(-pgid, syscall.SIGTERM)
	time.Sleep(100 * time.Millisecond)
	_ = syscall.Kill(-pgid, syscall.SIGKILL)
	_ = cmd.Process.Kill()
	time.Sleep(200 * time.Millisecond)
}

// Func to create the main feature in order to print
func CreateScript(filename string) (string, error) {
	originalEscapedName := EscapeFilename(filename)
	basename := strings.TrimSuffix(originalEscapedName, filepath.Ext(originalEscapedName))

	// Prefijo para archivos temporales
	tempEscapedName := "dccprint-temp-" + originalEscapedName
	cfg := config.Load()
	username := cfg.Account
	printer := cfg.Printer
	mode := cfg.Mode

	// Create a temporal copy pdf to use in SCP and SSH if the escaped name is different than the filename
	useTempCopy := tempEscapedName != filename
	if useTempCopy {
		input, err := os.ReadFile(filename)
		if err != nil {
			return "", fmt.Errorf("Error leyendo el archivo seleccionado: %w", err)
		}
		err = os.WriteFile(tempEscapedName, input, 0644)
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
	scriptContent += fmt.Sprintf("echo -e \"${GREEN}1. Subiendo archivo %s...${NC}\"\n", tempEscapedName)
	scriptContent += fmt.Sprintf("scp %q %s@anakena.dcc.uchile.cl:~\n", tempEscapedName, username)
	scriptContent += "if [ $? -ne 0 ]; then\n"
	scriptContent += "  echo -e \"${RED}ERROR: Falló la subida del archivo a anakena. Verifica tu conexión y vuelve a intentar.${NC}\"\n"
	scriptContent += "  exit 1\nfi\n\n"

	// SSH step
	scriptContent += "echo -e \"${GREEN}2. Conectando a anakena y ejecutando comandos...${NC}\"\n"
	scriptContent += fmt.Sprintf("ssh %s@anakena.dcc.uchile.cl << 'EOF'\n", username)

	var printCommand string
	pdfname := filepath.Base(tempEscapedName)
	psname := strings.TrimSuffix(pdfname, filepath.Ext(pdfname)) + ".ps"

	// Validate the PDF directly with Ghostscript for speed and reliability
	if err := ValidatePDFWithGhostscript(filename); err != nil {
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
		scriptContent += fmt.Sprintf("\nrm -- %q\n", tempEscapedName)
	}

	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		return "", fmt.Errorf("error escribiendo archivo: %w", err)
	}

	return scriptPath, nil
}

func CopyToClipboard(text string) error {
	importedErr := clipboard.WriteAll(text)
	if importedErr != nil {
		return fmt.Errorf("error copiando al portapapeles: %w", importedErr)
	}
	return nil
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
