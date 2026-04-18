package stats_test

import (
	"bytes"
	"strings"
	"testing"

	"logdrift/internal/stats"
)

func TestRecordTotal_IncreasesCount(t *testing.T) {
	tr := stats.New()
	tr.RecordTotal("svc-a")
	tr.RecordTotal("svc-a")

	var buf bytes.Buffer
	tr.Print(&buf)
	if !strings.Contains(buf.String(), "2") {
		t.Fatalf("expected total 2 in output, got: %s", buf.String())
	}
}

func TestRecordMatch_IncreasesMatchedCount(t *testing.T) {
	tr := stats.New()
	tr.RecordTotal("svc-b")
	tr.RecordTotal("svc-b")
	tr.RecordMatch("svc-b")

	var buf bytes.Buffer
	tr.Print(&buf)
	output := buf.String()
	if !strings.Contains(output, "svc-b") {
		t.Fatal("expected svc-b in output")
	}
	// line should show total=2 matched=1
	if !strings.Contains(output, "2") || !strings.Contains(output, "1") {
		t.Fatalf("unexpected counts in output: %s", output)
	}
}

func TestPrint_MultipleServices_SortedAlphabetically(t *testing.T) {
	tr := stats.New()
	for _, svc := range []string{"zebra", "alpha", "mango"} {
		tr.RecordTotal(svc)
	}

	var buf bytes.Buffer
	tr.Print(&buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// lines[0] is header
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines (header + 3 services), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[1], "alpha") {
		t.Errorf("expected first service to be alpha, got: %s", lines[1])
	}
	if !strings.HasPrefix(lines[3], "zebra") {
		t.Errorf("expected last service to be zebra, got: %s", lines[3])
	}
}

func TestNew_EmptyTracker_PrintsHeaderOnly(t *testing.T) {
	tr := stats.New()
	var buf bytes.Buffer
	tr.Print(&buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected only header line, got %d lines", len(lines))
	}
}
