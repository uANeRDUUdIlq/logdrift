package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Format controls the output style.
type Format int

const (
	FormatPretty Format = iota
	FormatRaw
)

// Formatter renders log lines for display.
type Formatter struct {
	format  Format
	service string
	color   string
}

var colors = []string{
	"\033[36m", // cyan
	"\033[32m", // green
	"\033[33m", // yellow
	"\033[35m", // magenta
	"\033[34m", // blue
}

var colorIndex int

// New returns a Formatter for the given service name and format.
func New(service string, format Format) *Formatter {
	c := colors[colorIndex%len(colors)]
	colorIndex++
	return &Formatter{format: format, service: service, color: c}
}

// Render formats a raw JSON log line for output.
func (f *Formatter) Render(line string) string {
	if f.format == FormatRaw {
		return fmt.Sprintf("%s[%s]\033[0m %s", f.color, f.service, line)
	}

	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return fmt.Sprintf("%s[%s]\033[0m %s", f.color, f.service, line)
	}

	ts := extractTime(fields)
	level := extractString(fields, "level", "info")
	msg := extractString(fields, "msg", extractString(fields, "message", line))

	var extras []string
	skip := map[string]bool{"level": true, "msg": true, "message": true, "time": true, "ts": true, "timestamp": true}
	for k, v := range fields {
		if !skip[k] {
			extras = append(extras, fmt.Sprintf("%s=%v", k, v))
		}
	}

	extra := ""
	if len(extras) > 0 {
		extra = " " + strings.Join(extras, " ")
	}

	return fmt.Sprintf("%s[%s]\033[0m %s %-5s %s%s", f.color, f.service, ts, strings.ToUpper(level), msg, extra)
}

func extractString(fields map[string]interface{}, key, def string) string {
	if v, ok := fields[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

func extractTime(fields map[string]interface{}) string {
	for _, key := range []string{"time", "ts", "timestamp"} {
		if v, ok := fields[key]; ok {
			switch val := v.(type) {
			case string:
				return val
			case float64:
				return time.Unix(int64(val), 0).UTC().Format(time.RFC3339)
			}
		}
	}
	return time.Now().UTC().Format(time.RFC3339)
}
