// Package config loads and validates logdrift YAML configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Service describes a single log source.
type Service struct {
	Name         string            `yaml:"name"`
	Path         string            `yaml:"path"`
	SampleRate   float64           `yaml:"sample_rate"`   // 0.0–1.0; 0 means unset → default
	Fields       map[string]string `yaml:"fields"`        // field filter rules
	RedactFields []string          `yaml:"redact_fields"`
}

// Config is the top-level configuration structure.
type Config struct {
	PollInterval time.Duration `yaml:"poll_interval"`
	DefaultSampleRate float64  `yaml:"default_sample_rate"`
	Services     []Service     `yaml:"services"`
	OutputFormat string        `yaml:"output_format"` // "pretty" | "raw"
	MaxLineLen   int           `yaml:"max_line_len"`
}

const defaultPollInterval = 200 * time.Millisecond

// Load reads and validates a YAML config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse: %w", err)
	}

	if len(cfg.Services) == 0 {
		return nil, errors.New("config: at least one service must be defined")
	}
	for i, svc := range cfg.Services {
		if svc.Name == "" {
			return nil, fmt.Errorf("config: service[%d]: name is required", i)
		}
		if svc.Path == "" {
			return nil, fmt.Errorf("config: service[%d] %q: path is required", i, svc.Name)
		}
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = defaultPollInterval
	}
	if cfg.DefaultSampleRate == 0 {
		cfg.DefaultSampleRate = 1.0
	}
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = "pretty"
	}
	return &cfg, nil
}
