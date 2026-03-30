package sync

import (
	"context"
	"path/filepath"
	"testing"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/envfile"
	"dotenv-sync/internal/provider"
)

func TestPlanPushDocsClassifiesChanges(t *testing.T) {
	cfg := config.Config{
		ItemName:    "my-repo",
		StorageMode: config.StorageModeNoteJSON,
		EnvFile:     filepath.Join(t.TempDir(), ".env"),
	}
	schema := envfile.ParseBytes(".env.example", envfile.KindSchema, []byte("DATABASE_URL=\nJWT_SECRET=\nPORT=8080\n"))
	local := envfile.ParseBytes(".env", envfile.KindLocal, []byte("DATABASE_URL=postgres://vault/dev\nJWT_SECRET=supersecret\nPORT=8080\nEXTRA_FLAG=true\n"))
	prov := &fakeProvider{
		payload: provider.EnvPayload{
			ItemName: "my-repo",
			Exists:   true,
			Env: map[string]string{
				"DATABASE_URL": "postgres://vault/old",
				"PORT":         "8080",
				"EXTRA_FLAG":   "false",
			},
		},
	}

	plan, target, err := PlanPushDocs(context.Background(), cfg, schema, local, prov)
	if err != nil {
		t.Fatalf("plan push docs: %v", err)
	}
	if !plan.WriteRequired {
		t.Fatal("expected write required")
	}
	if summary := Summarize(plan.Changes); summary.Added != 1 || summary.Updated != 1 || summary.Unchanged != 1 || summary.Extra != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	if target.ItemName != "my-repo" || !NoteJSONEqual(target.Env, CanonicalDocumentEnv(local)) {
		t.Fatalf("unexpected target payload: %+v", target)
	}
}

func TestPlanPushDocsRejectsFieldsMode(t *testing.T) {
	cfg := config.Config{StorageMode: config.StorageModeFields}
	schema := envfile.ParseBytes(".env.example", envfile.KindSchema, []byte("DATABASE_URL=\n"))
	local := envfile.ParseBytes(".env", envfile.KindLocal, []byte("DATABASE_URL=value\n"))

	if _, _, err := PlanPushDocs(context.Background(), cfg, schema, local, &fakeProvider{}); err == nil {
		t.Fatal("expected fields mode rejection")
	}
}
