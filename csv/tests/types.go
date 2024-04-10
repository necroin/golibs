package csv_tests

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
	SecondHeaderValue string `csv:"Header2"`
	ThirdHeaderValue  string `csv:"Header3"`
}

type NestedRow struct {
	FirstHeaderValue string `csv:"Header1"`
	NestedValue      NestedRowValue
}

type TypedRow struct {
	IntValue    int     `csv:"Int"`
	UintValue   uint    `csv:"Uint"`
	FloatValue  float64 `csv:"Float"`
	StringValue string  `csv:"String"`
}
