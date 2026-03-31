package integration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncCommandIntegration(t *testing.T) {
	bin := buildCLI(t)
	project := t.TempDir()
	itemName := filepath.Base(project)
	env := writeRBWStub(t, "unlocked", map[string]string{
		rbwLookupKey(itemName, "DATABASE_URL"): "postgres://vault/dev",
		rbwLookupKey(itemName, "JWT_SECRET"):   "supersecret",
	})
	writeFile(t, filepath.Join(project, ".env.example"), "# heading\r\nDATABASE_URL=\r\nJWT_SECRET=\r\nPORT=8080\r\n")

	stdout, stderr, code := runCLI(t, bin, project, env, "sync")
	if code != 0 || stderr != "" {
		t.Fatalf("sync failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
	}
	data, err := os.ReadFile(filepath.Join(project, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "\r\n") {
		t.Fatalf("expected CRLF preservation, got %q", content)
	}
	if !strings.Contains(content, "JWT_SECRET=supersecret") {
		t.Fatalf("sync content missing resolved value: %s", content)
	}
	for _, want := range []string{"ADD DATABASE_URL [RESOLVED]", "ADD JWT_SECRET [RESOLVED]", "ADD PORT [STATIC]", "WRITTEN "} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("initial sync output missing %q\n%s", want, stdout)
		}
	}

	writeFile(t, filepath.Join(project, ".env"), "# heading\r\nDATABASE_URL=\r\nJWT_SECRET=stale\r\nPORT=8080\r\n")
	stdout, stderr, code = runCLI(t, bin, project, env, "sync")
	if code != 0 || stderr != "" {
		t.Fatalf("update sync failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
	}
	after, err := os.ReadFile(filepath.Join(project, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(after), "JWT_SECRET=supersecret") {
		t.Fatalf("expected updated secret after sync, got %q", string(after))
	}
	if !strings.Contains(stdout, "UPDATE JWT_SECRET [RESOLVED]") {
		t.Fatalf("expected updated-key output, got %q", stdout)
	}

	before := string(after)
	stdout, stderr, code = runCLI(t, bin, project, env, "sync")
	if code != 0 || stderr != "" {
		t.Fatalf("third sync failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
	}
	final, err := os.ReadFile(filepath.Join(project, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	if string(final) != before {
		t.Fatalf("expected no-op file stability")
	}
}
