package document_test

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/necroin/golibs/libs/document"
	"github.com/necroin/golibs/utils"
	"gopkg.in/yaml.v3"
)

type SectionInterface interface {
	Section(path string, opts ...document.SectionOption) (*document.Section, error)
	Key(path string, opts ...document.SectionOption) (*document.Key, error)
}

func TestDocument_Yaml_Map_splited_path(t *testing.T) {
	paths := [][]string{
		{"section1", "key1"},
		{"section1", "key2"},
		{"section1", "section11", "key3"},
		{"section1", "section11", "key4"},
		{"section1", "section12", "key3"},
		{"section1", "section12", "key4"},
		{"section2", "key1"},
		{"section2", "key2"},
	}

	expected := []string{
		"s1k1",
		"s1k2",
		"s1s11k3",
		"s1s11k4",
		"s1s12k3",
		"s1s12k4",
		"s2k1",
		"s2k2",
	}

	document, err := document.NewMapDocument("example_map.yaml", func(reader io.Reader) document.Decoder { return yaml.NewDecoder(reader) })
	if err != nil {
		t.Fatalf("failed create document: %s", err)
	}

	for pathIndex, path := range paths {
		var iterator SectionInterface = document
		expectedValue := expected[pathIndex]
		for partIndex, part := range path {
			if partIndex == len(path)-1 {
				key, err := iterator.Key(part)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Printf("Compare %s (real) and %s (expected)\n", key.String(), expectedValue)
				if key.String() != expectedValue {
					t.Fatalf("values mismatch: %s (real) != %s (expected)", key.String(), expectedValue)
				}
				continue
			}
			fmt.Printf("Get %s section\n", part)
			iterator, err = iterator.Section(part)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestDocument_Yaml_Map_full_path(t *testing.T) {
	paths := []string{
		"section1.key1",
		"section1.key2",
		"section1.section11.key3",
		"section1.section11.key4",
		"section1.section12.key3",
		"section1.section12.key4",
		"section2.key1",
		"section2.key2",
	}

	expected := []string{
		"s1k1",
		"s1k2",
		"s1s11k3",
		"s1s11k4",
		"s1s12k3",
		"s1s12k4",
		"s2k1",
		"s2k2",
	}

	document, err := document.NewMapDocument("example_map.yaml", func(reader io.Reader) document.Decoder { return yaml.NewDecoder(reader) })
	if err != nil {
		t.Fatalf("failed create document: %s", err)
	}

	for pathIndex, path := range paths {
		expectedValue := expected[pathIndex]
		key, err := document.Key(path)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("Compare %s (real) and %s (expected)\n", key.String(), expectedValue)
		if key.String() != expectedValue {
			t.Fatalf("values mismatch: %s (real) != %s (expected)", key.String(), expectedValue)
		}
	}
}

func TestDocument_Yaml_List(t *testing.T) {
	paths := []string{
		"section1",
		"section2",
	}

	expected := []string{
		"",
		"",
	}

	document, err := document.NewMapDocument("example_slice.yaml", func(reader io.Reader) document.Decoder { return yaml.NewDecoder(reader) })
	if err != nil {
		t.Fatalf("failed create document: %s", err)
	}

	documentJson, _ := json.Marshal(document)
	utils.SaveToFile("document.json", documentJson)

	listSection, err := document.Section("listSection")
	if err != nil {
		t.Fatal(err)
	}

	for pathIndex, path := range paths {
		expectedValue := expected[pathIndex]
		section, err := listSection.Index(pathIndex)
		if err != nil {
			t.Fatal(err)
		}
		key, err := section.Key(path)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("Compare %s (real) and %s (expected)\n", key.String(), expectedValue)
		if key.String() != expectedValue {
			t.Fatalf("values mismatch: %s (real) != %s (expected)", key.String(), expectedValue)
		}
	}
}
