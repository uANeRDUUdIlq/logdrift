package jsonexpand_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/jsonexpand"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("invalid JSON %q: %v", line, err)
	}
	return m
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	e := jsonexpand.New(nil)
	const raw = "not json at all"
	if got := e.Apply(raw); got != raw {
		t.Fatalf("expected %q, got %q", raw, got)
	}
}

func TestApply_NoDotsNoFields_Unchanged(t *testing.T) {
	e := jsonexpand.New(nil)
	const raw = `{"level":"info","msg":"hello"}`
	got := decode(t, e.Apply(raw))
	if got["level"] != "info" || got["msg"] != "hello" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestApply_AutoExpand_DotKey(t *testing.T) {
	e := jsonexpand.New(nil) // expand all dot-keys
	raw := `{"a.b":"hello","level":"info"}`
	result := decode(t, e.Apply(raw))

	if _, exists := result["a.b"]; exists {
		t.Fatal("dot key should have been removed from top level")
	}
	aMap, ok := result["a"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nested map under 'a', got %T", result["a"])
	}
	if aMap["b"] != "hello" {
		t.Fatalf("expected 'hello', got %v", aMap["b"])
	}
}

func TestApply_DeepNesting(t *testing.T) {
	e := jsonexpand.New(nil)
	raw := `{"a.b.c":"deep"}`
	result := decode(t, e.Apply(raw))

	aMap := result["a"].(map[string]interface{})
	bMap := aMap["b"].(map[string]interface{})
	if bMap["c"] != "deep" {
		t.Fatalf("expected 'deep', got %v", bMap["c"])
	}
}

func TestApply_SelectedFields_OnlyExpandsNamed(t *testing.T) {
	e := jsonexpand.New([]string{"meta.env"})
	raw := `{"meta.env":"prod","other.key":"stays"}`
	result := decode(t, e.Apply(raw))

	if _, exists := result["meta.env"]; exists {
		t.Fatal("meta.env should have been expanded")
	}
	if result["other.key"] != "stays" {
		t.Fatalf("other.key should remain unexpanded, got %v", result["other.key"])
	}
	meta := result["meta"].(map[string]interface{})
	if meta["env"] != "prod" {
		t.Fatalf("expected 'prod', got %v", meta["env"])
	}
}

func TestApply_NonStringValue_NotExpanded(t *testing.T) {
	e := jsonexpand.New(nil)
	raw := `{"a.b":42}`
	result := decode(t, e.Apply(raw))
	// numeric value under a dot-key should remain as-is
	if result["a.b"] != float64(42) {
		t.Fatalf("expected 42 to remain unexpanded, got %v", result["a.b"])
	}
}
