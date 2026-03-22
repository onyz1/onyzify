package onyzify

import "errors"

var (
	ErrNilSchemaData    = errors.New("schema data cannot be nil")
	ErrNilArgs          = errors.New("args cannot be nil")
	ErrNilWizardOptions = errors.New("wizard options cannot be nil when wizard mode is enabled")
)
