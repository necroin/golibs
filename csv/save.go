package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"

	"github.com/necroin/golibs/csv/utils"
)

func Marshal[T any](dataWriter io.Writer, data []T) error {
	writer := csv.NewWriter(dataWriter)
	headers, err := findHeaders(reflect.ValueOf(utils.InstantiateSliceElement(&data)))
	if err != nil {
		return err
	}
	writer.Write(headers)
	writer.Flush()
	for _, row := range data {
		record, err := buildRecord(reflect.ValueOf(row))
		if err != nil {
			return err
		}
		writer.Write(record)
		writer.Flush()
	}
	return nil
}

func findHeaders(value reflect.Value) ([]string, error) {
	result := []string{}

	if value.Type().Kind() == reflect.Pointer {
		value = value.Elem()
	}

	for fieldIndex := 0; fieldIndex < value.NumField(); fieldIndex++ {
		rvfield := value.Field(fieldIndex)
		rtField := value.Type().Field(fieldIndex)
		if utils.IsStruct(rvfield) {
			record, err := findHeaders(rvfield)
			if err != nil {
				return result, err
			}
			result = append(result, record...)
		}

		tag := rtField.Tag.Get("csv")
		if tag == "" || tag == "-" {
			continue
		}
		result = append(result, tag)
	}

	return result, nil
}

func buildRecord(value reflect.Value) ([]string, error) {
	result := []string{}

	if value.Type().Kind() == reflect.Pointer {
		value = value.Elem()
	}

	for fieldIndex := 0; fieldIndex < value.NumField(); fieldIndex++ {
		rvField := value.Field(fieldIndex)
		rtField := value.Type().Field(fieldIndex)

		if utils.IsStruct(rvField) {
			record, err := buildRecord(rvField)
			if err != nil {
				return result, err
			}
			result = append(result, record...)
		}

		tag := rtField.Tag.Get("csv")
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
