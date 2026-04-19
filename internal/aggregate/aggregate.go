// Package aggregate counts occurrences of a JSON field value across log lines.
package aggregate

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"sync"
)

// Aggregator counts log lines grouped by a JSON field value.
type Aggregator struct {
	mu     sync.Mutex
	field  string
	counts map[string]int
}

// New creates an Aggregator that groups by the given JSON field.
func New(field string) *Aggregator {
	return &Aggregator{
		field:  field,
		counts: make(map[string]int),
	}
}

// Record parses line as JSON and increments the counter for the field value.
// Non-JSON lines or lines missing the field are counted under "<unknown>".
func (a *Aggregator) Record(line string) {
	key := "<unknown>"
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err == nil {
		if v, ok := m[a.field]; ok {
			key = fmt.Sprintf("%v", v)
		}
	}
	a.mu.Lock()
	a.counts[key]++
	a.mu.Unlock()
}

// Counts returns a snapshot of the current counts.
func (a *Aggregator) Counts() map[string]int {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make(map[string]int, len(a.counts))
	for k, v := range a.counts {
		out[k] = v
	}
	return out
}

// Print writes a sorted summary table to w.
func (a *Aggregator) Print(w io.Writer) {
	counts := a.Counts()
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Fprintf(w, "%-30s %s\n", a.field, "count")
	fmt.Fprintf(w, "%-30s %s\n", "-----", "-----")
	for _, k := range keys {
		fmt.Fprintf(w, "%-30s %d\n", k, counts[k])
	}
}
