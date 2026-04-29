// Package jsonroute routes log lines to different writers based on
// field-value matching rules. Each rule specifies a JSON field, an
// expected value, and a destination writer. Lines that do not match
// any rule are forwarded to an optional fallback writer.
package jsonroute

import (
	"encoding/json"
	"io"
)

// Rule describes a single routing rule.
type Rule struct {
	// Field is the top-level JSON key to inspect.
	Field string
	// Value is the string representation to match against.
	Value string
	// Writer is the destination for matching lines.
	Writer io.Writer
}

// Router routes lines to writers based on JSON field values.
type Router struct {
	rules    []Rule
	fallback io.Writer
}

// New creates a Router with the given rules and an optional fallback
// writer. Pass nil as fallback to silently drop unmatched lines.
func New(rules []Rule, fallback io.Writer) (*Router, error) {
	for i, r := range rules {
		if r.Field == "" {
			return nil, &ErrEmptyField{Index: i}
		}
		if r.Writer == nil {
			return nil, &ErrNilWriter{Index: i}
		}
	}
	return &Router{rules: rules, fallback: fallback}, nil
}

// Route inspects line and writes it to the first matching writer.
// If no rule matches and a fallback writer is configured, line is
// written there instead. Route returns any write error encountered.
func (r *Router) Route(line string) error {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return r.writeFallback(line)
	}

	for _, rule := range r.rules {
		v, ok := obj[rule.Field]
		if !ok {
			continue
		}
		var sv string
		switch val := v.(type) {
		case string:
			sv = val
		default:
			b, _ := json.Marshal(v)
			sv = string(b)
		}
		if sv == rule.Value {
			_, err := io.WriteString(rule.Writer, line+"\n")
			return err
		}
	}
	return r.writeFallback(line)
}

func (r *Router) writeFallback(line string) error {
	if r.fallback == nil {
		return nil
	}
	_, err := io.WriteString(r.fallback, line+"\n")
	return err
}
