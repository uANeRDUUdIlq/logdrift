// Package teewriter provides a TeeWriter that appends every log line it
// receives to one or more files on disk.
//
// It is intended to be used alongside the normal stdout formatter so that
// operators can simultaneously stream logs to the terminal and persist a
// filtered subset for later analysis.
//
// Usage:
//
//	tw, err := teewriter.New([]string{"/var/log/logdrift/api.log"})
//	if err != nil { ... }
//	defer tw.Close()
//	tw.Write(line)
package teewriter
