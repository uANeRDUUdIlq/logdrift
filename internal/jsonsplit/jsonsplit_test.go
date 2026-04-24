package jsonsplit_test

import (
	"encoding/json"
	"testing"

	"logdrift/internal/jsonsplit"
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
	_, err := jsonsplit.New("", "")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_ValidField_NoError(t *testing.T) {
	_, err := jsonsplit.New("items", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	s, _ := jsonsplit.New("items", "")
	result := s.Apply("not json")
	if len(result) != 1 || result[0] != "not json" {
		t.Fatalf("expected passthrough, got %v", result)
	}
}

func TestApply_FieldAbsent_ReturnsSameLine(t *testing.T) {
	s, _ := jsonsplit.New("items", "")
	line := `{"service":"api"}`
	result := s.Apply(line)
	if len(result) != 1 || result[0] != line {
		t.Fatalf("expected passthrough, got %v", result)
	}
}

func TestApply_FieldNotArray_ReturnsSameLine(t *testing.T) {
	s, _ := jsonsplit.New("items", "")
	line := `{"items":"scalar"}`
	result := s.Apply(line)
	if len(result) != 1 || result[0] != line {
		t.Fatalf("expected passthrough, got %v", result)
	}
}

func TestApply_ArrayField_SplitsIntoMultipleLines(t *testing.T) {
	s, _ := jsonsplit.New("items", "")
	line := `{"service":"api","items":["a","b","c"]}`
	result := s.Apply(line)
	if len(result) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result))
	}
	for _, r := range result {
		m := decode(t, r)
		if m["service"] != "api" {
			t.Errorf("base field missing in %s", r)
		}
		if _, ok := m["items"]; !ok {
			t.Errorf("element field missing in %s", r)
		}
	}
}

func TestApply_WithPrefix_UsesPrefix(t *testing.T) {
	s, _ := jsonsplit.New("errors", "error")
	line := `{"svc":"x","errors":[{"code":1},{"code":2}]}`
	result := s.Apply(line)
	if len(result) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(result))
	}
	for _, r := range result {
		m := decode(t, r)
		if _, ok := m["error"]; !ok {
			t.Errorf("prefix key 'error' missing in %s", r)
		}
		if _, ok := m["errors"]; ok {
			t.Errorf("original key 'errors' should be absent in %s", r)
		}
	}
}

func TestApply_EmptyArray_ReturnsSameLine(t *testing.T) {
	s, _ := jsonsplit.New("items", "")
	line := `{"service":"api","items":[]}`
	result := s.Apply(line)
	if len(result) != 1 || result[0] != line {
		t.Fatalf("expected passthrough for empty array, got %v", result)
	}
}
