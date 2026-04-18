// Package config loads and validates logdrift YAML configuration.
package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Service describes a single tailed log source.
type Service struct {
	Name     string            `yaml:"name"`
	Path     string            `yaml:"path"`
	Filters  []string          `yaml:"filters"`
	Transform TransformConfig  `yaml:"transform"`
}

// TransformConfig holds per-service field transformation rules.
type TransformConfig struct {
	Rename map[string]string `yaml:"rename"`
	Add    map[string]string `yaml:"add"`
}

// Config is the top-level logdrift configuration.
type Config struct {
	PollInterval time.Duration `yaml:"poll_interval"`
	Services     []Service     `yaml:"services"`
	Format       string        `yaml:"format"`
	MaxLen       int           `yaml:"max_len"`
	MinLevel     string        `yaml:"min_level"`
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
		return nil, errors.New("config: at least one service is required")
	}
	for i, s := range cfg.Services {
		if s.Name == "" {
			return nil, errors.New("config: service name is required")
		}
		if s.Path == "" {
			return nil, errors.New("config: service path is required")
		}
		_ = i
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = defaultPollInterval
	}
	if cfg.Format == "" {
		cfg.Format = "pretty"
	}
	return &cfg, nil
}
