// Package contextlines buffers log lines and emits surrounding context
// lines when a match occurs, similar to grep -B / -A.
package contextlines

import "container/ring"

// Buffer holds a sliding window of recent lines and a look-ahead queue
// so callers can emit N lines before and after a matching line.
type Buffer struct {
	before int
	after  int
	ring   *ring.Ring
	pending []pendingLine
}

type pendingLine struct {
	line      string
	service   string
	countdown int
}

// New creates a Buffer that keeps up to before lines of pre-match context
// and emits up to after lines of post-match context.
func New(before, after int) *Buffer {
	if before < 0 {
		before = 0
	}
	if after < 0 {
		after = 0
	}
	r := ring.New(before + 1)
	return &Buffer{before: before, after: after, ring: r}
}

// Push records a line in the pre-match ring buffer and decrements any
// active post-match countdowns. It returns lines that should be emitted.
func (b *Buffer) Push(line, service string, matched bool) []string {
	var out []string

	// Flush post-match context lines whose countdown reaches this line.
	next := b.pending[:0]
	for _, p := range b.pending {
		out = append(out, p.line)
		p.countdown--
		if p.countdown > 0 {
			next = append(next, p)
		}
	}
	b.pending = next

	if matched {
		// Emit buffered pre-match lines.
		b.ring.Do(func(v any) {
			if v != nil {
				out = append(out, v.(string))
			}
		})
		out = append(out, line)
		if b.after > 0 {
			b.pending = append(b.pending, pendingLine{line: "", service: service, countdown: b.after})
		}
		// Reset ring so pre-match lines aren't re-emitted.
		b.ring = ring.New(b.before + 1)
		return out
	}

	// Store in ring for potential future pre-match context.
	b.ring.Value = line
	b.ring = b.ring.Next()

	// If we are in post-match mode, emit this line.
	if len(b.pending) > 0 {
		out = append(out, line)
	}
	return out
}

// Reset clears all buffered state.
func (b *Buffer) Reset() {
	b.ring = ring.New(b.before + 1)
	b.pending = nil
}
