// Package transform provides field-level JSON log mutations for logdrift.
//
// A Transformer can rename existing fields and inject static key-value pairs
// into every log line that passes through the pipeline. Non-JSON lines are
// passed through unchanged.
//
// # Rules
//
// A Rule specifies two types of mutations:
//
//   - Rename: maps old field names to new field names. If the old field does
//     not exist in a log line, the rename is silently skipped.
//
//   - Add: injects static key-value pairs into every log line. Existing fields
//     are not overwritten.
//
// # Example usage
//
//	tr := transform.New(transform.Rule{
//		Rename: map[string]string{"msg": "message"},
//		Add:    map[string]string{"env": "staging"},
//	})
//	output := tr.Apply(line)
package transform
