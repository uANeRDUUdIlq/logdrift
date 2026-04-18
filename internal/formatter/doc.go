// Package formatter provides rendering logic for structured JSON log lines.
//
// Each service gets a Formatter instance with a distinct ANSI color assigned
// in round-robin order. Two output formats are supported:
//
//   - FormatPretty: parses the JSON and prints a human-friendly single line
//     with timestamp, level, message, and any additional fields.
//   - FormatRaw: prefixes the original line with the service tag and skips
//     any JSON parsing.
//
// Example:
//
//	f := formatter.New("api-gateway", formatter.FormatPretty)
//	fmt.Println(f.Render(`{"level":"info","msg":"request received","latency_ms":42}`))
package formatter
