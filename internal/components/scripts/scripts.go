package scripts

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func CreateScript(filename string) {
	basename := strings.TrimSuffix(filename, filepath.Ext(filename))
	psname := basename + ".ps"

	// var printer string
	// if m.printConfig.Location == 1 {
	// 	printer = "hp-335"
	// } else {
	// 	printer = "hp"
	// }

	// var lprCommand string
	// if m.printConfig.Copies > 1 {
	// 	lprCommand = fmt.Sprintf("duplex %q | lpr -P %s -#%d", psname, printer, m.printConfig.Copies)
	// } else {
	// 	lprCommand = fmt.Sprintf("duplex %q | lpr -P %s", psname, printer)
	// }

	scriptContent := `#!usr/bin/env bash
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
echo 'Solo falta ingresar dos veces tu contraseña cuenta anakena para imprimir'
echo '==============================================================='

`
	scriptContent += fmt.Sprintf("echo '1. Subiendo archivo %s...\n", filename)
	// scriptContent += fmt.Sprintf("scp %q %s@anakena.dcc.uchile.cl:~\n\n", filename, m.userConfig.Username)

	scriptContent += "echo '2. Conectando a anakena y ejecutando comandos...\n"
	// scriptContent += fmt.Sprintf("ssh %s@anakena.dcc.uchile.cl << 'EOF'\n", m.userConfig.Username)
	scriptContent += fmt.Sprintf("pdf2ps %q %q\n", filename, psname)
	// scriptContent += lprCommand + "\n"
	// scriptContent += fmt.Sprintf("lpq -P %s\n", printer)
	scriptContent += "papel\n"
	scriptContent += "EOF\n\n"
	scriptContent += `echo -e "${GREEN} IMPRESION COMPLETADA!${NC}"`
	scriptContent += `echo -e "Recuerda que puedes usar ssh usuario@dcc.anakena.uchile.cl y el comando papel para ver cuantas impresiones te quedan. El papel se actualiza después de que termina la impresión en curso."`
	scriptPath := "print-" + basename + ".sh"
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		return
	}

}
