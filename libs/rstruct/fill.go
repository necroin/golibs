package rstruct

import (
	"fmt"
	"reflect"

	"github.com/necroin/golibs/utils"
)

type TypeSetter func(value string, dst reflect.Value) error

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

func (message *RVStruct) FillStruct(setters map[string]TypeSetter, fillData any, withClear bool) error {
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
			}
			if err := messageField.AsStruct().FillStruct(setters, rvField.Interface(), withClear); err != nil {
				return err
			}
			continue
		}

		srcValue := messageField.String()
		if messageField.IsPointer() {
			srcValue = fmt.Sprintf("%v", utils.DerefValueOf(messageField.Get()).Interface())
		}

		if srcValue == "*" {
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

		if err := setByType(setters, srcValue, rvField, typeName); err != nil {
			return fmt.Errorf("[FillStruct] failed set value for %s field: %s", rtField.Name, err)
		}

		if withClear {
			messageField.Set(nil)
		}
	}

	return nil
}
