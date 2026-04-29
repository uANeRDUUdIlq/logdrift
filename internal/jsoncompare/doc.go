// Package jsoncompare provides a filter that watches one or more JSON fields
// across log lines and emits a line only when at least one tracked field
// changes value relative to the previous line seen for the same service.
//
// When a change is detected the emitted line is annotated with:
//
//	_prev_<field>  – the value before the change (omitted on first occurrence)
//	_next_<field>  – the new value
//
// This is useful for surfacing state transitions (e.g. HTTP status codes
// flipping from 200 to 500) without drowning in repeated identical lines.
//
// Non-JSON lines and lines that do not contain any tracked field always pass
// through unchanged.
package jsoncompare
