package contract_test

import (
	"strings"
	"testing"
)

func TestContractReleaseWorkflow(t *testing.T) {
	content := readRepoFile(t, ".github", "workflows", "release.yml")

	for _, want := range []string{
		"push:",
		"branches:",
		"- main",
		"timeout-minutes: 15",
		"Wait for older release runs",
		"actions/workflows/release.yml/runs",
		"go test ./...",
		"go run ./scripts/nextversion --bump minor",
		"already_released",
		"Commit already released by tag",
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
		"workflow_dispatch:",
		"tags:",
		"v*.*.*",
		"concurrency:",
	} {
		if strings.Contains(content, unwanted) {
			t.Fatalf("workflow unexpectedly contains %q", unwanted)
		}
	}

	waitIndex := strings.Index(content, "Wait for older release runs")
	resolveIndex := strings.Index(content, "go run ./scripts/nextversion --bump minor")
	testIndex := strings.Index(content, "go test ./...")
	buildIndex := strings.Index(content, "Build release artifacts")
	verifyIndex := strings.Index(content, "ds --version")
	publishIndex := strings.Index(content, "gh release view")
	releaseIndex := strings.Index(content, "gh release create")
	if waitIndex == -1 || resolveIndex == -1 || testIndex == -1 || buildIndex == -1 || verifyIndex == -1 || publishIndex == -1 || releaseIndex == -1 {
		t.Fatal("workflow missing required ordering markers")
	}
	if waitIndex > resolveIndex || resolveIndex > testIndex || testIndex > buildIndex || buildIndex > verifyIndex || verifyIndex > publishIndex || publishIndex > releaseIndex {
		t.Fatalf("workflow must run tests before creating the release")
	}
}
