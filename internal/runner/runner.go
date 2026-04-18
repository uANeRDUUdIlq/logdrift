// Package runner wires all logdrift components together and drives the main
// processing loop.
package runner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/logdrift/internal/config"
	"github.com/yourorg/logdrift/internal/dedupe"
	"github.com/yourorg/logdrift/internal/filter"
	"github.com/yourorg/logdrift/internal/formatter"
	"github.com/yourorg/logdrift/internal/highlight"
	"github.com/yourorg/logdrift/internal/multiplexer"
	"github.com/yourorg/logdrift/internal/ratelimit"
	"github.com/yourorg/logdrift/internal/redact"
	"github.com/yourorg/logdrift/internal/stats"
	"github.com/yourorg/logdrift/internal/truncate"
)

// Run loads config from cfgPath and starts the tail/filter/render loop.
func Run(cfgPath string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	f, err := filter.New(cfg.Filters)
	if err != nil {
		return fmt.Errorf("filter: %w", err)
	}

	hl := highlight.New(cfg.Highlight)
	tr := truncate.New(cfg.MaxLineLen)
	rl := ratelimit.New(cfg.RateLimit)
	dd := dedupe.New(cfg.DedupeWindow)
	rd := redact.New(cfg.RedactFields, cfg.RedactMask)
	st := stats.New()

	mux, err := multiplexer.Build(cfg)
	if err != nil {
		return fmt.Errorf("multiplexer: %w", err)
	}

	fmts := make([]*formatter.Formatter, len(cfg.Services))
	for i, svc := range cfg.Services {
		fmts[i] = formatter.New(svc.Name, i, cfg.Pretty)
	}
	fmtIndex := make(map[string]*formatter.Formatter, len(cfg.Services))
	for i, svc := range cfg.Services {
		fmtIndex[svc.Name] = fmts[i]
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case entry, ok := <-mux:
			if !ok {
				st.Print(os.Stderr)
				return nil
			}
			st.RecordTotal(entry.Service)
			line := rd.Apply(entry.Line)
			line = tr.Apply(line)
			if !f.Match(line) {
				continue
			}
			if !rl.Allow(entry.Service) {
				continue
			}
			if dd.IsDuplicate(entry.Service, line) {
				continue
			}
			st.RecordMatch(entry.Service)
			line = hl.Apply(line)
			fmt.Println(fmtIndex[entry.Service].Render(line))
		case <-sig:
			st.Print(os.Stderr)
			return nil
		}
	}
}
