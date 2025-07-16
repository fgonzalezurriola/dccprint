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

	cfg := config.Load()
	username := cfg.Account
	printer := cfg.Printer
	mode := cfg.Mode

	if err := ValidatePDFWithGhostscript(filename); err != nil {
		return "", err
	}

	scriptContent := `#!/usr/bin/env bash
ORANGE='\033[38;5;208m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${ORANGE}"
echo '
  ██████╗   ██████╗  ██████╗ ██████╗  ██████╗  ██╗ ███╗   ██╗ ████████╗
  ██╔══██╗ ██╔════╝ ██╔════╝ ██╔══██╗ ██╔══██╗ ██║ ████╗  ██║ ╚══██╔══╝
  ██║  ██║ ██║      ██║      ██████╔╝ ██████╔╝ ██║ ██╔██╗ ██║    ██║
  ██║  ██║ ██║      ██║      ██╔═══╝  ██╔══██╗ ██║ ██║╚██╗██║    ██║
  ██████╔╝ ╚██████╗ ╚██████╗ ██║      ██║  ██║ ██║ ██║ ╚████║    ██║
  ╚═════╝   ╚═════╝  ╚═════╝ ╚═╝      ╚═╝  ╚═╝ ╚═╝ ╚═╝  ╚═══╝    ╚═╝
'
echo -e "${NC}"

echo '==============================================================='
echo 'DCC PRINT - SCRIPT GENERADO'
echo 'Este script:'    
echo '1. Se conecta a Anakena con SSH'
echo '2. Transfiere el archivo PDF con cat y ejecuta el comando de impresión'
echo '==============================================================='

`
	pdfname := "dccprint-" + basename + ".pdf"
	psname := "dccprint-" + basename + ".ps"

	var printCommand string
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

	var queueCommand string
	switch printer {
	case "Salita":
		queueCommand = "lpq -P hp-335"
	case "Toqui":
		queueCommand = "lpq"
	}

	// SSH + cat sandwich to avoid asking two times the password
	scriptContent += "echo -e \"${GREEN}Conectando a anakena y procesando archivo...${NC}\"\n"
	// scriptContent += fmt.Sprintf("cat %q | ssh %s@anakena.dcc.uchile.cl 'cat > %s && %s && %s'\n",
	// 	filename, username, pdfname, printCommand, queueCommand)

	// Todo: test this to avoid trash in anakena
	scriptContent += fmt.Sprintf("cat %q | ssh %s@anakena.dcc.uchile.cl 'cat > %s && %s && %s && rm %s %s'\n",
	    filename, username, pdfname, printCommand, queueCommand, pdfname, psname)

	scriptContent += "if [ $? -ne 0 ]; then\n"
	scriptContent += "  echo -e \"${RED}ERROR: Falló la conexión o ejecución de comandos en anakena.${NC}\"\n"
	scriptContent += "  exit 1\nfi\n\n"

	scriptContent += "echo -e \"${GREEN}¡IMPRESIÓN COMPLETADA!${NC}\"\n"
	scriptContent += fmt.Sprintf("echo -e \"Recuerda: usa 'ssh %s@anakena.dcc.uchile.cl' y el comando 'papel' para ver impresiones restantes.\"\n", username)
	scriptContent += "echo -e \"Nota: El comando papel se actualiza después de haber finalizado la impresión\"\n"

	scriptPath := "dccprint-" + basename + ".sh"
	// selfdestruction of script after use
	scriptContent += `rm -- "$0"`

	writeErr := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if writeErr != nil {
		return "", fmt.Errorf("error escribiendo archivo: %w", writeErr)
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
