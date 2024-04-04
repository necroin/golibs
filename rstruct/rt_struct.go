package rstruct

import (
	"fmt"
	"reflect"
)

type RTStruct struct {
	fields       []*RTField
	fieldsByName map[string]*RTField
}

func NewStruct() *RTStruct {
	return &RTStruct{
		fields:       []*RTField{},
		fieldsByName: map[string]*RTField{},
	}
}

func (rts *RTStruct) New() *RVStruct {
	vFields := []*RVField{}
	vFieldsByName := map[string]*RVField{}

	for _, tField := range rts.fields {
		vField := &RVField{
			value: tField.defaultValue,
		}
		vFields = append(vFields, vField)
		vFieldsByName[tField.name] = vField
	}

	return &RVStruct{
		fields:       vFields,
		fieldsByName: vFieldsByName,
	}
}

func (rts *RTStruct) AddField(field *RTField) error {
	_, ok := rts.fieldsByName[field.name]
	if ok {
		return fmt.Errorf("[RTStruct] field with name '%s' already exists", field.name)
	}

	rts.fields = append(rts.fields, field)
	rts.fieldsByName[field.name] = field

	return nil
}

func (rts *RTStruct) AddFields(fields ...*RTField) error {
	for _, field := range fields {
		if err := rts.AddField(field); err != nil {
			return err
		}
	}
	return nil
}

func (rts *RTStruct) Extend(values ...any) error {
	for _, value := range values {
		rvValue := reflect.ValueOf(value)
		rtValue := reflect.TypeOf(value)

		for i := 0; i < rvValue.NumField(); i++ {
			field := rtValue.Field(i)
			newField := NewRTField(field.Name, reflect.Zero(field.Type).Interface())
			rts.AddField(newField)

			// tags := string(field.Tag)
			// tagsList := strings.Split(tags, " ")
		}
	}
	return nil
}
