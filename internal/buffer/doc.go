// Package buffer implements a fixed-size ring buffer for recent log lines.
//
// Buffer retains the last N lines pushed to it, evicting the oldest entry
// when capacity is exceeded. It is safe for concurrent use.
//
// Typical usage:
//
//	b := buffer.New(100)
//	b.Push(line)
//	for _, l := range b.Lines() {
//		fmt.Println(l)
//	}
package buffer
