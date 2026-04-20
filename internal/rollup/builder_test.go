package rollup

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestBuild_EmptyField_ReturnsError(t *testing.T) {
	_, err := Build(Config{Field: "", Window: time.Second}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestBuild_NegativeWindow_ReturnsError(t *testing.T) {
	_, err := Build(Config{Field: "level", Window: -1}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestBuild_Valid_ReturnsRollup(t *testing.T) {
	r, err := Build(Config{Field: "level", Window: 0}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil rollup")
	}
}

func TestBuild_WriterReceivesFlushOutput(t *testing.T) {
	var buf bytes.Buffer
	r, _ := Build(Config{Field: "level", Window: 0}, &buf)
	r.Record("svc", `{"level":"error"}`)
	r.Flush()
	// drain the entry so the goroutine writes
	<-r.Entries()
	// give goroutine time to write
	time.Sleep(20 * time.Millisecond)
	if !strings.Contains(buf.String(), "level=error") {
		t.Fatalf("expected output to contain level=error, got: %q", buf.String())
	}
}
