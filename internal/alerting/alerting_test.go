package alerting

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
	"time"
)

func rule(name, pattern string, cooldown time.Duration) Rule {
	return Rule{Name: name, Pattern: regexp.MustCompile(pattern), Cooldown: cooldown}
}

func TestCheck_NoRules_NoOutput(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf, nil)
	a.Check("svc", `{"level":"error"}`)
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestCheck_MatchingRule_WritesAlert(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf, []Rule{rule("err-alert", "error", 0)})
	a.Check("api", `{"level":"error","msg":"boom"}`)
	if !strings.Contains(buf.String(), "err-alert") {
		t.Fatalf("expected alert name in output, got %q", buf.String())
	}
	if !strings.Contains(buf.String(), "api") {
		t.Fatalf("expected service name in output, got %q", buf.String())
	}
}

func TestCheck_NonMatchingRule_NoOutput(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf, []Rule{rule("err-alert", "panic", 0)})
	a.Check("api", `{"level":"info"}`)
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestCheck_Cooldown_SuppressesDuplicate(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf, []Rule{rule("err-alert", "error", 10*time.Second)})
	a.Check("api", "error occurred")
	a.Check("api", "error occurred")
	count := strings.Count(buf.String(), "[ALERT]")
	if count != 1 {
		t.Fatalf("expected 1 alert due to cooldown, got %d", count)
	}
}

func TestCheck_CooldownExpired_FiresAgain(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf, []Rule{rule("err-alert", "error", time.Millisecond)})
	a.Check("api", "error")
	time.Sleep(5 * time.Millisecond)
	a.Check("api", "error")
	count := strings.Count(buf.String(), "[ALERT]")
	if count != 2 {
		t.Fatalf("expected 2 alerts, got %d", count)
	}
}

func TestCheck_DifferentServices_Independent(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf, []Rule{rule("err-alert", "error", 10*time.Second)})
	a.Check("api", "error")
	a.Check("worker", "error")
	count := strings.Count(buf.String(), "[ALERT]")
	if count != 2 {
		t.Fatalf("expected 2 alerts for different services, got %d", count)
	}
}

func TestReset_ClearsCooldowns(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf, []Rule{rule("err-alert", "error", 10*time.Second)})
	a.Check("api", "error")
	a.Reset()
	a.Check("api", "error")
	count := strings.Count(buf.String(), "[ALERT]")
	if count != 2 {
		t.Fatalf("expected 2 alerts after reset, got %d", count)
	}
}
