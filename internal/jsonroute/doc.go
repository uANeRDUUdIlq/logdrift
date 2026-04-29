// Package jsonroute provides field-based routing for structured JSON
// log lines.
//
// Lines are inspected for a configurable JSON field and routed to the
// writer whose rule value matches. Non-JSON lines and unmatched lines
// are forwarded to an optional fallback writer.
//
// Example:
//
//	router, err := jsonroute.New([]jsonroute.Rule{
//		{Field: "level", Value: "error", Writer: errWriter},
//		{Field: "level", Value: "warn",  Writer: warnWriter},
//	}, fallbackWriter)
package jsonroute
