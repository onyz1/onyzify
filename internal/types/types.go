package types

import (
	"fmt"
	"strings"
)

// TypeKind represents the specific kind of a type, such as int, string, bool, etc.
// It is used to identify the underlying data type of a Value.
type TypeKind uint8

// Constants for the different TypeKind values, representing various data types.
const (
	TypeInvalid TypeKind = iota

	TypeInt
	TypeInt64

	TypeUint
	TypeUint64

	TypeString
	TypeByte

	TypeBool

	TypeFloat64

	TypeList

	TypeTimestamp
)

// Type represents the type of a value. For list types, Elem holds the element type.
type Type struct {
	Kind TypeKind
	Elem *Type // For list types, Elem holds the type of the list elements.
}

var typeKindToString = map[TypeKind]string{
	TypeInvalid: "invalid",

	TypeInt:   "int",
	TypeInt64: "int64",

	TypeUint:   "uint",
	TypeUint64: "uint64",

	TypeString: "string",
	TypeByte:   "byte",

	TypeBool: "bool",

	TypeFloat64: "float64",

	TypeList: "list",

	TypeTimestamp: "timestamp",
}

// String returns a string representation of the Type. For basic types, it returns the
// type name (e.g., "int", "string"). For list types, it returns a string in the format
// "[]<elemType>", where <elemType> is the string representation of the element type.
// If the type kind is unknown, it returns "unknown".
func (t *Type) String() string {
	var sb strings.Builder

	switch t.Kind {
	case TypeList:
		sb.WriteString(typeKindToString[TypeList])
		sb.WriteString("[")
		if t.Elem != nil {
			sb.WriteString(t.Elem.String())
		} else {
			sb.WriteString("unknown")
		}
		sb.WriteString("]")

	default:
		if s, ok := typeKindToString[t.Kind]; ok {
			sb.WriteString(s)
		} else {
			sb.WriteString("unknown")
		}
	}

	return sb.String()
}

var stringToTypeKind = map[string]TypeKind{
	"invalid": TypeInvalid,

	"int":   TypeInt,
	"int64": TypeInt64,

	"uint":   TypeUint,
	"uint64": TypeUint64,

	"string": TypeString,
	"byte":   TypeByte,

	"bool": TypeBool,

	"float64": TypeFloat64,

	"list": TypeList,

	"timestamp": TypeTimestamp,
}

// ParseType takes a string representation of a type and converts it into a Type struct.
// It supports basic types as well as list types (e.g., "[]int", "[]string", etc.) and nested types (e.g., "[][]int").
// If the input string does not correspond to a known type, it returns an error.
func ParseType(s string) (*Type, error) {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, typeKindToString[TypeList]+"[") && strings.HasSuffix(s, "]") {
		elemTypeStr := strings.TrimSuffix(strings.TrimPrefix(s, typeKindToString[TypeList]+"["), "]")
		elemType, err := ParseType(elemTypeStr)
		if err != nil {
			return nil, fmt.Errorf("type: %q: %w", s, err)
		}
		return &Type{Kind: TypeList, Elem: elemType}, nil
	}

	if tk, ok := stringToTypeKind[s]; ok {
		return &Type{Kind: tk}, nil
	}

	return nil, fmt.Errorf("type: %q: %w", s, ErrUnknownType)
}
