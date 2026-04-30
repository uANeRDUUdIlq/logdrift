// Package jsonschema validates log lines against a simple field-type schema.
// Lines that do not conform are either dropped or annotated with a validation
// error field, depending on configuration.
package jsonschema

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldType represents the expected JSON type for a field.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeNumber  FieldType = "number"
	TypeBoolean FieldType = "boolean"
	TypeObject  FieldType = "object"
	TypeArray   FieldType = "array"
)

// Rule pairs a field name with its required type.
type Rule struct {
	Field string
	Type  FieldType
}

// Validator checks JSON log lines against a set of type rules.
type Validator struct {
	rules      []Rule
	dropOnFail bool
	errorField string
}

// New creates a Validator. If dropOnFail is true, non-conforming lines are
// dropped; otherwise the original line is passed through with an optional
// errorField annotation (when errorField is non-empty).
func New(rules []Rule, dropOnFail bool, errorField string) *Validator {
	return &Validator{rules: rules, dropOnFail: dropOnFail, errorField: errorField}
}

// Apply validates line against the configured rules. It returns the
// (possibly annotated) line and true, or an empty string and false when the
// line should be dropped.
func (v *Validator) Apply(line string) (string, bool) {
	if len(v.rules) == 0 {
		return line, true
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, true // non-JSON passes through unchanged
	}
	var violations []string
	for _, r := range v.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		if !matchesType(val, r.Type) {
			violations = append(violations, fmt.Sprintf("%s must be %s", r.Field, r.Type))
		}
	}
	if len(violations) == 0 {
		return line, true
	}
	if v.dropOnFail {
		return "", false
	}
	if v.errorField != "" {
		obj[v.errorField] = strings.Join(violations, "; ")
		b, err := json.Marshal(obj)
		if err == nil {
			return string(b), true
		}
	}
	return line, true
}

func matchesType(val any, t FieldType) bool {
	switch t {
	case TypeString:
		_, ok := val.(string)
		return ok
	case TypeNumber:
		_, ok := val.(float64)
		return ok
	case TypeBoolean:
		_, ok := val.(bool)
		return ok
	case TypeObject:
		_, ok := val.(map[string]any)
		return ok
	case TypeArray:
		_, ok := val.([]any)
		return ok
	}
	return false
}
