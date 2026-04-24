package jsonreorder

import "fmt"

// Config holds the configuration for building a Reorderer.
type Config struct {
	// Keys is the ordered list of field names that should appear first in
	// every output JSON object. Must contain at least one entry.
	Keys []string `yaml:"keys"`
}

// Build validates cfg and returns a ready-to-use Reorderer.
func Build(cfg Config) (*Reorderer, error) {
	if len(cfg.Keys) == 0 {
		return nil, fmt.Errorf("jsonreorder: keys list must not be empty")
	}
	seen := make(map[string]struct{}, len(cfg.Keys))
	for _, k := range cfg.Keys {
		if k == "" {
			return nil, fmt.Errorf("jsonreorder: key must not be empty string")
		}
		if _, dup := seen[k]; dup {
			return nil, fmt.Errorf("jsonreorder: duplicate key %q", k)
		}
		seen[k] = struct{}{}
	}
	return New(cfg.Keys), nil
}
