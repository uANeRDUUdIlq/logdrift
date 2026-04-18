package transform_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/transform"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	tr := transform.New(transform.Rule{})
	got := tr.Apply("plain text")
	if got != "plain text" {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestApply_RenameField(t *testing.T) {
	tr := transform.New(transform.Rule{
		Rename: map[string]string{"msg": "message"},
	})
	out := decode(t, tr.Apply(`{"msg":"hello"}`))
	if out["message"] != "hello" {
		t.Fatalf("expected message=hello, got %v", out)
	}
	if _, ok := out["msg"]; ok {
		t.Fatal("old key should be removed")
	}
}

func TestApply_AddField(t *testing.T) {
	tr := transform.New(transform.Rule{
		Add: map[string]string{"env": "production"},
	})
	out := decode(t, tr.Apply(`{"level":"info"}`))
	if out["env"] != "production" {
		t.Fatalf("expected env=production, got %v", out)
	}
}

func TestApply_RenameNonExistentKey_NoOp(t *testing.T) {
	tr := transform.New(transform.Rule{
		Rename: map[string]string{"missing": "new"},
	})
	out := decode(t, tr.Apply(`{"level":"warn"}`))
	if _, ok := out["new"]; ok {
		t.Fatal("should not add key for missing source")
	}
}

func TestApply_EmptyRules_ReturnsSameJSON(t *testing.T) {
	tr := transform.New(transform.Rule{})
	out := decode(t, tr.Apply(`{"level":"debug","msg":"ok"}`))
	if out["level"] != "debug" {
		t.Fatalf("unexpected mutation: %v", out)
	}
}
