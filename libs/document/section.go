package document

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/necroin/golibs/utils"
)

type Section struct {
	Sections map[string]*Section
	Keys     map[string]*Key
}

func NewSection(data any) *Section {
	sectionsByName := map[string]*Section{}
	keys := map[string]*Key{}

	if data == nil {
		return &Section{
			Sections: sectionsByName,
			Keys:     keys,
		}
	}

	rvData := utils.DerefValueOf(data)

	if utils.IsMap(rvData) {
		dataIterator := rvData.MapRange()
		for dataIterator.Next() {
			dataKey := dataIterator.Key().String()
			dataValue := dataIterator.Value()
			if dataValue.Interface() == nil {
				keys[dataKey] = &Key{name: dataKey, value: nil}
				continue
			}
			dataValue = reflect.ValueOf(dataValue.Interface())
			if utils.IsMap(dataValue) || utils.IsSlice(dataValue) {
				subSection := NewSection(dataValue.Interface())
				sectionsByName[dataKey] = subSection
				continue
			}
			keys[dataKey] = &Key{name: dataKey, value: dataValue.Interface()}
		}
	}

	if utils.IsSlice(rvData) {
		for index := range rvData.Len() {
			dataKey := formatIndex(index)
			dataValue := rvData.Index(index)
			if dataValue.Interface() == nil {
				keys[dataKey] = &Key{name: dataKey, value: nil}
				continue
			}
			dataValue = reflect.ValueOf(dataValue.Interface())
			if utils.IsMap(dataValue) || utils.IsSlice(dataValue) {
				subSection := NewSection(dataValue.Interface())
				sectionsByName[dataKey] = subSection
				continue
			}
			keys[dataKey] = &Key{name: dataKey, value: dataValue.Interface()}
		}
	}

	return &Section{
		Sections: sectionsByName,
		Keys:     keys,
	}
}

func (section *Section) Section(path string, opts ...SectionOption) (*Section, error) {
	if path == "" {
		return section, nil
	}

	options := DefaultSectionOptions()
	options.Apply(opts...)

	result := section

	parts := strings.Split(path, options.delimiter)
	for partIndex, part := range parts {
		result = result.Sections[part]
		if result == nil {
			return nil, fmt.Errorf("section %s data not found", strings.Join(parts[:partIndex+1], options.delimiter))
		}
	}

	return result, nil
}

func (section *Section) Index(index int, opts ...SectionOption) (*Section, error) {
	return section.Section(formatIndex(index), opts...)
}

func (section *Section) Key(path string, opts ...SectionOption) (*Key, error) {
	options := DefaultSectionOptions()
	options.Apply(opts...)

	parts := strings.Split(path, options.delimiter)

	if len(parts) > 1 {
		lastPartIndex := len(parts) - 1

		nestedSection, err := section.Section(strings.Join(parts[:lastPartIndex], options.delimiter), opts...)
		if err != nil {
			return nil, err
		}

		return nestedSection.Key(parts[lastPartIndex], opts...)
	}

	result, ok := section.Keys[path]
	if !ok {
		return nil, fmt.Errorf("no key: %s", path)
	}
	return result, nil
}

func (section *Section) SectionsNames() []string {
	return utils.MapKeys(section.Sections)
}

func (section *Section) KeysNames() []string {
	return utils.MapKeys(section.Keys)
}

func (section *Section) SectionsCount() int {
	return len(section.Sections)
}

func (section *Section) KeysCount() int {
	return len(section.Keys)
}
