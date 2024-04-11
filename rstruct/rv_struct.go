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

func (rvs *RVStruct) String() string {
	return fmt.Sprintf("%v", rvs.fieldsByName)
}
func (rvs *RVStruct) ToJson(tag string) ([]byte, error) {
	jsonFieldsByName := map[string]any{}

	for fieldName, field := range rvs.fieldsByName {
		tagValue, ok := field.rtField.GetTag(tag)
		if !ok {
			tagValue = fieldName
		}
		jsonValue := field.value
		castedFieldValue, ok := field.value.(*RVStruct)
		if ok {
			jsonValueData, _ := castedFieldValue.ToJson(tag)
			jsonValue = string(jsonValueData)
		}
		jsonFieldsByName[tagValue] = jsonValue
	}

	return json.Marshal(jsonFieldsByName)
}

func (rvs *RVStruct) FieldByIndex(index int) *RVField {
	return rvs.fields[index]
}

func (rvs *RVStruct) FieldByName(name string) *RVField {
	return rvs.fieldsByName[name]
}
