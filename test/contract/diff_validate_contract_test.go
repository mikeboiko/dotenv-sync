package contract_test

import (
	"strings"
	"testing"
)

func TestContractDiffValidateAndMissing(t *testing.T) {
	bin := buildCLI(t)
	project := setupProject(t,
		"DATABASE_URL=\nJWT_SECRET=\nPORT=8080\n",
		"DATABASE_URL=postgres://vault/dev\nPORT=9090\nEXTRA_KEY=value\n",
		"item_name: shared-dev\nmapping:\n  DATABASE_URL: database-url\n  JWT_SECRET: jwt-secret\n",
	)
	_, env := writeRBWStub(t, "unlocked", map[string]string{
		rbwLookupKey("shared-dev", "database-url"): "postgres://vault/dev",
	}, rbwLookupKey("shared-dev", "jwt-secret"))

	stdout, _, code := runCLI(t, bin, project, env, "diff")
	if code != 2 {
		t.Fatalf("diff exit code=%d stdout=%s", code, stdout)
	}
	for _, want := range []string{"MISSING JWT_SECRET [MISSING]", "UPDATE PORT [STATIC]", "EXTRA EXTRA_KEY"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("diff output missing %q\n%s", want, stdout)
		}
	}

	stdout, _, code = runCLI(t, bin, project, env, "validate")
	if code != 2 {
		t.Fatalf("validate exit code=%d stdout=%s", code, stdout)
	}
	for _, want := range []string{"E005 JWT_SECRET", "UPDATE PORT [STATIC]", "EXTRA EXTRA_KEY"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("validate output missing %q\n%s", want, stdout)
		}
	}

	stdout, _, code = runCLI(t, bin, project, env, "missing")
	if code != 2 {
		t.Fatalf("missing exit code=%d stdout=%s", code, stdout)
	}
	if !strings.Contains(stdout, "MISSING JWT_SECRET") {
		t.Fatalf("missing output missing unresolved key: %s", stdout)
	}
}
