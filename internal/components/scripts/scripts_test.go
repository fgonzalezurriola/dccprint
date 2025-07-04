package scripts

import (
	"os"
	"testing"
)

func TestRemoveGeneratedScripts(t *testing.T) {
	testScript := "print_test.sh"
	content := "#!/bin/bash\n"
	content += `echo "hola"\n`

}
