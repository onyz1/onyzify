package schema

import (
	"fmt"

	"github.com/onyz1/onyzify/internal/yaml"
)

// Load parses the schema data and returns a Schema representation of it.
func Load(data []byte) (Schema, error) {
	var sch Schema
	if err := yaml.Load(data, &sch); err != nil {
		return nil, fmt.Errorf("load schema: %w", err)
	}

	for fieldName := range sch {
		sch[fieldName].Name = fieldName
	}

	return sch, nil

}

// LoadFile reads a schema file from the specified path and returns a Schema representation of its contents.
func LoadFile(path string) (Schema, error) {
	var sch Schema

	if err := yaml.LoadFile(path, &sch); err != nil {
		return nil, fmt.Errorf("load schema file: %w", err)
	}

	for fieldName := range sch {
		sch[fieldName].Name = fieldName
	}

	return sch, nil
}
