package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_MissingFile_EmptyOffsets(t *testing.T) {
	s, err := New(tmpPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get("/var/log/app.log"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSetAndGet_ReturnsOffset(t *testing.T) {
	s, _ := New(tmpPath(t))
	s.Set("/var/log/app.log", 1024)
	if got := s.Get("/var/log/app.log"); got != 1024 {
		t.Fatalf("expected 1024, got %d", got)
	}
}

func TestFlush_PersistsOffsets(t *testing.T) {
	p := tmpPath(t)
	s, _ := New(p)
	s.Set("/srv/logs/svc.log", 4096)
	if err := s.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}
	s2, err := New(p)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if got := s2.Get("/srv/logs/svc.log"); got != 4096 {
		t.Fatalf("expected 4096, got %d", got)
	}
}

func TestFlush_AtomicRename_NoDirtyFile(t *testing.T) {
	p := tmpPath(t)
	s, _ := New(p)
	s.Set("a", 1)
	_ = s.Flush()
	if _, err := os.Stat(p + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("tmp file should not exist after flush")
	}
}

func TestNew_InvalidJSON_ReturnsError(t *testing.T) {
	p := tmpPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0o644)
	if _, err := New(p); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestGet_UnknownFile_ReturnsZero(t *testing.T) {
	s, _ := New(tmpPath(t))
	if got := s.Get("unknown"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
