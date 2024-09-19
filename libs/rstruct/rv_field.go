package rstruct

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type RVField struct {
	rtField  *RTField
	value    any
	isStruct bool
}

func (rvf *RVField) Set(value any) {
	rvf.value = value
}

func (rvf *RVField) Get() any {
	return rvf.value
}

func (rvf *RVField) Type() *RTField {
	return rvf.rtField
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

func (rvf *RVField) IsInterface() bool {
	return rvf.Kind() == reflect.Interface
}

func (rvf *RVField) IsStruct() bool {
	return rvf.isStruct
}

func (rvf *RVField) IsSlice() bool {
	return rvf.Kind() == reflect.Slice
}

func (rvf *RVField) IsMap() bool {
	return rvf.Kind() == reflect.Map
}

func (rvf *RVField) AsStruct() *RVStruct {
	return rvf.value.(*RVStruct)
}
