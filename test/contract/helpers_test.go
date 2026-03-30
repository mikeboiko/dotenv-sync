package contract_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, "../.."))
}

func buildCLI(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "ds")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/ds")
	cmd.Dir = repoRoot(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build cli: %v\n%s", err, out)
	}
	return bin
}

func buildCLIWithLdflags(t *testing.T, version, commit, buildTime string) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "ds")
	ldflags := fmt.Sprintf("-X dotenv-sync/pkg/dotenvsync.Version=%s -X dotenv-sync/pkg/dotenvsync.Commit=%s -X dotenv-sync/pkg/dotenvsync.BuildTime=%s", version, commit, buildTime)
	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", bin, "./cmd/ds")
	cmd.Dir = repoRoot(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build cli with ldflags: %v\n%s", err, out)
	}
	return bin
}

func runCLI(t *testing.T, bin, dir string, extraEnv []string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), extraEnv...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		return stdout.String(), stderr.String(), 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return stdout.String(), stderr.String(), exitErr.ExitCode()
	}
	t.Fatalf("run cli: %v", err)
	return "", "", 0
}

func readRepoFile(t *testing.T, parts ...string) string {
	t.Helper()
	path := filepath.Join(append([]string{repoRoot(t)}, parts...)...)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func readGoldenFile(t *testing.T, name string) string {
	t.Helper()
	return readRepoFile(t, "test", "testdata", "golden", name)
}

func readReleaseFixtureLines(t *testing.T, name string) []string {
	t.Helper()
	content := readRepoFile(t, "test", "testdata", "release", name)
	lines := strings.Split(content, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		filtered = append(filtered, line)
	}
	return filtered
}

func renderTemplate(input string, replacements map[string]string) string {
	output := input
	for old, newValue := range replacements {
		output = strings.ReplaceAll(output, old, newValue)
	}
	return output
}

func compactJSON(t *testing.T, input string) string {
	t.Helper()
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(input)); err != nil {
		t.Fatalf("compact json: %v", err)
	}
	return buf.String()
}

func currentPlatform() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func setupProject(t *testing.T, schema, env, cfg string) string {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".env.example"), schema)
	if env != "" {
		writeFile(t, filepath.Join(dir, ".env"), env)
	}
	if cfg != "" {
		writeFile(t, filepath.Join(dir, ".envsync.yaml"), cfg)
	}
	return dir
}

func rbwLookupKey(item, field string) string {
	return item + "::" + field
}

type rbwStubItem struct {
	Notes    string
	Password string
}

type rbwStubOptions struct {
	Status  string
	Fields  map[string]string
	Missing []string
	Items   map[string]rbwStubItem
}

type rbwStub struct {
	path    string
	env     []string
	logFile string
	items   string
}

func (s rbwStub) Env() []string {
	return append([]string{}, s.env...)
}

func (s rbwStub) Log(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(s.logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		t.Fatal(err)
	}
	return string(data)
}

func (s rbwStub) Note(t *testing.T, item string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(s.items, item+".notes"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func (s rbwStub) Password(t *testing.T, item string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(s.items, item+".password"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func writeRBWStub(t *testing.T, status string, get map[string]string, missing ...string) (string, []string) {
	stub := writeRBWStubWithOptions(t, rbwStubOptions{Status: status, Fields: get, Missing: missing})
	return stub.path, stub.Env()
}

func writeRBWStubWithOptions(t *testing.T, opts rbwStubOptions) rbwStub {
	t.Helper()
	stubDir := t.TempDir()
	itemsDir := filepath.Join(stubDir, "items")
	logFile := filepath.Join(stubDir, "rbw.log")
	if err := os.MkdirAll(itemsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for item, state := range opts.Items {
		writeFile(t, filepath.Join(itemsDir, item+".notes"), state.Notes)
		writeFile(t, filepath.Join(itemsDir, item+".password"), state.Password)
	}
	missingSet := map[string]bool{}
	for _, key := range opts.Missing {
		missingSet[key] = true
	}
	var script strings.Builder
	script.WriteString("#!/bin/sh\n")
	script.WriteString("LOG_FILE='" + logFile + "'\n")
	script.WriteString("ITEMS_DIR='" + itemsDir + "'\n")
	script.WriteString("echo \"$@\" >> \"$LOG_FILE\"\n")
	script.WriteString("item_file() { printf '%s/%s.%s' \"$ITEMS_DIR\" \"$1\" \"$2\"; }\n")
	script.WriteString("write_from_editor() {\n")
	script.WriteString("  item=\"$1\"\n")
	script.WriteString("  tmp=$(mktemp)\n")
	script.WriteString("  \"${VISUAL:-$EDITOR}\" \"$tmp\"\n")
	script.WriteString("  python3 - \"$tmp\" \"$(item_file \"$item\" password)\" \"$(item_file \"$item\" notes)\" <<'PY'\n")
	script.WriteString("from pathlib import Path\n")
	script.WriteString("import sys\n")
	script.WriteString("src, password_path, notes_path = sys.argv[1:4]\n")
	script.WriteString("lines = Path(src).read_text(encoding='utf-8').splitlines()\n")
	script.WriteString("password = lines[0] if lines else ''\n")
	script.WriteString("rest = lines[1:]\n")
	script.WriteString("while rest and rest[0] == '':\n")
	script.WriteString("    rest = rest[1:]\n")
	script.WriteString("notes = '\\n'.join(line for line in rest if not line.startswith('#'))\n")
	script.WriteString("Path(password_path).write_text(password, encoding='utf-8')\n")
	script.WriteString("Path(notes_path).write_text(notes, encoding='utf-8')\n")
	script.WriteString("PY\n")
	script.WriteString("}\n")
	script.WriteString("emit_raw() {\n")
	script.WriteString("  item=\"$1\"\n")
	script.WriteString("  python3 - \"$item\" \"$(item_file \"$item\" password)\" \"$(item_file \"$item\" notes)\" <<'PY'\n")
	script.WriteString("from pathlib import Path\n")
	script.WriteString("import json, os, sys\n")
	script.WriteString("item, password_path, notes_path = sys.argv[1:4]\n")
	script.WriteString("password = Path(password_path).read_text(encoding='utf-8') if os.path.exists(password_path) else ''\n")
	script.WriteString("notes = Path(notes_path).read_text(encoding='utf-8') if os.path.exists(notes_path) else ''\n")
	script.WriteString("print(json.dumps({'name': item, 'notes': notes, 'data': {'password': password}}))\n")
	script.WriteString("PY\n")
	script.WriteString("}\n")
	script.WriteString("case \"$1\" in\n")
	switch opts.Status {
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
		script.WriteString(fmt.Sprintf("list) printf '%%s\\n' '%s' ;;\n", opts.Status))
	}
	script.WriteString(fmt.Sprintf("status) printf '%%s\\n' '%s' ;;\n", opts.Status))
	script.WriteString("get)\n")
	script.WriteString("field=''\n")
	script.WriteString("item=''\n")
	script.WriteString("shift\n")
	script.WriteString("if [ \"$1\" = \"--raw\" ]; then\n")
	script.WriteString("  item=\"$2\"\n")
	script.WriteString("  if [ -f \"$(item_file \"$item\" notes)\" ] || [ -f \"$(item_file \"$item\" password)\" ]; then\n")
	script.WriteString("    emit_raw \"$item\"\n")
	script.WriteString("    exit 0\n")
	script.WriteString("  fi\n")
	script.WriteString("  echo 'not found' >&2\n")
	script.WriteString("  exit 1\n")
	script.WriteString("fi\n")
	script.WriteString("while [ \"$#\" -gt 0 ]; do\n")
	script.WriteString("  case \"$1\" in\n")
	script.WriteString("    --field) field=\"$2\"; shift 2 ;;\n")
	script.WriteString("    --*) shift ;;\n")
	script.WriteString("    *) item=\"$1\"; shift ;;\n")
	script.WriteString("  esac\n")
	script.WriteString("done\n")
	script.WriteString("case \"$item::$field\" in\n")
	for key, value := range opts.Fields {
		script.WriteString(fmt.Sprintf("%s) printf '%%s\\n' '%s' ;;\n", key, value))
	}
	for key := range missingSet {
		script.WriteString(fmt.Sprintf("%s) echo 'not found' >&2; exit 1 ;;\n", key))
	}
	script.WriteString("*) echo 'not found' >&2; exit 1 ;;\nesac\n;;\n")
	script.WriteString("add)\nshift\nwrite_from_editor \"$1\"\n;;\n")
	script.WriteString("edit)\nshift\nif [ ! -f \"$(item_file \"$1\" notes)\" ] && [ ! -f \"$(item_file \"$1\" password)\" ]; then echo 'not found' >&2; exit 1; fi\nwrite_from_editor \"$1\"\n;;\n")
	script.WriteString("sync) exit 0 ;;\n")
	script.WriteString("*) echo 'unsupported' >&2; exit 1 ;;\nesac\n")
	path := filepath.Join(stubDir, "rbw")
	writeFile(t, path, script.String())
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatal(err)
	}
	return rbwStub{
		path:    path,
		env:     []string{"PATH=" + stubDir + string(os.PathListSeparator) + os.Getenv("PATH")},
		logFile: logFile,
		items:   itemsDir,
	}
}
