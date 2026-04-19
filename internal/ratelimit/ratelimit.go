// Package ratelimit provides a simple token-bucket rate limiter for log lines.
// It can be used to suppress noisy services that emit too many lines per second.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks per-service token buckets.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     int // max lines per second
	window   time.Duration
}

type bucket struct {
	count    int
	resetAt  time.Time
}

// New creates a Limiter that allows up to rate lines per second per service.
// A rate of 0 disables limiting.
func New(rate int) *Limiter {
	return &Limiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		window:  time.Second,
	}
}

// Allow returns true if the line from the given service should be allowed through.
func (l *Limiter) Allow(service string) bool {
	if l.rate <= 0 {
		return true
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[service]
	if !ok || now.After(b.resetAt) {
		l.buckets[service] = &bucket{count: 1, resetAt: now.Add(l.window)}
		return true
	}

	if b.count >= l.rate {
		return false
	}

	b.count++
	return true
}

// Status returns the current line count and reset time for the given service.
// If the service has no active bucket, count is 0 and resetAt is zero.
func (l *Limiter) Status(service string) (count int, resetAt time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.buckets[service]
	if !ok || time.Now().After(b.resetAt) {
		return 0, time.Time{}
	}
	return b.count, b.resetAt
}

// Reset clears all buckets, useful for testing.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buckets = make(map[string]*bucket)
}
