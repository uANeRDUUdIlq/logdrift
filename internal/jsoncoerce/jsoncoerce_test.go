package jsoncoerce

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules_ReturnsSameLine(t *testing.T) {
	c := New(nil)
	line := `{"level":"info","code":200}`
	if got := c.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	c := New([]Rule{{Field: "code", Target: TypeString}})
	line := "not json at all"
	if got := c.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_NumberToString(t *testing.T) {
	c := New([]Rule{{Field: "code", Target: TypeString}})
	got := decode(t, c.Apply(`{"code":404}`))
	if v, ok := got["code"].(string); !ok || v != "404" {
		t.Errorf("expected string \"404\", got %v (%T)", got["code"], got["code"])
	}
}

func TestApply_StringToNumber(t *testing.T) {
	c := New([]Rule{{Field: "latency", Target: TypeNumber}})
	got := decode(t, c.Apply(`{"latency":"3.14"}`))
	if v, ok := got["latency"].(float64); !ok || v != 3.14 {
		t.Errorf("expected float64 3.14, got %v (%T)", got["latency"], got["latency"])
	}
}

func TestApply_StringToBool(t *testing.T) {
	c := New([]Rule{{Field: "ok", Target: TypeBool}})
	got := decode(t, c.Apply(`{"ok":"true"}`))
	if v, ok := got["ok"].(bool); !ok || !v {
		t.Errorf("expected bool true, got %v (%T)", got["ok"], got["ok"])
	}
}

func TestApply_BoolToNumber(t *testing.T) {
	c := New([]Rule{{Field: "flag", Target: TypeNumber}})
	got := decode(t, c.Apply(`{"flag":true}`))
	if v, ok := got["flag"].(float64); !ok || v != 1 {
		t.Errorf("expected float64 1, got %v (%T)", got["flag"], got["flag"])
	}
}

func TestApply_FieldAbsent_NoChange(t *testing.T) {
	c := New([]Rule{{Field: "missing", Target: TypeString}})
	line := `{"level":"info"}`
	if got := c.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_UncoercibleValue_LeftAsIs(t *testing.T) {
	c := New([]Rule{{Field: "tags", Target: TypeNumber}})
	line := `{"tags":["a","b"]}`
	if got := c.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_MultipleRules_AllApplied(t *testing.T) {
	c := New([]Rule{
		{Field: "code", Target: TypeString},
		{Field: "ok", Target: TypeBool},
	})
	got := decode(t, c.Apply(`{"code":200,"ok":"false"}`))
	if _, ok := got["code"].(string); !ok {
		t.Errorf("expected code to be string, got %T", got["code"])
	}
	if v, ok := got["ok"].(bool); !ok || v {
		t.Errorf("expected ok to be bool false, got %v (%T)", got["ok"], got["ok"])
	}
}
