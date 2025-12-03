package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/necroin/golibs/utils"
)

func UnmarshalData[T any](data []byte, result *[]T) error {
	return UnmarshalDataWithOptions(data, result, Options{})
}

func Unmarshal[T any](dataReader io.Reader, result *[]T) error {
	return UnmarshalWithOptions(dataReader, result, Options{})
}

func UnmarshalDataWithOptions[T any](data []byte, result *[]T, options Options) error {
	return UnmarshalWithOptions(bytes.NewReader(data), result, options)
}

func MakeColumns(columnsList []string) (map[string]int, error) {
	columns := map[string]int{}
	for index, column := range columnsList {
		if _, ok := columns[column]; ok {
			return nil, fmt.Errorf("[CSV] [Error] failed read columns, multiple column definition: %s", column)
		}
		columns[column] = index
	}

	return columns, nil
}

func MakeColumnsFromReader(reader *csv.Reader) (map[string]int, error) {
	columnsList, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("[CSV] [Error] failed read columns: %s", err)
	}
	return MakeColumns(columnsList)
}

func UnmarshalWithOptions[T any](dataReader io.Reader, result *[]T, options Options) error {
	options.SetDefaults()

	reader := csv.NewReader(dataReader)

	if options.Delimiter != 0 {
		reader.Comma = options.Delimiter
	}
	reader.Comment = options.Comment
	reader.FieldsPerRecord = options.FieldsPerRecord
	reader.LazyQuotes = options.LazyQuotes
	reader.TrimLeadingSpace = options.TrimLeadingSpace

	columns, err := MakeColumnsFromReader(reader)
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("[CSV] [Error] failed read data: %s", err)
		}

		if options.HeadersRedeclarePattern != "" && len(record) > 0 && record[0] == options.HeadersRedeclarePattern {
			columns, err = MakeColumnsFromReader(reader)
			if err != nil {
				return fmt.Errorf("[CSV] [Error] failed redeclare headers: %s", err)
			}
			continue
		}

		if err := AddRecord(result, record, columns, options); err != nil {
			return err
		}
	}

	return nil
}

func AddRecord[T any](result *[]T, data []string, columns map[string]int, options Options) error {
	record := utils.InstantiateSliceElement(result)
	rvRecord := reflect.ValueOf(record)
	rvRecordIndirect := reflect.Indirect(rvRecord)
	adapter := options.AdapterFunc(rvRecordIndirect)

	if adapter.IsStruct() {
		if err := fillStruct(adapter, data, columns, options); err != nil {
			return err
		}
	}

	if adapter.Kind() == reflect.Map {
		mapRecord := reflect.MakeMap(rvRecordIndirect.Type()).Interface().(T)
		record = &mapRecord
		rvRecord = reflect.ValueOf(mapRecord)
		for column, index := range columns {
			rvRecord.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(data[index]))
		}
	}

	*result = append(*result, *record)
	return nil
}

func fillStruct(structValue Adapter, data []string, columns map[string]int, options Options) error {
	if structValue.IsPointer() {
		structValue.Set(structValue.New().Get())
		structValue = structValue.Deref()
	}

	for fieldIndex := 0; fieldIndex < structValue.NumField(); fieldIndex++ {
		field := structValue.Field(fieldIndex)

		if field.IsStruct() {
			if err := fillStruct(structValue.Field(fieldIndex), data, columns, options); err != nil {
				return err
			}
			continue
		}

		tag := field.GetTag(options.Tag)
		if tag == "" || tag == "-" {
			continue
		}
		tag = utils.CleanTag(tag)

		columnIndex, ok := columns[tag]
		if !ok {
			continue
		}

		if columnIndex >= len(data) {
			continue
		}

		if err := setValue(structValue.Field(fieldIndex), data[columnIndex], options); err != nil {
			return err
		}
	}
	return nil
}

func setValue(field Adapter, data string, options Options) error {
	if options.TrimSpace {
		data = strings.Trim(data, "\t ")
	}

	if data == "" {
		return nil
	}

	if options.TrimQuotes {
		data = strings.Trim(data, "\"")
	}

	if field.IsPointer() {
		field.Set(field.New().Get())
		field = field.Deref()
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
