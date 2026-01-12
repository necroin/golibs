package table

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

type Table struct {
	tabwriter       *tabwriter.Writer
	headerFormatter Formatter

	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}

func New(headers []string, opts ...Option) *Table {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	return &Table{
		tabwriter:       tabwriter.NewWriter(options.out, 1, 1, 2, ' ', 0),
		headerFormatter: options.headerFormatter,
		Headers:         headers,
		Rows:            [][]string{},
	}
}

func (table *Table) InsertRow(values ...any) *Table {
	row := []string{}
	for _, value := range values {
		row = append(row, fmt.Sprint(value))
	}

	table.Rows = append(table.Rows, row)

	return table
}

func (table *Table) printValues(formatter Formatter, values ...string) {
	// for _, value := range values {
	// 	table.tabwriter.Write([]byte(value))
	// 	table.tabwriter.Write([]byte("\t"))
	// }

	row := strings.Join(values, "\t")
	if formatter != nil {
		row = formatter(row)
	}

	row = row + "\n"

	table.tabwriter.Write([]byte(row))
}

func (table *Table) printHeaders() {
	table.printValues(table.headerFormatter, table.Headers...)
}

func (table *Table) Print() {
	table.printHeaders()
	for _, row := range table.Rows {
		table.printValues(nil, row...)
	}
	table.tabwriter.Flush()

}
