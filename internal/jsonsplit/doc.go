// Package jsonsplit provides a Splitter that expands a JSON array field
// into multiple individual log lines — one per array element.
//
// This is useful when a single structured log entry contains a batch of
// events (e.g. a list of errors, spans, or records) and downstream
// processors expect one event per line.
//
// Example input:
//
//	{"service":"api","errors":[{"code":404},{"code":500}]}
//
// Example output (two lines):
//
//	{"service":"api","errors":{"code":404}}
//	{"service":"api","errors":{"code":500}}
//
// Non-JSON lines, lines missing the target field, and lines where the
// field value is not an array are passed through unchanged.
package jsonsplit
