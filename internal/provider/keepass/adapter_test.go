package keepass

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
)

// adapterWithStub builds an Adapter wired to a stub binary with the given password already set.
// This bypasses the interactive password prompt for testing.
func adapterWithStub(t *testing.T, bin, dbPath, group, password string) *Adapter {
	t.Helper()
	cfg := config.Config{
		KeePassDatabase: dbPath,
		KeePassGroup:    group,
	}
	a := NewAdapter(cfg)
	a.client.Bin = bin
	a.client.Password = password // pre-set so ensurePassword is a no-op
	return a
}

func TestAdapterNameIsKeepass(t *testing.T) {
	a := NewAdapter(config.Config{})
	if a.Name() != "keepass" {
		t.Fatalf("expected name=keepass, got %q", a.Name())
	}
}

func TestAdapterResolveReturnsValue(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.kdbx")
	// Create a dummy db file so the path exists
	if err := os.WriteFile(dbPath, []byte("dummy"), 0o600); err != nil {
		t.Fatal(err)
	}
	bin := stubBin(t, dir, "keepassxc-cli", `
if [ "$1" = "show" ] && [ "$2" = "-s" ]; then
  printf "Title: DATABASE_URL\nUserName: \nPassword: postgres://vault/dev\nURL: \nNotes: \n"
  exit 0
fi
exit 1
`)
	a := adapterWithStub(t, bin, dbPath, "dotenv", "masterpassword")
	res, err := a.Resolve(context.Background(), "DATABASE_URL", "")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if res.Source != "provider" || res.Value != "postgres://vault/dev" {
		t.Fatalf("unexpected resolution: %+v", res)
	}
}

func TestAdapterResolveUsesCacheOnSecondCall(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.kdbx")
	if err := os.WriteFile(dbPath, []byte("dummy"), 0o600); err != nil {
		t.Fatal(err)
	}
	callCount := 0
	// We can't easily count calls inside a shell script, so we use a log file.
	logFile := filepath.Join(dir, "calls.log")
	bin := stubBin(t, dir, "keepassxc-cli", `
echo "$@" >> '`+logFile+`'
if [ "$1" = "show" ]; then
  printf "Title: KEY\nPassword: somevalue\nNotes: \n"
  exit 0
fi
exit 1
`)
	_ = callCount
	a := adapterWithStub(t, bin, dbPath, "dotenv", "pw")

	// First call hits the binary
	if _, err := a.Resolve(context.Background(), "MY_KEY", ""); err != nil {
		t.Fatalf("first resolve: %v", err)
	}
	// Second call should use cache — binary not invoked again
	if _, err := a.Resolve(context.Background(), "MY_KEY", ""); err != nil {
		t.Fatalf("second resolve: %v", err)
	}

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}
	// Count how many times show was called
	lines := 0
	for _, line := range splitLines(string(data)) {
		if line != "" {
			lines++
		}
	}
	if lines != 1 {
		t.Fatalf("expected 1 CLI invocation (cache hit on 2nd), got %d\nlog: %s", lines, data)
	}
}

func TestAdapterResolveMissingKeyReturnsE005(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.kdbx")
	if err := os.WriteFile(dbPath, []byte("dummy"), 0o600); err != nil {
		t.Fatal(err)
	}
	bin := stubBin(t, dir, "keepassxc-cli", `
echo "Entry not found." >&2
exit 1
`)
	a := adapterWithStub(t, bin, dbPath, "dotenv", "pw")
	res, err := a.Resolve(context.Background(), "MISSING_KEY", "")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if res.Source != "missing" || res.IssueCode != "E005" {
		t.Fatalf("expected missing/E005, got %+v", res)
	}
}

func TestAdapterResolveManyReturnsAllResolutions(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.kdbx")
	if err := os.WriteFile(dbPath, []byte("dummy"), 0o600); err != nil {
		t.Fatal(err)
	}
	bin := stubBin(t, dir, "keepassxc-cli", `
if [ "$1" = "show" ] && [ "$4" = "dotenv/DATABASE_URL" ]; then
  printf "Title: DATABASE_URL\nPassword: postgres://vault/dev\nNotes: \n"
  exit 0
fi
if [ "$1" = "show" ] && [ "$4" = "dotenv/JWT_SECRET" ]; then
  printf "Title: JWT_SECRET\nPassword: topsecret\nNotes: \n"
  exit 0
fi
echo "Entry not found." >&2
exit 1
`)
	a := adapterWithStub(t, bin, dbPath, "dotenv", "pw")
	results, err := a.ResolveMany(context.Background(), map[string]string{
		"DATABASE_URL": "",
		"JWT_SECRET":   "",
	})
	if err != nil {
		t.Fatalf("ResolveMany: %v", err)
	}
	if results["DATABASE_URL"].Value != "postgres://vault/dev" {
		t.Fatalf("unexpected DATABASE_URL: %+v", results["DATABASE_URL"])
	}
	if results["JWT_SECRET"].Value != "topsecret" {
		t.Fatalf("unexpected JWT_SECRET: %+v", results["JWT_SECRET"])
	}
}

func TestAdapterCheckReadinessBinaryMissing(t *testing.T) {
	a := NewAdapter(config.Config{
		KeePassDatabase: "/some/db.kdbx",
		KeePassGroup:    "dotenv",
	})
	a.client.Bin = "/nonexistent/keepassxc-cli"

	status, err := a.CheckReadiness(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.CLIInstalled {
		t.Fatal("expected CLIInstalled=false")
	}
	if status.Code != "E001" {
		t.Fatalf("expected E001, got %q", status.Code)
	}
}

func TestAdapterCheckReadinessDatabaseMissing(t *testing.T) {
	dir := t.TempDir()
	bin := stubBin(t, dir, "keepassxc-cli", `exit 0`)
	a := NewAdapter(config.Config{
		KeePassDatabase: "/nonexistent/db.kdbx",
		KeePassGroup:    "dotenv",
	})
	a.client.Bin = bin
	a.client.Password = "pw"

	status, err := a.CheckReadiness(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Code != "E002" {
		t.Fatalf("expected E002 for missing db, got %q", status.Code)
	}
}

func TestAdapterCheckReadinessGroupMissing(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.kdbx")
	if err := os.WriteFile(dbPath, []byte("dummy"), 0o600); err != nil {
		t.Fatal(err)
	}
	bin := stubBin(t, dir, "keepassxc-cli", `
echo "Could not find entry nosuchgroup." >&2
exit 1
`)
	a := adapterWithStub(t, bin, dbPath, "nosuchgroup", "pw")

	status, err := a.CheckReadiness(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Code != "E004" {
		t.Fatalf("expected E004 for missing group, got %q", status.Code)
	}
}

func TestAdapterCheckReadinessReady(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.kdbx")
	if err := os.WriteFile(dbPath, []byte("dummy"), 0o600); err != nil {
		t.Fatal(err)
	}
	bin := stubBin(t, dir, "keepassxc-cli", `
if [ "$1" = "ls" ]; then
  printf "DATABASE_URL\n"
  exit 0
fi
exit 1
`)
	a := adapterWithStub(t, bin, dbPath, "dotenv", "pw")

	status, err := a.CheckReadiness(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Authenticated || !status.Unlocked {
		t.Fatalf("expected ready status, got %+v", status)
	}
	if status.Code != "" {
		t.Fatalf("expected no error code, got %q", status.Code)
	}
}

// Verify Adapter satisfies the provider.Provider interface at compile time.
var _ provider.Provider = (*Adapter)(nil)

// splitLines is a small helper to split a string into non-empty lines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}