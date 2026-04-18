package levelfilter

import (
	"testing"
)

func TestAllow_NonJSON_AlwaysPasses(t *testing.T) {
	f := New("error", "level")
	if !f.Allow("plain text log line") {
		t.Fatal("expected non-JSON line to pass")
	}
}

func TestAllow_NoLevelField_Passes(t *testing.T) {
	f := New("error", "level")
	if !f.Allow(`{"msg":"hello"}`) {
		t.Fatal("expected line without level field to pass")
	}
}

func TestAllow_BelowMin_Dropped(t *testing.T) {
	f := New("warn", "level")
	if f.Allow(`{"level":"debug","msg":"verbose"}`) {
		t.Fatal("expected debug line to be dropped when min=warn")
	}
	if f.Allow(`{"level":"info","msg":"info msg"}`) {
		t.Fatal("expected info line to be dropped when min=warn")
	}
}

func TestAllow_AtMin_Passes(t *testing.T) {
	f := New("warn", "level")
	if !f.Allow(`{"level":"warn","msg":"watch out"}`) {
		t.Fatal("expected warn line to pass when min=warn")
	}
}

func TestAllow_AboveMin_Passes(t *testing.T) {
	f := New("info", "level")
	if !f.Allow(`{"level":"error","msg":"boom"}`) {
		t.Fatal("expected error line to pass when min=info")
	}
}

func TestAllow_CaseInsensitive(t *testing.T) {
	f := New("INFO", "level")
	if !f.Allow(`{"level":"WARN","msg":"ok"}`) {
		t.Fatal("expected case-insensitive match to pass")
	}
	if f.Allow(`{"level":"DEBUG","msg":"low"}`) {
		t.Fatal("expected DEBUG to be dropped when min=info")
	}
}

func TestAllow_CustomLevelKey(t *testing.T) {
	f := New("error", "severity")
	if f.Allow(`{"severity":"info","msg":"meh"}`) {
		t.Fatal("expected info to be dropped with custom key and min=error")
	}
	if !f.Allow(`{"severity":"fatal","msg":"bad"}`) {
		t.Fatal("expected fatal to pass with min=error")
	}
}

func TestNew_EmptyMinLevel_AllowsAll(t *testing.T) {
	f := New("", "level")
	if !f.Allow(`{"level":"debug","msg":"trace"}`) {
		t.Fatal("expected all levels to pass when minLevel is empty")
	}
}
