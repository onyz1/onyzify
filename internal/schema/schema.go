package schema

import (
	"context"
	"fmt"
)

// Schema represents the overall structure of the schema configuration,
// containing a map of field names to their corresponding Field definitions.
type Schema map[string]*Field

// CompiledSchema represents a Schema that has been processed,
// with all fields compiled and ready for use.
type CompiledSchema map[string]*CompiledField

// Compile processes the Schema's configuration,
// validates its properties, and returns a CompiledSchema
// instance that contains the compiled fields ready for use.
func (s Schema) Compile(parent context.Context) (CompiledSchema, error) {
	compiled := make(CompiledSchema, len(s))

	for fieldName, field := range s {
		cf, err := field.Compile(parent)
		if err != nil {
			return nil, fmt.Errorf("compile field %q: %w", fieldName, err)
		}
		compiled[fieldName] = cf
	}

	return compiled, nil
}
