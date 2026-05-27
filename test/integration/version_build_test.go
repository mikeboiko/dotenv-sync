package integration_test

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"dotenv-sync/internal/release"
)

func TestVersionedBuildIntegration(t *testing.T) {
	scriptPath := "./scripts/nextversion"
	moduleRoot := repoRoot(t)

	t.Run("local preview matches binary metadata", func(t *testing.T) {
		project := t.TempDir()
		initGitRepo(t, project)
		writeFile(t, filepath.Join(project, "README.md"), "release tests\n")
		commitAll(t, project, "initial")
		if _, stderr, code := runCommand(t, project, "git", "tag", "v0.4.2"); code != 0 {
			t.Fatalf("tag v0.4.2 failed: %s", stderr)
		}
		writeFile(t, filepath.Join(project, "CHANGELOG.md"), "next release\n")
		commitAll(t, project, "next release")

		stdout, stderr, code := runGoMain(t, moduleRoot, scriptPath, "--dir", project)
		if code != 0 || stderr != "" {
			t.Fatalf("preview next version failed: code=%d stderr=%q", code, stderr)
		}
		version := strings.TrimSpace(stdout)
		if version != "v0.4.3" {
			t.Fatalf("preview version = %q", version)
		}

		bin := buildCLIWithLdflags(t, version, "feedbee", "2026-03-30T13:00:00Z")
		stdout, stderr, code = runCLI(t, bin, t.TempDir(), nil, "--version")
		if code != 0 || stderr != "" {
			t.Fatalf("short version failed: code=%d stderr=%q", code, stderr)
		}
		if strings.TrimSpace(stdout) != "ds "+version {
			t.Fatalf("short version = %q", stdout)
		}

		stdout, stderr, code = runCLI(t, bin, t.TempDir(), nil, "version")
		if code != 0 || stderr != "" {
			t.Fatalf("detailed version failed: code=%d stderr=%q", code, stderr)
		}
		wantDetailed := "Version: " + version + "\nCommit: feedbee\nBuilt: 2026-03-30T13:00:00Z\nPlatform: " + currentPlatform() + "\n"
		if stdout != wantDetailed {
			t.Fatalf("detailed version = %q want %q", stdout, wantDetailed)
		}
	})

	t.Run("asset names keep version parity across release targets", func(t *testing.T) {
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
	})
}
