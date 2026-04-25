package jsonlimit

import "fmt"

// Config holds the builder configuration for a Limiter.
type Config struct {
	// DefaultMax is the maximum number of array elements allowed for any
	// array field not explicitly listed in Fields. Zero disables the default.
	DefaultMax int

	// Fields maps specific JSON field names to their individual max lengths.
	// A value of zero or negative for a field is ignored.
	Fields map[string]int
}

// Build validates cfg and returns a ready-to-use Limiter.
func Build(cfg Config) (*Limiter, error) {
	if cfg.DefaultMax < 0 {
		return nil, fmt.Errorf("jsonlimit: default_max must be >= 0, got %d", cfg.DefaultMax)
	}
	for field, max := range cfg.Fields {
		if max < 0 {
			return nil, fmt.Errorf("jsonlimit: field %q max must be >= 0, got %d", field, max)
		}
	}
	return New(cfg.DefaultMax, cfg.Fields), nil
}
