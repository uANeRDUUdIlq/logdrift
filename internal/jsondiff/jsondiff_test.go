package jsondiff

import (
	"testing"
)

func TestChanged_NoFields_AlwaysTrue(t *testing.T) {
	d := New(nil)
	for i := 0; i < 3; i++ {
		if !d.Changed("svc", `{"level":"info"}`) {
			t.Fatal("expected true with no watched fields")
		}
	}
}

func TestChanged_NonJSON_AlwaysTrue(t *testing.T) {
	d := New([]string{"level"})
	if !d.Changed("svc", "not json at all") {
		t.Fatal("expected true for non-JSON line")
	}
}

func TestChanged_FirstSeen_ReturnsTrue(t *testing.T) {
	d := New([]string{"level"})
	if !d.Changed("svc", `{"level":"info"}`) {
		t.Fatal("first occurrence must return true")
	}
}

func TestChanged_SameValue_ReturnsFalse(t *testing.T) {
	d := New([]string{"level"})
	d.Changed("svc", `{"level":"info"}`)
	if d.Changed("svc", `{"level":"info"}`) {
		t.Fatal("unchanged field must return false")
	}
}

func TestChanged_DifferentValue_ReturnsTrue(t *testing.T) {
	d := New([]string{"level"})
	d.Changed("svc", `{"level":"info"}`)
	if !d.Changed("svc", `{"level":"error"}`) {
		t.Fatal("changed field must return true")
	}
}

func TestChanged_DifferentServices_Independent(t *testing.T) {
	d := New([]string{"level"})
	d.Changed("svcA", `{"level":"info"}`)
	d.Changed("svcB", `{"level":"info"}`)
	// svcA same value -> false
	if d.Changed("svcA", `{"level":"info"}`) {
		t.Fatal("svcA should not detect change")
	}
	// svcB different value -> true
	if !d.Changed("svcB", `{"level":"warn"}`) {
		t.Fatal("svcB should detect change")
	}
}

func TestChanged_MissingField_TreatedAsEmpty(t *testing.T) {
	d := New([]string{"level"})
	d.Changed("svc", `{"msg":"hello"}`)
	// still no level field -> same empty value -> false
	if d.Changed("svc", `{"msg":"world"}`) {
		t.Fatal("absent watched field should not trigger change")
	}
}

func TestReset_ClearsSnapshots(t *testing.T) {
	d := New([]string{"level"})
	d.Changed("svc", `{"level":"info"}`)
	d.Reset()
	// after reset the next identical line is treated as first-seen
	if !d.Changed("svc", `{"level":"info"}`) {
		t.Fatal("expected true after reset")
	}
}
