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
		script string
		code   string
	}{
		{
			name: "unlocked",
			script: "#!/bin/sh\ncase \"$1\" in\n" +
				"unlocked) exit 0 ;;\n" +
				"list) printf 'DATABASE_URL\\n' ;;\n" +
				"*) exit 1 ;;\n" +
				"esac\n",
			code: "",
		},
		{
			name: "locked",
			script: "#!/bin/sh\ncase \"$1\" in\n" +
				"unlocked) exit 1 ;;\n" +
				"list) echo 'database is locked' >&2; exit 1 ;;\n" +
				"*) exit 1 ;;\n" +
				"esac\n",
			code: "E003",
		},
		{
			name: "logged-out",
			script: "#!/bin/sh\ncase \"$1\" in\n" +
				"unlocked) exit 1 ;;\n" +
				"list) echo 'not logged in' >&2; exit 1 ;;\n" +
				"*) exit 1 ;;\n" +
				"esac\n",
			code: "E002",
		},
		{
			name: "legacy-status-fallback",
			script: "#!/bin/sh\ncase \"$1\" in\n" +
				"unlocked) echo \"error: unrecognized subcommand 'unlocked'\" >&2; exit 2 ;;\n" +
				"status) printf 'unlocked\\n' ;;\n" +
				"*) exit 1 ;;\n" +
				"esac\n",
			code: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			bin := filepath.Join(dir, "rbw")
			if err := os.WriteFile(bin, []byte(tc.script), 0o755); err != nil {
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
