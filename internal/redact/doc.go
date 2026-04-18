// Package redact provides field-level redaction for structured JSON log lines.
//
// A Redactor is initialised with a list of field names and an optional mask
// string. When Apply is called with a raw log line the package attempts to
// parse the line as JSON; any top-level key whose lower-cased name matches a
// configured field is replaced with the mask value before the object is
// re-serialised and returned.
//
// Non-JSON lines pass through unmodified, making the redactor safe to use in
// pipelines that may contain mixed log formats.
package redact
