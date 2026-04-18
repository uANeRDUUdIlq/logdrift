// Package runner wires all logdrift components and drives the main pipeline.
package runner

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/logdrift/internal/config"
	"github.com/yourorg/logdrift/internal/formatter"
	"github.com/yourorg/logdrift/internal/multiplexer"
	"github.com/yourorg/logdrift/internal/redact"
	"github.com/yourorg/logdrift/internal/stats"
	"github.com/yourorg/logdrift/internal/transform"
	"github.com/yourorg/logdrift/internal/truncate"
)

// Run loads config from cfgPath and starts the tailing pipeline, writing
// formatted output to out until interrupted or an error occurs.
func Run(cfgPath string, out io.Writer) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("runner: %w", err)
	}

	tracker := stats.New()
	trunc := truncate.New(cfg.MaxLen)

	serviceNames := make([]string, len(cfg.Services))
	for i, s := range cfg.Services {
		serviceNames[i] = s.Name
	}

	ch, err := multiplexer.Build(cfg)
	if err != nil {
		return fmt.Errorf("runner: %w", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case entry, ok := <-ch:
			if !ok {
				tracker.Print(out)
				return nil
			}
			tracker.RecordTotal(entry.Service)

			line := entry.Line

			// per-service transform
			for _, svc := range cfg.Services {
				if svc.Name == entry.Service {
					tr := transform.New(transform.Rule{
						Rename: svc.Transform.Rename,
						Add:    svc.Transform.Add,
					})
					line = tr.Apply(line)
				}
			}

			rd := redact.New(nil, "")
			line = rd.Apply(line)
			line = trunc.Apply(line)

			fmt_ := formatter.New(serviceNames, cfg.Format)
			tracker.RecordMatch(entry.Service)
			fmt.Fprintln(out, fmt_.Render(entry.Service, line))

		case <-sig:
			tracker.Print(out)
			return nil
		}
	}
}
