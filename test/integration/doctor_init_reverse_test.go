package integration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorInitAndReverseIntegration(t *testing.T) {
	bin := buildCLI(t)

	t.Run("doctor readiness failure", func(t *testing.T) {
		env := writeRBWStub(t, "locked", map[string]string{})
		project := t.TempDir()
		writeFile(t, filepath.Join(project, ".env.example"), "API_KEY=\n")
		_, stderr, code := runCLI(t, bin, project, env, "doctor")
		if code != 1 || !strings.Contains(stderr, "ERROR E003") {
			t.Fatalf("doctor stderr=%q code=%d", stderr, code)
		}
	})

	t.Run("init and reverse keep schema secret-safe", func(t *testing.T) {
		project := t.TempDir()
		writeFile(t, filepath.Join(project, ".env"), "DATABASE_URL=postgres://localhost/dev\nPORT=8080\nNEW_SECRET=value\n")
		stdout, stderr, code := runCLI(t, bin, project, nil, "init")
		if code != 0 || stderr != "" {
			t.Fatalf("init failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
		}
		data, err := os.ReadFile(filepath.Join(project, ".env.example"))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "postgres://") || strings.Contains(string(data), "value") {
			t.Fatalf("init wrote secret-looking values: %s", data)
		}

		writeFile(t, filepath.Join(project, ".env"), "DATABASE_URL=postgres://localhost/dev\nPORT=8080\nNEW_SECRET=value\nEXTRA_FLAG=true\n")
		stdout, stderr, code = runCLI(t, bin, project, nil, "reverse")
		if code != 0 || stderr != "" {
			t.Fatalf("reverse failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
		}
		data, err = os.ReadFile(filepath.Join(project, ".env.example"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "EXTRA_FLAG=") || strings.Contains(string(data), "true") {
			t.Fatalf("reverse schema content unexpected: %s", data)
		}
	})

	t.Run("init duplicate key surfaces error details", func(t *testing.T) {
		project := t.TempDir()
		writeFile(t, filepath.Join(project, ".env"), "API_KEY=first\nAPI_KEY=second\n")
		stdout, stderr, code := runCLI(t, bin, project, nil, "init")
		if code != 2 {
			t.Fatalf("init duplicate code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		if !strings.Contains(stderr, "ERROR E008") || !strings.Contains(stderr, "duplicate key detected: API_KEY") {
			t.Fatalf("unexpected init duplicate stderr=%q", stderr)
		}
	})
}
