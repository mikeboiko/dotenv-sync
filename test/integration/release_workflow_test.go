package integration_test

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseWorkflowIntegration(t *testing.T) {
	scriptPath := "./scripts/nextversion"
	moduleRoot := repoRoot(t)

	t.Run("patch baseline with no tags", func(t *testing.T) {
		project := t.TempDir()
		initGitRepo(t, project)
		writeFile(t, filepath.Join(project, "README.md"), "release tests\n")
		commitAll(t, project, "initial")

		stdout, stderr, code := runGoMain(t, moduleRoot, scriptPath, "--bump", "patch", "--dir", project)
		if code != 0 || stderr != "" {
			t.Fatalf("nextversion patch failed: code=%d stderr=%q", code, stderr)
		}
		if strings.TrimSpace(stdout) != "v0.0.1" {
			t.Fatalf("patch baseline = %q", stdout)
		}
	})

	t.Run("minor bump uses latest semver tag", func(t *testing.T) {
		project := t.TempDir()
		initGitRepo(t, project)
		writeFile(t, filepath.Join(project, "README.md"), "release tests\n")
		commitAll(t, project, "initial")
		if _, stderr, code := runCommand(t, project, "git", "tag", "v0.4.2"); code != 0 {
			t.Fatalf("tag v0.4.2 failed: %s", stderr)
		}
		if _, stderr, code := runCommand(t, project, "git", "tag", "notes"); code != 0 {
			t.Fatalf("tag notes failed: %s", stderr)
		}

		stdout, stderr, code := runGoMain(t, moduleRoot, scriptPath, "--bump", "minor", "--dir", project)
		if code != 0 || stderr != "" {
			t.Fatalf("nextversion minor failed: code=%d stderr=%q", code, stderr)
		}
		if strings.TrimSpace(stdout) != "v0.5.0" {
			t.Fatalf("minor bump = %q", stdout)
		}
	})

	t.Run("invalid bump fails with actionable error", func(t *testing.T) {
		project := t.TempDir()
		initGitRepo(t, project)
		writeFile(t, filepath.Join(project, "README.md"), "release tests\n")
		commitAll(t, project, "initial")

		_, stderr, code := runGoMain(t, moduleRoot, scriptPath, "--bump", "build", "--dir", project)
		if code == 0 {
			t.Fatal("expected invalid bump failure")
		}
		if !strings.Contains(stderr, "unsupported bump") {
			t.Fatalf("unexpected stderr: %q", stderr)
		}
	})
}
