// Package jsondefault fills in missing JSON fields with configured default values.
// Fields that already exist in the log line are left untouched.
package jsondefault

import (
	"encoding/json"
)

// Defaulter applies default values to JSON log lines.
type Defaulter struct {
	defaults map[string]any
}

// New creates a Defaulter that will inject the provided key/value pairs into
// any JSON line that is missing those keys. Values must be JSON-serialisable;
// non-serialisable values are silently skipped during construction.
func New(defaults map[string]any) *Defaulter {
	safe := make(map[string]any, len(defaults))
	for k, v := range defaults {
		if _, err := json.Marshal(v); err == nil {
			safe[k] = v
		}
	}
	return &Defaulter{defaults: safe}
}

// Apply returns the (possibly modified) log line. If the line is not valid JSON
// it is returned unchanged. Only keys absent from the original object are
// injected; existing keys are never overwritten.
func (d *Defaulter) Apply(line string) string {
	if len(d.defaults) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	modified := false
	for k, v := range d.defaults {
		if _, exists := obj[k]; !exists {
			raw, _ := json.Marshal(v)
			obj[k] = json.RawMessage(raw)
			modified = true
		}
	}

	if !modified {
		return line
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
