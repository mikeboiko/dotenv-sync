package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
)

type fakePushProvider struct {
	payload     provider.EnvPayload
	storeCount  int
	loadCount   int
	stored      provider.EnvPayload
	readiness   provider.Status
	readinessOK bool
}

func (f *fakePushProvider) Name() string { return "fake" }

func (f *fakePushProvider) CheckReadiness(_ context.Context) (provider.Status, error) {
	if !f.readinessOK {
		return provider.Status{Provider: "bitwarden", CLIInstalled: true, Authenticated: true, Unlocked: true}, nil
	}
	return f.readiness, nil
}

func (f *fakePushProvider) Resolve(_ context.Context, key, ref string) (provider.Resolution, error) {
	return provider.Resolution{Key: key, Ref: ref}, nil
}

func (f *fakePushProvider) ResolveMany(_ context.Context, refs map[string]string) (map[string]provider.Resolution, error) {
	return map[string]provider.Resolution{}, nil
}

func (f *fakePushProvider) LoadEnvPayload(_ context.Context) (provider.EnvPayload, error) {
	f.loadCount++
	return f.payload, nil
}

func (f *fakePushProvider) StoreEnvPayload(_ context.Context, payload provider.EnvPayload) (provider.WriteResult, error) {
	f.storeCount++
	f.stored = payload
	return provider.WriteResult{ItemName: payload.ItemName, Created: !payload.Exists, Updated: payload.Exists}, nil
}

func TestPushCommandDryRunPrintsRedactedPreview(t *testing.T) {
	project := t.TempDir()
	writeFileForCommandTest(t, filepath.Join(project, ".env.example"), "DATABASE_URL=\nJWT_SECRET=\n")
	writeFileForCommandTest(t, filepath.Join(project, ".env"), "DATABASE_URL=postgres://vault/dev\nJWT_SECRET=supersecret\n")
	writeFileForCommandTest(t, filepath.Join(project, ".envsync.yaml"), "storage_mode: note_json\nitem_name: Jesse\n")

	fake := &fakePushProvider{payload: provider.EnvPayload{ItemName: "Jesse", Exists: false, Env: map[string]string{}}}
	restore := swapPushProviderFactory(func(config.Config) provider.PushProvider { return fake })
	defer restore()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)
	if err := os.Chdir(project); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	cmd := newPushCommand(streams{stdout: &stdout, stderr: &stderr}, &rootOptions{})
	cmd.SetArgs([]string{"--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute push --dry-run: %v", err)
	}
	if fake.storeCount != 0 {
		t.Fatalf("expected no store in dry-run, got %d", fake.storeCount)
	}
	for _, want := range []string{"ADD DATABASE_URL [REDACTED]", "ADD JWT_SECRET [REDACTED]", "CHECKED bitwarden:Jesse"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q\n%s", want, stdout.String())
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func TestPushCommandRejectsFieldsMode(t *testing.T) {
	project := t.TempDir()
	writeFileForCommandTest(t, filepath.Join(project, ".env.example"), "DATABASE_URL=\n")
	writeFileForCommandTest(t, filepath.Join(project, ".env"), "DATABASE_URL=value\n")

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)
	if err := os.Chdir(project); err != nil {
		t.Fatal(err)
	}

	cmd := newPushCommand(streams{stdout: &bytes.Buffer{}, stderr: &bytes.Buffer{}}, &rootOptions{})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "storage_mode: note_json") {
		t.Fatalf("expected fields mode error, got %v", err)
	}
}

func writeFileForCommandTest(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func swapPushProviderFactory(factory func(config.Config) provider.PushProvider) func() {
	previous := pushProviderFactory
	pushProviderFactory = factory
	return func() {
		pushProviderFactory = previous
	}
}
