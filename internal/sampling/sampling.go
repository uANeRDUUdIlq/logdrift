// Package sampling provides probabilistic log sampling for logdrift.
// It allows retaining only a fraction of log lines per service to reduce noise.
package sampling

import (
	"math/rand"
	"sync"
)

// Sampler decides whether a log line should be kept based on a per-service rate.
type Sampler struct {
	mu       sync.Mutex
	rates    map[string]float64 // 0.0 = drop all, 1.0 = keep all
	default_ float64
	rng      *rand.Rand
}

// New creates a Sampler. defaultRate applies to services not listed in rates.
// Values are clamped to [0.0, 1.0].
func New(defaultRate float64, rates map[string]float64) *Sampler {
	if rates == nil {
		rates = make(map[string]float64)
	}
	return &Sampler{
		rates:    rates,
		default_: clamp(defaultRate),
		rng:      rand.New(rand.NewSource(42)),
	}
}

// Allow returns true if the log line for the given service should be kept.
func (s *Sampler) Allow(service string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	rate, ok := s.rates[service]
	if !ok {
		rate = s.default_
	}
	rate = clamp(rate)
	if rate <= 0.0 {
		return false
	}
	if rate >= 1.0 {
		return true
	}
	return s.rng.Float64() < rate
}

// SetRate updates the sampling rate for a specific service at runtime.
func (s *Sampler) SetRate(service string, rate float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rates[service] = clamp(rate)
}

func clamp(v float64) float64 {
	if v < 0.0 {
		return 0.0
	}
	if v > 1.0 {
		return 1.0
	}
	return v
}
