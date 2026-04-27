package jsonpivot

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyArray_ReturnsError(t *testing.T) {
	_, err := New("", "name", "value")
	if err == nil {
		t.Fatal("expected error for empty array field")
	}
}

func TestNew_EmptyKeyField_ReturnsError(t *testing.T) {
	_, err := New("items", "", "value")
	if err == nil {
		t.Fatal("expected error for empty key field")
	}
}

func TestNew_EmptyValField_ReturnsError(t *testing.T) {
	_, err := New("items", "name", "")
	if err == nil {
		t.Fatal("expected error for empty value field")
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	p, _ := New("metrics", "name", "value")
	const line = "not json at all"
	if got := p.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_MissingArray_ReturnsSameLine(t *testing.T) {
	p, _ := New("metrics", "name", "value")
	const line = `{"svc":"api"}`
	if got := p.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_PivotsArray(t *testing.T) {
	p, _ := New("tags", "k", "v")
	line := `{"svc":"api","tags":[{"k":"env","v":"prod"},{"k":"region","v":"us-east"}]}`
	out := p.Apply(line)
	m := decode(t, out)

	if _, hasArr := m["tags"]; hasArr {
		t.Error("pivot array field should be removed from output")
	}
	if m["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", m["env"])
	}
	if m["region"] != "us-east" {
		t.Errorf("expected region=us-east, got %v", m["region"])
	}
	if m["svc"] != "api" {
		t.Errorf("expected svc=api, got %v", m["svc"])
	}
}

func TestApply_NumericValue_Preserved(t *testing.T) {
	p, _ := New("metrics", "name", "value")
	line := `{"metrics":[{"name":"cpu","value":0.42}]}`
	out := p.Apply(line)
	m := decode(t, out)
	if m["cpu"] != 0.42 {
		t.Errorf("expected cpu=0.42, got %v", m["cpu"])
	}
}

func TestApply_ItemMissingKeyField_Skipped(t *testing.T) {
	p, _ := New("items", "name", "val")
	line := `{"items":[{"val":"x"},{"name":"ok","val":"y"}]}`
	out := p.Apply(line)
	m := decode(t, out)
	if m["ok"] != "y" {
		t.Errorf("expected ok=y, got %v", m["ok"])
	}
	// entry without "name" should not produce any extra key beyond "ok"
	if len(m) != 1 {
		t.Errorf("expected 1 key, got %d: %v", len(m), m)
	}
}
