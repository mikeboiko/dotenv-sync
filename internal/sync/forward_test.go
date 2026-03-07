package sync

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/envfile"
	"dotenv-sync/internal/provider"
)

func TestPlanForwardDocsClassifiesChanges(t *testing.T) {
	cfg := config.Config{EnvFile: filepath.Join(t.TempDir(), ".env"), Mapping: map[string]string{"DATABASE_URL": "database-url", "JWT_SECRET": "jwt-secret"}}
	schema := envfile.ParseBytes(".env.example", envfile.KindSchema, []byte("DATABASE_URL=\nJWT_SECRET=\nPORT=8080\n"))
	local := envfile.ParseBytes(".env", envfile.KindLocal, []byte("DATABASE_URL=postgres://vault/dev\nPORT=9090\n"))
	prov := fakeProvider{resolutions: map[string]provider.Resolution{
		"DATABASE_URL": {Source: "provider", Value: "postgres://vault/dev"},
		"JWT_SECRET":   {Source: "provider", Value: "topsecret"},
	}}

	plan, target, err := PlanForwardDocs(context.Background(), cfg, schema, local, prov)
	if err != nil {
		t.Fatalf("plan forward docs returned error: %v", err)
	}
	if !plan.WriteRequired {
		t.Fatalf("expected write required")
	}
	got := string(envfile.Render(target))
	for _, want := range []string{"DATABASE_URL=postgres://vault/dev", "JWT_SECRET=topsecret", "PORT=8080"} {
		if !strings.Contains(got, want) {
			t.Fatalf("target output missing %q\n%s", want, got)
		}
	}

	summary := Summarize(plan.Changes)
	if summary.Added != 1 || summary.Updated != 1 || summary.Unchanged != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
}
