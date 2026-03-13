package table

import (
	"fmt"
	"strings"
)

type Table struct {
	options   *Options
	maxWidths []int
	Headers   []string   `json:"headers"`
	Rows      [][]string `json:"rows"`
}

func New(headers ...string) *Table {
	return &Table{
		options:   DefaultOptions(),
		maxWidths: make([]int, len(headers)),
		Headers:   headers,
		Rows:      [][]string{},
	}
}

func (table *Table) SetOptions(opts ...Option) *Table {
	for _, opt := range opts {
		opt(table.options)
	}
	return table
}

func (table *Table) InsertRow(values ...any) *Table {
	row := []string{}
	for i, value := range values {
		valueStr := fmt.Sprint(value)
		valueLen := len(valueStr)
		row = append(row, valueStr)
		if valueLen > table.maxWidths[i] {
			table.maxWidths[i] = valueLen
		}
	}

	table.Rows = append(table.Rows, row)

	return table
}

func (table *Table) printValues(formatter Formatter, values ...string) {
	row := ""

	for index := range table.Headers {
		value := ""
		if index < len(values) {
			value = values[index]
		}

		row += value

		valueLen := len(value)
		padding := table.maxWidths[index] - valueLen + table.options.padding - 1
		if index == len(table.Headers)-1 {
			padding = table.maxWidths[index] - valueLen
		}
		row += strings.Repeat(" ", padding)
		row += string(table.options.padchar)
	}

	if formatter != nil {
		row = formatter(row)
	}

	row = row + "\n"

	table.options.out.Write([]byte(row))
}

func (table *Table) printHeaders() {
	table.printValues(table.options.headerFormatter, table.Headers...)
}

func (table *Table) Print() {
	table.printHeaders()
	for _, row := range table.Rows {
		table.printValues(nil, row...)
	}
}
