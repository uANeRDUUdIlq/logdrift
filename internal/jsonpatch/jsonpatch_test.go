package jsonpatch_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logdrift/internal/jsonpatch"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoFields_ReturnsSameLine(t *testing.T) {
	p, err := jsonpatch.New(nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	line := `{"level":"info","msg":"hello"}`
	if got := p.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	p, _ := jsonpatch.New(map[string]string{"env": "prod"})
	line := "not json at all"
	if got := p.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_AddsNewField(t *testing.T) {
	p, _ := jsonpatch.New(map[string]string{"env": "staging"})
	line := `{"level":"info"}`
	got := decode(t, p.Apply(line))
	if got["env"] != "staging" {
		t.Errorf("expected env=staging, got %v", got["env"])
	}
}

func TestApply_OverwritesExistingField(t *testing.T) {
	p, _ := jsonpatch.New(map[string]string{"level": "error"})
	line := `{"level":"debug","msg":"test"}`
	got := decode(t, p.Apply(line))
	if got["level"] != "error" {
		t.Errorf("expected level=error, got %v", got["level"])
	}
	if got["msg"] != "test" {
		t.Errorf("msg field should be preserved")
	}
}

func TestApply_MultipleFields_AllApplied(t *testing.T) {
	p, _ := jsonpatch.New(map[string]string{"env": "prod", "region": "us-east-1"})
	line := `{"msg":"ok"}`
	got := decode(t, p.Apply(line))
	if got["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", got["env"])
	}
	if got["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %v", got["region"])
	}
}

func TestApply_PreservesUnpatchedFields(t *testing.T) {
	p, _ := jsonpatch.New(map[string]string{"env": "prod"})
	line := `{"service":"api","level":"warn"}`
	got := decode(t, p.Apply(line))
	if got["service"] != "api" {
		t.Errorf("service field should be preserved")
	}
	if got["level"] != "warn" {
		t.Errorf("level field should be preserved")
	}
}
