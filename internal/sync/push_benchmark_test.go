package sync

import (
	"context"
	"path/filepath"
	"testing"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/envfile"
	"dotenv-sync/internal/provider"
)

func BenchmarkPlanPushDocs(b *testing.B) {
	cfg := config.Config{
		ItemName:    "bench-repo",
		StorageMode: config.StorageModeNoteJSON,
		EnvFile:     filepath.Join(b.TempDir(), ".env"),
	}
	schema := envfile.ParseBytes(".env.example", envfile.KindSchema, benchmarkSchema(500))
	local := envfile.ParseBytes(".env", envfile.KindLocal, benchmarkLocal(500))
	envValues := CanonicalDocumentEnv(local)
	prov := &fakeProvider{
		payload: provider.EnvPayload{
			ItemName: "bench-repo",
			Exists:   true,
			Env:      envValues,
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := PlanPushDocs(context.Background(), cfg, schema, local, prov); err != nil {
			b.Fatal(err)
		}
	}
}
