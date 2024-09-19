package csv

import (
	"reflect"

	"github.com/necroin/golibs/utils"
)

type ReflectAdapter struct {
	value     reflect.Value
	fieldType reflect.StructField
}

func NewReflectAdapter(value reflect.Value) Adapter {
	return &ReflectAdapter{
		value: value,
	}
}

func (ra *ReflectAdapter) Kind() reflect.Kind {
	return ra.value.Kind()
}

func (ra *ReflectAdapter) IsPointer() bool {
	return utils.IsPointer(ra.value)
}

func (ra *ReflectAdapter) IsNil() bool {
	return utils.IsNil(ra.value)
}

func (ra *ReflectAdapter) IsStruct() bool {
	return utils.IsStruct(ra.value)
}

func (ra *ReflectAdapter) Set(value any) {
	ra.value.Set(reflect.ValueOf(value))
}

func (ra *ReflectAdapter) Get() any {
	return ra.value.Interface()
}

func (ra *ReflectAdapter) New() Adapter {
	return &ReflectAdapter{value: reflect.New(ra.value.Type().Elem())}
}

func (ra *ReflectAdapter) Deref() Adapter {
	return &ReflectAdapter{
		value: ra.value.Elem(),
	}
}

func (ra *ReflectAdapter) Field(index int) Adapter {
	return &ReflectAdapter{value: ra.value.Field(index), fieldType: ra.value.Type().Field(index)}
}

func (ra *ReflectAdapter) NumField() int {
	return ra.value.NumField()
}

func (ra *ReflectAdapter) GetTag(key string) string {
	return ra.fieldType.Tag.Get(key)
}

func (ra *ReflectAdapter) SetString(value string) {
	ra.value.SetString(value)
}

func (ra *ReflectAdapter) SetInt(value int64) {
	ra.value.SetInt(value)
}

func (ra *ReflectAdapter) SetUint(value uint64) {
	ra.value.SetUint(value)
}

func (ra *ReflectAdapter) SetFloat(value float64) {
	ra.value.SetFloat(value)
}

func (ra *ReflectAdapter) SetBool(value bool) {
	ra.value.SetBool(value)
}
