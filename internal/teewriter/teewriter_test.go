package teewriter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/logdrift/logdrift/internal/teewriter"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tee-*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestNew_NoPaths_NoError(t *testing.T) {
	tw, err := teewriter.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tw.Close()
	// Write should be a no-op.
	if err := tw.Write("hello"); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
}

func TestNew_BadPath_ReturnsError(t *testing.T) {
	_, err := teewriter.New([]string{"/no/such/dir/file.log"})
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}

func TestWrite_SingleFile_ContainsLine(t *testing.T) {
	p := tmpFile(t)
	tw, err := teewriter.New([]string{p})
	if err != nil {
		t.Fatal(err)
	}
	defer tw.Close()

	if err := tw.Write(`{"level":"info","msg":"hello"}`); err != nil {
		t.Fatal(err)
	}
	tw.Close()

	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "hello") {
		t.Errorf("expected file to contain 'hello', got: %s", data)
	}
}

func TestWrite_MultipleFiles_AllReceiveLine(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.log")
	p2 := filepath.Join(dir, "b.log")

	tw, err := teewriter.New([]string{p1, p2})
	if err != nil {
		t.Fatal(err)
	}
	const line = `{"service":"api","msg":"started"}`
	if err := tw.Write(line); err != nil {
		t.Fatal(err)
	}
	tw.Close()

	for _, p := range []string{p1, p2} {
		data, _ := os.ReadFile(p)
		if !strings.Contains(string(data), "started") {
			t.Errorf("%s: expected 'started', got: %s", p, data)
		}
	}
}

func TestClose_Idempotent(t *testing.T) {
	p := tmpFile(t)
	tw, _ := teewriter.New([]string{p})
	tw.Close()
	// Second close should not panic.
	tw.Close()
}
