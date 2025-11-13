package rstruct

type RTField struct {
	Name         string            `json:"name"`
	DefaultValue any               `json:"default_value"`
	Tags         map[string]string `json:"tags"`
}

func NewRTField(name string, defaultValue any) *RTField {
	return &RTField{
		Name:         name,
		DefaultValue: defaultValue,
		Tags:         map[string]string{},
	}
}

func (rtf *RTField) SetTag(name string, value string) *RTField {
	rtf.Tags[name] = value
	return rtf
}

func (rtf *RTField) RemoveTag(name string) {
	delete(rtf.Tags, name)
}

func (rtf *RTField) GetTag(name string) (string, bool) {
	value, ok := rtf.Tags[name]
	return value, ok
}

func (rvf *RTField) IsStruct() bool {
	_, ok := rvf.DefaultValue.(*RTStruct)
	return ok
}

func (rvf *RTField) AsStruct() *RTStruct {
	return rvf.DefaultValue.(*RTStruct)
}
