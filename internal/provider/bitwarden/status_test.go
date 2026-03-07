package bitwarden

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckReadinessWithMissingBinary(t *testing.T) {
	status, err := checkReadinessWithClient(context.Background(), &RBWClient{Bin: filepath.Join(t.TempDir(), "rbw-missing")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Code != "E001" {
		t.Fatalf("expected E001, got %+v", status)
	}
}

func TestCheckReadinessWithStubStatuses(t *testing.T) {
	cases := []struct {
		name   string
		status string
		code   string
	}{
		{name: "unlocked", status: "unlocked", code: ""},
		{name: "locked", status: "locked", code: "E003"},
		{name: "logged-out", status: "logged out", code: "E002"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			bin := filepath.Join(dir, "rbw")
			script := "#!/bin/sh\nif [ \"$1\" = \"status\" ]; then printf '%s\\n' '" + tc.status + "'; exit 0; fi\nexit 1\n"
			if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
				t.Fatal(err)
			}
			status, err := checkReadinessWithClient(context.Background(), &RBWClient{Bin: bin})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if status.Code != tc.code {
				t.Fatalf("expected code %q, got %+v", tc.code, status)
			}
		})
	}
}
