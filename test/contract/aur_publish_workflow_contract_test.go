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
		"workflow_dispatch:",
		"release_tag:",
		"pkgrel:",
		"fetch-depth: 0",
		"AUR_SSH_PRIVATE_KEY",
		"gh release download",
		"git show \"${{ steps.meta.outputs.version }}:LICENSE\"",
		"git show \"${{ steps.meta.outputs.version }}:README.md\"",
		"--pkgrel",
		"--license-sha256",
		"--readme-sha256",
		"go run ./scripts/aurpkg",
		"ssh://aur@aur.archlinux.org/dotenv-sync-bin.git",
		"upgpkg: dotenv-sync-bin ${PKGVER}-${PKGREL}",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("workflow missing %q", want)
		}
	}

	for _, unwanted := range []string{
		"main",
	} {
		if strings.Contains(content, unwanted) {
			t.Fatalf("workflow unexpectedly contains %q", unwanted)
		}
	}
}
