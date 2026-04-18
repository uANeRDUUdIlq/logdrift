package truncate_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/truncate"
)

func TestApply_Disabled_ZeroMaxLen(t *testing.T) {
	tr := truncate.New(0, "...")
	line := "this is a long line that should not be truncated"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_Disabled_NegativeMaxLen(t *testing.T) {
	tr := truncate.New(-1, "...")
	line := "another line"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_ShortLine_NotTruncated(t *testing.T) {
	tr := truncate.New(80, "...")
	line := "short"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_ExactLength_NotTruncated(t *testing.T) {
	tr := truncate.New(5, "...")
	line := "hello"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_LongLine_Truncated(t *testing.T) {
	tr := truncate.New(10, "...")
	line := "hello world this is long"
	got := tr.Apply(line)
	if len([]rune(got)) != 10 {
		t.Errorf("expected 10 runes, got %d: %q", len([]rune(got)), got)
	}
	if got[len(got)-3:] != "..." {
		t.Errorf("expected suffix '...', got %q", got)
	}
}

func TestApply_DefaultSuffix(t *testing.T) {
	tr := truncate.New(6, "")
	line := "abcdefgh"
	got := tr.Apply(line)
	if got != "abc..." {
		t.Errorf("expected 'abc...', got %q", got)
	}
}

func TestApply_UnicodeRunes(t *testing.T) {
	tr := truncate.New(4, "…")
	// Each char is one rune
	line := "日本語テスト"
	got := tr.Apply(line)
	runes := []rune(got)
	if len(runes) != 4 {
		t.Errorf("expected 4 runes, got %d: %q", len(runes), got)
	}
}
