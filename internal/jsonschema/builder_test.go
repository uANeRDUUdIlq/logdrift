package jsonschema

import "testing"

func TestBuild_Valid_ReturnsValidator(t *testing.T) {
	v, err := Build(Config{
		Rules:      []Rule{{Field: "level", Type: TypeString}},
		DropOnFail: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected non-nil validator")
	}
}

func TestBuild_EmptyField_ReturnsError(t *testing.T) {
	_, err := Build(Config{
		Rules: []Rule{{Field: "", Type: TypeString}},
	})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestBuild_UnknownType_ReturnsError(t *testing.T) {
	_, err := Build(Config{
		Rules: []Rule{{Field: "x", Type: "uuid"}},
	})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestBuild_NoRules_Valid(t *testing.T) {
	v, err := Build(Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"msg":"hello"}`
	out, ok := v.Apply(line)
	if !ok || out != line {
		t.Fatalf("expected pass-through, got ok=%v out=%q", ok, out)
	}
}

func TestBuild_DropOnFail_WorksEndToEnd(t *testing.T) {
	v, err := Build(Config{
		Rules:      []Rule{{Field: "status", Type: TypeNumber}},
		DropOnFail: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := v.Apply(`{"status":"bad"}`)
	if ok {
		t.Fatal("expected line to be dropped")
	}
}
