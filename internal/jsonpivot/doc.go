// Package jsonpivot provides a transformer that pivots an array of
// key/value objects inside a JSON log line into a flat top-level map.
//
// This is useful when services emit structured metrics or labels as an array
// of {"name":"...","value":...} entries and you want them as first-class
// fields for downstream filtering or rendering.
//
// Example
//
//	input:  {"svc":"api","tags":[{"k":"env","v":"prod"},{"k":"region","v":"us-east"}]}
//	output: {"svc":"api","env":"prod","region":"us-east"}
package jsonpivot
