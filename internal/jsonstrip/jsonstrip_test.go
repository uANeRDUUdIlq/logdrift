package jsonstrip

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoKeys_ReturnsSameLine(t *testing.T) {
	s := New(nil)
	line := `{"level":"info","msg":"hello"}`
	if got := s.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	s := New([]string{"secret"})
	line := "not json at all"
	if got := s.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_RemovesTopLevelKey(t *testing.T) {
	s := New([]string{"secret"})
	line := `{"level":"info","secret":"abc","msg":"hello"}`
	got := decode(t, s.Apply(line))
	if _, ok := got["secret"]; ok {
		t.Fatal("expected 'secret' to be removed")
	}
	if got["msg"] != "hello" {
		t.Fatal("expected 'msg' to be preserved")
	}
}

func TestApply_KeyAbsent_NoChange(t *testing.T) {
	s := New([]string{"ghost"})
	line := `{"level":"info","msg":"hello"}`
	got := s.Apply(line)
	if decode(t, got)["msg"] != "hello" {
		t.Fatal("unexpected mutation")
	}
}

func TestApply_RemovesNestedKey(t *testing.T) {
	s := New([]string{"meta.internal"})
	line := `{"msg":"hi","meta":{"internal":true,"version":"1"}}`
	got := decode(t, s.Apply(line))
	meta, ok := got["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("expected meta to remain")
	}
	if _, ok := meta["internal"]; ok {
		t.Fatal("expected 'meta.internal' to be removed")
	}
	if meta["version"] != "1" {
		t.Fatal("expected 'meta.version' to be preserved")
	}
}

func TestApply_MultipleKeys_AllRemoved(t *testing.T) {
	s := New([]string{"token", "password"})
	line := `{"user":"alice","token":"xyz","password":"s3cr3t"}`
	got := decode(t, s.Apply(line))
	for _, k := range []string{"token", "password"} {
		if _, ok := got[k]; ok {
			t.Fatalf("expected %q to be removed", k)
		}
	}
	if got["user"] != "alice" {
		t.Fatal("expected 'user' to be preserved")
	}
}
