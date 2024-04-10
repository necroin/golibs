package rstruct

import (
	"fmt"
	"reflect"
)

type ExtendData struct {
	Value any
	Tags  []string
}

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
			rtField: tField,
			value:   tField.defaultValue,
		}
		vFields = append(vFields, vField)
		vFieldsByName[tField.name] = vField
	}

	return &RVStruct{
		rtStruct:     rts,
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

func (rts *RTStruct) Extend(extendValues ...ExtendData) error {
	for _, extendValue := range extendValues {
		rvValue := reflect.ValueOf(extendValue.Value)
		rtValue := reflect.TypeOf(extendValue.Value)

		for i := 0; i < rvValue.NumField(); i++ {
			field := rtValue.Field(i)
			newField := NewRTField(field.Name, reflect.Zero(field.Type).Interface())

			for _, tagName := range extendValue.Tags {
				tag := field.Tag.Get(tagName)
				if tag == "" || tag == "-" {
					continue
				}
				newField.SetTag(tagName, tag)
			}

			rts.AddField(newField)
		}
	}
	return nil
}
