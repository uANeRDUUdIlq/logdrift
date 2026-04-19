// Package alerting emits a notification when a log line matches a named alert rule.
package alerting

import (
	"fmt"
	"io"
	"regexp"
	"time"
)

// Rule defines a single alert condition.
type Rule struct {
	Name    string
	Pattern *regexp.Regexp
	Cooldown time.Duration
}

// Alerter checks log lines against alert rules and writes notifications.
type Alerter struct {
	rules    []Rule
	out      io.Writer
	lastFire map[string]time.Time
}

// New returns an Alerter that writes alert notifications to out.
func New(out io.Writer, rules []Rule) *Alerter {
	return &Alerter{
		rules:    rules,
		out:      out,
		lastFire: make(map[string]time.Time),
	}
}

// Check tests line against every rule and writes an alert for each match
// that is not within the rule's cooldown window.
func (a *Alerter) Check(service, line string) {
	now := time.Now()
	for _, r := range a.rules {
		if !r.Pattern.MatchString(line) {
			continue
		}
		key := r.Name + "|" + service
		if last, ok := a.lastFire[key]; ok && r.Cooldown > 0 && now.Sub(last) < r.Cooldown {
			continue
		}
		a.lastFire[key] = now
		fmt.Fprintf(a.out, "[ALERT] rule=%q service=%q\n", r.Name, service)
	}
}

// Reset clears the cooldown state for all rules.
func (a *Alerter) Reset() {
	a.lastFire = make(map[string]time.Time)
}
