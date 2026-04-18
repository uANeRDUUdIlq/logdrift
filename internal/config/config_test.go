package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/logdrift/internal/config"
)

func writeYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "logdrift.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	p := writeYAML(t, `
poll_interval: 500ms
services:
  - name: api
    path: /var/log/api.log
    format: pretty
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PollInterval != 500*time.Millisecond {
		t.Errorf("poll_interval: got %v, want 500ms", cfg.PollInterval)
	}
	if len(cfg.Services) != 1 || cfg.Services[0].Name != "api" {
		t.Errorf("unexpected services: %+v", cfg.Services)
	}
}

func TestLoad_DefaultPollInterval(t *testing.T) {
	p := writeYAML(t, `services:\n  - name: svc\n    path: /tmp/x.log\n`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PollInterval != 250*time.Millisecond {
		t.Errorf("default poll_interval: got %v", cfg.PollInterval)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/logdrift.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_NoServices(t *testing.T) {
	p := writeYAML(t, `services: []\n`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for empty services")
	}
}

func TestLoad_InvalidFormat(t *testing.T) {
	p := writeYAML(t, `services:\n  - name: svc\n    path: /tmp/x.log\n    format: json\n`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for bad format")
	}
}
