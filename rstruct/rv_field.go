package rstruct

import (
	"encoding/json"
	"fmt"
)

type RVField struct {
	value any
}

func (rvf *RVField) Set(valua any) {
	rvf.value = valua
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
