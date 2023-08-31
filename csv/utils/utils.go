package utils

import (
	"reflect"
)

func IsStruct(rValue reflect.Value) bool {
	return rValue.Type().Kind() == reflect.Struct || (rValue.Type().Kind() == reflect.Pointer && rValue.Type().Elem().Kind() == reflect.Struct)
}

func PointerOf[T any](value T) *T {
	return &value
}

func InstantiateSliceElement[T any](value *[]T) *T {
	return new(T)
}
