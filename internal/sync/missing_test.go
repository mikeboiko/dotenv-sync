package sync

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/report"
)

func TestPlanMissingReturnsExitTwoWhenKeysAreUnresolved(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{SchemaFile: filepath.Join(dir, ".env.example"), EnvFile: filepath.Join(dir, ".env")}
	os.WriteFile(cfg.SchemaFile, []byte("JWT_SECRET=\n"), 0o644)
	os.WriteFile(cfg.EnvFile, []byte("JWT_SECRET=\n"), 0o644)
	plan, err := PlanMissing(context.Background(), cfg, fakeProvider{})
	if report.ExitCode(err) != report.ExitValidation {
		t.Fatalf("expected validation exit, got %v", err)
	}
	if len(plan.Changes) == 0 || plan.Changes[0].ChangeType != "missing" {
		t.Fatalf("expected missing change, got %#v", plan.Changes)
	}
}
