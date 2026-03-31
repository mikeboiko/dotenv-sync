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

func TestAdapterLoadEnvPayloadSupportsFieldsMode(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "rbw")
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"get\" ] && [ \"$2\" = \"--raw\" ] && [ \"$3\" = \"shared-dev\" ]; then\n" +
		"  printf '%s\\n' '{\"name\":\"shared-dev\",\"notes\":\"keep-me\",\"data\":{\"password\":\"shared-secret\"}}'\n" +
		"  exit 0\n" +
		"fi\n" +
		"echo 'not found' >&2\n" +
		"exit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	adapter := &Adapter{
		client: &RBWClient{Bin: bin},
		cfg:    config.Config{ItemName: "shared-dev", StorageMode: config.StorageModeFields},
		cache:  map[string]provider.Resolution{},
	}

	payload, err := adapter.LoadEnvPayload(context.Background())
	if err != nil {
		t.Fatalf("load fields payload: %v", err)
	}
	if !payload.Exists || payload.Password != "shared-secret" || payload.Notes != "keep-me" {
		t.Fatalf("unexpected fields payload: %+v", payload)
	}
}

func TestAdapterStoreEnvPayloadSupportsFieldsModePasswordWrites(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "rbw")
	logPath := filepath.Join(dir, "rbw.log")
	editCapture := filepath.Join(dir, "edit.txt")
	script := "#!/bin/sh\n" +
		"echo \"$@\" >> '" + logPath + "'\n" +
		"tmp=$(mktemp)\n" +
		"case \"$1\" in\n" +
		"edit)\n" +
		"  \"${VISUAL:-$EDITOR}\" \"$tmp\"\n" +
		"  cat \"$tmp\" > '" + editCapture + "'\n" +
		"  ;;\n" +
		"sync) exit 0 ;;\n" +
		"*) echo 'unsupported' >&2; exit 1 ;;\n" +
		"esac\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	ptyWrapper := filepath.Join(dir, "script")
	ptyWrapperContent := "#!/bin/sh\n" +
		"exec /bin/sh -c \"$3\"\n"
	if err := os.WriteFile(ptyWrapper, []byte(ptyWrapperContent), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	adapter := &Adapter{
		client: &RBWClient{Bin: bin},
		cfg: config.Config{
			ItemName:    "shared-dev",
			StorageMode: config.StorageModeFields,
			Mapping: map[string]string{
				"DB_PASSWD": "password",
			},
		},
		cache: map[string]provider.Resolution{},
	}

	_, err := adapter.StoreEnvPayload(context.Background(), provider.EnvPayload{
		ItemName:    "shared-dev",
		StorageMode: config.StorageModeFields,
		Exists:      true,
		Notes:       "keep-me",
		Password:    "old-secret",
		Env: map[string]string{
			"DB_PASSWD": "rotated-secret",
		},
	})
	if err != nil {
		t.Fatalf("store fields payload: %v", err)
	}
	editData, err := os.ReadFile(editCapture)
	if err != nil {
		t.Fatal(err)
	}
	if string(editData) != "rotated-secret\n\nkeep-me\n" {
		t.Fatalf("unexpected editor content: %q", string(editData))
	}
	logData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(logData), "edit shared-dev") || !strings.Contains(string(logData), "sync") {
		t.Fatalf("unexpected rbw log: %s", logData)
	}
}

func TestAdapterResolveTreatsNoEntryFoundAsMissing(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "rbw")
	script := "#!/bin/sh\n" +
		"echo \"rbw get: couldn't find entry for '$4': no entry found\" >&2\n" +
		"exit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	adapter := &Adapter{
		client: &RBWClient{Bin: bin},
		cfg:    config.Config{ItemName: "my-repo"},
		cache:  map[string]provider.Resolution{},
	}

	resolution, err := adapter.Resolve(context.Background(), "DATABASE_URL", "DATABASE_URL")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if resolution.Source != "missing" || resolution.IssueCode != "E005" {
		t.Fatalf("expected missing resolution, got %+v", resolution)
	}
}
