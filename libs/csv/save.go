package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"

	"github.com/necroin/golibs/utils"
)

func Marshal[T any](dataWriter io.Writer, data []T) error {
	return MarshalWithOptions(dataWriter, data, Options{})
}

func MarshalWithOptions[T any](dataWriter io.Writer, data []T, options Options) error {
	options.SetDefaults()

	writer := csv.NewWriter(dataWriter)

	if options.Delimiter != 0 {
		writer.Comma = options.Delimiter
	}
	writer.UseCRLF = options.UseCRLF

	headers, err := findHeaders(options.AdapterFunc(reflect.Indirect(reflect.ValueOf(utils.InstantiateSliceElement(&data)))), options)
	if err != nil {
		return err
	}
	writer.Write(headers)
	writer.Flush()
	for _, row := range data {
		rowPointer := &row
		record, err := buildRecord(options.AdapterFunc(reflect.Indirect(reflect.ValueOf(rowPointer))), options)
		if err != nil {
			return err
		}
		writer.Write(record)
		writer.Flush()
	}
	return nil
}

func findHeaders(value Adapter, options Options) ([]string, error) {
	result := []string{}

	if value.IsPointer() {
		value = value.Deref()
	}

	for fieldIndex := 0; fieldIndex < value.NumField(); fieldIndex++ {
		field := value.Field(fieldIndex)
		if field.IsStruct() {
			record, err := findHeaders(field, options)
			if err != nil {
				return result, err
			}
			result = append(result, record...)
		}

		tag := field.GetTag(options.Tag)
		if tag == "" || tag == "-" {
			continue
		}
		result = append(result, tag)
	}

	return result, nil
}

func buildRecord(value Adapter, options Options) ([]string, error) {
	result := []string{}

	if value.IsPointer() {
		value = value.Deref()
	}

	for fieldIndex := 0; fieldIndex < value.NumField(); fieldIndex++ {
		field := value.Field(fieldIndex)

		if field.IsStruct() {
			record, err := buildRecord(field, options)
			if err != nil {
				return result, err
			}
			result = append(result, record...)
		}

		tag := field.GetTag(options.Tag)
		if tag == "" || tag == "-" {
			continue
		}

		if field.IsPointer() {
			if field.IsNil() {
				result = append(result, "")
				continue
			}
			field = field.Deref()
		}
		result = append(result, fmt.Sprintf("%v", field.Get()))
	}
	return result, nil
}
