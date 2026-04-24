// Package jsonrename provides a transformer that bulk-renames JSON keys
// according to a static mapping before the line is forwarded downstream.
package jsonrename

import (
	"encoding/json"
)

// Renamer rewrites JSON keys according to a caller-supplied mapping.
// Keys not present in the mapping are left unchanged. Non-JSON lines
// pass through unmodified.
type Renamer struct {
	// mapping is old-key → new-key.
	mapping map[string]string
}

// New returns a Renamer that will replace every key found in mapping with
// its corresponding value. An empty mapping is valid and results in a
// no-op transformer.
func New(mapping map[string]string) *Renamer {
	copy := make(map[string]string, len(mapping))
	for k, v := range mapping {
		copy[k] = v
	}
	return &Renamer{mapping: copy}
}

// Apply renames keys in line according to the configured mapping.
// If line is not valid JSON, or the mapping is empty, line is returned
// unchanged.
func (r *Renamer) Apply(line string) string {
	if len(r.mapping) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for oldKey, newKey := range r.mapping {
		if oldKey == newKey {
			continue
		}
		val, ok := obj[oldKey]
		if !ok {
			continue
		}
		delete(obj, oldKey)
		obj[newKey] = val
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
