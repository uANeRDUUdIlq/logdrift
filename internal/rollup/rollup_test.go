package rollup

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func jsonLine(fields map[string]interface{}) string {
	b, _ := json.Marshal(fields)
	return string(b)
}

func TestRecord_NonJSON_Ignored(t *testing.T) {
	r := New("level", 0)
	r.Record("svc", "not json")
	r.Flush()
	if len(r.Entries()) != 0 {
		t.Fatal("expected no entries")
	}
}

func TestRecord_MissingField_Ignored(t *testing.T) {
	r := New("level", 0)
	r.Record("svc", jsonLine(map[string]interface{}{"msg": "hi"}))
	r.Flush()
	if len(r.Entries()) != 0 {
		t.Fatal("expected no entries")
	}
}

func TestFlush_CountsAggregated(t *testing.T) {
	r := New("level", 0)
	for i := 0; i < 5; i++ {
		r.Record("svc", jsonLine(map[string]interface{}{"level": "error"}))
	}
	r.Flush()
	entry := <-r.Entries()
	if entry.Count != 5 {
		t.Fatalf("want 5, got %d", entry.Count)
	}
	if entry.Value != "error" {
		t.Fatalf("want error, got %s", entry.Value)
	}
}

func TestFlush_MultipleValues_SeparateEntries(t *testing.T) {
	r := New("level", 0)
	r.Record("svc", jsonLine(map[string]interface{}{"level": "info"}))
	r.Record("svc", jsonLine(map[string]interface{}{"level": "error"}))
	r.Record("svc", jsonLine(map[string]interface{}{"level": "info"}))
	r.Flush()
	counts := map[string]int{}
	for i := 0; i < 2; i++ {
		e := <-r.Entries()
		counts[e.Value] = e.Count
	}
	if counts["info"] != 2 || counts["error"] != 1 {
		t.Fatalf("unexpected counts: %v", counts)
	}
}

func TestFlush_ResetsAfterFlush(t *testing.T) {
	r := New("level", 0)
	r.Record("svc", jsonLine(map[string]interface{}{"level": "warn"}))
	r.Flush()
	<-r.Entries()
	r.Flush()
	if len(r.Entries()) != 0 {
		t.Fatal("expected empty after second flush")
	}
}

func TestFlush_DifferentServices_Independent(t *testing.T) {
	r := New("level", 0)
	r.Record("alpha", jsonLine(map[string]interface{}{"level": "info"}))
	r.Record("beta", jsonLine(map[string]interface{}{"level": "info"}))
	r.Flush()
	svcs := map[string]bool{}
	for i := 0; i < 2; i++ {
		e := <-r.Entries()
		svcs[e.Service] = true
	}
	if !svcs["alpha"] || !svcs["beta"] {
		t.Fatal("expected entries for both services")
	}
}

func TestTicker_AutoFlush(t *testing.T) {
	r := New("level", 50*time.Millisecond)
	defer r.Stop()
	r.Record("svc", jsonLine(map[string]interface{}{"level": "debug"}))
	select {
	case e := <-r.Entries():
		if e.Count != 1 {
			t.Fatalf("want 1, got %d", e.Count)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for auto-flush")
	}
	_ = fmt.Sprintf("") // suppress import
}
