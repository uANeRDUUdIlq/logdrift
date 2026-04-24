// Package jsoncount counts occurrences of a specific JSON field value
// across log lines and periodically emits a summary line.
package jsoncount

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"
)

// Counter tracks value frequencies for a named JSON field.
type Counter struct {
	mu      sync.Mutex
	field   string
	counts  map[string]int
	window  time.Duration
	writer  io.Writer
	stopCh  chan struct{}
}

// New creates a Counter that watches field, flushes counts every window
// to writer, and starts its background flush goroutine.
func New(field string, window time.Duration, w io.Writer) *Counter {
	c := &Counter{
		field:  field,
		counts: make(map[string]int),
		window: window,
		writer: w,
		stopCh: make(chan struct{}),
	}
	if window > 0 {
		go c.loop()
	}
	return c
}

// Record parses line as JSON and increments the counter for the value
// found at field. Non-JSON lines and lines missing the field are ignored.
func (c *Counter) Record(line string) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return
	}
	v, ok := obj[c.field]
	if !ok {
		return
	}
	key := fmt.Sprintf("%v", v)
	c.mu.Lock()
	c.counts[key]++
	c.mu.Unlock()
}

// Flush writes a sorted summary of current counts to the writer and resets
// the internal map. It is safe to call concurrently.
func (c *Counter) Flush() {
	c.mu.Lock()
	snap := make(map[string]int, len(c.counts))
	for k, v := range c.counts {
		snap[k] = v
	}
	c.counts = make(map[string]int)
	c.mu.Unlock()

	if len(snap) == 0 {
		return
	}

	keys := make([]string, 0, len(snap))
	for k := range snap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(c.writer, "[jsoncount] field=%q\n", c.field)
	for _, k := range keys {
		fmt.Fprintf(c.writer, "  %-40s %d\n", k, snap[k])
	}
}

// Stop halts the background flush goroutine.
func (c *Counter) Stop() {
	if c.window > 0 {
		close(c.stopCh)
	}
}

func (c *Counter) loop() {
	ticker := time.NewTicker(c.window)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.Flush()
		case <-c.stopCh:
			return
		}
	}
}
