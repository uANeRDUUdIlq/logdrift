// Package highlight applies configurable keyword-based color highlighting
// to rendered log lines.
//
// Rules map plain-text keywords to named terminal colors (red, yellow, green,
// blue, magenta, cyan). All occurrences of a keyword within a line are
// replaced with their ANSI-colored equivalents before output.
//
// Example usage:
//
//	h := highlight.New([]highlight.Rule{
//		{Word: "ERROR", Color: "red"},
//		{Word: "WARN",  Color: "yellow"},
//	})
//	fmt.Println(h.Apply(line))
package highlight
