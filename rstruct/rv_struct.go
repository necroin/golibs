package rstruct

import (
	"encoding/json"
	"fmt"
)

type RVStruct struct {
	rtStruct     *RTStruct
	fields       []*RVField
	fieldsByName map[string]*RVField
}

func (rvs *RVStruct) FieldByIndex(index int) *RVField {
	if index < 0 || index >= len(rvs.fields) {
		return nil
	}
	return rvs.fields[index]
}

func (rvs *RVStruct) FieldByName(name string) *RVField {
	field, ok := rvs.fieldsByName[name]
	if !ok {
		return nil
	}
	return field
}

func (rvs *RVStruct) Type() *RTStruct {
	return rvs.rtStruct
}

func (rvs *RVStruct) String() string {
	return fmt.Sprintf("%v", rvs.fieldsByName)
}

func (rvs *RVStruct) ToMap(tag string) map[string]any {
	result := map[string]any{}

	for fieldName, field := range rvs.fieldsByName {
		tagValue, ok := field.rtField.GetTag(tag)
		if !ok {
			tagValue = fieldName
		}
		fieldValue := field.value
		rvsValue, ok := fieldValue.(*RVStruct)
		if ok {
			fieldValue = rvsValue.ToMap(tag)
		}
		result[tagValue] = fieldValue
	}
	return result
}

func (rvs *RVStruct) ToJson(tag string) ([]byte, error) {
	jsonFieldsByName := rvs.ToMap(tag)
	return json.Marshal(jsonFieldsByName)
}
