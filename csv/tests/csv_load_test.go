package csv_tests

import (
	"os"
	"path"
	"testing"

	"github.com/necroin/golibs/csv"
	"github.com/necroin/golibs/utils"

	"github.com/google/go-cmp/cmp"
)

const (
	dataPath = "../assets"
)

func LoadAssert[M any, N any](t *testing.T, rows []M, expected []N) {
	if cmp.Equal(rows, expected) == false {
		t.Error(rows)
	}
}

func TestLoad_Common(t *testing.T) {
	data, err := os.ReadFile(path.Join(dataPath, "common.csv"))
	if err != nil {
		t.Error(err)
	}

	rows := []CommonRow{}
	if err := csv.UnmarshalData(data, &rows); err != nil {
		t.Error(err)
	}

	expected := []CommonRow{
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

	LoadAssert(t, rows, expected)
}

func TestLoad_Pointer(t *testing.T) {
	data, err := os.ReadFile(path.Join(dataPath, "common.csv"))
	if err != nil {
		t.Error(err)
	}

	rows := []PointerRow{}
	if err := csv.UnmarshalData(data, &rows); err != nil {
		t.Error(err)
	}

	expected := []PointerRow{
		{
			FirstHeaderValue:  utils.PointerOf("R1V1"),
			SecondHeaderValue: utils.PointerOf("R1V2"),
			ThirdHeaderValue:  utils.PointerOf("R1V3"),
		},
		{
			FirstHeaderValue:  utils.PointerOf("R2V1"),
			SecondHeaderValue: utils.PointerOf("R2V2"),
			ThirdHeaderValue:  utils.PointerOf("R2V3"),
		},
		{
			FirstHeaderValue:  utils.PointerOf("R3V1"),
			SecondHeaderValue: utils.PointerOf("R3V2"),
			ThirdHeaderValue:  utils.PointerOf("R3V3"),
		},
	}

	LoadAssert(t, rows, expected)
}

func TestLoad_Pointer_Nil(t *testing.T) {
	data, err := os.ReadFile(path.Join(dataPath, "pointer_nil.csv"))
	if err != nil {
		t.Error(err)
	}

	rows := []PointerRow{}
	if err := csv.UnmarshalData(data, &rows); err != nil {
		t.Error(err)
	}

	expected := []PointerRow{
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

	LoadAssert(t, rows, expected)
}

func TestLoad_Nested(t *testing.T) {
	data, err := os.ReadFile(path.Join(dataPath, "common.csv"))
	if err != nil {
		t.Error(err)
	}

	rows := []NestedRow{}
	if err := csv.UnmarshalData(data, &rows); err != nil {
		t.Error(err)
	}

	expected := []NestedRow{
		{
			FirstHeaderValue: "R1V1",
			NestedValue: NestedRowValue{
				SecondHeaderValue: "R1V2",
				ThirdHeaderValue:  "R1V3",
			},
		},
		{
			FirstHeaderValue: "R2V1",
			NestedValue: NestedRowValue{
				SecondHeaderValue: "R2V2",
				ThirdHeaderValue:  "R2V3",
			},
		},
		{
			FirstHeaderValue: "R3V1",
			NestedValue: NestedRowValue{
				SecondHeaderValue: "R3V2",
				ThirdHeaderValue:  "R3V3",
			},
		},
	}

	LoadAssert(t, rows, expected)
}

func TestLoad_Typed(t *testing.T) {
	data, err := os.ReadFile(path.Join(dataPath, "typed.csv"))
	if err != nil {
		t.Error(err)
	}

	rows := []TypedRow{}
	if err := csv.UnmarshalData(data, &rows); err != nil {
		t.Error(err)
	}

	expected := []TypedRow{
		{
			IntValue:    1,
			UintValue:   1,
			FloatValue:  1.1,
			StringValue: "value1",
		},
	}

	LoadAssert(t, rows, expected)
}

func TestLoad_DoubledColumn(t *testing.T) {
	data, err := os.ReadFile(path.Join(dataPath, "doubled_column.csv"))
	if err != nil {
		t.Error(err)
	}

	rows := []CommonRow{}
	err = csv.UnmarshalData(data, &rows)

	if err == nil {
		t.Error("Must be error: multiple column definition")
	}

	if err.Error() != "[CSV] [Error] failed read columns, multiple column definition: Header2" {
		t.Error("Must be error: multiple column definition")
	}
}
