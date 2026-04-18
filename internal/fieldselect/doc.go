// Package fieldselect implements field projection for structured JSON log lines.
//
// A Selector is configured with a list of top-level field names to retain.
// Any field not in that list is stripped before the line is forwarded to the
// formatter or output sink.  Lines that are not valid JSON are passed through
// unchanged so that raw / plaintext logs are never silently dropped.
//
// Example usage:
//
//	s := fieldselect.New([]string{"time", "level", "msg"})
//	filtered := s.Apply(rawLine)
package fieldselect
