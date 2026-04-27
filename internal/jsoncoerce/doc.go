// Package jsoncoerce provides a Coercer pipeline stage that converts JSON
// field values to a specified target type.
//
// Supported target types:
//
//	"string" — converts numbers and booleans to their string representation.
//	"number" — parses string values as float64; converts booleans (true→1, false→0).
//	"bool"   — parses string values via strconv.ParseBool; converts numbers (0→false, else→true).
//
// Non-JSON lines and fields that cannot be coerced are passed through
// unchanged. This stage is useful for normalising inconsistent field types
// that appear across different log producers before downstream processing.
package jsoncoerce
