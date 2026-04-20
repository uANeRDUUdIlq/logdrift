package jsonmerge_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/jsonmerge"
)

func decode(t *testing.T, s string) map[string]string {
	t.Helper()
	var m map[string]string
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoFields_ReturnsSameLine(t *testing.T) {
	m := jsonmerge.New(nil)
	line := `{"msg":"hello"}`
	if got := m.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	m := jsonmerge.New(map[string]string{"env": "prod"})
	line := "not json at all"
	if got := m.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_InjectsField(t *testing.T) {
	m := jsonmerge.New(map[string]string{"env": "prod"})
	got := decode(t, m.Apply(`{"msg":"hello"}`))
	if got["env"] != "prod" {
		t.Fatalf("expected env=prod, got %q", got["env"])
	}
	if got["msg"] != "hello" {
		t.Fatalf("original field missing")
	}
}

func TestApply_DoesNotOverwriteExistingField(t *testing.T) {
	m := jsonmerge.New(map[string]string{"env": "prod"})
	got := decode(t, m.Apply(`{"env":"staging","msg":"hi"}`))
	if got["env"] != "staging" {
		t.Fatalf("existing field was overwritten; got %q", got["env"])
	}
}

func TestApply_MultipleFields_AllInjected(t *testing.T) {
	m := jsonmerge.New(map[string]string{"env": "prod", "region": "us-east-1"})
	got := decode(t, m.Apply(`{"msg":"ok"}`))
	if got["env"] != "prod" {
		t.Fatalf("env missing")
	}
	if got["region"] != "us-east-1" {
		t.Fatalf("region missing")
	}
}

func TestNew_FieldsAreCopied(t *testing.T) {
	src := map[string]string{"env": "prod"}
	m := jsonmerge.New(src)
	src["env"] = "mutated"
	got := decode(t, m.Apply(`{"msg":"x"}`))
	if got["env"] != "prod" {
		t.Fatalf("merger was affected by mutation of source map")
	}
}
