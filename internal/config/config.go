package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Service describes a single log source.
type Service struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Format string `yaml:"format"` // "pretty" or "raw"
}

// Config is the top-level logdrift configuration.
type Config struct {
	PollInterval time.Duration `yaml:"poll_interval"`
	Filters      []string      `yaml:"filters"`
	Services     []Service     `yaml:"services"`
}

// Load reads and parses a YAML config file from path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = 250 * time.Millisecond
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Services) == 0 {
		return fmt.Errorf("config: at least one service must be defined")
	}
	for i, svc := range c.Services {
		if svc.Name == "" {
			return fmt.Errorf("config: service[%d]: name is required", i)
		}
		if svc.Path == "" {
			return fmt.Errorf("config: service[%d] %q: path is required", i, svc.Name)
		}
		if svc.Format != "" && svc.Format != "pretty" && svc.Format != "raw" {
			return fmt.Errorf("config: service[%d] %q: format must be \"pretty\" or \"raw\"", i, svc.Name)
		}
	}
	return nil
}
