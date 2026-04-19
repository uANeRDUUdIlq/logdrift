// Package labelinject merges static key-value labels into every JSON log line.
package labelinject

import (
	"encoding/json"
)

// Injector merges a fixed set of labels into JSON log lines.
type Injector struct {
	labels map[string]string
}

// New returns an Injector that will add the given labels to every line.
// If labels is empty the Injector is a no-op.
func New(labels map[string]string) *Injector {
	copy := make(map[string]string, len(labels))
	for k, v := range labels {
		copy[k] = v
	}
	return &Injector{labels: copy}
}

// Apply merges the configured labels into line.
// Non-JSON lines are returned unchanged.
// Existing fields in the line are NOT overwritten.
func (inj *Injector) Apply(line string) string {
	if len(inj.labels) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for k, v := range inj.labels {
		if _, exists := obj[k]; !exists {
			encoded, _ := json.Marshal(v)
			obj[k] = json.RawMessage(encoded)
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
