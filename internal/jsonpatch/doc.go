// Package jsonpatch provides a Patcher that applies a fixed set of key/value
// overrides to structured JSON log lines at read time.
//
// This is useful for injecting metadata such as environment, region, or
// deployment version into log lines that were emitted without that context.
//
// Non-JSON lines are passed through unchanged so that raw or plaintext
// services are not broken by the pipeline.
//
// Example usage:
//
//	p, err := jsonpatch.New(map[string]string{"env": "prod", "region": "eu-west-1"})
//	patched := p.Apply(rawLine)
package jsonpatch
