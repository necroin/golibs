package tests

import (
	"fmt"
	"testing"

	"github.com/necroin/golibs/findex"
)

func TestMain(t *testing.T) {
	indexedFile, err := findex.NewFile[string](
		"data.csv",
		func(data []string) string {
			return data[0]
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	indexedFile.Index()
	fmt.Println(indexedFile.Find("110"))
}
