// Package jsonsorter provides a processor that sorts the keys of a JSON log
// line alphabetically before forwarding it downstream. This makes log output
// deterministic and easier to diff across services.
package jsonsorter

import (
	"encoding/json"
	"sort"
)

// Sorter sorts the top-level keys of a JSON object alphabetically.
type Sorter struct {
	descending bool
}

// New returns a Sorter. When descending is true keys are sorted Z→A.
func New(descending bool) *Sorter {
	return &Sorter{descending: descending}
}

// Apply returns the input line with its JSON keys sorted. Non-JSON lines are
// returned unchanged.
func (s *Sorter) Apply(line string) string {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return line
	}

	keys := make([]string, 0, len(raw))
	for k := range raw {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if s.descending {
			return keys[i] > keys[j]
		}
		return keys[i] < keys[j]
	})

	// Encode manually so key order is preserved.
	out := make([]byte, 0, len(line)+2)
	out = append(out, '{')
	for idx, k := range keys {
		keyBytes, err := json.Marshal(k)
		if err != nil {
			return line
		}
		out = append(out, keyBytes...)
		out = append(out, ':')
		out = append(out, raw[k]...)
		if idx < len(keys)-1 {
			out = append(out, ',')
		}
	}
	out = append(out, '}')
	return string(out)
}
