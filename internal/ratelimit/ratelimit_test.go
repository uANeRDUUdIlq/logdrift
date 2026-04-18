package ratelimit_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/ratelimit"
)

func TestAllow_ZeroRate_AlwaysAllows(t *testing.T) {
	l := ratelimit.New(0)
	for i := 0; i < 1000; i++ {
		if !l.Allow("svc") {
			t.Fatal("expected all lines allowed when rate is 0")
		}
	}
}

func TestAllow_WithinRate_Allowed(t *testing.T) {
	l := ratelimit.New(5)
	for i := 0; i < 5; i++ {
		if !l.Allow("svc") {
			t.Fatalf("expected line %d to be allowed", i+1)
		}
	}
}

func TestAllow_ExceedsRate_Dropped(t *testing.T) {
	l := ratelimit.New(3)
	for i := 0; i < 3; i++ {
		l.Allow("svc")
	}
	if l.Allow("svc") {
		t.Fatal("expected 4th line to be dropped")
	}
}

func TestAllow_DifferentServices_IndependentBuckets(t *testing.T) {
	l := ratelimit.New(2)
	l.Allow("a")
	l.Allow("a")
	// service a is now at limit
	if l.Allow("a") {
		t.Fatal("expected service a to be rate limited")
	}
	// service b should still be allowed
	if !l.Allow("b") {
		t.Fatal("expected service b to be allowed independently")
	}
}

func TestReset_ClearsBuckets(t *testing.T) {
	l := ratelimit.New(1)
	l.Allow("svc")
	if l.Allow("svc") {
		t.Fatal("expected svc to be rate limited before reset")
	}
	l.Reset()
	if !l.Allow("svc") {
		t.Fatal("expected svc to be allowed after reset")
	}
}
