package findex

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CSVFile[K comparable, V any, R Row[V]] struct {
	file       *os.File
	data       map[K]int64
	idata      map[int64]int64
	rowsCount  int64
	rowHandler func([]string) K
}

func NewCSV[K comparable, V any, R Row[V]](path string, rowHandler func(data []string) K) (File[K, V], error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("[findex] [NewCSV] failed open file: %s", err)
	}

	return &CSVFile[K, V, R]{
		file:       file,
		data:       map[K]int64{},
		idata:      map[int64]int64{},
		rowHandler: rowHandler,
	}, nil
}

func (file *CSVFile[K, V, R]) Close() {
	file.file.Close()
}

func (file *CSVFile[K, V, R]) Index() error {
	file.data = map[K]int64{}
	file.idata = map[int64]int64{}

	if _, err := file.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("[findex] [CSVFile] [Index] failed reset file offset: %s", err)
	}

	rowsCount := int64(0)
	csvReader := csv.NewReader(file.file)
	for {
		offset := csvReader.InputOffset()
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("[findex] [CSVFile] [Index] failed read data: %s", err)
		}

		rowsCount = rowsCount + 1

		if file.rowHandler != nil {
			file.data[file.rowHandler(record)] = offset
		}

		file.idata[rowsCount-1] = offset
	}

	file.rowsCount = rowsCount
	return nil
}

func (file *CSVFile[K, V, R]) RowCount() int64 {
	return file.rowsCount
}

func (file *CSVFile[K, V, R]) FindOffset(offset int64) (*V, error) {
	file.file.Seek(int64(offset), io.SeekStart)

	data, err := csv.NewReader(file.file).Read()
	if err != nil {
		return nil, fmt.Errorf("[findex] [CSVFile] [FindOffset] failed read data: %s", err)
	}

	result := new(V)
	var iResult R = result
	if err := iResult.UnmarshalCsv(data); err != nil {
		return nil, fmt.Errorf("[findex] [CSVFile] [FindOffset] failed decode data: %s", err)
	}

	return result, nil
}

func (file *CSVFile[K, V, R]) FindKey(key K) (*V, error) {
	offset, ok := file.data[key]
	if !ok {
		return nil, fmt.Errorf("[findex] [CSVFile] [FindKey] no data for key %v", key)
	}

	result, err := file.FindOffset(offset)
	if err != nil {
		return nil, fmt.Errorf("[findex] [CSVFile] [FindKey] -> %s", err)
	}

	return result, nil
}

func (file *CSVFile[K, V, R]) FindIndex(index int64) (*V, error) {
	offset, ok := file.idata[index]
	if !ok {
		return nil, fmt.Errorf("[findex] [CSVFile] [FindIndex] no data for index %d", index)
	}

	result, err := file.FindOffset(offset)
	if err != nil {
		return nil, fmt.Errorf("[findex] [CSVFile] [FindIndex] -> %s", err)
	}

	return result, nil
}
