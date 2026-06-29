package rstruct

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/necroin/golibs/utils"
)

type FillOption func(options *FillOptions)

type FillOptions struct {
	WithClear  bool
	SkipValues []string
}

func FillWithClear() FillOption {
	return func(options *FillOptions) {
		options.WithClear = true
	}
}

func WithSkipValues(values ...string) FillOption {
	return func(options *FillOptions) {
		options.SkipValues = append(options.SkipValues, values...)
	}
}

func setByType(setters map[string]TypeSetter, src string, dst reflect.Value, typeName string) error {
	if src == "null" {
		return nil
	}

	setter, ok := setters[typeName]
	if !ok {
		return fmt.Errorf("[SetByType] unknown type: %s", typeName)
	}
	if err := setter(src, dst); err != nil {
		return fmt.Errorf("[SetByType] -> %s", err)
	}
	return nil
}

func validateType(setters map[string]TypeSetter, typeName string) error {
	_, ok := setters[typeName]
	if !ok {
		return fmt.Errorf("[ValidateType] unknown type: %s", typeName)
	}
	return nil
}

func (message *RVStruct) FillStruct(setters map[string]TypeSetter, fillData any, opts ...FillOption) error {
	options := &FillOptions{}
	for _, opt := range opts {
		opt(options)
	}

	rvFillData := utils.DerefValueOf(fillData)
	rtFillData := utils.DerefTypeOf(fillData)

	for i := range rvFillData.NumField() {
		rvField := rvFillData.Field(i)
		rtField := rtFillData.Field(i)

		if !rtField.IsExported() {
			continue
		}

		messageField := message.FieldByName(rtField.Name)
		if messageField == nil {
			return fmt.Errorf("[FillStruct] missing field: %s", rtField.Name)
		}

		if messageField.IsNil() {
			continue
		}

		if utils.IsStruct(rvField) && messageField.IsStruct() {
			if messageField.AsStruct().IsNil() {
				continue
			}

			if utils.IsPointer(rvField) {
				rvField.Set(reflect.New(rvField.Type().Elem()))

				if err := messageField.AsStruct().FillStruct(setters, rvField.Interface(), opts...); err != nil {
					return err
				}
			} else {
				if err := messageField.AsStruct().FillStruct(setters, rvField.Addr().Interface(), opts...); err != nil {
					return err
				}
			}

			continue
		}

		srcValue := messageField.String()
		if messageField.IsPointer() {
			srcValue = fmt.Sprintf("%v", utils.DerefValueOf(messageField.Get()).Interface())
		}

		if slices.Contains(options.SkipValues, srcValue) {
			continue
		}

		typeName := utils.GetFullNameOfTypeReflect(rtField.Type)
		if typeName == "" {
			typeName = utils.GetFullNameOfTypeReflect(rtField.Type.Elem())
		}

		if utils.IsMap(rvField) {
			typeName = fmt.Sprintf("map[%s]%s", rtField.Type.Key().Name(), rtField.Type.Elem().Name())
		}

		if utils.IsSlice(rvField) {
			sliceType := rtField.Type.Elem()
			if sliceType.Kind() == reflect.Pointer {
				sliceType = sliceType.Elem()
			}
			typeName = fmt.Sprintf("[]%s", utils.GetFullNameOfTypeReflect(sliceType))

		}

		if err := validateType(setters, typeName); err != nil && !(utils.IsStruct(rvField) && messageField.IsStruct()) {
			return fmt.Errorf("[FillStruct] failed validate type: %s", err)
		}

		if err := setByType(setters, srcValue, rvField, typeName); err != nil {
			return fmt.Errorf("[FillStruct] failed set value for %s field: %s", rtField.Name, err)
		}

		if options.WithClear {
			messageField.Set(nil)
		}
	}

	return nil
}
