package fieldselect

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	return m
}

func TestApply_NoFields_ReturnsSameLine(t *testing.T) {
	s := New(nil)
	line := `{"level":"info","msg":"hello"}`
	if got := s.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	s := New([]string{"level"})
	line := "not json at all"
	if got := s.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_KeepsSelectedFields(t *testing.T) {
	s := New([]string{"level", "msg"})
	line := `{"level":"info","msg":"hello","service":"api","ts":"2024-01-01"}`
	got := decode(t, s.Apply(line))

	if _, ok := got["level"]; !ok {
		t.Error("expected 'level' to be present")
	}
	if _, ok := got["msg"]; !ok {
		t.Error("expected 'msg' to be present")
	}
	if _, ok := got["service"]; ok {
		t.Error("expected 'service' to be absent")
	}
	if _, ok := got["ts"]; ok {
		t.Error("expected 'ts' to be absent")
	}
}

func TestApply_FieldMissingInLine_Omitted(t *testing.T) {
	s := New([]string{"level", "trace_id"})
	line := `{"level":"warn","msg":"oops"}`
	got := decode(t, s.Apply(line))

	if _, ok := got["level"]; !ok {
		t.Error("expected 'level' to be present")
	}
	if _, ok := got["trace_id"]; ok {
		t.Error("'trace_id' should be absent when not in source")
	}
}

func TestApply_AllFieldsSelected_EquivalentContent(t *testing.T) {
	s := New([]string{"a", "b"})
	line := `{"a":1,"b":2}`
	got := decode(t, s.Apply(line))
	if len(got) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(got))
	}
}
