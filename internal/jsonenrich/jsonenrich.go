// Package jsonenrich attaches static or derived metadata fields to JSON log lines.
//
// It is useful for stamping every log line with environment, region, or
// host information before forwarding to a sink.
package jsonenrich

import (
	"encoding/json"
	"os"
	"strings"
)

// Enricher appends fixed key/value pairs to each JSON log line.
type Enricher struct {
	fields map[string]string
}

// Option configures an Enricher.
type Option func(*Enricher)

// WithField adds a static key/value pair to every line.
func WithField(key, value string) Option {
	return func(e *Enricher) {
		e.fields[key] = value
	}
}

// WithHostname adds the system hostname under the given key.
// If the hostname cannot be resolved the field is omitted.
func WithHostname(key string) Option {
	return func(e *Enricher) {
		if h, err := os.Hostname(); err == nil {
			e.fields[key] = h
		}
	}
}

// New returns an Enricher configured with the supplied options.
// Fields are only injected when the key is not already present in the line.
func New(opts ...Option) *Enricher {
	e := &Enricher{fields: make(map[string]string)}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Apply injects configured fields into line and returns the modified JSON.
// Non-JSON lines and lines with no configured fields are returned unchanged.
func (e *Enricher) Apply(line string) string {
	if len(e.fields) == 0 {
		return line
	}
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "{") {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(trimmed), &obj); err != nil {
		return line
	}

	changed := false
	for k, v := range e.fields {
		if _, exists := obj[k]; !exists {
			encoded, err := json.Marshal(v)
			if err != nil {
				continue
			}
			obj[k] = json.RawMessage(encoded)
			changed = true
		}
	}
	if !changed {
		return line
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
