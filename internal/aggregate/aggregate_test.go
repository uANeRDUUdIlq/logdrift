package aggregate

import (
	"bytes"
	"strings"
	"testing"
)

func TestRecord_NonJSON_CountsUnknown(t *testing.T) {
	a := New("level")
	a.Record("not json")
	if a.Counts()["<unknown>"] != 1 {
		t.Fatal("expected <unknown> count 1")
	}
}

func TestRecord_MissingField_CountsUnknown(t *testing.T) {
	a := New("level")
	a.Record(`{"msg":"hello"}`)
	if a.Counts()["<unknown>"] != 1 {
		t.Fatal("expected <unknown> count 1")
	}
}

func TestRecord_KnownField_CountsValue(t *testing.T) {
	a := New("level")
	a.Record(`{"level":"info"}`)
	a.Record(`{"level":"info"}`)
	a.Record(`{"level":"error"}`)
	c := a.Counts()
	if c["info"] != 2 {
		t.Fatalf("expected info=2, got %d", c["info"])
	}
	if c["error"] != 1 {
		t.Fatalf("expected error=1, got %d", c["error"])
	}
}

func TestCounts_ReturnsCopy(t *testing.T) {
	a := New("level")
	a.Record(`{"level":"info"}`)
	c := a.Counts()
	c["info"] = 999
	if a.Counts()["info"] != 1 {
		t.Fatal("Counts should return a copy")
	}
}

func TestPrint_SortedOutput(t *testing.T) {
	a := New("level")
	a.Record(`{"level":"warn"}`)
	a.Record(`{"level":"info"}`)
	a.Record(`{"level":"info"}`)
	var buf bytes.Buffer
	a.Print(&buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + separator + info + warn
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d: %v", len(lines), lines)
	}
	if !strings.Contains(lines[2], "info") {
		t.Errorf("expected info before warn, got: %s", lines[2])
	}
	if !strings.Contains(lines[3], "warn") {
		t.Errorf("expected warn last, got: %s", lines[3])
	}
}
