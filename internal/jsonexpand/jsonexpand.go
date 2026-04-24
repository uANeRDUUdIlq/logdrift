// Package jsonexpand expands dot-notation string values in a JSON log line
// into nested JSON objects. For example, a field "a.b.c" with value "v"
// becomes {"a":{"b":{"c":"v"}}} and is merged back into the top-level object.
package jsonexpand

import (
	"encoding/json"
	"strings"
)

// Expander expands dot-notation keys found in selected fields into nested objects.
type Expander struct {
	fields []string // if empty, expand all string-valued keys
}

// New returns an Expander that will expand the given fields. Pass an empty
// slice to expand every string-valued key whose name contains a dot.
func New(fields []string) *Expander {
	return &Expander{fields: fields}
}

// Apply parses line as JSON, expands configured fields, and returns the
// re-serialised line. Non-JSON lines are returned unchanged.
func (e *Expander) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	result := make(map[string]interface{})
	for k, v := range obj {
		if e.shouldExpand(k) {
			if s, ok := v.(string); ok {
				nested := buildNested(k, s)
				mergeInto(result, nested)
				continue
			}
		}
		result[k] = v
	}

	b, err := json.Marshal(result)
	if err != nil {
		return line
	}
	return string(b)
}

func (e *Expander) shouldExpand(key string) bool {
	if len(e.fields) == 0 {
		return strings.Contains(key, ".")
	}
	for _, f := range e.fields {
		if f == key {
			return true
		}
	}
	return false
}

// buildNested converts "a.b.c" + value into map["a"]map["b"]map["c"]value.
func buildNested(key, value string) map[string]interface{} {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) == 1 {
		return map[string]interface{}{key: value}
	}
	return map[string]interface{}{parts[0]: buildNested(parts[1], value)}
}

// mergeInto deep-merges src into dst.
func mergeInto(dst, src map[string]interface{}) {
	for k, sv := range src {
		if dv, ok := dst[k]; ok {
			if dMap, ok := dv.(map[string]interface{}); ok {
				if sMap, ok := sv.(map[string]interface{}); ok {
					mergeInto(dMap, sMap)
					continue
				}
			}
		}
		dst[k] = sv
	}
}
