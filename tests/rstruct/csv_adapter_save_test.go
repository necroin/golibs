package rstruct_tests

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/necroin/golibs/libs/csv"
	"github.com/necroin/golibs/libs/rstruct"
	csv_tests "github.com/necroin/golibs/tests/csv"
)

func TestCSVSave_Common(t *testing.T) {
	customStruct := rstruct.NewStruct()
	err := customStruct.Extend(rstruct.ExtendOption{
		Value: csv_tests.CommonRow{},
		Tags:  map[string]string{"csv": "csv"},
	})
	if err != nil {
		t.Fatal(err)
	}

	data := []csv_tests.CommonRow{
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

	rows := []rstruct.RVStruct{}

	for _, dataRow := range data {
		row := customStruct.New()

		row.FieldByName("FirstHeaderValue").Set(dataRow.FirstHeaderValue)
		row.FieldByName("SecondHeaderValue").Set(dataRow.SecondHeaderValue)
		row.FieldByName("ThirdHeaderValue").Set(dataRow.ThirdHeaderValue)

		rows = append(rows, *row)
	}

	file := &bytes.Buffer{}

	if err := csv.MarshalWithOptions(file, rows, csv.Options{
		AdapterFunc: func(value reflect.Value) csv.Adapter {
			return rstruct.NewCSVAdapter(customStruct, value)
		},
	}); err != nil {
		t.Fatal(err)
	}

	csv_tests.SaveAssert(t, file.String(), string(csv_tests.CommonData))
}
