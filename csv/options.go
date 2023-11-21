package csv

type Options struct {
	// Field delimiter (set to ',' by default)
	Delimiter rune
	// Comment, if not 0, is the comment character.
	// Lines beginning with the Comment character without preceding whitespace are ignored.
	// With leading whitespace the Comment character becomes part of the field, even if TrimLeadingSpace is true.
	// Comment must be a valid rune and must not be \r, \n, or the Unicode replacement character (0xFFFD).
	// It must also not be equal to Comma.
	Comment rune
	// FieldsPerRecord is the number of expected fields per record.
	// If FieldsPerRecord is positive, Read requires each record to have the given number of fields.
	// If FieldsPerRecord is 0, Read sets it to the number of fields in the first record, so that future records must have the same field count.
	// If FieldsPerRecord is negative, no check is made and records may have a variable number of fields.
	FieldsPerRecord int
	// If LazyQuotes is true, a quote may appear in an unquoted field and a non-doubled quote may appear in a quoted field.
	LazyQuotes bool
	// If TrimLeadingSpace is true, leading white space in a field is ignored.
	// This is done even if the field delimiter, Comma, is white space.
	TrimLeadingSpace bool
	// True to use \r\n as the line terminator
	UseCRLF bool
	// True for trim spaces while read.
	TrimSpace bool
	// True for trim spaces while read.
	TrimQuotes bool
	// Parse Tag
	Tag string
}
