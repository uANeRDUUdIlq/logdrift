package rollup

import (
	"fmt"
	"io"
	"time"
)

// Config holds configuration for building a Rollup pipeline stage.
type Config struct {
	Field  string        `yaml:"field"`
	Window time.Duration `yaml:"window"`
}

// Build constructs a Rollup from cfg and starts a goroutine that prints
// each flushed Entry to w.
func Build(cfg Config, w io.Writer) (*Rollup, error) {
	if cfg.Field == "" {
		return nil, fmt.Errorf("rollup: field must not be empty")
	}
	if cfg.Window < 0 {
		return nil, fmt.Errorf("rollup: window must be >= 0")
	}
	r := New(cfg.Field, cfg.Window)
	go func() {
		for e := range r.Entries() {
			fmt.Fprintf(w, "[rollup] service=%s %s=%s count=%d\n",
				e.Service, e.Field, e.Value, e.Count)
		}
	}()
	return r, nil
}
