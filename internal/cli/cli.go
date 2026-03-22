package cli

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/onyz1/infonyz"
	"github.com/onyz1/onyzify/internal/schema"
	"github.com/onyz1/onyzify/internal/types"
	"github.com/onyz1/onyzify/internal/value"
)

// FlagInput represents a command-line input that corresponds to a field defined in the schema.
type FlagInput struct {
	// Name is the name of the flag, which corresponds to the field name in the schema.
	Name string
	// Type indicates the data type of the flag input, which is used for parsing and validation.
	Type types.Type
	// Value is the actual value of the flag input, which can be set during runtime based on user input or other sources.
	Value value.Value
}

// Set parses the input string according to the specified Type and sets the Value field of the FlagInput.
//
// It is implemented to satisfy the [flag.Value] interface,
func (fi *FlagInput) Set(s string) error {
	var parsedVal *value.Value

	var err error
	if !types.IsListType(fi.Type) {
		parsedVal, err = value.ParseValue(fi.Type, s)
		if err != nil {
			return fmt.Errorf("parse value: %w", err)
		}
	} else {
		parsedVal, err = value.ParseValueList(fi.Type, strings.Split(strings.ReplaceAll(s, " ", ""), ","))
	}

	fi.Value = *parsedVal

	return nil
}

// String returns a string representation of the FlagInput's value,
// which is used for displaying the flag's value in help messages and other outputs.
//
// It is implemented to satisfy the [flag.Value] interface.
func (fi *FlagInput) String() string {
	return fi.Value.Stringify()
}

// FlagInputs is a mapping of field names to their corresponding FlagInput,
// allowing for easy access to the parsed command-line inputs based on the schema fields.
type FlagInputs map[string]*FlagInput

// ToAnyMap converts the FlagInputs into a [map[string]any] format, which can be used for further processing,
// such as validation, transformation, or integration with other components of the application.
// It iterates over each FlagInput in the map and extracts its value, creating a new map where the keys are the field names
// and the values are the corresponding flag input values.
func (fi FlagInputs) ToAnyMap() map[string]any {
	anyMap := make(map[string]any, len(fi))
	for fieldName, flagInput := range fi {
		anyMap[fieldName] = flagInput.Value.Interface()
	}

	return anyMap
}

// Build constructs a [flag.FlagSet] based on the provided [schema.CompiledSchema]
// and registers flags for each field defined in the schema.
//
// It returns the constructed [flag.FlagSet], a mapping of field names to their
// corresponding [FlagInput] instances, and any error encountered during the process.
//
// The function iterates over each field in the compiled schema and registers a
// flag for it based on its type (e.g., int, string, bool, float64).
// It also initializes a [FlagInput] for each field to store the parsed value later during the parsing phase.
// If an unsupported type is encountered, it returns an error indicating the specific field and type that caused the issue.
func Build(ctx context.Context, fs flag.FlagSet, sch schema.CompiledSchema) (*flag.FlagSet, FlagInputs, error) {
	log := infonyz.FromContext(ctx)
	debug := log.IsLevel(infonyz.DebugLevel)

	var flagInputs = make(FlagInputs, len(sch))

	for fieldName, field := range sch {
		flagInput := &FlagInput{Name: fieldName, Type: field.Type}

		if field.Default.Interface() != nil {
			flagInput.Value = field.Default
		} else {
			// Initialize the FlagInput with a zero value for its type to ensure it has a valid state before parsing.
			flagInput.Value = value.Value{Type: field.Type}
		}

		flagInputs[fieldName] = flagInput
		fs.Var(flagInput, fieldName, field.Description)

		if debug {
			log.Debug(
				"registered flag",
				infonyz.String("field", fieldName),
				infonyz.String("type", field.Type.String()),
				infonyz.String("default", field.Default.Stringify()),
				infonyz.Bool("required", field.Required),
				infonyz.Int("enum_count", len(field.Enum)),
				infonyz.F("enum_values", field.Enum.Strings()),
			)
		}
	}

	return &fs, flagInputs, nil
}

// Parse processes the command-line arguments using the provided [flag.FlagSet]
// and validates them against the [schema.CompiledSchema].
// It checks for required fields, parses the values according to their types,
// and validates enum constraints if applicable.
//
// If any required field is missing, if a value cannot be parsed according to its type,
// or if a value does not satisfy the enum constraints,
// it returns an error with details about the specific issue and the field involved.
func Parse(ctx context.Context, fs *flag.FlagSet, inputs FlagInputs, sch schema.CompiledSchema, args []string) error {
	log := infonyz.FromContext(ctx)
	debug := log.IsLevel(infonyz.DebugLevel)

	if debug {
		log.Debug(
			"parsing command-line arguments",
			infonyz.Int("arg_count", len(args)),
			infonyz.F("args", args),
		)
	}

	if fs == nil {
		return ErrNilFlagSet
	}

	if sch == nil {
		return nil
	}

	fs.Usage = func() {
		out := fs.Output()

		bold := "\033[1m"
		reset := "\033[0m"
		color := "\033[38;5;90m"
		gray := "\033[38;5;244m"

		fmt.Fprintf(out, "%s%sUsage%s\n", bold, color, reset)
		fmt.Fprintf(out, "  %s%s [options]%s\n\n", gray, fs.Name(), reset)

		printGroup := func(title string, required bool) {
			fmt.Fprintf(out, "%s%s%s\n", bold+color, title, reset)

			for name, field := range sch {
				if field.Required != required {
					continue
				}

				meta := ""
				if def := field.Default.Stringify(); def != "" {
					meta += fmt.Sprintf(" (default: %s)", def)
				}
				if enum := field.Enum; len(enum) > 0 {
					meta += fmt.Sprintf(" (allowed: %s)", field.Enum.Strings())
				}

				fmt.Fprintf(out, "  -%s %s%s%-12s %s%s%s%s\n", name, gray, field.Type.String(), reset, field.Description, gray, meta, reset)
			}
			fmt.Fprintln(out)
		}

		printGroup("Required", true)
		printGroup("Optional", false)
		fmt.Fprint(out, reset)
	}

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	// Create a set of flags that were actually provided in the command-line arguments for validation purposes.
	setFlags := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		setFlags[f.Name] = true
	})

	for fieldName, field := range sch {
		if debug {
			log.Debug(
				"validating field",
				infonyz.String("field", fieldName),
				infonyz.String("type", field.Type.String()),
				infonyz.Bool("required", field.Required),
				infonyz.Int("enum_count", len(field.Enum)),
				infonyz.F("enum_values", field.Enum.Strings()),
			)
		}

		flagInput, ok := inputs[fieldName]
		if !ok {
			fs.Usage()
			return fmt.Errorf("field %q: no corresponding flag input found: %w", fieldName, ErrFlagInputNotFound)
		}

		_, isSet := setFlags[fieldName]

		if err := field.CheckVal(&flagInput.Value, isSet); err != nil {
			fs.Usage()
			return fmt.Errorf("field %q: check value validity: %w", fieldName, err)
		}

		if debug {
			log.Debug(
				"parsed flag",
				infonyz.String("field", fieldName),
				infonyz.String("type", field.Type.String()),
				infonyz.String("value", flagInput.Value.Stringify()),
			)
		}
	}

	return nil
}
