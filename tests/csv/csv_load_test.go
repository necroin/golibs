package csv_tests

import (
	"testing"

	"github.com/necroin/golibs/libs/csv"
	"github.com/necroin/golibs/utils"
)

func TestLoad_Common(t *testing.T) {
	rows := []CommonRow{}
	if err := csv.UnmarshalData(CommonData, &rows); err != nil {
		t.Fatal(err)
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

func TestLoad_Common_HeaderRedeclare(t *testing.T) {
	rows := []CommonRow{}
	if err := csv.UnmarshalDataWithOptions(HeaderRedeclareCommonData, &rows, csv.Options{HeadersRedeclarePattern: "__header_redeclare__"}); err != nil {
		t.Fatal(err)
	}

	expected := []CommonRow{
		{
			FirstHeaderValue:  "R1V1",
			SecondHeaderValue: "R1V2",
		},
		{
			SecondHeaderValue: "R2V2",
			ThirdHeaderValue:  "R2V3",
		},
	}

	LoadAssert(t, rows, expected)
}

func TestLoad_Pointer(t *testing.T) {
	rows := []PointerRow{}
	if err := csv.UnmarshalData(CommonData, &rows); err != nil {
		t.Fatal(err)
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
	rows := []PointerRow{}
	if err := csv.UnmarshalData(PointerNilData, &rows); err != nil {
		t.Fatal(err)
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
	rows := []NestedRow{}
	if err := csv.UnmarshalData(CommonData, &rows); err != nil {
		t.Fatal(err)
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
	rows := []TypedRow{}
	if err := csv.UnmarshalData(TypedData, &rows); err != nil {
		t.Fatal(err)
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
	rows := []CommonRow{}
	err := csv.UnmarshalData(DoubledColumnData, &rows)

	if err == nil {
		t.Fatal("Must be error: multiple column definition")
	}

	if err.Error() != "[CSV] [Error] failed read columns, multiple column definition: Header2" {
		t.Fatal("Must be error: multiple column definition")
	}
}

func TestLoad_Map(t *testing.T) {
	rows := []map[string]string{}
	if err := csv.UnmarshalData(CommonData, &rows); err != nil {
		t.Fatal(err)
	}

	expected := []map[string]string{
		{
			"Header1": "R1V1",
			"Header2": "R1V2",
			"Header3": "R1V3",
		},
		{
			"Header1": "R2V1",
			"Header2": "R2V2",
			"Header3": "R2V3",
		},
		{
			"Header1": "R3V1",
			"Header2": "R3V2",
			"Header3": "R3V3",
		},
	}

	LoadAssert(t, rows, expected)
}
