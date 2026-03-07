package contract_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestContractSync(t *testing.T) {
	bin := buildCLI(t)
	_, env := writeRBWStub(t, "unlocked", map[string]string{
		"database-url": "postgres://vault/dev",
		"jwt-secret":   "topsecret",
	})
	project := setupProject(t, "# Application settings\nDATABASE_URL=\nJWT_SECRET=\nPORT=8080\n", "", "mapping:\n  DATABASE_URL: database-url\n  JWT_SECRET: jwt-secret\n")

	stdout, stderr, code := runCLI(t, bin, project, env, "sync", "--dry-run")
	if code != 0 || stderr != "" {
		t.Fatalf("dry-run sync failed: code=%d stderr=%q", code, stderr)
	}
	for _, want := range []string{"ADD DATABASE_URL [RESOLVED]", "ADD JWT_SECRET [RESOLVED]", "ADD PORT [STATIC]", "CHECKED "} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("dry-run output missing %q\n%s", want, stdout)
		}
	}
	if strings.Contains(stdout, "topsecret") || strings.Contains(stdout, "postgres://vault/dev") {
		t.Fatalf("dry-run leaked secret output: %s", stdout)
	}

	stdout, stderr, code = runCLI(t, bin, project, env, "sync")
	if code != 0 || stderr != "" {
		t.Fatalf("sync failed: code=%d stderr=%q", code, stderr)
	}
	data, err := os.ReadFile(filepath.Join(project, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, want := range []string{"DATABASE_URL=postgres://vault/dev", "JWT_SECRET=topsecret", "PORT=8080"} {
		if !strings.Contains(content, want) {
			t.Fatalf("env file missing %q\n%s", want, content)
		}
	}
	if strings.Contains(stdout, "topsecret") || strings.Contains(stdout, "postgres://vault/dev") {
		t.Fatalf("sync stdout leaked secret output: %s", stdout)
	}

	stdout, stderr, code = runCLI(t, bin, project, env, "sync")
	if code != 0 || stderr != "" {
		t.Fatalf("second sync failed: code=%d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "UNCHANGED") {
		t.Fatalf("expected unchanged summary, got:\n%s", stdout)
	}
}
