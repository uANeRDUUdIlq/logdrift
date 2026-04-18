// Package transform applies field-level mutations to JSON log lines,
// such as renaming or adding static fields before output.
package transform

import (
	"encoding/json"
)

// Rule describes a single transformation.
type Rule struct {
	Rename map[string]string `yaml:"rename"` // old_key -> new_key
	Add    map[string]string `yaml:"add"`    // key -> static value
}

// Transformer mutates JSON log lines according to configured rules.
type Transformer struct {
	rules Rule
}

// New returns a Transformer configured with the given Rule.
func New(r Rule) *Transformer {
	return &Transformer{rules: r}
}

// Apply returns the transformed line. Non-JSON lines are returned unchanged.
func (t *Transformer) Apply(line string) string {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for oldKey, newKey := range t.rules.Rename {
		if val, ok := obj[oldKey]; ok {
			obj[newKey] = val
			delete(obj, oldKey)
		}
	}

	for k, v := range t.rules.Add {
		obj[k] = v
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
