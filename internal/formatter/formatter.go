package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/onyz1/onyzify/internal/schema"
)

// Formatter defines a function type that takes a [schema.CompiledField] and an [io.Writer] as parameters.
// The function is responsible for generating a formatted prompt for the given field and writing it to the provided writer.
// This allows for customizable formatting of prompts when collecting user input in a wizard-like interface.
type Formatter func(field *schema.CompiledField, dst io.Writer)

// StructuredFormatter is a concrete implementation of the [Formatter] type that generates a structured prompt for a given field.
// It includes the field name, type, whether it's required, and any additional metadata such as description, default value, and allowed enum values.
// The prompt is formatted in a clear and organized manner to provide users with all necessary information about the field before they input their value.
func StructuredFormatter(field *schema.CompiledField, dst io.Writer) {
	// Header
	fmt.Fprintf(dst, "Field: %s\n", field.Name)
	fmt.Fprintf(dst, "  Type: %s\n", field.Type.String())
	fmt.Fprintf(dst, "  Required: %t\n", field.Required)

	// Optional metadata
	if field.Description != "" {
		fmt.Fprintf(dst, "  Description: %s\n", field.Description)
	}

	if field.Default.Stringify() != "" {
		fmt.Fprintf(dst, "  Default: %s\n", field.Default.Stringify())
	}

	if len(field.Enum) > 0 {
		fmt.Fprintf(dst, "  Allowed values: %s\n", strings.Join(field.Enum.Strings(), ", "))
	}

	// Input hint
	fmt.Fprintf(dst, "Enter value: ")
}

// UsageFormatter is another implementation of the [Formatter] type that generates a concise prompt for a given field.
// It includes the field name, type, description, and any relevant metadata such as default value and allowed enum values.
// The prompt is designed to be compact while still providing users with the necessary information to input their value correctly.
func UsageFormatter(field *schema.CompiledField, dst io.Writer) {
	fmt.Fprintf(dst, "%s %-12s %s", field.Name, field.Type.String(), field.Description)

	if def := field.Default.Stringify(); def != "" {
		fmt.Fprintf(dst, " (default: %s)", def)
	}

	if enum := field.Enum.Strings(); len(enum) > 0 {
		fmt.Fprintf(dst, " (allowed: [%s])", strings.Join(enum, ", "))
	}

	fmt.Fprintln(dst)
}
