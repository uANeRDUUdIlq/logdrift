// Package jsonmerge provides a Merger that merges a static set of key/value
// fields into every JSON log line that passes through it.  Non-JSON lines are
// returned unchanged.
package jsonmerge

import (
	"encoding/json"
)

// Merger merges a fixed map of fields into JSON log lines.
type Merger struct {
	fields map[string]string
}

// New returns a Merger that will inject fields into every JSON line.
// If fields is empty the returned Merger is a no-op.
func New(fields map[string]string) *Merger {
	copy := make(map[string]string, len(fields))
	for k, v := range fields {
		copy[k] = v
	}
	return &Merger{fields: copy}
}

// Apply merges the configured fields into line.  Existing keys in the line are
// NOT overwritten — the original value wins.  Non-JSON lines are returned as-is.
func (m *Merger) Apply(line string) string {
	if len(m.fields) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for k, v := range m.fields {
		if _, exists := obj[k]; exists {
			continue
		}
		encoded, err := json.Marshal(v)
		if err != nil {
			continue
		}
		obj[k] = json.RawMessage(encoded)
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
