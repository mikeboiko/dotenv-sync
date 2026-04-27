package contract_test

import (
	"strings"
	"testing"
)

func TestContractReleaseWorkflow(t *testing.T) {
	content := readRepoFile(t, ".github", "workflows", "release.yml")

	for _, want := range []string{
		"push:",
		"tags:",
		"v*.*.*",
		"concurrency:",
		"cancel-in-progress: false",
		"timeout-minutes: 15",
		"go test ./...",
		"stable semver tag",
		"CGO_ENABLED=0",
		"cp LICENSE README.md",
		"ds_${VERSION}_SHA256SUMS",
		"gh release view",
		"gh release upload",
		"gh release create",
		"ds --version",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("workflow missing %q", want)
		}
	}

	for _, unwanted := range []string{
		"branches:",
		"- main",
		"workflow_dispatch:",
		"go run ./scripts/nextversion",
		"release_required",
		"skip_reason",
		"commit already released by tag",
		"repair the GitHub release record manually",
	} {
		if strings.Contains(content, unwanted) {
			t.Fatalf("workflow unexpectedly contains %q", unwanted)
		}
	}

	testIndex := strings.Index(content, "go test ./...")
	buildIndex := strings.Index(content, "Build release artifacts")
	verifyIndex := strings.Index(content, "ds --version")
	publishIndex := strings.Index(content, "gh release view")
	releaseIndex := strings.Index(content, "gh release create")
	if testIndex == -1 || buildIndex == -1 || verifyIndex == -1 || publishIndex == -1 || releaseIndex == -1 {
		t.Fatal("workflow missing required ordering markers")
	}
	if testIndex > buildIndex || buildIndex > verifyIndex || verifyIndex > publishIndex || publishIndex > releaseIndex {
		t.Fatalf("workflow must run tests before creating the release")
	}
}
