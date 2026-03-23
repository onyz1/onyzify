package onyzify

import "errors"

var (
	// ErrNilSchemaData is returned when the schema data provided to [Options] is nil.
	ErrNilSchemaData = errors.New("schema data cannot be nil")
	// ErrNilArgs is returned when the args provided to [Options] are nil.
	ErrNilArgs = errors.New("args cannot be nil")
	// ErrNilWizardOptions is returned when the wizard options provided to [Options] are nil while [Options.Wizard] is true.
	ErrNilWizardOptions = errors.New("wizard options cannot be nil when wizard mode is enabled")
)
