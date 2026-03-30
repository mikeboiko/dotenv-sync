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

func TestAdapterResolveManyUsesNoteJSONPayloadOnce(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "rbw.log")
	bin := filepath.Join(dir, "rbw")
	script := "#!/bin/sh\n" +
		"echo \"$@\" >> '" + logFile + "'\n" +
		"if [ \"$1\" = \"get\" ] && [ \"$2\" = \"--raw\" ] && [ \"$3\" = \"my-repo\" ]; then\n" +
		"  cat <<'EOF'\n" +
		"{\"name\":\"my-repo\",\"notes\":\"{\\\"format\\\":\\\"dotenv-sync/note-json@v1\\\",\\\"env\\\":{\\\"DATABASE_URL\\\":\\\"postgres://vault/dev\\\",\\\"JWT_SECRET\\\":\\\"supersecret\\\"}}\",\"data\":{\"password\":\"keep-me\"}}\n" +
		"EOF\n" +
		"  exit 0\n" +
		"fi\n" +
		"echo 'not found' >&2\n" +
		"exit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	adapter := &Adapter{
		client: &RBWClient{Bin: bin},
		cfg:    config.Config{ItemName: "my-repo", StorageMode: config.StorageModeNoteJSON},
		cache:  map[string]provider.Resolution{},
	}

	results, err := adapter.ResolveMany(context.Background(), map[string]string{
		"DATABASE_URL": "ignored-mapping",
		"JWT_SECRET":   "also-ignored",
		"MISSING_KEY":  "still-ignored",
	})
	if err != nil {
		t.Fatalf("resolve many: %v", err)
	}
	if results["DATABASE_URL"].Value != "postgres://vault/dev" || results["JWT_SECRET"].Value != "supersecret" {
		t.Fatalf("unexpected results: %+v", results)
	}
	if results["MISSING_KEY"].IssueCode != "E005" {
		t.Fatalf("expected missing note_json key to use E005, got %+v", results["MISSING_KEY"])
	}
	if _, err := adapter.Resolve(context.Background(), "DATABASE_URL", "ignored"); err != nil {
		t.Fatalf("resolve from cached payload: %v", err)
	}

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(data), "get --raw my-repo"); got != 1 {
		t.Fatalf("expected one raw rbw invocation, got %d log=%q", got, string(data))
	}
}

func TestAdapterLoadEnvPayloadRejectsMalformedNotes(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "rbw")
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"get\" ] && [ \"$2\" = \"--raw\" ] && [ \"$3\" = \"my-repo\" ]; then\n" +
		"  printf '%s\\n' '{\"name\":\"my-repo\",\"notes\":\"not-json\",\"data\":{\"password\":\"keep-me\"}}'\n" +
		"  exit 0\n" +
		"fi\n" +
		"echo 'not found' >&2\n" +
		"exit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	adapter := &Adapter{
		client: &RBWClient{Bin: bin},
		cfg:    config.Config{ItemName: "my-repo", StorageMode: config.StorageModeNoteJSON},
		cache:  map[string]provider.Resolution{},
	}

	if _, err := adapter.LoadEnvPayload(context.Background()); err == nil || !strings.Contains(err.Error(), "provider note_json payload is malformed") {
		t.Fatalf("expected malformed payload error, got %v", err)
	}
}
