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
func Build(cfgs []ServiceConfig) (*Multiplexer, map[string]<-chan string, error) {
	tailers := make(map[string]*tailer.Tailer, len(cfgs))
	lines := make(map[string]<-chan string, len(cfgs))

	for _, cfg := range cfgs {
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
