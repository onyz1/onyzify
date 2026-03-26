package env

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Inputs represents a map of input values, where the key is the field name and the value is the corresponding input value.
type Inputs map[string]any

// Build takes a map of input values and constructs a string representation of environment variable definitions.
// Each key in the inputs map is transformed into an environment variable name (uppercase with spaces replaced by underscores),
// and the corresponding value is set as the value of that environment variable.
// The resulting string contains lines in the format "ENV_VAR_NAME=value".
//
// For array values, the elements are joined with commas. For [time.Time] values, they are formatted in [time.RFC3339] format.
func Build(inputs Inputs) string {
	var sb strings.Builder

	for fieldName, value := range inputs {
		envVarName := strings.ToUpper(strings.ReplaceAll(fieldName, " ", "_"))
		sb.WriteString(envVarName)
		sb.WriteString("=")

		switch v := value.(type) {
		case []any:
			for i, val := range v {
				if i > 0 {
					sb.WriteString(",")
				}
				fmt.Fprintf(&sb, "%v", val)
			}

		case time.Time:
			fmt.Fprintf(&sb, "%v", v.Format(time.RFC3339))

		default:
			fmt.Fprintf(&sb, "%v", value)
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// BuildFile takes a map of input values and a file path, constructs the environment variable definitions using the [Build] function,
// and writes the resulting string to a file at the specified path.
//
// If there is an error during the file writing process, it returns an error with a descriptive message.
func BuildFile(inputs Inputs, path string) error {
	content := Build(inputs)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}

	return nil
}
