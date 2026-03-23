package types

import "errors"

var (
	// ErrUnknownType is returned when an unknown type is encountered.
	ErrUnknownType = errors.New("unknown type")
)
