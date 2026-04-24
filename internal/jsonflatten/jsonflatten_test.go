package jsonflatten_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logdrift/internal/jsonflatten"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	f := jsonflatten.New(".", "")
	const raw = "not json at all"
	if got := f.Apply(raw); got != raw {
		t.Fatalf("want %q, got %q", raw, got)
	}
}

func TestApply_AlreadyFlat_Unchanged(t *testing.T) {
	f := jsonflatten.New(".", "")
	const line = `{"level":"info","msg":"hello"}`
	got := decode(t, f.Apply(line))
	if got["level"] != "info" || got["msg"] != "hello" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestApply_NestedObject_Flattened(t *testing.T) {
	f := jsonflatten.New(".", "")
	line := `{"a":{"b":{"c":1}}}`
	got := decode(t, f.Apply(line))
	val, ok := got["a.b.c"]
	if !ok {
		t.Fatalf("key a.b.c missing; got %v", got)
	}
	// JSON numbers decoded as float64 by default in map[string]any
	if val != int64(1) {
		t.Fatalf("want 1, got %v (%T)", val, val)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	f := jsonflatten.New("_", "")
	line := `{"http":{"status":200}}`
	got := decode(t, f.Apply(line))
	if _, ok := got["http_status"]; !ok {
		t.Fatalf("key http_status missing; got %v", got)
	}
}

func TestApply_Prefix_PrependedToKeys(t *testing.T) {
	f := jsonflatten.New(".", "log")
	line := `{"level":"warn"}`
	got := decode(t, f.Apply(line))
	if _, ok := got["log.level"]; !ok {
		t.Fatalf("key log.level missing; got %v", got)
	}
}

func TestApply_EmptySeparator_DefaultsToDot(t *testing.T) {
	f := jsonflatten.New("", "")
	line := `{"x":{"y":true}}`
	got := decode(t, f.Apply(line))
	if _, ok := got["x.y"]; !ok {
		t.Fatalf("key x.y missing; got %v", got)
	}
}

func TestApply_MixedDepths(t *testing.T) {
	f := jsonflatten.New(".", "")
	line := `{"a":1,"b":{"c":2,"d":{"e":3}}}`
	got := decode(t, f.Apply(line))
	expected := []string{"a", "b.c", "b.d.e"}
	for _, k := range expected {
		if _, ok := got[k]; !ok {
			t.Errorf("key %q missing; got %v", k, got)
		}
	}
	if _, bad := got["b"]; bad {
		t.Errorf("intermediate key 'b' should not appear")
	}
}
