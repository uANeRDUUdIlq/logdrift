// Package writeto provides an output sink that writes rendered log lines
// to one or more io.Writer targets (e.g. stdout, a file).
package writeto

import (
	"fmt"
	"io"
	"sync"
)

// Sink fans rendered lines out to one or more writers.
type Sink struct {
	mu      sync.Mutex
	writers []io.Writer
}

// New returns a Sink that writes to each of the supplied writers.
// At least one writer must be provided.
func New(w ...io.Writer) *Sink {
	if len(w) == 0 {
		panic("writeto: at least one writer required")
	}
	return &Sink{writers: w}
}

// Write sends line to every registered writer.
// Errors from individual writers are collected and returned as a single error.
func (s *Sink) Write(line string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var first error
	for _, w := range s.writers {
		if _, err := fmt.Fprintln(w, line); err != nil && first == nil {
			first = err
		}
	}
	return first
}

// Add appends an additional writer at runtime.
func (s *Sink) Add(w io.Writer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writers = append(s.writers, w)
}
