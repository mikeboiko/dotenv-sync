package bitwarden

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
)

func TestAdapterResolveUsesRepoItemAndFieldOverride(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "rbw.log")
	bin := filepath.Join(dir, "rbw")
	script := "#!/bin/sh\n" +
		"echo \"$@\" >> '" + logFile + "'\n" +
		"if [ \"$1\" = \"get\" ] && [ \"$2\" = \"--field\" ] && [ \"$3\" = \"db_url\" ] && [ \"$4\" = \"my-repo\" ]; then\n" +
		"  printf 'postgres://vault/dev\\n'\n" +
		"  exit 0\n" +
		"fi\n" +
		"echo 'not found' >&2\n" +
		"exit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	adapter := &Adapter{
		client: &RBWClient{Bin: bin},
		cfg:    config.Config{ItemName: "my-repo"},
		cache:  map[string]provider.Resolution{},
	}

	resolution, err := adapter.Resolve(context.Background(), "DATABASE_URL", "db_url")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if resolution.Source != "provider" || resolution.Value != "postgres://vault/dev" {
		t.Fatalf("unexpected resolution: %+v", resolution)
	}

	if _, err := adapter.Resolve(context.Background(), "DATABASE_URL", "db_url"); err != nil {
		t.Fatalf("resolve from cache: %v", err)
	}

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Fields(strings.TrimSpace(string(data)))
	if got := strings.Count(string(data), "get --field db_url my-repo"); got != 1 {
		t.Fatalf("expected one rbw invocation, log=%q parsed=%v", string(data), lines)
	}
}
