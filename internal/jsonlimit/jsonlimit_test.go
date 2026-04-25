package jsonlimit_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logdrift/internal/jsonlimit"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	l := jsonlimit.New(3, nil)
	got := l.Apply("not json")
	if got != "not json" {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestApply_NoLimits_ReturnsSameLine(t *testing.T) {
	l := jsonlimit.New(0, nil)
	input := `{"items":[1,2,3,4,5]}`
	got := l.Apply(input)
	if got != input {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestApply_DefaultMax_TruncatesArrays(t *testing.T) {
	l := jsonlimit.New(2, nil)
	got := l.Apply(`{"items":[1,2,3,4]}`)
	m := decode(t, got)
	arr, ok := m["items"].([]interface{})
	if !ok {
		t.Fatal("items is not an array")
	}
	if len(arr) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(arr))
	}
}

func TestApply_FieldSpecific_OverridesDefault(t *testing.T) {
	l := jsonlimit.New(2, map[string]int{"tags": 4})
	got := l.Apply(`{"tags":["a","b","c","d","e"],"items":[1,2,3]}`)
	m := decode(t, got)

	tags := m["tags"].([]interface{})
	if len(tags) != 4 {
		t.Fatalf("tags: expected 4, got %d", len(tags))
	}

	items := m["items"].([]interface{})
	if len(items) != 2 {
		t.Fatalf("items: expected 2 (default), got %d", len(items))
	}
}

func TestApply_ArrayAlreadyWithinLimit_Unchanged(t *testing.T) {
	l := jsonlimit.New(10, nil)
	input := `{"items":[1,2,3]}`
	got := l.Apply(input)
	m := decode(t, got)
	if len(m["items"].([]interface{})) != 3 {
		t.Fatal("array should not be truncated")
	}
}

func TestApply_NonArrayField_Untouched(t *testing.T) {
	l := jsonlimit.New(2, nil)
	input := `{"level":"info","msg":"hello"}`
	got := l.Apply(input)
	m := decode(t, got)
	if m["level"] != "info" || m["msg"] != "hello" {
		t.Fatalf("non-array fields should be untouched, got %v", m)
	}
}

func TestApply_NegativeDefaultMax_TreatedAsZero(t *testing.T) {
	l := jsonlimit.New(-5, nil)
	input := `{"items":[1,2,3,4,5]}`
	got := l.Apply(input)
	if got != input {
		t.Fatalf("expected unchanged, got %q", got)
	}
}
