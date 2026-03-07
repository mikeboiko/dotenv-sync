package contract_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func writeRBWStub(t *testing.T, status string, get map[string]string, missing ...string) (string, []string) {
	t.Helper()
	stubDir := t.TempDir()
	missingSet := map[string]bool{}
	for _, key := range missing {
		missingSet[key] = true
	}
	var script strings.Builder
	script.WriteString("#!/bin/sh\n")
	script.WriteString("case \"$1\" in\n")
	script.WriteString(fmt.Sprintf("status) printf '%%s\\n' '%s' ;;\n", status))
	script.WriteString("get)\ncase \"$2\" in\n")
	for key, value := range get {
		script.WriteString(fmt.Sprintf("%s) printf '%%s\\n' '%s' ;;\n", key, value))
	}
	for key := range missingSet {
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
	return path, []string{"PATH=" + stubDir + string(os.PathListSeparator) + os.Getenv("PATH")}
}
