// Package throttle implements per-service line-rate throttling for logdrift.
//
// A Throttler enforces a maximum number of log lines per second for each
// service. Lines that exceed the configured rate are dropped silently.
//
// A rate of 0 disables throttling entirely, allowing all lines through.
//
// Example usage:
//
//	th := throttle.New(100) // allow 100 lines/sec per service
//	if th.Allow(entry.Service) {
//		// emit line
//	}
package throttle
