// Package jsonpath provides dot-notation field extraction from JSON log lines.
package jsonpath

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Extractor extracts values from JSON using dot-notation paths.
type Extractor struct{}

// New returns a new Extractor.
func New() *Extractor {
	return &Extractor{}
}

// Get retrieves a value from a JSON line using a dot-notation path.
// Returns the value as a string and true if found, or "", false otherwise.
func (e *Extractor) Get(line, path string) (string, bool) {
	if line == "" || path == "" {
		return "", false
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", false
	}

	parts := strings.Split(path, ".")
	var current interface{} = obj

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return "", false
		}
		current, ok = m[part]
		if !ok {
			return "", false
		}
	}

	switch v := current.(type) {
	case string:
		return v, true
	case float64:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", v), "0"), "."), true
	case bool:
		if v {
			return "true", true
		}
		return "false", true
	case nil:
		return "null", true
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", false
		}
		return string(b), true
	}
}
