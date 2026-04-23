// Package jsonpatch applies RFC-7396-style merge patches to JSON log lines.
// Non-JSON lines are passed through unchanged.
package jsonpatch

import (
	"encoding/json"
)

// Patcher applies a fixed set of key/value overrides to every JSON log line.
type Patcher struct {
	patch map[string]json.RawMessage
}

// New returns a Patcher that will apply the given patch map to each line.
// Values in patch must be valid JSON (strings, numbers, booleans, null, objects).
// An empty patch map is valid and results in a no-op patcher.
func New(patch map[string]string) (*Patcher, error) {
	p := make(map[string]json.RawMessage, len(patch))
	for k, v := range patch {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		p[k] = b
	}
	return &Patcher{patch: p}, nil
}

// Apply merges the configured patch into line.
// If line is not valid JSON, it is returned unchanged.
// Existing keys are overwritten; new keys are added.
func (p *Patcher) Apply(line string) string {
	if len(p.patch) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for k, v := range p.patch {
		obj[k] = v
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(b)
}
