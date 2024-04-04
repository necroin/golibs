package rstruct

import (
	"reflect"

	"github.com/necroin/golibs/csv"
)

type CSVAdapter struct {
	structValue *RVStruct
	structType  *RTStruct

	fieldValue *RVField
	fieldType  *RTField
}

func NewCSVAdapter(structType *RTStruct, value reflect.Value) csv.Adapter {
	castedValue, ok := value.Addr().Interface().(*RVStruct)
	if !ok {
		panic("[rstruct] [CSVAdapter] failed cast reflect interface to RVStruct")
	}

	instance := structType.New()
	castedValue.fields = instance.fields
	castedValue.fieldsByName = instance.fieldsByName
	return &CSVAdapter{
		structValue: castedValue,
		structType:  structType,
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
	}
	return csva
}

func (csva *CSVAdapter) Deref() csv.Adapter {
	return csva
}

func (csva *CSVAdapter) Field(index int) csv.Adapter {
	return &CSVAdapter{
		structValue: nil,
		structType:  nil,
		fieldValue:  csva.structValue.fields[index],
		fieldType:   csva.structType.fields[index],
	}
}

func (csva *CSVAdapter) NumField() int {
	return len(csva.structValue.fields)
}

func (csva *CSVAdapter) GetTag(key string) string {
	return csva.fieldType.tags[key]
}

func (csva *CSVAdapter) SetString(value string) {
	csva.fieldValue.Set(value)
}

func (csva *CSVAdapter) SetInt(value int64) {
	csva.fieldValue.Set(value)
}

func (csva *CSVAdapter) SetUint(value uint64) {
	csva.fieldValue.Set(value)
}

func (csva *CSVAdapter) SetFloat(value float64) {
	csva.fieldValue.Set(value)
}

func (csva *CSVAdapter) SetBool(value bool) {
	csva.fieldValue.Set(value)
}
