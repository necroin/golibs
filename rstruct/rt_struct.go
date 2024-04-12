package rstruct

import (
	"fmt"
	"reflect"

	"github.com/necroin/golibs/utils"
)

const (
	CommonFlatMode = iota
	NestedFlatMode
)

type ExtendOption struct {
	Value           any
	Tags            map[string]string
	TagsPrefix      map[string]string
	PrefixDelimiter rune
	IsFlat          bool
	FlatMode        int
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

		nestedRTStruct, ok := tField.defaultValue.(*RTStruct)
		if ok {
			vField.value = nestedRTStruct.New()
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
		return fmt.Errorf("[RTStruct] [AddField] field with name '%s' already exists", field.name)
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

func (rts *RTStruct) FieldByIndex(index int) *RTField {
	if index < 0 || index >= len(rts.fields) {
		return nil
	}
	return rts.fields[index]
}

func (rts *RTStruct) FieldByName(name string) *RTField {
	field, ok := rts.fieldsByName[name]
	if !ok {
		return nil
	}
	return field
}

func (rts *RTStruct) Extend(extendOptions ...ExtendOption) error {
	for extendOptionNumber, extendOption := range extendOptions {
		if extendOption.Value == nil {
			return fmt.Errorf("[RTStruct] [Extend] %d extend value is nil", extendOptionNumber)
		}

		if !utils.IsStruct(reflect.ValueOf(extendOption.Value)) {
			return fmt.Errorf("[RTStruct] [Extend] %d extend value is not a struct", extendOptionNumber)
		}

		if extendOption.PrefixDelimiter == 0 {
			extendOption.PrefixDelimiter = '.'
		}

		rvExValue := reflect.ValueOf(extendOption.Value)
		rtExValue := reflect.TypeOf(extendOption.Value)

		for i := 0; i < rvExValue.NumField(); i++ {
			rvExField := rvExValue.Field(i)
			rtExField := rtExValue.Field(i)

			if extendOption.IsFlat && utils.IsStruct(rvExField) {
				flatExtendOption := ExtendOption{
					Value:      rvExField.Interface(),
					Tags:       extendOption.Tags,
					TagsPrefix: utils.MapCopy(extendOption.TagsPrefix),
					IsFlat:     extendOption.IsFlat,
					FlatMode:   extendOption.FlatMode,
				}

				if extendOption.FlatMode == NestedFlatMode {
					for exTagName, rtfTagName := range flatExtendOption.Tags {
						tag := rtExField.Tag.Get(exTagName)
						if tag == "" || tag == "-" {
							continue
						}
						prefix, ok := extendOption.TagsPrefix[rtfTagName]
						if !ok {
							flatExtendOption.TagsPrefix[rtfTagName] = tag
						} else {
							flatExtendOption.TagsPrefix[rtfTagName] = prefix + string(extendOption.PrefixDelimiter) + tag
						}
					}
				}

				if err := rts.Extend(flatExtendOption); err != nil {
					return fmt.Errorf("[RTStruct] [Extend] failed flat extend: %s", err)
				}
				continue
			}

			rtsField := rts.FieldByName(rtExField.Name)
			if rtsField == nil {
				if utils.IsStruct(rvExField) {
					nestedStruct := NewStruct()
					nestedStruct.Extend(ExtendOption{
						Value:      rvExField.Interface(),
						Tags:       extendOption.Tags,
						TagsPrefix: utils.MapCopy(extendOption.TagsPrefix),
						IsFlat:     extendOption.IsFlat,
						FlatMode:   extendOption.FlatMode,
					})
					rtsField = NewRTField(rtExField.Name, nestedStruct)
				} else {
					rtsField = NewRTField(rtExField.Name, reflect.Zero(rtExField.Type).Interface())
				}

				if err := rts.AddField(rtsField); err != nil {
					return fmt.Errorf("[RTStruct] [Extend] failed add field: %s", err)
				}
			}

			for exTagName, rtfTagName := range extendOption.Tags {
				tag := rtExField.Tag.Get(exTagName)
				if tag == "" || tag == "-" {
					continue
				}

				prefix, ok := extendOption.TagsPrefix[rtfTagName]
				if ok {
					tag = prefix + string(extendOption.PrefixDelimiter) + tag
				}
				rtsField.SetTag(rtfTagName, tag)
			}
		}
	}
	return nil
}
