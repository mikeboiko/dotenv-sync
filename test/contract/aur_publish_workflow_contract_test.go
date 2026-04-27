package contract_test

import (
	"strings"
	"testing"
)

func TestContractAurPublishWorkflow(t *testing.T) {
	content := readRepoFile(t, ".github", "workflows", "aur-publish.yml")

	for _, want := range []string{
		"release:",
		"published",
		"ref: ${{ github.event.release.tag_name }}",
		"AUR_SSH_PRIVATE_KEY",
		"gh release download",
		"--license-sha256",
		"--readme-sha256",
		"go run ./scripts/aurpkg",
		"ssh://aur@aur.archlinux.org/dotenv-sync-bin.git",
		"upgpkg: dotenv-sync-bin",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("workflow missing %q", want)
		}
	}

	for _, unwanted := range []string{
		"workflow_dispatch:",
		"main",
	} {
		if strings.Contains(content, unwanted) {
			t.Fatalf("workflow unexpectedly contains %q", unwanted)
		}
	}
}
