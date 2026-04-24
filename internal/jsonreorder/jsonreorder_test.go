package jsonreorder_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/jsonreorder"
)

func decode(t *testing.T, s string) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func keys(t *testing.T, s string) []string {
	t.Helper()
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		t.Fatalf("keys: %v", err)
	}
	// re-decode preserving insertion order via a custom approach
	type kv struct{ K string }
	// Use json.Decoder token approach to get ordered keys
	dec := json.NewDecoder(strings.NewReader(s))
	dec.Token() // '{'
	out := []string{}
	for dec.More() {
		tok, _ := dec.Token()
		out = append(out, tok.(string))
		var skip json.RawMessage
		dec.Decode(&skip)
	}
	return out
}

import "strings"

func TestApply_NoKeys_ReturnsSameLine(t *testing.T) {
	r := jsonreorder.New(nil)
	line := `{"z":1,"a":2}`
	if got := r.Apply(line); got != line {
		t.Fatalf("expected unchanged, got %s", got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	r := jsonreorder.New([]string{"ts", "level"})
	line := "not json at all"
	if got := r.Apply(line); got != line {
		t.Fatalf("expected unchanged, got %s", got)
	}
}

func TestApply_PriorityKeysFirst(t *testing.T) {
	r := jsonreorder.New([]string{"ts", "level", "msg"})
	line := `{"msg":"hello","service":"api","ts":"2024-01-01","level":"info"}`
	got := r.Apply(line)
	ks := keys(t, got)
	if ks[0] != "ts" || ks[1] != "level" || ks[2] != "msg" {
		t.Fatalf("unexpected key order: %v", ks)
	}
}

func TestApply_UnspecifiedKeysAppended(t *testing.T) {
	r := jsonreorder.New([]string{"ts"})
	line := `{"z":99,"a":1,"ts":"now"}`
	got := r.Apply(line)
	ks := keys(t, got)
	if ks[0] != "ts" {
		t.Fatalf("expected ts first, got %v", ks)
	}
	if len(ks) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(ks))
	}
}

func TestApply_MissingPriorityKey_Skipped(t *testing.T) {
	r := jsonreorder.New([]string{"ts", "level", "msg"})
	line := `{"msg":"hi","code":200}`
	got := r.Apply(line)
	m := decode(t, got)
	if _, ok := m["ts"]; ok {
		t.Fatal("ts should not be present")
	}
	if _, ok := m["msg"]; !ok {
		t.Fatal("msg should be present")
	}
}

func TestApply_AllFieldsAccountedFor(t *testing.T) {
	r := jsonreorder.New([]string{"a", "b"})
	line := `{"c":3,"a":1,"b":2}`
	got := r.Apply(line)
	m := decode(t, got)
	if len(m) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(m))
	}
}
