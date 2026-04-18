// Package multiplexer fans-in log lines from multiple tailers into a single channel.
package multiplexer

import (
	"sync"

	"github.com/user/logdrift/internal/tailer"
)

// Entry holds a log line together with the service name it originated from.
type Entry struct {
	Service string
	Line    string
}

// Multiplexer fans-in lines from multiple tailers.
type Multiplexer struct {
	services map[string]*tailer.Tailer
	out      chan Entry
}

// New creates a Multiplexer for the given map of service-name → tailer.
func New(services map[string]*tailer.Tailer) *Multiplexer {
	return &Multiplexer{
		services: services,
		out:      make(chan Entry, 256),
	}
}

// Run starts one goroutine per tailer and merges their output into Out().
// It closes Out() once every tailer's channel is drained.
func (m *Multiplexer) Run(lines map[string]<-chan string) {
	var wg sync.WaitGroup
	for svc, ch := range lines {
		wg.Add(1)
		go func(service string, src <-chan string) {
			defer wg.Done()
			for line := range src {
				m.out <- Entry{Service: service, Line: line}
			}
		}(svc, ch)
	}
	go func() {
		wg.Wait()
		close(m.out)
	}()
}

// Out returns the merged output channel.
func (m *Multiplexer) Out() <-chan Entry {
	return m.out
}
