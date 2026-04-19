// Package retrywriter provides a writer that retries failed writes
// up to a configurable number of attempts before returning an error.
package retrywriter

import (
	"fmt"
	"io"
	"time"
)

// Writer wraps an io.Writer and retries failed writes.
type Writer struct {
	w       io.Writer
	attempts int
	delay    time.Duration
}

// New returns a Writer that retries up to attempts times, waiting delay
// between each attempt. attempts must be >= 1; values below 1 are clamped to 1.
func New(w io.Writer, attempts int, delay time.Duration) *Writer {
	if attempts < 1 {
		attempts = 1
	}
	return &Writer{w: w, attempts: attempts, delay: delay}
}

// Write attempts to write p to the underlying writer, retrying on error.
func (rw *Writer) Write(p []byte) (int, error) {
	var lastErr error
	for i := 0; i < rw.attempts; i++ {
		n, err := rw.w.Write(p)
		if err == nil {
			return n, nil
		}
		lastErr = err
		if i < rw.attempts-1 && rw.delay > 0 {
			time.Sleep(rw.delay)
		}
	}
	return 0, fmt.Errorf("retrywriter: all %d attempts failed: %w", rw.attempts, lastErr)
}
