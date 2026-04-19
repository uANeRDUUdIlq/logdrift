package retrywriter_test

import (
	"errors"
	"io"
	"testing"
	"time"

	"logdrift/internal/retrywriter"
)

type countingWriter struct {
	failFor int
	calls   int
	buf     []byte
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	cw.calls++
	if cw.calls <= cw.failFor {
		return 0, errors.New("transient error")
	}
	cw.buf = append(cw.buf, p...)
	return len(p), nil
}

func TestWrite_SucceedsFirstAttempt(t *testing.T) {
	cw := &countingWriter{}
	rw := retrywriter.New(cw, 3, 0)
	_, err := rw.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cw.calls != 1 {
		t.Fatalf("expected 1 call, got %d", cw.calls)
	}
}

func TestWrite_RetriesOnTransientError(t *testing.T) {
	cw := &countingWriter{failFor: 2}
	rw := retrywriter.New(cw, 5, 0)
	_, err := rw.Write([]byte("line"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cw.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", cw.calls)
	}
}

func TestWrite_AllAttemptsFail_ReturnsError(t *testing.T) {
	cw := &countingWriter{failFor: 10}
	rw := retrywriter.New(cw, 3, 0)
	_, err := rw.Write([]byte("data"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if cw.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", cw.calls)
	}
}

func TestNew_ZeroAttempts_ClampsToOne(t *testing.T) {
	cw := &countingWriter{failFor: 1}
	rw := retrywriter.New(cw, 0, 0)
	_, err := rw.Write([]byte("x"))
	if err == nil {
		t.Fatal("expected error with only 1 attempt against failing writer")
	}
	if cw.calls != 1 {
		t.Fatalf("expected 1 call, got %d", cw.calls)
	}
}

func TestWrite_DelayBetweenRetries(t *testing.T) {
	cw := &countingWriter{failFor: 1}
	rw := retrywriter.New(cw, 3, 10*time.Millisecond)
	start := time.Now()
	_, err := rw.Write([]byte("timed"))
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed < 10*time.Millisecond {
		t.Fatalf("expected delay between retries, elapsed=%v", elapsed)
	}
}

func TestWrite_ImplementsIOWriter(t *testing.T) {
	var _ io.Writer = retrywriter.New(io.Discard, 1, 0)
}
