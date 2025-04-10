package document

import (
	"fmt"
	"io"
	"os"
)

type Document struct {
	Root *Section
}

func NewDocument[T any](filepath string, container T, newDecoder func(reader io.Reader) Decoder) (*Document, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed open file: %w", err)
	}

	if err := newDecoder(file).Decode(&container); err != nil {
		return nil, fmt.Errorf("failed decode file data: %w", err)
	}

	return &Document{
		Root: NewSection(container),
	}, nil
}

func NewMapDocument(filepath string, newDecoder func(reader io.Reader) Decoder) (*Document, error) {
	return NewDocument(filepath, map[string]any{}, newDecoder)
}

func (document *Document) Section(path string, opts ...SectionOption) (*Section, error) {
	return document.Root.Section(path, opts...)
}

func (document *Document) Key(path string, opts ...SectionOption) (*Key, error) {
	return document.Root.Key(path, opts...)
}
