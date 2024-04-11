package utils

import (
	"reflect"
)

func IsPointer(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Pointer
}

func IsInterface(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Interface
}

func IsStruct(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Struct || (IsPointer(value) && value.Elem().Kind() == reflect.Struct)
}

func IsSlice(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Slice || (IsPointer(value) && value.Type().Elem().Kind() == reflect.Slice)
}

func IsMap(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Map || (IsPointer(value) && value.Type().Elem().Kind() == reflect.Map)
}

func IsNil(value reflect.Value) bool {
	return (IsPointer(value) || IsMap(value) || IsSlice(value) || IsInterface(value)) && value.IsNil()
}

func DerefValueOf(value interface{}) reflect.Value {
	result := reflect.ValueOf(value)
	if IsPointer(result) {
		result = result.Elem()
	}
	return result
}

func DerefTypeOf(value interface{}) reflect.Type {
	result := reflect.TypeOf(value)
	if result.Kind() == reflect.Pointer {
		result = result.Elem()
	}
	return result
}
