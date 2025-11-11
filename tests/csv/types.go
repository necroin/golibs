package csv_tests

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	CommonData                = []byte("Header1,Header2,Header3\nR1V1,R1V2,R1V3\nR2V1,R2V2,R2V3\nR3V1,R3V2,R3V3\n")
	PointerNilData            = []byte("Header1,Header2,Header3\n,R1V2,R1V3\nR2V1,,R2V3\nR3V1,R3V2,\n,,\n")
	TypedData                 = []byte("Int,Uint,Float,String\n1,1,1.1,value1\n")
	DoubledColumnData         = []byte("Header1,Header2,Header2\nR1V1,R1V2,R1V3\nR2V1,R2V2,R2V3\nR3V1,R3V2,R3V3\n")
	HeaderRedeclareCommonData = []byte("Header1,Header2\nR1V1,R1V2\n__header_redeclare__,\nHeader2,Header3\nR2V2,R2V3\n")
)

type CommonRow struct {
	FirstHeaderValue  string `csv:"Header1" json:"Header1"`
	SecondHeaderValue string `csv:"Header2" json:"Header2"`
	ThirdHeaderValue  string `csv:"Header3" json:"Header3"`
}

type PointerRow struct {
	FirstHeaderValue  *string `csv:"Header1" json:"Header1"`
	SecondHeaderValue *string `csv:"Header2" json:"Header2"`
	ThirdHeaderValue  *string `csv:"Header3" json:"Header3"`
}

type NestedRowValue struct {
	SecondHeaderValue string `csv:"Header2" json:"Header2"`
	ThirdHeaderValue  string `csv:"Header3" json:"Header3"`
}

type NestedRow struct {
	FirstHeaderValue string `csv:"Header1" json:"Header1"`
	NestedValue      NestedRowValue
}

type TypedRow struct {
	IntValue    int     `csv:"Int"`
	UintValue   uint    `csv:"Uint"`
	FloatValue  float64 `csv:"Float"`
	StringValue string  `csv:"String"`
}

func LoadAssert[M any, N any](t *testing.T, rows []M, expected []N) {
	if !cmp.Equal(rows, expected) {
		t.Fatal(rows)
	}
}

func SaveAssert(t *testing.T, fact string, expected string) {
	if fact != expected {
		t.Fatalf("%s (fact) != %s (expected)", fact, expected)
	}
}
