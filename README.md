# Go Libraries
- [CSV](#CSV) - Reads csv files uses tags.
- [Concurrent](#Concurrent) - Provides thread-safe containers and atomic types.
- [Metrics](#Metrics) - Provides thread-safe metrics.
- [FSM](#FSM) - Provides finite state machine logic.
- [Tokenizer](#Tokenizer) - Parses given text to tokens.
- [RStruct](#RStruct) - Provides interface for custom struct.

## CSV
Reads csv files uses tags.
### Install
```sh
go get github.com/necroin/golibs/libs/csv
```
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

## Concurrent
Provides thread-safe containers and atomic types.
### Install
```sh
go get github.com/necroin/golibs/libs/concurrent
```
### Types
- `AtomicValue[T]`
	- Functions:
		- `NewAtomicValue[T any]() *AtomicValue[T]`
	- Methods:
		- `Get() T`
		- `Set(value T)`
- `AtomicNumber[T]` 
	- Types:
		- `float32` | `float64` 
		- `int` | `int8` | `int16` | `int32` | `int64` 
		- `uint` | `uint8` | `uint16` | `uint32` | `uint64`
	- Functions:
		- `NewAtomicNumber[T Number]() *AtomicNumber[T]`
	- Methods:
		- `Get() T`
		- `Set(value T)`
		- `Add(value T) T`
		- `Sub(value T) T`
		- `Inc() T`
		- `Dec() T`
- `ConcurrentMap[K,V]`
	- Functions:
		- `NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V]` - Constructs a new container.
	- Methods:
		- `Insert(key K, value V)` - Inserts element into the container, replace if the container already contain an element with an equivalent key.
		- `Find(key K) (V, bool)` - Finds an element with key equivalent to key.
		- `Erase(key K) (V, bool)` - Removes specified element from the container.
		- `Iterate(handler func(K, V))` - Iterates over elements of the container with specified handler.
		- `Size() int` - Returns the number of elements in the container.
		- `IsEmpty() bool` - Checks if the container has no elements.
		- `Keys() []K` - Returns slice of map keys.
		- `Values() []V` - Returns slice of map values.
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

## Metrics
Provides thread-safe metrics.
- `Counter`
- `Gauge`
- `Label`
- `Histogram`
- `The above vectorized metrics with labels`
### Install
```sh
go get github.com/necroin/golibs/libs/metrics
```
### Types
- `Counter` and `CounterVector`
```Go
package main

var (
	counter       = metrics.NewCounter(metrics.CounterOpts{Name: "test_counter", Help: "Counter help information"})
	counterVector = metrics.NewCounterVector(
		metrics.CounterOpts{Name: "test_counter_vector", Help: "Counter vector help information"},
		"label1", "label2",
	)
)

func main() {
	counter.Inc()
	counter.Add(rand.Float64())
	counterVector.WithLabelValues("test11", "test12").Inc()
	counterVector.WithLabelValues("test11", "test12").Add(rand.Float64())

	fmt.Println(counter.Get())
	fmt.Println(counterVector.WithLabelValues("test11", "test12").Get())
}
```

- `Gauge` and `GaugeVector`
```Go
package main

var (
	gauge       = metrics.NewGauge(metrics.GaugeOpts{Name: "test_gauge", Help: "Gauge help information"})
	gaugeVector = metrics.NewGaugeVector(
		metrics.GaugeOpts{Name: "test_gauge_vector", Help: "Gauge vector help information"},
		"label1", "label2",
	)
)

func main() {
	gauge.Set(rand.Float64())
	gauge.Add(rand.Float64())
	gauge.Sub(rand.Float64())
	gauge.Inc()
	gauge.Dec()

	gaugeVector.WithLabelValues("test11", "test12").Set(rand.Float64())
	gaugeVector.WithLabelValues("test11", "test12").Add(rand.Float64())
	gaugeVector.WithLabelValues("test11", "test12").Sub(rand.Float64())
	gaugeVector.WithLabelValues("test11", "test12").Inc()
	gaugeVector.WithLabelValues("test11", "test12").Dec()

	fmt.Println(gauge.Get())
	fmt.Println(gaugeVector.WithLabelValues("test11", "test12").Get())
}
```

- `Label` and `LabelVector`
```Go
var (
	label       = metrics.NewLabel(metrics.LabelOpts{Name: "test_label", Help: "Label help information"})
	labelVector = metrics.NewLabelVector(
		metrics.LabelOpts{Name: "test_label_vector", Help: "Label vector help information"},
		"label1", "label2",
	)
)

func main() {
	label.Set(RandomString(10))
	labelVector.WithLabelValues("test11", "test12").Set(RandomString(10))

	fmt.Println(label.Get())
	fmt.Println(labelVector.WithLabelValues("test11", "test12").Get())
}
```

- `Histogram` and `HistogramVector`
```Go
package main

var (
	histogram = metrics.NewHistogram(metrics.HistogramOpts{
		Name: "test_histogram", Help: "Histogram help information",
		Buckets: metrics.Buckets{Start: 0, Range: 10, Count: 10},
	})
	histogramVector = metrics.NewHistogramVector(
		metrics.HistogramOpts{
			Name: "test_histogram_vector", Help: "Histogram vector help information",
			Buckets: metrics.Buckets{Start: 0, Range: 10, Count: 10},
		},
		"label1", "label2",
	)
)

func main() {
	histogram.Observe(rand.Float64() * 100)
	histogramVector.WithLabelValues("test11", "test12").Observe(rand.Float64() * 100)
}
```

### Metrics server
```Go
func main() {
	registry := metrics.NewRegistry()
	registry.Register(counter)
	registry.Register(counterVector)
	registry.Register(gauge)
	registry.Register(gaugeVector)
	registry.Register(histogram)
	registry.Register(histogramVector)

	http.Handle("/metrics", registry.Handler())
	http.Handle("/metrics/json", registry.JsonHandler())
	http.ListenAndServe("localhost:3301", nil)
}
```

## FSM
Provides finite state machine logic.
### Install
```sh
go get github.com/necroin/golibs/libs/fsm
```
Types:
- `FSM[Args]`
	- Functions:
		- `NewFSM[Args]` - Creates new FSM.
	- Methoods:
		- `AddState(action func()) *State[Args]`
		- `SetCurrentState(state *State[Args])`
		- `Handle(args Args)`
		- `Execute()`
- `State[Args]`
	- Methoods:
		- `AddTransition(handler func(Args) *State[Args])`

## Tokenizer
Parses given text to tokens.
### Install
```sh
go get github.com/necroin/golibs/libs/tokenizer
```
Types:
- `Tokenizer`
	- Functions:
		- `NewTokenizer(tokens ...*Token)` - Creates new Tokenizer.
	- Methoods:
		- `Find(text []byte) (*Token, error)`
		- `Parse(text []byte) ([]*Token, error)`
		- `SetIgnoreSpaces(value bool)`
		- `SetIgnoreTabs(value bool)`
- `Token`
	- Functions:
		- `NewToken(name string, pattern string)` - Creates new Token.
	- Methoods:
		- `String() string`
		- `Name() string`
		- `Value() string`
		- `ValueInt() (int, error)`

## RStruct
Provides interface for custom struct.
### Install
```sh
go get github.com/necroin/golibs/libs/rstruct
```
Types:
- `RTField` - Describes the component of the field type.
	- Functions:
		- `NewRTField(name string, defaultValue any)` - Creates new Reflect Type Field.
	- Methoods:
		- `Name() string` - Returns the field name.
		- `Tags() map[string]string` - Returns a (name) - (value) tag table.
		- `SetTag(name string, value string) *RTField` - Sets the tag value by name.
		- `RemoveTag(name string)` - Deletes a tag by name.
		- `GetTag(name string) (string, bool)` - Gets the tag value by name.
		- `IsStruct() bool` - Checks RTStruct is RTStruct.
		- `AsStruct() *RTStruct` - Casts RTField to RTStruct.
- `RVField` - Describes the component of the field value.
	- Methoods:
		- `Set(value any)` - Sets the value of the field.
		- `Get() any` - Gets the value of the field.
		- `Type() *RTField` - Returns the field type.
		- `String() string` - Returns the string view of the field.
		- `ToJson() ([]byte, error)` - Returns the json view of the field.
		- `IsNil() bool` - Checks the field for nil.
		- `Kind() reflect.Kind` - Returns the kind of field.
		- `IsPointer() bool` - Checks the field for a pointer.
		- `IsInterface() bool` - Checks the field for the interface.
		- `IsStruct() bool` - Checks the field for structure.
		- `IsSlice() bool` - Checks the field for a slice.
		- `IsMap() bool` - Checks the field for the map.
		- `AsStruct() *RVStruct` - Casts RVField to RVStruct.
- `RTStruct` - Describes the component of the structure type.
	- Functions:
		- `NewStruct()` - Creates new Reflect Type Structure.
	- Methoods:
		- `New() *RVStruct` - Creates a new value for the structure.
		- `AddField(field *RTField) error` - Adds a new field.
		- `AddFields(fields ...*RTField) error` - Adds new fields.
		- `NumField() int` - Returns the number of fields.
		- `FieldByIndex(index int) *RTField` - Returns the field by index.
		- `FieldByName(name string) *RTField` - Returns a field by name.
		- `Extend(extendOptions ...ExtendOption) error` - Extends the structure using fields from another structure.
		- `String() string` - Returns the string view of the structure.
		- `SortedString() string` - Returns the string view of the structure with sorted fields.
- `RVStruct` - Describes the component of the structure value.
	- Methoods:
		- `FieldByIndex(index int) *RVField` - Returns the field by index.
		- `FieldByName(name string) *RVField` - Returns a field by name.
		- `FieldsListByTag(tag string) []*RVField` - Returns a list of fields with the specified tag.
		- `FieldsMapByTag(tag string) map[string]*RVField` - Returns a map of fields with the specified tag, where the key is the tag value.
		- `Type() *RTStruct` - Returns the field type.
		- `String() string` - Returns the string view of the structure.
		- `ToMap(tag string) map[string]any` - Returns the (field name) - (value) table.
		- `ToJson(tag string)` - Returns the json view of the structure.