// Package sampling implements probabilistic log sampling for logdrift.
//
// A Sampler holds a default rate and optional per-service rates in the range
// [0.0, 1.0]. Calling Allow(service) returns true with probability equal to
// the configured rate, allowing callers to discard a fraction of log lines
// before they are formatted or written to output.
//
// Rates outside [0.0, 1.0] are silently clamped. A rate of 0.0 drops every
// line; a rate of 1.0 keeps every line. Rates can be updated at runtime via
// SetRate without restarting the process.
package sampling
