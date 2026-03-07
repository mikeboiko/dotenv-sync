package sync

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
	"dotenv-sync/internal/report"
)

func TestPlanDiffReportsValidationExitOnDrift(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{SchemaFile: filepath.Join(dir, ".env.example"), EnvFile: filepath.Join(dir, ".env"), Mapping: map[string]string{"DATABASE_URL": "database-url"}}
	os.WriteFile(cfg.SchemaFile, []byte("DATABASE_URL=\nPORT=8080\n"), 0o644)
	os.WriteFile(cfg.EnvFile, []byte("DATABASE_URL=postgres://vault/dev\nPORT=9090\n"), 0o644)
	prov := fakeProvider{resolutions: map[string]provider.Resolution{"DATABASE_URL": {Source: "provider", Value: "postgres://vault/dev"}}}
	plan, err := PlanDiff(context.Background(), cfg, prov)
	if report.ExitCode(err) != report.ExitValidation {
		t.Fatalf("expected validation exit, got %v", err)
	}
	if Summarize(plan.Changes).Updated == 0 {
		t.Fatalf("expected update change in diff plan")
	}
}
