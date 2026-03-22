package onyzify

import (
	"context"
	"fmt"

	"github.com/onyz1/onyzify/internal/schema"
)

// loadSchema loads a schema from the provided byte data and returns a [schema.CompiledSchema] instance.
func loadSchema(parent context.Context, data []byte) (schema.CompiledSchema, error) {
	sch, err := schema.Load(data)
	if err != nil {
		return nil, fmt.Errorf("load schema: %w", err)
	}

	compiledSch, err := sch.Compile(parent)
	if err != nil {
		return nil, fmt.Errorf("compile schema: %w", err)
	}

	return compiledSch, nil
}
