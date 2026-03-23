package yaml

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// YAMLOutput is the output target type accepted by the YAML helpers.
type YAMLOutput any

// Load parses the YAML data and populates the provided output structure with the corresponding values.
func Load(data []byte, out YAMLOutput) error {
	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	return nil
}

// LoadFile reads a YAML file from the specified path and populates
// the provided output structure with the corresponding values.
func LoadFile(path string, out YAMLOutput) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read yaml file: %w", err)
	}

	return Load(data, out)
}

// YAMLInputs represents a collection of YAML inputs,
// where each key is a field name and the corresponding value
// is the value to be serialized in YAML format.
type YAMLInputs map[string]any

// Save converts the YAMLInputs into a YAML byte slice.
func Save(input YAMLInputs) ([]byte, error) {
	data, err := yaml.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal yaml: %w", err)
	}

	return data, nil
}

// SaveFile takes a file path and YAMLInputs, converts the inputs into YAML format,
// and writes the resulting data to the specified file.
func SaveFile(path string, input YAMLInputs) error {
	data, err := Save(input)
	if err != nil {
		return fmt.Errorf("save yaml: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("write yaml file: %w", err)
	}

	return nil
}
