package tests

import (
	"os"
	"path"
	"testing"

	"github.com/necroin/golibs/csv"
	"github.com/necroin/golibs/csv/utils"
)

const (
	outDataPath = "../assets/out"
)

func TestSave_Common(t *testing.T) {
	file, err := os.Create(path.Join(outDataPath, "common.csv"))
	if err != nil {
		t.Error(err)
	}

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
		t.Error(err)
	}
}

func TestSave_Pointer(t *testing.T) {
	file, err := os.Create(path.Join(outDataPath, "pointer.csv"))
	if err != nil {
		t.Error(err)
	}

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
		t.Error(err)
	}
}

func TestSave_Pointer_Nil(t *testing.T) {
	file, err := os.Create(path.Join(outDataPath, "pointer_nil.csv"))
	if err != nil {
		t.Error(err)
	}

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
		t.Error(err)
	}
}

func TestSave_Nested(t *testing.T) {
	file, err := os.Create(path.Join(outDataPath, "nested.csv"))
	if err != nil {
		t.Error(err)
	}

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
		t.Error(err)
	}
}

func TestSave_Typed(t *testing.T) {
	file, err := os.Create(path.Join(outDataPath, "typed.csv"))
	if err != nil {
		t.Error(err)
	}

	data := []TypedRow{
		{
			IntValue:    1,
			UintValue:   1,
			FloatValue:  1.1,
			StringValue: "value1",
		},
	}

	if err := csv.Marshal(file, data); err != nil {
		t.Error(err)
	}
}
