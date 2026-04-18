// Package config loads and validates logdrift configuration files.
package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Service describes a single log source.
type Service struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// HighlightRule maps a keyword to a display color.
type HighlightRule struct {
	Word  string `yaml:"word"`
	Color string `yaml:"color"`
}

// Config is the top-level logdrift configuration.
type Config struct {
	PollInterval time.Duration   `yaml:"poll_interval"`
	Pretty       bool            `yaml:"pretty"`
	RawOutput    bool            `yaml:"raw_output"`
	FilterExpr   string          `yaml:"filter"`
	Highlight    []HighlightRule `yaml:"highlight"`
	Services     []Service       `yaml:"services"`
}

const defaultPollInterval = 500 * time.Millisecond

// Load reads and validates a YAML config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.Services) == 0 {
		return nil, errors.New("config: at least one service must be defined")
	}
	for i, s := range cfg.Services {
		if s.Name == "" {
			return nil, errors.New("config: service name must not be empty")
		}
		if s.Path == "" {
			return nil, errors.New("config: service path must not be empty")
		}
		_ = i
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = defaultPollInterval
	}
	return &cfg, nil
}
