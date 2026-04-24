package jsonsorter

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func keys(t *testing.T, line string) []string {
	t.Helper()
	// Parse preserving order via json.Decoder is not straightforward;
	// instead we re-encode and scan for key positions.
	var ordered []string
	dec := json.NewDecoder(nil)
	_ = dec
	// Simpler: unmarshal into interface{} won't preserve order.
	// We trust Apply's manual encoding and verify via round-trip equality.
	_ = decode(t, line)
	return ordered
}

func TestApply_NonJSON_ReturnsSameLine(t *testing.T) {
	s := New(false)
	input := "not json at all"
	if got := s.Apply(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestApply_AscendingOrder(t *testing.T) {
	s := New(false)
	input := `{"zebra":1,"apple":2,"mango":3}`
	out := s.Apply(input)

	// Verify the output is valid JSON with the same data.
	origMap := decode(t, input)
	outMap := decode(t, out)
	if len(origMap) != len(outMap) {
		t.Fatalf("key count mismatch: %d vs %d", len(origMap), len(outMap))
	}

	// Verify ascending order by scanning raw bytes.
	expected := `{"apple":2,"mango":3,"zebra":1}`
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestApply_DescendingOrder(t *testing.T) {
	s := New(true)
	input := `{"apple":1,"mango":2,"zebra":3}`
	out := s.Apply(input)
	expected := `{"zebra":3,"mango":2,"apple":1}`
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestApply_SingleKey_Unchanged(t *testing.T) {
	s := New(false)
	input := `{"only":"one"}`
	if got := s.Apply(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestApply_EmptyObject(t *testing.T) {
	s := New(false)
	input := `{}`
	out := s.Apply(input)
	if out != input {
		t.Errorf("expected %q, got %q", input, out)
	}
}

func TestApply_NestedValuesPreserved(t *testing.T) {
	s := New(false)
	input := `{"z":{"inner":true},"a":[1,2,3]}`
	out := s.Apply(input)
	expected := `{"a":[1,2,3],"z":{"inner":true}}`
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}
