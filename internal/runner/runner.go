// Package runner wires together config, tailers, multiplexer, filter, and
// formatter into a single blocking Run call.
package runner

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/example/logdrift/internal/config"
	"github.com/example/logdrift/internal/filter"
	"github.com/example/logdrift/internal/formatter"
	"github.com/example/logdrift/internal/multiplexer"
)

// Options holds CLI-level overrides.
type Options struct {
	ConfigPath string
	RawOutput  bool
	Filters    []string // extra filter expressions from flags
	Output     io.Writer
}

// Run loads configuration and streams log lines until ctx is cancelled.
func Run(ctx context.Context, opts Options) error {
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	f, err := filter.New(append(cfg.Filters, opts.Filters...))
	if err != nil {
		return fmt.Errorf("build filter: %w", err)
	}

	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	formatters := make([]*formatter.Formatter, len(cfg.Services))
	for i, svc := range cfg.Services {
		formatters[i] = formatter.New(svc.Name, i)
	}

	ch, err := multiplexer.Build(ctx, cfg)
	if err != nil {
		return fmt.Errorf("build multiplexer: %w", err)
	}

	// index formatters by service name for O(1) lookup
	fmtByService := make(map[string]*formatter.Formatter, len(formatters))
	for i, svc := range cfg.Services {
		fmtByService[svc.Name] = formatters[i]
	}

	for entry := range ch {
		if !f.Match(entry.Line) {
			continue
		}
		fmt := fmtByService[entry.Service]
		if fmt == nil {
			continue
		}
		line := fmt.Render(entry.Line, !opts.RawOutput)
		_, _ = io.WriteString(out, line+"\n")
	}
	return nil
}
