// Package rollup groups log lines by a JSON field value within a time
// window and emits a single summary line when the window closes.
package rollup

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Entry is a rolled-up summary emitted after a window closes.
type Entry struct {
	Service string
	Field   string
	Value   string
	Count   int
	First   time.Time
	Last    time.Time
}

type bucket struct {
	count int
	first time.Time
	last  time.Time
}

// Rollup accumulates counts keyed by (service, fieldValue).
type Rollup struct {
	mu      sync.Mutex
	field   string
	window  time.Duration
	buckets map[string]*bucket // key: service+"\x00"+value
	ticker  *time.Ticker
	out     chan Entry
	stop    chan struct{}
}

// New creates a Rollup that groups by field over window duration.
// window=0 disables periodic flushing (manual Flush only).
func New(field string, window time.Duration) *Rollup {
	r := &Rollup{
		field:   field,
		window:  window,
		buckets: make(map[string]*bucket),
		out:     make(chan Entry, 256),
		stop:    make(chan struct{}),
	}
	if window > 0 {
		r.ticker = time.NewTicker(window)
		go r.loop()
	}
	return r
}

// Record ingests a raw log line for the given service.
func (r *Rollup) Record(service, line string) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return
	}
	v, ok := m[r.field]
	if !ok {
		return
	}
	val := fmt.Sprintf("%v", v)
	now := time.Now()
	key := service + "\x00" + val
	r.mu.Lock()
	b, exists := r.buckets[key]
	if !exists {
		b = &bucket{first: now}
		r.buckets[key] = b
	}
	b.count++
	b.last = now
	r.mu.Unlock()
}

// Flush emits all current buckets and resets state.
func (r *Rollup) Flush() {
	r.mu.Lock()
	snap := r.buckets
	r.buckets = make(map[string]*bucket)
	r.mu.Unlock()
	for key, b := range snap {
		var svc, val string
		for i, c := range key {
			if c == 0 {
				svc = key[:i]
				val = key[i+1:]
				break
			}
		}
		r.out <- Entry{Service: svc, Field: r.field, Value: val, Count: b.count, First: b.first, Last: b.last}
	}
}

// Entries returns the output channel.
func (r *Rollup) Entries() <-chan Entry { return r.out }

// Stop halts the background ticker.
func (r *Rollup) Stop() {
	if r.ticker != nil {
		r.ticker.Stop()
		close(r.stop)
	}
}

func (r *Rollup) loop() {
	for {
		select {
		case <-r.ticker.C:
			r.Flush()
		case <-r.stop:
			return
		}
	}
}
