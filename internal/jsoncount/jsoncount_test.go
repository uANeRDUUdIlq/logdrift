package jsoncount

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRecord_NonJSON_Ignored(t *testing.T) {
	var buf bytes.Buffer
	c := New("level", 0, &buf)
	c.Record("not json at all")
	c.Flush()
	if buf.Len() != 0 {
		t.Fatalf("expected no output for non-JSON line, got %q", buf.String())
	}
}

func TestRecord_MissingField_Ignored(t *testing.T) {
	var buf bytes.Buffer
	c := New("level", 0, &buf)
	c.Record(`{"msg":"hello"}`)
	c.Flush()
	if buf.Len() != 0 {
		t.Fatalf("expected no output when field absent, got %q", buf.String())
	}
}

func TestRecord_KnownField_CountsValue(t *testing.T) {
	var buf bytes.Buffer
	c := New("level", 0, &buf)
	c.Record(`{"level":"info","msg":"a"}`)
	c.Record(`{"level":"info","msg":"b"}`)
	c.Record(`{"level":"error","msg":"c"}`)
	c.Flush()

	out := buf.String()
	if !strings.Contains(out, "info") {
		t.Errorf("expected 'info' in output, got %q", out)
	}
	if !strings.Contains(out, "error") {
		t.Errorf("expected 'error' in output, got %q", out)
	}
}

func TestFlush_ResetsCounters(t *testing.T) {
	var buf bytes.Buffer
	c := New("level", 0, &buf)
	c.Record(`{"level":"warn"}`)
	c.Flush()
	buf.Reset()
	c.Flush() // second flush should produce no output
	if buf.Len() != 0 {
		t.Fatalf("expected empty output after reset, got %q", buf.String())
	}
}

func TestFlush_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	c := New("level", 0, &buf)
	for _, lvl := range []string{"warn", "error", "debug", "info"} {
		c.Record(`{"level":"` + lvl + `"}`)
	}
	c.Flush()
	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	// lines[0] is header; values start at index 1
	if len(lines) < 5 {
		t.Fatalf("expected at least 5 lines, got %d: %q", len(lines), out)
	}
	order := []string{"debug", "error", "info", "warn"}
	for i, want := range order {
		if !strings.Contains(lines[i+1], want) {
			t.Errorf("line %d: want %q, got %q", i+1, want, lines[i+1])
		}
	}
}

func TestStop_HaltsBackgroundFlush(t *testing.T) {
	var buf bytes.Buffer
	c := New("level", 50*time.Millisecond, &buf)
	c.Stop()
	time.Sleep(120 * time.Millisecond)
	// No panic and no goroutine leak is the success condition.
}
