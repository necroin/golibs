package rstruct

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type RVField struct {
	rtField *RTField
	value   any
	kind    reflect.Kind
}

func (rvf *RVField) Set(value any) {
	if value != nil {
		rvf.kind = reflect.TypeOf(value).Kind()
	}
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

func (rvf *RVField) ToJSON() ([]byte, error) {
	return json.Marshal(rvf.value)
}

func (rvf *RVField) MarshalText() ([]byte, error) {
	return rvf.ToJSON()
}

func (rvf *RVField) IsNil() bool {
	return rvf.value == nil
}

func (rvf *RVField) Kind() reflect.Kind {
	if rvf.kind != 0 {
		return rvf.kind
	}
	return reflect.TypeOf(rvf.value).Kind()
}

func (rvf *RVField) IsPointer() bool {
	if rvf.IsNil() {
		return true
	}
	return rvf.Kind() == reflect.Pointer
}

func (rvf *RVField) IsStruct() bool {
	if rvf.IsNil() {
		return false
	}
	return rvf.Kind() == reflect.Struct
}
