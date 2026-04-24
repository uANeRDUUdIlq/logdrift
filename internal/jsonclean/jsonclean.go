// Package jsonclean removes null, empty-string, and empty-object/array
// fields from a JSON log line before it reaches downstream processors.
package jsonclean

import (
	"encoding/json"
)

// Options controls which kinds of empty values are stripped.
type Options struct {
	// RemoveNull removes fields whose value is JSON null.
	RemoveNull bool
	// RemoveEmptyString removes fields whose value is "".
	RemoveEmptyString bool
	// RemoveEmptyObject removes fields whose value is {}.
	RemoveEmptyObject bool
	// RemoveEmptyArray removes fields whose value is [].
	RemoveEmptyArray bool
}

// Cleaner strips unwanted empty values from JSON log lines.
type Cleaner struct {
	opts Options
}

// New returns a Cleaner configured with opts.
// If no removal flags are set the Cleaner is a no-op.
func New(opts Options) *Cleaner {
	return &Cleaner{opts: opts}
}

// Apply returns the line with matching empty fields removed.
// Non-JSON lines are returned unchanged.
func (c *Cleaner) Apply(line string) string {
	if !c.opts.RemoveNull && !c.opts.RemoveEmptyString &&
		!c.opts.RemoveEmptyObject && !c.opts.RemoveEmptyArray {
		return line
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for k, v := range obj {
		if c.opts.RemoveNull && string(v) == "null" {
			delete(obj, k)
			continue
		}
		if c.opts.RemoveEmptyString && string(v) == `""` {
			delete(obj, k)
			continue
		}
		if c.opts.RemoveEmptyObject && string(v) == "{}" {
			delete(obj, k)
			continue
		}
		if c.opts.RemoveEmptyArray && string(v) == "[]" {
			delete(obj, k)
			continue
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
