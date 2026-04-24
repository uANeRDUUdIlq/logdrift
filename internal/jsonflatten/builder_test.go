package jsonflatten_test

import (
	"testing"

	"github.com/user/logdrift/internal/jsonflatten"
)

func TestBuild_Defaults_ReturnsFlattenner(t *testing.T) {
	f, err := jsonflatten.Build(jsonflatten.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Flattener")
	}
}

func TestBuild_NegativeMaxDepth_ReturnsError(t *testing.T) {
	_, err := jsonflatten.Build(jsonflatten.Config{MaxDepth: -1})
	if err == nil {
		t.Fatal("expected error for negative max_depth")
	}
}

func TestBuild_CustomSeparatorAndPrefix(t *testing.T) {
	f, err := jsonflatten.Build(jsonflatten.Config{Separator: "_", Prefix: "svc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"level":"info"}`
	got := decode(t, f.Apply(line))
	if _, ok := got["svc_level"]; !ok {
		t.Fatalf("key svc_level missing; got %v", got)
	}
}

func TestBuild_ZeroMaxDepth_IsValid(t *testing.T) {
	f, err := jsonflatten.Build(jsonflatten.Config{MaxDepth: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Flattener")
	}
}
