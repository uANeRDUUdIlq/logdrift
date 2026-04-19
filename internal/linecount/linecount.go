// Package linecount tracks the number of lines seen and emitted per service.
package linecount

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

// Counter tracks seen and emitted line counts per service.
type Counter struct {
	mu      sync.Mutex
	seen    map[string]int
	emitted map[string]int
}

// New returns an initialised Counter.
func New() *Counter {
	return &Counter{
		seen:    make(map[string]int),
		emitted: make(map[string]int),
	}
}

// RecordSeen increments the seen count for service.
func (c *Counter) RecordSeen(service string) {
	c.mu.Lock()
	c.seen[service]++
	c.mu.Unlock()
}

// RecordEmitted increments the emitted count for service.
func (c *Counter) RecordEmitted(service string) {
	c.mu.Lock()
	c.emitted[service]++
	c.mu.Unlock()
}

// Snapshot returns a copy of seen and emitted maps.
func (c *Counter) Snapshot() (seen, emitted map[string]int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	seen = make(map[string]int, len(c.seen))
	emitted = make(map[string]int, len(c.emitted))
	for k, v := range c.seen {
		seen[k] = v
	}
	for k, v := range c.emitted {
		emitted[k] = v
	}
	return
}

// Print writes a summary table to w.
func (c *Counter) Print(w io.Writer) {
	seen, emitted := c.Snapshot()
	services := make([]string, 0, len(seen))
	for k := range seen {
		services = append(services, k)
	}
	sort.Strings(services)
	fmt.Fprintf(w, "%-20s %8s %8s\n", "SERVICE", "SEEN", "EMITTED")
	for _, s := range services {
		fmt.Fprintf(w, "%-20s %8d %8d\n", s, seen[s], emitted[s])
	}
}
