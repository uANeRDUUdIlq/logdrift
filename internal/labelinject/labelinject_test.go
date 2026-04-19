package labelinject_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logdrift/internal/labelinject"
)

func decode(t *testing.T, s string) map[string]string {
	t.Helper()
	var m map[string]string
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoLabels_ReturnsSameLine(t *testing.T) {
	inj := labelinject.New(nil)
	line := `{"msg":"hello"}`
	if got := inj.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	inj := labelinject.New(map[string]string{"env": "prod"})
	line := "plain text log"
	if got := inj.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestApply_InjectsLabel(t *testing.T) {
	inj := labelinject.New(map[string]string{"env": "prod"})
	got := decode(t, inj.Apply(`{"msg":"hi"}`))
	if got["env"] != "prod" {
		t.Fatalf("expected env=prod, got %q", got["env"])
	}
}

func TestApply_DoesNotOverwriteExistingField(t *testing.T) {
	inj := labelinject.New(map[string]string{"env": "prod"})
	got := decode(t, inj.Apply(`{"env":"staging","msg":"hi"}`))
	if got["env"] != "staging" {
		t.Fatalf("existing field should not be overwritten, got %q", got["env"])
	}
}

func TestApply_MultipleLabels_AllInjected(t *testing.T) {
	inj := labelinject.New(map[string]string{"env": "prod", "region": "us-east-1"})
	got := decode(t, inj.Apply(`{"msg":"hi"}`))
	if got["env"] != "prod" {
		t.Fatalf("missing env label")
	}
	if got["region"] != "us-east-1" {
		t.Fatalf("missing region label")
	}
}
