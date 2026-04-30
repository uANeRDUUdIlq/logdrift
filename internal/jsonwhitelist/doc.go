// Package jsonwhitelist provides a Whitelister that strips all JSON keys
// not present in a configured allow-list.
//
// Usage:
//
//	w := jsonwhitelist.New([]string{"level", "msg", "ts"})
//	filtered := w.Apply(rawLine)
//
// Non-JSON lines are passed through unchanged. If the allowed field list
// is empty, every line is passed through unchanged (no-op mode).
package jsonwhitelist
