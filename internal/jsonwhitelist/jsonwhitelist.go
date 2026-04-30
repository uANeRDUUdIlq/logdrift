// Package jsonwhitelist drops any JSON fields not present in an allowed set.
// Non-JSON lines are passed through unchanged.
package jsonwhitelist

import (
	"encoding/json"
)

// Whitelister removes all JSON keys not present in the allowed set.
type Whitelister struct {
	allowed map[string]struct{}
}

// New returns a Whitelister that retains only the given fields.
// If fields is empty the Whitelister is a no-op.
func New(fields []string) *Whitelister {
	allowed := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		allowed[f] = struct{}{}
	}
	return &Whitelister{allowed: allowed}
}

// Apply returns the line with only whitelisted fields kept.
// Non-JSON lines and empty field lists are returned unchanged.
func (w *Whitelister) Apply(line string) string {
	if len(w.allowed) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	out := make(map[string]json.RawMessage, len(w.allowed))
	for k, v := range obj {
		if _, ok := w.allowed[k]; ok {
			out[k] = v
		}
	}

	b, err := json.Marshal(out)
	if err != nil {
		return line
	}
	return string(b)
}
