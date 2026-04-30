package jsoncast

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

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ValidField_ReturnsNonNil(t *testing.T) {
	c, err := New("items")
	if err != nil || c == nil {
		t.Fatalf("unexpected: err=%v c=%v", err, c)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	c, _ := New("items")
	out := c.Apply("not json")
	if len(out) != 1 || out[0] != "not json" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestApply_FieldAbsent_ReturnsSameLine(t *testing.T) {
	c, _ := New("items")
	line := `{"service":"api"}`
	out := c.Apply(line)
	if len(out) != 1 || out[0] != line {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestApply_FieldNotArray_ReturnsSameLine(t *testing.T) {
	c, _ := New("items")
	line := `{"items":"scalar"}`
	out := c.Apply(line)
	if len(out) != 1 || out[0] != line {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestApply_ObjectElements_MergedIntoBase(t *testing.T) {
	c, _ := New("items")
	line := `{"service":"api","items":[{"id":1},{"id":2}]}`
	out := c.Apply(line)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	for i, o := range out {
		m := decode(t, o)
		if m["service"] != "api" {
			t.Errorf("line %d: missing service field", i)
		}
		if _, ok := m["id"]; !ok {
			t.Errorf("line %d: missing id field", i)
		}
		if _, ok := m["items"]; ok {
			t.Errorf("line %d: items field should be absent", i)
		}
	}
}

func TestApply_ScalarElements_StoredUnderFieldName(t *testing.T) {
	c, _ := New("tags")
	line := `{"service":"api","tags":["alpha","beta"]}`
	out := c.Apply(line)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	m0 := decode(t, out[0])
	if m0["tags"] != "alpha" {
		t.Errorf("expected tags=alpha, got %v", m0["tags"])
	}
	m1 := decode(t, out[1])
	if m1["tags"] != "beta" {
		t.Errorf("expected tags=beta, got %v", m1["tags"])
	}
}

func TestApply_EmptyArray_ReturnsSameLine(t *testing.T) {
	c, _ := New("items")
	line := `{"service":"api","items":[]}`
	out := c.Apply(line)
	if len(out) != 1 || out[0] != line {
		t.Fatalf("expected passthrough for empty array, got %v", out)
	}
}
