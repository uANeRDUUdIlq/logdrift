package jsondefault_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logdrift/internal/jsondefault"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoDefaults_ReturnsSameLine(t *testing.T) {
	d := jsondefault.New(nil)
	line := `{"level":"info","msg":"hello"}`
	if got := d.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	d := jsondefault.New(map[string]any{"env": "prod"})
	line := "not json at all"
	if got := d.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_InjectsDefault(t *testing.T) {
	d := jsondefault.New(map[string]any{"env": "prod"})
	got := decode(t, d.Apply(`{"level":"info"}`))
	if got["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", got["env"])
	}
}

func TestApply_DoesNotOverwriteExistingField(t *testing.T) {
	d := jsondefault.New(map[string]any{"env": "prod"})
	got := decode(t, d.Apply(`{"env":"staging"}`))
	if got["env"] != "staging" {
		t.Errorf("expected env=staging to be preserved, got %v", got["env"])
	}
}

func TestApply_MultipleDefaults_AllInjected(t *testing.T) {
	d := jsondefault.New(map[string]any{"env": "prod", "region": "us-east-1"})
	got := decode(t, d.Apply(`{"level":"warn"}`))
	if got["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", got["env"])
	}
	if got["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %v", got["region"])
	}
}

func TestApply_PartialDefaults_OnlyMissingInjected(t *testing.T) {
	d := jsondefault.New(map[string]any{"env": "prod", "region": "us-east-1"})
	got := decode(t, d.Apply(`{"env":"dev"}`))
	if got["env"] != "dev" {
		t.Errorf("expected env=dev to be preserved, got %v", got["env"])
	}
	if got["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1 to be injected, got %v", got["region"])
	}
}

func TestApply_NumericDefault(t *testing.T) {
	d := jsondefault.New(map[string]any{"retries": 3})
	got := decode(t, d.Apply(`{"level":"error"}`))
	if got["retries"] != float64(3) {
		t.Errorf("expected retries=3, got %v", got["retries"])
	}
}
