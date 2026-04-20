// Package jsonmerge provides a Merger that enriches structured JSON log lines
// with a static set of key/value string fields at processing time.
//
// A common use-case is stamping every log line with deployment metadata such as
// environment name, region, or cluster identifier that is not already present in
// the raw log output.
//
// Existing fields in the log line are never overwritten; the original value
// always takes precedence over the configured merge fields.
//
// Non-JSON lines pass through unmodified.
package jsonmerge
