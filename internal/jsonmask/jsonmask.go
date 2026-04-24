// Package jsonmask selectively masks JSON fields whose values match a given
// regular expression, replacing the matched portion with a configurable mask
// string. Unlike redact (which targets field names), jsonmask targets field
// values, making it useful for scrubbing sensitive patterns such as credit-card
// numbers, tokens, or UUIDs embedded anywhere in a log line.
package jsonmask

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Rule pairs a compiled pattern with the mask string to substitute.
type Rule struct {
	Pattern *regexp.Regexp
	Mask    string
}

// Masker applies value-pattern masking to JSON log lines.
type Masker struct {
	rules []Rule
}

// New returns a Masker built from the supplied rules. Rules with a nil Pattern
// are silently skipped.
func New(rules []Rule) *Masker {
	filtered := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.Pattern != nil {
			if r.Mask == "" {
				r.Mask = "[MASKED]"
			}
			filtered = append(filtered, r)
		}
	}
	return &Masker{rules: filtered}
}

// Apply scans every string value in the top-level JSON object and replaces
// any sub-string that matches a rule pattern with the corresponding mask.
// Non-JSON lines and lines with no matching values are returned unchanged.
func (m *Masker) Apply(line string) string {
	if len(m.rules) == 0 {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	changed := false
	for k, raw := range obj {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			continue
		}
		masked := s
		for _, r := range m.rules {
			masked = r.Pattern.ReplaceAllString(masked, r.Mask)
		}
		if masked != s {
			encoded, err := json.Marshal(masked)
			if err == nil {
				obj[k] = json.RawMessage(encoded)
				changed = true
			}
		}
	}

	if !changed {
		return line
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return strings.TrimSpace(string(out))
}
