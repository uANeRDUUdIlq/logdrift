// Package jsonlimit truncates JSON array fields to a maximum number of elements.
package jsonlimit

import (
	"encoding/json"
)

// Limiter truncates JSON array fields that exceed a configured maximum length.
type Limiter struct {
	fields map[string]int
	defaultMax int
}

// New creates a Limiter. defaultMax applies to all array fields not explicitly
// listed in fields. A defaultMax of 0 means no default limit is applied.
func New(defaultMax int, fields map[string]int) *Limiter {
	f := make(map[string]int, len(fields))
	for k, v := range fields {
		if v > 0 {
			f[k] = v
		}
	}
	if defaultMax < 0 {
		defaultMax = 0
	}
	return &Limiter{fields: f, defaultMax: defaultMax}
}

// Apply parses line as JSON, truncates any array fields that exceed their
// configured limit, and returns the re-encoded line. Non-JSON lines are
// returned unchanged.
func (l *Limiter) Apply(line string) string {
	if len(l.fields) == 0 && l.defaultMax == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	modified := false
	for key, raw := range obj {
		var arr []json.RawMessage
		if err := json.Unmarshal(raw, &arr); err != nil {
			continue // not an array
		}

		max, ok := l.fields[key]
		if !ok {
			if l.defaultMax == 0 {
				continue
			}
			max = l.defaultMax
		}
		if max <= 0 || len(arr) <= max {
			continue
		}

		truncated, err := json.Marshal(arr[:max])
		if err != nil {
			continue
		}
		obj[key] = json.RawMessage(truncated)
		modified = true
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
