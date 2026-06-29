package rstruct

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/necroin/golibs/utils"
)

type TypeSetter func(value string, dst reflect.Value) error
type Getter[T any] func(value *string) (*T, error)
type Converter[T any] func(value *string) (*T, error)

var (
	DefaultTypeSetters = map[string]TypeSetter{
		"bool":    SetBool,
		"string":  SetString,
		"float32": SetFloat32,
		"float64": SetFloat64,
		"int":     SetInt,
		"int32":   SetInt32,
		"int64":   SetInt64,
		"uint":    SetUInt,
		"uint32":  SetUInt32,
		"uint64":  SetUInt64,

		"[]bool":    SetHandler(MakeMultyGetter(MakeJsonSliceGetter[bool]("Get_Bool_Slice"), MakeStringSliceGetter(utils.CastToBool))),
		"[]string":  SetHandler(MakeMultyGetter(MakeJsonSliceGetter[string]("Get_String_Slice"), MakeStringSliceGetter(func(value *string) (*string, error) { return value, nil }))),
		"[]float32": SetHandler(MakeMultyGetter(MakeJsonSliceGetter[float32]("Get_Float32_Slice"), MakeStringSliceGetter(utils.CastToFloat32))),
		"[]float64": SetHandler(MakeMultyGetter(MakeJsonSliceGetter[float64]("Get_Float64_Slice"), MakeStringSliceGetter(utils.CastToFloat64))),
		"[]int":     SetHandler(MakeMultyGetter(MakeJsonSliceGetter[int]("Get_Int_Slice"), MakeStringSliceGetter(utils.CastToInt))),
		"[]int32":   SetHandler(MakeMultyGetter(MakeJsonSliceGetter[int32]("Get_Int32_Slice"), MakeStringSliceGetter(utils.CastToInt32))),
		"[]int64":   SetHandler(MakeMultyGetter(MakeJsonSliceGetter[int64]("Get_Int64_Slice"), MakeStringSliceGetter(utils.CastToInt64))),
		"[]uint":    SetHandler(MakeMultyGetter(MakeJsonSliceGetter[uint]("Get_Uint_Slice"), MakeStringSliceGetter(utils.CastToUInt))),
		"[]uint32":  SetHandler(MakeMultyGetter(MakeJsonSliceGetter[uint32]("Get_Uint32_Slice"), MakeStringSliceGetter(utils.CastToUInt32))),
		"[]uint64":  SetHandler(MakeMultyGetter(MakeJsonSliceGetter[uint64]("Get_Uint64_Slice"), MakeStringSliceGetter(utils.CastToUInt64))),

		"map[string]bool":    SetHandler(MakeJsonGetter[map[string]bool]("Get_map[string]bool")),
		"map[string]string":  SetHandler(MakeJsonGetter[map[string]string]("Get_map[string]string")),
		"map[string]float32": SetHandler(MakeJsonGetter[map[string]float32]("Get_map[string]float32")),
		"map[string]float64": SetHandler(MakeJsonGetter[map[string]float64]("Get_map[string]float64")),
		"map[string]int":     SetHandler(MakeJsonGetter[map[string]int]("Get_map[string]int")),
		"map[string]int32":   SetHandler(MakeJsonGetter[map[string]int32]("Get_map[string]int32")),
		"map[string]int64":   SetHandler(MakeJsonGetter[map[string]int64]("Get_map[string]int64")),
		"map[string]uint":    SetHandler(MakeJsonGetter[map[string]uint]("Get_map[string]uint")),
		"map[string]uint32":  SetHandler(MakeJsonGetter[map[string]uint32]("Get_map[string]uint32")),
		"map[string]uint64":  SetHandler(MakeJsonGetter[map[string]uint64]("Get_map[string]uint64")),

		"map[int]bool":    SetHandler(MakeJsonGetter[map[int]bool]("Get_map[int]bool")),
		"map[int]string":  SetHandler(MakeJsonGetter[map[int]string]("Get_map[int]string")),
		"map[int]float32": SetHandler(MakeJsonGetter[map[int]float32]("Get_map[int]float32")),
		"map[int]float64": SetHandler(MakeJsonGetter[map[int]float64]("Get_map[int]float64")),
		"map[int]int":     SetHandler(MakeJsonGetter[map[int]int]("Get_map[int]int")),
		"map[int]int32":   SetHandler(MakeJsonGetter[map[int]int32]("Get_map[int]int32")),
		"map[int]int64":   SetHandler(MakeJsonGetter[map[int]int64]("Get_map[int]int64")),
		"map[int]uint":    SetHandler(MakeJsonGetter[map[int]uint]("Get_map[int]uint")),
		"map[int]uint32":  SetHandler(MakeJsonGetter[map[int]uint32]("Get_map[int]uint32")),
		"map[int]uint64":  SetHandler(MakeJsonGetter[map[int]uint64]("Get_map[int]uint64")),

		"map[int32]bool":    SetHandler(MakeJsonGetter[map[int32]bool]("Get_map[int]bool")),
		"map[int32]string":  SetHandler(MakeJsonGetter[map[int32]string]("Get_map[int]string")),
		"map[int32]float32": SetHandler(MakeJsonGetter[map[int32]float32]("Get_map[int]float32")),
		"map[int32]float64": SetHandler(MakeJsonGetter[map[int32]float64]("Get_map[int]float64")),
		"map[int32]int":     SetHandler(MakeJsonGetter[map[int32]int]("Get_map[int]int")),
		"map[int32]int32":   SetHandler(MakeJsonGetter[map[int32]int32]("Get_map[int]int32")),
		"map[int32]int64":   SetHandler(MakeJsonGetter[map[int32]int64]("Get_map[int]int64")),
		"map[int32]uint":    SetHandler(MakeJsonGetter[map[int32]uint]("Get_map[int]uint")),
		"map[int32]uint32":  SetHandler(MakeJsonGetter[map[int32]uint32]("Get_map[int]uint32")),
		"map[int32]uint64":  SetHandler(MakeJsonGetter[map[int32]uint64]("Get_map[int]uint64")),
	}
)

func SetTemplate[T any](value string, dst reflect.Value, castFunc func(value *string) (*T, error)) error {
	castedPointerValue, err := castFunc(&value)
	if err != nil {
		return err
	}

	if utils.IsPointer(dst) {
		dst.Set(reflect.ValueOf(castedPointerValue))
		return nil
	}
	dst.Set(reflect.ValueOf(utils.GetPointerValueOrDefault(castedPointerValue)))
	return nil
}

func SetHandler[T any](handler func(value *string) (*T, error)) TypeSetter {
	return func(value string, dst reflect.Value) error {
		return SetTemplate(value, dst, handler)
	}
}

func SetString(value string, dst reflect.Value) error {
	return SetTemplate(value, dst, func(value *string) (*string, error) { return value, nil })
}

func SetInt(value string, dst reflect.Value) error {
	return SetTemplate(value, dst, utils.CastToInt)
}

func SetInt32(value string, dst reflect.Value) error {
	return SetTemplate(value, dst, utils.CastToInt32)
}

func SetInt64(value string, dst reflect.Value) error {
	return SetTemplate(value, dst, utils.CastToInt64)
}

func SetUInt(value string, dst reflect.Value) error {
	return SetTemplate(value, dst, utils.CastToUInt)
}

func SetUInt32(value string, dst reflect.Value) error {
	return SetTemplate(value, dst, utils.CastToUInt32)
}

func SetUInt64(value string, dst reflect.Value) error {
	return SetTemplate(value, dst, utils.CastToUInt64)
}

func SetFloat32(value string, dst reflect.Value) error {
	return SetTemplate(value, dst, utils.CastToFloat32)
}

func SetFloat64(value string, dst reflect.Value) error {
	return SetTemplate(value, dst, utils.CastToFloat64)
}

func SetBool(value string, dst reflect.Value) error {
	return SetTemplate(value, dst, utils.CastToBool)
}
func MakeStringSliceGetter[T any](converter Converter[T]) Getter[[]T] {
	return func(value *string) (*[]T, error) {
		if value == nil {
			return nil, nil
		}

		result := []T{}

		parts := utils.ParseStringSlice(*value)
		for _, part := range parts {
			resultPart, err := converter(&part)
			if err != nil {
				return nil, fmt.Errorf("[GetStringSlice] failed convert %s value", part)
			}
			result = append(result, *resultPart)
		}

		return &result, nil
	}
}

func GetJsonTemplate[T any](value *string, tag string) (*T, error) {
	if value == nil {
		return nil, nil
	}
	result := new(T)
	if err := json.Unmarshal([]byte(*value), result); err != nil {
		return result, fmt.Errorf("[%s] json decode error: %s", tag, err)
	}
	return result, nil
}

func MakeJsonGetter[T any](tag string) Getter[T] {
	return func(value *string) (*T, error) {
		return GetJsonTemplate[T](value, tag)
	}
}

func GetJsonSliceTemplate[T any](value *string, tag string) (*[]T, error) {
	if value == nil {
		return nil, nil
	}
	result := []T{}
	if err := json.Unmarshal([]byte(*value), &result); err != nil {
		return &result, fmt.Errorf("[%s] json decode error for %s: %s", tag, *value, err)
	}
	return &result, nil
}

func MakeJsonSliceGetter[T any](tag string) Getter[[]T] {
	return func(value *string) (*[]T, error) {
		return GetJsonSliceTemplate[T](value, tag)
	}
}

func MakeMultyGetter[T any](getters ...Getter[T]) Getter[T] {
	return func(value *string) (*T, error) {
		errors := []error{}
		for _, getter := range getters {
			result, err := getter(value)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			return result, nil
		}

		return nil, fmt.Errorf("[Multy Getter] failed get value: %s", errors)
	}
}
