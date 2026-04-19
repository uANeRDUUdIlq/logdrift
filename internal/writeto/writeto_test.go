package writeto_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logdrift/internal/writeto"
)

func TestWrite_SingleWriter(t *testing.T) {
	var buf bytes.Buffer
	s := writeto.New(&buf)
	if err := s.Write("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "hello" {
		t.Fatalf("want 'hello', got %q", got)
	}
}

func TestWrite_MultipleWriters(t *testing.T) {
	var a, b bytes.Buffer
	s := writeto.New(&a, &b)
	_ = s.Write("ping")
	if strings.TrimSpace(a.String()) != "ping" {
		t.Errorf("writer a: got %q", a.String())
	}
	if strings.TrimSpace(b.String()) != "ping" {
		t.Errorf("writer b: got %q", b.String())
	}
}

func TestWrite_AppendedWriter(t *testing.T) {
	var a, b bytes.Buffer
	s := writeto.New(&a)
	s.Add(&b)
	_ = s.Write("added")
	if strings.TrimSpace(b.String()) != "added" {
		t.Errorf("added writer: got %q", b.String())
	}
}

func TestNew_PanicsWithNoWriters(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with no writers")
		}
	}()
	writeto.New()
}

func TestWrite_FailingWriter_ReturnsError(t *testing.T) {
	s := writeto.New(&failWriter{})
	if err := s.Write("boom"); err == nil {
		t.Error("expected error from failing writer")
	}
}

type failWriter struct{}

func (f *failWriter) Write(_ []byte) (int, error) {
	return 0, bytes.ErrTooLarge
}
