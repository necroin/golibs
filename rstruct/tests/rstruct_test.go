package rstruct_tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/necroin/golibs/rstruct"
	"github.com/necroin/golibs/utils"
)

func TestMain(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.AddFields(
		rstruct.NewRTField("StringField", "StringData").
			SetTag("json", "string_field"),
		rstruct.NewRTField("IntField", 10).
			SetTag("json", "int_field"),
		rstruct.NewRTField("PointerField", utils.PointerOf("PointerFieldData")).
			SetTag("json", "pointer_field"),
		rstruct.NewRTField("NilPointerField", nil).
			SetTag("json", "nil_pointer_field"),
		rstruct.NewRTField("StructField", struct{}{}).
			SetTag("json", "struct_field"),
	)
	if err != nil {
		t.Fatal(err)
	}

	instance := customStruct.New()

	if instance.FieldByName("StringField").Kind() != reflect.String {
		t.Fatalf("invalid string field kind: %s", instance.FieldByName("StringField").Kind())
	}

	if instance.FieldByName("IntField").Kind() != reflect.Int {
		t.Fatalf("invalid int field kind: %s", instance.FieldByName("IntField").Kind())
	}

	if instance.FieldByName("PointerField").Kind() != reflect.Pointer {
		t.Fatalf("invalid pointer field kind: %s", instance.FieldByName("PointerField").Kind())
	}

	if instance.FieldByName("NilPointerField").Kind() != reflect.Pointer {
		t.Fatalf("invalid nil pointer field kind: %s", instance.FieldByName("NilPointerField").Kind())
	}

	if instance.FieldByName("StructField").Kind() != reflect.Struct {
		t.Fatalf("invalid struct field kind: %s", instance.FieldByName("StructField").Kind())
	}

	if instance.String() != fmt.Sprintf("map[IntField:10 NilPointerField: PointerField:%v StringField:StringData StructField:{}]", instance.FieldByName("PointerField").Get()) {
		t.Fatalf("invalid string implementation: %s", instance.String())
	}

	jsonData, _ := instance.ToJson("json")
	if string(jsonData) != `{"int_field":10,"nil_pointer_field":null,"pointer_field":"PointerFieldData","string_field":"StringData","struct_field":{}}` {
		t.Fatalf("invalid json implementation: %s", string(jsonData))
	}
}

func TestExtend_Common(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.Extend(rstruct.ExtendOption{
		Value: CommonExtendStruct{},
		Tags:  map[string]string{"json": "json"},
	})
	if err != nil {
		t.Fatal(err)
	}

	instance := customStruct.New()

	if instance.String() != "map[FirstField: SecondField:0 ThirdField:false]" {
		t.Fatalf("invalid string result: %s", instance.String())
	}
	jsonData, _ := instance.ToJson("json")
	if string(jsonData) != `{"first_field":"","second_field":0,"third_field":false}` {
		t.Fatalf("invalid json result: %s", string(jsonData))
	}
}

func TestExtend_Common_Prefix(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.Extend(rstruct.ExtendOption{
		Value:      CommonExtendStruct{},
		Tags:       map[string]string{"json": "common"},
		TagsPrefix: map[string]string{"common": "common"},
	})
	if err != nil {
		t.Fatal(err)
	}

	instance := customStruct.New()

	if instance.String() != "map[FirstField: SecondField:0 ThirdField:false]" {
		t.Fatalf("invalid string result: %s", instance.String())
	}
	jsonData, _ := instance.ToJson("common")
	if string(jsonData) != `{"common.first_field":"","common.second_field":0,"common.third_field":false}` {
		t.Fatalf("invalid json result: %s", string(jsonData))
	}
}

func TestExtend_Pointer(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.Extend(rstruct.ExtendOption{
		Value: PointerExtendStruct{},
		Tags:  map[string]string{"json": "json"},
	})
	if err != nil {
		t.Fatal(err)
	}

	instance := customStruct.New()

	if instance.String() != "map[FirstField:<nil> SecondField:<nil> ThirdField:<nil>]" {
		t.Fatalf("invalid string result: %s", instance.String())
	}
	jsonData, _ := instance.ToJson("json")
	if string(jsonData) != `{"first_field":null,"second_field":null,"third_field":null}` {
		t.Fatalf("invalid json result: %s", string(jsonData))
	}
}

func TestExtend_Nested(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.Extend(rstruct.ExtendOption{
		Value: NestedExtendStruct{},
		Tags:  map[string]string{"json": "json"},
	})
	if err != nil {
		t.Fatal(err)
	}

	instance := customStruct.New()

	if instance.String() != "map[NestedFirstField:{} NestedSecondField:{0} NestedThirdField:{false} NotNestedField:]" {
		t.Fatalf("invalid string result: %s", instance.String())
	}
	jsonData, _ := instance.ToJson("json")
	if string(jsonData) != `{"nested_first_field":{"first_field":""},"nested_second_field":{"second_field":0},"nested_third_field":{"third_field":false},"not_nested_field":""}` {
		t.Fatalf("invalid json result: %s", string(jsonData))
	}
}

func TestExtend_Nested_Flat_CommonMode(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.Extend(rstruct.ExtendOption{
		Value:  NestedExtendStruct{},
		Tags:   map[string]string{"json": "json"},
		IsFlat: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	instance := customStruct.New()

	if instance.String() != "map[FirstField: NotNestedField: SecondField:0 ThirdField:false]" {
		t.Fatalf("invalid string result: %s", instance.String())
	}
	jsonData, _ := instance.ToJson("json")
	if string(jsonData) != `{"first_field":"","not_nested_field":"","second_field":0,"third_field":false}` {
		t.Fatalf("invalid json result: %s", string(jsonData))
	}
}

func TestExtend_Nested_Flat_NestedMode(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.Extend(rstruct.ExtendOption{
		Value:    NestedExtendStruct{},
		Tags:     map[string]string{"json": "json"},
		IsFlat:   true,
		FlatMode: rstruct.NestedFlatMode,
	})
	if err != nil {
		t.Fatal(err)
	}

	instance := customStruct.New()

	if instance.String() != "map[FirstField: NotNestedField: SecondField:0 ThirdField:false]" {
		t.Fatalf("invalid string result: %s", instance.String())
	}
	jsonData, _ := instance.ToJson("json")
	if string(jsonData) != `{"nested_first_field.first_field":"","nested_second_field.second_field":0,"nested_third_field.third_field":false,"not_nested_field":""}` {
		t.Fatalf("invalid json result: %s", string(jsonData))
	}
}
