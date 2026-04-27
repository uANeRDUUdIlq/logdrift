// Package jsoncoerce provides a pipeline stage that coerces JSON field
// values to a target type (string, number, bool). Non-JSON lines are passed
// through unchanged. Fields that cannot be coerced are left as-is.
package jsoncoerce

import (
	"encoding/json"
	"strconv"
	"strings"
)

// TargetType represents the destination type for a coercion rule.
type TargetType string

const (
	TypeString TargetType = "string"
	TypeNumber TargetType = "number"
	TypeBool   TargetType = "bool"
)

// Rule maps a JSON field key to the desired target type.
type Rule struct {
	Field  string
	Target TargetType
}

// Coercer applies type coercion rules to JSON log lines.
type Coercer struct {
	rules []Rule
}

// New returns a Coercer that applies the given rules.
func New(rules []Rule) *Coercer {
	return &Coercer{rules: rules}
}

// Apply coerces fields in a JSON line according to the configured rules.
// Non-JSON input is returned unchanged.
func (c *Coercer) Apply(line string) string {
	if len(c.rules) == 0 {
		return line
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	changed := false
	for _, r := range c.rules {
		v, ok := obj[r.Field]
		if !ok {
			continue
		}
		coerced, ok := coerce(v, r.Target)
		if !ok {
			continue
		}
		obj[r.Field] = coerced
		changed = true
	}
	if !changed {
		return line
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(b)
}

func coerce(v interface{}, t TargetType) (interface{}, bool) {
	switch t {
	case TypeString:
		switch val := v.(type) {
		case string:
			return val, true
		case float64:
			return strconv.FormatFloat(val, 'f', -1, 64), true
		case bool:
			return strconv.FormatBool(val), true
		}
	case TypeNumber:
		switch val := v.(type) {
		case float64:
			return val, true
		case string:
			f, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
			if err == nil {
				return f, true
			}
		case bool:
			if val {
				return float64(1), true
			}
			return float64(0), true
		}
	case TypeBool:
		switch val := v.(type) {
		case bool:
			return val, true
		case string:
			b, err := strconv.ParseBool(strings.TrimSpace(val))
			if err == nil {
				return b, true
			}
		case float64:
			return val != 0, true
		}
	}
	return nil, false
}
