// Package aggregate provides a thread-safe counter that groups structured log
// lines by the value of a chosen JSON field.
//
// Usage:
//
//	a := aggregate.New("level")
//	for line := range lines {
//		a.Record(line)
//	}
//	a.Print(os.Stdout)
package aggregate
