package throttle_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/throttle"
)

func TestAllow_ZeroRate_AlwaysAllows(t *testing.T) {
	th := throttle.New(0)
	for i := 0; i < 100; i++ {
		if !th.Allow("svc") {
			t.Fatal("expected allow with zero rate")
		}
	}
}

func TestAllow_WithinRate_Allowed(t *testing.T) {
	th := throttle.New(5)
	for i := 0; i < 5; i++ {
		if !th.Allow("svc") {
			t.Fatalf("expected allow on call %d", i)
		}
	}
}

func TestAllow_ExceedsRate_Dropped(t *testing.T) {
	th := throttle.New(3)
	allowed := 0
	for i := 0; i < 10; i++ {
		if th.Allow("svc") {
			allowed++
		}
	}
	if allowed != 3 {
		t.Fatalf("expected 3 allowed, got %d", allowed)
	}
}

func TestAllow_DifferentServices_Independent(t *testing.T) {
	th := throttle.New(2)
	if !th.Allow("a") { t.Fatal("a call 1") }
	if !th.Allow("a") { t.Fatal("a call 2") }
	if th.Allow("a") { t.Fatal("a call 3 should be dropped") }
	if !th.Allow("b") { t.Fatal("b call 1") }
	if !th.Allow("b") { t.Fatal("b call 2") }
	if th.Allow("b") { t.Fatal("b call 3 should be dropped") }
}

func TestReset_ClearsBuckets(t *testing.T) {
	th := throttle.New(1)
	th.Allow("svc") // consume token
	if th.Allow("svc") {
		t.Fatal("expected drop before reset")
	}
	th.Reset()
	if !th.Allow("svc") {
		t.Fatal("expected allow after reset")
	}
}
