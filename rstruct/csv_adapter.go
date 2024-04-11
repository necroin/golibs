package rstruct

import (
	"reflect"

	"github.com/necroin/golibs/csv"
)

type CSVAdapter struct {
	structValue *RVStruct
	fieldValue  *RVField
}

func NewCSVAdapter(structType *RTStruct, value reflect.Value) csv.Adapter {
	castedValue, ok := value.Addr().Interface().(*RVStruct)
	if !ok {
		panic("[rstruct] [CSVAdapter] failed cast reflect interface to RVStruct")
	}

	instance := structType.New()
	castedValue.rtStruct = structType
	castedValue.fields = instance.fields
	castedValue.fieldsByName = instance.fieldsByName
	return &CSVAdapter{
		structValue: castedValue,
		fieldValue:  nil,
	}
}

func (csva *CSVAdapter) Kind() reflect.Kind {
	if csva.fieldValue != nil {
		return csva.fieldValue.Kind()
	}
	return reflect.Struct
}

func (csva *CSVAdapter) IsPointer() bool {
	if csva.fieldValue != nil {
		return csva.fieldValue.IsPointer()
	}
	return false
}

func (csva *CSVAdapter) IsNil() bool {
	if csva.fieldValue != nil {
		return csva.fieldValue.IsNil()
	}
	return csva.structValue == nil
}

func (csva *CSVAdapter) IsStruct() bool {
	if csva.fieldValue != nil {
		return csva.fieldValue.IsStruct()
	}
	return true
}

func (csva *CSVAdapter) Set(value any) {
	if csva.fieldValue != nil {
		csva.fieldValue.value = value
	} else {
		castedValue, ok := value.(*RVStruct)
		if !ok {
			panic("[rstruct] [CSVAdapter] [Set] failed cast reflect interface to RVStruct")
		}
		csva.structValue = castedValue
	}

}

func (csva *CSVAdapter) Get() any {
	if csva.fieldValue != nil {
		return csva.fieldValue.value
	}
	return csva.structValue
}

func (csva *CSVAdapter) New() csv.Adapter {
	if csva.fieldValue != nil {
		csva.fieldValue.value = ""
		return csva
	}
	return &CSVAdapter{
		structValue: csva.structValue.rtStruct.New(),
		fieldValue:  nil,
	}
}

func (csva *CSVAdapter) Deref() csv.Adapter {
	return csva
}

func (csva *CSVAdapter) Field(index int) csv.Adapter {
	field := csva.structValue.fields[index]

	castedFieldValue, ok := field.value.(*RVStruct)
	if ok {
		return &CSVAdapter{
			structValue: castedFieldValue,
			fieldValue:  nil,
		}
	}
	return &CSVAdapter{
		structValue: nil,
		fieldValue:  field,
	}
}

func (csva *CSVAdapter) NumField() int {
	return len(csva.structValue.fields)
}

func (csva *CSVAdapter) GetTag(key string) string {
	return csva.fieldValue.rtField.tags[key]
}

func (csva *CSVAdapter) SetValue(value any) {
	csva.fieldValue.Set(value)
}

func (csva *CSVAdapter) SetString(value string) {
	csva.SetValue(value)
}

func (csva *CSVAdapter) SetInt(value int64) {
	csva.SetValue(value)
}

func (csva *CSVAdapter) SetUint(value uint64) {
	csva.SetValue(value)
}

func (csva *CSVAdapter) SetFloat(value float64) {
	csva.SetValue(value)
}

func (csva *CSVAdapter) SetBool(value bool) {
	csva.SetValue(value)
}
