# Go Libraries
- [CSV](#CSV)

## CSV
Reads csv files uses tags.
___
### Install
```sh
go get github.com/necroin/golibs/csv
```
___
### Load file
1. Use "csv" tag in struct.
- String read support
```Go
type CommonRow struct {
	FirstHeaderValue  string `csv:"Header1"`
	SecondHeaderValue string `csv:"Header2"`
	ThirdHeaderValue  string `csv:"Header3"`
}
```
- Pointer support (`nil` by default)
```Go
type PointerRow struct {
	FirstHeaderValue  *string `csv:"Header1"`
	SecondHeaderValue *string `csv:"Header2"`
	ThirdHeaderValue  *string `csv:"Header3"`
}
```
- Nested struct support
```Go
type NestedRowValue struct {
	SecondHeaderValue string `csv:"Header2"`
	ThirdHeaderValue  string `csv:"Header3"`
}

type NestedRow struct {
	FirstHeaderValue string `csv:"Header1"`
	NestedValue      NestedRowValue
}
```
- Standart types support
```Go
type TypedRow struct {
	IntValue    int     `csv:"Int"`
	UintValue   uint    `csv:"Uint"`
	FloatValue  float64 `csv:"Float"`
	StringValue string  `csv:"String"`
}
```
2. Unmarshal data
```Go
data, err := os.ReadFile("data.csv")
if err != nil {
	// error handle
}

rows := []CommonRow{}
if err := csv.UnmarshalData(data, &rows); err != nil {
	// error handle
}
```
___
### Save data to file
```Go
file, err := os.Create("data.csv")
if err != nil {
	// error handle
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
	// error handle
}
```
___