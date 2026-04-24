// Package jsonflatten provides a pipeline stage that flattens nested JSON
// objects into a single-level map using a configurable separator.
//
// Example:
//
//	{"a": {"b": {"c": 1}}} → {"a.b.c": 1}
package jsonflatten

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Flattener flattens nested JSON log lines.
type Flattener struct {
	separator string
	prefix    string
}

// New returns a Flattener that joins nested keys with sep.
// If sep is empty, "." is used. prefix is prepended to every top-level key
// (may be empty).
func New(sep, prefix string) *Flattener {
	if sep == "" {
		sep = "."
	}
	return &Flattener{separator: sep, prefix: prefix}
}

// Apply flattens the JSON object in line. Non-JSON lines are returned
// unchanged.
func (f *Flattener) Apply(line string) string {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	flat := make(map[string]any)
	flatten(obj, f.prefix, f.separator, flat)

	b, err := json.Marshal(flat)
	if err != nil {
		return line
	}
	return string(b)
}

// flatten recursively walks v, building dotted keys into dst.
func flatten(v map[string]any, prefix, sep string, dst map[string]any) {
	for k, val := range v {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}
		switch child := val.(type) {
		case map[string]any:
			flatten(child, key, sep, dst)
		default:
			dst[key] = formatScalar(val)
		}
	}
}

// formatScalar converts numeric JSON values back to their natural form so that
// integer-like floats are not serialised with a decimal point.
func formatScalar(v any) any {
	switch n := v.(type) {
	case float64:
		if n == float64(int64(n)) {
			return int64(n)
		}
		return fmt.Sprintf("%g", n)
		// keep strings, bools, nil as-is
	}
	_ = strings.Contains // keep import used via doc only
	return v
}
