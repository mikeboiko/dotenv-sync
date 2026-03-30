package release

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLatestVersionFromTagsIgnoresNonSemver(t *testing.T) {
	got, ok := LatestVersionFromTags(readReleaseFixtureLines(t, "mixed-tags.txt"))
	if !ok {
		t.Fatal("expected latest semver tag")
	}
	if got != "v1.4.10" {
		t.Fatalf("latest tag = %q", got)
	}
}

func TestNextPatchVersion(t *testing.T) {
	got, err := NextPatchVersion("v1.4.10")
	if err != nil {
		t.Fatalf("next patch version: %v", err)
	}
	if got != "v1.4.11" {
		t.Fatalf("next patch = %q", got)
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
	want := readReleaseFixtureLines(t, "expected-assets.txt")
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

func TestNextPatchVersionForRepoUsesBaselineWhenNoTags(t *testing.T) {
	repo := t.TempDir()
	initGitRepository(t, repo)
	writeRepoFile(t, repo, "README.md", "version tests\n")
	commitRepo(t, repo, "initial")

	got, err := NextPatchVersionForRepo(context.Background(), repo)
	if err != nil {
		t.Fatalf("next patch version from repo: %v", err)
	}
	if got != "v0.0.1" {
		t.Fatalf("baseline patch version = %q", got)
	}
}

func TestReleaseTagForRefDetectsAlreadyReleasedCommit(t *testing.T) {
	repo := t.TempDir()
	initGitRepository(t, repo)
	writeRepoFile(t, repo, "README.md", "version tests\n")
	commitRepo(t, repo, "initial")
	tagRepo(t, repo, readSingleReleaseFixture(t, "already-released.txt"))
	tagRepo(t, repo, "notes")

	got, ok, err := ReleaseTagForRef(context.Background(), repo, "HEAD")
	if err != nil {
		t.Fatalf("release tag for ref: %v", err)
	}
	if !ok {
		t.Fatal("expected semver release tag for HEAD")
	}
	if got != "v0.4.3" {
		t.Fatalf("release tag for HEAD = %q", got)
	}
}

func TestVersionTagExists(t *testing.T) {
	repo := t.TempDir()
	initGitRepository(t, repo)
	writeRepoFile(t, repo, "README.md", "version tests\n")
	commitRepo(t, repo, "initial")
	tagRepo(t, repo, "v0.4.3")

	exists, err := VersionTagExists(context.Background(), repo, "v0.4.3")
	if err != nil {
		t.Fatalf("version tag exists: %v", err)
	}
	if !exists {
		t.Fatal("expected version tag to exist")
	}

	exists, err = VersionTagExists(context.Background(), repo, "v0.4.4")
	if err != nil {
		t.Fatalf("version tag exists: %v", err)
	}
	if exists {
		t.Fatal("did not expect version tag to exist")
	}
}

func TestVersionTagExistsFindsUnmergedFutureTag(t *testing.T) {
	repo := t.TempDir()
	initGitRepository(t, repo)
	writeRepoFile(t, repo, "README.md", "version tests\n")
	commitRepo(t, repo, "initial")
	tagRepo(t, repo, "v0.4.2")

	runGit(t, repo, "switch", "-c", "side-release")
	writeRepoFile(t, repo, "side.txt", "side release\n")
	commitRepo(t, repo, "side release")
	tagRepo(t, repo, "v0.4.3")

	runGit(t, repo, "switch", "main")
	writeRepoFile(t, repo, "main.txt", "main release\n")
	commitRepo(t, repo, "main release")

	latest, ok, err := LatestVersionFromRepo(context.Background(), repo)
	if err != nil {
		t.Fatalf("latest version from repo: %v", err)
	}
	if !ok {
		t.Fatal("expected reachable semver tag")
	}
	if latest != "v0.4.2" {
		t.Fatalf("reachable version = %q", latest)
	}

	exists, err := VersionTagExists(context.Background(), repo, "v0.4.3")
	if err != nil {
		t.Fatalf("version tag exists: %v", err)
	}
	if !exists {
		t.Fatal("expected unmerged future tag to exist globally")
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

func readReleaseFixtureLines(t *testing.T, name string) []string {
	t.Helper()
	path := filepath.Join("..", "..", "test", "testdata", "release", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		filtered = append(filtered, line)
	}
	return filtered
}

func readSingleReleaseFixture(t *testing.T, name string) string {
	t.Helper()
	lines := readReleaseFixtureLines(t, name)
	if len(lines) != 1 {
		t.Fatalf("expected exactly one release fixture line in %s, got %d", name, len(lines))
	}
	return lines[0]
}
