// Package jsonstrip removes specified top-level or nested keys from JSON log lines.
package jsonstrip

import (
	"encoding/json"
	"strings"
)

// Stripper removes a fixed set of JSON keys from each log line.
type Stripper struct {
	keys map[string]struct{}
}

// New returns a Stripper that will remove all keys listed in keys.
// Keys may use dot notation to target nested fields (e.g. "meta.internal").
// If keys is empty every line is returned unchanged.
func New(keys []string) *Stripper {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return &Stripper{keys: m}
}

// Apply removes the configured keys from line and returns the result.
// Non-JSON lines and lines that contain none of the keys are returned as-is.
func (s *Stripper) Apply(line string) string {
	if len(s.keys) == 0 {
		return line
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	changed := stripKeys(obj, s.keys, "")
	if !changed {
		return line
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// stripKeys recursively removes matching keys and reports whether any were removed.
func stripKeys(obj map[string]interface{}, keys map[string]struct{}, prefix string) bool {
	changed := false
	for k, v := range obj {
		full := k
		if prefix != "" {
			full = prefix + "." + k
		}
		if _, ok := keys[full]; ok {
			delete(obj, k)
			changed = true
			continue
		}
		if nested, ok := v.(map[string]interface{}); ok {
			if stripKeys(nested, keys, full) {
				changed = true
			}
		}
	}
	_ = strings.Contains // keep import used via dot-notation split above
	return changed
}
