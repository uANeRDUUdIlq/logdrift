// Package jsonreorder reorders JSON object keys into a deterministic
// user-specified order, placing any unspecified keys at the end.
package jsonreorder

import (
	"encoding/json"
)

// Reorderer moves a fixed list of keys to the front of each JSON log line.
type Reorderer struct {
	keys []string
	index map[string]int
}

// New returns a Reorderer that will place keys in the given order at the
// front of every JSON object it processes. Keys not listed are appended
// in their original iteration order.
func New(keys []string) *Reorderer {
	idx := make(map[string]int, len(keys))
	for i, k := range keys {
		idx[k] = i
	}
	return &Reorderer{keys: keys, index: idx}
}

// Apply reorders the keys in line according to the configured priority list.
// Non-JSON lines are returned unchanged.
func (r *Reorderer) Apply(line string) string {
	if len(r.keys) == 0 {
		return line
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return line
	}

	// Build ordered output: priority keys first, remainder after.
	type kv struct {
		key string
		val json.RawMessage
	}

	ordered := make([]kv, 0, len(raw))
	for _, k := range r.keys {
		if v, ok := raw[k]; ok {
			ordered = append(ordered, kv{k, v})
		}
	}
	for k, v := range raw {
		if _, prioritised := r.index[k]; !prioritised {
			ordered = append(ordered, kv{k, v})
		}
	}

	out, err := marshalOrdered(ordered)
	if err != nil {
		return line
	}
	return string(out)
}

func marshalOrdered(pairs []struct {
	key string
	val json.RawMessage
}) ([]byte, error) {
	buf := []byte{'{'}
	for i, p := range pairs {
		if i > 0 {
			buf = append(buf, ',')
		}
		k, err := json.Marshal(p.key)
		if err != nil {
			return nil, err
		}
		buf = append(buf, k...)
		buf = append(buf, ':')
		buf = append(buf, p.val...)
	}
	buf = append(buf, '}')
	return buf, nil
}
