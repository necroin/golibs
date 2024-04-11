package rstruct_tests

import (
	"encoding/json"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/necroin/golibs/csv"
	csv_tests "github.com/necroin/golibs/csv/tests"
	"github.com/necroin/golibs/rstruct"
	"github.com/necroin/golibs/utils"
)

const (
	csvDataPath = "../../csv/assets"
)

func LoadAssertFunc[T any](t *testing.T, rows []rstruct.RVStruct, cmpResult []T) {
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		cmpRow := cmpResult[i]

		csvRowData, _ := row.ToJson("csv")
		jsonCmpRowData, _ := json.Marshal(cmpRow)
		if string(csvRowData) != string(jsonCmpRowData) {
			t.Fatalf("%s != %s", string(csvRowData), string(jsonCmpRowData))
		}
	}
}

func TestCSVLoad_Common(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.Extend(rstruct.ExtendOption{
		Value: csv_tests.CommonRow{},
		Tags:  map[string]string{"csv": "csv"},
	})
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path.Join(csvDataPath, "common.csv"))
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

	cmpResult := []csv_tests.CommonRow{
		{
			FirstHeaderValue:  "R1V1",
			SecondHeaderValue: "R1V2",
			ThirdHeaderValue:  "R1V3",
		},
		{
			FirstHeaderValue:  "R2V1",
			SecondHeaderValue: "R2V2",
			ThirdHeaderValue:  "R2V3",
		},
		{
			FirstHeaderValue:  "R3V1",
			SecondHeaderValue: "R3V2",
			ThirdHeaderValue:  "R3V3",
		},
	}

	CompareCsvResults(t, rows, cmpResult)
}

func TestLoad_Pointer_Nil(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.Extend(rstruct.ExtendOption{
		Value: csv_tests.PointerRow{},
		Tags:  map[string]string{"csv": "csv"},
	})
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path.Join(csvDataPath, "pointer_nil.csv"))
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

	cmpResult := []csv_tests.PointerRow{
		{
			FirstHeaderValue:  nil,
			SecondHeaderValue: utils.PointerOf("R1V2"),
			ThirdHeaderValue:  utils.PointerOf("R1V3"),
		},
		{
			FirstHeaderValue:  utils.PointerOf("R2V1"),
			SecondHeaderValue: nil,
			ThirdHeaderValue:  utils.PointerOf("R2V3"),
		},
		{
			FirstHeaderValue:  utils.PointerOf("R3V1"),
			SecondHeaderValue: utils.PointerOf("R3V2"),
			ThirdHeaderValue:  nil,
		},
		{
			FirstHeaderValue:  nil,
			SecondHeaderValue: nil,
			ThirdHeaderValue:  nil,
		},
	}

	CompareCsvResults(t, rows, cmpResult)
}
