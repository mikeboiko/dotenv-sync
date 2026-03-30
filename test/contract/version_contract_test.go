package contract_test

import (
	"strings"
	"testing"
)

func TestContractVersionCommands(t *testing.T) {
	project := t.TempDir()

	t.Run("short version uses concise output", func(t *testing.T) {
		bin := buildCLIWithLdflags(t, "v1.2.3", "abc1234", "2026-03-30T12:00:00Z")
		stdout, stderr, code := runCLI(t, bin, project, nil, "--version")
		if code != 0 || stderr != "" {
			t.Fatalf("short version failed: code=%d stderr=%q", code, stderr)
		}
		want := renderTemplate(readRepoFile(t, "test", "testdata", "golden", "version-short.txt"), map[string]string{
			"{{VERSION}}": "v1.2.3",
		})
		if stdout != want {
			t.Fatalf("short version stdout = %q want %q", stdout, want)
		}
	})

	t.Run("detailed version uses shared metadata", func(t *testing.T) {
		bin := buildCLIWithLdflags(t, "v1.2.3", "abc1234", "2026-03-30T12:00:00Z")
		stdout, stderr, code := runCLI(t, bin, project, nil, "version")
		if code != 0 || stderr != "" {
			t.Fatalf("version command failed: code=%d stderr=%q", code, stderr)
		}
		want := renderTemplate(readRepoFile(t, "test", "testdata", "golden", "version-detailed.txt"), map[string]string{
			"{{VERSION}}":    "v1.2.3",
			"{{COMMIT}}":     "abc1234",
			"{{BUILD_TIME}}": "2026-03-30T12:00:00Z",
			"{{PLATFORM}}":   currentPlatform(),
		})
		if stdout != want {
			t.Fatalf("version stdout = %q want %q", stdout, want)
		}
	})

	t.Run("version rejects extra args", func(t *testing.T) {
		bin := buildCLI(t)
		_, stderr, code := runCLI(t, bin, project, nil, "version", "extra")
		if code == 0 {
			t.Fatal("expected non-zero exit for extra args")
		}
		if !strings.Contains(stderr, "unknown command \"extra\" for \"ds version\"") {
			t.Fatalf("unexpected stderr: %s", stderr)
		}
	})
}
