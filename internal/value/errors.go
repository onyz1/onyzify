package value

import "errors"

var (
	// ErrUnsupportedType is returned when an unsupported type is encountered.
	ErrUnsupportedType = errors.New("unsupported type")
	// ErrMustHaveElemType is returned when a type that requires an element type (like a list) is missing the Elem field.
	ErrMustHaveElemType = errors.New("must have an element type")
)
