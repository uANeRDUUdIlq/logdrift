package jsonvalidate_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/jsonvalidate"
)

func TestAllow_NoRequired_ValidJSON_Passes(t *testing.T) {
	v := jsonvalidate.New(nil)
	if !v.Allow(`{"level":"info","msg":"hello"}`) {
		t.Fatal("expected valid JSON with no required fields to pass")
	}
}

func TestAllow_NoRequired_NonJSON_Dropped(t *testing.T) {
	v := jsonvalidate.New(nil)
	if v.Allow(`not json at all`) {
		t.Fatal("expected non-JSON to be dropped by default")
	}
}

func TestAllow_NonJSON_WithAllowNonJSON_Passes(t *testing.T) {
	v := jsonvalidate.New(nil, jsonvalidate.WithAllowNonJSON())
	if !v.Allow(`plain text log line`) {
		t.Fatal("expected non-JSON to pass when WithAllowNonJSON is set")
	}
}

func TestAllow_RequiredFieldPresent_Passes(t *testing.T) {
	v := jsonvalidate.New([]string{"level", "msg"})
	if !v.Allow(`{"level":"error","msg":"boom","svc":"api"}`) {
		t.Fatal("expected line with all required fields to pass")
	}
}

func TestAllow_RequiredFieldMissing_Dropped(t *testing.T) {
	v := jsonvalidate.New([]string{"level", "msg", "trace_id"})
	if v.Allow(`{"level":"info","msg":"ok"}`) {
		t.Fatal("expected line missing trace_id to be dropped")
	}
}

func TestAllow_EmptyJSON_NoRequired_Passes(t *testing.T) {
	v := jsonvalidate.New([]string{})
	if !v.Allow(`{}`) {
		t.Fatal("expected empty JSON object to pass when no fields required")
	}
}

func TestAllow_EmptyJSON_RequiredField_Dropped(t *testing.T) {
	v := jsonvalidate.New([]string{"service"})
	if v.Allow(`{}`) {
		t.Fatal("expected empty JSON object to be dropped when field required")
	}
}

func TestAllow_FieldValueNull_StillPresent_Passes(t *testing.T) {
	// A field with a null value is still considered present.
	v := jsonvalidate.New([]string{"trace_id"})
	if !v.Allow(`{"trace_id":null,"msg":"test"}`) {
		t.Fatal("expected null-valued required field to count as present")
	}
}
