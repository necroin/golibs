package tests

import (
	"slices"
	"strconv"
	"testing"

	"github.com/necroin/golibs/libs/findex"
)

type Worker struct {
	Id         string
	FirstName  string
	SecondName string
	Mail       string
	Gmail      string
	Job        string
}

func (worker *Worker) UnmarshalCsv(data []string) error {
	worker.Id = data[0]
	worker.FirstName = data[1]
	worker.SecondName = data[2]
	worker.Mail = data[3]
	worker.Gmail = data[4]
	worker.Job = data[5]

	return nil
}

func (worker *Worker) MarshalCsv() []string {
	return []string{
		worker.Id,
		worker.FirstName,
		worker.SecondName,
		worker.Mail,
		worker.Gmail,
		worker.Job,
	}
}

func TestMain(t *testing.T) {
	indexedFile, err := findex.NewCSV[string, Worker](
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
		keyIndexedData, err := indexedFile.FindKey(index)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(keyIndexedData.MarshalCsv(), data) {
			t.Fatalf("%s != %s", keyIndexedData, data)
		}

		intIndex, _ := strconv.Atoi(index)
		indexedData, err := indexedFile.FindIndex(int64(intIndex) - 100)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(indexedData.MarshalCsv(), data) {
			t.Fatalf("%s != %s", indexedData, data)
		}
	}
}
