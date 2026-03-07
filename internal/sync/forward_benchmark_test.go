package sync

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/envfile"
	"dotenv-sync/internal/provider"
)

func benchmarkSchema(size int) []byte {
	var b strings.Builder
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			fmt.Fprintf(&b, "KEY_%03d=static\n", i)
		} else {
			fmt.Fprintf(&b, "KEY_%03d=\n", i)
		}
	}
	return []byte(b.String())
}

func benchmarkLocal(size int) []byte {
	var b strings.Builder
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			fmt.Fprintf(&b, "KEY_%03d=static\n", i)
		} else {
			fmt.Fprintf(&b, "KEY_%03d=resolved\n", i)
		}
	}
	return []byte(b.String())
}

func benchmarkProvider(size int) fakeProvider {
	resolutions := make(map[string]provider.Resolution, size)
	for i := 0; i < size; i++ {
		if i%2 == 1 {
			key := fmt.Sprintf("KEY_%03d", i)
			resolutions[key] = provider.Resolution{Source: "provider", Value: "resolved"}
		}
	}
	return fakeProvider{resolutions: resolutions}
}

func BenchmarkPlanForwardDocs(b *testing.B) {
	schema := envfile.ParseBytes(".env.example", envfile.KindSchema, benchmarkSchema(500))
	local := envfile.ParseBytes(".env", envfile.KindLocal, benchmarkLocal(500))
	cfg := config.Config{EnvFile: filepath.Join(b.TempDir(), ".env")}
	prov := benchmarkProvider(500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := PlanForwardDocs(context.Background(), cfg, schema, local, prov)
		if err != nil {
			b.Fatal(err)
		}
	}
}
