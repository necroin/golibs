package csv

import "reflect"

type Adapter interface {
	Kind() reflect.Kind
	IsPointer() bool
	IsNil() bool
	IsStruct() bool
	Set(value any)
	Get() any
	New() Adapter
	Deref() Adapter
	Field(index int) Adapter
	NumField() int
	GetTag(key string) string
	SetString(value string)
	SetInt(value int64)
	SetUint(value uint64)
	SetFloat(value float64)
	SetBool(value bool)
}
