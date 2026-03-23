package schema

import "errors"

var (
	// ErrFieldNameRequired is returned when a field name is missing in the schema definition.
	ErrFieldNameRequired = errors.New("field name is required")
	// ErrFieldTypeRequired is returned when a field type is missing in the schema definition.
	ErrFieldTypeRequired = errors.New("field type is required")
	// ErrDefaultNotInEnum is returned when a default value is provided that is not in the enum values.
	ErrDefaultNotInEnum = errors.New("default value is not in enum values")
	// ErrValueRequired is returned when a value is required but not provided.
	ErrValueRequired = errors.New("value is required")
	// ErrValueNotInEnum is returned when a provided value is not in the enum values.
	ErrValueNotInEnum = errors.New("value is not in enum values")
)
