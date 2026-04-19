package buffer_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/buffer"
)

func TestNew_ZeroSize_PushNoOp(t *testing.T) {
	b := buffer.New(0)
	b.Push("line")
	if b.Len() != 0 {
		t.Fatalf("expected 0, got %d", b.Len())
	}
	if b.Lines() != nil {
		t.Fatal("expected nil lines")
	}
}

func TestPush_UnderCapacity(t *testing.T) {
	b := buffer.New(5)
	b.Push("a")
	b.Push("b")
	lines := b.Lines()
	if len(lines) != 2 || lines[0] != "a" || lines[1] != "b" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestPush_OverCapacity_EvictsOldest(t *testing.T) {
	b := buffer.New(3)
	for _, l := range []string{"a", "b", "c", "d"} {
		b.Push(l)
	}
	lines := b.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "b" || lines[1] != "c" || lines[2] != "d" {
		t.Fatalf("unexpected order: %v", lines)
	}
}

func TestLines_OrderIsOldestFirst(t *testing.T) {
	b := buffer.New(4)
	for _, l := range []string{"x", "y", "z"} {
		b.Push(l)
	}
	lines := b.Lines()
	for i, want := range []string{"x", "y", "z"} {
		if lines[i] != want {
			t.Fatalf("index %d: want %q got %q", i, want, lines[i])
		}
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	b := buffer.New(4)
	b.Push("a")
	b.Push("b")
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", b.Len())
	}
	if b.Lines() != nil {
		t.Fatal("expected nil after reset")
	}
}

func TestPush_ExactCapacity_NoEviction(t *testing.T) {
	b := buffer.New(3)
	b.Push("p")
	b.Push("q")
	b.Push("r")
	lines := b.Lines()
	if len(lines) != 3 || lines[0] != "p" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}
