package onyzify

import (
	"context"
	"fmt"

	"github.com/onyz1/infonyz"
	"github.com/onyz1/onyzify/internal/cli"
	"github.com/onyz1/onyzify/internal/formatter"
	"github.com/onyz1/onyzify/internal/wizard"
)

// Engine represents the core component of the onyzify application,
// responsible for processing the provided options and executing the main logic of the application.
type Engine struct {
	opts *Options
}

// New creates a new instance of the Engine using the provided Options.
// It validates the options and returns an error if any required fields are missing or invalid.
// If the options are valid, it initializes and returns a new Engine instance.
func New(opts *Options) (*Engine, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %w", err)
	}

	return &Engine{opts: opts}, nil
}

// Run executes the main logic of the Engine. It processes the provided options, loads the schema,
// and either builds a CLI or runs a wizard based on the options. It returns the result of the execution
// or an error if any step fails during the execution.
func (e *Engine) Run(parent context.Context) (*Result, error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	ctx = infonyz.WithLogger(ctx, e.opts.GetLogger())

	compiledSch, err := loadSchema(ctx, e.opts.SchemaData)
	if err != nil {
		return nil, fmt.Errorf("load schema: %w", err)
	}

	var inputs map[string]any
	if !e.opts.Wizard {
		fs, flagInputs, err := cli.Build(ctx, e.opts.FlagSet, compiledSch)
		if err != nil {
			return nil, fmt.Errorf("build CLI: %w", err)
		}

		if err := cli.Parse(ctx, compiledSch, flagInputs, formatter.UsageFormatter, fs, e.opts.Args); err != nil {
			return nil, fmt.Errorf("parse CLI: %w", err)
		}
		inputs = flagInputs.ToAnyMap()
	} else {
		wizInputs, err := wizard.Run(
			compiledSch,
			formatter.StructuredFormatter,
			e.opts.WizardOptions.Dst,
			e.opts.WizardOptions.Src,
		)
		if err != nil {
			return nil, fmt.Errorf("run wizard: %w", err)
		}
		inputs = wizInputs.ToAnyMap()
	}

	return &Result{
		inputs: inputs,
	}, nil
}
