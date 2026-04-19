package jsonpath_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/jsonpath"
)

func TestGet_NonJSON_ReturnsFalse(t *testing.T) {
	e := jsonpath.New()
	_, ok := e.Get("not json", "field")
	if ok {
		t.Fatal("expected false for non-JSON input")
	}
}

func TestGet_EmptyPath_ReturnsFalse(t *testing.T) {
	e := jsonpath.New()
	_, ok := e.Get(`{"a":"b"}`, "")
	if ok {
		t.Fatal("expected false for empty path")
	}
}

func TestGet_TopLevelField(t *testing.T) {
	e := jsonpath.New()
	v, ok := e.Get(`{"level":"info"}`, "level")
	if !ok || v != "info" {
		t.Fatalf("expected 'info', got %q ok=%v", v, ok)
	}
}

func TestGet_NestedField(t *testing.T) {
	e := jsonpath.New()
	v, ok := e.Get(`{"meta":{"service":"auth"}}`, "meta.service")
	if !ok || v != "auth" {
		t.Fatalf("expected 'auth', got %q ok=%v", v, ok)
	}
}

func TestGet_MissingField_ReturnsFalse(t *testing.T) {
	e := jsonpath.New()
	_, ok := e.Get(`{"level":"info"}`, "missing")
	if ok {
		t.Fatal("expected false for missing field")
	}
}

func TestGet_NumericField(t *testing.T) {
	e := jsonpath.New()
	v, ok := e.Get(`{"code":200}`, "code")
	if !ok || v != "200" {
		t.Fatalf("expected '200', got %q ok=%v", v, ok)
	}
}

func TestGet_BoolField(t *testing.T) {
	e := jsonpath.New()
	v, ok := e.Get(`{"ok":true}`, "ok")
	if !ok || v != "true" {
		t.Fatalf("expected 'true', got %q ok=%v", v, ok)
	}
}

func TestGet_NullField(t *testing.T) {
	e := jsonpath.New()
	v, ok := e.Get(`{"x":null}`, "x")
	if !ok || v != "null" {
		t.Fatalf("expected 'null', got %q ok=%v", v, ok)
	}
}

func TestGet_DeepNesting_MissingIntermediate(t *testing.T) {
	e := jsonpath.New()
	_, ok := e.Get(`{"a":{"b":"c"}}`, "a.z.deep")
	if ok {
		t.Fatal("expected false for missing intermediate key")
	}
}
