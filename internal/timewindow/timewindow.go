// Package timewindow provides filtering of log lines by a time range.
// Lines are parsed for a timestamp field and dropped if they fall outside
// the configured [Since, Until] window.
package timewindow

import (
	"encoding/json"
	"time"
)

// Filter drops log lines whose timestamp falls outside the window.
type Filter struct {
	field string
	since time.Time
	until time.Time
}

// New creates a Filter. field is the JSON key holding the timestamp.
// Zero values for since/until mean "unbounded" on that side.
func New(field string, since, until time.Time) *Filter {
	if field == "" {
		field = "time"
	}
	return &Filter{field: field, since: since, until: until}
}

// Allow returns true when the line should be kept.
func (f *Filter) Allow(line string) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return true // non-JSON passes through
	}

	raw, ok := obj[f.field]
	if !ok {
		return true // no timestamp field – pass through
	}

	var t time.Time
	switch v := raw.(type) {
	case string:
		var err error
		for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05"} {
			t, err = time.Parse(layout, v)
			if err == nil {
				break
			}
		}
		if t.IsZero() {
			return true
		}
	case float64:
		t = time.Unix(int64(v), 0).UTC()
	default:
		return true
	}

	if !f.since.IsZero() && t.Before(f.since) {
		return false
	}
	if !f.until.IsZero() && t.After(f.until) {
		return false
	}
	return true
}
