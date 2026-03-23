package cli

import (
	"flag"
	"fmt"

	"github.com/onyz1/onyzify/internal/formatter"
	"github.com/onyz1/onyzify/internal/io"
	"github.com/onyz1/onyzify/internal/schema"
)

func Build(fs flag.FlagSet, sch schema.CompiledSchema) (*flag.FlagSet, io.Inputs, error) {
	inputs := io.New(len(sch))

	for fieldName, field := range sch {
		flagInput := &io.Input{
			Name:  fieldName,
			Type:  field.Type,
			Value: field.Default,
		}
		inputs[fieldName] = flagInput

		fs.Var(flagInput, fieldName, field.Description)
	}

	return &fs, inputs, nil
}

func Parse(sch schema.CompiledSchema, inputs io.Inputs, formatter formatter.Formatter, fs *flag.FlagSet, args []string) error {
	fs.Usage = func() {
		usage(sch, formatter, fs)
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
	}

	return nil
}
