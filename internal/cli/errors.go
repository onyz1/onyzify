package cli

import "errors"

var (
	// ErrFlagInputNotFound is returned when a flag input cannot be found for a struct field.
	ErrFlagInputNotFound = errors.New("flag input not found for field")
)
