package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"golibs/csv/utils"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func UnmarshalData[T any](data []byte, result *[]T) error {
	return UnmarshalDataWithOptions[T](data, result, Options{})
}

func Unmarshal[T any](dataReader io.Reader, result *[]T) error {
	return UnmarshalWithOptions[T](dataReader, result, Options{})
}

func UnmarshalDataWithOptions[T any](data []byte, result *[]T, options Options) error {
	return UnmarshalWithOptions[T](bytes.NewReader(data), result, options)
}

func UnmarshalWithOptions[T any](dataReader io.Reader, result *[]T, options Options) error {
	reader := csv.NewReader(dataReader)

	if options.Delimiter != 0 {
		reader.Comma = options.Delimiter
	}
	reader.Comment = options.Comment
	reader.FieldsPerRecord = options.FieldsPerRecord
	reader.LazyQuotes = options.LazyQuotes
	reader.TrimLeadingSpace = options.TrimLeadingSpace

	columnsList, err := reader.Read()
	if err != nil {
		return fmt.Errorf("[CSV] [Error] failed read columns: %s", err)
	}

	columns := map[string]int{}
	for index, column := range columnsList {
		columns[column] = index
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("[CSV] [Error] failed read data: %s", err)
		}
		if err := AddRecord(result, record, columns); err != nil {
			return err
		}
	}

	return nil
}

func AddRecord[T any](result *[]T, data []string, columns map[string]int) error {
	record := utils.InstantiateSliceElement(result)
	rvRecord := reflect.Indirect(reflect.ValueOf(record))

	if utils.IsStruct(rvRecord) {
		if err := fillStruct(rvRecord, data, columns); err != nil {
			return err
		}
	}

	*result = append(*result, *record)
	return nil
}

func fillStruct(rValue reflect.Value, data []string, columns map[string]int) error {
	if rValue.Type().Kind() == reflect.Pointer {
		rValue.Set(reflect.New(rValue.Type().Elem()))
		rValue = rValue.Elem()
	}

	for fieldIndex := 0; fieldIndex < rValue.NumField(); fieldIndex++ {
		rvField := rValue.Field(fieldIndex)
		rtField := rValue.Type().Field(fieldIndex)

		if utils.IsStruct(rvField) {
			if err := fillStruct(rValue.Field(fieldIndex), data, columns); err != nil {
				return err
			}
			continue
		}

		tag := rtField.Tag.Get("csv")
		if tag == "" || tag == "-" {
			continue
		}

		columnIndex, ok := columns[tag]
		if !ok {
			continue
		}

		if columnIndex >= len(data) {
			continue
		}

		if err := setValue(rValue.Field(fieldIndex), data[columnIndex]); err != nil {
			return err
		}
	}
	return nil
}

func setValue(field reflect.Value, data string) error {
	data = strings.TrimSpace(data)
	if data == "" {
		return nil
	}

	if field.Kind() == reflect.Pointer {
		field.Set(reflect.New(field.Type().Elem()))
		field = field.Elem()
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(data)
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		intData, err := strconv.ParseInt(data, 10, 64)
		if err != nil && data != "" {
			return fmt.Errorf("[CSV] [Error] failed parse int '%s': %s", data, err)
		}
		field.SetInt(intData)
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		uintData, err := strconv.ParseUint(data, 10, 64)
		if err != nil && data != "" {
			return fmt.Errorf("[CSV] [Error] failed parse uint '%s': %s", data, err)
		}
		field.SetUint(uintData)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		floatData, err := strconv.ParseFloat(data, 64)
		if err != nil && data != "" {
			return fmt.Errorf("[CSV] [Error] failed parse float '%s': %s", data, err)
		}
		field.SetFloat(floatData)
	case reflect.Bool:
		boolData, err := strconv.ParseBool(data)
		if err != nil && data != "" {
			return fmt.Errorf("[CSV] [Error] failed parse bool '%s': %s", data, err)
		}
		field.SetBool(boolData)
	}

	return nil
}
