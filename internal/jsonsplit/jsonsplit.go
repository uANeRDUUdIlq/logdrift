// Package jsonsplit splits a single JSON log line into multiple lines
// by expanding an array field into one output line per element.
package jsonsplit

import (
	"encoding/json"
	"fmt"
)

// Splitter expands a named array field in a JSON log line into multiple
// individual JSON lines, one per array element. Non-JSON lines and lines
// where the target field is absent or not an array are passed through
// unchanged as a single-element slice.
type Splitter struct {
	field  string
	prefix string
}

// New returns a Splitter that expands the given array field.
// prefix is an optional key name used to nest the element value; when
// empty the element is merged into the parent object.
func New(field, prefix string) (*Splitter, error) {
	if field == "" {
		return nil, fmt.Errorf("jsonsplit: field must not be empty")
	}
	return &Splitter{field: field, prefix: prefix}, nil
}

// Apply expands the array field and returns one JSON string per element.
// If the line cannot be split it is returned as-is in a one-element slice.
func (s *Splitter) Apply(line string) []string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return []string{line}
	}

	raw, ok := obj[s.field]
	if !ok {
		return []string{line}
	}

	var elems []json.RawMessage
	if err := json.Unmarshal(raw, &elems); err != nil {
		return []string{line}
	}

	// Build a base object without the array field.
	base := make(map[string]json.RawMessage, len(obj)-1)
	for k, v := range obj {
		if k != s.field {
			base[k] = v
		}
	}

	out := make([]string, 0, len(elems))
	for _, elem := range elems {
		row := make(map[string]json.RawMessage, len(base)+1)
		for k, v := range base {
			row[k] = v
		}
		key := s.field
		if s.prefix != "" {
			key = s.prefix
		}
		row[key] = elem
		b, err := json.Marshal(row)
		if err != nil {
			continue
		}
		out = append(out, string(b))
	}

	if len(out) == 0 {
		return []string{line}
	}
	return out
}
