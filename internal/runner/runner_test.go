package runner_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/logdrift/internal/runner"
)

func writeConfig(t *testing.T, logFile string) string {
	t.Helper()
	content := `services:
  - name: svc1
    path: ` + logFile + `
poll_interval: 20ms
`
	p := filepath.Join(t.TempDir(), "logdrift.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_OutputsMatchingLines(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "app.log")
	if err := os.WriteFile(tmp, []byte(`{"level":"info","msg":"hello"}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfgPath := writeConfig(t, tmp)
	var buf bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	err := runner.Run(ctx, runner.Options{
		ConfigPath: cfgPath,
		RawOutput:  true,
		Output:     &buf,
	})
	// context cancellation is expected, not an error we surface
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("hello")) {
		t.Errorf("expected 'hello' in output, got: %s", buf.String())
	}
}

func TestRun_FilterDropsLines(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "app.log")
	lines := `{"level":"debug","msg":"verbose"}` + "\n" +
		`{"level":"error","msg":"boom"}` + "\n"
	if err := os.WriteFile(tmp, []byte(lines), 0o644); err != nil {
		t.Fatal(err)
	}

	cfgPath := writeConfig(t, tmp)
	var buf bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	_ = runner.Run(ctx, runner.Options{
		ConfigPath: cfgPath,
		RawOutput:  true,
		Filters:    []string{`level == "error"`},
		Output:     &buf,
	})

	if bytes.Contains(buf.Bytes(), []byte("verbose")) {
		t.Error("debug line should have been filtered out")
	}
	if !bytes.Contains(buf.Bytes(), []byte("boom")) {
		t.Error("error line should be present")
	}
}

func TestRun_BadConfig(t *testing.T) {
	err := runner.Run(context.Background(), runner.Options{
		ConfigPath: "/nonexistent/path.yaml",
	})
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}
