package wizard

import (
	"fmt"
	"strings"

	"github.com/onyz1/onyzify/internal/schema"
	"github.com/onyz1/onyzify/internal/types"
)

// Formatter is a function type that takes a field name and a
// pointer to a [schema.CompiledField] and returns a formatted string.
// This function is used to generate prompts for each field when
// running the wizard, allowing for customizable formatting based on the field's properties.
type Formatter func(field *schema.CompiledField) string

// DefaultFormatter is a built-in implementation of the [Formatter] type that generates a prompt string
// based on the properties of the provided [schema.CompiledField]. It includes the field name, type,
// whether it's required, and any optional metadata such as description, default value, and allowed enum values.
// The generated prompt is designed to be informative and user-friendly when displayed during the wizard execution.
func DefaultFormatter(field *schema.CompiledField) string {
	var b strings.Builder

	// Header
	if !types.IsListType(field.Type) {
		fmt.Fprintf(&b, "Field: %s\n", field.Name)
	} else {
		fmt.Fprintf(&b, "Field: %s (comma separated)\n", field.Name)
	}
	fmt.Fprintf(&b, "  Type: %s\n", field.Type.String())
	fmt.Fprintf(&b, "  Required: %t\n", field.Required)

	// Optional metadata
	if field.Description != "" {
		fmt.Fprintf(&b, "  Description: %s\n", field.Description)
	}

	if field.Default.Stringify() != "" {
		fmt.Fprintf(&b, "  Default: %s\n", field.Default.Stringify())
	}

	if len(field.Enum) > 0 {
		fmt.Fprintf(&b, "  Allowed values: %s\n", strings.Join(field.Enum.Strings(), ", "))
	}

	// Input hint
	b.WriteString("Enter value: ")

	return b.String()
}
