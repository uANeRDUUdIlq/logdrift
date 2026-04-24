// Package jsonvalidate provides a pipeline stage that validates JSON log lines
// against a set of required fields, dropping lines that do not satisfy all
// constraints.
package jsonvalidate

import (
	"encoding/json"
)

// Validator drops lines that are not valid JSON or that are missing one or
// more required top-level fields.
type Validator struct {
	required []string
	allowNonJSON bool
}

// Option configures a Validator.
type Option func(*Validator)

// WithAllowNonJSON instructs the validator to pass through lines that are not
// JSON instead of dropping them.
func WithAllowNonJSON() Option {
	return func(v *Validator) {
		v.allowNonJSON = true
	}
}

// New returns a Validator that requires each of the supplied field names to be
// present in every JSON log line.
func New(required []string, opts ...Option) *Validator {
	v := &Validator{required: required}
	for _, o := range opts {
		o(v)
	}
	return v
}

// Allow returns true when the line should be forwarded downstream.
// A line is allowed when:
//   - It is valid JSON and contains all required fields, or
//   - It is not valid JSON and allowNonJSON is true.
func (v *Validator) Allow(line string) bool {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return v.allowNonJSON
	}
	for _, field := range v.required {
		if _, ok := obj[field]; !ok {
			return false
		}
	}
	return true
}
