# Go Libraries
- [CSV](#CSV) - Reads csv files uses tags.
- [Concurrent](#Concurrent) - Provides thread-safe containers and atomic types.
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
## Concurrent
Provides thread-safe containers and atomic types.
- `AtomicBool`
	- Methods:
		- `Get() bool`
		- `Set(value bool)`
		- `Equal(other *AtomicBool) bool`
		- `NotEqual(other *AtomicBool) bool`
- `AtomicNumber[T]` 
	- Types:
		- `float32` | `float64` 
		- `int` | `int8` | `int16` | `int32` | `int64` 
		- `uint` | `uint8` | `uint16` | `uint32` | `uint64`
	- Methods:
		- `Get() T`
		- `Set(value T)`
		- `Add(value T)`
		- `Sub(value T)`
		- `Inc()`
		- `Dec()`
- `ConcurrentMap[K,V]`
	- Functions:
		- `NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V]` - Constructs a new container.
	- Methods:
		- `Insert(key K, value V)` - Inserts element into the container, replace if the container already contain an element with an equivalent key.
		- `Find(key K) (V, bool)` - Finds an element with key equivalent to key.
		- `Erase(key K) (V, bool)` - Removes specified element from the container.
		- `Iterate(handler func(K, V))` - Iterates over elements of the container with specified handler.
- `ConcurrentSlice[V]` and `ConcurrentSliceIterator[V]`
	- Functions:
		- `NewConcurrentSlice[V any]() *ConcurrentSlice[V]` - Constructs a new container.
	- Methods:
		- `Insert(index uint, value V) error` - Inserts element at the specified location in the container.
		- `Append(values ...V)` - Appends the given elements value to the end of the container.
		- `At(index uint) (V, error) ` - Returns the element at specified location index, with bounds checking.
										 If index is not within the range of the container, an error is returned.
		- `Erase(index uint) error` - Erases the specified element from the container.
		- `Size() int` - Returns the number of elements in the container.
		- `IsEmpty() bool` - Checks if the container has no elements.
		- `Front() V` - Returns the first element in the container.
						Calling front on an empty container causes undefined behavior.
		- `Back() V` - Returns the last element in the container.
					   Calling back on an empty container causes undefined behavior.
		- `Begin() *ConcurrentSliceIterator[V]` - Returns an iterator to the first element of the container.
		- `End() *ConcurrentSliceIterator[V]` - Returns an iterator to the element following the last element of the container.
	- `ConcurrentSliceIterator[V]` Methods:
		- `Next() *ConcurrentSliceIterator[V]`
		- `Get() (V, error)`
		- `Pos() uint`
		- `Set(value V) error `
		- `Equal(other *ConcurrentSliceIterator[V]) bool`