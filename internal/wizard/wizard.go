package wizard

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/onyz1/onyzify/internal/schema"
	"github.com/onyz1/onyzify/internal/types"
	"github.com/onyz1/onyzify/internal/value"
)

// WizardInput represents a single input field in the wizard, containing the field name and its corresponding value.
type WizardInput struct {
	// Name is the name of the field, which corresponds to the field name defined in the schema.
	Name string
	// Type indicates the data type of the input, which is used for parsing and validation.
	Type types.Type
	// Value is the actual value of the input, which can be set during runtime based
	// on user input or other sources. It is of type [value.Value],
	// which allows for flexible handling of different data types as defined in the schema.
	Value value.Value
}

// Set parses the input string according to the specified Type and sets the Value field of the WizardInput.
func (wi *WizardInput) Set(s string) error {
	var parsedVal *value.Value

	var err error
	if !types.IsListType(wi.Type) {
		parsedVal, err = value.ParseValue(wi.Type, s)
		if err != nil {
			return fmt.Errorf("parse value: %w", err)
		}
	} else {
		parsedVal, err = value.ParseValueList(wi.Type, strings.Split(strings.ReplaceAll(s, " ", ""), ","))
	}

	wi.Value = *parsedVal

	return nil
}

// WizardInputs is a mapping of field names to their corresponding [WizardInput] instances,
// allowing for easy access to the parsed inputs based on the schema fields.
type WizardInputs map[string]*WizardInput

// ToAnyMap converts the [WizardInputs] mapping into a standard [map[string]any] format,
// where the keys are the field names and the values are the actual input values.
// This is useful for further processing or outputting the parsed inputs in a more generic format.
func (wi WizardInputs) ToAnyMap() map[string]any {
	anyMap := make(map[string]any, len(wi))
	for key, input := range wi {
		anyMap[key] = input.Value.Interface()
	}

	return anyMap
}

// Run executes the wizard by prompting the user for input based on the provided compiled schema.
// It uses the provided formatter function to generate prompts for each field and reads user input from the specified source.
// The collected inputs are returned as a [WizardInputs] mapping, or an error if any issues occur during the process.
func Run(sch schema.CompiledSchema, formatter Formatter, dst io.Writer, src io.Reader) (WizardInputs, error) {
	var inputs = make(WizardInputs, len(sch))

	for fieldName, field := range sch {
		input := &WizardInput{Name: fieldName, Type: field.Type}

		if field.Default.Interface() != nil {
			input.Value = field.Default
		} else {
			// Initialize the Value field with the correct Type to ensure it is ready for parsing user input.
			input.Value = value.Value{Type: field.Type}
		}

		inputs[fieldName] = input

		reader := bufio.NewReader(src)

		var isSet bool
		for {
			fmt.Fprint(dst, formatter(field))

			userInput, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("read user input: %w", err)
			}

			userInput = strings.TrimSpace(userInput)

			if userInput == "" {
				break
			}

			err = input.Set(userInput)
			if err != nil {
				fmt.Fprintf(dst, "Invalid input: %v. Please try again.\n", err)
				continue
			}

			isSet = true

			break
		}

		if err := field.CheckVal(&input.Value, isSet); err != nil {
			return nil, fmt.Errorf("field: %q: check value validity: %w", fieldName, err)
		}
	}

	return inputs, nil
}
