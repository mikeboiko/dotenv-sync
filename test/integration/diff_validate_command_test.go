package integration_test

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDiffValidateAndMissingIntegration(t *testing.T) {
	bin := buildCLI(t)

	t.Run("drift and missing secrets", func(t *testing.T) {
		project := t.TempDir()
		itemName := filepath.Base(project)
		env := writeRBWStub(t, "unlocked", map[string]string{
			rbwLookupKey(itemName, "database-url"): "postgres://vault/dev",
		}, rbwLookupKey(itemName, "jwt-secret"))
		writeFile(t, filepath.Join(project, ".env.example"), "DATABASE_URL=\nJWT_SECRET=\nPORT=8080\n")
		writeFile(t, filepath.Join(project, ".env"), "DATABASE_URL=postgres://vault/dev\nPORT=9090\nEXTRA_KEY=value\n")
		writeFile(t, filepath.Join(project, ".envsync.yaml"), "mapping:\n  DATABASE_URL: database-url\n  JWT_SECRET: jwt-secret\n")

		stdout, _, code := runCLI(t, bin, project, env, "diff")
		if code != 2 {
			t.Fatalf("diff exit code=%d stdout=%s", code, stdout)
		}
		for _, want := range []string{"MISSING JWT_SECRET [MISSING]", "UPDATE PORT [STATIC]", "EXTRA EXTRA_KEY"} {
			if !strings.Contains(stdout, want) {
				t.Fatalf("diff output missing %q\n%s", want, stdout)
			}
		}

		stdout, _, code = runCLI(t, bin, project, env, "missing")
		if code != 2 || !strings.Contains(stdout, "MISSING JWT_SECRET") {
			t.Fatalf("missing output unexpected: code=%d stdout=%s", code, stdout)
		}
	})

	t.Run("duplicate and malformed schema", func(t *testing.T) {
		project := t.TempDir()
		itemName := filepath.Base(project)
		env := writeRBWStub(t, "unlocked", map[string]string{
			rbwLookupKey(itemName, "DATABASE_URL"): "value",
		})
		writeFile(t, filepath.Join(project, ".env.example"), "DATABASE_URL=\nDATABASE_URL=\nBAD LINE\n")
		writeFile(t, filepath.Join(project, ".env"), "DATABASE_URL=value\n")
		stdout, _, code := runCLI(t, bin, project, env, "validate")
		if code != 2 {
			t.Fatalf("validate exit code=%d stdout=%s", code, stdout)
		}
		for _, want := range []string{"E006", "E008 DATABASE_URL"} {
			if !strings.Contains(stdout, want) {
				t.Fatalf("validate output missing %q\n%s", want, stdout)
			}
		}
	})
}
