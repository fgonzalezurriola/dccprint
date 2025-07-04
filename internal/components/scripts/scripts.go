package scripts

import (
	"log"
	"os"
	"strings"
)

// Func to remove al Shell Scripts that starts with "printdcc-"
func removeGeneratedScripts() error {

}

// Func to retrieve all pdfs in the current dir
func getPDFFiles() []string {
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
