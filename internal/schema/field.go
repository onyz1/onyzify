package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/onyz1/infonyz"
	"github.com/onyz1/onyzify/internal/types"
	"github.com/onyz1/onyzify/internal/value"
	"gopkg.in/yaml.v3"
)

// Field represents the structure of a field as defined in the schema configuration.
type Field struct {
	// Name is the name of the field, which serves as the key in the schema.
	Name string
	// Type is the string representation of the field's type,
	// which will be parsed into a [types.Type] value during compilation.
	Type string `yaml:"type"`

	// Required indicates whether the field is mandatory or optional.
	Required bool `yaml:"required,omitempty"`
	// Description provides a human-readable explanation of the field's purpose and usage.
	Description string `yaml:"description,omitempty"`

	// Default holds the default value for the field, if specified,
	// as a string that will be parsed into a [value.Value] based on the field's type during compilation.
	Default string `yaml:"default,omitempty"`
	// Enum contains a list of allowed values for the field, if specified,
	// as strings that will be parsed into [value.Value] instances based on the field's type during compilation.
	Enum []string `yaml:"enum,omitempty"`
}

// UnmarshalYAML implements the [yaml.Unmarshaler] interface for the Field struct.
// It allows the [Field] struct to be properly populated when unmarshaling YAML data.
//
// The method handles the decoding of the YAML node into the Field struct's properties,
// including special handling for the [Field.Default] and [Field.Enum] fields to ensure they are stored as strings.
func (f *Field) UnmarshalYAML(value *yaml.Node) error {
	// temporary struct to decode the YAML data into, allowing us to handle the Default and Enum fields as any type
	type anyField struct {
		Name        string `yaml:"-"`
		Type        string `yaml:"type"`
		Required    bool   `yaml:"required,omitempty"`
		Description string `yaml:"description,omitempty"`
		Default     any    `yaml:"default,omitempty"`
		Enum        []any  `yaml:"enum,omitempty"`
	}

	var alias anyField
	if err := value.Decode(&alias); err != nil {
		return fmt.Errorf("decode field: %w", err)
	}

	f.Name = alias.Name
	f.Type = alias.Type
	f.Required = alias.Required
	f.Description = alias.Description

	if alias.Default != nil {
		switch v := alias.Default.(type) {
		case string:
			f.Default = v

		case []any:
			out, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("marshal default value: %w", err)
			}
			f.Default = string(out)

		default:
			return fmt.Errorf("type: %T: %w", alias.Default, ErrUnsupportedDefaultType)
		}
	}

	if len(alias.Enum) > 0 {
		for _, item := range alias.Enum {
			switch v := item.(type) {
			case string:
				f.Enum = append(f.Enum, v)

			case []any:
				out, err := json.Marshal(v)
				if err != nil {
					return fmt.Errorf("marshal enum value: %w", err)
				}
				f.Enum = append(f.Enum, string(out))

			default:
				return fmt.Errorf("type: %T: %w", item, ErrUnsupportedEnumType)
			}
		}
	}

	return nil
}

// Compile processes the Field's configuration,
// validates its properties, and returns a [CompiledField]
// instance that contains the parsed type and values ready for use.
func (f *Field) Compile(ctx context.Context) (*CompiledField, error) {
	log := infonyz.FromContext(ctx)
	debug := log.IsLevel(infonyz.DebugLevel)

	if debug {
		log.Debug(
			"compiling field",
			infonyz.String("field", f.Name),
			infonyz.String("type", f.Type),
			infonyz.String("default", f.Default),
			infonyz.Bool("required", f.Required),
			infonyz.Int("enum_count", len(f.Enum)),
			infonyz.F("enum_values", f.Enum),
		)
	}

	if f.Name == "" {
		return nil, ErrFieldNameRequired
	}

	if f.Type == "" {
		return nil, ErrFieldTypeRequired
	}

	if len(f.Enum) > 0 && f.Default != "" {
		if !slices.Contains(f.Enum, f.Default) {
			return nil, fmt.Errorf("default value %q is not in enum %v: %w", f.Default, f.Enum, ErrDefaultNotInEnum)
		}
	}

	compiled := new(CompiledField)

	compiled.Name = f.Name

	var err error

	parsedType, err := types.ParseType(f.Type)
	if err != nil {
		return nil, fmt.Errorf("parse type: %w", err)
	}

	compiled.Type = *parsedType
	compiled.Required = f.Required
	compiled.Description = f.Description

	if f.Default != "" {
		val, err := value.ParseValue(compiled.Type, f.Default)
		if err != nil {
			return nil, fmt.Errorf("parse default value: %w", err)
		}
		compiled.Default = *val
	}

	if len(f.Enum) > 0 {
		for _, enumStr := range f.Enum {
			enumVal, err := value.ParseValue(compiled.Type, enumStr)
			if err != nil {
				return nil, fmt.Errorf("parse enum value %q: %w", enumStr, err)
			}
			compiled.Enum = append(compiled.Enum, *enumVal)
		}
	}

	if debug {
		log.Debug(
			"compiled field",
			infonyz.String("field", compiled.Name),
			infonyz.String("type", compiled.Type.String()),
			infonyz.String("default", compiled.Default.Stringify()),
			infonyz.Bool("required", compiled.Required),
			infonyz.Int("enum_count", len(f.Enum)),
			infonyz.F("enum_values", f.Enum),
		)
	}

	return compiled, nil
}
