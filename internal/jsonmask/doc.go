// Package jsonmask provides value-pattern masking for structured JSON log lines.
//
// Unlike field-name redaction, jsonmask inspects the string values of every
// top-level JSON key and replaces sub-strings that match a configured regular
// expression with a mask token (default: "[MASKED]").
//
// Typical use cases include scrubbing API tokens, credit-card numbers, social
// security numbers, or any other sensitive pattern that may appear as part of
// a field value rather than being identifiable by field name alone.
//
// Example:
//
//	rules := []jsonmask.Rule{
//		{Pattern: regexp.MustCompile(`tok_[a-z0-9]+`), Mask: "[TOKEN]"},
//	}
//	m := jsonmask.New(rules)
//	masked := m.Apply(`{"token":"tok_abc123","user":"alice"}`)
//	// masked => {"token":"[TOKEN]","user":"alice"}
package jsonmask
