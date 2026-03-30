package bitwarden

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRBWClientGetRawItemParsesJSON(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "rbw")
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"get\" ] && [ \"$2\" = \"--raw\" ] && [ \"$3\" = \"repo\" ]; then\n" +
		"  cat <<'EOF'\n" +
		"{\"name\":\"repo\",\"notes\":\"{\\\"format\\\":\\\"dotenv-sync/note-json@v1\\\",\\\"env\\\":{}}\",\"data\":{\"password\":\"keep-me\"}}\n" +
		"EOF\n" +
		"  exit 0\n" +
		"fi\n" +
		"echo 'not found' >&2\n" +
		"exit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	item, err := (&RBWClient{Bin: bin}).GetRawItem(context.Background(), "repo")
	if err != nil {
		t.Fatalf("get raw item: %v", err)
	}
	if item.Name != "repo" || item.Password != "keep-me" {
		t.Fatalf("unexpected raw item: %+v", item)
	}
}

func TestRBWClientMutationsUseScriptedEditor(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "rbw")
	logPath := filepath.Join(dir, "rbw.log")
	addCapture := filepath.Join(dir, "add.txt")
	editCapture := filepath.Join(dir, "edit.txt")
	script := "#!/bin/sh\n" +
		"echo \"$@\" >> '" + logPath + "'\n" +
		"tmp=$(mktemp)\n" +
		"case \"$1\" in\n" +
		"add)\n" +
		"  \"${VISUAL:-$EDITOR}\" \"$tmp\"\n" +
		"  cat \"$tmp\" > '" + addCapture + "'\n" +
		"  ;;\n" +
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

	client := &RBWClient{Bin: bin}
	if err := client.AddItem(context.Background(), "repo", "", "{\"env\":{}}"); err != nil {
		t.Fatalf("add item: %v", err)
	}
	if err := client.EditItem(context.Background(), "repo", "secret-password", "{\"env\":{\"A\":\"1\"}}"); err != nil {
		t.Fatalf("edit item: %v", err)
	}

	addData, err := os.ReadFile(addCapture)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(addData), "{\"env\":{}}") {
		t.Fatalf("add editor content missing notes: %q", string(addData))
	}
	editData, err := os.ReadFile(editCapture)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(editData), "secret-password\n") {
		t.Fatalf("edit editor content missing password: %q", string(editData))
	}
	logData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(logData), "add repo") || !strings.Contains(string(logData), "edit repo") {
		t.Fatalf("unexpected rbw log: %s", logData)
	}
}

func TestRBWClientGetRawItemTreatsNoEntryFoundAsMissing(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "rbw")
	script := "#!/bin/sh\n" +
		"echo \"rbw get: couldn't find entry for '$3': no entry found\" >&2\n" +
		"exit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := (&RBWClient{Bin: bin}).GetRawItem(context.Background(), "repo")
	if !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("expected ErrItemNotFound, got %v", err)
	}
}
