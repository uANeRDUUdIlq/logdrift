package jsonclip_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/jsonclip"
)

func decode(t *testing.T, line string) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	c := jsonclip.New(3)
	input := "not json at all"
	if got := c.Apply(input); got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestApply_DisabledZeroMaxKeys_ReturnsSameLine(t *testing.T) {
	c := jsonclip.New(0)
	input := `{"a":1,"b":2,"c":3,"d":4}`
	if got := c.Apply(input); got != input {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_DisabledNegativeMaxKeys_ReturnsSameLine(t *testing.T) {
	c := jsonclip.New(-5)
	input := `{"a":1,"b":2}`
	if got := c.Apply(input); got != input {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_FewerKeysThanLimit_ReturnsSameLine(t *testing.T) {
	c := jsonclip.New(10)
	input := `{"x":1,"y":2}`
	if got := c.Apply(input); got != input {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_ExactLimit_ReturnsSameLine(t *testing.T) {
	c := jsonclip.New(3)
	input := `{"a":1,"b":2,"c":3}`
	if got := c.Apply(input); got != input {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_ExceedsLimit_DropsExtraKeys(t *testing.T) {
	c := jsonclip.New(2)
	input := `{"level":"info","msg":"hello","extra":"drop_me"}`
	got := c.Apply(input)
	m := decode(t, got)
	if len(m) != 2 {
		t.Fatalf("expected 2 keys, got %d in %q", len(m), got)
	}
	if _, ok := m["extra"]; ok {
		t.Fatalf("'extra' should have been dropped, got %q", got)
	}
}

func TestApply_LimitOne_KeepsFirstKey(t *testing.T) {
	c := jsonclip.New(1)
	input := `{"first":"keep","second":"gone"}`
	got := c.Apply(input)
	m := decode(t, got)
	if len(m) != 1 {
		t.Fatalf("expected 1 key, got %d in %q", len(m), got)
	}
	if _, ok := m["first"]; !ok {
		t.Fatalf("expected 'first' key to be present in %q", got)
	}
}
