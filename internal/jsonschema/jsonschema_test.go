package jsonschema

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules_ReturnsSameLine(t *testing.T) {
	v := New(nil, false, "")
	out, ok := v.Apply(`{"level":"info"}`)
	if !ok || out != `{"level":"info"}` {
		t.Fatalf("unexpected: ok=%v out=%q", ok, out)
	}
}

func TestApply_NonJSON_PassesThrough(t *testing.T) {
	v := New([]Rule{{Field: "level", Type: TypeString}}, true, "")
	out, ok := v.Apply("not json")
	if !ok || out != "not json" {
		t.Fatalf("unexpected: ok=%v out=%q", ok, out)
	}
}

func TestApply_ConformingLine_Passes(t *testing.T) {
	v := New([]Rule{{Field: "code", Type: TypeNumber}}, true, "")
	line := `{"code":200,"msg":"ok"}`
	out, ok := v.Apply(line)
	if !ok || out != line {
		t.Fatalf("unexpected: ok=%v out=%q", ok, out)
	}
}

func TestApply_ViolationDropOnFail_DropsLine(t *testing.T) {
	v := New([]Rule{{Field: "code", Type: TypeNumber}}, true, "")
	out, ok := v.Apply(`{"code":"two-hundred"}`)
	if ok || out != "" {
		t.Fatalf("expected drop, got ok=%v out=%q", ok, out)
	}
}

func TestApply_ViolationAnnotate_AddsErrorField(t *testing.T) {
	v := New([]Rule{{Field: "code", Type: TypeNumber}}, false, "_schema_error")
	out, ok := v.Apply(`{"code":"oops"}`)
	if !ok {
		t.Fatal("expected line to pass through")
	}
	m := decode(t, out)
	if _, exists := m["_schema_error"]; !exists {
		t.Fatalf("expected _schema_error field in %q", out)
	}
}

func TestApply_MissingField_IsIgnored(t *testing.T) {
	v := New([]Rule{{Field: "missing", Type: TypeBoolean}}, true, "")
	line := `{"level":"info"}`
	out, ok := v.Apply(line)
	if !ok || out != line {
		t.Fatalf("unexpected: ok=%v out=%q", ok, out)
	}
}

func TestApply_MultipleRules_AllViolationsAnnotated(t *testing.T) {
	v := New([]Rule{
		{Field: "code", Type: TypeNumber},
		{Field: "ok", Type: TypeBoolean},
	}, false, "_err")
	out, ok := v.Apply(`{"code":"x","ok":"yes"}`)
	if !ok {
		t.Fatal("expected pass-through")
	}
	m := decode(t, out)
	errVal, _ := m["_err"].(string)
	if errVal == "" {
		t.Fatalf("expected _err annotation, got %q", out)
	}
}
