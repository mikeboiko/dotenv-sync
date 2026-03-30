package contract_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestContractPush(t *testing.T) {
	bin := buildCLI(t)
	itemName := "Jesse"
	configYAML := "storage_mode: note_json\nitem_name: Jesse\n"
	t.Run("dry run previews redacted changes", func(t *testing.T) {
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
		want := renderTemplate(readRepoFile(t, "test", "testdata", "golden", "push-dry-run.txt"), map[string]string{"{{ITEM}}": itemName})
		if stdout != want {
			t.Fatalf("push dry-run stdout = %q want %q", stdout, want)
		}
		if strings.Contains(stdout, "supersecret") {
			t.Fatalf("push dry-run leaked secret: %s", stdout)
		}
		if log := stub.Log(t); strings.Contains(log, "add ") || strings.Contains(log, "edit ") || strings.Contains(log, "sync") {
			t.Fatalf("dry-run should not mutate provider, log=%q", log)
		}
	})

	t.Run("create writes repo-scoped note_json payload", func(t *testing.T) {
		project := setupProject(t,
			"DATABASE_URL=\nJWT_SECRET=\nPORT=8080\n",
			"DATABASE_URL=postgres://vault/dev\nJWT_SECRET=supersecret\nPORT=8080\n",
			configYAML,
		)
		stub := writeRBWStubWithOptions(t, rbwStubOptions{Status: "unlocked"})

		stdout, stderr, code := runCLI(t, bin, project, stub.Env(), "push")
		if code != 0 || stderr != "" {
			t.Fatalf("push failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
		}
		if !strings.Contains(stdout, "WRITTEN bitwarden:Jesse (added: 3)") {
			t.Fatalf("unexpected push summary: %s", stdout)
		}
		notes := strings.TrimSpace(stub.Note(t, itemName))
		wantNotes := strings.TrimSpace(readRepoFile(t, "test", "testdata", "provider", "note-json-valid.json"))
		if notes != wantNotes {
			t.Fatalf("unexpected provider notes = %q want %q", notes, wantNotes)
		}
		if got := stub.Password(t, itemName); got != "" {
			t.Fatalf("expected blank password, got %q", got)
		}
		log := stub.Log(t)
		if !strings.Contains(log, "add Jesse") || !strings.Contains(log, "sync") {
			t.Fatalf("expected add + sync in log, got %q", log)
		}
	})

	t.Run("unchanged skips provider mutation", func(t *testing.T) {
		project := setupProject(t,
			"DATABASE_URL=\nJWT_SECRET=\nPORT=8080\n",
			"DATABASE_URL=postgres://vault/dev\nJWT_SECRET=supersecret\nPORT=8080\n",
			configYAML,
		)
		stub := writeRBWStubWithOptions(t, rbwStubOptions{
			Status: "unlocked",
			Items: map[string]rbwStubItem{
				itemName: {Notes: strings.TrimSpace(readRepoFile(t, "test", "testdata", "provider", "note-json-valid.json")), Password: "keep-me"},
			},
		})

		stdout, stderr, code := runCLI(t, bin, project, stub.Env(), "push")
		if code != 0 || stderr != "" {
			t.Fatalf("push unchanged failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
		}
		want := renderTemplate(readRepoFile(t, "test", "testdata", "golden", "push-unchanged.txt"), map[string]string{"{{ITEM}}": itemName})
		if stdout != want {
			t.Fatalf("push unchanged stdout = %q want %q", stdout, want)
		}
		if log := stub.Log(t); strings.Contains(log, "add ") || strings.Contains(log, "edit ") || strings.Contains(log, "sync") {
			t.Fatalf("unchanged run should not mutate provider, log=%q", log)
		}
	})

	t.Run("fields mode explains migration path", func(t *testing.T) {
		project := setupProject(t,
			"DATABASE_URL=\n",
			"DATABASE_URL=postgres://vault/dev\n",
			"item_name: Jesse\n",
		)
		stub := writeRBWStubWithOptions(t, rbwStubOptions{Status: "unlocked"})

		stdout, stderr, code := runCLI(t, bin, project, stub.Env(), "push")
		if code != 1 {
			t.Fatalf("push fields mode exit code=%d stdout=%q stderr=%q", code, stdout, stderr)
		}
		want := readRepoFile(t, "test", "testdata", "golden", "push-fields-mode-error.txt")
		if stderr != want {
			t.Fatalf("push fields stderr = %q want %q", stderr, want)
		}
		if stdout != "" {
			t.Fatalf("expected empty stdout, got %q", stdout)
		}
		if strings.TrimSpace(stub.Log(t)) != "" {
			t.Fatalf("expected no provider calls in fields mode, log=%q", stub.Log(t))
		}
	})

	t.Run("extra env keys stay redacted and do not touch schema", func(t *testing.T) {
		project := setupProject(t,
			"DATABASE_URL=\n",
			"DATABASE_URL=postgres://vault/dev\nEXTRA_FLAG=true\n",
			configYAML,
		)
		stub := writeRBWStubWithOptions(t, rbwStubOptions{Status: "unlocked"})

		stdout, stderr, code := runCLI(t, bin, project, stub.Env(), "push", "--dry-run")
		if code != 0 || stderr != "" {
			t.Fatalf("push extra dry-run failed: code=%d stderr=%q stdout=%q", code, stderr, stdout)
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
