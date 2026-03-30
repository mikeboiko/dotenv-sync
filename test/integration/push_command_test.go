package integration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPushCommandIntegration(t *testing.T) {
	bin := buildCLI(t)
	itemName := "Jesse"
	configYAML := "storage_mode: note_json\nitem_name: Jesse\n"

	t.Run("dry run create unchanged and update flows", func(t *testing.T) {
		project := setupProject(t,
			"DATABASE_URL=\nJWT_SECRET=\nPORT=8080\n",
			"DATABASE_URL=postgres://vault/dev\nJWT_SECRET=supersecret\nPORT=8080\n",
			configYAML,
		)
		stub := writeRBWStubWithOptions(t, rbwStubOptions{Status: "unlocked"})

		stdout, stderr, code := runCLI(t, bin, project, stub.Env(), "push", "--dry-run")
		if code != 0 || stderr != "" {
			t.Fatalf("push dry-run failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
		}
		wantDryRun := renderTemplate(readRepoFile(t, "test", "testdata", "golden", "push-dry-run.txt"), map[string]string{"{{ITEM}}": itemName})
		if stdout != wantDryRun {
			t.Fatalf("push dry-run stdout = %q want %q", stdout, wantDryRun)
		}
		if log := stub.Log(t); strings.Contains(log, "add ") || strings.Contains(log, "edit ") || strings.Contains(log, "sync") {
			t.Fatalf("dry-run should not mutate provider, log=%q", log)
		}

		stdout, stderr, code = runCLI(t, bin, project, stub.Env(), "push")
		if code != 0 || stderr != "" {
			t.Fatalf("push create failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
		}
		if !strings.Contains(stdout, "WRITTEN bitwarden:Jesse (added: 3)") {
			t.Fatalf("unexpected create summary: %s", stdout)
		}
		wantNotes := strings.TrimSpace(readRepoFile(t, "test", "testdata", "provider", "note-json-valid.json"))
		if notes := strings.TrimSpace(stub.Note(t, itemName)); notes != wantNotes {
			t.Fatalf("unexpected created notes = %q want %q", notes, wantNotes)
		}
		if got := strings.Count(stub.Log(t), "add Jesse"); got != 1 {
			t.Fatalf("expected one add, got %d log=%q", got, stub.Log(t))
		}
		if got := strings.Count(stub.Log(t), "sync"); got != 1 {
			t.Fatalf("expected one sync, got %d log=%q", got, stub.Log(t))
		}

		stdout, stderr, code = runCLI(t, bin, project, stub.Env(), "push")
		if code != 0 || stderr != "" {
			t.Fatalf("push unchanged failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
		}
		wantUnchanged := renderTemplate(readRepoFile(t, "test", "testdata", "golden", "push-unchanged.txt"), map[string]string{"{{ITEM}}": itemName})
		if stdout != wantUnchanged {
			t.Fatalf("push unchanged stdout = %q want %q", stdout, wantUnchanged)
		}
		if got := strings.Count(stub.Log(t), "add Jesse"); got != 1 {
			t.Fatalf("unchanged run should not add again, log=%q", stub.Log(t))
		}
		if got := strings.Count(stub.Log(t), "edit Jesse"); got != 0 {
			t.Fatalf("unchanged run should not edit, log=%q", stub.Log(t))
		}

		writeFile(t, filepath.Join(project, ".env"), "DATABASE_URL=postgres://vault/dev\nJWT_SECRET=rotated\nPORT=8080\n")
		stdout, stderr, code = runCLI(t, bin, project, stub.Env(), "push")
		if code != 0 || stderr != "" {
			t.Fatalf("push update failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
		}
		if !strings.Contains(stdout, "WRITTEN bitwarden:Jesse (updated: 1, unchanged: 2)") {
			t.Fatalf("unexpected update summary: %s", stdout)
		}
		if !strings.Contains(stub.Note(t, itemName), "\"JWT_SECRET\":\"rotated\"") {
			t.Fatalf("expected rotated secret in notes, got %q", stub.Note(t, itemName))
		}
		if got := strings.Count(stub.Log(t), "edit Jesse"); got != 1 {
			t.Fatalf("expected one edit, got %d log=%q", got, stub.Log(t))
		}
		if got := strings.Count(stub.Log(t), "sync"); got != 2 {
			t.Fatalf("expected second sync after update, got %d log=%q", got, stub.Log(t))
		}
	})

	t.Run("fields mode rejection and extra previews", func(t *testing.T) {
		project := setupProject(t,
			"DATABASE_URL=\n",
			"DATABASE_URL=postgres://vault/dev\nEXTRA_FLAG=true\n",
			"item_name: Jesse\n",
		)
		stub := writeRBWStubWithOptions(t, rbwStubOptions{Status: "unlocked"})

		stdout, stderr, code := runCLI(t, bin, project, stub.Env(), "push")
		if code != 1 || !strings.Contains(stderr, "ERROR E009") {
			t.Fatalf("expected fields mode error, code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		if strings.TrimSpace(stub.Log(t)) != "" {
			t.Fatalf("expected no provider calls, log=%q", stub.Log(t))
		}

		writeFile(t, filepath.Join(project, ".envsync.yaml"), configYAML)
		stdout, stderr, code = runCLI(t, bin, project, stub.Env(), "push", "--dry-run")
		if code != 0 || stderr != "" {
			t.Fatalf("push extra preview failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
		}
		if !strings.Contains(stdout, "EXTRA EXTRA_FLAG [REDACTED]") {
			t.Fatalf("expected extra preview, got %s", stdout)
		}
		data, err := os.ReadFile(filepath.Join(project, ".env.example"))
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "DATABASE_URL=\n" {
			t.Fatalf("schema should remain unchanged, got %q", string(data))
		}
	})
}
