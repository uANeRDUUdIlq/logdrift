package jsonflatten

import "fmt"

// Config holds the configuration used by Build to construct a Flattener.
type Config struct {
	// Separator is the string inserted between parent and child key segments.
	// Defaults to "." when empty.
	Separator string `yaml:"separator"`

	// Prefix is prepended to every top-level key in the flattened output.
	// May be empty.
	Prefix string `yaml:"prefix"`

	// MaxDepth, when > 0, limits how many levels deep the flattener will
	// recurse. Objects beyond MaxDepth are serialised as JSON strings.
	MaxDepth int `yaml:"max_depth"`
}

// Build validates cfg and returns a ready-to-use Flattener.
func Build(cfg Config) (*Flattener, error) {
	if cfg.MaxDepth < 0 {
		return nil, fmt.Errorf("jsonflatten: max_depth must be >= 0, got %d", cfg.MaxDepth)
	}
	sep := cfg.Separator
	if sep == "" {
		sep = "."
	}
	return New(sep, cfg.Prefix), nil
}
