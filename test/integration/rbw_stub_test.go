package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRBWStub(t *testing.T, status string, get map[string]string, missing ...string) []string {
	t.Helper()
	stubDir := t.TempDir()
	var script strings.Builder
	script.WriteString("#!/bin/sh\n")
	script.WriteString("case \"$1\" in\n")
	script.WriteString(fmt.Sprintf("status) printf '%%s\\n' '%s' ;;\n", status))
	script.WriteString("get)\ncase \"$2\" in\n")
	for key, value := range get {
		script.WriteString(fmt.Sprintf("%s) printf '%%s\\n' '%s' ;;\n", key, value))
	}
	for _, key := range missing {
		script.WriteString(fmt.Sprintf("%s) echo 'not found' >&2; exit 1 ;;\n", key))
	}
	script.WriteString("*) echo 'not found' >&2; exit 1 ;;\nesac\n;;\n")
	script.WriteString("list) printf 'DATABASE_URL\\nJWT_SECRET\\n' ;;\n")
	script.WriteString("*) echo 'unsupported' >&2; exit 1 ;;\nesac\n")
	path := filepath.Join(stubDir, "rbw")
	writeFile(t, path, script.String())
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatal(err)
	}
	return []string{"PATH=" + stubDir + string(os.PathListSeparator) + os.Getenv("PATH")}
}
