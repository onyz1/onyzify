package value

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/onyz1/onyzify/internal/types"
)

// Value represents a value of a specific type, encapsulating the actual data based on the defined Type.
type Value struct {
	// Type indicates the type of the value, which determines how the underlying data is stored and accessed.
	Type types.Type

	// Int holds the integer value when Type is TypeInt.
	Int int
	// IntList holds a slice of integers when Type is TypeIntList.
	IntList []int
	// Int64 holds the int64 value when Type is TypeInt64.
	Int64 int64
	// Int64List holds a slice of int64 values when Type is TypeInt64List.
	Int64List []int64

	// Uint holds the unsigned integer value when Type is TypeUint.
	Uint uint
	// UintList holds a slice of unsigned integers when Type is TypeUintList.
	UintList []uint
	// Uint64 holds the uint64 value when Type is TypeUint64.
	Uint64 uint64
	// Uint64List holds a slice of uint64 values when Type is TypeUint64List.
	Uint64List []uint64

	// String holds the string value when Type is TypeString.
	String string
	// StringList holds a slice of strings when Type is TypeStringList.
	StringList []string

	// Byte holds the byte value when Type is TypeByte.
	Byte byte
	// ByteList holds a slice of bytes when Type is TypeByteList, typically used for binary data.
	ByteList []byte

	// Bool holds the boolean value when Type is TypeBool.
	Bool bool
	// BoolList holds a slice of boolean values when Type is TypeBoolList.
	BoolList []bool

	// Float64 holds the float64 value when Type is TypeFloat64.
	Float64 float64
	// Float64List holds a slice of float64 values when Type is TypeFloat64List.
	Float64List []float64
}

// Stringify converts the Value to its string representation based on its Type.
// It uses the appropriate formatting for each type to ensure that the output
// is human-readable and correctly represents the underlying data.
//
// TODO: Consider implementing a more robust stringification method that
// handles lists more properly by converting each element to a [Value]
// and calling Stringify on each element, then joining all.
func (v *Value) Stringify() string {
	switch v.Type {
	case types.TypeInt:
		return strconv.Itoa(v.Int)

	case types.TypeIntList:
		return fmt.Sprint(v.IntList)

	case types.TypeInt64:
		return strconv.FormatInt(v.Int64, 10)

	case types.TypeInt64List:
		return fmt.Sprint(v.Int64List)

	case types.TypeUint:
		return strconv.FormatUint(uint64(v.Uint), 10)

	case types.TypeUintList:
		return fmt.Sprint(v.UintList)

	case types.TypeUint64:
		return strconv.FormatUint(uint64(v.Uint64), 10)

	case types.TypeUint64List:
		return fmt.Sprint(v.Uint64List)

	case types.TypeString:
		return v.String

	case types.TypeStringList:
		return strings.Join(v.StringList, ", ")

	case types.TypeByte:
		return strconv.QuoteRune(rune(v.Byte))

	case types.TypeByteList:
		return fmt.Sprint(v.ByteList)

	case types.TypeBool:
		return strconv.FormatBool(v.Bool)

	case types.TypeBoolList:
		return fmt.Sprint(v.BoolList)

	case types.TypeFloat64:
		return strconv.FormatFloat(v.Float64, 'f', -1, 64)

	default:
		return ""
	}
}

// Interface returns the underlying value of the Value struct as an empty interface (any),
// allowing for dynamic type handling based on the Value's Type.
// It checks the Type of the Value and returns the corresponding field as an interface{}.
func (v *Value) Interface() any {
	switch v.Type {
	case types.TypeInt:
		return v.Int

	case types.TypeIntList:
		return v.IntList

	case types.TypeInt64:
		return v.Int64

	case types.TypeInt64List:
		return v.Int64List

	case types.TypeUint:
		return v.Uint

	case types.TypeUintList:
		return v.UintList

	case types.TypeUint64:
		return v.Uint64

	case types.TypeUint64List:
		return v.Uint64List

	case types.TypeString:
		return v.String

	case types.TypeStringList:
		return v.StringList

	case types.TypeByte:
		return v.Byte

	case types.TypeByteList:
		return v.ByteList

	case types.TypeBool:
		return v.Bool

	case types.TypeBoolList:
		return v.BoolList

	case types.TypeFloat64:
		return v.Float64

	case types.TypeFloat64List:
		return v.Float64List

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

	switch v.Type {
	case types.TypeInt:
		return v.Int == other.Int

	case types.TypeIntList:
		if len(v.IntList) != len(other.IntList) {
			return false
		}
		for i := range v.IntList {
			if v.IntList[i] != other.IntList[i] {
				return false
			}
		}
		return true

	case types.TypeInt64:
		return v.Int64 == other.Int64

	case types.TypeInt64List:
		if len(v.Int64List) != len(other.Int64List) {
			return false
		}
		for i := range v.Int64List {
			if v.Int64List[i] != other.Int64List[i] {
				return false
			}
		}
		return true

	case types.TypeUint:
		return v.Uint == other.Uint

	case types.TypeUintList:
		if len(v.UintList) != len(other.UintList) {
			return false
		}
		for i := range v.UintList {
			if v.UintList[i] != other.UintList[i] {
				return false
			}
		}
		return true

	case types.TypeUint64:
		return v.Uint64 == other.Uint64

	case types.TypeUint64List:
		if len(v.Uint64List) != len(other.Uint64List) {
			return false
		}
		for i := range v.Uint64List {
			if v.Uint64List[i] != other.Uint64List[i] {
				return false
			}
		}
		return true

	case types.TypeString:
		return v.String == other.String

	case types.TypeStringList:
		if len(v.StringList) != len(other.StringList) {
			return false
		}
		for i := range v.StringList {
			if v.StringList[i] != other.StringList[i] {
				return false
			}
		}
		return true

	case types.TypeByte:
		return v.Byte == other.Byte

	case types.TypeByteList:
		if len(v.ByteList) != len(other.ByteList) {
			return false
		}
		for i := range v.ByteList {
			if v.ByteList[i] != other.ByteList[i] {
				return false
			}
		}
		return true

	case types.TypeBool:
		return v.Bool == other.Bool

	case types.TypeBoolList:
		if len(v.BoolList) != len(other.BoolList) {
			return false
		}
		for i := range v.BoolList {
			if v.BoolList[i] != other.BoolList[i] {
				return false
			}
		}
		return true

	case types.TypeFloat64:
		return v.Float64 == other.Float64

	case types.TypeFloat64List:
		if len(v.Float64List) != len(other.Float64List) {
			return false
		}
		for i := range v.Float64List {
			if v.Float64List[i] != other.Float64List[i] {
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
	switch v.Type {
	case types.TypeInt:
		return v.Int == 0

	case types.TypeIntList:
		return len(v.IntList) == 0

	case types.TypeInt64:
		return v.Int64 == 0

	case types.TypeInt64List:
		return len(v.Int64List) == 0

	case types.TypeUint:
		return v.Uint == 0

	case types.TypeUintList:
		return len(v.UintList) == 0

	case types.TypeUint64:
		return v.Uint64 == 0

	case types.TypeUint64List:
		return len(v.Uint64List) == 0

	case types.TypeString:
		return v.String == ""

	case types.TypeStringList:
		return len(v.StringList) == 0

	case types.TypeByte:
		return v.Byte == 0

	case types.TypeByteList:
		return len(v.ByteList) == 0

	case types.TypeBool:
		return v.Bool == false

	case types.TypeBoolList:
		return len(v.BoolList) == 0

	case types.TypeFloat64:
		return v.Float64 == 0.0

	case types.TypeFloat64List:
		return len(v.Float64List) == 0

	default:
		return true
	}
}

// ParseValue takes a [types.Type] and an input string,
// attempts to parse the input according to the specified [types.Type],
// and returns a Value containing the parsed data or an error if parsing fails or if the [types.Type] is unsupported.
func ParseValue(t types.Type, input string) (*Value, error) {
	var v = new(Value)
	v.Type = t

	if types.IsListType(t) {
		return nil, fmt.Errorf("parse list type %q: see ParseValueList for parsing list types: %w", v.Type, ErrUnsupportedType)
	}

	switch v.Type {
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

	default:
		return nil, ErrUnsupportedType
	}

	return v, nil
}

// ParseValueList takes a [types.Type] that represents a list type
// (e.g., [TypeIntList], [TypeStringList], etc.) and a slice of input strings,
// attempts to parse each input string according to the specified list type,
// and returns a Value containing the parsed list data or an error if parsing fails or if the [types.Type] is unsupported.
func ParseValueList(t types.Type, input []string) (*Value, error) {
	var v = new(Value)
	v.Type = t

	switch v.Type {
	case types.TypeIntList:
		for _, s := range input {
			parsed, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("parse int %q: %w", s, err)
			}
			v.IntList = append(v.IntList, parsed)
		}

	case types.TypeInt64List:
		for _, s := range input {
			parsed, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse int64 %q: %w", s, err)
			}
			v.Int64List = append(v.Int64List, parsed)
		}

	case types.TypeUintList:
		for _, s := range input {
			parsed, err := strconv.ParseUint(s, 10, 0)
			if err != nil {
				return nil, fmt.Errorf("parse uint %q: %w", s, err)
			}
			v.UintList = append(v.UintList, uint(parsed))
		}

	case types.TypeUint64List:
		for _, s := range input {
			parsed, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse uint64 %q: %w", s, err)
			}
			v.Uint64List = append(v.Uint64List, parsed)
		}

	case types.TypeStringList:
		v.StringList = input

	case types.TypeByteList:
		for _, s := range input {
			parsed, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("parse byte %q: %w", s, err)
			}
			v.ByteList = append(v.ByteList, byte(parsed))
		}

	case types.TypeBoolList:
		for _, s := range input {
			parsed, err := strconv.ParseBool(s)
			if err != nil {
				return nil, fmt.Errorf("parse bool %q: %w", s, err)
			}
			v.BoolList = append(v.BoolList, parsed)
		}

	case types.TypeFloat64List:
		for _, s := range input {
			parsed, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, fmt.Errorf("parse float64 %q: %w", s, err)
			}
			v.Float64List = append(v.Float64List, parsed)
		}

	default:
		return nil, ErrUnsupportedType
	}

	return v, nil
}
