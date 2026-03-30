package integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRootFromTB(tb testing.TB) string {
	tb.Helper()
	wd, err := os.Getwd()
	if err != nil {
		tb.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, "../.."))
}

func repoRoot(t *testing.T) string {
	return repoRootFromTB(t)
}

func buildCLI(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "ds")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/ds")
	cmd.Dir = repoRootFromTB(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build cli: %v\n%s", err, out)
	}
	return bin
}

func buildCLIWithLdflags(t *testing.T, version, commit, buildTime string) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "ds")
	ldflags := "-X dotenv-sync/pkg/dotenvsync.Version=" + version +
		" -X dotenv-sync/pkg/dotenvsync.Commit=" + commit +
		" -X dotenv-sync/pkg/dotenvsync.BuildTime=" + buildTime
	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", bin, "./cmd/ds")
	cmd.Dir = repoRootFromTB(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build cli with ldflags: %v\n%s", err, out)
	}
	return bin
}

func runCLI(t *testing.T, bin, dir string, extraEnv []string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), extraEnv...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		return stdout.String(), stderr.String(), 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return stdout.String(), stderr.String(), exitErr.ExitCode()
	}
	t.Fatalf("run cli: %v", err)
	return "", "", 0
}

func runGoMain(t *testing.T, dir, pkg string, args ...string) (string, string, int) {
	t.Helper()
	cmdArgs := append([]string{"run", pkg}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		return stdout.String(), stderr.String(), 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return stdout.String(), stderr.String(), exitErr.ExitCode()
	}
	t.Fatalf("run go main: %v", err)
	return "", "", 0
}

func runCommand(t *testing.T, dir, name string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		return stdout.String(), stderr.String(), 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return stdout.String(), stderr.String(), exitErr.ExitCode()
	}
	t.Fatalf("run command: %v", err)
	return "", "", 0
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	commands := [][]string{
		{"init", "-b", "main"},
		{"config", "user.name", "Test User"},
		{"config", "user.email", "test@example.com"},
	}
	for _, args := range commands {
		if stdout, stderr, code := runCommand(t, dir, "git", args...); code != 0 {
			t.Fatalf("git %v failed: code=%d stdout=%q stderr=%q", args, code, stdout, stderr)
		}
	}
}

func commitAll(t *testing.T, dir, message string) string {
	t.Helper()
	if stdout, stderr, code := runCommand(t, dir, "git", "add", "."); code != 0 {
		t.Fatalf("git add failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if stdout, stderr, code := runCommand(t, dir, "git", "commit", "-m", message); code != 0 {
		t.Fatalf("git commit failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	stdout, stderr, code := runCommand(t, dir, "git", "rev-parse", "--short", "HEAD")
	if code != 0 {
		t.Fatalf("git rev-parse failed: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	return strings.TrimSpace(stdout)
}

func currentPlatform() string {
	return runtime.GOOS + "/" + runtime.GOARCH
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
