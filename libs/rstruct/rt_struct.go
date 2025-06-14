package rstruct

import (
	"bytes"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/necroin/golibs/utils"
)

const (
	CommonFlatMode = iota
	NestedFlatMode
)

const (
	ZeroDefaultValueMode = iota
	NilDefaultValueMode
)

func GetDefaultValue(mode int, valueType reflect.Type) any {
	if mode == NilDefaultValueMode {
		return nil
	}
	return reflect.Zero(valueType).Interface()
}

type ExtendOption struct {
	// Value of extend type.
	Value any
	// Tag conversion map.
	Tags map[string]string
	// Tags prefix map.
	TagsPrefix map[string]string
	// Use true to remove all tags modifiers.
	IsPureTag bool
	// Uses '.' by default.
	PrefixDelimiter rune
	// Makes the nested structure flat.
	IsFlat bool
	// Use NestedFlatMode to preserve nesting in tags.
	FlatMode int
	// ZeroDefaultValueMode to fill values with zero type values, NilDefaultValueMode to fill all fields with nil.
	DefaultValueMode int
	// List of types that will be ignored in nested logic.
	IgnoreNested []any
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
			rtField:  tField,
			value:    tField.defaultValue,
			isStruct: false,
		}

		nestedRTStruct, ok := tField.defaultValue.(*RTStruct)
		if ok {
			vField.value = nestedRTStruct.New()
			vField.isStruct = true
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

func (rts *RTStruct) SetField(field *RTField) error {
	rtsField, ok := rts.fieldsByName[field.name]
	if !ok {
		return rts.AddField(field)
	}

	rtsField.defaultValue = field.defaultValue
	rtsField.tags = map[string]string{}

	return nil
}

func (rts *RTStruct) NumField() int {
	return len(rts.fields)
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

		ignoreNestedNames := []string{}
		for _, ignoreType := range extendOption.IgnoreNested {
			rIgnoreType := reflect.TypeOf(ignoreType)
			ignoreNestedNames = append(ignoreNestedNames, rIgnoreType.PkgPath()+"/"+rIgnoreType.Name())
		}

		rvExValue := reflect.ValueOf(extendOption.Value)
		rtExValue := reflect.TypeOf(extendOption.Value)

		for i := 0; i < rvExValue.NumField(); i++ {
			rvExField := rvExValue.Field(i)
			rtExField := rtExValue.Field(i)

			if !rtExField.IsExported() {
				continue
			}

			if utils.IsPointer(rvExField) {
				rvExField = reflect.New(rtExField.Type.Elem())
				rvExField = rvExField.Elem()
			}

			nestedExtendOption := ExtendOption{
				Value:            rvExField.Interface(),
				Tags:             extendOption.Tags,
				TagsPrefix:       utils.MapCopy(extendOption.TagsPrefix),
				IsPureTag:        extendOption.IsPureTag,
				PrefixDelimiter:  extendOption.PrefixDelimiter,
				IsFlat:           extendOption.IsFlat,
				FlatMode:         extendOption.FlatMode,
				DefaultValueMode: extendOption.DefaultValueMode,
				IgnoreNested:     extendOption.IgnoreNested,
			}

			inIgnoreNestedList := slices.Contains(ignoreNestedNames, rvExField.Type().PkgPath()+"/"+rvExField.Type().Name())

			if extendOption.IsFlat && utils.IsStruct(rvExField) && !inIgnoreNestedList {
				if extendOption.FlatMode == NestedFlatMode {
					for exTagName, rtfTagName := range nestedExtendOption.Tags {
						tag := rtExField.Tag.Get(exTagName)
						if tag == "" || tag == "-" {
							continue
						}

						if extendOption.IsPureTag {
							tag = utils.CleanTag(tag)
						}

						prefix, ok := extendOption.TagsPrefix[rtfTagName]
						if !ok {
							nestedExtendOption.TagsPrefix[rtfTagName] = tag
						} else {
							nestedExtendOption.TagsPrefix[rtfTagName] = prefix + string(extendOption.PrefixDelimiter) + tag
						}
					}
				}

				if err := rts.Extend(nestedExtendOption); err != nil {
					return fmt.Errorf("[RTStruct] [Extend] failed flat extend: %s", err)
				}
				continue
			}

			rtsField := rts.FieldByName(rtExField.Name)
			if rtsField == nil {
				if utils.IsStruct(rvExField) && !inIgnoreNestedList {
					nestedStruct := NewStruct()
					nestedStruct.Extend(nestedExtendOption)
					rtsField = NewRTField(rtExField.Name, nestedStruct)
				} else {
					rtsField = NewRTField(rtExField.Name, GetDefaultValue(extendOption.DefaultValueMode, rtExField.Type))
				}

				if err := rts.AddField(rtsField); err != nil {
					return fmt.Errorf("[RTStruct] [Extend] failed add field: %s", err)
				}
			}

			if utils.IsStruct(rvExField) && rtsField.IsStruct() && !inIgnoreNestedList {
				rtsField.AsStruct().Extend(nestedExtendOption)
			}

			for exTagName, rtfTagName := range extendOption.Tags {
				tag := rtExField.Tag.Get(exTagName)
				if tag == "" || tag == "-" {
					continue
				}

				if extendOption.IsPureTag {
					tag = utils.CleanTag(tag)
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

func (rts *RTStruct) write(writer *tabwriter.Writer, level int, isSorted bool) {
	levelOffset := strings.Repeat("\t", level)

	if level == 0 {
		fmt.Fprintf(writer, "{\n")
	} else {
		fmt.Fprintf(writer, "\t{\n")
	}

	fields := rts.fields
	if isSorted {
		fields = append(rts.fields[:0:0], rts.fields...)
		sort.Slice(fields, func(i, j int) bool { return fields[i].name < fields[j].name })
	}

	for _, field := range fields {
		nestedRTStruct, ok := field.defaultValue.(*RTStruct)
		if ok {
			fmt.Fprintf(writer, "%s %s", levelOffset, field.name)
			nestedRTStruct.write(writer, level+1, isSorted)
		} else {
			fmt.Fprintf(writer, "%s %s\t%v\n", levelOffset, field.name, field.defaultValue)
		}
		for tagName, tagContent := range field.tags {
			fmt.Fprintf(writer, "%s\t`%s\t: %s`\n", levelOffset, tagName, tagContent)
		}
	}
	fmt.Fprintf(writer, "%s}\n", levelOffset)
}

func (rts *RTStruct) string(isSorted bool) string {
	buffer := &bytes.Buffer{}
	writer := tabwriter.NewWriter(buffer, 1, 1, 1, ' ', 0)
	rts.write(writer, 0, isSorted)
	writer.Flush()
	return buffer.String()
}

func (rts *RTStruct) String() string {
	return rts.string(false)
}

func (rts *RTStruct) SortedString() string {
	return rts.string(true)
}
