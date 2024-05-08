package findex

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CSVFile[K comparable] struct {
	file       *os.File
	data       map[K]int64
	rowHandler func([]string) K
}

func NewCSV[K comparable](path string, rowHandler func(data []string) K) (File[K], error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("[findex] [NewCSV] failed open file: %s", err)
	}

	return &CSVFile[K]{
		file:       file,
		data:       map[K]int64{},
		rowHandler: rowHandler,
	}, nil
}

func NewCSVFromFile[K comparable](file *os.File, rowHandler func(data []string) K) (File[K], error) {
	if file == nil {
		return nil, fmt.Errorf("[findex] [NewCSVFromFile] file is nil")
	}

	return &CSVFile[K]{
		file:       file,
		data:       map[K]int64{},
		rowHandler: rowHandler,
	}, nil
}

func (file *CSVFile[K]) Close() {
	file.file.Close()
}

func (file *CSVFile[K]) Index() error {
	file.data = map[K]int64{}
	if _, err := file.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("[findex] [CSVFile] [Index] failed reset file offset: %s", err)
	}

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
		file.data[file.rowHandler(record)] = offset
	}
	return nil
}

func (file *CSVFile[K]) Find(key K) ([]string, error) {
	offset, ok := file.data[key]
	if !ok {
		return nil, fmt.Errorf("[findex] [CSVFile] [Find] no data for key %v", key)
	}
	file.file.Seek(int64(offset), io.SeekStart)
	return csv.NewReader(file.file).Read()
}
