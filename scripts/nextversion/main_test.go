package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDefaultsToPatchPreview(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	writeFile(t, filepath.Join(repo, "README.md"), "release tests\n")
	commitAll(t, repo, "initial")
	tagRepo(t, repo, "v0.4.2")
	writeFile(t, filepath.Join(repo, "CHANGELOG.md"), "next change\n")
	commitAll(t, repo, "next")

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--dir", repo}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run exit code = %d stderr=%q", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != strings.TrimSpace(readGoldenFile(t, "release-published.txt", map[string]string{
		"{{VERSION}}": "v0.4.3",
	})) {
		t.Fatalf("preview stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func TestRunReturnsAlreadyReleasedExitCode(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	writeFile(t, filepath.Join(repo, "README.md"), "release tests\n")
	commitAll(t, repo, "initial")
	tagRepo(t, repo, "v0.4.3")

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--dir", repo}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("run exit code = %d want 2", code)
	}
	if strings.TrimSpace(stdout.String()) != "v0.4.3" {
		t.Fatalf("released stdout = %q", stdout.String())
	}
	wantStderr := strings.TrimSpace(readGoldenFile(t, "release-skipped.txt", map[string]string{
		"{{VERSION}}": "v0.4.3",
	}))
	if strings.TrimSpace(stderr.String()) != wantStderr {
		t.Fatalf("released stderr = %q want %q", stderr.String(), wantStderr)
	}
}

func TestRunRejectsPositionalArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"unexpected"}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("expected positional arg failure")
	}
	if !strings.Contains(stderr.String(), "does not accept positional arguments") {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func TestRunFailsWhenNextTagAlreadyExistsOutsideMainHistory(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	writeFile(t, filepath.Join(repo, "README.md"), "release tests\n")
	commitAll(t, repo, "initial")
	tagRepo(t, repo, "v0.4.2")

	runGit(t, repo, "switch", "-c", "side-release")
	writeFile(t, filepath.Join(repo, "side.txt"), "side release\n")
	commitAll(t, repo, "side release")
	tagRepo(t, repo, "v0.4.3")

	runGit(t, repo, "switch", "main")
	writeFile(t, filepath.Join(repo, "main.txt"), "main release\n")
	commitAll(t, repo, "main release")

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"--dir", repo}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("expected duplicate tag failure")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "next release tag v0.4.3 already exists") {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func readGoldenFile(t *testing.T, name string, replacements map[string]string) string {
	t.Helper()
	path := filepath.Join("..", "..", "test", "testdata", "golden", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for old, newValue := range replacements {
		content = strings.ReplaceAll(content, old, newValue)
	}
	return content
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	commands := [][]string{
		{"init", "-b", "main"},
		{"config", "user.name", "Test User"},
		{"config", "user.email", "test@example.com"},
	}
	for _, args := range commands {
		runGit(t, dir, args...)
	}
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

func commitAll(t *testing.T, dir, message string) {
	t.Helper()
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", message)
}

func tagRepo(t *testing.T, dir, tag string) {
	t.Helper()
	runGit(t, dir, "tag", tag)
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
