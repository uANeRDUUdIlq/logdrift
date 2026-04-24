package jsonrename_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logdrift/internal/jsonrename"
)

func decode(t *testing.T, line string) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoMapping_ReturnsSameLine(t *testing.T) {
	r := jsonrename.New(nil)
	line := `{"level":"info","msg":"hello"}`
	if got := r.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	r := jsonrename.New(map[string]string{"a": "b"})
	line := "not json at all"
	if got := r.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_RenamesKey(t *testing.T) {
	r := jsonrename.New(map[string]string{"msg": "message"})
	line := `{"level":"info","msg":"hello"}`
	got := r.Apply(line)
	m := decode(t, got)
	if _, ok := m["message"]; !ok {
		t.Error("expected key 'message' to exist")
	}
	if _, ok := m["msg"]; ok {
		t.Error("expected old key 'msg' to be removed")
	}
}

func TestApply_MissingKey_IsNoOp(t *testing.T) {
	r := jsonrename.New(map[string]string{"nonexistent": "other"})
	line := `{"level":"warn"}`
	got := r.Apply(line)
	m := decode(t, got)
	if _, ok := m["level"]; !ok {
		t.Error("expected 'level' to be preserved")
	}
	if _, ok := m["other"]; ok {
		t.Error("unexpected key 'other'")
	}
}

func TestApply_MultipleRenames(t *testing.T) {
	r := jsonrename.New(map[string]string{
		"msg":   "message",
		"level": "severity",
	})
	line := `{"level":"error","msg":"boom","ts":"2024-01-01"}`
	got := r.Apply(line)
	m := decode(t, got)
	for _, want := range []string{"message", "severity", "ts"} {
		if _, ok := m[want]; !ok {
			t.Errorf("expected key %q to exist", want)
		}
	}
	for _, gone := range []string{"msg", "level"} {
		if _, ok := m[gone]; ok {
			t.Errorf("expected old key %q to be removed", gone)
		}
	}
}

func TestApply_SameKeyMapping_IsNoOp(t *testing.T) {
	r := jsonrename.New(map[string]string{"level": "level"})
	line := `{"level":"debug"}`
	got := r.Apply(line)
	m := decode(t, got)
	if _, ok := m["level"]; !ok {
		t.Error("expected 'level' to be preserved when old==new")
	}
}
