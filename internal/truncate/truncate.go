// Package truncate provides line-length truncation for log output.
package truncate

import "unicode/utf8"

// Truncator truncates log lines that exceed a maximum rune length.
type Truncator struct {
	maxLen  int
	suffix  string
}

// New returns a Truncator. If maxLen is zero or negative, truncation is
// disabled. suffix is appended to truncated lines (e.g. "...").
func New(maxLen int, suffix string) *Truncator {
	if suffix == "" {
		suffix = "..."
	}
	return &Truncator{maxLen: maxLen, suffix: suffix}
}

// Apply truncates line if it exceeds the configured maximum rune count.
// If truncation is disabled (maxLen <= 0) the original line is returned.
func (t *Truncator) Apply(line string) string {
	if t.maxLen <= 0 {
		return line
	}
	count := utf8.RuneCountInString(line)
	if count <= t.maxLen {
		return line
	}
	// Trim to maxLen - len(suffix runes) runes, then append suffix.
	suffixLen := utf8.RuneCountInString(t.suffix)
	keep := t.maxLen - suffixLen
	if keep <= 0 {
		return t.suffix
	}
	runes := []rune(line)
	return string(runes[:keep]) + t.suffix
}
