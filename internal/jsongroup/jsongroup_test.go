package jsongroup

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestPush_NonJSON_PassesThrough(t *testing.T) {
	g := New("service", "")
	out := g.Push("not json")
	if len(out) != 1 || out[0] != "not json" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestPush_SameGroup_NoEmit(t *testing.T) {
	g := New("svc", "_count")
	out := g.Push(`{"svc":"api","msg":"a"}`)
	if len(out) != 0 {
		t.Fatalf("expected no output, got %v", out)
	}
	out = g.Push(`{"svc":"api","msg":"b"}`)
	if len(out) != 0 {
		t.Fatalf("expected no output on same group, got %v", out)
	}
}

func TestPush_GroupChange_EmitsSummary(t *testing.T) {
	g := New("svc", "_count")
	g.Push(`{"svc":"api","msg":"a"}`)
	g.Push(`{"svc":"api","msg":"b"}`)
	out := g.Push(`{"svc":"worker","msg":"c"}`)
	if len(out) != 1 {
		t.Fatalf("expected 1 summary line, got %d", len(out))
	}
	m := decode(t, out[0])
	if m["_count"].(float64) != 2 {
		t.Errorf("expected _count=2, got %v", m["_count"])
	}
	if m["svc"] != "api" {
		t.Errorf("expected svc=api, got %v", m["svc"])
	}
}

func TestFlush_EmptyBuffer_ReturnsNil(t *testing.T) {
	g := New("svc", "_count")
	if out := g.Flush(); out != nil {
		t.Fatalf("expected nil, got %v", out)
	}
}

func TestFlush_PendingLines_EmitsSummary(t *testing.T) {
	g := New("svc", "_count")
	g.Push(`{"svc":"db","latency":5}`)
	g.Push(`{"svc":"db","latency":10}`)
	out := g.Flush()
	if len(out) != 1 {
		t.Fatalf("expected 1 line, got %d", len(out))
	}
	m := decode(t, out[0])
	if m["_count"].(float64) != 2 {
		t.Errorf("_count mismatch: %v", m["_count"])
	}
}

func TestNew_DefaultCountKey(t *testing.T) {
	g := New("svc", "")
	if g.countKey != "_count" {
		t.Errorf("expected default countKey _count, got %q", g.countKey)
	}
}

func TestPush_NonJSON_FlushesPending(t *testing.T) {
	g := New("svc", "_count")
	g.Push(`{"svc":"api","x":1}`)
	out := g.Push("plain text")
	// Should emit the buffered group summary AND the plain line.
	if len(out) != 2 {
		t.Fatalf("expected 2 lines (summary + passthrough), got %v", out)
	}
	if out[1] != "plain text" {
		t.Errorf("expected passthrough as second element, got %q", out[1])
	}
}
