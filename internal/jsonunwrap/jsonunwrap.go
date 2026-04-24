// Package jsonunwrap promotes fields from a nested JSON object to the top level.
//
// Given a line like:
//
//	{"meta":{"host":"web-1","env":"prod"},"msg":"ok"}
//
// Unwrapping "meta" yields:
//
//	{"host":"web-1","env":"prod","msg":"ok"}
//
// Existing top-level keys are never overwritten.
package jsonunwrap

import (
	"encoding/json"
)

// Unwrapper promotes fields from a nested key to the top-level object.
type Unwrapper struct {
	fields []string // nested keys whose children should be promoted
}

// New returns an Unwrapper that will promote the contents of each named
// nested field. If fields is empty, Apply is a no-op.
func New(fields []string) *Unwrapper {
	copy := make([]string, len(fields))
	for i, f := range fields {
		copy[i] = f
	}
	return &Unwrapper{fields: copy}
}

// Apply unwraps the configured nested fields in line and returns the result.
// If line is not valid JSON or no configured fields are present, line is
// returned unchanged.
func (u *Unwrapper) Apply(line string) string {
	if len(u.fields) == 0 {
		return line
	}

	var top map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &top); err != nil {
		return line
	}

	modified := false
	for _, field := range u.fields {
		raw, ok := top[field]
		if !ok {
			continue
		}
		var nested map[string]json.RawMessage
		if err := json.Unmarshal(raw, &nested); err != nil {
			continue
		}
		delete(top, field)
		for k, v := range nested {
			if _, exists := top[k]; !exists {
				top[k] = v
			}
		}
		modified = true
	}

	if !modified {
		return line
	}

	out, err := json.Marshal(top)
	if err != nil {
		return line
	}
	return string(out)
}
