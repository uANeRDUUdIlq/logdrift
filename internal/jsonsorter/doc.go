// Package jsonsorter sorts the top-level keys of a JSON log line
// alphabetically (ascending or descending) before passing the line
// downstream.
//
// Non-JSON lines are forwarded unchanged.
//
// Usage:
//
//	sorter := jsonsorter.New(false) // ascending
//	output := sorter.Apply(line)
//
// This is useful when diffing log streams from multiple services where
// field insertion order may differ between implementations.
package jsonsorter
