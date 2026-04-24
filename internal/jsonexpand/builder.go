package jsonexpand

import "fmt"

// Config holds the configuration used by Build to construct an Expander.
type Config struct {
	// Fields lists the dot-notation keys to expand. Leave nil or empty to
	// automatically expand every string-valued key whose name contains a dot.
	Fields []string `yaml:"fields"`

	// MaxDepth limits how many levels of nesting are created when expanding a
	// key. Zero or negative means unlimited.
	MaxDepth int `yaml:"max_depth"`
}

// Build validates cfg and returns a ready-to-use Expander.
func Build(cfg Config) (*Expander, error) {
	if cfg.MaxDepth < 0 {
		return nil, fmt.Errorf("jsonexpand: max_depth must be >= 0, got %d", cfg.MaxDepth)
	}
	return New(cfg.Fields), nil
}
