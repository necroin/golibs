package rstruct_tests

import (
	"testing"

	"github.com/necroin/golibs/libs/rstruct"
)

func Test_Error_Overwrite_Value_Struct(t *testing.T) {
	type Message1 struct {
		Field1 string `json:"field1"`
	}

	type Message2 struct {
		Field1 struct {
			Field2 string `json:"field2"`
		} `json:"field1"`
	}

	customStruct := rstruct.NewStruct()

	if err := customStruct.Extend(rstruct.ExtendOption{
		Value:            Message1{},
		Tags:             map[string]string{"json": "json"},
		DefaultValueMode: rstruct.NilDefaultValueMode,
		IsPureTag:        true,
	}); err != nil {
		t.Fatalf("failed extend message 1: %s", err)
	}

	err := customStruct.Extend(rstruct.ExtendOption{
		Value:            Message2{},
		Tags:             map[string]string{"json": "json"},
		DefaultValueMode: rstruct.NilDefaultValueMode,
		IsPureTag:        true,
	})

	if err.Error() != "[RTStruct] [Extend] Field1 field overwrite from value to struct" {
		t.Fatalf("failed extend message 2: %s", err)
	}
}

func Test_Error_Overwrite_Struct_Value(t *testing.T) {
	type Message1 struct {
		Field1 string `json:"field1"`
	}

	type Message2 struct {
		Field1 struct {
			Field2 string `json:"field2"`
		} `json:"field1"`
	}

	customStruct := rstruct.NewStruct()

	if err := customStruct.Extend(rstruct.ExtendOption{
		Value:            Message2{},
		Tags:             map[string]string{"json": "json"},
		DefaultValueMode: rstruct.NilDefaultValueMode,
		IsPureTag:        true,
	}); err != nil {
		t.Fatalf("failed extend message 2: %s", err)
	}

	err := customStruct.Extend(rstruct.ExtendOption{
		Value:            Message1{},
		Tags:             map[string]string{"json": "json"},
		DefaultValueMode: rstruct.NilDefaultValueMode,
		IsPureTag:        true,
	})

	if err.Error() != "[RTStruct] [Extend] Field1 field overwrite from struct to value" {
		t.Fatalf("failed extend message 1: %s", err)
	}
}
