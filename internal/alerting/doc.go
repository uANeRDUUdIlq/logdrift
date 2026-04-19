// Package alerting provides pattern-based alert rules for log lines.
//
// Each Rule pairs a name with a compiled regular expression and an optional
// cooldown duration. When a log line matches, an alert notification is written
// to the configured io.Writer. Repeated matches within the cooldown window are
// suppressed to avoid notification storms.
//
// Usage:
//
//	rules := []alerting.Rule{
//		{Name: "oom", Pattern: regexp.MustCompile(`out of memory`), Cooldown: time.Minute},
//	}
//	a := alerting.New(os.Stderr, rules)
//	a.Check("api", logLine)
package alerting
