package onyzify

import (
	"fmt"

	"github.com/onyz1/onyzify/internal/env"
	"github.com/onyz1/onyzify/internal/yaml"
)

// Result represents the output of the onyzify engine after processing the inputs and schema.
type Result struct {
	inputs map[string]any
}

// String returns a string representation of the parsed inputs in Result.
// YAML converts the parsed inputs into YAML-formatted bytes.
// SaveYAML saves the parsed inputs as YAML to the specified file.
func (r *Result) String() string {
	return fmt.Sprintf("%v", r.inputs)
}

// YAML converts the map representation of the parsed command-line inputs
// into a YAML-formatted byte slice. If there is an error during the conversion process,
// it returns an error with a descriptive message.
func (r *Result) YAML() ([]byte, error) {
	data, err := yaml.Save(r.inputs)
	if err != nil {
		return nil, fmt.Errorf("convert to YAML: %w", err)
	}

	return data, nil
}

// SaveYAML saves the YAML representation of the parsed command-line inputs to a file at the specified path.
// If there is an error during the file saving process, it returns an error with a descriptive message.
func (r *Result) SaveYAML(path string) error {
	if err := yaml.SaveFile(path, r.inputs); err != nil {
		return fmt.Errorf("save YAML file: %w", err)
	}
	return nil

}

// ENV converts the map representation of the parsed command-line inputs into a byte slice formatted as environment variable definitions.
// Each key in the inputs map is transformed into an environment variable name (uppercase with spaces replaced by underscores),
// and the corresponding value is set as the value of that environment variable.
func (r *Result) ENV() []byte {
	return []byte(env.Build(r.inputs))
}

// SaveENV saves the environment variable definitions generated from the parsed command-line inputs to a file at the specified path.
// If there is an error during the file saving process, it returns an error with a descriptive message.
func (r *Result) SaveENV(path string) error {
	if err := env.BuildFile(r.inputs, path); err != nil {
		return fmt.Errorf("save env file: %w", err)
	}
	return nil
}
