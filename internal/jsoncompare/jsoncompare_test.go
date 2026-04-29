package jsoncompare_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/jsoncompare"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoFields_PassesThrough(t *testing.T) {
	c := jsoncompare.New(nil)
	line := `{"level":"info"}`
	out, ok := c.Apply("svc", line)
	if !ok || out != line {
		t.Fatalf("expected passthrough, got ok=%v out=%q", ok, out)
	}
}

func TestApply_NonJSON_PassesThrough(t *testing.T) {
	c := jsoncompare.New([]string{"level"})
	out, ok := c.Apply("svc", "not json")
	if !ok || out != "not json" {
		t.Fatalf("expected passthrough for non-JSON, got ok=%v out=%q", ok, out)
	}
}

func TestApply_FirstSeen_Emitted(t *testing.T) {
	c := jsoncompare.New([]string{"status"})
	out, ok := c.Apply("svc", `{"status":"ok"}`)
	if !ok {
		t.Fatal("expected first-seen line to be emitted")
	}
	m := decode(t, out)
	if m["_next_status"] != "ok" {
		t.Fatalf("expected _next_status=ok, got %v", m["_next_status"])
	}
	if _, has := m["_prev_status"]; has {
		t.Fatal("expected no _prev_status on first-seen")
	}
}

func TestApply_SameValue_Suppressed(t *testing.T) {
	c := jsoncompare.New([]string{"status"})
	c.Apply("svc", `{"status":"ok"}`)
	_, ok := c.Apply("svc", `{"status":"ok"}`)
	if ok {
		t.Fatal("expected duplicate value to be suppressed")
	}
}

func TestApply_ChangedValue_EmittedWithAnnotations(t *testing.T) {
	c := jsoncompare.New([]string{"status"})
	c.Apply("svc", `{"status":"ok"}`)
	out, ok := c.Apply("svc", `{"status":"error"}`)
	if !ok {
		t.Fatal("expected changed value to be emitted")
	}
	m := decode(t, out)
	if m["_prev_status"] != "ok" {
		t.Fatalf("expected _prev_status=ok, got %v", m["_prev_status"])
	}
	if m["_next_status"] != "error" {
		t.Fatalf("expected _next_status=error, got %v", m["_next_status"])
	}
}

func TestApply_DifferentServices_Independent(t *testing.T) {
	c := jsoncompare.New([]string{"code"})
	c.Apply("svcA", `{"code":200}`)
	// svcB has never seen code=200, so it should emit
	_, ok := c.Apply("svcB", `{"code":200}`)
	if !ok {
		t.Fatal("expected svcB first-seen to be emitted independently")
	}
}

func TestApply_MissingField_DoesNotChangeState(t *testing.T) {
	c := jsoncompare.New([]string{"level"})
	// First emit establishes state
	c.Apply("svc", `{"level":"info"}`)
	// Line without the field should be suppressed (no change detected)
	_, ok := c.Apply("svc", `{"msg":"hello"}`)
	if ok {
		t.Fatal("expected line without tracked field to be suppressed")
	}
}
