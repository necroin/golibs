package rstruct_tests

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/necroin/golibs/csv"
	"github.com/necroin/golibs/rstruct"
)

const (
	dataPath = "../../csv/assets"
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
	jsonData, _ := instance.Encode("json")
	if string(jsonData) != `{"age":"10","name":"\"MyName\""}` {
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
	jsonData, _ := instance.Encode("json")
	if string(jsonData) != `{"first_field":"\"\"","second_field":"0","third_field":"false"}` {
		t.Fatalf("invalid json result: %s", string(jsonData))
	}
}

func TestCSV(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.AddFields(
		rstruct.NewRTField("FirstHeaderValue", nil).
			SetTag("csv", "Header1"),
		rstruct.NewRTField("SecondHeaderValue", "").
			SetTag("csv", "Header2"),
		rstruct.NewRTField("ThirdHeaderValue", "").
			SetTag("csv", "Header3"),
	)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path.Join(dataPath, "common.csv"))
	if err != nil {
		t.Error(err)
	}

	rows := []rstruct.RVStruct{}
	if err := csv.UnmarshalDataWithOptions(data, &rows, csv.Options{
		AdapterFunc: func(value reflect.Value) csv.Adapter {
			return rstruct.NewCSVAdapter(customStruct, value)
		},
	}); err != nil {
		t.Error(err)
	}
	for _, row := range rows {
		output, _ := row.Encode("csv")
		fmt.Println(string(output))
	}
}
