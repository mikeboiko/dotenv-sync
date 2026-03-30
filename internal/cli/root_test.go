package cli

import (
	"bytes"
	"testing"
)

func TestNewRootCommandRegistersVersionSurface(t *testing.T) {
	cmd := NewRootCommand(streams{stdout: &bytes.Buffer{}, stderr: &bytes.Buffer{}})

	if cmd.Version == "" {
		t.Fatal("expected root version to be configured")
	}
	if cmd.Name() != "ds" {
		t.Fatalf("command name = %q", cmd.Name())
	}
	found := false
	for _, child := range cmd.Commands() {
		if child.Name() == "version" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected version subcommand to be registered")
	}
}
