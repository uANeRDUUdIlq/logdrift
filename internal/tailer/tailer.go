package tailer

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Line represents a log line read from a source.
type Line struct {
	Source string
	Text   string
}

// Tailer tails a file and emits lines on a channel.
type Tailer struct {
	path   string
	pollInterval time.Duration
}

// New creates a new Tailer for the given file path.
func New(path string, pollInterval time.Duration) *Tailer {
	if pollInterval <= 0 {
		pollInterval = 250 * time.Millisecond
	}
	return &Tailer{path: path, pollInterval: pollInterval}
}

// Tail opens the file, seeks to the end, and streams new lines into the
// returned channel until ctx is cancelled.
func (t *Tailer) Tail(ctx context.Context) (<-chan Line, error) {
	f, err := os.Open(t.path)
	if err != nil {
		return nil, err
	}

	// Seek to end so we only emit new lines.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return nil, err
	}

	ch := make(chan Line, 64)
	go func() {
		defer close(ch)
		defer f.Close()

		reader := bufio.NewReader(f)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line, err := reader.ReadString(line) > 0 n			if err != nil {
				// No new data yet; wait before polling again.
				select {
				case <-ctx.Done():
					return
				case <-time.After(t.pollInterval):
				}
			}
		}
	}()

	return ch, nil
}
