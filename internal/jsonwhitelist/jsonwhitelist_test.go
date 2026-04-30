package jsonwhitelist_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/jsonwhitelist"
)

func decode(t *testing.T, line string) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoFields_ReturnsSameLine(t *testing.T) {
	w := jsonwhitelist.New(nil)
	input := `{"level":"info","msg":"hello"}`
	if got := w.Apply(input); got != input {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	w := jsonwhitelist.New([]string{"level"})
	input := "not json at all"
	if got := w.Apply(input); got != input {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_KeepsAllowedFields(t *testing.T) {
	w := jsonwhitelist.New([]string{"level", "msg"})
	input := `{"level":"info","msg":"hello","ts":"2024-01-01","caller":"main.go"}`
	got := w.Apply(input)
	m := decode(t, got)
	if _, ok := m["level"]; !ok {
		t.Error("expected 'level' to be kept")
	}
	if _, ok := m["msg"]; !ok {
		t.Error("expected 'msg' to be kept")
	}
	if _, ok := m["ts"]; ok {
		t.Error("expected 'ts' to be removed")
	}
	if _, ok := m["caller"]; ok {
		t.Error("expected 'caller' to be removed")
	}
}

func TestApply_AllFieldsMissing_EmitsEmptyObject(t *testing.T) {
	w := jsonwhitelist.New([]string{"nonexistent"})
	input := `{"level":"info","msg":"hello"}`
	got := w.Apply(input)
	m := decode(t, got)
	if len(m) != 0 {
		t.Fatalf("expected empty object, got %v", m)
	}
}

func TestApply_SingleField_Retained(t *testing.T) {
	w := jsonwhitelist.New([]string{"msg"})
	input := `{"level":"warn","msg":"oops","code":500}`
	got := w.Apply(input)
	m := decode(t, got)
	if len(m) != 1 {
		t.Fatalf("expected 1 field, got %d", len(m))
	}
	if _, ok := m["msg"]; !ok {
		t.Error("expected 'msg' to be kept")
	}
}

func TestApply_EmptyFieldList_IsNoOp(t *testing.T) {
	w := jsonwhitelist.New([]string{})
	input := `{"level":"debug","msg":"trace"}`
	if got := w.Apply(input); got != input {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}
