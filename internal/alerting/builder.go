package alerting

import (
	"fmt"
	"io"
	"regexp"
	"time"
)

// RuleConfig is the serialisable form of a Rule, suitable for YAML/JSON config.
type RuleConfig struct {
	Name     string `yaml:"name"     json:"name"`
	Pattern  string `yaml:"pattern"  json:"pattern"`
	Cooldown string `yaml:"cooldown" json:"cooldown"`
}

// Build constructs an Alerter from a slice of RuleConfig values.
// It returns an error if any pattern fails to compile or any cooldown
// duration cannot be parsed.
func Build(out io.Writer, cfgs []RuleConfig) (*Alerter, error) {
	rules := make([]Rule, 0, len(cfgs))
	for _, c := range cfgs {
		if c.Name == "" {
			return nil, fmt.Errorf("alerting: rule missing name")
		}
		re, err := regexp.Compile(c.Pattern)
		if err != nil {
			return nil, fmt.Errorf("alerting: rule %q bad pattern: %w", c.Name, err)
		}
		var cd time.Duration
		if c.Cooldown != "" {
			cd, err = time.ParseDuration(c.Cooldown)
			if err != nil {
				return nil, fmt.Errorf("alerting: rule %q bad cooldown: %w", c.Name, err)
			}
		}
		rules = append(rules, Rule{Name: c.Name, Pattern: re, Cooldown: cd})
	}
	return New(out, rules), nil
}
