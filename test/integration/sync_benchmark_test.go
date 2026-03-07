package integration_test

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkLargeFixturePreparation(b *testing.B) {
	schema := loadLargeSchema(b)
	for i := 0; i < b.N; i++ {
		project := b.TempDir()
		if err := os.WriteFile(filepath.Join(project, ".env.example"), []byte(schema), 0o644); err != nil {
			b.Fatal(err)
		}
	}
}
