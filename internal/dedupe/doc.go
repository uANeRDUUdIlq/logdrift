// Package dedupe implements a sliding-window deduplication filter for log
// lines. It is designed to suppress bursts of identical messages that commonly
// appear when a service repeatedly logs the same error within a short period.
//
// Usage:
//
//	dd := dedupe.New(5 * time.Second)
//	if !dd.IsDuplicate(serviceName, rawLine) {
//	    // forward line to output
//	}
//
// Each (service, line) pair is tracked independently so that identical text
// from different services is never conflated.
package dedupe
