package contract_test

import (
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"dotenv-sync/internal/release"
)

func TestContractNextPatchPreviewFeedsVersionedBuild(t *testing.T) {
	project := t.TempDir()
	initVersionRepo(t, project)
	writeFile(t, filepath.Join(project, "README.md"), "release tests\n")
	commitVersionRepo(t, project, "initial")
	tagVersionRepo(t, project, "v0.4.2")
	writeFile(t, filepath.Join(project, "CHANGELOG.md"), "next release\n")
	commitVersionRepo(t, project, "next release")

	version := previewNextVersion(t, project)
	if version != "v0.4.3" {
		t.Fatalf("preview version = %q", version)
	}

	bin := buildCLIWithLdflags(t, version, "feedbee", "2026-03-30T13:00:00Z")

	stdout, stderr, code := runCLI(t, bin, t.TempDir(), nil, "--version")
	if code != 0 || stderr != "" {
		t.Fatalf("short version failed: code=%d stderr=%q", code, stderr)
	}
	wantShort := renderTemplate(readGoldenFile(t, "version-short.txt"), map[string]string{
		"{{VERSION}}": version,
	})
	if stdout != wantShort {
		t.Fatalf("short version stdout = %q want %q", stdout, wantShort)
	}

	stdout, stderr, code = runCLI(t, bin, t.TempDir(), nil, "version")
	if code != 0 || stderr != "" {
		t.Fatalf("version command failed: code=%d stderr=%q", code, stderr)
	}
	wantDetailed := renderTemplate(readGoldenFile(t, "version-detailed.txt"), map[string]string{
		"{{VERSION}}":    version,
		"{{COMMIT}}":     "feedbee",
		"{{BUILD_TIME}}": "2026-03-30T13:00:00Z",
		"{{PLATFORM}}":   currentPlatform(),
	})
	if stdout != wantDetailed {
		t.Fatalf("version stdout = %q want %q", stdout, wantDetailed)
	}
}

func TestContractReleaseArtifactsKeepVersionParityAcrossTargets(t *testing.T) {
	version := "v1.2.3"
	got := []string{
		release.AssetName(version, "linux", "amd64"),
		release.AssetName(version, "linux", "arm64"),
		release.AssetName(version, "darwin", "amd64"),
		release.AssetName(version, "darwin", "arm64"),
		release.AssetName(version, "windows", "amd64"),
	}
	want := readReleaseFixtureLines(t, "expected-assets.txt")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("asset names = %#v want %#v", got, want)
	}
	for _, asset := range got {
		if !strings.Contains(asset, version) {
			t.Fatalf("asset %q does not include version %q", asset, version)
		}
	}
}

func previewNextVersion(t *testing.T, repoDir string) string {
	t.Helper()
	cmd := exec.Command("go", "run", "./scripts/nextversion", "--dir", repoDir)
	cmd.Dir = repoRoot(t)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("preview next version: %v\n%s", err, out)
	}
	return strings.TrimSpace(string(out))
}

func initVersionRepo(t *testing.T, dir string) {
	t.Helper()
	for _, args := range [][]string{
		{"init", "-b", "main"},
		{"config", "user.name", "Test User"},
		{"config", "user.email", "test@example.com"},
	} {
		runVersionGit(t, dir, args...)
	}
}

func commitVersionRepo(t *testing.T, dir, message string) {
	t.Helper()
	runVersionGit(t, dir, "add", ".")
	runVersionGit(t, dir, "commit", "-m", message)
}

func tagVersionRepo(t *testing.T, dir, tag string) {
	t.Helper()
	runVersionGit(t, dir, "tag", tag)
}

func runVersionGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
