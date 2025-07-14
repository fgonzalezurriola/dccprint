package scripts

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fgonzalezurriola/dccprint/internal/config"
)

// Func to remove al Shell Scripts that starts with "printdcc_"
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

// Func to create the main feature in order to print
func CreateScript(filename string) (string, error) {
	basename := strings.TrimSuffix(filename, filepath.Ext(filename))
	psname := basename + ".ps"
	cfg := config.Load()
	username := cfg.Account
	printer := cfg.Printer
	mode := cfg.Mode

	scriptContent := `#!/usr/bin/env bash
ORANGE='\033[38;5;208m'
GREEN='\033[0;32m'
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
echo 'Este script generado:'	
echo 'Con SSH se conecta a Anakena con tu cuenta, pidiendo tu contraseña'
echo 'Con SCP el archivo en formato .ps (postscript)'
echo 'Usa el comando para imprimir según la configuración seleccionada'
echo 'A continuación ingresa tu contraseña para ingresar por SSH'
echo '==============================================================='

`
	scriptContent += fmt.Sprintf("echo '1. Subiendo archivo %s...\n", filename)
	scriptContent += fmt.Sprintf("scp %q %s@anakena.dcc.uchile.cl:~\n", filename, username)
	scriptContent += "if [ $? -ne 0 ]; then\n  echo -e \"\\033[0;31mERROR: Falló la subida del archivo a anakena. Verifica tu conexión y vuelve a intentar.\\033[0m\"\n  exit 1\nfi\n\n"
	scriptContent += "echo '2. Conectando a anakena y ejecutando comandos...\n"
	scriptContent += fmt.Sprintf("ssh %s@anakena.dcc.uchile.cl << 'EOF'\n", username)
	scriptContent += fmt.Sprintf("pdf2ps %q %q\n", filename, psname)

	var printCommand string

	switch printer {
	case "Toqui":
		switch mode {
		case "Simple (Reverso en blanco)":
			printCommand = fmt.Sprintf("lpr %s", psname)
		case "Doble cara, Borde largo (Recomendado)":
			printCommand = fmt.Sprintf("duplex %s | lpr", psname)
		case "Doble cara, Borde corto":
			printCommand = fmt.Sprintf("duplex -l %s | lpr", psname)
		}
	case "Salita":
		switch mode {
		case "Simple (Reverso en blanco)":
			printCommand = fmt.Sprintf("lpr -P hp-335 %s", psname)
		case "Doble cara, Borde largo (Recomendado)":
			printCommand = fmt.Sprintf("duplex %s | lpr -P hp-335", psname)
		case "Doble cara, Borde corto":
			printCommand = fmt.Sprintf("duplex -l %s | lpr -P hp-335", psname)
		}
	}

	scriptContent += printCommand + "\n"

	// Output verifier
	switch printer {
	case "Salita":
		scriptContent += "lpq -P hp-335\n"
	case "Toqui":
		scriptContent += "lpq\n"
	}

	scriptContent += "papel\n"
	scriptContent += "EOF\n"
	scriptContent += "\nif [ $? -ne 0 ]; then\n  echo -e \"\\033[0;31mERROR: Falló la conexión o ejecución de comandos en anakena. Verifica tu conexión y vuelve a intentar.\\033[0m\"\n  exit 1\nfi\n"
	scriptContent += `echo -e "${GREEN}IMPRESION COMPLETADA!${NC}"`
	scriptContent += "\n"
	scriptContent += `echo -e "Recuerda: usa 'ssh usuario@anakena.dcc.uchile.cl' y el comando 'papel' para ver impresiones restantes."`

	scriptPath := "print-" + basename + ".sh"
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
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
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
