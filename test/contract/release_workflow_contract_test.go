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
		"concurrency:",
		"cancel-in-progress: false",
		"timeout-minutes: 15",
		"go test ./...",
		"go run ./scripts/nextversion",
		"release_required",
		"skip_reason",
		"gh release view",
		"commit already released by tag",
		"repair the GitHub release record manually",
		"gh release create",
		"ds --version",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("workflow missing %q", want)
		}
	}

	for _, unwanted := range []string{
		"workflow_dispatch:",
		"release_notes:",
		"--bump",
		"inputs:",
	} {
		if strings.Contains(content, unwanted) {
			t.Fatalf("workflow unexpectedly contains %q", unwanted)
		}
	}

	testIndex := strings.Index(content, "go test ./...")
	buildIndex := strings.Index(content, "Build release artifacts")
	verifyIndex := strings.Index(content, "ds --version")
	releaseIndex := strings.Index(content, "gh release create")
	if testIndex == -1 || buildIndex == -1 || verifyIndex == -1 || releaseIndex == -1 {
		t.Fatal("workflow missing required ordering markers")
	}
	if testIndex > buildIndex || buildIndex > verifyIndex || verifyIndex > releaseIndex {
		t.Fatalf("workflow must run tests before creating the release")
	}
	if strings.Index(content, "gh release view") > releaseIndex {
		t.Fatalf("workflow must handle rerun repair checks before creating the release")
	}
}
