package csv_tests

import (
	"bytes"
	"testing"

	"github.com/necroin/golibs/libs/csv"
	"github.com/necroin/golibs/utils"
)

func TestSave_Common(t *testing.T) {
	file := &bytes.Buffer{}

	data := []CommonRow{
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

	if err := csv.Marshal(file, data); err != nil {
		t.Fatal(err)
	}

	SaveAssert(t, file.String(), string(CommonData))
}

func TestSave_Pointer(t *testing.T) {
	file := &bytes.Buffer{}

	data := []PointerRow{
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

	if err := csv.Marshal(file, data); err != nil {
		t.Fatal(err)
	}

	SaveAssert(t, file.String(), string(CommonData))
}

func TestSave_Pointer_Nil(t *testing.T) {
	file := &bytes.Buffer{}

	data := []PointerRow{
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

	if err := csv.Marshal(file, data); err != nil {
		t.Fatal(err)
	}

	SaveAssert(t, file.String(), string(PointerNilData))
}

func TestSave_Nested(t *testing.T) {
	file := &bytes.Buffer{}

	data := []NestedRow{
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

	if err := csv.Marshal(file, data); err != nil {
		t.Fatal(err)
	}

	SaveAssert(t, file.String(), string(CommonData))
}

func TestSave_Typed(t *testing.T) {
	file := &bytes.Buffer{}

	data := []TypedRow{
		{
			IntValue:    1,
			UintValue:   1,
			FloatValue:  1.1,
			StringValue: "value1",
		},
	}

	if err := csv.Marshal(file, data); err != nil {
		t.Fatal(err)
	}

	SaveAssert(t, file.String(), string(TypedData))
}
