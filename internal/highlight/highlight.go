// Package highlight provides keyword highlighting for log output.
package highlight

import (
	"strings"

	"github.com/fatih/color"
)

// Highlighter holds compiled highlight rules.
type Highlighter struct {
	keywords []keyword
}

type keyword struct {
	word  string
	color *color.Color
}

// Rule maps a keyword to a named color.
type Rule struct {
	Word  string
	Color string
}

var namedColors = map[string]*color.Color{
	"red":     color.New(color.FgRed),
	"yellow":  color.New(color.FgYellow),
	"green":   color.New(color.FgGreen),
	"blue":    color.New(color.FgBlue),
	"magenta": color.New(color.FgMagenta),
	"cyan":    color.New(color.FgCyan),
}

// New creates a Highlighter from a list of Rules.
// Unknown color names fall back to yellow.
func New(rules []Rule) *Highlighter {
	h := &Highlighter{}
	for _, r := range rules {
		c, ok := namedColors[strings.ToLower(r.Color)]
		if !ok {
			c = namedColors["yellow"]
		}
		h.keywords = append(h.keywords, keyword{word: r.Word, color: c})
	}
	return h
}

// Apply replaces all keyword occurrences in line with their colored versions.
// If no rules are configured the original line is returned unchanged.
func (h *Highlighter) Apply(line string) string {
	if len(h.keywords) == 0 {
		return line
	}
	for _, kw := range h.keywords {
		colored := kw.color.Sprint(kw.word)
		line = strings.ReplaceAll(line, kw.word, colored)
	}
	return line
}
