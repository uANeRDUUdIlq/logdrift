package filter

import (
	"testing"
)

func TestNew_ValidExpressions(t *testing.T) {
	f, err := New([]string{"level=error", "service=api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(f.Rules))
	}
}

func TestNew_InvalidExpression(t *testing.T) {
	_, err := New([]string{"badexpr"})
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestMatch_NoRules(t *testing.T) {
	f, _ := New(nil)
	if !f.Match([]byte(`{"level":"info"}`)) {
		t.Fatal("expected match with no rules")
	}
}

func TestMatch_Hit(t *testing.T) {
	f, _ := New([]string{"level=error"})
	if !f.Match([]byte(`{"level":"error","msg":"boom"}`)) {
		t.Fatal("expected match")
	}
}

func TestMatch_Miss(t *testing.T) {
	f, _ := New([]string{"level=error"})
	if f.Match([]byte(`{"level":"info","msg":"ok"}`)) {
		t.Fatal("expected no match")
	}
}

func TestMatch_MissingField(t *testing.T) {
	f, _ := New([]string{"service=api"})
	if f.Match([]byte(`{"level":"info"}`)) {
		t.Fatal("expected no match when field absent")
	}
}

func TestMatch_InvalidJSON(t *testing.T) {
	f, _ := New([]string{"level=error"})
	if f.Match([]byte(`not json`)) {
		t.Fatal("expected no match for invalid JSON")
	}
}

func TestMatch_CaseInsensitive(t *testing.T) {
	f, _ := New([]string{"level=ERROR"})
	if !f.Match([]byte(`{"level":"error"}`)) {
		t.Fatal("expected case-insensitive match")
	}
}
