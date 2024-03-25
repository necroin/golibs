package findex

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

const (
	libName = "findex"
)

type File[K comparable] struct {
	file       *os.File
	data       map[K]int64
	rawHandler func([]string) K
}

func NewFile[K comparable](path string, rawHandler func(data []string) K) (*File[K], error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed open file: %s", libName, err)
	}

	return &File[K]{
		file:       file,
		data:       map[K]int64{},
		rawHandler: rawHandler,
	}, nil
}

func (file File[K]) Close() {
	file.file.Close()
}

func (file File[K]) Index() error {
	csvReader := csv.NewReader(file.file)
	for {
		offset := csvReader.InputOffset()
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("[%s] [Index] failed read data: %s", libName, err)
		}
		file.data[file.rawHandler(record)] = offset
	}
	return nil
}

func (file File[K]) Find(key K) ([]string, error) {
	offset, ok := file.data[key]
	if !ok {
		return nil, fmt.Errorf("[%s] [Find] no data for key", libName)
	}
	file.file.Seek(int64(offset), io.SeekStart)
	return csv.NewReader(file.file).Read()
}
