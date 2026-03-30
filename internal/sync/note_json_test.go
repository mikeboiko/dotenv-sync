package sync

import (
	"strings"
	"testing"

	"dotenv-sync/internal/envfile"
)

func TestRenderAndParseNoteJSONRoundTrip(t *testing.T) {
	doc := envfile.ParseBytes(".env", envfile.KindLocal, []byte("DATABASE_URL=postgres://vault/dev\nPRIVATE_KEY=\"line1\\nline2\"\n"))

	rendered, err := RenderNoteJSON(CanonicalDocumentEnv(doc))
	if err != nil {
		t.Fatalf("render note_json: %v", err)
	}
	if !strings.Contains(rendered, "\"format\":\""+NoteJSONFormat+"\"") {
		t.Fatalf("missing format marker in %s", rendered)
	}

	parsed, err := ParseNoteJSON(rendered)
	if err != nil {
		t.Fatalf("parse note_json: %v", err)
	}
	if !NoteJSONEqual(parsed.Env, map[string]string{
		"DATABASE_URL": "postgres://vault/dev",
		"PRIVATE_KEY":  "\"line1\\nline2\"",
	}) {
		t.Fatalf("unexpected payload env: %+v", parsed.Env)
	}
}

func TestParseNoteJSONRejectsMalformedPayload(t *testing.T) {
	if _, err := ParseNoteJSON("not-json"); err == nil {
		t.Fatal("expected malformed payload error")
	}
}

func TestNoteJSONEqualIgnoresMapOrder(t *testing.T) {
	left := map[string]string{"B": "2", "A": "1"}
	right := map[string]string{"A": "1", "B": "2"}
	if !NoteJSONEqual(left, right) {
		t.Fatal("expected semantic equality")
	}
}
