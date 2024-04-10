package rstruct_tests

import (
	"testing"

	"github.com/necroin/golibs/rstruct"
)

func TestMain(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.AddFields(
		rstruct.NewRTField("Name", "MyName").
			SetTag("json", "name"),
		rstruct.NewRTField("Age", 10).
			SetTag("json", "age"),
	)
	if err != nil {
		t.Fatal(err)
	}

	instance := customStruct.New()
	if instance.String() != "map[Age:10 Name:MyName]" {
		t.Fatalf("invalid string implementation: %s", instance.String())
	}
	jsonData, _ := instance.ToJson("json")
	if string(jsonData) != `{"age":10,"name":"MyName"}` {
		t.Fatalf("invalid json implementation: %s", string(jsonData))
	}
}

func TestExtend(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.Extend(rstruct.ExtendData{
		Value: ExtendTestStruct{},
		Tags:  []string{"json"},
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
