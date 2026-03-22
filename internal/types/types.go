package types

import "fmt"

// Type represents the basic data types.
type Type uint8

// Predefined types.
const (
	TypeInvalid Type = iota

	TypeInt
	TypeIntList
	TypeInt64
	TypeInt64List

	TypeUint
	TypeUintList
	TypeUint64
	TypeUint64List

	TypeString
	TypeStringList
	TypeByte
	// TypeByteList is a special type that represents a list of bytes,
	// typically used for binary data it corresponds to [[]byte].
	TypeByteList

	TypeBool
	TypeBoolList

	TypeFloat64
	TypeFloat64List
)

// typeFromString maps string representations of types to their corresponding Type values.
var typeFromString = map[string]Type{
	"int":     TypeInt,
	"[]int":   TypeIntList,
	"int64":   TypeInt64,
	"[]int64": TypeInt64List,

	"uint":     TypeUint,
	"[]uint":   TypeUintList,
	"uint64":   TypeUint64,
	"[]uint64": TypeUint64List,

	"string":          TypeString,
	"[]string":        TypeStringList,
	"byte":            TypeByte,
	"[]byte (binary)": TypeByteList,

	"bool":   TypeBool,
	"[]bool": TypeBoolList,

	"float64":   TypeFloat64,
	"[]float64": TypeFloat64List,
}

// stringFromType maps Type values to their corresponding string representations.
var stringFromType = map[Type]string{
	TypeInt:       "int",
	TypeIntList:   "[]int",
	TypeInt64:     "int64",
	TypeInt64List: "[]int64",

	TypeUint:       "uint",
	TypeUintList:   "[]uint",
	TypeUint64:     "uint64",
	TypeUint64List: "[]uint64",

	TypeString:     "string",
	TypeStringList: "[]string",
	TypeByte:       "byte",
	TypeByteList:   "[]byte (binary)",

	TypeBool:     "bool",
	TypeBoolList: "[]bool",

	TypeFloat64:     "float64",
	TypeFloat64List: "[]float64",
}

// ParseType converts a string representation of a type to its corresponding Type value.
func ParseType(typeStr string) (Type, error) {
	if t, ok := typeFromString[typeStr]; ok {
		return t, nil
	}
	return TypeInvalid, fmt.Errorf("unknown type: %q", typeStr)
}

// IsListType checks if the given Type is a list type (e.g., []int, []string, etc.).
func IsListType(t Type) bool {
	return t == TypeIntList ||
		t == TypeInt64List ||
		t == TypeUintList ||
		t == TypeUint64List ||
		t == TypeStringList ||
		t == TypeByteList ||
		t == TypeBoolList ||
		t == TypeFloat64List
}

// String returns the string representation of the Type.
func (t Type) String() string {
	if s, ok := stringFromType[t]; ok {
		return s
	}
	return "invalid"
}
