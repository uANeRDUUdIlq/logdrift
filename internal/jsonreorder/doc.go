// Package jsonreorder provides a transformer that reorders the keys of JSON
// log lines so that a configured set of high-priority keys (e.g. "ts",
// "level", "msg") always appear first.
//
// This is useful when piping logdrift output to human-readable terminals or
// downstream tools that expect a canonical field ordering.
//
// Non-JSON lines are passed through without modification. Keys not listed in
// the priority list are appended after the priority keys in their original
// iteration order.
package jsonreorder
