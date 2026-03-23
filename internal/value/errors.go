package value

import "errors"

var (
	// ErrUnsupportedType is returned when an unsupported type is encountered.
	ErrUnsupportedType = errors.New("unsupported type")
	// ErrMustHaveElemType is returned when a list type does not have an element type defined.
	ErrMustHaveElemType = errors.New("list type must have an element type")
)
