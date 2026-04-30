// Package jsonclip trims a JSON log line so that it contains only the first N
// top-level keys, in document order.  Keys beyond the limit are silently
// dropped.  Non-JSON lines are passed through unchanged.
package jsonclip

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Clipper drops top-level keys from a JSON object beyond a configured limit.
type Clipper struct {
	maxKeys int
}

// New returns a Clipper that keeps at most maxKeys top-level keys.
// If maxKeys is zero or negative the Clipper is a no-op.
func New(maxKeys int) *Clipper {
	return &Clipper{maxKeys: maxKeys}
}

// Apply returns a version of line that contains at most maxKeys top-level keys.
// The original line is returned unchanged when:
//   - maxKeys <= 0 (disabled)
//   - line is not a JSON object
//   - the object has fewer or equal keys than the limit
func (c *Clipper) Apply(line string) string {
	if c.maxKeys <= 0 {
		return line
	}

	var probe map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &probe); err != nil {
		return line
	}
	if len(probe) <= c.maxKeys {
		return line
	}

	// Walk the token stream to preserve insertion order.
	dec := json.NewDecoder(strings.NewReader(line))
	dec.Token() // consume '{'

	type pair struct {
		key string
		val json.RawMessage
	}
	var pairs []pair
	for dec.More() && len(pairs) < c.maxKeys {
		keyTok, err := dec.Token()
		if err != nil {
			return line
		}
		key := fmt.Sprintf("%v", keyTok)
		var val json.RawMessage
		if err := dec.Decode(&val); err != nil {
			return line
		}
		pairs = append(pairs, pair{key, val})
	}

	var sb strings.Builder
	sb.WriteByte('{')
	for i, p := range pairs {
		if i > 0 {
			sb.WriteByte(',')
		}
		keyBytes, _ := json.Marshal(p.key)
		sb.Write(keyBytes)
		sb.WriteByte(':')
		sb.Write(p.val)
	}
	sb.WriteByte('}')
	return sb.String()
}
