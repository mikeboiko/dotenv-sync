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

		stdout, stderr, code := runGoMain(t, moduleRoot, scriptPath, "--dir", project)
		if code != 0 || stderr != "" {
			t.Fatalf("nextversion patch failed: code=%d stderr=%q", code, stderr)
		}
		if strings.TrimSpace(stdout) != "v0.0.1" {
			t.Fatalf("patch baseline = %q", stdout)
		}
	})

	t.Run("successive main pushes advance patch versions", func(t *testing.T) {
		expected := readReleaseFixtureLines(t, "concurrent-main-push.txt")
		project := t.TempDir()
		initGitRepo(t, project)
		writeFile(t, filepath.Join(project, "README.md"), "release tests\n")
		commitAll(t, project, "initial")
		if _, stderr, code := runCommand(t, project, "git", "tag", "v0.4.2"); code != 0 {
			t.Fatalf("tag v0.4.2 failed: %s", stderr)
		}

		writeFile(t, filepath.Join(project, "CHANGELOG.md"), "release 1\n")
		commitAll(t, project, "release 1")
		stdout, stderr, code := runGoMain(t, moduleRoot, scriptPath, "--dir", project)
		if code != 0 || stderr != "" {
			t.Fatalf("nextversion first patch failed: code=%d stderr=%q", code, stderr)
		}
		if strings.TrimSpace(stdout) != expected[0] {
			t.Fatalf("first patch preview = %q", stdout)
		}
		if _, stderr, code := runCommand(t, project, "git", "tag", expected[0]); code != 0 {
			t.Fatalf("tag %s failed: %s", expected[0], stderr)
		}

		writeFile(t, filepath.Join(project, "NEXT.md"), "release 2\n")
		commitAll(t, project, "release 2")
		stdout, stderr, code = runGoMain(t, moduleRoot, scriptPath, "--dir", project)
		if code != 0 || stderr != "" {
			t.Fatalf("nextversion second patch failed: code=%d stderr=%q", code, stderr)
		}
		if strings.TrimSpace(stdout) != expected[1] {
			t.Fatalf("second patch preview = %q", stdout)
		}
	})

	t.Run("already tagged commit returns released exit code", func(t *testing.T) {
		project := t.TempDir()
		initGitRepo(t, project)
		writeFile(t, filepath.Join(project, "README.md"), "release tests\n")
		commitAll(t, project, "initial")
		releasedTag := readReleaseFixtureLines(t, "already-released.txt")[0]
		if _, stderr, code := runCommand(t, project, "git", "tag", releasedTag); code != 0 {
			t.Fatalf("tag %s failed: %s", releasedTag, stderr)
		}

		stdout, stderr, code := runGoMain(t, moduleRoot, scriptPath, "--dir", project)
		if code != 2 {
			t.Fatalf("already released exit code = %d stderr=%q", code, stderr)
		}
		if strings.TrimSpace(stdout) != releasedTag {
			t.Fatalf("already released stdout = %q", stdout)
		}
		wantStderr := strings.TrimSpace(strings.ReplaceAll(readGoldenFile(t, "release-skipped.txt"), "{{VERSION}}", releasedTag))
		if strings.TrimSpace(stderr) != wantStderr {
			t.Fatalf("already released stderr = %q want %q", stderr, wantStderr)
		}
	})

	t.Run("non semver tags are ignored", func(t *testing.T) {
		project := t.TempDir()
		initGitRepo(t, project)
		writeFile(t, filepath.Join(project, "README.md"), "release tests\n")
		commitAll(t, project, "initial")
		if _, stderr, code := runCommand(t, project, "git", "tag", "v0.4.2"); code != 0 {
			t.Fatalf("tag v0.4.2 failed: %s", stderr)
		}
		for _, tag := range []string{"release-2026-03-30", "notes", "v2.0.0-beta"} {
			if _, stderr, code := runCommand(t, project, "git", "tag", tag); code != 0 {
				t.Fatalf("tag %s failed: %s", tag, stderr)
			}
		}
		writeFile(t, filepath.Join(project, "CHANGELOG.md"), "next change\n")
		commitAll(t, project, "next")

		stdout, stderr, code := runGoMain(t, moduleRoot, scriptPath, "--dir", project)
		if code != 0 || stderr != "" {
			t.Fatalf("non-semver preview failed: code=%d stderr=%q", code, stderr)
		}
		if strings.TrimSpace(stdout) != "v0.4.3" {
			t.Fatalf("non-semver preview = %q", stdout)
		}
	})

	t.Run("positional args fail with actionable error", func(t *testing.T) {
		project := t.TempDir()
		initGitRepo(t, project)
		writeFile(t, filepath.Join(project, "README.md"), "release tests\n")
		commitAll(t, project, "initial")

		_, stderr, code := runGoMain(t, moduleRoot, scriptPath, "--dir", project, "unexpected")
		if code == 0 {
			t.Fatal("expected positional argument failure")
		}
		if !strings.Contains(stderr, "does not accept positional arguments") {
			t.Fatalf("unexpected stderr: %q", stderr)
		}
	})

	t.Run("duplicate future tag fails with actionable error", func(t *testing.T) {
		project := t.TempDir()
		initGitRepo(t, project)
		writeFile(t, filepath.Join(project, "README.md"), "release tests\n")
		commitAll(t, project, "initial")
		if _, stderr, code := runCommand(t, project, "git", "tag", "v0.4.2"); code != 0 {
			t.Fatalf("tag v0.4.2 failed: %s", stderr)
		}

		if _, stderr, code := runCommand(t, project, "git", "switch", "-c", "side-release"); code != 0 {
			t.Fatalf("switch side-release failed: %s", stderr)
		}
		writeFile(t, filepath.Join(project, "side.txt"), "side release\n")
		commitAll(t, project, "side release")
		if _, stderr, code := runCommand(t, project, "git", "tag", "v0.4.3"); code != 0 {
			t.Fatalf("tag v0.4.3 failed: %s", stderr)
		}

		if _, stderr, code := runCommand(t, project, "git", "switch", "main"); code != 0 {
			t.Fatalf("switch main failed: %s", stderr)
		}
		writeFile(t, filepath.Join(project, "main.txt"), "main release\n")
		commitAll(t, project, "main release")

		stdout, stderr, code := runGoMain(t, moduleRoot, scriptPath, "--dir", project)
		if code == 0 {
			t.Fatal("expected duplicate tag failure")
		}
		if strings.TrimSpace(stdout) != "" {
			t.Fatalf("unexpected stdout: %q", stdout)
		}
		if !strings.Contains(stderr, "next release tag v0.4.3 already exists") {
			t.Fatalf("unexpected stderr: %q", stderr)
		}
	})
}
