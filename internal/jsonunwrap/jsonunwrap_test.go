package jsonunwrap_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logdrift/internal/jsonunwrap"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoFields_ReturnsSameLine(t *testing.T) {
	u := jsonunwrap.New(nil)
	line := `{"meta":{"host":"web-1"},"msg":"ok"}`
	if got := u.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	u := jsonunwrap.New([]string{"meta"})
	line := "not json at all"
	if got := u.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_UnwrapsNestedField(t *testing.T) {
	u := jsonunwrap.New([]string{"meta"})
	line := `{"meta":{"host":"web-1","env":"prod"},"msg":"ok"}`
	got := decode(t, u.Apply(line))

	if _, ok := got["meta"]; ok {
		t.Fatal("expected 'meta' key to be removed")
	}
	if got["host"] != "web-1" {
		t.Fatalf("expected host=web-1, got %v", got["host"])
	}
	if got["env"] != "prod" {
		t.Fatalf("expected env=prod, got %v", got["env"])
	}
	if got["msg"] != "ok" {
		t.Fatalf("expected msg=ok, got %v", got["msg"])
	}
}

func TestApply_DoesNotOverwriteExistingField(t *testing.T) {
	u := jsonunwrap.New([]string{"meta"})
	// top-level "host" should win over meta.host
	line := `{"host":"original","meta":{"host":"nested","env":"prod"},"msg":"ok"}`
	got := decode(t, u.Apply(line))

	if got["host"] != "original" {
		t.Fatalf("expected host=original, got %v", got["host"])
	}
	if got["env"] != "prod" {
		t.Fatalf("expected env=prod, got %v", got["env"])
	}
}

func TestApply_MissingField_NoChange(t *testing.T) {
	u := jsonunwrap.New([]string{"meta"})
	line := `{"msg":"ok"}`
	if got := u.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_NestedFieldNotObject_NoChange(t *testing.T) {
	u := jsonunwrap.New([]string{"meta"})
	line := `{"meta":"string-not-object","msg":"ok"}`
	got := decode(t, u.Apply(line))
	// meta should still be present since it wasn't an object
	if got["meta"] != "string-not-object" {
		t.Fatalf("expected meta to remain, got %v", got["meta"])
	}
}

func TestApply_MultipleFields_AllUnwrapped(t *testing.T) {
	u := jsonunwrap.New([]string{"meta", "tags"})
	line := `{"meta":{"host":"web-1"},"tags":{"region":"us-east"},"msg":"ok"}`
	got := decode(t, u.Apply(line))

	for _, key := range []string{"meta", "tags"} {
		if _, ok := got[key]; ok {
			t.Fatalf("expected %q to be removed", key)
		}
	}
	if got["host"] != "web-1" {
		t.Fatalf("expected host=web-1, got %v", got["host"])
	}
	if got["region"] != "us-east" {
		t.Fatalf("expected region=us-east, got %v", got["region"])
	}
}
