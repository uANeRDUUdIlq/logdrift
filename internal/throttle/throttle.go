// Package throttle provides a token-bucket throttler that limits the number
// of log lines emitted per service per second, dropping excess lines.
package throttle

import (
	"sync"
	"time"
)

// Throttler tracks per-service token buckets.
type Throttler struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     int // max lines per second (0 = unlimited)
	interval time.Duration
}

type bucket struct {
	tokens    int
	resetAt   time.Time
}

// New creates a Throttler with the given lines-per-second rate.
// A rate of 0 disables throttling.
func New(rate int) *Throttler {
	return &Throttler{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		interval: time.Second,
	}
}

// Allow returns true if the line from the given service should be emitted.
func (t *Throttler) Allow(service string) bool {
	if t.rate <= 0 {
		return true
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	b, ok := t.buckets[service]
	if !ok || now.After(b.resetAt) {
		t.buckets[service] = &bucket{
			tokens:  t.rate - 1,
			resetAt: now.Add(t.interval),
		}
		return true
	}
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

// Reset clears all buckets.
func (t *Throttler) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.buckets = make(map[string]*bucket)
}
