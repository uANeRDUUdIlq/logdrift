package jsonschema

import (
	"errors"
	"fmt"
)

// Config holds the builder parameters for a Validator.
type Config struct {
	// Rules is the list of field-type constraints.
	Rules []Rule
	// DropOnFail causes non-conforming lines to be dropped when true.
	// When false the line passes through, optionally annotated.
	DropOnFail bool
	// ErrorField is the JSON key injected with a violation summary when
	// DropOnFail is false and the field name is non-empty.
	ErrorField string
}

// Build validates cfg and returns a ready-to-use Validator.
func Build(cfg Config) (*Validator, error) {
	validTypes := map[FieldType]bool{
		TypeString:  true,
		TypeNumber:  true,
		TypeBoolean: true,
		TypeObject:  true,
		TypeArray:   true,
	}
	for i, r := range cfg.Rules {
		if r.Field == "" {
			return nil, fmt.Errorf("rule[%d]: field must not be empty", i)
		}
		if !validTypes[r.Type] {
			return nil, fmt.Errorf("rule[%d]: unknown type %q", i, r.Type)
		}
	}
	if !cfg.DropOnFail && cfg.ErrorField == "" {
		// Perfectly valid — violations are silently ignored.
		_ = errors.New // keep import if needed later
	}
	return New(cfg.Rules, cfg.DropOnFail, cfg.ErrorField), nil
}
