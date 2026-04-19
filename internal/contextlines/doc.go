// Package contextlines provides a Buffer that captures N lines before and
// after a matching log line, mirroring the behaviour of grep -B / -A.
//
// Usage:
//
//	buf := contextlines.New(2, 3) // 2 before, 3 after
//	for _, entry := range entries {
//		matched := filter.Match(entry.Line)
//		for _, l := range buf.Push(entry.Line, entry.Service, matched) {
//			fmt.Println(l)
//		}
//	}
package contextlines
