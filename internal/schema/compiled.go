package schema

import (
	"fmt"
	"slices"

	"github.com/onyz1/onyzify/internal/types"
	"github.com/onyz1/onyzify/internal/value"
)

// enumValue represents a slice of [value.Value] instances that correspond to the enum values defined for a field.
type enumValue []value.Value

// Strings returns a slice of string representations of the enum values.
func (e enumValue) Strings() []string {
	strs := make([]string, len(e))
	for i, enumVal := range e {
		strs[i] = enumVal.Stringify()
	}
	return strs
}

// CompiledField represents a Field that has been processed and validated,
// with its type and values parsed and ready for use.
type CompiledField struct {
	// Name is the name of the field, which serves as the key in the schema.
	Name string
	// Type is the parsed type of the field, represented as a [types.Type] value.
	Type types.Type

	// Required indicates whether the field is mandatory or optional.
	Required bool
	// Description provides a human-readable explanation of the field's purpose and usage.
	Description string

	// Default holds the default value for the field, if specified,
	// parsed as a [value.Value] based on the field's type.
	Default value.Value
	// Enum contains a list of allowed values for the field, if specified,
	// each parsed as a [value.Value] based on the field's type.
	Enum enumValue
}

// CheckVal validates the provided value against the field's constraints, such as required status and enum values.
// It returns an error if the value is invalid according to the field's properties.
//
// isSet indicates whether the value was explicitly set by the user (e.g., via a flag or wizard input).
// This is important for distinguishing between zero values and missing values, especially for required fields.
func (f *CompiledField) CheckVal(val *value.Value, isSet bool) error {
	if !isSet && f.Required && val.IsZero() {
		return fmt.Errorf("value is required: %w", ErrValueRequired)
	}

	if len(f.Enum) > 0 {
		if !slices.ContainsFunc(f.Enum, val.Equal) {
			return fmt.Errorf("value %q is not in enum %v: %w", val.Stringify(), f.Enum.Strings(), ErrValueNotInEnum)
		}
	}

	return nil
}

// CheckValString is a helper method that takes a string input, parses it according to the field's type,
// and then validates it using the CheckVal method.
// This is useful for validating raw string inputs from command-line flags
// or wizard prompts before they are parsed into value.Value instances.
//
// Check [CompiledField.CheckVal] for more details on the validation logic and error handling.
func (f *CompiledField) CheckValString(val string, isSet bool) error {
	// If the input string is empty, we treat it as a zero value and check it directly.
	if val == "" {
		return f.CheckVal(&value.Value{}, isSet)
	}

	parsedVal, err := value.ParseValue(f.Type, val)
	if err != nil {
		return fmt.Errorf("parse value: %w", err)
	}

	return f.CheckVal(parsedVal, isSet)
}
