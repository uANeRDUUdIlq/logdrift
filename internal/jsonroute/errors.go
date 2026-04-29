package jsonroute

import "fmt"

// ErrEmptyField is returned when a rule has a blank Field value.
type ErrEmptyField struct{ Index int }

func (e *ErrEmptyField) Error() string {
	return fmt.Sprintf("jsonroute: rule %d has an empty field", e.Index)
}

// ErrNilWriter is returned when a rule has a nil Writer.
type ErrNilWriter struct{ Index int }

func (e *ErrNilWriter) Error() string {
	return fmt.Sprintf("jsonroute: rule %d has a nil writer", e.Index)
}
