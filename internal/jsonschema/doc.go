// Package jsonschema provides a lightweight field-type validator for structured
// JSON log lines.
//
// A Validator is constructed with a list of Rules, each mapping a top-level
// JSON field to an expected type (string, number, boolean, object, or array).
// Fields that are absent in a log line are not checked — only present fields
// are validated.
//
// When a violation is detected the Validator can either:
//   - Drop the line entirely (DropOnFail: true), or
//   - Pass the line through, optionally injecting a human-readable summary
//     into a configurable error field (e.g. "_schema_error").
//
// Non-JSON lines are always passed through unchanged.
package jsonschema
