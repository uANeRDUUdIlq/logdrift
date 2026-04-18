package sampling

import (
	"testing"
)

func TestAllow_ZeroRate_DropsAll(t *testing.T) {
	s := New(0.0, nil)
	for i := 0; i < 100; i++ {
		if s.Allow("svc") {
			t.Fatal("expected all lines to be dropped")
		}
	}
}

func TestAllow_FullRate_KeepsAll(t *testing.T) {
	s := New(1.0, nil)
	for i := 0; i < 100; i++ {
		if !s.Allow("svc") {
			t.Fatal("expected all lines to be kept")
		}
	}
}

func TestAllow_PerServiceRate_OverridesDefault(t *testing.T) {
	s := New(1.0, map[string]float64{"noisy": 0.0})
	if s.Allow("noisy") {
		t.Fatal("expected noisy service to be dropped")
	}
	if !s.Allow("other") {
		t.Fatal("expected other service to be kept")
	}
}

func TestAllow_PartialRate_StatisticallyCorrect(t *testing.T) {
	s := New(0.5, nil)
	kept := 0
	const n = 10000
	for i := 0; i < n; i++ {
		if s.Allow("svc") {
			kept++
		}
	}
	ratio := float64(kept) / float64(n)
	if ratio < 0.45 || ratio > 0.55 {
		t.Fatalf("expected ~50%% kept, got %.2f%%", ratio*100)
	}
}

func TestSetRate_UpdatesAtRuntime(t *testing.T) {
	s := New(1.0, nil)
	s.SetRate("svc", 0.0)
	for i := 0; i < 50; i++ {
		if s.Allow("svc") {
			t.Fatal("expected dropped after SetRate(0)")
		}
	}
}

func TestClamp_NegativeBecomesZero(t *testing.T) {
	s := New(-5.0, nil)
	for i := 0; i < 50; i++ {
		if s.Allow("svc") {
			t.Fatal("negative rate should clamp to 0")
		}
	}
}

func TestClamp_OverOneBecomesOne(t *testing.T) {
	s := New(99.0, nil)
	for i := 0; i < 50; i++ {
		if !s.Allow("svc") {
			t.Fatal("rate >1 should clamp to 1")
		}
	}
}
