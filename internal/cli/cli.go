package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/onyz1/infonyz"
	"github.com/onyz1/onyzify/internal/formatter"
	"github.com/onyz1/onyzify/internal/io"
	"github.com/onyz1/onyzify/internal/schema"
)

func Build(ctx context.Context, fs flag.FlagSet, sch schema.CompiledSchema) (*flag.FlagSet, io.Inputs, error) {
	log := infonyz.FromContext(ctx)
	debug := log.IsLevel(infonyz.DebugLevel)

	inputs := io.New(len(sch))

	for fieldName, field := range sch {
		flagInput := &io.Input{
			Name:  fieldName,
			Type:  field.Type,
			Value: field.Default,
		}
		inputs[fieldName] = flagInput

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

	return &fs, inputs, nil
}

func Parse(ctx context.Context, sch schema.CompiledSchema, inputs io.Inputs, formatter formatter.Formatter, fs *flag.FlagSet, args []string) error {
	log := infonyz.FromContext(ctx)
	debug := log.IsLevel(infonyz.DebugLevel)

	if debug {
		log.Debug(
			"parsing command-line arguments",
			infonyz.Int("arg_count", len(args)),
			infonyz.F("args", args),
		)
	}

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
