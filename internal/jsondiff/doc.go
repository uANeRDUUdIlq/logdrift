// Package jsondiff provides a Differ that forwards a log line only when one
// or more nominated JSON fields change value relative to the previously seen
// line for the same service.
//
// Typical usage:
//
//	d := jsondiff.New([]string{"level", "status"})
//	if d.Changed(entry.Service, entry.Line) {
//		// emit the line
//	}
//
// Non-JSON lines are always considered changed and passed through.  An empty
// field list disables filtering entirely.
package jsondiff
