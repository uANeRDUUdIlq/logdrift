// Package jsongroup groups consecutive JSON log lines by a shared field value,
// emitting a summary line once the group changes or is flushed.
package jsongroup

import (
	"encoding/json"
	"fmt"
)

// Grouper accumulates consecutive lines that share the same value for a
// configured key field and emits a single merged summary when the group ends.
type Grouper struct {
	field    string
	countKey string
	current  string
	buf      []map[string]any
}

// New returns a Grouper that groups lines by field. countKey is the JSON key
// used to record how many lines were merged in the emitted summary.
func New(field, countKey string) *Grouper {
	if countKey == "" {
		countKey = "_count"
	}
	return &Grouper{field: field, countKey: countKey}
}

// Push accepts a raw log line. If the group key changes, the previous group is
// flushed and the new line starts a fresh group. Returns zero or more summary
// lines ready to be emitted.
func (g *Grouper) Push(line string) []string {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		// Non-JSON: flush pending group and pass through unchanged.
		out := g.Flush()
		return append(out, line)
	}

	val := ""
	if v, ok := obj[g.field]; ok {
		val = fmt.Sprintf("%v", v)
	}

	var out []string
	if val != g.current && len(g.buf) > 0 {
		out = g.Flush()
	}
	g.current = val
	g.buf = append(g.buf, obj)
	return out
}

// Flush finalises the current group and returns a summary line (if any lines
// were buffered). The buffer is cleared after the call.
func (g *Grouper) Flush() []string {
	if len(g.buf) == 0 {
		return nil
	}
	// Merge all buffered objects; last writer wins for duplicate keys.
	merged := make(map[string]any)
	for _, obj := range g.buf {
		for k, v := range obj {
			merged[k] = v
		}
	}
	merged[g.countKey] = len(g.buf)
	g.buf = g.buf[:0]
	g.current = ""

	b, err := json.Marshal(merged)
	if err != nil {
		return nil
	}
	return []string{string(b)}
}
