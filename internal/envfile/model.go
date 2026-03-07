package envfile

import "slices"

type Kind string

type LineType string

const (
	KindSchema Kind = "schema"
	KindLocal  Kind = "local"

	LineBlank      LineType = "blank"
	LineComment    LineType = "comment"
	LineAssignment LineType = "assignment"
	LineInvalid    LineType = "invalid"
)

type EnvironmentLine struct {
	Index             int
	Raw               string
	LineType          LineType
	Key               string
	Value             string
	InlineComment     string
	ManagedByProvider bool
	Prefix            string
	Suffix            string
}

type Document struct {
	Path            string
	Kind            Kind
	LineEnding      string
	TrailingNewline bool
	Lines           []EnvironmentLine
	ParseErrors     []string
	Duplicates      []string
	Raw             string
}

func (d Document) Clone() Document {
	clone := d
	clone.Lines = slices.Clone(d.Lines)
	clone.ParseErrors = slices.Clone(d.ParseErrors)
	clone.Duplicates = slices.Clone(d.Duplicates)
	return clone
}

func (d Document) AssignmentMap() map[string]EnvironmentLine {
	result := map[string]EnvironmentLine{}
	for _, line := range d.Lines {
		if line.LineType == LineAssignment {
			result[line.Key] = line
		}
	}
	return result
}

func (d Document) Keys() []string {
	keys := make([]string, 0)
	for _, line := range d.Lines {
		if line.LineType == LineAssignment {
			keys = append(keys, line.Key)
		}
	}
	return keys
}
