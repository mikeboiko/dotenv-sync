package contract_test

import "testing"

func TestContractVersionedLocalBuild(t *testing.T) {
	bin := buildCLIWithLdflags(t, "v3.4.5", "feedbee", "2026-03-30T13:00:00Z")
	project := t.TempDir()

	stdout, stderr, code := runCLI(t, bin, project, nil, "version")
	if code != 0 || stderr != "" {
		t.Fatalf("version command failed: code=%d stderr=%q", code, stderr)
	}

	want := renderTemplate(readRepoFile(t, "test", "testdata", "golden", "version-detailed.txt"), map[string]string{
		"{{VERSION}}":    "v3.4.5",
		"{{COMMIT}}":     "feedbee",
		"{{BUILD_TIME}}": "2026-03-30T13:00:00Z",
		"{{PLATFORM}}":   currentPlatform(),
	})
	if stdout != want {
		t.Fatalf("version stdout = %q want %q", stdout, want)
	}
}
