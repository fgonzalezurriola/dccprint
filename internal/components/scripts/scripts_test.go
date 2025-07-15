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
	content := []byte("#!/usr/bin/env bash\necho 'hola'")
	err := os.WriteFile(testScriptPath, content, 0755)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	err = RemoveGeneratedScripts(tempDir)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
}

func TestEscapeFilename(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"archivo con espacios.pdf", "archivoconespacios.pdf"},
		{"áéíóúñÑ.pdf", "aeiounn.pdf"},
		{"file@#$.pdf", "file.pdf"},
		{"normal-file.pdf", "normal-file.pdf"},
		{"", "file.pdf"},
		{"á rch ivo.pdf", "archivo.pdf"},
	}
	for _, c := range cases {
		out := EscapeFilename(c.input)
		if out != c.expected {
			t.Errorf("EscapeFilename(%q) = %q; want %q", c.input, out, c.expected)
		}
	}
}

func TestCopyToClipboard(t *testing.T) {
	text := "Prueba de Clipboard"
	err := CopyToClipboard(text)
	if err != nil {
		t.Errorf("CopyToClipboard devolvió error: %v", err)
	}
}
