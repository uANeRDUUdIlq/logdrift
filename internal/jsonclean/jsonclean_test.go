package jsonclean_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/jsonclean"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoOptions_ReturnsSameLine(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{})
	line := `{"a":null,"b":"","c":{}}`
	if got := c.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{RemoveNull: true})
	line := "not json at all"
	if got := c.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_RemoveNull(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{RemoveNull: true})
	got := decode(t, c.Apply(`{"level":"info","err":null}`))
	if _, ok := got["err"]; ok {
		t.Error("expected 'err' field to be removed")
	}
	if got["level"] != "info" {
		t.Error("expected 'level' field to be preserved")
	}
}

func TestApply_RemoveEmptyString(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{RemoveEmptyString: true})
	got := decode(t, c.Apply(`{"msg":"hello","trace_id":""}`))
	if _, ok := got["trace_id"]; ok {
		t.Error("expected 'trace_id' field to be removed")
	}
	if got["msg"] != "hello" {
		t.Error("expected 'msg' field to be preserved")
	}
}

func TestApply_RemoveEmptyObject(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{RemoveEmptyObject: true})
	got := decode(t, c.Apply(`{"meta":{},"svc":"api"}`))
	if _, ok := got["meta"]; ok {
		t.Error("expected 'meta' field to be removed")
	}
}

func TestApply_RemoveEmptyArray(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{RemoveEmptyArray: true})
	got := decode(t, c.Apply(`{"tags":[],"level":"warn"}`))
	if _, ok := got["tags"]; ok {
		t.Error("expected 'tags' field to be removed")
	}
}

func TestApply_AllOptions_RemovesAllEmptyKinds(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{
		RemoveNull:        true,
		RemoveEmptyString: true,
		RemoveEmptyObject: true,
		RemoveEmptyArray:  true,
	})
	input := `{"a":null,"b":"","c":{},"d":[],"keep":"yes"}`
	got := decode(t, c.Apply(input))
	for _, field := range []string{"a", "b", "c", "d"} {
		if _, ok := got[field]; ok {
			t.Errorf("expected field %q to be removed", field)
		}
	}
	if got["keep"] != "yes" {
		t.Error("expected 'keep' field to be preserved")
	}
}
