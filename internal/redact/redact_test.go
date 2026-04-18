package redact_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/redact"
)

func TestApply_NoFields_ReturnsSameLine(t *testing.T) {
	r := redact.New(nil, "")
	line := `{"level":"info","msg":"hello"}`
	if got := r.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	r := redact.New([]string{"password"}, "")
	line := "plain text log line"
	if got := r.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_RedactsMatchingField(t *testing.T) {
	r := redact.New([]string{"password"}, "***")
	line := `{"user":"alice","password":"secret"}`
	out := r.Apply(line)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if obj["password"] != "***" {
		t.Fatalf("expected password redacted, got %v", obj["password"])
	}
	if obj["user"] != "alice" {
		t.Fatalf("expected user preserved, got %v", obj["user"])
	}
}

func TestApply_CaseInsensitiveField(t *testing.T) {
	r := redact.New([]string{"Authorization"}, "")
	line := `{"authorization":"Bearer token123"}`
	out := r.Apply(line)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if obj["authorization"] != "[REDACTED]" {
		t.Fatalf("expected redaction, got %v", obj["authorization"])
	}
}

func TestApply_DefaultMask(t *testing.T) {
	r := redact.New([]string{"token"}, "")
	line := `{"token":"abc"}`
	out := r.Apply(line)

	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["token"] != "[REDACTED]" {
		t.Fatalf("expected default mask, got %v", obj["token"])
	}
}

func TestApply_NoMatchingField_ReturnsSameLine(t *testing.T) {
	r := redact.New([]string{"secret"}, "")
	line := `{"level":"debug","msg":"ok"}`
	out := r.Apply(line)
	if out != line {
		t.Fatalf("expected unchanged, got %q", out)
	}
}
