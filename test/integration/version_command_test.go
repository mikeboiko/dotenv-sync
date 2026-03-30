package integration_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVersionCommandsIntegration(t *testing.T) {
	t.Run("development build uses fallback metadata without touching env files", func(t *testing.T) {
		bin := buildCLI(t)
		project := t.TempDir()
		schemaPath := filepath.Join(project, ".env.example")
		envPath := filepath.Join(project, ".env")
		writeFile(t, schemaPath, "DATABASE_URL=\n")
		writeFile(t, envPath, "DATABASE_URL=postgres://localhost/dev\n")

		beforeSchema, err := os.ReadFile(schemaPath)
		if err != nil {
			t.Fatal(err)
		}
		beforeEnv, err := os.ReadFile(envPath)
		if err != nil {
			t.Fatal(err)
		}

		stdout, stderr, code := runCLI(t, bin, project, nil, "--version")
		if code != 0 || stderr != "" {
			t.Fatalf("short version failed: code=%d stderr=%q", code, stderr)
		}
		if stdout != "ds dev\n" {
			t.Fatalf("unexpected short version stdout: %q", stdout)
		}

		stdout, stderr, code = runCLI(t, bin, project, nil, "version")
		if code != 0 || stderr != "" {
			t.Fatalf("version command failed: code=%d stderr=%q", code, stderr)
		}
		if stdout == "" {
			t.Fatal("expected detailed version output")
		}

		afterSchema, err := os.ReadFile(schemaPath)
		if err != nil {
			t.Fatal(err)
		}
		afterEnv, err := os.ReadFile(envPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(beforeSchema) != string(afterSchema) {
			t.Fatalf("schema file changed: before=%q after=%q", beforeSchema, afterSchema)
		}
		if string(beforeEnv) != string(afterEnv) {
			t.Fatalf("env file changed: before=%q after=%q", beforeEnv, afterEnv)
		}
	})

	t.Run("ldflags build reports injected metadata", func(t *testing.T) {
		bin := buildCLIWithLdflags(t, "v2.3.4", "deadbee", "2026-03-30T12:34:56Z")
		project := t.TempDir()

		stdout, stderr, code := runCLI(t, bin, project, nil, "version")
		if code != 0 || stderr != "" {
			t.Fatalf("version command failed: code=%d stderr=%q", code, stderr)
		}
		want := "Version: v2.3.4\nCommit: deadbee\nBuilt: 2026-03-30T12:34:56Z\nPlatform: " + currentPlatform() + "\n"
		if stdout != want {
			t.Fatalf("detailed version stdout = %q want %q", stdout, want)
		}
	})
}
