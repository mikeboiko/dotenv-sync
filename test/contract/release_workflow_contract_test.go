package contract_test

import (
	"strings"
	"testing"
)

func TestContractReleaseWorkflow(t *testing.T) {
	content := readRepoFile(t, ".github", "workflows", "release.yml")

	for _, want := range []string{
		"workflow_dispatch:",
		"bump:",
		"- major",
		"- minor",
		"- patch",
		"release_notes:",
		"go test ./...",
		"go run ./scripts/nextversion --bump",
		"gh release create",
		"github.event.repository.default_branch",
		"Ensure the release tag does not already exist",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("workflow missing %q", want)
		}
	}

	testIndex := strings.Index(content, "go test ./...")
	tagCheckIndex := strings.Index(content, "Ensure the release tag does not already exist")
	releaseIndex := strings.Index(content, "gh release create")
	if testIndex == -1 || releaseIndex == -1 || testIndex > releaseIndex {
		t.Fatalf("workflow must run tests before creating the release")
	}
	if tagCheckIndex == -1 || tagCheckIndex > releaseIndex {
		t.Fatalf("workflow must check for duplicate tags before creating the release")
	}
}
