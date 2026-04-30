// Package jsonclip provides a Clipper that limits the number of top-level keys
// in a JSON log line.
//
// This is useful when downstream systems impose field-count quotas or when you
// want to strip noisy fields that appear at the tail of a log object without
// enumerating them explicitly (complement to jsonstrip).
//
// Keys are retained in their original document order; excess keys are silently
// discarded.  Non-JSON lines are always passed through unchanged.
package jsonclip
