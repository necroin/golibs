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
		jsonFieldsByName[tagValue] = field.value
	}

	return json.Marshal(jsonFieldsByName)
}

func (rvs *RVStruct) SetByIndex(index int, value any) {
	rvs.fields[index].value = value
}

func (rvs *RVStruct) SetByName(name string, value any) {
	rvs.fieldsByName[name].value = value
}

func (rvs *RVStruct) GetByIndex(index int) *RVField {
	return rvs.fields[index]
}

func (rvs *RVStruct) GetByName(name string) *RVField {
	return rvs.fieldsByName[name]
}
