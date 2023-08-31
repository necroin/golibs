package csv

type Options struct {
	Delimiter        rune
	Comment          rune
	FieldsPerRecord  int
	LazyQuotes       bool
	TrimLeadingSpace bool
	Tag              string
}
