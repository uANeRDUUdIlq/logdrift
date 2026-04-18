package highlight_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logdrift/internal/highlight"
)

func TestNew_NoRules_ReturnsHighlighter(t *testing.T) {
	h := highlight.New(nil)
	if h == nil {
		t.Fatal("expected non-nil Highlighter")
	}
}

func TestApply_NoRules_ReturnsSameLine(t *testing.T) {
	h := highlight.New(nil)
	line := "hello world"
	if got := h.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_KnownColor_ContainsKeyword(t *testing.T) {
	rules := []highlight.Rule{{Word: "ERROR", Color: "red"}}
	h := highlight.New(rules)
	result := h.Apply("level=ERROR msg=boom")
	if !strings.Contains(result, "ERROR") {
		t.Errorf("expected result to contain ERROR, got %q", result)
	}
}

func TestApply_UnknownColor_FallsBackToYellow(t *testing.T) {
	rules := []highlight.Rule{{Word: "WARN", Color: "neon"}}
	h := highlight.New(rules)
	result := h.Apply("level=WARN msg=watch out")
	if !strings.Contains(result, "WARN") {
		t.Errorf("expected result to contain WARN, got %q", result)
	}
}

func TestApply_MultipleRules_AllReplaced(t *testing.T) {
	rules := []highlight.Rule{
		{Word: "ERROR", Color: "red"},
		{Word: "timeout", Color: "yellow"},
	}
	h := highlight.New(rules)
	result := h.Apply("ERROR: timeout occurred")
	if !strings.Contains(result, "ERROR") {
		t.Error("expected ERROR in result")
	}
	if !strings.Contains(result, "timeout") {
		t.Error("expected timeout in result")
	}
}

func TestApply_NoMatch_ReturnsSameLine(t *testing.T) {
	rules := []highlight.Rule{{Word: "FATAL", Color: "red"}}
	h := highlight.New(rules)
	line := "level=INFO msg=all good"
	if got := h.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}
