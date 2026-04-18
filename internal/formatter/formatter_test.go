package formatter

import (
	"strings"
	"testing"
)

func TestNew_AssignsService(t *testing.T) {
	// reset global index for deterministic colors
	colorIndex = 0
	f := New("svc-a", FormatPretty)
	if f.service != "svc-a" {
		t.Fatalf("expected service svc-a, got %s", f.service)
	}
}

func TestNew_CyclesColors(t *testing.T) {
	colorIndex = 0
	f1 := New("a", FormatPretty)
	f2 := New("b", FormatPretty)
	if f1.color == f2.color && len(colors) > 1 {
		t.Fatal("expected different colors for consecutive services")
	}
}

func TestRender_Raw(t *testing.T) {
	colorIndex = 0
	f := New("api", FormatRaw)
	line := `{"msg":"hello"}`
	out := f.Render(line)
	if !strings.Contains(out, "[api]") {
		t.Fatalf("expected service tag in output, got: %s", out)
	}
	if !strings.Contains(out, line) {
		t.Fatalf("expected raw line in output, got: %s", out)
	}
}

func TestRender_Pretty_ValidJSON(t *testing.T) {
	colorIndex = 0
	f := New("worker", FormatPretty)
	line := `{"level":"error","msg":"something broke","time":"2024-01-01T00:00:00Z","code":500}`
	out := f.Render(line)
	if !strings.Contains(out, "[worker]") {
		t.Fatalf("missing service tag: %s", out)
	}
	if !strings.Contains(out, "ERROR") {
		t.Fatalf("expected level ERROR in output: %s", out)
	}
	if !strings.Contains(out, "something broke") {
		t.Fatalf("expected message in output: %s", out)
	}
	if !strings.Contains(out, "code=500") {
		t.Fatalf("expected extra field code=500: %s", out)
	}
}

func TestRender_Pretty_InvalidJSON(t *testing.T) {
	colorIndex = 0
	f := New("db", FormatPretty)
	line := "not json at all"
	out := f.Render(line)
	if !strings.Contains(out, line) {
		t.Fatalf("expected raw line passthrough, got: %s", out)
	}
}

func TestRender_Pretty_UnixTimestamp(t *testing.T) {
	colorIndex = 0
	f := New("svc", FormatPretty)
	line := `{"level":"info","msg":"tick","ts":1700000000}`
	out := f.Render(line)
	if !strings.Contains(out, "2023") {
		t.Fatalf("expected formatted timestamp year in output: %s", out)
	}
}
