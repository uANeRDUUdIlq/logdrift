package jsonpatch

import "fmt"

// Config holds the configuration for building a Patcher from external config.
type Config struct {
	// Fields maps JSON field names to the string values that should be injected
	// or overwritten in every log line.
	Fields map[string]string `yaml:"fields"`
}

// Build validates cfg and returns a ready-to-use Patcher.
// An empty Fields map is allowed and produces a no-op Patcher.
func Build(cfg Config) (*Patcher, error) {
	for k := range cfg.Fields {
		if k == "" {
			return nil, fmt.Errorf("jsonpatch: field key must not be empty")
		}
	}
	return New(cfg.Fields)
}
