package value

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/onyz1/onyzify/internal/types"
)

// Value represents a value of a specific type, encapsulating the actual data based on the defined Type.
type Value struct {
	// Type indicates the type of the value, which determines how the underlying data is stored and accessed.
	Type types.Type

	// Int holds the integer value when Type is TypeInt.
	Int int
	// Int64 holds the int64 value when Type is TypeInt64.
	Int64 int64

	// Uint holds the unsigned integer value when Type is TypeUint.
	Uint uint
	// Uint64 holds the uint64 value when Type is TypeUint64.
	Uint64 uint64

	// String holds the string value when Type is TypeString.
	String string

	// Byte holds the byte value when Type is TypeByte.
	Byte byte

	// Bool holds the boolean value when Type is TypeBool.
	Bool bool

	// Float64 holds the float64 value when Type is TypeFloat64.
	Float64 float64

	// Timestamp holds the time value when Type is TypeTimestamp.
	Timestamp time.Time

	// List holds a slice of Value when Type is a list type (e.g., []int, []string, etc.).
	List []*Value
}

// Stringify converts the Value to its string representation based on its Type.
// It uses the appropriate formatting for each type to ensure that the output
// is human-readable and correctly represents the underlying data.
func (v *Value) Stringify() string {
	switch v.Type.Kind {
	case types.TypeInt:
		return strconv.Itoa(v.Int)

	case types.TypeInt64:
		return strconv.FormatInt(v.Int64, 10)

	case types.TypeUint:
		return strconv.FormatUint(uint64(v.Uint), 10)

	case types.TypeUint64:
		return strconv.FormatUint(uint64(v.Uint64), 10)

	case types.TypeString:
		return v.String

	case types.TypeByte:
		return strconv.QuoteRune(rune(v.Byte))

	case types.TypeBool:
		return strconv.FormatBool(v.Bool)

	case types.TypeFloat64:
		return strconv.FormatFloat(v.Float64, 'f', -1, 64)

	case types.TypeTimestamp:
		return v.Timestamp.Format(time.RFC3339)

	case types.TypeList:
		var sb strings.Builder
		sb.WriteString("[")
		for i, item := range v.List {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(item.Stringify())
		}
		sb.WriteString("]")
		return sb.String()

	default:
		return ""
	}
}

// Interface returns the underlying value of the Value struct as an empty interface (any),
// allowing for dynamic type handling based on the Value's Type.
// It checks the Type of the Value and returns the corresponding field as an interface{}.
func (v *Value) Interface() any {
	switch v.Type.Kind {
	case types.TypeInt:
		return v.Int

	case types.TypeInt64:
		return v.Int64

	case types.TypeUint:
		return v.Uint

	case types.TypeUint64:
		return v.Uint64

	case types.TypeString:
		return v.String

	case types.TypeByte:
		return v.Byte

	case types.TypeBool:
		return v.Bool

	case types.TypeFloat64:
		return v.Float64

	case types.TypeList:
		var list []any
		for _, item := range v.List {
			list = append(list, item.Interface())
		}
		return list

	case types.TypeTimestamp:
		return v.Timestamp

	default:
		return nil
	}
}

// Equal compares the current [Value] instance with another [Value] instance for equality.
// It first checks if the [Value.Type] of both [Value] instances is the same.
// If the Types differ, it returns false immediately.
// If the Types are the same, it compares the corresponding value fields based on the Type and returns true if they are equal, or false otherwise.
func (v *Value) Equal(other Value) bool {
	if v.Type != other.Type {
		return false
	}

	switch v.Type.Kind {
	case types.TypeInt:
		return v.Int == other.Int

	case types.TypeInt64:
		return v.Int64 == other.Int64

	case types.TypeUint:
		return v.Uint == other.Uint

	case types.TypeUint64:
		return v.Uint64 == other.Uint64

	case types.TypeString:
		return v.String == other.String

	case types.TypeByte:
		return v.Byte == other.Byte

	case types.TypeBool:
		return v.Bool == other.Bool

	case types.TypeFloat64:
		return v.Float64 == other.Float64

	case types.TypeTimestamp:
		return v.Timestamp.Equal(other.Timestamp)

	case types.TypeList:
		if len(v.List) != len(other.List) {
			return false
		}
		for i := range v.List {
			if !v.List[i].Equal(*other.List[i]) {
				return false
			}
		}
		return true

	default:
		return false
	}
}

// IsZero checks if the Value is considered a "zero" value based on its Type.
// For example, for [TypeInt], it checks if the Int field is 0; for [TypeString],
// it checks if the String field is empty; and so on.
// If the Type is unrecognized, it defaults to returning true, indicating that the Value is considered zero.
func (v *Value) IsZero() bool {
	switch v.Type.Kind {
	case types.TypeInt:
		return v.Int == 0

	case types.TypeInt64:
		return v.Int64 == 0

	case types.TypeUint:
		return v.Uint == 0

	case types.TypeUint64:
		return v.Uint64 == 0

	case types.TypeString:
		return v.String == ""

	case types.TypeByte:
		return v.Byte == 0

	case types.TypeBool:
		return v.Bool == false

	case types.TypeFloat64:
		return v.Float64 == 0.0

	case types.TypeTimestamp:
		return v.Timestamp.IsZero()

	case types.TypeList:
		if len(v.List) == 0 {
			return true
		}
		for _, item := range v.List {
			if !item.IsZero() {
				return false
			}
		}
		return true

	default:
		return true
	}
}

// UnmarshalJSON implements the [json.Unmarshaler] interface for the [Value] struct.
// It takes a JSON byte slice as input and attempts to unmarshal it into the [Value] struct based on the [Value.Type].
//
// The function first unmarshals the JSON into an empty interface (any) to determine the raw data type,
// and then it uses a switch statement to handle different Type kinds, assigning the appropriate field in the [Value] struct.
// If the JSON data does not match the expected type for the [Value], it returns an error indicating the mismatch.
func (v *Value) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	switch v.Type.Kind {
	case types.TypeInt:
		if num, ok := raw.(float64); ok {
			v.Int = int(num)
			return nil
		}
		return fmt.Errorf("expected int, got %T", raw)

	case types.TypeInt64:
		if num, ok := raw.(float64); ok {
			v.Int64 = int64(num)
			return nil
		}
		return fmt.Errorf("expected int64, got %T", raw)

	case types.TypeUint:
		if num, ok := raw.(float64); ok {
			v.Uint = uint(num)
			return nil
		}
		return fmt.Errorf("expected uint, got %T", raw)

	case types.TypeUint64:
		if num, ok := raw.(float64); ok {
			v.Uint64 = uint64(num)
			return nil
		}
		return fmt.Errorf("expected uint64, got %T", raw)

	case types.TypeString:
		if str, ok := raw.(string); ok {
			v.String = str
			return nil
		}
		return fmt.Errorf("expected string, got %T", raw)

	case types.TypeByte:
		if num, ok := raw.(float64); ok {
			v.Byte = byte(num)
			return nil
		}
		return fmt.Errorf("expected byte, got %T", raw)

	case types.TypeBool:
		if b, ok := raw.(bool); ok {
			v.Bool = b
			return nil
		}
		return fmt.Errorf("expected bool, got %T", raw)

	case types.TypeFloat64:
		if num, ok := raw.(float64); ok {
			v.Float64 = num
			return nil
		}
		return fmt.Errorf("expected float64, got %T", raw)

	case types.TypeTimestamp:
		if str, ok := raw.(string); ok {
			t, err := time.Parse(time.RFC3339, str)
			if err != nil {
				return fmt.Errorf("parse timestamp: %w", err)
			}
			v.Timestamp = t
			return nil
		}
		return fmt.Errorf("expected timestamp string, got %T", raw)

	case types.TypeList:
		if arr, ok := raw.([]any); ok {
			for _, item := range arr {
				elemValue := &Value{Type: *v.Type.Elem}
				itemData, err := json.Marshal(item)
				if err != nil {
					return fmt.Errorf("marshal list item: %w", err)
				}
				if err := elemValue.UnmarshalJSON(itemData); err != nil {
					return fmt.Errorf("unmarshal list item: %w", err)
				}
				v.List = append(v.List, elemValue)
			}
			return nil
		}
		return fmt.Errorf("expected list, got %T", raw)

	default:
		return fmt.Errorf("unsupported type: %s", v.Type.String())
	}
}

// ParseValue takes a [types.Type] and an input string,
// attempts to parse the input according to the specified [types.Type],
// and returns a Value containing the parsed data or an error if parsing fails
// or if the [types.Type] is unsupported.
//
// For list types, it expects the input to be a JSON array string (e.g., "[1, 2, 3]")
// and parses each element according to the element type defined in the [types.Type].
func ParseValue(t types.Type, input string) (*Value, error) {
	var v = new(Value)
	v.Type = t

	switch v.Type.Kind {
	case types.TypeInt:
		parsed, err := strconv.Atoi(input)
		if err != nil {
			return nil, fmt.Errorf("parse int %q: %w", input, err)
		}

		v.Int = parsed

	case types.TypeInt64:
		parsed, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse int64 %q: %w", input, err)
		}

		v.Int64 = parsed

	case types.TypeUint:
		parsed, err := strconv.ParseUint(input, 10, 0)
		if err != nil {
			return nil, fmt.Errorf("parse uint %q: %w", input, err)
		}

		v.Uint = uint(parsed)

	case types.TypeUint64:
		parsed, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse uint64 %q: %w", input, err)
		}

		v.Uint64 = parsed

	case types.TypeString:
		v.String = input

	case types.TypeByte:
		parsed, err := strconv.ParseUint(input, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("parse byte %q: %w", input, err)
		}

		v.Byte = byte(parsed)

	case types.TypeBool:
		parsed, err := strconv.ParseBool(input)
		if err != nil {
			return nil, fmt.Errorf("parse bool %q: %w", input, err)
		}

		v.Bool = parsed

	case types.TypeFloat64:
		parsed, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return nil, fmt.Errorf("parse float64 %q: %w", input, err)
		}

		v.Float64 = parsed

	case types.TypeTimestamp:
		parsed, err := time.Parse(time.RFC3339, input)
		if err != nil {
			return nil, fmt.Errorf("parse timestamp %q: %w", input, err)
		}

		v.Timestamp = parsed

	case types.TypeList:
		if v.Type.Elem == nil {
			return nil, ErrMustHaveElemType
		}

		var raw []json.RawMessage
		if err := json.Unmarshal([]byte(input), &raw); err != nil {
			return nil, fmt.Errorf("unmarshal list: %w", err)
		}

		for _, item := range raw {
			elemValue := &Value{Type: *v.Type.Elem}

			if err := elemValue.UnmarshalJSON(item); err != nil {
				return nil, fmt.Errorf("unmarshal list item: %w", err)
			}

			v.List = append(v.List, elemValue)
		}

	default:
		return nil, ErrUnsupportedType
	}

	return v, nil
}
