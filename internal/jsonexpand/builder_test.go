package jsonexpand_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/jsonexpand"
)

func TestBuild_Defaults_ReturnsExpander(t *testing.T) {
	e, err := jsonexpand.Build(jsonexpand.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil expander")
	}
}

func TestBuild_NegativeMaxDepth_ReturnsError(t *testing.T) {
	_, err := jsonexpand.Build(jsonexpand.Config{MaxDepth: -1})
	if err == nil {
		t.Fatal("expected error for negative max_depth")
	}
}

func TestBuild_WithFields_ExpandsOnlyNamed(t *testing.T) {
	e, err := jsonexpand.Build(jsonexpand.Config{
		Fields:   []string{"http.method"},
		MaxDepth: 0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw := `{"http.method":"GET","other.key":"untouched"}`
	result := decode(t, e.Apply(raw))

	if result["other.key"] != "untouched" {
		t.Fatalf("other.key should remain, got %v", result["other.key"])
	}
	http, ok := result["http"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nested map under 'http'")
	}
	if http["method"] != "GET" {
		t.Fatalf("expected GET, got %v", http["method"])
	}
}

func TestBuild_ZeroMaxDepth_IsValid(t *testing.T) {
	_, err := jsonexpand.Build(jsonexpand.Config{MaxDepth: 0})
	if err != nil {
		t.Fatalf("zero max_depth should be valid, got: %v", err)
	}
}
