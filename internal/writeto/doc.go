// Package writeto provides a fan-out output sink for logdrift.
//
// A Sink accepts rendered log lines and writes them to one or more
// io.Writer targets concurrently-safe. Typical targets are os.Stdout
// and an open log file, allowing logdrift to mirror output to the
// terminal and a persistent file simultaneously.
//
// Usage:
//
//	sink := writeto.New(os.Stdout, logFile)
//	sink.Write(renderedLine)
package writeto
