package jsonmask_test

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/yourorg/logdrift/internal/jsonmask"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules_ReturnsSameLine(t *testing.T) {
	m := jsonmask.New(nil)
	line := `{"msg":"hello","token":"secret-abc"}`
	if got := m.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	m := jsonmask.New([]jsonmask.Rule{
		{Pattern: regexp.MustCompile(`\d+`), Mask: "[NUM]"},
	})
	line := "plain text log 12345"
	if got := m.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_MasksMatchingValue(t *testing.T) {
	m := jsonmask.New([]jsonmask.Rule{
		{Pattern: regexp.MustCompile(`tok_[a-z0-9]+`), Mask: "[TOKEN]"},
	})
	line := `{"service":"auth","token":"tok_abc123"}`
	got := m.Apply(line)
	obj := decode(t, got)
	if obj["token"] != "[TOKEN]" {
		t.Errorf("expected token to be masked, got %v", obj["token"])
	}
	if obj["service"] != "auth" {
		t.Errorf("expected service to be preserved, got %v", obj["service"])
	}
}

func TestApply_NoMatchingValue_ReturnsSameLine(t *testing.T) {
	m := jsonmask.New([]jsonmask.Rule{
		{Pattern: regexp.MustCompile(`tok_[a-z0-9]+`), Mask: "[TOKEN]"},
	})
	line := `{"msg":"nothing sensitive here"}`
	if got := m.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_DefaultMask_WhenMaskEmpty(t *testing.T) {
	m := jsonmask.New([]jsonmask.Rule{
		{Pattern: regexp.MustCompile(`\d{4}-\d{4}-\d{4}-\d{4}`)},
	})
	line := `{"cc":"1234-5678-9012-3456"}`
	got := m.Apply(line)
	obj := decode(t, got)
	if obj["cc"] != "[MASKED]" {
		t.Errorf("expected default mask, got %v", obj["cc"])
	}
}

func TestApply_MultipleRules_AllApplied(t *testing.T) {
	m := jsonmask.New([]jsonmask.Rule{
		{Pattern: regexp.MustCompile(`tok_\w+`), Mask: "[TOKEN]"},
		{Pattern: regexp.MustCompile(`\d{3}-\d{2}-\d{4}`), Mask: "[SSN]"},
	})
	line := `{"a":"tok_xyz","b":"123-45-6789"}`
	got := m.Apply(line)
	obj := decode(t, got)
	if obj["a"] != "[TOKEN]" {
		t.Errorf("expected a=[TOKEN], got %v", obj["a"])
	}
	if obj["b"] != "[SSN]" {
		t.Errorf("expected b=[SSN], got %v", obj["b"])
	}
}

func TestApply_NilPatternRule_Skipped(t *testing.T) {
	m := jsonmask.New([]jsonmask.Rule{
		{Pattern: nil, Mask: "[SHOULD NOT APPLY]"},
		{Pattern: regexp.MustCompile(`secret`), Mask: "[REDACTED]"},
	})
	line := `{"msg":"my secret value"}`
	got := m.Apply(line)
	obj := decode(t, got)
	if obj["msg"] != "my [REDACTED] value" {
		t.Errorf("unexpected msg: %v", obj["msg"])
	}
}

func TestApply_NestedObject_Unchanged(t *testing.T) {
	// Masking applies only to string values; nested objects should pass through intact.
	m := jsonmask.New([]jsonmask.Rule{
		{Pattern: regexp.MustCompile(`tok_\w+`), Mask: "[TOKEN]"},
	})
	line := `{"meta":{"user":"alice"},"token":"tok_xyz"}`
	got := m.Apply(line)
	obj := decode(t, got)
	if obj["token"] != "[TOKEN]" {
		t.Errorf("expected token masked, got %v", obj["token"])
	}
	meta, ok := obj["meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected meta to be an object, got %T", obj["meta"])
	}
	if meta["user"] != "alice" {
		t.Errorf("expected meta.user=alice, got %v", meta["user"])
	}
}
