package integration_test

import (
	"strings"
	"testing"

	"dotenv-sync/internal/release"
)

func TestVersionedBuildIntegration(t *testing.T) {
	bin := buildCLIWithLdflags(t, "v3.4.5", "feedbee", "2026-03-30T13:00:00Z")
	project := t.TempDir()

	stdout, stderr, code := runCLI(t, bin, project, nil, "--version")
	if code != 0 || stderr != "" {
		t.Fatalf("short version failed: code=%d stderr=%q", code, stderr)
	}
	if strings.TrimSpace(stdout) != "ds v3.4.5" {
		t.Fatalf("short version = %q", stdout)
	}

	assetName := release.AssetName("v3.4.5", "linux", "amd64")
	if !strings.Contains(assetName, "v3.4.5") {
		t.Fatalf("asset name missing version: %q", assetName)
	}
}
