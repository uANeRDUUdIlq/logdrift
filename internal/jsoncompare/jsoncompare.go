// Package jsoncompare emits a line only when a specified JSON field's value
// changes relative to the previous line for a given key. Unlike jsondiff,
// jsoncompare captures the old and new values and injects them into the
// output line as metadata fields, making it easy to trace value transitions.
package jsoncompare

import (
	"encoding/json"
	"sync"
)

// Comparer watches one or more JSON fields and, when a value changes, injects
// "_prev_<field>" and "_next_<field>" annotations into the emitted line.
type Comparer struct {
	fields []string
	mu     sync.Mutex
	last   map[string]map[string]any // field -> service -> last value
}

// New returns a Comparer that tracks the given fields.
// If fields is empty every line passes through unchanged.
func New(fields []string) *Comparer {
	last := make(map[string]map[string]any, len(fields))
	for _, f := range fields {
		last[f] = make(map[string]any)
	}
	return &Comparer{fields: fields, last: last}
}

// Apply returns the (possibly annotated) line and true when the line should be
// emitted, or "", false when it should be suppressed.
//
// A line is emitted when at least one tracked field has changed since the last
// call for the same service. Non-JSON lines always pass through unchanged.
func (c *Comparer) Apply(service, line string) (string, bool) {
	if len(c.fields) == 0 {
		return line, true
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, true
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	changed := false
	annotations := map[string]any{}

	for _, field := range c.fields {
		curVal, hasCur := obj[field]
		prevVal, hasPrev := c.last[field][service]

		if !hasCur {
			continue
		}

		if !hasPrev || !jsonEqual(prevVal, curVal) {
			changed = true
			if hasPrev {
				annotations["_prev_"+field] = prevVal
			}
			annotations["_next_"+field] = curVal
			c.last[field][service] = curVal
		}
	}

	if !changed {
		return "", false
	}

	for k, v := range annotations {
		obj[k] = v
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return line, true
	}
	return string(b), true
}

func jsonEqual(a, b any) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}
