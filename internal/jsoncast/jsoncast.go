// Package jsoncast re-emits a JSON log line once for each element of a named
// array field, promoting each element's fields into the top-level object.
// Lines that are not valid JSON, or whose target field is not an array, are
// passed through unchanged.
package jsoncast

import (
	"encoding/json"
)

// Caster fans out array-valued fields into individual lines.
type Caster struct {
	field string
}

// New returns a Caster that fans out the given array field.
// field must be non-empty; if it is empty New returns an error.
func New(field string) (*Caster, error) {
	if field == "" {
		return nil, errEmptyField
	}
	return &Caster{field: field}, nil
}

// Apply fans out the named array field of line into zero or more JSON strings.
// If line is not valid JSON, or the field is absent or not an array, Apply
// returns []string{line} so the original line is preserved.
func (c *Caster) Apply(line string) []string {
	var root map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &root); err != nil {
		return []string{line}
	}

	raw, ok := root[c.field]
	if !ok {
		return []string{line}
	}

	var elems []json.RawMessage
	if err := json.Unmarshal(raw, &elems); err != nil {
		// Field exists but is not an array – pass through.
		return []string{line}
	}

	// Build a base object without the array field.
	base := make(map[string]json.RawMessage, len(root)-1)
	for k, v := range root {
		if k != c.field {
			base[k] = v
		}
	}

	results := make([]string, 0, len(elems))
	for _, elem := range elems {
		// If the element is an object, merge its keys into the base.
		var sub map[string]json.RawMessage
		merged := make(map[string]json.RawMessage, len(base))
		for k, v := range base {
			merged[k] = v
		}
		if err := json.Unmarshal(elem, &sub); err == nil {
			for k, v := range sub {
				merged[k] = v
			}
		} else {
			// Scalar element: store under the original field name.
			merged[c.field] = elem
		}
		b, err := json.Marshal(merged)
		if err != nil {
			continue
		}
		results = append(results, string(b))
	}

	if len(results) == 0 {
		return []string{line}
	}
	return results
}

var errEmptyField = fmt.Errorf("jsoncast: field must not be empty")
