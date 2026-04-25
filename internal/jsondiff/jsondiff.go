// Package jsondiff emits a line only when one or more tracked JSON fields
// change value compared to the previously seen line for that service.
package jsondiff

import (
	"encoding/json"
	"sync"
)

// Differ holds the last-seen values for each (service, field) pair and
// forwards a line only when at least one watched field has changed.
type Differ struct {
	mu     sync.Mutex
	fields []string
	// last maps service -> field -> last value (JSON-encoded scalar).
	last map[string]map[string]string
}

// New returns a Differ that watches the given list of top-level JSON fields.
// If fields is empty every line is forwarded unchanged.
func New(fields []string) *Differ {
	copy := make([]string, len(fields))
	for i, f := range fields {
		copy[i] = f
	}
	return &Differ{
		fields: copy,
		last:   make(map[string]map[string]string),
	}
}

// Changed returns true when at least one watched field in line differs from
// the value seen in the previous call for the same service, and updates the
// stored snapshot.  Non-JSON lines always return true.
func (d *Differ) Changed(service, line string) bool {
	if len(d.fields) == 0 {
		return true
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return true
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	prev, ok := d.last[service]
	if !ok {
		prev = make(map[string]string)
		d.last[service] = prev
	}

	changed := false
	for _, f := range d.fields {
		current := ""
		if raw, exists := obj[f]; exists {
			current = string(raw)
		}
		if prev[f] != current {
			changed = true
			prev[f] = current
		}
	}
	return changed
}

// Reset clears all stored snapshots.
func (d *Differ) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.last = make(map[string]map[string]string)
}
