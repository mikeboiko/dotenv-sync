package envfile

import "testing"

func TestRenderPreservesCRLFAndTrailingNewline(t *testing.T) {
	doc := ParseBytes(".env.example", KindSchema, []byte("A=1\r\n# comment\r\nB=\r\n"))
	doc.Lines[2].Value = "2"
	got := string(Render(doc))
	want := "A=1\r\n# comment\r\nB=2\r\n"
	if got != want {
		t.Fatalf("render mismatch\nwant: %q\ngot:  %q", want, got)
	}
}
