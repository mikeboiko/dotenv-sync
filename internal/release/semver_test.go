package release

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLatestVersionFromTagsIgnoresNonSemver(t *testing.T) {
	got, ok := LatestVersionFromTags([]string{
		"release-2026-03-30",
		"v1.4.2",
		"notes",
		"v1.4.10",
		"v2.0.0-beta",
	})
	if !ok {
		t.Fatal("expected latest semver tag")
	}
	if got != "v1.4.10" {
		t.Fatalf("latest tag = %q", got)
	}
}

func TestNextVersionBumpsSemver(t *testing.T) {
	cases := map[string]string{
		"patch": "v1.4.11",
		"minor": "v1.5.0",
		"major": "v2.0.0",
	}
	for bump, want := range cases {
		got, err := NextVersion("v1.4.10", bump)
		if err != nil {
			t.Fatalf("next version %s: %v", bump, err)
		}
		if got != want {
			t.Fatalf("%s bump = %q", bump, got)
		}
	}
}

func TestReleaseAssetNames(t *testing.T) {
	got := []string{
		AssetName("v1.2.3", "linux", "amd64"),
		AssetName("v1.2.3", "linux", "arm64"),
		AssetName("v1.2.3", "darwin", "amd64"),
		AssetName("v1.2.3", "darwin", "arm64"),
		AssetName("v1.2.3", "windows", "amd64"),
	}
	want := []string{
		"ds_v1.2.3_linux_amd64.tar.gz",
		"ds_v1.2.3_linux_arm64.tar.gz",
		"ds_v1.2.3_darwin_amd64.tar.gz",
		"ds_v1.2.3_darwin_arm64.tar.gz",
		"ds_v1.2.3_windows_amd64.zip",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("asset names = %#v", got)
	}
}

func TestLatestVersionFromRepoUsesReachableSemverTags(t *testing.T) {
	repo := t.TempDir()
	initGitRepository(t, repo)
	writeRepoFile(t, repo, "README.md", "version tests\n")
	commitRepo(t, repo, "initial")
	tagRepo(t, repo, "v0.4.2")
	tagRepo(t, repo, "notes")
	writeRepoFile(t, repo, "CHANGELOG.md", "next change\n")
	commitRepo(t, repo, "next")
	tagRepo(t, repo, "v0.5.1")

	got, ok, err := LatestVersionFromRepo(context.Background(), repo)
	if err != nil {
		t.Fatalf("latest version from repo: %v", err)
	}
	if !ok {
		t.Fatal("expected latest version from repo")
	}
	if got != "v0.5.1" {
		t.Fatalf("latest version from repo = %q", got)
	}
}

func TestValidateReleaseBranch(t *testing.T) {
	if err := ValidateReleaseBranch("main", "main"); err != nil {
		t.Fatalf("unexpected branch validation error: %v", err)
	}
	if err := ValidateReleaseBranch("version", "main"); err == nil {
		t.Fatal("expected branch validation failure")
	}
}

func initGitRepository(t *testing.T, dir string) {
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

func writeRepoFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commitRepo(t *testing.T, dir, message string) {
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
