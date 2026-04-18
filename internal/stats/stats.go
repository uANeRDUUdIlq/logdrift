// Package stats tracks per-service log line counters and match rates.
package stats

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"sync/atomic"
)

// Counter holds atomic counters for a single service.
type Counter struct {
	Total   atomic.Int64
	Matched atomic.Int64
}

// Tracker maintains counters for multiple services.
type Tracker struct {
	mu       sync.RWMutex
	services map[string]*Counter
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{services: make(map[string]*Counter)}
}

// counter returns (creating if needed) the Counter for a service.
func (t *Tracker) counter(service string) *Counter {
	t.mu.RLock()
	c, ok := t.services[service]
	t.mu.RUnlock()
	if ok {
		return c
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if c, ok = t.services[service]; ok {
		return c
	}
	c = &Counter{}
	t.services[service] = c
	return c
}

// RecordTotal increments the total line count for service.
func (t *Tracker) RecordTotal(service string) {
	t.counter(service).Total.Add(1)
}

// RecordMatch increments the matched line count for service.
func (t *Tracker) RecordMatch(service string) {
	t.counter(service).Matched.Add(1)
}

// Print writes a summary table to w.
func (t *Tracker) Print(w io.Writer) {
	t.mu.RLock()
	names := make([]string, 0, len(t.services))
	for k := range t.services {
		names = append(names, k)
	}
	t.mu.RUnlock()

	sort.Strings(names)
	fmt.Fprintf(w, "%-20s %10s %10s\n", "SERVICE", "TOTAL", "MATCHED")
	for _, name := range names {
		c := t.counter(name)
		fmt.Fprintf(w, "%-20s %10d %10d\n", name, c.Total.Load(), c.Matched.Load())
	}
}
