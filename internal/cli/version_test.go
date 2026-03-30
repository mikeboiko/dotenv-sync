package cli

import (
	"bytes"
	"testing"

	"dotenv-sync/pkg/dotenvsync"
)

func TestVersionCommandWritesDetailedMetadata(t *testing.T) {
	restore := snapshotVersionMetadata()
	defer restore()
	dotenvsync.Version = "v1.2.3"
	dotenvsync.Commit = "abc1234"
	dotenvsync.BuildTime = "2026-03-30T12:00:00Z"

	var stdout, stderr bytes.Buffer
	cmd := newVersionCommand(streams{stdout: &stdout, stderr: &stderr})
	cmd.SetArgs(nil)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute version command: %v", err)
	}
	if got := stdout.String(); got == "" {
		t.Fatal("expected version output")
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func TestVersionCommandRejectsExtraArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	cmd := newVersionCommand(streams{stdout: &stdout, stderr: &stderr})
	cmd.SetArgs([]string{"extra"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected extra args error")
	}
}

func snapshotVersionMetadata() func() {
	version, commit, buildTime := dotenvsync.Version, dotenvsync.Commit, dotenvsync.BuildTime
	return func() {
		dotenvsync.Version = version
		dotenvsync.Commit = commit
		dotenvsync.BuildTime = buildTime
	}
}
