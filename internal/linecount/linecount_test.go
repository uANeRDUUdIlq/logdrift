package linecount

import (
	"bytes"
	"strings"
	"testing"
)

func TestRecordSeen_IncreasesCount(t *testing.T) {
	c := New()
	c.RecordSeen("api")
	c.RecordSeen("api")
	seen, _ := c.Snapshot()
	if seen["api"] != 2 {
		t.Fatalf("expected 2, got %d", seen["api"])
	}
}

func TestRecordEmitted_IncreasesCount(t *testing.T) {
	c := New()
	c.RecordEmitted("worker")
	_, emitted := c.Snapshot()
	if emitted["worker"] != 1 {
		t.Fatalf("expected 1, got %d", emitted["worker"])
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	c := New()
	c.RecordSeen("svc")
	seen, _ := c.Snapshot()
	seen["svc"] = 999
	seen2, _ := c.Snapshot()
	if seen2["svc"] != 1 {
		t.Fatal("snapshot should be a copy")
	}
}

func TestSnapshot_IndependentServices(t *testing.T) {
	c := New()
	c.RecordSeen("a")
	c.RecordSeen("b")
	c.RecordSeen("b")
	seen, _ := c.Snapshot()
	if seen["a"] != 1 || seen["b"] != 2 {
		t.Fatalf("unexpected counts: %v", seen)
	}
}

func TestPrint_SortedOutput(t *testing.T) {
	c := New()
	c.RecordSeen("zebra")
	c.RecordSeen("alpha")
	c.RecordEmitted("alpha")
	var buf bytes.Buffer
	c.Print(&buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[1], "alpha") {
		t.Errorf("expected alpha first, got: %s", lines[1])
	}
	if !strings.Contains(lines[2], "zebra") {
		t.Errorf("expected zebra second, got: %s", lines[2])
	}
}

func TestPrint_HeaderAlwaysPresent(t *testing.T) {
	c := New()
	var buf bytes.Buffer
	c.Print(&buf)
	if !strings.Contains(buf.String(), "SERVICE") {
		t.Error("expected header line")
	}
}
