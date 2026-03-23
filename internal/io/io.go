package io

import (
	"fmt"

	"github.com/onyz1/onyzify/internal/types"
	"github.com/onyz1/onyzify/internal/value"
)

// Input represents a single input field, containing the field name, its corresponding type, and the actual value.
type Input struct {
	// Name is the name of the field, which corresponds to the field name defined in the schema.
	Name string
	// Type indicates the data type of the input, which is used for parsing and validation.
	Type types.Type
	// Value is the actual value of the input,
	// which can be set during runtime based on user input or other sources. It is of type [value.Value],
	// which allows for flexible handling of different data types as defined in the schema.
	Value value.Value
}

// Set parses the input string according to the specified Type and sets the Value field of the Input.
func (i *Input) Set(s string) error {
	parsedVal, err := value.ParseValue(i.Type, s)
	if err != nil {
		return fmt.Errorf("parse value: %w", err)
	}

	i.Value = *parsedVal

	return nil
}

// String returns a string representation of the Input's value.
func (i *Input) String() string {
	return i.Value.Stringify()
}

// Inputs is a mapping of field names to their corresponding [Input] instances,
// allowing for easy access to the parsed inputs based on the schema fields.
type Inputs map[string]*Input

// ToAnyMap converts the [Inputs] mapping into a standard [map[string]any] format,
// where the keys are the field names and the values are the actual input values.
// This is useful for further processing or outputting the parsed inputs in a more generic format.
func (i Inputs) ToAnyMap() map[string]any {
	anyMap := make(map[string]any, len(i))
	for key, input := range i {
		anyMap[key] = input.Value.Interface()
	}

	return anyMap
}

// New creates a new [Inputs] mapping with the specified size,
// which can be used to store parsed inputs based on the schema fields.
func New(size int) Inputs {
	return make(Inputs, size)
}
