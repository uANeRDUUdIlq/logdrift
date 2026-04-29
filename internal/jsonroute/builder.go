package jsonroute

import (
	"fmt"
	"io"
)

// BuilderRule is the configuration representation of a routing rule,
// suitable for use in YAML/JSON config files.
type BuilderRule struct {
	Field  string
	Value  string
	Writer io.Writer
}

// Build constructs a Router from a slice of BuilderRules and a
// fallback writer. It validates each rule and returns an error if any
// rule is misconfigured.
func Build(rules []BuilderRule, fallback io.Writer) (*Router, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("jsonroute: at least one rule is required")
	}
	routing := make([]Rule, len(rules))
	for i, br := range rules {
		routing[i] = Rule{
			Field:  br.Field,
			Value:  br.Value,
			Writer: br.Writer,
		}
	}
	return New(routing, fallback)
}
