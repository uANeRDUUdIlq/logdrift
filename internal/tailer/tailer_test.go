package tailer

import (
	"context"
	"os"
	"testing"
	"time"
)

func writeTmp(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logdrift-*.log")
	if err != nil {
		t.Fatal(err)
	}
	if content != "" {
		if _, err := f.WriteString(content); err != nil {
			t.Fatal(err)
		}
	}
	return f
}

func TestNew_DefaultPollInterval(t *testing.T) {
	tlr := New("/dev/null", 0)
	if tlr.pollInterval != 250*time.Millisecond {
		t.Errorf("expected 250ms, got %v", tlr.pollInterval)
	}
}

func TestTail_ReceivesNewLines(t *testing.T) {
	f := writeTmp(t, "")
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tlr := New(f.Name(), 10*time.Millisecond)
	ch, err := tlr.Tail(ctx)
	if err != nil {
		t.Fatal(err)
	}

	want := "{\"level\":\"info\"}\n"
	if _, err := f.WriteString(want); err != nil {
		t.Fatal(err)
	}

	select {
	case line := <-ch:
		if line.Text != want {
			t.Errorf("got %q, want %q", line.Text, want)
		}
		if line.Source != f.Name() {
			t.Errorf("source: got %q, want %q", line.Source, f.Name())
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for line")
	}
}

func TestTail_FileNotFound(t *testing.T) {
	tlr := New("/nonexistent/path/file.log", 10*time.Millisecond)
	_, err := tlr.Tail(context.Background())
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
