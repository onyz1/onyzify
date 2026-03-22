package schema

import "errors"

var (
	ErrFieldNameRequired = errors.New("field name is required")
	ErrFieldTypeRequired = errors.New("field type is required")
	ErrDefaultNotInEnum  = errors.New("default value is not in enum values")
	ErrValueRequired     = errors.New("value is required")
	ErrValueNotInEnum    = errors.New("value is not in enum values")
)
