// Package teewriter duplicates matched log lines into a separate file sink,
// useful for capturing a filtered subset of logs to disk while still printing
// to stdout.
package teewriter

import (
	"fmt"
	"os"
	"sync"
)

// TeeWriter writes every line it receives to one or more backing files.
type TeeWriter struct {
	mu    sync.Mutex
	files []*os.File
}

// New opens each path for append-only writing (creating if absent) and returns
// a TeeWriter that fans out to all of them. The caller must call Close when
// done.
func New(paths []string) (*TeeWriter, error) {
	if len(paths) == 0 {
		return &TeeWriter{}, nil
	}
	files := make([]*os.File, 0, len(paths))
	for _, p := range paths {
		f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			// Close already-opened files before returning.
			for _, already := range files {
				_ = already.Close()
			}
			return nil, fmt.Errorf("teewriter: open %q: %w", p, err)
		}
		files = append(files, f)
	}
	return &TeeWriter{files: files}, nil
}

// Write appends line (with a trailing newline) to every backing file.
// Errors from individual files are joined and returned.
func (t *TeeWriter) Write(line string) error {
	if len(t.files) == 0 {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	var first error
	for _, f := range t.files {
		if _, err := fmt.Fprintln(f, line); err != nil && first == nil {
			first = fmt.Errorf("teewriter: write to %s: %w", f.Name(), err)
		}
	}
	return first
}

// Close flushes and closes all backing files.
func (t *TeeWriter) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	var first error
	for _, f := range t.files {
		if err := f.Close(); err != nil && first == nil {
			first = err
		}
	}
	t.files = nil
	return first
}
