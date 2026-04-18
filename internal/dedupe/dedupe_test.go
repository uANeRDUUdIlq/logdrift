package dedupe

import (
	"testing"
	"time"
)

func TestIsDuplicate_ZeroWindow_AlwaysFalse(t *testing.T) {
	d := New(0)
	if d.IsDuplicate("svc", "line") {
		t.Fatal("expected false for zero window")
	}
	if d.IsDuplicate("svc", "line") {
		t.Fatal("expected false for zero window on repeat")
	}
}

func TestIsDuplicate_FirstSeen_ReturnsFalse(t *testing.T) {
	d := New(time.Minute)
	if d.IsDuplicate("svc", "hello") {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicate_SecondSeen_ReturnsTrue(t *testing.T) {
	d := New(time.Minute)
	d.IsDuplicate("svc", "hello")
	if !d.IsDuplicate("svc", "hello") {
		t.Fatal("second occurrence within window should be a duplicate")
	}
}

func TestIsDuplicate_ExpiredEntry_ReturnsFalse(t *testing.T) {
	d := New(50 * time.Millisecond)
	now := time.Now()
	d.nowFunc = func() time.Time { return now }

	d.IsDuplicate("svc", "hello")

	// advance clock past window
	d.nowFunc = func() time.Time { return now.Add(100 * time.Millisecond) }

	if d.IsDuplicate("svc", "hello") {
		t.Fatal("entry should have expired")
	}
}

func TestIsDuplicate_DifferentServices_Independent(t *testing.T) {
	d := New(time.Minute)
	d.IsDuplicate("svcA", "line")
	if d.IsDuplicate("svcB", "line") {
		t.Fatal("different services should be tracked independently")
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	d := New(time.Minute)
	d.IsDuplicate("svc", "line")
	d.Reset()
	if d.IsDuplicate("svc", "line") {
		t.Fatal("after reset, line should not be a duplicate")
	}
}
