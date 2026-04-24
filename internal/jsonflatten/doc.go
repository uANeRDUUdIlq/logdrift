// Package jsonflatten provides a Flattener pipeline stage that converts
// deeply-nested JSON log objects into a single-level key/value map.
//
// Nested keys are joined with a configurable separator (default ".").
// An optional prefix can be prepended to every resulting key, which is
// useful when merging multiple log sources that share field names.
//
// Non-JSON lines pass through unchanged so the stage is safe to include in
// pipelines that may contain mixed log formats.
//
// Usage:
//
//	f := jsonflatten.New(".", "")
//	output := f.Apply(line)
package jsonflatten
