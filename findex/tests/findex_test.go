package tests

import (
	"slices"
	"testing"

	"github.com/necroin/golibs/findex"
)

func TestMain(t *testing.T) {
	indexedFile, err := findex.NewCSVFile[string](
		"data.csv",
		func(data []string) string {
			return data[0]
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	indexedFile.Index()

	compareDataMap := map[string][]string{
		"100": {"100", "Marylou", "Judye", "Marylou.Judye@yopmail.com", "Marylou.Judye@gmail.com", "police officer"},
		"101": {"101", "Robbi", "Scammon", "Robbi.Scammon@yopmail.com", "Robbi.Scammon@gmail.com", "firefighter"},
		"123": {"123", "Steffane", "Atonsah", "Steffane.Atonsah@yopmail.com", "Steffane.Atonsah@gmail.com", "police officer"},
		"225": {"225", "Viviene", "Roche", "Viviene.Roche@yopmail.com", "Viviene.Roche@gmail.com", "developer"},
		"254": {"254", "Pollyanna", "Chapland", "Pollyanna.Chapland@yopmail.com", "Pollyanna.Chapland@gmail.com", "police officer"},
		"291": {"291", "Sissy", "Eachern", "Sissy.Eachern@yopmail.com", "Sissy.Eachern@gmail.com", "police officer"},
		"350": {"350", "Sashenka", "Ferino", "Sashenka.Ferino@yopmail.com", "Sashenka.Ferino@gmail.com", "worker"},
	}

	for index, data := range compareDataMap {
		indexedData, err := indexedFile.Find(index)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(indexedData, data) {
			t.Fatalf("%s != %s", indexedData, data)
		}
	}
}
