// Package redact provides field-level redaction for structured JSON log lines.
package redact

import (
	"encoding/json"
	"strings"
)

// Redactor replaces sensitive field values in parsed JSON log lines.
type Redactor struct {
	fields map[string]struct{}
	mask   string
}

// New creates a Redactor that masks the given field names with the provided mask string.
// If mask is empty, "[REDACTED]" is used.
func New(fields []string, mask string) *Redactor {
	if mask == "" {
		mask = "[REDACTED]"
	}
	fm := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		fm[strings.ToLower(f)] = struct{}{}
	}
	return &Redactor{fields: fm, mask: mask}
}

// Apply redacts sensitive fields from a raw JSON log line.
// Non-JSON lines are returned unchanged.
func (r *Redactor) Apply(line string) string {
	if len(r.fields) == 0 {
		return line
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	modified := false
	for k := range obj {
		if _, ok := r.fields[strings.ToLower(k)]; ok {
			obj[k] = r.mask
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
