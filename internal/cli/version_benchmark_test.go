package cli

import (
	"io"
	"testing"

	"dotenv-sync/pkg/dotenvsync"
)

func BenchmarkVersionCommand(b *testing.B) {
	restore := snapshotBenchmarkMetadata()
	defer restore()
	dotenvsync.Version = "v1.2.3"
	dotenvsync.Commit = "abc1234"
	dotenvsync.BuildTime = "2026-03-30T12:00:00Z"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := newVersionCommand(streams{stdout: io.Discard, stderr: io.Discard})
		cmd.SetArgs(nil)
		if err := cmd.Execute(); err != nil {
			b.Fatal(err)
		}
	}
}

func snapshotBenchmarkMetadata() func() {
	version, commit, buildTime := dotenvsync.Version, dotenvsync.Commit, dotenvsync.BuildTime
	return func() {
		dotenvsync.Version = version
		dotenvsync.Commit = commit
		dotenvsync.BuildTime = buildTime
	}
}
