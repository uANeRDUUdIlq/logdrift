// Package buffer provides a fixed-size ring buffer for log lines,
// allowing recent lines to be replayed on demand.
package buffer

import "sync"

// Buffer is a thread-safe ring buffer storing the last N log lines.
type Buffer struct {
	mu   sync.Mutex
	data []string
	size int
	head int
	count int
}

// New creates a Buffer that retains at most size lines.
// If size <= 0, Push is a no-op and Lines always returns nil.
func New(size int) *Buffer {
	if size < 0 {
		size = 0
	}
	return &Buffer{
		data: make([]string, size),
		size: size,
	}
}

// Push adds a line to the buffer, evicting the oldest entry when full.
func (b *Buffer) Push(line string) {
	if b.size == 0 {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data[b.head] = line
	b.head = (b.head + 1) % b.size
	if b.count < b.size {
		b.count++
	}
}

// Lines returns a snapshot of buffered lines in insertion order (oldest first).
func (b *Buffer) Lines() []string {
	if b.size == 0 {
		return nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.count == 0 {
		return nil
	}
	out := make([]string, b.count)
	start := (b.head - b.count + b.size) % b.size
	for i := 0; i < b.count; i++ {
		out[i] = b.data[(start+i)%b.size]
	}
	return out
}

// Len returns the current number of lines stored.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.count
}

// Reset clears all stored lines.
func (b *Buffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.head = 0
	b.count = 0
}
