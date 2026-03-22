package cli

import "errors"

var (
	ErrFlagInputNotFound    = errors.New("flag input not found for field")
	ErrNilFlagSet           = errors.New("nil flag set")
	ErrUnsupportedType      = errors.New("unsupported type")
	ErrMissingRequiredField = errors.New("missing required field")
	ErrValueNotInEnum       = errors.New("value not in enum")
	ErrSingleCharExpected   = errors.New("single character expected")
)
