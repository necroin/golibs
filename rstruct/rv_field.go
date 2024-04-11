package rstruct

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type RVField struct {
	rtField *RTField
	value   any
}

func (rvf *RVField) Set(value any) {
	rvf.value = value
}

func (rvf *RVField) Get() any {
	return rvf.value
}

func (rvf *RVField) String() string {
	if rvf.value == nil {
		return ""
	}
	return fmt.Sprintf("%v", rvf.value)
}

func (rvf *RVField) ToJson() ([]byte, error) {
	return json.Marshal(rvf.value)
}

func (rvf *RVField) IsNil() bool {
	return rvf.value == nil
}

func (rvf *RVField) Kind() reflect.Kind {
	if rvf.IsNil() {
		return reflect.Pointer
	}
	return reflect.TypeOf(rvf.value).Kind()
}

func (rvf *RVField) IsPointer() bool {
	return rvf.Kind() == reflect.Pointer
}

func (rvf *RVField) IsStruct() bool {
	return rvf.Kind() == reflect.Struct
}
