// Package jsongroup provides a Grouper that buffers consecutive JSON log lines
// sharing the same value for a configurable field and emits a single merged
// summary line when the group ends.
//
// This is useful for collapsing bursts of identical-context log entries into
// one representative line, reducing noise in high-volume tails while preserving
// the most recent field values and recording a line count.
//
// Usage:
//
//	g := jsongroup.New("service", "_count")
//	for _, raw := range lines {
//		for _, summary := range g.Push(raw) {
//			fmt.Println(summary)
//		}
//	}
//	for _, summary := range g.Flush() {
//		fmt.Println(summary)
//	}
package jsongroup
