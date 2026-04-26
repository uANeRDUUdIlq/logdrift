// Package jsontemplate renders log lines using Go text/template expressions
// against parsed JSON fields. Lines that are not valid JSON are passed through
// unchanged. If the template fails to execute the original line is returned.
package jsontemplate

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

// Renderer applies a Go template to each JSON log line, producing a new
// string. The template receives the decoded JSON object as its dot value
// (map[string]any), so fields are accessed with {{.fieldName}}.
//
// Example template: "{{.level}} {{.service}} — {{.msg}}"
type Renderer struct {
	tmpl *template.Template
}

// New compiles tplStr as a Go text/template and returns a Renderer.
// An error is returned when the template cannot be parsed.
func New(tplStr string) (*Renderer, error) {
	funcMap := template.FuncMap{
		// upper / lower are convenience helpers available inside templates.
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"trim":  strings.TrimSpace,
		// default returns the second argument when the first is the zero value.
		"default": func(def, val any) any {
			if val == nil || val == "" {
				return def
			}
			return val
		},
	}

	tmpl, err := template.New("logline").Funcs(funcMap).Parse(tplStr)
	if err != nil {
		return nil, err
	}
	return &Renderer{tmpl: tmpl}, nil
}

// Apply renders line through the compiled template.
//
// If line is not valid JSON, or if template execution fails, the original
// line is returned unmodified so that the pipeline is never interrupted.
func (r *Renderer) Apply(line string) string {
	var fields map[string]any
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		// Non-JSON lines pass through.
		return line
	}

	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, fields); err != nil {
		// Template execution error — fall back to original.
		return line
	}
	return buf.String()
}
