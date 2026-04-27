package jsoncoerce

import "fmt"

// Config holds the configuration for building a Coercer from external
// sources such as a YAML config file.
type Config struct {
	// Rules is the ordered list of coercion rules to apply.
	Rules []RuleConfig `yaml:"rules"`
}

// RuleConfig is the serialisable form of a single coercion rule.
type RuleConfig struct {
	Field  string `yaml:"field"`
	Target string `yaml:"target"`
}

// Build validates cfg and returns a ready-to-use Coercer.
// It returns an error if any rule contains an unknown target type or an
// empty field name.
func Build(cfg Config) (*Coercer, error) {
	rules := make([]Rule, 0, len(cfg.Rules))
	for _, rc := range cfg.Rules {
		if rc.Field == "" {
			return nil, fmt.Errorf("jsoncoerce: rule has empty field name")
		}
		var tt TargetType
		switch TargetType(rc.Target) {
		case TypeString:
			tt = TypeString
		case TypeNumber:
			tt = TypeNumber
		case TypeBool:
			tt = TypeBool
		default:
			return nil, fmt.Errorf("jsoncoerce: unknown target type %q for field %q", rc.Target, rc.Field)
		}
		rules = append(rules, Rule{Field: rc.Field, Target: tt})
	}
	return New(rules), nil
}
