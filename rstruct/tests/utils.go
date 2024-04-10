package rstruct_tests

import (
	"encoding/json"
	"testing"

	"github.com/necroin/golibs/rstruct"
)

func CompareCsvResults[T any](t *testing.T, rows []rstruct.RVStruct, cmpResult []T) {
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		cmpRow := cmpResult[i]

		csvRowData, _ := row.ToJson("csv")
		jsonCmpRowData, _ := json.Marshal(cmpRow)
		if string(csvRowData) != string(jsonCmpRowData) {
			t.Fatalf("%s != %s", string(csvRowData), string(jsonCmpRowData))
		}
	}
}
