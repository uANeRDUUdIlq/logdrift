package jsontruncate_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/logdrift/internal/jsontruncate"
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
	tr := jsontruncate.New(10)
	const line = "not json at all"
	if got := tr.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_DisabledZeroMaxLen_ReturnsSameLine(t *testing.T) {
	tr := jsontruncate.New(0)
	line := `{"msg":"this is a very long message that should not be touched"}`
	if got := tr.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_ShortValue_NotTruncated(t *testing.T) {
	tr := jsontruncate.New(50)
	line := `{"msg":"hello"}`
	got := tr.Apply(line)
	m := decode(t, got)
	if m["msg"] != "hello" {
		t.Fatalf("expected 'hello', got %v", m["msg"])
	}
}

func TestApply_LongValue_Truncated(t *testing.T) {
	tr := jsontruncate.New(5)
	line := `{"msg":"hello world"}`
	got := tr.Apply(line)
	m := decode(t, got)
	v, _ := m["msg"].(string)
	if !strings.HasPrefix(v, "hello") {
		t.Fatalf("expected prefix 'hello', got %q", v)
	}
	if !strings.HasSuffix(v, "...[truncated]") {
		t.Fatalf("expected truncation suffix, got %q", v)
	}
}

func TestApply_CustomSuffix(t *testing.T) {
	tr := jsontruncate.New(4, jsontruncate.WithSuffix("…"))
	line := `{"msg":"hello world"}`
	got := tr.Apply(line)
	m := decode(t, got)
	v, _ := m["msg"].(string)
	if !strings.HasSuffix(v, "…") {
		t.Fatalf("expected '…' suffix, got %q", v)
	}
}

func TestApply_WithFields_OnlyNamedFieldTruncated(t *testing.T) {
	tr := jsontruncate.New(5, jsontruncate.WithFields("msg"))
	long := "abcdefghij"
	line, _ := json.Marshal(map[string]interface{}{"msg": long, "other": long})
	got := tr.Apply(string(line))
	m := decode(t, got)

	msg, _ := m["msg"].(string)
	if !strings.HasSuffix(msg, "...[truncated]") {
		t.Fatalf("expected msg truncated, got %q", msg)
	}
	other, _ := m["other"].(string)
	if other != long {
		t.Fatalf("expected other unchanged, got %q", other)
	}
}

func TestApply_NonStringValue_Unchanged(t *testing.T) {
	tr := jsontruncate.New(2)
	line := `{"count":12345}`
	got := tr.Apply(line)
	m := decode(t, got)
	if m["count"].(float64) != 12345 {
		t.Fatalf("expected count 12345, got %v", m["count"])
	}
}
