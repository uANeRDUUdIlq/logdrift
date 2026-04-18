// Package fieldselect provides a Selector that keeps only a specified subset
// of top-level JSON fields in each log line, discarding the rest.
package fieldselect

import (
	"encoding/json"
)

// Selector filters JSON log lines to include only the configured fields.
type Selector struct {
	fields map[string]struct{}
}

// New creates a Selector that retains only the given field names.
// If fields is empty the Selector is a no-op.
func New(fields []string) *Selector {
	m := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		m[f] = struct{}{}
	}
	return &Selector{fields: m}
}

// Apply returns a new JSON line containing only the selected fields.
// If no fields are configured, or the line is not valid JSON, the original
// line is returned unchanged.
func (s *Selector) Apply(line string) string {
	if len(s.fields) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	out := make(map[string]json.RawMessage, len(s.fields))
	for k, v := range obj {
		if _, ok := s.fields[k]; ok {
			out[k] = v
		}
	}

	b, err := json.Marshal(out)
	if err != nil {
		return line
	}
	return string(b)
}
