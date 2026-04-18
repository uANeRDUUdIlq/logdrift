package timewindow_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/timewindow"
)

var (
	base  = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	early = base.Add(-time.Hour)
	late  = base.Add(time.Hour)
)

func jsonLine(ts time.Time) string {
	return fmt.Sprintf(`{"time":%q,"msg":"hello"}`, ts.Format(time.RFC3339))
}

func TestAllow_NonJSON_PassesThrough(t *testing.T) {
	f := timewindow.New("time", early, late)
	if !f.Allow("not json at all") {
		t.Fatal("expected non-JSON to pass")
	}
}

func TestAllow_NoTimestampField_PassesThrough(t *testing.T) {
	f := timewindow.New("time", early, late)
	if !f.Allow(`{"msg":"no ts here"}`) {
		t.Fatal("expected line without ts field to pass")
	}
}

func TestAllow_WithinWindow_Passes(t *testing.T) {
	f := timewindow.New("time", early, late)
	if !f.Allow(jsonLine(base)) {
		t.Fatal("expected line within window to pass")
	}
}

func TestAllow_BeforeSince_Dropped(t *testing.T) {
	f := timewindow.New("time", base, late)
	if f.Allow(jsonLine(early)) {
		t.Fatal("expected line before since to be dropped")
	}
}

func TestAllow_AfterUntil_Dropped(t *testing.T) {
	f := timewindow.New("time", early, base)
	if f.Allow(jsonLine(late)) {
		t.Fatal("expected line after until to be dropped")
	}
}

func TestAllow_ZeroSince_UnboundedLow(t *testing.T) {
	f := timewindow.New("time", time.Time{}, late)
	if !f.Allow(jsonLine(early)) {
		t.Fatal("expected zero since to allow any early line")
	}
}

func TestAllow_ZeroUntil_UnboundedHigh(t *testing.T) {
	f := timewindow.New("time", early, time.Time{})
	if !f.Allow(jsonLine(late)) {
		t.Fatal("expected zero until to allow any late line")
	}
}

func TestAllow_UnixTimestamp_Passes(t *testing.T) {
	f := timewindow.New("ts", early, late)
	line := fmt.Sprintf(`{"ts":%d,"msg":"unix"}`, base.Unix())
	if !f.Allow(line) {
		t.Fatal("expected unix timestamp within window to pass")
	}
}

func TestNew_DefaultField(t *testing.T) {
	f := timewindow.New("", early, late)
	if !f.Allow(jsonLine(base)) {
		t.Fatal("expected default field 'time' to be used")
	}
}
