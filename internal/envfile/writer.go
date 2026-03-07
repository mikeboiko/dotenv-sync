package envfile

import (
	"bytes"
	"os"
	"strings"

	"dotenv-sync/internal/fs"
)

func Render(doc Document) []byte {
	lines := make([]string, 0, len(doc.Lines))
	for _, line := range doc.Lines {
		switch line.LineType {
		case LineAssignment:
			lines = append(lines, line.Prefix+line.Value+line.Suffix)
		default:
			lines = append(lines, line.Raw)
		}
	}
	output := strings.Join(lines, doc.LineEnding)
	if doc.TrailingNewline {
		output += doc.LineEnding
	}
	return []byte(output)
}

func EqualContent(doc Document, data []byte) bool {
	return bytes.Equal(Render(doc), data)
}

func WriteDocument(path string, doc Document) (bool, error) {
	perm := os.FileMode(0o644)
	return fs.WriteFileAtomicIfChanged(path, Render(doc), perm)
}
