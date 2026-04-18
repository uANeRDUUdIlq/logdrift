package multiplexer_test

import (
	"sort"
	"testing"

	"github.com/user/logdrift/internal/multiplexer"
)

func makeLineChans(data map[string][]string) map[string]<-chan string {
	result := make(map[string]<-chan string, len(data))
	for svc, lines := range data {
		ch := make(chan string, len(lines))
		for _, l := range lines {
			ch <- l
		}
		close(ch)
		result[svc] = ch
	}
	return result
}

func TestRun_MergesAllEntries(t *testing.T) {
	data := map[string][]string{
		"api":    {`{"msg":"a1"}`, `{"msg":"a2"}`},
		"worker": {`{"msg":"w1"}`},
	}
	m := multiplexer.New(nil)
	m.Run(makeLineChans(data))

	var got []multiplexer.Entry
	for e := range m.Out() {
		got = append(got, e)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestRun_ServiceNamesPreserved(t *testing.T) {
	data := map[string][]string{
		"svcA": {"line1"},
		"svcB": {"line2"},
	}
	m := multiplexer.New(nil)
	m.Run(makeLineChans(data))

	svcs := map[string]bool{}
	for e := range m.Out() {
		svcs[e.Service] = true
	}

	keys := []string{"svcA", "svcB"}
	sort.Strings(keys)
	for _, k := range keys {
		if !svcs[k] {
			t.Errorf("expected service %q in output", k)
		}
	}
}

func TestRun_EmptyInput_ClosesChannel(t *testing.T) {
	m := multiplexer.New(nil)
	m.Run(map[string]<-chan string{})
	count := 0
	for range m.Out() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 entries, got %d", count)
	}
}
