package sync

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
	"dotenv-sync/internal/report"
)

type failingProvider struct {
	err error
}

func (f failingProvider) Name() string { return "failing" }

func (f failingProvider) CheckReadiness(context.Context) (provider.Status, error) {
	return provider.Status{Provider: "bitwarden", CLIInstalled: true, Authenticated: true, Unlocked: true}, nil
}

func (f failingProvider) Resolve(context.Context, string, string) (provider.Resolution, error) {
	return provider.Resolution{}, f.err
}

func (f failingProvider) ResolveMany(context.Context, map[string]string) (map[string]provider.Resolution, error) {
	return nil, f.err
}

func TestPlanValidatePassesOnMatchingFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{SchemaFile: filepath.Join(dir, ".env.example"), EnvFile: filepath.Join(dir, ".env"), Mapping: map[string]string{"DATABASE_URL": "database-url"}}
	os.WriteFile(cfg.SchemaFile, []byte("DATABASE_URL=\nPORT=8080\n"), 0o644)
	os.WriteFile(cfg.EnvFile, []byte("DATABASE_URL=postgres://vault/dev\nPORT=8080\n"), 0o644)
	prov := fakeProvider{resolutions: map[string]provider.Resolution{"DATABASE_URL": {Source: "provider", Value: "postgres://vault/dev"}}}
	_, err := PlanValidate(context.Background(), cfg, prov)
	if err != nil {
		t.Fatalf("expected validate success, got %v", err)
	}
}

func TestPlanValidateReturnsExitTwoOnMissingSecret(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{SchemaFile: filepath.Join(dir, ".env.example"), EnvFile: filepath.Join(dir, ".env")}
	os.WriteFile(cfg.SchemaFile, []byte("JWT_SECRET=\n"), 0o644)
	os.WriteFile(cfg.EnvFile, []byte("JWT_SECRET=\n"), 0o644)
	_, err := PlanValidate(context.Background(), cfg, fakeProvider{})
	if report.ExitCode(err) != report.ExitValidation {
		t.Fatalf("expected validation exit, got %v", err)
	}
}

func TestPlanValidatePreservesMalformedProviderPayloadError(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{SchemaFile: filepath.Join(dir, ".env.example"), EnvFile: filepath.Join(dir, ".env")}
	os.WriteFile(cfg.SchemaFile, []byte("JWT_SECRET=\n"), 0o644)
	os.WriteFile(cfg.EnvFile, []byte("JWT_SECRET=value\n"), 0o644)

	_, err := PlanValidate(context.Background(), cfg, failingProvider{
		err: report.NewAppError("E010", report.ExitOperational, "provider note_json payload is malformed", "sync cannot trust the repo-scoped payload", "repair the Bitwarden notes and retry", nil),
	})
	var appErr *report.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != "E010" {
		t.Fatalf("expected E010, got %+v", appErr)
	}
}
