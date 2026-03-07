package envfile

import (
	"fmt"
	"os"
	"strings"
)

func ParseFile(path string, kind Kind) (Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}
	return ParseBytes(path, kind, data), nil
}

func ParseBytes(path string, kind Kind, data []byte) Document {
	raw := string(data)
	lineEnding := "\n"
	if strings.Contains(raw, "\r\n") {
		lineEnding = "\r\n"
	}
	normalized := strings.ReplaceAll(raw, "\r\n", "\n")
	parts := strings.Split(normalized, "\n")
	trailing := strings.HasSuffix(normalized, "\n")
	if trailing {
		parts = parts[:len(parts)-1]
	}
	doc := Document{Path: path, Kind: kind, LineEnding: lineEnding, TrailingNewline: trailing, Raw: raw}
	seen := map[string]int{}
	for i, rawLine := range parts {
		line := parseLine(i, rawLine, kind)
		doc.Lines = append(doc.Lines, line)
		if line.LineType == LineInvalid {
			doc.ParseErrors = append(doc.ParseErrors, fmt.Sprintf("line %d: malformed dotenv entry", i+1))
		}
		if line.LineType == LineAssignment {
			seen[line.Key]++
			if seen[line.Key] > 1 {
				doc.Duplicates = append(doc.Duplicates, line.Key)
			}
		}
	}
	return doc
}

func parseLine(index int, raw string, kind Kind) EnvironmentLine {
	line := EnvironmentLine{Index: index, Raw: raw}
	trimmed := strings.TrimSpace(raw)
	switch {
	case trimmed == "":
		line.LineType = LineBlank
		return line
	case strings.HasPrefix(trimmed, "#"):
		line.LineType = LineComment
		return line
	}
	eq := strings.Index(raw, "=")
	if eq <= 0 {
		line.LineType = LineInvalid
		return line
	}
	keyPart := raw[:eq]
	valuePart := raw[eq+1:]
	key := strings.TrimSpace(keyPart)
	value, suffix := splitInlineComment(valuePart)
	line.LineType = LineAssignment
	line.Key = key
	line.Value = value
	line.Prefix = raw[:eq+1]
	line.Suffix = suffix
	line.InlineComment = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(suffix), "#"))
	line.ManagedByProvider = kind == KindSchema && strings.TrimSpace(value) == ""
	return line
}

func splitInlineComment(value string) (string, string) {
	inSingle := false
	inDouble := false
	for i := 0; i < len(value); i++ {
		switch value[i] {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble {
				if i == 0 || value[i-1] == ' ' || value[i-1] == '\t' {
					return strings.TrimRight(value[:i], " \t"), value[i:]
				}
			}
		}
	}
	return value, ""
}
