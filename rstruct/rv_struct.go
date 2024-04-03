package rstruct

import (
	"encoding/json"
	"fmt"
)

type RVStruct struct {
	fields       []*RVField
	fieldsByName map[string]*RVField
}

func (rvs *RVStruct) String() string {
	return fmt.Sprintf("%v", rvs.fieldsByName)
}
func (rvs *RVStruct) ToJSON() ([]byte, error) {
	return json.Marshal(rvs.fieldsByName)
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
