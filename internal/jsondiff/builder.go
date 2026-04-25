package jsondiff

import "fmt"

// Config holds the parameters needed to construct a Differ via Build.
type Config struct {
	// Fields is the list of top-level JSON keys to watch for changes.
	// An empty slice disables change detection (every line passes).
	Fields []string
}

// Build validates cfg and returns a ready-to-use Differ.
// It returns an error if any field name is an empty string.
func Build(cfg Config) (*Differ, error) {
	for i, f := range cfg.Fields {
		if f == "" {
			return nil, fmt.Errorf("jsondiff: field at index %d must not be empty", i)
		}
	}
	return New(cfg.Fields), nil
}
