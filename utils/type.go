package utils

import (
	"fmt"
	"reflect"
)

// Reports whether value is pointer.
func IsPointer(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Pointer
}

// Reports whether value is interface.
func IsInterface(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Interface
}

// Reports whether value is struct.
func IsStruct(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Struct || (IsPointer(value) && value.Elem().Kind() == reflect.Struct)
}

// Reports whether value is slice.
func IsSlice(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Slice || (IsPointer(value) && value.Type().Elem().Kind() == reflect.Slice)
}

// Reports whether value is map.
func IsMap(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Map || (IsPointer(value) && value.Type().Elem().Kind() == reflect.Map)
}

// Reports whether value is nil.
func IsNil(value reflect.Value) bool {
	return value.Interface() == nil || (IsPointer(value) || IsMap(value) || IsSlice(value) || IsInterface(value)) && value.IsNil()
}

// Dereferences the value if it is a pointer.
func DerefValue(value reflect.Value) reflect.Value {
	result := value
	if IsPointer(result) {
		result = result.Elem()
	}
	return result
}

// Dereferences the value if it is a pointer.
func DerefType(value reflect.Type) reflect.Type {
	result := value
	if result.Kind() == reflect.Pointer {
		result = result.Elem()
	}
	return result
}

// Dereferences the value if it is a pointer.
func DerefValueOf(value any) reflect.Value {
	return DerefValue(reflect.ValueOf(value))
}

// Dereferences the value if it is a pointer.
func DerefTypeOf(value any) reflect.Type {
	return DerefType(reflect.TypeOf(value))
}

// Returns v's length.
func Len(value reflect.Value) int {
	if IsSlice(value) || IsMap(value) {
		return value.Len()
	}
	return 0
}

// Returns v's length.
func LenOf(value any) int {
	if value == nil {
		return 0
	}
	return Len(DerefValueOf(value))
}

// Returns the boolean value represented by the string.
func ParseBool(value string) (bool, error) {
	switch value {
	case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "y", "ON", "on", "On":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "n", "OFF", "off", "Off":
		return false, nil
	}
	return false, fmt.Errorf("parsing \"%s\": invalid syntax", value)
}
