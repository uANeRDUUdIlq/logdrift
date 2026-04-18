package multiplexer

import (
	"fmt"
	"time"

	"github.com/user/logdrift/internal/tailer"
)

// ServiceConfig describes a single service to tail.
type ServiceConfig struct {
	Name     string
	FilePath string
	Poll     time.Duration
}

// Build creates a Multiplexer and the underlying tailers from a slice of
// ServiceConfig values. It returns an error if any tailer cannot be created.
// Callers must invoke Run on the returned Multiplexer with the returned line
// channels to begin receiving entries.
//
// If two ServiceConfig entries share the same Name, the last one wins.
func Build(cfgs []ServiceConfig) (*Multiplexer, map[string]<-chan string, error) {
	if len(cfgs) == 0 {
		return nil, nil, fmt.Errorf("multiplexer: at least one service config is required")
	}

	tailers := make(map[string]*tailer.Tailer, len(cfgs))
	lines := make(map[string]<-chan string, len(cfgs))

	for _, cfg := range cfgs {
		if cfg.Name == "" {
			return nil, nil, fmt.Errorf("multiplexer: service name must not be empty")
		}
		t, err := tailer.New(cfg.FilePath, cfg.Poll)
		if err != nil {
			return nil, nil, fmt.Errorf("multiplexer: service %q: %w", cfg.Name, err)
		}
		tailers[cfg.Name] = t
		lines[cfg.Name] = t.Lines()
	}

	m := New(tailers)
	return m, lines, nil
}
