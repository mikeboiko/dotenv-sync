package keepass

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// stubBin writes a shell script to dir/<name> and returns its path.
// The script is made executable so exec.LookPath and exec.Command can use it.
func stubBin(t *testing.T, dir, name, script string) string {
	t.Helper()
	bin := filepath.Join(dir, name)
	if err := os.WriteFile(bin, []byte("#!/bin/sh\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
	return bin
}

// clientWithBin returns a KPXCClient wired to a stub binary with the given password.
func clientWithBin(bin, password string) *KPXCClient {
	return &KPXCClient{Bin: bin, Password: password}
}

// --- parsePasswordField ---

func TestParsePasswordFieldExtractsValue(t *testing.T) {
	output := "Title: DATABASE_URL\nUserName: \nPassword: postgres://localhost:5432/testdb\nURL: \nNotes: \n"
	got, err := parsePasswordField(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "postgres://localhost:5432/testdb" {
		t.Fatalf("expected postgres URL, got %q", got)
	}
}

func TestParsePasswordFieldHandlesEmptyValue(t *testing.T) {
	output := "Title: MY_KEY\nPassword: \nNotes: \n"
	got, err := parsePasswordField(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestParsePasswordFieldErrorsWhenMissing(t *testing.T) {
	output := "Title: MY_KEY\nNotes: no password line here\n"
	_, err := parsePasswordField(output)
	if err == nil {
		t.Fatal("expected error when Password field missing")
	}
}

// --- isNotFoundText ---

func TestIsNotFoundTextMatchesKnownPhrases(t *testing.T) {
	cases := []string{
		"Entry not found.",
		"Could not find entry dotenv/MISSING_KEY",
		"no such entry",
		"no entry for this path",
	}
	for _, c := range cases {
		if !isNotFoundText(c) {
			t.Errorf("expected isNotFoundText(%q) = true", c)
		}
	}
}

func TestIsNotFoundTextIgnoresUnrelatedErrors(t *testing.T) {
	cases := []string{
		"Invalid credentials were provided",
		"HMAC mismatch",
		"database file is locked",
	}
	for _, c := range cases {
		if isNotFoundText(c) {
			t.Errorf("expected isNotFoundText(%q) = false", c)
		}
	}
}

// --- KPXCClient.Show ---

func TestKPXCClientShowReturnsPasswordValue(t *testing.T) {
	dir := t.TempDir()
	bin := stubBin(t, dir, "keepassxc-cli", `
if [ "$1" = "show" ] && [ "$2" = "-s" ]; then
  printf "Title: DATABASE_URL\nUserName: \nPassword: postgres://vault/dev\nURL: \nNotes: \n"
  exit 0
fi
echo "unexpected args" >&2
exit 1
`)
	client := clientWithBin(bin, "masterpassword")
	value, err := client.Show(context.Background(), "test.kdbx", "dotenv/DATABASE_URL")
	if err != nil {
		t.Fatalf("Show: %v", err)
	}
	if value != "postgres://vault/dev" {
		t.Fatalf("expected postgres URL, got %q", value)
	}
}

func TestKPXCClientShowReturnsErrItemNotFound(t *testing.T) {
	dir := t.TempDir()
	bin := stubBin(t, dir, "keepassxc-cli", `
echo "Entry not found." >&2
exit 1
`)
	client := clientWithBin(bin, "masterpassword")
	_, err := client.Show(context.Background(), "test.kdbx", "dotenv/MISSING")
	if !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("expected ErrItemNotFound, got %v", err)
	}
}

func TestKPXCClientShowReturnsErrBinaryMissing(t *testing.T) {
	client := clientWithBin("/nonexistent/keepassxc-cli", "pw")
	_, err := client.Show(context.Background(), "test.kdbx", "dotenv/KEY")
	if !errors.Is(err, ErrBinaryMissing) {
		t.Fatalf("expected ErrBinaryMissing, got %v", err)
	}
}

// --- KPXCClient.ListGroup ---

func TestKPXCClientListGroupReturnsEntries(t *testing.T) {
	dir := t.TempDir()
	bin := stubBin(t, dir, "keepassxc-cli", `
if [ "$1" = "ls" ]; then
  printf "DATABASE_URL\nJWT_SECRET\nsubgroup/\n"
  exit 0
fi
exit 1
`)
	client := clientWithBin(bin, "masterpassword")
	entries, err := client.ListGroup(context.Background(), "test.kdbx", "dotenv")
	if err != nil {
		t.Fatalf("ListGroup: %v", err)
	}
	// subgroup/ should be filtered out, leaving two entries
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(entries), entries)
	}
	if entries[0] != "DATABASE_URL" || entries[1] != "JWT_SECRET" {
		t.Fatalf("unexpected entries: %v", entries)
	}
}

func TestKPXCClientListGroupReturnsErrItemNotFound(t *testing.T) {
	dir := t.TempDir()
	bin := stubBin(t, dir, "keepassxc-cli", `
echo "Could not find entry nosuchgroup." >&2
exit 1
`)
	client := clientWithBin(bin, "masterpassword")
	_, err := client.ListGroup(context.Background(), "test.kdbx", "nosuchgroup")
	if !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("expected ErrItemNotFound, got %v", err)
	}
}