package dotenvsync

import (
	"runtime"
	"testing"
)

func TestCurrentUsesFallbackMetadata(t *testing.T) {
	restore := snapshotMetadata()
	defer restore()
	Version = ""
	Commit = ""
	BuildTime = ""

	got := Current()
	if got.Version != "dev" {
		t.Fatalf("version fallback = %q", got.Version)
	}
	if got.Commit != "none" {
		t.Fatalf("commit fallback = %q", got.Commit)
	}
	if got.BuildTime != "unknown" {
		t.Fatalf("build time fallback = %q", got.BuildTime)
	}
	if got.Platform != runtime.GOOS+"/"+runtime.GOARCH {
		t.Fatalf("platform = %q", got.Platform)
	}
}

func TestRenderersUseInjectedMetadata(t *testing.T) {
	restore := snapshotMetadata()
	defer restore()
	Version = "v1.2.3"
	Commit = "abc1234"
	BuildTime = "2026-03-30T12:00:00Z"

	if got := Short("ds"); got != "ds v1.2.3" {
		t.Fatalf("short version = %q", got)
	}

	want := "Version: v1.2.3\nCommit: abc1234\nBuilt: 2026-03-30T12:00:00Z\nPlatform: " + runtime.GOOS + "/" + runtime.GOARCH
	if got := Detailed(); got != want {
		t.Fatalf("detailed version = %q", got)
	}
}

func snapshotMetadata() func() {
	version, commit, buildTime := Version, Commit, BuildTime
	return func() {
		Version = version
		Commit = commit
		BuildTime = buildTime
	}
}
