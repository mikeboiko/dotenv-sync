package contract_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestContractDoctorInitReverse(t *testing.T) {
	bin := buildCLI(t)

	t.Run("doctor reports locked rbw", func(t *testing.T) {
		_, env := writeRBWStub(t, "locked", map[string]string{})
		project := setupProject(t, "DATABASE_URL=\n", "", "")
		stdout, stderr, code := runCLI(t, bin, project, env, "doctor")
		if code != 1 {
			t.Fatalf("doctor exit code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		for _, want := range []string{"ERROR E003", "Problem: Bitwarden database is locked", "Action: run 'rbw unlock' and retry"} {
			if !strings.Contains(stderr, want) {
				t.Fatalf("doctor stderr missing %q\n%s", want, stderr)
			}
		}
	})

	t.Run("init dry run blanks secrets", func(t *testing.T) {
		project := setupProject(t, "", "DATABASE_URL=postgres://localhost/dev\nJWT_SECRET=supersecret\nPORT=8080\n", "")
		stdout, stderr, code := runCLI(t, bin, project, nil, "init", "--dry-run")
		if code != 0 || stderr != "" {
			t.Fatalf("init dry-run failed: code=%d stderr=%q", code, stderr)
		}
		for _, want := range []string{"ADD DATABASE_URL [MISSING]", "ADD JWT_SECRET [MISSING]", "ADD PORT [STATIC]"} {
			if !strings.Contains(stdout, want) {
				t.Fatalf("init output missing %q\n%s", want, stdout)
			}
		}
		if strings.Contains(stdout, "supersecret") {
			t.Fatalf("init leaked secret in stdout: %s", stdout)
		}
	})

	t.Run("reverse dry run adds blank placeholders", func(t *testing.T) {
		project := setupProject(t, "PORT=8080\n", "PORT=8080\nNEW_API_KEY=abc123\n", "")
		stdout, stderr, code := runCLI(t, bin, project, nil, "reverse", "--dry-run")
		if code != 0 || stderr != "" {
			t.Fatalf("reverse dry-run failed: code=%d stderr=%q", code, stderr)
		}
		if !strings.Contains(stdout, "ADD NEW_API_KEY [MISSING]") {
			t.Fatalf("reverse output missing placeholder add:\n%s", stdout)
		}
		if strings.Contains(stdout, "abc123") {
			t.Fatalf("reverse leaked env value: %s", stdout)
		}
	})

	t.Run("reverse writes blank placeholder only", func(t *testing.T) {
		project := setupProject(t, "PORT=8080\n", "PORT=8080\nNEW_API_KEY=abc123\n", "")
		_, _, code := runCLI(t, bin, project, nil, "reverse")
		if code != 0 {
			t.Fatalf("reverse exit code=%d", code)
		}
		data, err := os.ReadFile(filepath.Join(project, ".env.example"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "NEW_API_KEY=") || strings.Contains(string(data), "abc123") {
			t.Fatalf("reverse schema content unexpected:\n%s", data)
		}
	})
}
