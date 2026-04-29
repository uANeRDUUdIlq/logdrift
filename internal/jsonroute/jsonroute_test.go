package jsonroute_test

import (
	"bytes"
	"strings"
	"testing"

	"logdrift/internal/jsonroute"
)

func TestNew_EmptyField_ReturnsError(t *testing.T) {
	_, err := jsonroute.New([]jsonroute.Rule{
		{Field: "", Value: "error", Writer: &bytes.Buffer{}},
	}, nil)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_NilWriter_ReturnsError(t *testing.T) {
	_, err := jsonroute.New([]jsonroute.Rule{
		{Field: "level", Value: "error", Writer: nil},
	}, nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestRoute_MatchingRule_WritesToCorrectWriter(t *testing.T) {
	errBuf := &bytes.Buffer{}
	warnBuf := &bytes.Buffer{}
	router, err := jsonroute.New([]jsonroute.Rule{
		{Field: "level", Value: "error", Writer: errBuf},
		{Field: "level", Value: "warn", Writer: warnBuf},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := router.Route(`{"level":"error","msg":"boom"}`); err != nil {
		t.Fatalf("Route error: %v", err)
	}
	if !strings.Contains(errBuf.String(), "boom") {
		t.Errorf("expected error writer to contain 'boom', got %q", errBuf.String())
	}
	if warnBuf.Len() != 0 {
		t.Errorf("expected warn writer to be empty")
	}
}

func TestRoute_NoMatch_WritesFallback(t *testing.T) {
	fallback := &bytes.Buffer{}
	router, _ := jsonroute.New([]jsonroute.Rule{
		{Field: "level", Value: "error", Writer: &bytes.Buffer{}},
	}, fallback)

	_ = router.Route(`{"level":"info","msg":"hello"}`)
	if !strings.Contains(fallback.String(), "hello") {
		t.Errorf("expected fallback to contain 'hello', got %q", fallback.String())
	}
}

func TestRoute_NoMatch_NilFallback_DropsLine(t *testing.T) {
	dst := &bytes.Buffer{}
	router, _ := jsonroute.New([]jsonroute.Rule{
		{Field: "level", Value: "error", Writer: dst},
	}, nil)

	if err := router.Route(`{"level":"debug","msg":"quiet"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Len() != 0 {
		t.Errorf("expected destination to be empty")
	}
}

func TestRoute_NonJSON_WritesFallback(t *testing.T) {
	fallback := &bytes.Buffer{}
	router, _ := jsonroute.New([]jsonroute.Rule{
		{Field: "level", Value: "error", Writer: &bytes.Buffer{}},
	}, fallback)

	_ = router.Route("plain text line")
	if !strings.Contains(fallback.String(), "plain text line") {
		t.Errorf("expected fallback to contain plain text, got %q", fallback.String())
	}
}

func TestRoute_NumericFieldValue_Matches(t *testing.T) {
	dst := &bytes.Buffer{}
	router, _ := jsonroute.New([]jsonroute.Rule{
		{Field: "code", Value: "404", Writer: dst},
	}, nil)

	_ = router.Route(`{"code":404,"msg":"not found"}`)
	if !strings.Contains(dst.String(), "not found") {
		t.Errorf("expected dst to contain 'not found', got %q", dst.String())
	}
}

func TestBuild_NoRules_ReturnsError(t *testing.T) {
	_, err := jsonroute.Build(nil, nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestBuild_Valid_ReturnsRouter(t *testing.T) {
	dst := &bytes.Buffer{}
	router, err := jsonroute.Build([]jsonroute.BuilderRule{
		{Field: "level", Value: "error", Writer: dst},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if router == nil {
		t.Fatal("expected non-nil router")
	}
}
