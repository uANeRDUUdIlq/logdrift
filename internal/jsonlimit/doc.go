// Package jsonlimit provides a transformer that truncates JSON array fields
// to a configured maximum number of elements.
//
// This is useful when log lines contain large arrays (e.g. stack frames,
// batch IDs, tag lists) that would otherwise overwhelm the output.
//
// A per-field limit takes precedence over the global default. Fields that are
// not arrays, or whose length is already within the limit, are left unchanged.
// Non-JSON lines pass through unmodified.
package jsonlimit
