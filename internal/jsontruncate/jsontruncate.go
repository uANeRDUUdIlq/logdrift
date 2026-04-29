// Package jsontruncate truncates string values in JSON log lines that exceed
// a configured maximum length, replacing the excess with a configurable suffix.
package jsontruncate

import (
	"encoding/json"
)

const defaultSuffix = "...[truncated]"

// Truncater truncates long string values within JSON objects.
type Truncater struct {
	maxLen int
	suffix string
	fields map[string]bool // if non-empty, only truncate these fields
}

// Option configures a Truncater.
type Option func(*Truncater)

// WithSuffix overrides the default truncation suffix.
func WithSuffix(s string) Option {
	return func(t *Truncater) { t.suffix = s }
}

// WithFields restricts truncation to the named top-level fields.
func WithFields(fields ...string) Option {
	return func(t *Truncater) {
		for _, f := range fields {
			t.fields[f] = true
		}
	}
}

// New creates a Truncater that shortens string values longer than maxLen.
// If maxLen <= 0 truncation is disabled and Apply is a no-op.
func New(maxLen int, opts ...Option) *Truncater {
	t := &Truncater{
		maxLen: maxLen,
		suffix: defaultSuffix,
		fields: make(map[string]bool),
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply truncates string values in the JSON line and returns the result.
// Non-JSON lines and lines where maxLen <= 0 are returned unchanged.
func (t *Truncater) Apply(line string) string {
	if t.maxLen <= 0 {
		return line
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	changed := false
	for k, v := range obj {
		if len(t.fields) > 0 && !t.fields[k] {
			continue
		}
		s, ok := v.(string)
		if !ok || len(s) <= t.maxLen {
			continue
		}
		obj[k] = s[:t.maxLen] + t.suffix
		changed = true
	}

	if !changed {
		return line
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(b)
}
