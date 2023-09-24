package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"

	"github.com/necroin/golibs/csv/utils"
)

func Marshal[T any](dataWriter io.Writer, data []T) error {
	return MarshalWithOptions(dataWriter, data, Options{})
}

func MarshalWithOptions[T any](dataWriter io.Writer, data []T, options Options) error {
	writer := csv.NewWriter(dataWriter)

	if options.Delimiter != 0 {
		writer.Comma = options.Delimiter
	}
	writer.UseCRLF = options.UseCRLF

	if options.Tag == "" {
		options.Tag = "csv"
	}

	headers, err := findHeaders(reflect.ValueOf(utils.InstantiateSliceElement(&data)), options)
	if err != nil {
		return err
	}
	writer.Write(headers)
	writer.Flush()
	for _, row := range data {
		record, err := buildRecord(reflect.ValueOf(row), options)
		if err != nil {
			return err
		}
		writer.Write(record)
		writer.Flush()
	}
	return nil
}

func findHeaders(value reflect.Value, options Options) ([]string, error) {
	result := []string{}

	if value.Type().Kind() == reflect.Pointer {
		value = value.Elem()
	}

	for fieldIndex := 0; fieldIndex < value.NumField(); fieldIndex++ {
		rvfield := value.Field(fieldIndex)
		rtField := value.Type().Field(fieldIndex)
		if utils.IsStruct(rvfield) {
			record, err := findHeaders(rvfield, options)
			if err != nil {
				return result, err
			}
			result = append(result, record...)
		}

		tag := rtField.Tag.Get(options.Tag)
		if tag == "" || tag == "-" {
			continue
		}
		result = append(result, tag)
	}

	return result, nil
}

func buildRecord(value reflect.Value, options Options) ([]string, error) {
	result := []string{}

	if value.Type().Kind() == reflect.Pointer {
		value = value.Elem()
	}

	for fieldIndex := 0; fieldIndex < value.NumField(); fieldIndex++ {
		rvField := value.Field(fieldIndex)
		rtField := value.Type().Field(fieldIndex)

		if utils.IsStruct(rvField) {
			record, err := buildRecord(rvField, options)
			if err != nil {
				return result, err
			}
			result = append(result, record...)
		}

		tag := rtField.Tag.Get(options.Tag)
		if tag == "" || tag == "-" {
			continue
		}

		if rvField.Type().Kind() == reflect.Pointer {
			if rvField.IsNil() {
				result = append(result, "")
				continue
			}
			rvField = rvField.Elem()
		}
		result = append(result, fmt.Sprintf("%v", rvField.Interface()))
	}
	return result, nil
}
