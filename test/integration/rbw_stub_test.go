package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func rbwLookupKey(item, field string) string {
	return item + "::" + field
}

func writeRBWStub(t *testing.T, status string, get map[string]string, missing ...string) []string {
	t.Helper()
	stubDir := t.TempDir()
	var script strings.Builder
	script.WriteString("#!/bin/sh\n")
	script.WriteString("case \"$1\" in\n")
	switch status {
	case "unlocked":
		script.WriteString("unlocked) exit 0 ;;\n")
		script.WriteString("list) printf 'DATABASE_URL\\nJWT_SECRET\\n' ;;\n")
	case "locked":
		script.WriteString("unlocked) exit 1 ;;\n")
		script.WriteString("list) echo 'database is locked' >&2; exit 1 ;;\n")
	case "logged out":
		script.WriteString("unlocked) exit 1 ;;\n")
		script.WriteString("list) echo 'not logged in' >&2; exit 1 ;;\n")
	default:
		script.WriteString("unlocked) exit 1 ;;\n")
		script.WriteString(fmt.Sprintf("list) printf '%%s\\n' '%s' ;;\n", status))
	}
	script.WriteString(fmt.Sprintf("status) printf '%%s\\n' '%s' ;;\n", status))
	script.WriteString("get)\n")
	script.WriteString("field=''\n")
	script.WriteString("item=''\n")
	script.WriteString("shift\n")
	script.WriteString("while [ \"$#\" -gt 0 ]; do\n")
	script.WriteString("  case \"$1\" in\n")
	script.WriteString("    --field) field=\"$2\"; shift 2 ;;\n")
	script.WriteString("    --*) shift ;;\n")
	script.WriteString("    *) item=\"$1\"; shift ;;\n")
	script.WriteString("  esac\n")
	script.WriteString("done\n")
	script.WriteString("case \"$item::$field\" in\n")
	for key, value := range get {
		script.WriteString(fmt.Sprintf("%s) printf '%%s\\n' '%s' ;;\n", key, value))
	}
	for _, key := range missing {
		script.WriteString(fmt.Sprintf("%s) echo 'not found' >&2; exit 1 ;;\n", key))
	}
	script.WriteString("*) echo 'not found' >&2; exit 1 ;;\nesac\n;;\n")
	script.WriteString("*) echo 'unsupported' >&2; exit 1 ;;\nesac\n")
	path := filepath.Join(stubDir, "rbw")
	writeFile(t, path, script.String())
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatal(err)
	}
	return []string{"PATH=" + stubDir + string(os.PathListSeparator) + os.Getenv("PATH")}
}
