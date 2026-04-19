package multiplexer

import (
	"time"

	"github.com/yourorg/logdrift/internal/buffer"
	"github.com/yourorg/logdrift/internal/config"
	"github.com/yourorg/logdrift/internal/tailer"
)

// ServiceBuffer pairs a service name with its recent-line buffer.
type ServiceBuffer struct {
	Service string
	Buf     *buffer.Buffer
}

// Build constructs a Multiplexer and per-service Buffers from cfg.
// Each service gets a tailer and a ring buffer of bufferSize lines.
func Build(cfg *config.Config, bufferSize int) (*Multiplexer, []ServiceBuffer, error) {
	chans := make([]<-chan Entry, 0, len(cfg.Services))
	bufs := make([]ServiceBuffer, 0, len(cfg.Services))

	for _, svc := range cfg.Services {
		poll := time.Duration(cfg.PollInterval) * time.Millisecond
		t, err := tailer.New(svc.Path, poll)
		if err != nil {
			return nil, nil, err
		}

		buf := buffer.New(bufferSize)
		bufs = append(bufs, ServiceBuffer{Service: svc.Name, Buf: buf})

		raw := t.Tail()
		entryCh := make(chan Entry)
		go func(name string, in <-chan string, out chan<- Entry, b *buffer.Buffer) {
			for line := range in {
				b.Push(line)
				out <- Entry{Service: name, Line: line}
			}
			close(out)
		}(svc.Name, raw, entryCh, buf)

		chans = append(chans, entryCh)
	}

	return New(chans), bufs, nil
}
