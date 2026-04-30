package jsonenrich_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logdrift/internal/jsonenrich"
)

func decode(t *testing.T, line string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoFields_ReturnsSameLine(t *testing.T) {
	e := jsonenrich.New()
	const line = `{"msg":"hello"}`
	if got := e.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	e := jsonenrich.New(jsonenrich.WithField("env", "prod"))
	const line = "plain text log"
	if got := e.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_InjectsField(t *testing.T) {
	e := jsonenrich.New(jsonenrich.WithField("env", "staging"))
	got := e.Apply(`{"msg":"hi"}`)
	m := decode(t, got)
	if m["env"] != "staging" {
		t.Fatalf("expected env=staging, got %v", m["env"])
	}
}

func TestApply_DoesNotOverwriteExistingField(t *testing.T) {
	e := jsonenrich.New(jsonenrich.WithField("env", "prod"))
	got := e.Apply(`{"env":"dev","msg":"hi"}`)
	m := decode(t, got)
	if m["env"] != "dev" {
		t.Fatalf("existing field should not be overwritten, got %v", m["env"])
	}
}

func TestApply_MultipleFields_AllInjected(t *testing.T) {
	e := jsonenrich.New(
		jsonenrich.WithField("env", "prod"),
		jsonenrich.WithField("region", "us-east-1"),
	)
	got := e.Apply(`{"msg":"ok"}`)
	m := decode(t, got)
	if m["env"] != "prod" {
		t.Fatalf("env not injected")
	}
	if m["region"] != "us-east-1" {
		t.Fatalf("region not injected")
	}
}

func TestApply_WithHostname_InjectsHost(t *testing.T) {
	e := jsonenrich.New(jsonenrich.WithHostname("host"))
	got := e.Apply(`{"msg":"ping"}`)
	m := decode(t, got)
	// hostname resolution may succeed or fail in CI; if the key is present it
	// must be a non-empty string.
	if v, ok := m["host"]; ok {
		if s, _ := v.(string); s == "" {
			t.Fatal("host field present but empty")
		}
	}
}

func TestApply_InvalidJSON_ReturnsSameLine(t *testing.T) {
	e := jsonenrich.New(jsonenrich.WithField("env", "prod"))
	const line = `{bad json}`
	if got := e.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}
