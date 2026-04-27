// Package jsonpivot rotates a repeated field value into a top-level key.
//
// Given a line such as:
//
//	{"metrics":[{"name":"cpu","value":0.42},{"name":"mem","value":0.71}]}
//
// Pivoting on field "metrics", key "name", value "value" produces:
//
//	{"cpu":0.42,"mem":0.71}
package jsonpivot

import (
	"encoding/json"
)

// Pivoter rotates array entries into a flat map keyed by a named field.
type Pivoter struct {
	array string
	keyField string
	valField string
}

// New returns a Pivoter that reads arr (an array field), uses keyField as the
// map key and valField as the map value. Returns an error when any argument is
// empty.
func New(array, keyField, valField string) (*Pivoter, error) {
	if array == "" || keyField == "" || valField == "" {
		return nil, errEmptyField
	}
	return &Pivoter{array: array, keyField: keyField, valField: valField}, nil
}

// Apply transforms line and returns the pivoted JSON. Non-JSON lines and lines
// missing the target array are returned unchanged.
func (p *Pivoter) Apply(line string) string {
	var root map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &root); err != nil {
		return line
	}

	rawArr, ok := root[p.array]
	if !ok {
		return line
	}

	var items []map[string]json.RawMessage
	if err := json.Unmarshal(rawArr, &items); err != nil {
		return line
	}

	result := make(map[string]json.RawMessage, len(root)-1+len(items))
	for k, v := range root {
		if k != p.array {
			result[k] = v
		}
	}

	for _, item := range items {
		keyRaw, hasKey := item[p.keyField]
		valRaw, hasVal := item[p.valField]
		if !hasKey || !hasVal {
			continue
		}
		var keyStr string
		if err := json.Unmarshal(keyRaw, &keyStr); err != nil {
			continue
		}
		result[keyStr] = valRaw
	}

	out, err := json.Marshal(result)
	if err != nil {
		return line
	}
	return string(out)
}
