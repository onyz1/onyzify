package onyzify

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/onyz1/infonyz"
)

// WizardOptions represents the configuration options for the wizard mode of the onyzify application.
type WizardOptions struct {
	Dst io.Writer
	Src io.Reader
}

// Options represents the configuration options for the onyzify application.
type Options struct {
	// SchemaData is the raw byte data of the schema,
	// which can be used to define the structure and validation rules for the application.
	SchemaData []byte
	// Args is a slice of strings representing the command-line arguments passed to the application.
	Args []string
	// Wizard is a boolean flag that indicates whether to enable wizard mode, which provides an interactive CLI experience.
	Wizard bool
	// WizardOptions is an instance of [WizardOptions],
	// which is used to configure the wizard mode.
	WizardOptions WizardOptions

	// FlagSet is an instance of [flag.FlagSet],
	// which is used for defining and parsing command-line flags.
	// If not provided, it will be initialized with a default FlagSet.
	FlagSet flag.FlagSet
	// Logger is an instance of [infonyz.Logger],
	// which is used for logging messages, errors, and other information throughout the application.
	// If not provided, it will be initialized with a NoopLogger that discards all log messages.
	Logger infonyz.Logger
}

// Validate checks that [Options.SchemaData] is set, [Options.Args] is provided in CLI mode,
// and [Options.WizardOptions.Dst]/[Options.WizardOptions.Src] are provided in wizard mode.
func (o *Options) Validate() error {
	if o.SchemaData == nil {
		return ErrNilSchemaData
	}

	if o.Args == nil && !o.Wizard {
		return ErrNilArgs
	}

	if o.Wizard && (o.WizardOptions.Dst == nil || o.WizardOptions.Src == nil) {
		return ErrNilWizardOptions
	}

	return nil
}

// GetLogger returns the Logger instance from the Options.
// If the Logger is nil, it returns a NoopLogger that discards all log messages.
func (o *Options) GetLogger() infonyz.Logger {
	if o.Logger == nil {
		return infonyz.NoopLogger()
	}
	return o.Logger
}

// WithSchemaFile reads the schema data from the specified file path and
// returns a new [Options] instance with the [Options.SchemaData] field populated.
// If there is an error reading the file, it returns an error with a descriptive message.
//
// Example usage:
//
//	opts := &Options{Args: os.Args[1:]}
//	opts, err := opts.WithSchemaFile("schema.yaml")
//	if err != nil {
//	    // handle error
//	}
//
//	if err := opts.Validate(); err != nil {
//	    // handle validation error
//	}
func (o Options) WithSchemaFile(path string) (*Options, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read schema file: %w", err)
	}

	o.SchemaData = data

	return &o, nil
}
