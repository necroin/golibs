package rstruct

type RTField struct {
	name         string
	defaultValue any
	tags         map[string]string
}

func NewRTField(name string, defaultValue any) *RTField {
	return &RTField{
		name:         name,
		defaultValue: defaultValue,
		tags:         map[string]string{},
	}
}

func (rtf *RTField) Name() string {
	return rtf.name
}

func (rtf *RTField) Tags() map[string]string {
	return rtf.tags
}

func (rtf *RTField) SetTag(name string, value string) *RTField {
	rtf.tags[name] = value
	return rtf
}

func (rtf *RTField) RemoveTag(name string) {
	delete(rtf.tags, name)
}

func (rtf *RTField) GetTag(name string) (string, bool) {
	value, ok := rtf.tags[name]
	return value, ok
}

func (rvf *RTField) IsStruct() bool {
	_, ok := rvf.defaultValue.(*RTStruct)
	return ok
}

func (rvf *RTField) AsStruct() *RTStruct {
	return rvf.defaultValue.(*RTStruct)
}
