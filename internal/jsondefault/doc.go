// Package jsondefault provides a Defaulter that injects missing fields into
// structured JSON log lines.
//
// Fields already present in the log line are never overwritten, making this
// safe to use as a fallback enrichment step in a processing pipeline.
//
// Example usage:
//
//	d := jsondefault.New(map[string]any{
//		"env":    "production",
//		"region": "us-east-1",
//	})
//	out := d.Apply(line)
package jsondefault
