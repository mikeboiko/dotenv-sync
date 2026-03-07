package sync

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"dotenv-sync/internal/config"
)

func BenchmarkPlanDiff(b *testing.B) {
	dir := b.TempDir()
	cfg := config.Config{SchemaFile: filepath.Join(dir, ".env.example"), EnvFile: filepath.Join(dir, ".env")}
	if err := os.WriteFile(cfg.SchemaFile, benchmarkSchema(500), 0o644); err != nil {
		b.Fatal(err)
	}
	if err := os.WriteFile(cfg.EnvFile, benchmarkLocal(500), 0o644); err != nil {
		b.Fatal(err)
	}
	prov := benchmarkProvider(500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := PlanDiff(context.Background(), cfg, prov)
		if err != nil {
			b.Fatal(err)
		}
	}
}
