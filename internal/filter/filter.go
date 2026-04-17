package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule defines a single filter condition on a JSON log field.
type Rule struct {
	Field string
	Value string
}

// Filter holds a set of rules that log entries must satisfy.
type Filter struct {
	Rules []Rule
}

// New creates a Filter from a slice of "field=value" expressions.
func New(exprs []string) (*Filter, error) {
	f := &Filter{}
	for _, expr := range exprs {
		parts := strings.SplitN(expr, "=", 2)
		if len(parts) != 2 {
			return nil, &InvalidExprError{Expr: expr}
		}
		f.Rules = append(f.Rules, Rule{Field: parts[0], Value: parts[1]})
	}
	return f, nil
}

// Match reports whether the raw JSON line satisfies all filter rules.
func (f *Filter) Match(line []byte) bool {
	if len(f.Rules) == 0 {
		return true
	}
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return false
	}
	for _, rule := range f.Rules {
		v, ok := entry[rule.Field]
		if !ok {
			return false
		}
		if !strings.EqualFold(fmt.Sprintf("%v", v), rule.Value) {
			return false
		}
	}
	return true
}

// InvalidExprError is returned when a filter expression is malformed.
type InvalidExprError struct {
	Expr string
}

func (e *InvalidExprError) Error() string {
	return "invalid filter expression: " + e.Expr
}
