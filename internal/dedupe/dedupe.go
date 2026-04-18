// Package dedupe provides a sliding-window deduplication filter that
// suppresses repeated log lines within a configurable time window.
package dedupe

import (
	"sync"
	"time"
)

// Deduper tracks recently seen log lines and reports whether a line is a
// duplicate within the configured window.
type Deduper struct {
	mu      sync.Mutex
	window  time.Duration
	seen    map[string]time.Time
	nowFunc func() time.Time
}

// New creates a Deduper that suppresses identical lines seen within window.
// A zero or negative window disables deduplication (all lines are allowed).
func New(window time.Duration) *Deduper {
	return &Deduper{
		window:  window,
		seen:    make(map[string]time.Time),
		nowFunc: time.Now,
	}
}

// IsDuplicate returns true if line was already seen within the window.
// It also records the line if it is new or has expired.
func (d *Deduper) IsDuplicate(service, line string) bool {
	if d.window <= 0 {
		return false
	}

	key := service + "\x00" + line
	now := d.nowFunc()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.evict(now)

	if _, exists := d.seen[key]; exists {
		return true
	}
	d.seen[key] = now
	return false
}

// evict removes entries older than the window. Must be called with mu held.
func (d *Deduper) evict(now time.Time) {
	for k, t := range d.seen {
		if now.Sub(t) > d.window {
			delete(d.seen, k)
		}
	}
}

// Reset clears all tracked entries.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}
