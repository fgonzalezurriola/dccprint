package scripts

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveGeneratedScripts(t *testing.T) {
	tempDir := t.TempDir()
	testScriptPath := filepath.Join(tempDir, "dccprint_test.sh")
	content := []byte("#!/bin/bash\necho 'hola'")
	err := os.WriteFile(testScriptPath, content, 0755)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	err = removeGeneratedScripts(tempDir)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

}
